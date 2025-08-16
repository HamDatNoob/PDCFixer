package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"

	"github.com/schollz/progressbar/v3"
)

// CheckForge assumes that Forge is correctly installed based on the existence of the version directory
func CheckForge(mcVersion string) bool {
	forgeVersion := mcVersion + "-forge" + ForgeVersion[mcVersion]
	forgeVersionDirPath := path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftVersionsDir, forgeVersion)
	_, err := os.Stat(forgeVersionDirPath)
	return os.IsNotExist(err)
}

func DownloadForge(mcVersion string) {
	req, _ := http.NewRequest("GET", fmt.Sprintf(ForgeDownload, ForgeVersion[mcVersion], ForgeVersion[mcVersion]), nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	file, _ := os.OpenFile(path.Join(tempDir, fmt.Sprintf(ForgeFile, ForgeVersion[mcVersion])), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Forge",
	)

	io.Copy(io.MultiWriter(file, bar), resp.Body)
}

func InstallForge(mcVersion string) {
	archive, _ := zip.OpenReader(path.Join(tempDir, fmt.Sprintf(ForgeFile, ForgeVersion[mcVersion])))
	defer archive.Close()

	// Copying the version profile
	profileFile, _ := archive.Open("install_profile.json")
	defer profileFile.Close()
	decoder := json.NewDecoder(profileFile)
	var profileJson map[string]interface{}
	decoder.Decode(&profileJson)

	profileJsonString, _ := json.MarshalIndent(profileJson["versionInfo"], "", "  ")
	forgeVersion := mcVersion + "-forge" + ForgeVersion[mcVersion]

	forgeVersionDirPath := path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftVersionsDir, forgeVersion)
	os.Mkdir(forgeVersionDirPath, os.ModePerm)

	forgeProfileFile, _ := os.OpenFile(path.Join(forgeVersionDirPath, forgeVersion+".json"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer forgeProfileFile.Close()

	forgeProfileFile.Write(profileJsonString)

	// Copying the universal jar
	baseJarName := "forge-" + ForgeVersion[mcVersion]
	universalJarName := baseJarName + "-universal.jar"
	libraryJarName := baseJarName + ".jar"
	universalJar, _ := archive.Open(universalJarName)
	defer universalJar.Close()

	libDir := path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftLibrariesDir, ForgeLibDir, ForgeVersion[mcVersion])
	_, err := os.Stat(libDir)
	if os.IsNotExist(err) {
		os.MkdirAll(libDir, os.ModePerm)
	}

	libraryJarFile, _ := os.OpenFile(path.Join(libDir, libraryJarName), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer libraryJarFile.Close()

	universalJarStat, _ := universalJar.Stat()

	bar := progressbar.DefaultBytes(
		universalJarStat.Size(),
		"Copying universal Forge jar",
	)
	io.Copy(io.MultiWriter(libraryJarFile, bar), universalJar)
}
