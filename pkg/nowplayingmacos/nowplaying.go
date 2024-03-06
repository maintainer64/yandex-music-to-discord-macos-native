package nowplayingmacos

import (
	"os/exec"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type TrackInfo struct {
	Artist       string
	Album        string
	Title        string
	Duration     float64
	ElapsedTime  float64
	PlaybackRate bool
}

func parseLineToString(line string) string {
	clearString := strings.TrimSpace(line)
	if strings.ToLower(clearString) == "null" {
		return ""
	}

	return clearString
}

func parseLineToFloat(line string) float64 {
	digit, err := strconv.ParseFloat(line, 64)
	if err != nil {
		return 0
	}

	return digit
}

func parseLineToBoolean(line string) bool {
	return line == "1"
}

func GetCurrentTrackInfo() *TrackInfo {
	cmd := exec.Command("nowplaying-cli", "get", "artist", "album", "title", "duration", "elapsedTime", "playbackRate")

	output, err := cmd.Output()
	if err != nil {
		log.Panic("Error:", err)
		return nil
	}

	lines := strings.Split(string(output), "\n")
	track := &TrackInfo{}

	track.Artist = parseLineToString(lines[0])
	track.Album = parseLineToString(lines[1])
	track.Title = parseLineToString(lines[2])
	track.Duration = parseLineToFloat(parseLineToString(lines[3]))
	track.ElapsedTime = parseLineToFloat(parseLineToString(lines[4]))
	track.PlaybackRate = parseLineToBoolean(parseLineToString(lines[5]))

	return track
}
