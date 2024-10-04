package config

import (
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/command"
	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/go-libs/fs"
)

func Config() error {
	manager := command.NewManager(command.ManagerConfig{Searchable: true})
	items, err := list()
	if err != nil {
		return err
	}

	sort.Sort(items)

	for _, item := range items {
		manager.Register(item.AppName, fmt.Sprintf("update %s configuration", item.AppName), UseConfig(item))
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
		url := fmt.Sprintf("https://raw.githubusercontent.com/lspaccatrosi16/config/master/%s.zip", item.AppName)
		fmt.Printf("Downloading %s config from lspaccatrosi16/config\n", item.AppName)

		resp, err := http.Get(url)

		if err != nil {
			return err
		}

		homeDir, found := os.LookupEnv("HOME")
		if !found {
			return fmt.Errorf("could not resolve HOME env var")
		}

		parsedPath := strings.ReplaceAll(item.CfgPath, "$HOME", homeDir)
		if item.Replace {
			err = os.RemoveAll(parsedPath)
			if err != nil {
				return err
			}
		}

		defPath, err := input.GetConfirmSelection("Use the default config path?")
		if err != nil {
			return err
		}

		if !defPath {
			parsedPath = input.GetInput("Output path")
		}

		fmt.Printf("Using path: %s\n", parsedPath)

		err = fs.Decompress(resp.Body, parsedPath, fs.Zip)
		if err != nil {
			return err
		}

		fmt.Printf("Successfuly got %s config\n", item.AppName)
		return nil
	}

}
