package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/schollz/progressbar/v3"
)

const OptifineVersionProfileJson = `{
	  "id": "%s",
	  "inheritsFrom": "%s",
	  "time": "%s",
      "releaseTime": "%s",
	  "type": "release",
	  "libraries": [
	    {
          "name": "optifine:OptiFine:%s"
	    },
	    {
	      "name": "%s"
        }
	  ],
	  "mainClass": "net.minecraft.launchwrapper.Launch",
	  "minecraftArguments": "--username ${auth_player_name} --version ${version_name} --gameDir ${game_directory} --assetsDir ${assets_root} --assetIndex ${assets_index_name} --uuid ${auth_uuid} --accessToken ${auth_access_token} --userType ${user_type} --versionType ${version_type}  --tweakClass optifine.OptiFineTweaker"
	}`

// CheckOptifine assumes that Optifine is correctly installed based on the existence of the version directory
func CheckOptifine(mcVersion string) bool {
	mcVersion = FixVersionOptifine(mcVersion)

	optifineVersion := fmt.Sprintf(OptifineVersionDir, mcVersion, OptifineVersion[mcVersion])
	optifineVersionDirPath := path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftVersionsDir, optifineVersion)
	_, err := os.Stat(optifineVersionDirPath)
	return os.IsNotExist(err)
}

func FixVersionOptifine(mcVersion string) string {
	if mcVersion == "1.9" {
		return "1.9.0"
	}
	return mcVersion
}

func DownloadOptifine(mcVersion string) {
	ofVersion := FixVersionOptifine(mcVersion)

	fmt.Println("Fetching Optifine download link...")
	htmlReq, _ := http.NewRequest("GET", fmt.Sprintf(OptifineDownload, ofVersion, OptifineVersion[mcVersion]), nil)
	htmlResp, _ := http.DefaultClient.Do(htmlReq)
	defer htmlResp.Body.Close()

	bytes, _ := io.ReadAll(htmlResp.Body)
	body := string(bytes)
	downloadUrl := strings.Split(strings.Split(body, "<a href='")[1], "'")[0]

	file, _ := os.OpenFile(path.Join(tempDir, fmt.Sprintf(OptifineFile, mcVersion, OptifineVersion[mcVersion])), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	dlUrl := OptifineUrl + downloadUrl
	fmt.Printf("Optifine download link: %s\n", dlUrl)
	req, _ := http.NewRequest("GET", dlUrl, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Optifine",
	)

	io.Copy(io.MultiWriter(file, bar), resp.Body)
}

func InstallOptifine(mcVersion string) {
	// Patch Optifine using their stupid system
	javaPath := path.Join(tempDir, JavaPath[runtime.GOOS])
	optifineDlPath := path.Join(tempDir, fmt.Sprintf(OptifineFile, mcVersion, OptifineVersion[mcVersion]))
	optifineLibPath := path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftLibrariesDir, OptifineLibraryDir,
		fmt.Sprintf(OptifineLibraryVersionDir, mcVersion, OptifineVersion[mcVersion]))
	_, err := os.Stat(optifineLibPath)
	if os.IsNotExist(err) {
		os.MkdirAll(optifineLibPath, os.ModePerm)
	}

	fmt.Println("Installing Optifine...")
	cmd := exec.Command(
		javaPath,
		"-cp",
		optifineDlPath,
		"optifine.Patcher",
		path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftVersionsDir, mcVersion, mcVersion+".jar"),
		optifineDlPath,
		path.Join(optifineLibPath, fmt.Sprintf(OptifineFile, mcVersion, OptifineVersion[mcVersion])),
	)
	cmd.Run()

	// Install Optifine's own launchwrapper
	optifineLWDir := path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftLibrariesDir, OptifineLaunchwrapperDir)
	_, err = os.Stat(optifineLWDir)
	if os.IsNotExist(err) {
		os.MkdirAll(optifineLWDir, os.ModePerm)
	}

	archive, _ := zip.OpenReader(optifineDlPath)
	defer archive.Close()

	launchwrapperArchive, err := archive.Open(OptifineLaunchwrapperFile)

	if err == nil {
		stat, _ := launchwrapperArchive.Stat()

		lwFile, _ := os.OpenFile(path.Join(optifineLWDir, OptifineLaunchwrapperFile), os.O_CREATE|os.O_WRONLY, os.ModePerm)
		defer lwFile.Close()

		bar := progressbar.DefaultBytes(
			stat.Size(),
			"Extracting Launchwrapper",
		)

		io.Copy(io.MultiWriter(lwFile, bar), launchwrapperArchive)
		defer launchwrapperArchive.Close()
	}

	// Create version profile
	optifineVerDir := fmt.Sprintf(OptifineVersionDir, mcVersion, OptifineVersion[mcVersion])

	optifineVersionProfileDir := path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftVersionsDir, optifineVerDir)
	_, err = os.Stat(optifineVersionProfileDir)
	if os.IsNotExist(err) {
		os.MkdirAll(optifineVersionProfileDir, os.ModePerm)
	}

	timestamp := time.Now().Format(time.RFC3339)
	launchwrapper := "optifine:launchwrapper-of:2.2"
	if mcVersion == "1.8.8" {
		launchwrapper = "net.minecraft:launchwrapper:1.12"
	}

	optifineProfileJson := fmt.Sprintf(OptifineVersionProfileJson,
		fmt.Sprintf(OptifineVersionDir, mcVersion, OptifineVersion[mcVersion]),
		mcVersion,
		timestamp,
		timestamp,
		fmt.Sprintf(OptifineLibraryVersionDir, mcVersion, OptifineVersion[mcVersion]),
		launchwrapper,
	)

	optifineVersionJsonFile, _ := os.OpenFile(path.Join(optifineVersionProfileDir, optifineVerDir+".json"), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer optifineVersionJsonFile.Close()
	optifineVersionJsonFile.WriteString(optifineProfileJson)
}

func DownloadJre() {
	// Download the JRE
	req, _ := http.NewRequest("GET", JreDownloads[runtime.GOOS], nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	jrePath := path.Join(tempDir, "jre.zip")
	file, _ := os.OpenFile(jrePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Java",
	)

	io.Copy(io.MultiWriter(file, bar), resp.Body)

	// Unzip the JRE
	archive, _ := zip.OpenReader(jrePath)
	defer archive.Close()

	for _, archiveFile := range archive.File {
		path := path.Join(tempDir, archiveFile.Name)

		if archiveFile.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		extractBar := progressbar.DefaultBytes(
			archiveFile.FileInfo().Size(),
			fmt.Sprintf("Extracting %s...", archiveFile.Name),
		)

		destinationFile, _ := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		archiveItem, _ := archiveFile.Open()
		io.Copy(io.MultiWriter(destinationFile, extractBar), archiveItem)

		destinationFile.Close()
		archiveItem.Close()
	}
}

func DownloadLaunchwrapper() {
	req, _ := http.NewRequest("GET", MinecraftLaunchwrapperDownload, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	libPath := path.Join(configDir, MinecraftDir[runtime.GOOS], MinecraftLibrariesDir, MinecraftLaunchwrapperLocation)
	_, err := os.Stat(libPath)
	if os.IsNotExist(err) {
		os.MkdirAll(libPath, os.ModePerm)
	}

	file, _ := os.OpenFile(path.Join(libPath, MinecraftLaunchwrapperName), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"Downloading Mojang Launchwrapper",
	)

	io.Copy(io.MultiWriter(file, bar), resp.Body)
}
