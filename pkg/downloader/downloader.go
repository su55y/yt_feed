package downloader

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist) && err == nil
}

func download_file(url, path string, wg *sync.WaitGroup) {
	defer wg.Done()

	if exists(path) {
		return
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("download from '%s' error: %s\n", url, err)
		return
	}
	defer resp.Body.Close()

	out, err := os.Create(path)
	if err != nil {
		log.Printf("create '%s' file error: %s\n", path, err.Error())
		return
	}
	defer out.Close()

	io.Copy(out, resp.Body)
}

func DownloadAll(urls map[string]string) {
	var wg sync.WaitGroup
	for k, u := range urls {
		wg.Add(1)
		go download_file(u, k, &wg)
	}
	wg.Wait()
}
