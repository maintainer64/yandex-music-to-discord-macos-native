package main

import (
	"encoding/json"
	"os"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/maintainer64/yandex-music-to-discord-macos-native/pkg/configurator"
	log "github.com/sirupsen/logrus"
)

func initConfig(config *configurator.ProjectConfig) {
	if config == nil {
		return
	}
	if config.Discord == nil { //nolint:wsl
		config.Discord = &configurator.DiscordConfig{}
	}
	if config.Debug { //nolint:wsl
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}
}

type ViewSettings struct {
	app        fyne.App
	mainWindow fyne.Window
	config     *configurator.ProjectConfig
}

func (c *ViewSettings) init() {
	c.app = app.New()
	c.config = configurator.LoadConfig("config.json")
	initConfig(c.config)
	errorLabel := widget.NewLabel("") //nolint:wsl
	debugCheckbox := widget.NewCheck("Debug", func(checked bool) {
		c.config.Debug = checked
	})
	debugCheckbox.Checked = c.config.Debug
	clientIDEntry := widget.NewEntry()
	clientIDEntry.SetPlaceHolder("Discord application id")
	clientIDEntry.Text = c.config.Discord.ClientID
	maskedTokenEntry := widget.NewEntry()
	maskedTokenEntry.SetPlaceHolder("Super secret token discord app")
	maskedTokenEntry.Text = c.config.Discord.Token
	updateSecondsSpinBox := widget.NewEntry()
	updateSecondsSpinBox.SetPlaceHolder("How match update time seconds")
	updateSecondsSpinBox.Text = strconv.FormatInt(c.config.Discord.UpdateSeconds, 10)

	saveButton := widget.NewButton("Save", func() {
		c.config.Discord.ClientID = clientIDEntry.Text
		c.config.Discord.Token = maskedTokenEntry.Text
		updateSeconds, err := strconv.ParseInt(updateSecondsSpinBox.Text, 10, 64)
		if err != nil {
			errorLabel.SetText("Invalid seconds")
			return
		}
		c.config.Discord.UpdateSeconds = updateSeconds

		// Сохранение конфигурационного файла
		configData, err := json.MarshalIndent(c.config, "", "  ")
		if err != nil {
			errorLabel.SetText("Save configuration error")
			return
		}
		err = os.WriteFile("config.json", configData, 0600)
		if err != nil {
			errorLabel.SetText("Save write configuration error")
			return
		}
		initConfig(c.config)
	})

	// Вкладки для отображения настроек и запуска программы
	tabs := container.NewAppTabs(
		container.NewTabItem("Start", container.NewVBox(
			widget.NewButton("Start", func() {
				c.onClickButtonStart()
			}),
		)),
		container.NewTabItem("Settings", container.NewVBox(
			widget.NewLabel("Discord application ID:"),
			clientIDEntry,
			widget.NewLabel("Discord token:"),
			maskedTokenEntry,
			widget.NewLabel("Update seconds:"),
			updateSecondsSpinBox,
			debugCheckbox,
			saveButton,
		)),
	)

	// Создание контейнера с элементами интерфейса
	content := container.NewVBox(
		tabs,
		errorLabel,
	)

	// Создание окна приложения
	c.mainWindow = c.app.NewWindow("Rich Yandex.Music")
	c.mainWindow.SetMaster()
	c.mainWindow.SetContent(content)
	c.mainWindow.Resize(fyne.NewSize(400, 300))
	c.mainWindow.CenterOnScreen()
	c.mainWindow.SetFixedSize(true)
	c.mainWindow.SetPadded(true)
	c.mainWindow.ShowAndRun()
}

func (c *ViewSettings) onClickButtonStart() {
	log.Debug("Click button start")
	c.mainWindow.Hide()
	updater := NewRichUpdaterTrack(c.config.Discord)
	updater.ExecuteForever()
}

func (c *ViewSettings) Execute() {
	c.init()
}

func main() {
	view := ViewSettings{app: nil, mainWindow: nil, config: nil}
	view.Execute()
}
