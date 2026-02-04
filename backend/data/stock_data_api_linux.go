package data

import "os/exec"

// CheckBrowser 在 Linux 上查找可用浏览器路径
func CheckBrowser() (string, bool) {
	candidates := []string{
		"google-chrome",
		"google-chrome-stable",
		"chromium",
		"chromium-browser",
		"microsoft-edge",
		"brave-browser",
	}
	for _, name := range candidates {
		if path, err := exec.LookPath(name); err == nil {
			return path, true
		}
	}
	return "", false
}
