package blocks

import "github.com/su55y/yt_feed/internal/models"

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

func PrintVideos(videos []models.Video) []models.Line {
	lines := make([]models.Line, 0)
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

func PrintPlaylists(playlists map[string]models.Playlist) []models.Line {
	lines := make([]models.Line, 0)
	for _, v := range playlists {
		lines = append(lines, models.Line{
			Text: v.Title,
			Data: v.Id,
			Icon: v.ThumbnailPath,
		})
	}

	return lines
}
