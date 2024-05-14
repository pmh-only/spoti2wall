package utils

import (
	"fmt"
	"os/exec"

	"github.com/fatih/color"
)

func checkImageMagickExist() bool {
	cmd := exec.Command("convert", "-h")
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()

	return err == nil
}

func BlurImage(path string, blurOption int) {
	if !checkImageMagickExist() {
		color.Red("❌ You need to install `imagemagick` to use --blur flag.")
		return
	}

	cmd := exec.Command(
		"convert", path,
		"-blur",
		fmt.Sprintf("0x%d", blurOption),
		path)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}

func DarkImage(path string, darkOption int) {
	if !checkImageMagickExist() {
		color.Red("❌ You need to install `imagemagick` to use --dark flag.")
		return
	}

	cmd := exec.Command(
		"convert", path,
		"-fill", "black",
		"-colorize",
		fmt.Sprintf("%d%%", darkOption),
		path)

	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}
