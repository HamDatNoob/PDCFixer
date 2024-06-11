package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"time"
)

type Profile struct {
	Created       string `json:"created"`
	Icon          string `json:"icon"`
	JavaArgs      string `json:"javaArgs"`
	LastUsed      string `json:"lastUsed"`
	LastVersionId string `json:"lastVersionId"`
	Name          string `json:"name"`
	Type          string `json:"type"`
}

func CreateProfile(mcVersion string, baseVersion string, suffix string, profileVersion string) {
	filePath := path.Join(configDir, MinecraftDir[runtime.GOOS], "launcher_profiles.json")
	jsonFile, _ := os.Open(filePath)

	pdcDir := path.Join(configDir, MinecraftDir[runtime.GOOS], PDCDir)

	os.WriteFile(path.Join(pdcDir, "games.txt"), []byte(GameOptions), 0644)

	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var launcherProfiles map[string]interface{}
	json.Unmarshal(byteValue, &launcherProfiles)
	timestamp := time.Now().Format(time.RFC3339)

	pdcFile := path.Join(pdcDir, MinecraftVersionsDir, fmt.Sprintf(FileName, mcVersion))

	podcrashProfile := Profile{
		Name:          fmt.Sprintf(ProfileName, baseVersion, suffix),
		Type:          ProfileType,
		Created:       timestamp,
		LastUsed:      timestamp,
		Icon:          ProfileIcon,
		LastVersionId: profileVersion,
		JavaArgs:      fmt.Sprintf(JavaArgs, DefaultRAM, pdcDir, pdcFile),
	}

	profiles, _ := launcherProfiles["profiles"].(map[string]interface{})
	profiles["podcrash-"+baseVersion+"-"+suffix] = podcrashProfile

	jsonString, _ := json.MarshalIndent(launcherProfiles, "", "  ")
	os.WriteFile(filePath, jsonString, 0644)
}
