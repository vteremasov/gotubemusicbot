package youtube

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/vteremasov/go-music-bot/id3"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"os/exec"

	"github.com/cavaliercoder/grab"
	"github.com/rylio/ytdl"
)

const youtubeURLRegex = `(?:youtube\.com\/\S*(?:(?:\/e(?:mbed))?\/|watch\/?\?(?:\S*?&?v\=))|youtu\.be\/)([a-zA-Z0-9_-]{6,11})`

func GetYoutubeID(url string) string {
	r := regexp.MustCompile(youtubeURLRegex)

	return r.FindAllStringSubmatch(url, -1)[0][1]
}

func DownloadMp3(info ytdl.VideoInfo, dest string) string {
	if dest == "" {
		dest = "."
	}

	var format = info.Formats.Best(ytdl.FormatAudioEncodingKey)[0]

	videoURL, _ := info.GetDownloadURL(format)

	randomFileName := getRandomFileName()
	tempMp4Location := dest + "/" + randomFileName + ".mp4"
	tempMp3Location := dest + "/" + randomFileName + ".mp3"

	client := grab.NewClient()
	req, _ := grab.NewRequest(tempMp4Location, videoURL.String())

	resp := client.Do(req)

	if err := resp.Err(); err != nil {
		log.Fatal("Download failed: ", err)
	}

	var cmd = exec.Command("ffmpeg", "-i", tempMp4Location, "-q:a", "0", "-map", "a", tempMp3Location)

	err := cmd.Start()
	if err != nil {
		log.Panic(err)
	}

	err = cmd.Wait()
	if err != nil {
		log.Panic(err)
	}

	_ = os.Remove(tempMp4Location)

	return tempMp3Location
}

func getRandomFileName() string {
	randBytes := make([]byte, 16)
	_, err := rand.Read(randBytes)

	if err != nil {
		log.Panic(err)
	}

	return hex.EncodeToString(randBytes)
}

func RenameMp3File(mp3Location string, metadata *id3.Metadata) string {
	const seperatorsRegex = `[ &_=+:]`
	separators := regexp.MustCompile(seperatorsRegex)

	filename := separators.ReplaceAllString(metadata.Title, " ")
	finalLocation := filepath.Dir(mp3Location) + "/" + filename + ".mp3"

	err := os.Rename(mp3Location, finalLocation)
	if err != nil {
		log.Fatal("Error while renaming the file: ", err)
	}

	return finalLocation
}
