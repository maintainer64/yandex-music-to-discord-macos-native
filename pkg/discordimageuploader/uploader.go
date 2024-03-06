package discordimageuploader

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

func DownloadRawDataURL(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL) //nolint:gosec
	if err != nil {
		return nil, err //nolint:wrapcheck
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return imageData, nil
}

func (c *DiscordImageUploader) UploadByURL(imageURL string) (*DiscordAttachment, error) {
	imageData, err := DownloadRawDataURL(imageURL)
	if err != nil {
		log.Debug("Uploader. Failed to download image: ", err)
		return nil, err //nolint:wsl
	}
	base64String := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageData) //nolint:wsl

	// Создание новой Discord сессии
	dg, err := discordgo.New("Bot " + c.token)
	if err != nil {
		log.Debug("Uploader. Error create discord uploader session: ", err) //nolint:wrapcheck
		return nil, err                                                     //nolint:wsl,wrapcheck
	}
	// Загрузка изображения в качестве аватара
	user, err := dg.UserUpdate("", base64String)
	if err != nil {
		log.Debug("Ошибка загрузки изображения в качестве аватара:", err) //nolint:wrapcheck
		return nil, err                                                   //nolint:wsl,wrapcheck
	}

	// Отключение сессии Discord при завершении работы
	defer dg.Close()           //nolint:wsl
	return &DiscordAttachment{ //nolint:wsl
		ID:  user.Avatar,
		URL: fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s", c.clientID, user.Avatar),
	}, errors.New("error: image URL not found")
}
