package utils

import (
	"syscall"
	"unsafe"

	"github.com/fatih/color"
)

const (
	SPI_GETDESKWALLPAPER = 0x0073
	SPI_SETDESKWALLPAPER = 0x0014
	SPIF_UPDATEINIFILE   = 0x01
	SPIF_SENDCHANGE      = 0x02

	MAX_PATH = 260
)

type WallpaperManager struct {
	DefaultPath string
}

func NewWallpaperManager(defaultPath string) *WallpaperManager {
	return &WallpaperManager{
		DefaultPath: defaultPath,
	}
}

func GetWallpaperPath() (string, error) {
	var wallpaperPath [MAX_PATH]uint16

	ret, _, err := syscall.NewLazyDLL("user32.dll").NewProc("SystemParametersInfoW").Call(
		SPI_GETDESKWALLPAPER,
		MAX_PATH,
		uintptr(unsafe.Pointer(&wallpaperPath[0])),
		0,
	)
	if ret == 0 {
		return "", err
	}

	return syscall.UTF16ToString(wallpaperPath[:]), nil
}

func (w *WallpaperManager) ApplyWallpaper(path string) {
	wallpaperPathUTF16, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		panic(err)
	}

	proc := syscall.NewLazyDLL("user32.dll").NewProc("SystemParametersInfoW")

	ret, _, _ := proc.Call(
		SPI_SETDESKWALLPAPER,
		0,
		uintptr(unsafe.Pointer(wallpaperPathUTF16)),
		SPIF_UPDATEINIFILE|SPIF_SENDCHANGE,
	)

	if ret == 0 {
		color.Red("❌ Cannot apply wallpaper.")
	}
}

func (w *WallpaperManager) RestoreWallpaper() {
	w.ApplyWallpaper(w.DefaultPath)
}
