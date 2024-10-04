package aup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/go-cli-tools/config"
	"github.com/lspaccatrosi16/spac/lib/path"
)

func Install() error {
	assetLink := "https://github.com/lspaccatrosi16/aup/releases/latest/download/aup-linux"
	resp, err := http.Get(assetLink)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("Got aup binary of size %.2f MiB\n", float32(resp.ContentLength)/(1024*1024))
	binPath, err := getBinPath()
	if err != nil {
		return err
	}

	targetPath := filepath.Join(binPath, "aup")

	f, err := os.Create(targetPath)

	if err != nil {
		return err
	}

	defer f.Close()
	io.Copy(f, resp.Body)

	err = os.Chmod(targetPath, 0o755)
	if err != nil {
		return err
	}

	configDir, err := config.GetConfigPath("aup")
	if err != nil {
		return err
	}

	pathFile := filepath.Join(configDir, ".auprc")

	path.AddStringToPath(fmt.Sprintf("source %s", path.Abs(pathFile)))

	return nil
}

func getHome() (string, error) {
	homeDir, found := os.LookupEnv("HOME")
	if !found {
		return "", fmt.Errorf("$HOME is unset")
	}

	return homeDir, nil

}

func getBinPath() (string, error) {
	homeDir, err := getHome()

	if err != nil {
		return "", err
	}

	binDir := filepath.Join(homeDir, "bin")
	return binDir, nil

}
