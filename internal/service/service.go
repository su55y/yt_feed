package service

import (
	"context"
	"errors"
	"html"
	"log"
	"path/filepath"
	"strings"

	"github.com/su55y/yt_feed/internal/config"
	"github.com/su55y/yt_feed/internal/consts"
	"github.com/su55y/yt_feed/internal/models"
	"github.com/su55y/yt_feed/pkg/downloader"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Service struct {
	YT        *youtube.Service
	AppConfig *config.AppConfig
}

func New(ctx context.Context, conf *config.AppConfig) Service {
	yt, err := youtube.NewService(ctx, option.WithAPIKey(conf.API_KEY))
	if err != nil {
		log.Fatalf("Unable to create YouTube service: %s", err.Error())
	}
	return Service{
		YT:        yt,
		AppConfig: conf,
	}
}

// Get channels list request and download thumbnails for them
func (s *Service) GetChannels() ([]models.Channel, error) {
	call := s.YT.Channels.List([]string{"snippet"}).
		Id(strings.Join(s.AppConfig.Channels, ",")).
		MaxResults(50)

	res, err := call.Do()
	if err != nil {
		return nil, err
	}

	if res.Items == nil || len(res.Items) == 0 {
		return nil, errors.New("get channels list request failed")
	}

	channels := make([]models.Channel, 0)
	thumbnails := make(map[string]string, 0)

	for _, c := range res.Items {
		channelThumbnails := parseThumbnails(c.Snippet.Thumbnails)
		path := ""
		switch s.AppConfig.ThumbSize {
		case consts.SP_HIGH:
			path = s.getThumbnailsPath(c.Id, channelThumbnails[consts.SP_HIGH].URL)
			thumbnails[path] = channelThumbnails[consts.SP_HIGH].URL
		case consts.SP_MEDIUM:
			path = s.getThumbnailsPath(c.Id, channelThumbnails[consts.SP_MEDIUM].URL)
			thumbnails[path] = channelThumbnails[consts.SP_MEDIUM].URL
		default:
			path = s.getThumbnailsPath(c.Id, channelThumbnails[consts.SP_DEFAULT].URL)
			thumbnails[path] = channelThumbnails[consts.SP_DEFAULT].URL
		}
		channels = append(channels, models.Channel{
			Id:            c.Id,
			Title:         c.Snippet.Title,
			Thumbnails:    channelThumbnails,
			ThumbnailPath: path,
		})
	}

	if !s.AppConfig.ThumbOff {
		downloader.DownloadAll(thumbnails)
	}

	return channels, nil
}

func (s *Service) GetUploads(channelId string) ([]models.Video, error) {
	if updId, ok := s.getUploadsId(channelId); ok {
		res, err := s.getPlaylistVideos(updId)
		if err != nil {
			return nil, err
		}

		return s.parseVideos(res), nil
	}

	return nil, errors.New("can't get uploads it for channel " + channelId)
}

func (s *Service) GetVideos(playlistId string) ([]models.Video, error) {
	res, err := s.getPlaylistVideos(playlistId)
	if err != nil {
		return nil, err
	}

	return s.parseVideos(res), nil
}

func (s *Service) GetPlaylists(channelId string) ([]models.Playlist, error) {
	res, err := s.getPlaylists(channelId)
	if err != nil {
		return nil, err
	}

	return s.parsePlaylists(res), nil
}

func (s *Service) getPlaylists(channelId string) (*youtube.PlaylistListResponse, error) {
	call := s.YT.Playlists.List([]string{"snippet"}).
		ChannelId(channelId).
		MaxResults(50)
	return call.Do()
}

// return playlists slice by channel id
func (s *Service) parsePlaylists(res *youtube.PlaylistListResponse) []models.Playlist {
	thumbnails := make(map[string]string, 0)
	playlists := []models.Playlist{}
	for _, p := range res.Items {
		playlistThumbnails := parseThumbnails(p.Snippet.Thumbnails)
		path := ""
		switch s.AppConfig.ThumbSize {
		case consts.SP_HIGH:
			path = s.getThumbnailsPath(p.Id, playlistThumbnails[consts.SP_HIGH].URL)
			thumbnails[path] = playlistThumbnails[consts.SP_HIGH].URL
		case consts.SP_MEDIUM:
			path = s.getThumbnailsPath(p.Id, playlistThumbnails[consts.SP_MEDIUM].URL)
			thumbnails[path] = playlistThumbnails[consts.SP_MEDIUM].URL
		default:
			path = s.getThumbnailsPath(p.Id, playlistThumbnails[consts.SP_DEFAULT].URL)
			thumbnails[path] = playlistThumbnails[consts.SP_DEFAULT].URL
		}
		vidRes, err := s.getPlaylistVideos(p.Id)
		if err != nil {
			log.Printf("can't get videos for playlist %s", p.Id)
			continue
		}
		videos := s.parseVideos(vidRes)
		playlists = append(playlists, models.Playlist{
			Id:            p.Id,
			Title:         html.EscapeString(p.Snippet.Title),
			Videos:        videos,
			Thumbnails:    playlistThumbnails,
			ThumbnailPath: path,
		})
	}

	downloader.DownloadAll(thumbnails)
	return playlists
}

func (s *Service) parseVideos(res *youtube.PlaylistItemListResponse) []models.Video {
	videos := make([]models.Video, 0)
	thumbnails := make(map[string]string, 0)
	for _, v := range res.Items {
		videoThumbnails := parseThumbnails(v.Snippet.Thumbnails)
		path := ""
		switch s.AppConfig.ThumbSize {
		case consts.SP_HIGH:
			path = s.getThumbnailsPath(v.Id, videoThumbnails[consts.SP_HIGH].URL)
			thumbnails[path] = videoThumbnails[consts.SP_HIGH].URL
		case consts.SP_MEDIUM:
			path = s.getThumbnailsPath(v.Id, videoThumbnails[consts.SP_MEDIUM].URL)
			thumbnails[path] = videoThumbnails[consts.SP_MEDIUM].URL
		default:
			path = s.getThumbnailsPath(v.Id, videoThumbnails[consts.SP_DEFAULT].URL)
			if len(videoThumbnails[consts.SP_DEFAULT].URL) > 0 {
				thumbnails[path] = videoThumbnails[consts.SP_DEFAULT].URL
			}
		}
		videos = append(videos, models.Video{
			Id:            v.Snippet.ResourceId.VideoId,
			Title:         html.EscapeString(v.Snippet.Title),
			Thumbnails:    videoThumbnails,
			ThumbnailPath: path,
		})
	}

	if !s.AppConfig.ThumbOff {
		downloader.DownloadAll(thumbnails)
	}
	return videos
}

// returns latest playlist 50 videos
func (s *Service) getPlaylistVideos(playlistId string) (*youtube.PlaylistItemListResponse, error) {
	call := s.YT.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(playlistId).
		MaxResults(50)

	return call.Do()
}

// get channel uploads playlist id
func (s *Service) getUploadsId(channelId string) (string, bool) {
	call := s.YT.Channels.List([]string{"contentDetails"}).Id(channelId)
	res, err := call.Do()
	if err != nil {
		log.Printf("get channel content details error: %s\n", err.Error())
		return "", false
	}

	if len(res.Items) >= 1 {
		for _, item := range res.Items {
			if item.ContentDetails != nil &&
				item.ContentDetails.RelatedPlaylists != nil &&
				len(item.ContentDetails.RelatedPlaylists.Uploads) == 24 {
				return item.ContentDetails.RelatedPlaylists.Uploads, true
			}
		}
	}

	return "", false
}

// returns Thumbnails struct by *youtube.ThumbnailDetails
func parseThumbnails(t *youtube.ThumbnailDetails) map[string]models.Thumbnail {
	thumbnails := map[string]models.Thumbnail{
		consts.SP_HIGH:    {},
		consts.SP_MEDIUM:  {},
		consts.SP_DEFAULT: {},
	}

	if t != nil {
		if t.Default != nil {
			var d models.Thumbnail
			d.URL = t.Default.Url
			d.Height = int(t.Default.Height)
			d.Width = int(t.Default.Width)
			thumbnails[consts.SP_DEFAULT] = d
		}
		if t.Medium != nil {
			var m models.Thumbnail
			m.URL = t.Medium.Url
			m.Height = int(t.Medium.Height)
			m.Width = int(t.Medium.Width)
			thumbnails[consts.SP_MEDIUM] = m
		}
		if t.High != nil {
			var h models.Thumbnail
			h.URL = t.High.Url
			h.Height = int(t.High.Height)
			h.Width = int(t.High.Width)
			thumbnails[consts.SP_HIGH] = h
		}
	}

	return thumbnails
}

// Join path for thumbnails from cache path, id and extension, taken from url
func (s *Service) getThumbnailsPath(id, url string) string {
	ext := filepath.Ext(filepath.Base(url))
	if len(ext) == 0 && len(filepath.Base(url)) != 0 {
		ext = ".jpg"
	}
	sizePrefix := consts.SP_DEFAULT
	switch s.AppConfig.ThumbSize {
	case consts.SP_HIGH:
		sizePrefix = consts.SP_HIGH
	case consts.SP_MEDIUM:
		sizePrefix = consts.SP_MEDIUM

	}
	return filepath.Join(
		s.AppConfig.ThumbDir,
		sizePrefix+id+ext,
	)
}
