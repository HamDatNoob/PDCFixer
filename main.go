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
	mcVersion := "1.8.9"

	pdcVersionSelect := promptui.Select{
		Label: "PDC Version",
		Items: []string{"Vanilla", "Forge"},
	}

	_, pdcVersion, _ := pdcVersionSelect.Run()

	ramSelect := promptui.Prompt{
		Label:    "RAM (GB)",
		Default:  strconv.Itoa(DefaultRAM),
		Validate: validateRAM,
	}

	ramString, _ := ramSelect.Run()
	ramAmount, _ := strconv.Atoi(ramString)

	forge := false

	switch pdcVersion {
	case "Forge":
		forge = true
	}

	baseVersion := mcVersion
	suffix := "Vanilla"
	profileVersion := mcVersion

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
	fmt.Println("You must close and re-open the Minecraft Launcher if you haven't")
	fmt.Scanln()
}
