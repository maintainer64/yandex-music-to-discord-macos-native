package discordImageUploader

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

type DiscordImageUploader struct {
	token    string
	clientID string
}

type DiscordAttachment struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

func NewDiscordImageUploader(token string, clientID string) *DiscordImageUploader {
	return &DiscordImageUploader{
		token:    token,
		clientID: clientID,
	}
}

func DownloadRawDataUrl(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return imageData, nil
}

func (c *DiscordImageUploader) UploadByURL(imageURL string) (*DiscordAttachment, error) {
	imageData, err := DownloadRawDataUrl(imageURL)
	if err != nil {
		log.Debug("Uploader. Failed to download image: ", err)
		return nil, err
	}
	base64String := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageData)

	// Создание новой Discord сессии
	dg, err := discordgo.New("Bot " + c.token)
	if err != nil {
		log.Debug("Uploader. Error create discord uploader session: ", err)
		return nil, err
	}
	// Загрузка изображения в качестве аватара
	user, err := dg.UserUpdate("", base64String)
	if err != nil {
		log.Debug("Ошибка загрузки изображения в качестве аватара:", err)
		return nil, err
	}

	// Отключение сессии Discord при завершении работы
	defer dg.Close()
	return &DiscordAttachment{
		ID:  user.Avatar,
		URL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s", c.clientID, user.Avatar),
	}, errors.New("error: image URL not found")
}
