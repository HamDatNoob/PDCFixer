package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/manifoldco/promptui"
)

func validateRAM(input string) error {
	_, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("please enter a valid number")
	}
	return nil
}

func main() {
	mcVersionSelect := promptui.Select{
		Label: "Minecraft Version",
		Items: []string{"1.8.8", "1.8.9", "1.9"},
	}

	_, mcVersion, _ := mcVersionSelect.Run()

	pdcVersionSelect := promptui.Select{
		Label: "PDC Version",
		Items: []string{"Vanilla", "Optifine", "Forge"},
	}

	_, pdcVersion, _ := pdcVersionSelect.Run()

	ramSelect := promptui.Prompt{
		Label:    "RAM (GB)",
		Default:  strconv.Itoa(DefaultRAM),
		Validate: validateRAM,
	}

	ramString, _ := ramSelect.Run()
	ramAmount, _ := strconv.Atoi(ramString)

	optifine := false
	forge := false

	switch pdcVersion {
	case "Optifine":
		optifine = true
	case "Forge":
		forge = true
	}

	baseVersion := mcVersion
	suffix := "Vanilla"
	profileVersion := mcVersion

	if optifine {
		if CheckOptifine(mcVersion) {
			DownloadJre()
			DownloadLaunchwrapper()
			DownloadOptifine(mcVersion)
			InstallOptifine(mcVersion)
		}
		suffix = "Optifine"
		profileVersion = fmt.Sprintf(OptifineVersionDir, baseVersion, OptifineVersion[baseVersion])
	}

	if forge {
		if CheckForge(mcVersion) {
			DownloadForge(mcVersion)
			InstallForge(mcVersion)
		}
		mcVersion += "-forge"
		suffix = "Forge"
		profileVersion = fmt.Sprintf("%s-forge%s", baseVersion, ForgeVersion[baseVersion])
	}

	DownloadPDC(mcVersion)

	fmt.Println("Cleaning up temporary files...")
	os.RemoveAll(tempDir)

	CreateProfile(mcVersion, baseVersion, suffix, profileVersion, ramAmount)
	fmt.Println("PDC has been installed to the Minecraft Launcher")
	fmt.Scanln()
}
