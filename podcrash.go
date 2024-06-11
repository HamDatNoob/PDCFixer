package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
)

func DownloadPDC(mcVersion string) {
	pdcVersionsDir := path.Join(configDir, MinecraftDir[runtime.GOOS], PDCDir, MinecraftVersionsDir)
	_, err := os.Stat(pdcVersionsDir)
	if os.IsNotExist(err) {
		os.MkdirAll(pdcVersionsDir, 0644)
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf(PDCDownloadUrl, mcVersion), nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	file, _ := os.OpenFile(path.Join(pdcVersionsDir, fmt.Sprintf(FileName, mcVersion)), os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading PDC",
	)

	io.Copy(io.MultiWriter(file, bar), resp.Body)
}
