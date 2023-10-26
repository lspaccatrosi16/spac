package golang

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/spac/lib/install/sudo"
	"github.com/lspaccatrosi16/spac/lib/path"
)

func Install() error {
start:
	langVersion := input.GetValidatedInput("Golang version", validate)

	goUrl := fmt.Sprintf("https://go.dev/dl/go%s.linux-amd64.tar.gz", langVersion)

	resp, err := http.Get(goUrl)

	if resp.StatusCode == 404 {
		fmt.Printf("invalid go version: %s\n", langVersion)
		goto start
	}

	if err != nil {
		return err
	}

	fmt.Printf("successfuly got go v%s\n", langVersion)

	defer resp.Body.Close()

	tmpDir := os.TempDir()
	tmpFolder := filepath.Join(tmpDir, fmt.Sprintf("go-%s", langVersion))

	err = os.RemoveAll(tmpFolder)

	if err != nil {
		return err
	}

	err = os.MkdirAll(tmpFolder, 0o755)

	if err != nil {
		return err
	}

	tmpFile := filepath.Join(tmpFolder, "go.tar.gz")
	f, err := os.Create(tmpFile)

	if err != nil {
		return err
	}

	defer f.Close()

	io.Copy(f, resp.Body)

	if err != nil {
		return err
	}

	err = sudo.RunSudo("rm -rf /usr/local/go")

	if err != nil {
		return err
	}

	err = sudo.RunSudo(fmt.Sprintf("tar -C /usr/local -xzf %s", tmpFile))

	if err != nil {
		return err
	}

	err = os.RemoveAll(tmpFolder)

	if err != nil {
		return err
	}

	err = path.AddDirToPath("/usr/local/go/bin")

	if err != nil {
		return err
	}

	return nil
}

func validate(in string) error {
	components := strings.Split(in, ".")
	if len(components) != 3 {
		return fmt.Errorf("go version must be in the format x.x.x")
	}

	for i, c := range components {
		_, err := strconv.ParseInt(c, 10, 64)
		if err != nil {
			return fmt.Errorf("component %d is not an integer", i)
		}
	}
	return nil
}
