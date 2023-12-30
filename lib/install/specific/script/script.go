package script

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/spac/lib/install/sudo"
)

func UseScript(url string, name string) func() error {
	return func() error {
		tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("%s-install.sh", name))
		f, err := os.Create(tmpFile)
		if err != nil {
			return err
		}

		resp, err := http.Get(url)

		fmt.Printf("downloaded %s install script\n", name)

		if err != nil {
			return err
		}

		io.Copy(f, resp.Body)

		f.Close()
		resp.Body.Close()

		err = os.Chmod(tmpFile, 0o755)
		if err != nil {
			return err
		}

		err = sudo.RunSudo(fmt.Sprintf("sh %s", tmpFile))

		if err != nil {
			return err
		}

		err = os.Remove(tmpFile)

		if err != nil {
			return err
		}

		return nil
	}
}
