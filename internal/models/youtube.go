package models

import "time"

type Playlist struct {
	Id            string               `json:"id"`
	Title         string               `json:"title"`
	Videos        []Video              `json:"videos"`
	Thumbnails    map[string]Thumbnail `json:"thumb"`
	ThumbnailPath string               `json:"thumb_path"`
}

type Video struct {
	Id            string               `json:"id"`
	Title         string               `json:"title"`
	Thumbnails    map[string]Thumbnail `json:"thumb"`
	ThumbnailPath string               `json:"thumb_path"`
}

type Channel struct {
	Id            string               `json:"id"`
	Title         string               `json:"title"`
	Thumbnails    map[string]Thumbnail `json:"thumb"`
	ThumbnailPath string               `json:"thumb_path"`
	LastUpdate    time.Time            `json:"last_update"`
}

type Thumbnail struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}
