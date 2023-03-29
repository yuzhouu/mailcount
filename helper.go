package main

import (
	"os/exec"
	"runtime"
)

func openbrowser(url string) error {
	var cmd string
	switch runtime.GOOS {
	case "windows":
		cmd = "start"
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	return exec.Command(cmd, url).Start()
}
