package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/su55y/yt_feed/internal/blocks"
	"github.com/su55y/yt_feed/internal/config"
	"github.com/su55y/yt_feed/internal/consts"
	"github.com/su55y/yt_feed/internal/models"
	"github.com/su55y/yt_feed/internal/service"
	"github.com/su55y/yt_feed/internal/storage"
	"google.golang.org/api/youtube/v3"
)

var (
	yt      *youtube.Service
	conf    Config
	appConf config.AppConfig
	pageNum uint8 = 1
	rawIn   string
)

type Config struct {
	AppCachePath  string
	HomePath      string
	ConfPathRoot  string
	CachePathRoot string
	ConfFullPath  string
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist) && err == nil
}

func readEnv() {
	// set /home/user/.config
	if conf.ConfPathRoot = os.Getenv(consts.ENV_CONFIG_HOME); !exists(conf.ConfPathRoot) {
		conf.ConfPathRoot = filepath.Join(conf.HomePath, consts.DEF_CONFIG_PATH)
	}

	// set /home/user/.cache
	if conf.CachePathRoot = os.Getenv(consts.ENV_CACHE_HOME); !exists(conf.CachePathRoot) {
		conf.CachePathRoot = filepath.Join(conf.HomePath, consts.DEF_CACHE_PATH)
	}
}

func getAppConfig() {
	appConfDirPath := filepath.Join(
		conf.ConfPathRoot,
		consts.APP_NAME,
	)

	appConfFilePath := filepath.Join(
		appConfDirPath,
		consts.APP_CONFIG_NAME,
	)

	if !exists(appConfDirPath) {
		if err := os.MkdirAll(appConfDirPath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat(appConfFilePath); errors.Is(err, os.ErrNotExist) {
		log.Printf(consts.INF_NEW_CONFIG, appConfFilePath)
		ioutil.WriteFile(appConfFilePath, []byte(consts.DEF_CONFIG), 0666)
	}

	var err error
	appConf, err = config.GetAppConfig(appConfFilePath)
	if err != nil {
		log.Printf(consts.ERR_CONFIG_LOAD, err)
	}
}

func init() {
	var err error
	if conf.HomePath, err = os.UserHomeDir(); err != nil {
		log.Fatal(err)
	}

	readEnv()
	getAppConfig()

	if len(appConf.API_KEY) == 0 {
		if exists(appConf.ApiKeyPath) {
			apiBytes, err := ioutil.ReadFile(appConf.ApiKeyPath)
			if err != nil {
				log.Fatal(fmt.Errorf(consts.ERR_NO_API_KEY_FILE, appConf.ApiKeyPath, err))
			}

			if appConf.API_KEY = strings.TrimSpace(string(apiBytes)); len(appConf.API_KEY) == 0 {
				log.Fatal(fmt.Errorf(consts.ERR_API_KEY_FILE_READ, appConf.ApiKeyPath))
			}
		} else {
			if appConf.API_KEY = os.Getenv(consts.ENV_YT_API_KEY); len(appConf.API_KEY) == 0 {
				log.Fatal(fmt.Errorf("%s", consts.ERR_NO_API_KEY))
			}
		}
	}

	conf.AppCachePath = filepath.Join(conf.CachePathRoot, consts.APP_NAME)
	if !exists(conf.AppCachePath) {
		if err := os.MkdirAll(conf.AppCachePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if exists(appConf.CachePath) {
		conf.AppCachePath = appConf.CachePath
	} else {
		appConf.CachePath = conf.AppCachePath
	}

	appConf.ThumbDir = filepath.Join(appConf.CachePath, consts.THUMB_DIR_NAME)
	if !exists(appConf.ThumbDir) {
		if err := os.MkdirAll(appConf.ThumbDir, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}
}

func openInMPV(id string) bool {
	c := exec.Command("mpv", "https://www.youtube.com/watch?v="+id)
	if err := c.Start(); err != nil {
		log.Println(err.Error())
		return false
	}

	return c.Process.Pid > 0
}

func main() {
	f, err := os.OpenFile(
		filepath.Join(conf.AppCachePath, "log"),
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	// channels := make([]models.Channel, 0)
	blocksOutput := models.Blocks{}

	// updChan := make(chan bool)
	ytService := service.New(context.Background(), &appConf)
	stor := storage.New(&appConf, &ytService)

	channels, err := stor.ReadChannels()
	if err != nil {
		log.Fatal(err)
	}

	blocksOutput.Lines = blocks.PrintChannels(channels)
	blocksOutput.Message = "updating..."
	j, err := json.Marshal(&blocksOutput)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(j))

	// updRes := <-updChan
	// if updRes {
	// 	blocksOutput.Message = "updating...done"
	// } else {
	// 	blocksOutput.Message = "updating...failed"
	// }

	var runMPV bool

	inDecoder := json.NewDecoder(os.Stdin)
	blocksInput := models.BlocksIn{}
	currentChannel := ""

	for {
		if err := inDecoder.Decode(&blocksInput); err != nil {
			log.Fatal(err)
		}

		switch blocksInput.Name {
		case consts.IN_SELECT_ENTRY:
			if len(blocksInput.Data) == 24 {
				currentChannel = blocksInput.Data
			}
			switch blocksInput.Value {
			case "videos":
				if videos, err := stor.ReadUploads(blocksInput.Data, false); err != nil {
					blocksOutput.Message = fmt.Sprintf(
						"videos for %s not ready",
						blocksInput.Data,
					)
					log.Printf(
						"can't read %s videos due to error: %s",
						blocksInput.Data,
						err.Error(),
					)
				} else {
					blocksOutput.Lines = []models.Line{{Text: "back"}}
					blocksOutput.Lines = append(blocksOutput.Lines, blocks.PrintVideos(videos)...)
					blocksOutput.Message = fmt.Sprintf(
						"last %d videos of %s",
						len(videos),
						blocksInput.Data,
					)
				}
			case "playlists":
				if playlists, err := stor.ReadAllPlaylists(blocksInput.Data, false); err != nil {
					blocksOutput.Message = fmt.Sprintf("playlists %s not ready", blocksInput.Data)
					log.Printf("can't read playlists for %s due to error: %s",
						blocksInput.Data,
						err.Error(),
					)
				} else {
					blocksOutput.Lines = []models.Line{{Text: "back"}}
					blocksOutput.Lines = append(
						blocksOutput.Lines,
						blocks.PrintPlaylists(playlists)...,
					)
					blocksOutput.Message = fmt.Sprintf(
						"last %d playlists of %s",
						len(playlists),
						blocksInput.Data,
					)
				}
			case "update playlists":
				blocksOutput.Message = "updating playlists for " + blocksInput.Data
				j, _ := json.Marshal(&blocksOutput)
				fmt.Println(string(j))
				if _, err := stor.ReadAllPlaylists(blocksInput.Data, true); err != nil {
					blocksOutput.Message = "error while updating playlists"
				} else {
					blocksOutput.Message = "done"
				}

				blocksOutput.Lines = blocks.PrintChannelMenu(blocksInput.Data)
			case "update videos":
				blocksOutput.Message = "updating videos for " + blocksInput.Data
				j, _ := json.Marshal(&blocksOutput)
				fmt.Println(string(j))
				if _, err := stor.ReadUploads(blocksInput.Data, true); err != nil {
					blocksOutput.Message = "error while updating videos"
				} else {
					blocksOutput.Message = "done"
				}
				blocksOutput.Lines = blocks.PrintChannelMenu(blocksInput.Data)
			case "back":
				blocksOutput.Message = "channels list"
				blocksOutput.Lines = blocks.PrintChannels(channels)
			default:
				if len(blocksInput.Data) == 34 { // playlist
					if playlists, err := stor.ReadAllPlaylists(currentChannel, false); err != nil {
						blocksOutput.Message = "get playlist videos error"
					} else {
						blocksOutput.Lines = []models.Line{
							{Text: "back"},
						}
						blocksOutput.Lines = append(
							blocksOutput.Lines,
							blocks.PrintVideos(playlists[blocksInput.Data].Videos)...,
						)
						blocksOutput.Message = fmt.Sprintf("last %d videos of %s playlist",
							len(blocksOutput.Lines)-1, blocksInput.Data)
					}
				} else if len(blocksInput.Data) == 11 {
					if runMPV = openInMPV(blocksInput.Data); !runMPV {
						blocksOutput.Message += " : error"
					}
				} else {
					blocksOutput.Message = "choose action"
					blocksOutput.Lines = blocks.PrintChannelMenu(blocksInput.Data)
				}
			}
		}

		blocksOutput.Prompt = ""
		j, err := json.Marshal(&blocksOutput)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(j))

		if runMPV {
			time.Sleep(2 * time.Second)
			os.Exit(0)
		}
	}
}
