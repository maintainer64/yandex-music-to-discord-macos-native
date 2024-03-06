package main

import (
	"time"

	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/cache"
	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/configurator"
	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/discordImageUploader"
	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/nowPlayingMacOs"
	log "github.com/sirupsen/logrus"

	"github.com/hugolgst/rich-go/client"
)

type RichUpdaterTrack struct {
	cacheManager *cache.FixedQueue
	searcher     *discordImageUploader.YaMusicTrackSearcher
	uploader     *discordImageUploader.DiscordImageUploader
	config       *configurator.DiscordConfig
	timeStarted  time.Time
}

func NewRichUpdaterTrack(config *configurator.DiscordConfig) *RichUpdaterTrack {
	err := client.Login(config.ClientID)
	if err != nil {
		panic(err)
	}
	return &RichUpdaterTrack{
		cacheManager: cache.NewFixedQueue(5),
		searcher:     discordImageUploader.NewYaMusicTrackSearcher(),
		uploader:     discordImageUploader.NewDiscordImageUploader(config.Token, config.ClientID),
		config:       config,
		timeStarted:  time.Now(),
	}
}

func (c *RichUpdaterTrack) Execute() {
	track := nowPlayingMacOs.GetCurrentTrackInfo()
	if track.Album == "" {
		log.Debug("Track is not detected. No Album Found")
		return
	}
	buttons := make([]*client.Button, 0)
	yaTrack := c.getYandexTrack(track)
	yaTrackUrl := ""
	if yaTrack.Track != nil && yaTrack.Track.Id > 0 && yaTrack.Track.AlbumId > 0 {
		yaTrackUrl = discordImageUploader.GetTrackYandexMusicUri(yaTrack.Track.AlbumId, yaTrack.Track.Id)
	}
	if yaTrackUrl != "" {
		buttons = append(buttons, &client.Button{
			Label: "Track",
			Url:   yaTrackUrl,
		})
	}
	LargeImage := ""
	if yaTrack.DiscordAttachment != nil {
		LargeImage = yaTrack.DiscordAttachment.URL
	}
	end := c.timeStarted.Add(time.Hour)
	_ = client.SetActivity(client.Activity{
		State:      track.Artist + " — " + track.Album,
		Details:    track.Title,
		LargeImage: LargeImage,
		LargeText:  track.Title,
		Timestamps: &client.Timestamps{
			Start: &c.timeStarted,
			End:   &end,
		},
		Buttons: buttons,
	})
	log.Debug(
		"Get current track info\n",
		"title: ", track.Title,
		"\nurl: ", yaTrackUrl,
		"\ncover: ", yaTrack.DiscordAttachment,
	)
}

func (c *RichUpdaterTrack) getYandexTrack(track *nowPlayingMacOs.TrackInfo) *discordImageUploader.CacheItem {
	if track == nil {
		return nil
	}
	if track.Title == "" {
		return nil
	}
	cacheKey := track.Title + "_" + track.Artist + "_" + track.Album
	if trackInfo, ok := c.cacheManager.Get(cacheKey).(*discordImageUploader.CacheItem); ok {
		return trackInfo
	}
	yandexTrack := c.searcher.GetTrack(track.Title)
	var discordImage *discordImageUploader.DiscordAttachment = nil
	if yandexTrack != nil && yandexTrack.AlbumUri != "" {
		discordImage, _ = c.uploader.UploadByURL(
			discordImageUploader.GetCoverImageByUri(yandexTrack.AlbumUri),
		)
	} else if yandexTrack != nil && yandexTrack.ArtistUri != "" {
		discordImage, _ = c.uploader.UploadByURL(
			discordImageUploader.GetCoverImageByUri(yandexTrack.ArtistUri),
		)
	}
	trackInfoPush := &discordImageUploader.CacheItem{
		Key:               cacheKey,
		Track:             yandexTrack,
		DiscordAttachment: discordImage,
	}
	c.cacheManager.Push(trackInfoPush)
	return trackInfoPush
}

func (c *RichUpdaterTrack) ExecuteForever() {
	for {
		c.Execute()
		time.Sleep(time.Second * time.Duration(c.config.UpdateSeconds))
	}
}
