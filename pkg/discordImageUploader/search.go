package discordImageUploader

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
	Id         int    `json:"id,omitempty"`
	Title      string `json:"title,omitempty"`
	DurationMs int    `json:"durationMs,omitempty"`
	AlbumId    int    `json:"albumId,omitempty"`
	AlbumTitle string `json:"albumTitle,omitempty"`
	AlbumUri   string `json:"albumUri,omitempty"`
	ArtistId   int    `json:"artistId,omitempty"`
	ArtistName string `json:"artistName,omitempty"`
	ArtistUri  string `json:"artistUri,omitempty"`
}

func NewYaMusicTrackSearcher() *YaMusicTrackSearcher {
	return &YaMusicTrackSearcher{
		client: yamusic.NewClient(
			yamusic.HTTPClient(circuit.NewHTTPClient(time.Second*5, 10, nil)),
		),
	}
}

func GetCoverImageByUri(uri string) string {
	if uri == "" {
		return uri
	}
	if !strings.HasPrefix(uri, "https") {
		uri = fmt.Sprintf("https://%s", uri)
	}
	return strings.Replace(uri, "%%", "400x400", 1)
}

func GetTrackYandexMusicUri(albumId int, trackId int) string {
	return fmt.Sprintf("https://music.yandex.ru/album/%d/track/%d", albumId, trackId)
}

func (c *YaMusicTrackSearcher) GetTrack(query string) *YaMusicTrack {
	tracks, resp, err := c.client.Search().Tracks(context.Background(), query, nil)
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	searchResult := YaMusicTrack{}
	if len(tracks.Result.Tracks.Results) > 0 {
		track := tracks.Result.Tracks.Results[0]
		searchResult.Id = track.ID
		searchResult.Title = track.Title
		searchResult.DurationMs = track.DurationMs
		if track.Artists != nil && len(track.Artists) > 0 {
			artist := track.Artists[0]
			searchResult.ArtistId = artist.ID
			searchResult.ArtistName = artist.Name
			searchResult.ArtistUri = artist.Cover.URI
		}
		if track.Albums != nil && len(track.Albums) > 0 {
			album := track.Albums[0]
			searchResult.AlbumId = album.ID
			searchResult.AlbumTitle = album.Title
			searchResult.AlbumUri = album.CoverURI
		}
	}
	return &searchResult
}
