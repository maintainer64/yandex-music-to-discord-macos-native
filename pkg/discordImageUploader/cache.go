package discordImageUploader

type CacheItem struct {
	Key               string
	Track             *YaMusicTrack
	DiscordAttachment *DiscordAttachment
}

func (item *CacheItem) GetCacheKey() string {
	return item.Key
}
