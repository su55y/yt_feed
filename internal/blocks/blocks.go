package blocks

import (
	"fmt"

	"github.com/su55y/yt_feed/internal/models"
)

func PrintChannelMenu(channelId string) []models.Line {
	actions := []string{
		"back", "videos", "playlists", "update videos", "update playlists",
	}
	lines := make([]models.Line, 0)
	for _, a := range actions {
		lines = append(lines, models.Line{
			Text: a,
			Data: channelId,
		})
	}

	return lines
}

func PrintChannels(channels map[string]models.Channel) []models.Line {
	lines := make([]models.Line, 0)
	for _, c := range channels {
		lines = append(lines, models.Line{
			Text: c.Title,
			Data: c.Id,
			Icon: c.ThumbnailPath,
		})
	}
	return lines
}

func PrintVideos(playlist models.Playlist, channelsId string) models.Blocks {
	return models.Blocks{
		Lines:   getVideosLines(playlist.Videos, channelsId),
		Message: fmt.Sprintf("last %d videos of %s playlist", len(playlist.Videos), playlist.Title),
	}
}

func PrintPlaylists(playlists map[string]models.Playlist, channelId string) models.Blocks {
	return models.Blocks{
		Lines:   getPlaylistsLines(playlists, channelId),
		Message: fmt.Sprintf("last %d playlists", len(playlists)),
	}
}

func getVideosLines(videos []models.Video, channelId string) []models.Line {
	lines := []models.Line{{Text: "back", Data: "channel:" + channelId}}
	for _, v := range videos {
		if v.Title != "Private video" {
			lines = append(lines, models.Line{
				Text: v.Title,
				Data: v.Id,
				Icon: v.ThumbnailPath,
			})
		}
	}

	return lines
}

func getPlaylistsLines(playlists map[string]models.Playlist, channelId string) []models.Line {
	lines := []models.Line{{Text: "back", Data: "channel:" + channelId}}
	for _, v := range playlists {
		lines = append(lines, models.Line{
			Text: v.Title,
			Data: v.Id,
			Icon: v.ThumbnailPath,
		})
	}

	return lines
}
