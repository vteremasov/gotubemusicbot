package id3

import (
	"github.com/bogem/id3v2"
	"github.com/rylio/ytdl"
	"log"
)

// Metadata contains information collected from user
type Metadata struct {
	Title  string
	Artist string
	Album  string
}

// CollectMetadata will prompt the user to input mp3 title, artist, and album
func CollectMetadata(videoInfo *ytdl.VideoInfo) *Metadata {
	mp3Title := videoInfo.Title
	mp3Artist := videoInfo.Author
	mp3Album := "Unknown" // TODO: find a way to set album

	return &Metadata{Title: mp3Title, Artist: mp3Artist, Album: mp3Album}
}

// SetMetadata will apply metadata to the mp3 file's id3 tag
func SetMetadata(mp3Location string, metadata *Metadata) {
	tag, err := id3v2.Open(mp3Location, id3v2.Options{Parse: true})

	if err != nil {
		log.Fatal("Error while initializing a tag: ", err)
	}

	tag.SetTitle(metadata.Title)
	tag.SetArtist(metadata.Artist)
	tag.SetAlbum(metadata.Album)

	if err = tag.Save(); err != nil {
		log.Fatal("Error while saving a tag: ", err)
	}
}
