package config

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/command"
)

type configItem struct {
	Name string
	File string
	Path string
}

type configList []configItem

func (c configList) Len() int {
	return len(c)
}

func (c configList) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c configList) Less(i int, j int) bool {
	return c[i].Name < c[j].Name
}

func Config() error {
	manager := command.NewManager(command.ManagerConfig{Searchable: true})
	items := list()

	sort.Sort(items)

	for _, item := range items {
		manager.Register(item.Name, fmt.Sprintf("update %s configuration", item.Name), UseConfig(item))
	}

	for {
		exit := manager.Tui()

		if exit {
			break
		}
	}

	return nil
}

func UseConfig(item configItem) func() error {
	return func() error {
		url := fmt.Sprintf("https://raw.githubusercontent.com/lspaccatrosi16/luca-config/master/%s", item.File)

		fmt.Printf("Downloading %s config from lspaccatrosi16/luca-config\n", item.Name)

		resp, err := http.Get(url)

		if err != nil {
			return err
		}

		homeDir, found := os.LookupEnv("HOME")
		if !found {
			return fmt.Errorf("could not resolve HOME env var")
		}

		parsedBasePath := strings.ReplaceAll(item.Path, "$HOME", homeDir)

		targetPath := filepath.Join(parsedBasePath, item.File)

		f, err := os.Create(targetPath)

		if err != nil {
			return err
		}

		io.Copy(f, resp.Body)

		resp.Body.Close()
		f.Close()

		fmt.Printf("Successfuly got %s config\n", item.Name)

		return nil
	}

}
