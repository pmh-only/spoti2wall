package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/ini.v1"
)

const tmpPath = "/tmp/spoti2wall/"

func CalcTmpPath(url string, blur, dark int) string {
	imgId := strings.Split(url, "image/")[1]

	return fmt.Sprintf("%s%s-b%d-d%d.jpg", tmpPath, imgId, blur, dark)
}

func CheckImageExist(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}

func DownloadImage(url string, path string) error {
	if err := os.MkdirAll(tmpPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download image: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download image, got status code: %d", res.StatusCode)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, res.Body); err != nil {
		return fmt.Errorf("failed to write image to file: %v", err)
	}

	return nil
}

func GetConfigPath() string {
	configDir, _ := os.UserConfigDir()
	os.MkdirAll(configDir, os.ModePerm)
	return configDir + "/spoti2wall.ini"
}

func SaveRefreshToken(token string) {
	configPath := GetConfigPath()

	conf, _ := ini.Load(configPath)
	conf.Section("").Key("refresh_token").SetValue(token)

	conf.SaveTo(configPath)
}

func SaveClientId(id string) {
	configPath := GetConfigPath()

	conf, _ := ini.Load(configPath)
	conf.Section("").Key("client_id").SetValue(id)
	conf.SaveTo(configPath)
}

func SaveClientSecret(secret string) {
	configPath := GetConfigPath()

	conf, _ := ini.Load(configPath)
	conf.Section("").Key("client_secret").SetValue(secret)
	conf.SaveTo(configPath)
}

func ReadRefreshToken() (token string) {
	configPath := GetConfigPath()

	conf, _ := ini.Load(configPath)
	token = conf.Section("").Key("refresh_token").String()

	return
}

func OpenBrowser(url string) {
	cmd := exec.Command("xdg-open", url)
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Run()
}
