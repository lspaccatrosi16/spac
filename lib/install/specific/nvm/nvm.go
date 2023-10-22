package nvm

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/spac/lib/install/sudo"
)

func Install() error {
	scriptUrl := "https://github.com/nvm-sh/nvm/blob/master/install.sh"
	resp, err := http.Get(scriptUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	tmpDir := os.TempDir()

	tmpFile := filepath.Join(tmpDir, "nvm-install.sh")

	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}

	io.Copy(f, resp.Body)
	f.Close()

	err = os.Chmod(tmpFile, 0o755)
	if err != nil {
		return err
	}

	err = sudo.RunSudo(fmt.Sprintf("sh %s", tmpFile))

	if err != nil {
		return err
	}

	return nil
}
