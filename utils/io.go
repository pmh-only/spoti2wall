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

func DownloadImage(url string, path string) {
	os.MkdirAll(tmpPath, os.ModePerm)

	res, _ := http.Get(url)
	out, _ := os.Create(path)

	io.Copy(out, res.Body)
	res.Body.Close()
}

func SaveRefreshToken(token string) {
	configDir, _ := os.UserConfigDir()
	configPath := configDir + "/spoti2wall.ini"

	os.MkdirAll(configDir, os.ModePerm)
	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		os.Create(configPath)
	}

	conf, _ := ini.Load(configPath)
	conf.Section("").Key("refresh_token").SetValue(token)

	conf.SaveTo(configPath)
}

func ReadRefreshToken() (token string) {
	configDir, _ := os.UserConfigDir()
	configPath := configDir + "/spoti2wall.ini"

	os.MkdirAll(configDir, os.ModePerm)
	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		os.Create(configPath)
	}

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
