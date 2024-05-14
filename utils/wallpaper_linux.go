package utils

import (
	"os/exec"
)

func CheckWallpaperCliExist() bool {
	cmd := exec.Command("nitrogen", "-h")
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()

	return err == nil
}

func ApplyWallpaper(path string) {
	cmd := exec.Command("nitrogen", "--set-zoom-fill", path)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}

func RestoreWallpaper() {
	cmd := exec.Command("nitrogen", "--restore")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}
