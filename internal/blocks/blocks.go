package blocks

import (
	"fmt"

	"github.com/su55y/yt_feed/internal/models"
)

func PrintChannels(channels []models.Channel) []models.Line {
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

func printVideosLines(videos []models.Video) []models.Line {
	lines := []models.Line{{Text: "back"}}
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

func PrintVideos(playlist models.Playlist) models.Blocks {
	return models.Blocks{
		Lines:   printVideosLines(playlist.Videos),
		Message: fmt.Sprintf("last %d videos of %s playlist", len(playlist.Videos), playlist.Title),
	}
}

func getPlaylistsLines(playlists map[string]models.Playlist) []models.Line {
	lines := []models.Line{{Text: "back"}}
	for _, v := range playlists {
		lines = append(lines, models.Line{
			Text: v.Title,
			Data: v.Id,
			Icon: v.ThumbnailPath,
		})
	}

	return lines
}

func PrintPlaylists(playlists map[string]models.Playlist) models.Blocks {
	return models.Blocks{
		Lines:   getPlaylistsLines(playlists),
		Message: fmt.Sprintf("last %d playlists", len(playlists)),
	}
}
