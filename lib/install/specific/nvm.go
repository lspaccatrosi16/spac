package specific

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func InstallNvm() error {
	scriptUrl := "https://github.com/nvm-sh/nvm/blob/master/install.sh"
	resp, err := http.Get(scriptUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	f, err := os.CreateTemp("", "*")
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(f, resp.Body)

	stats, err := f.Stat()

	if err != nil {
		return err
	}

	fmt.Println(stats.Name())
	return nil

}
