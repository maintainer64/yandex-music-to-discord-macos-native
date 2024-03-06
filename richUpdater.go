package main

import (
	"time"

	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/cache"
	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/configurator"
	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/discordimageuploader"
	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/nowplayingmacos"
	log "github.com/sirupsen/logrus"

	"github.com/hugolgst/rich-go/client"
)

type RichUpdaterTrack struct {
	cacheManager *cache.FixedQueue
	searcher     *discordimageuploader.YaMusicTrackSearcher
	uploader     *discordimageuploader.DiscordImageUploader
	config       *configurator.DiscordConfig
	state        *RichUpdaterState
}

type RichUpdaterState struct {
	TimeListenTrack time.Time
	CurrentTrack    string
}

func (c *RichUpdaterState) ListenTrack(track string) *time.Time {
	if track == c.CurrentTrack {
		return &c.TimeListenTrack
	}
	c.CurrentTrack = track //nolint:wsl
	c.TimeListenTrack = time.Now()

	return &c.TimeListenTrack
}

func NewRichUpdaterTrack(config *configurator.DiscordConfig) *RichUpdaterTrack {
	err := client.Login(config.ClientID)
	if err != nil {
		panic(err)
	}
	state := &RichUpdaterState{ //nolint:wsl
		TimeListenTrack: time.Now(),
		CurrentTrack:    "",
	}
	return &RichUpdaterTrack{ //nolint:wsl
		cacheManager: cache.NewFixedQueue(5),
		searcher:     discordimageuploader.NewYaMusicTrackSearcher(),
		uploader:     discordimageuploader.NewDiscordImageUploader(config.Token, config.ClientID),
		config:       config,
		state:        state,
	}
}

func (c *RichUpdaterTrack) generateActivity(track *nowplayingmacos.TrackInfo, yaTrack *discordimageuploader.CacheItem) client.Activity {
	activity := client.Activity{}
	if track == nil || track.Album == "" { //nolint:wsl
		log.Debug("Track is not detected. No Album Found")
		return activity
	}
	activity.State = track.Artist + " — " + track.Album //nolint:wsl
	activity.Details = track.Title
	if c.state != nil && yaTrack != nil { //nolint:wsl
		activity.Timestamps = &client.Timestamps{
			Start: c.state.ListenTrack(yaTrack.Key),
		}
	}

	yaTrackURL := ""
	if yaTrack != nil && yaTrack.Track != nil && yaTrack.Track.ID > 0 && yaTrack.Track.AlbumID > 0 {
		yaTrackURL = discordimageuploader.GetTrackYandexMusicURI(yaTrack.Track.AlbumID, yaTrack.Track.ID)
	}
	if yaTrackURL != "" { //nolint:wsl
		activity.Buttons = []*client.Button{
			&client.Button{
				Label: "🎶Listen",
				Url:   yaTrackURL,
			},
		}
	}
	if yaTrack != nil && yaTrack.DiscordAttachment != nil { //nolint:wsl
		activity.LargeImage = yaTrack.DiscordAttachment.URL
		activity.LargeText = track.Title
	}

	return activity
}

func (c *RichUpdaterTrack) Execute() {
	track := nowplayingmacos.GetCurrentTrackInfo()
	yaTrack := c.getYandexTrack(track)
	activity := c.generateActivity(track, yaTrack)
	_ = client.SetActivity(activity)
	log.Debug(
		"Get current track info\n",
		"activity: ", activity,
	)
}

func (c *RichUpdaterTrack) getYandexTrack(track *nowplayingmacos.TrackInfo) *discordimageuploader.CacheItem {
	if track == nil {
		return nil
	}
	if len(track.Title) == 0 { //nolint:wsl
		return nil
	}
	cacheKey := track.Title + "_" + track.Artist + "_" + track.Album                         //nolint:wsl
	if trackInfo, ok := c.cacheManager.Get(cacheKey).(*discordimageuploader.CacheItem); ok { //nolint:wsl
		return trackInfo
	}
	yandexTrack := c.searcher.GetTrack(track.Title)          //nolint:wsl
	var discordImage *discordimageuploader.DiscordAttachment //nolint:wsl
	if yandexTrack != nil && yandexTrack.AlbumURI != "" {    //nolint:wsl
		discordImage, _ = c.uploader.UploadByURL(
			discordimageuploader.GetCoverImageByURI(yandexTrack.AlbumURI),
		)
	} else if yandexTrack != nil && yandexTrack.ArtistURI != "" {
		discordImage, _ = c.uploader.UploadByURL(
			discordimageuploader.GetCoverImageByURI(yandexTrack.ArtistURI),
		)
	}
	trackInfoPush := &discordimageuploader.CacheItem{ //nolint:wsl
		Key:               cacheKey,
		Track:             yandexTrack,
		DiscordAttachment: discordImage,
	}
	c.cacheManager.Push(trackInfoPush)
	return trackInfoPush //nolint:wsl
}

func (c *RichUpdaterTrack) ExecuteForever() {
	for {
		c.Execute()
		time.Sleep(time.Second * time.Duration(c.config.UpdateSeconds))
	}
}
