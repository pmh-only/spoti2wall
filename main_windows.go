package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/pmh-only/spoti2wall/rest"
	"github.com/pmh-only/spoti2wall/state"
	"github.com/pmh-only/spoti2wall/utils"
)

var blurFlag int
var darkFlag int
var reauthFlag bool

func init() {
	flag.IntVar(&blurFlag, "b", 0, "Blur image with blur option")
	flag.IntVar(&darkFlag, "d", 0, "Dark image with darker option")
	flag.BoolVar(&reauthFlag, "reauth", false, "Reauth with spotify")
	flag.Parse()

	rest.RefreshToken = utils.ReadRefreshToken()

	state.InitConfig()
}

func getClientId() string {
	return state.GlobalConfig.Section("").Key("client_id").String()
}

func main() {
	color.Magenta("🎵 spoti2wall started...")

	defaultWallpaperPath, err := utils.GetWallpaperPath()
	if err != nil {
		color.Red("❌ Cannot get default wallpaper path.")
	}

	if getClientId() == "" {
		color.Green("📝 Enter client id [default]: ")
		scanner := bufio.NewScanner(os.Stdin)
		_ = scanner.Scan()
		if scanner.Text() == "" {
			utils.SaveClientId(rest.ClientId)
		} else {
			utils.SaveClientId(scanner.Text())
		}

		color.Green("📝 Enter client secret [default]: ")
		_ = scanner.Scan()
		if scanner.Text() == "" {
			utils.SaveClientSecret(rest.ClientSecret)
		} else {
			utils.SaveClientSecret(scanner.Text())
		}

		// reinit for refresh values
		state.InitConfig()
	}

	wallpaperManager := utils.NewWallpaperManager(defaultWallpaperPath)

	if reauthFlag || rest.RefreshToken == "" {
		rest.StartAuthServer()
	} else {
		rest.AccessToken = rest.RefreshAccessToken(rest.RefreshToken)
		go rest.KeepRefreshToken()
	}

	prevImage := ""

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			image := rest.GetTrackImage()

			if image == "" {
				wallpaperManager.RestoreWallpaper()
				prevImage = ""
				continue
			}

			if image == prevImage {
				continue
			}

			prevImage = image

			filepath := utils.CalcTmpPath(image, blurFlag, darkFlag)

			if !utils.CheckImageExist(filepath) {

				color.Yellow("🔍 Downloading image...")
				err := utils.DownloadImage(image, filepath)
				if err != nil {
					color.Red("❌ Failed to download image: %v", err)
					continue
				}

				if blurFlag > 0 {
					utils.BlurImage(filepath, blurFlag)
				}

				if darkFlag > 0 {
					utils.DarkImage(filepath, darkFlag)
				}
			} else {
				color.Yellow("📁 Image already exists. Skipping download.")
			}

			wallpaperManager.ApplyWallpaper(fmt.Sprintf("C:%s", filepath))
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGTERM,
		syscall.SIGINT,
		os.Interrupt)

	<-sig

	color.Yellow("🔌 Shutting down...")
	wallpaperManager.RestoreWallpaper()
}
