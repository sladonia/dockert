package docker

import (
	"os"
	"runtime"
)

func IsRunningInDockerContainer() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

func IsDarwinOS() bool {
	return runtime.GOOS == "darwin"
}
