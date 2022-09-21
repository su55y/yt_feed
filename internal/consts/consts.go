package consts

const (
	// app consts
	APP_NAME        = "yt_feed"
	APP_CONFIG_NAME = "config.yaml"

	// blocks input action names
	IN_EXECUTE_CUSTOM_ITEM = "execute custom input"
	IN_CUSTOM_KEY          = "custom key"
	IN_SELECT_ENTRY        = "select entry"
	IN_ACTIVE_ENTRY        = "active entry"

	// app env names
	ENV_YT_API_KEY   = "YT_FEED_API_KEY"
	ENV_YT_CACHE_DIR = "YT_FEED_CACHE_DIR"

	// env
	ENV_CONFIG_HOME = "XDG_CONFIG_HOME"
	ENV_CACHE_HOME  = "XDG_CACHE_HOME"

	// defaults
	DEF_CACHE_PATH  = ".cache"
	DEF_CONFIG_PATH = ".config/yt_feed/config.yaml"
	DEF_CONFIG      = `# youtube api key (https://console.cloud.google.com/)
# api_key: "<YT_API_KEY>"
# api_key_path: "/path/to/api_key"

# max results
max_results: 100

# absolute path to alternative cache dir
# "/home/user/.cache/yt_feed" by default
# cache_dir: "/path/to/cache"

# thumbnails are loaded into the cache directory with format '(h/m/d)(t)(video_id).ext' 
# you can disable thumbnails loading
thumbnails_disable: false

# thumbnails size: high(~15-30k),medium(~8-15k),default(~3-4k)
thumbnails_size: "default"

# channels is an array of channels ids
# channels:
#   - "value1"
#   - "value2"
#   - "value3"
# or
# channels: [
#   "value1",
#   "value2",
#   "value3"
# ]`

	// info
	INF_NEW_CONFIG   = "new config written to %s"
	INF_SELECT_PARSE = "can't read video id"

	// errors
	ERR_NO_API_KEY        = "api key was not found either in the config or in the env variable"
	ERR_NO_API_KEY_FILE   = "api key not found in '%s': %v"
	ERR_API_KEY_FILE_READ = "no api key in '%s'"
	ERR_CONFIG_LOAD       = "load config error: %v"

	// size prefixes
	SP_HIGH    = "high"
	SP_MEDIUM  = "medium"
	SP_DEFAULT = "default"

	// files prefixes
	P_VIDEOS    = "videos"
	P_PLAYLISTS = "playlists"

	// extensions
	EXT_JSON = ".json"
	EXT_JPG  = ".jpg"

	// dirs
	THUMB_DIR_NAME = "thumbnails"
)
