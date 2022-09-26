package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/su55y/yt_feed/internal/config"
	"github.com/su55y/yt_feed/internal/consts"
	"github.com/su55y/yt_feed/internal/models"
	"github.com/su55y/yt_feed/internal/service"
)

const (
	channelsFile = "channels.json"
)

type Storage struct {
	AppConfig *config.AppConfig
	Service   *service.Service
}

func New(conf *config.AppConfig, serv *service.Service) Storage {
	return Storage{
		AppConfig: conf,
		Service:   serv,
	}
}

func (s *Storage) ReadChannels() (map[string]models.Channel, error) {
	path := filepath.Join(s.AppConfig.CachePath, channelsFile)
	if !exists(path) {
		channels, err := s.Service.GetChannels()
		if err != nil {
			return nil, err
		}
		if !s.writeChannelsToFile(channels) {
			return nil, errors.New("can't write channels to file")
		}

		return channels, nil
	}

	channels := make(map[string]models.Channel, 0)
	channelsRaw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("read channels from file error: %s\n", err.Error())
		return nil, err
	}

	if err := json.Unmarshal([]byte(channelsRaw), &channels); err != nil {
		log.Printf("channels unmarshal error: %s\n", err.Error())
		return nil, err
	}
	return channels, nil
}

func (s *Storage) ReadAllPlaylists(
	channelId string,
	update bool,
) (map[string]models.Playlist, error) {
	path := filepath.Join(
		s.AppConfig.CachePath,
		fmt.Sprintf("%s%s%s", consts.P_PLAYLISTS, channelId, consts.EXT_JSON),
	)

	if !exists(path) || update {
		playlists, err := s.Service.GetPlaylists(channelId)
		if err != nil {
			return nil, err
		}

		if !s.writePlaylistsToFile(channelId, playlists) {
			return nil, errors.New("can't write playlists to file")
		}

		playlistsMap := make(map[string]models.Playlist, 0)
		for _, p := range playlists {
			playlistsMap[p.Id] = p
		}
		return playlistsMap, nil
	}

	playlists := make([]models.Playlist, 0)
	playlistsRaw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("read playlists from file error: %s\n", err.Error())
		return nil, err
	}

	if err := json.Unmarshal([]byte(playlistsRaw), &playlists); err != nil {
		log.Printf("playlists unmarshal error: %s\n", err.Error())
		return nil, err
	}

	playlistsMap := make(map[string]models.Playlist, 0)
	for _, p := range playlists {
		playlistsMap[p.Id] = p
	}

	return playlistsMap, nil
}

func (s *Storage) ReadUploads(channelId string, update bool) ([]models.Video, error) {
	path := filepath.Join(
		s.AppConfig.CachePath,
		fmt.Sprintf("%s%s%s", consts.P_VIDEOS, channelId, consts.EXT_JSON),
	)

	if !exists(path) || update {
		videos, err := s.Service.GetUploads(channelId)
		if err != nil {
			return nil, err
		}

		if !s.writeVideosToFile(channelId, videos) {
			return nil, errors.New("can't write videos to file")
		}

		return videos, nil
	}

	videos := make([]models.Video, 0)
	videosRaw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("read videos from file error: %s\n", err.Error())
		return nil, err
	}

	if err := json.Unmarshal([]byte(videosRaw), &videos); err != nil {
		log.Printf("videos unmarshal error: %s\n", err.Error())
		return nil, err
	}
	return videos, nil
}

// read playlist videos
func (s *Storage) ReadPlaylist(playlistId string) ([]models.Video, error) {
	path := filepath.Join(
		s.AppConfig.CachePath,
		fmt.Sprintf("%s%s%s", consts.P_VIDEOS, playlistId, consts.EXT_JSON),
	)
	if !exists(path) {
		videos, err := s.Service.GetVideos(playlistId)
		if err != nil {
			return nil, err
		}

		if !s.writeVideosToFile(playlistId, videos) {
			return nil, errors.New("can't write videos to file")
		}

		return videos, nil
	}

	videos := make([]models.Video, 0)
	videosRaw, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("read videos from file error: %s\n", err.Error())
		return nil, err
	}

	if err := json.Unmarshal([]byte(videosRaw), &videos); err != nil {
		log.Printf("videos unmarshal error: %s\n", err.Error())
		return nil, err
	}
	return videos, nil
}

func (s *Storage) writeChannelsToFile(channels map[string]models.Channel) bool {
	path := filepath.Join(s.AppConfig.CachePath, channelsFile)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("open channels.json file error: %s\n", err.Error())
		return false
	}

	if err := json.NewEncoder(file).Encode(&channels); err != nil {
		log.Printf("write to channels.json file error: %s\n", err.Error())
	}

	return true
}

func (s *Storage) writePlaylistsToFile(channelId string, playlists []models.Playlist) bool {
	path := filepath.Join(
		s.AppConfig.CachePath,
		fmt.Sprintf("%s%s%s", consts.P_PLAYLISTS, channelId, consts.EXT_JSON),
	)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("open %#v file error: %s\n", path, err.Error())
		return false
	}

	if err := json.NewEncoder(file).Encode(&playlists); err != nil {
		log.Printf("write to %#v file error: %s\n", path, err.Error())
		return false
	}
	return true
}

func (s *Storage) writeVideosToFile(channelId string, videos []models.Video) bool {
	path := filepath.Join(
		s.AppConfig.CachePath,
		fmt.Sprintf("%s%s%s", consts.P_VIDEOS, channelId, consts.EXT_JSON),
	)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Printf("open %#v file error: %s\n", path, err.Error())
		return false
	}

	if err := json.NewEncoder(file).Encode(&videos); err != nil {
		log.Printf("write to %#v file error: %s\n", path, err.Error())
		return false
	}
	return true
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist) && err == nil
}
