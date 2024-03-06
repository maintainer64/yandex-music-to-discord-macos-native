package discordimageuploader

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ndrewnee/go-yamusic/yamusic"
	circuit "github.com/rubyist/circuitbreaker"
)

type YaMusicTrackSearcher struct {
	client *yamusic.Client
}

type YaMusicTrack struct {
	ID         int    `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	DurationMs int    `json:"durationMs,omitempty"`
	AlbumID    int    `json:"albumId,omitempty"`
	AlbumTitle string `json:"albumTitle,omitempty"`
	AlbumURI   string `json:"albumUri,omitempty"`
	ArtistID   int    `json:"artistId,omitempty"`
	ArtistName string `json:"artistName,omitempty"`
	ArtistURI  string `json:"artistUri,omitempty"`
}

func NewYaMusicTrackSearcher() *YaMusicTrackSearcher {
	return &YaMusicTrackSearcher{
		client: yamusic.NewClient(
			yamusic.HTTPClient(circuit.NewHTTPClient(time.Second*5, 10, nil)),
		),
	}
}

func GetCoverImageByURI(uri string) string {
	if uri == "" {
		return uri
	}
	if !strings.HasPrefix(uri, "https") { //nolint:wsl
		uri = fmt.Sprintf("https://%s", uri)
	}

	return strings.Replace(uri, "%%", "400x400", 1)
}

func GetTrackYandexMusicURI(albumID int, trackID int) string {
	return fmt.Sprintf("https://music.yandex.ru/album/%d/track/%d", albumID, trackID)
}

func (c *YaMusicTrackSearcher) GetTrack(query string) *YaMusicTrack {
	tracks, resp, err := c.client.Search().Tracks(context.Background(), query, nil) //nolint:bodyclose
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK { //nolint:wsl
		return nil
	}
	searchResult := YaMusicTrack{}             //nolint:wsl
	if len(tracks.Result.Tracks.Results) > 0 { //nolint:wsl
		track := tracks.Result.Tracks.Results[0]
		searchResult.ID = track.ID
		searchResult.Title = track.Title
		searchResult.DurationMs = track.DurationMs
		if track.Artists != nil && len(track.Artists) > 0 { //nolint:wsl
			artist := track.Artists[0]
			searchResult.ArtistID = artist.ID
			searchResult.ArtistName = artist.Name
			searchResult.ArtistURI = artist.Cover.URI
		}
		if track.Albums != nil && len(track.Albums) > 0 { //nolint:wsl
			album := track.Albums[0]
			searchResult.AlbumID = album.ID
			searchResult.AlbumTitle = album.Title
			searchResult.AlbumURI = album.CoverURI
		}
	}

	return &searchResult
}
