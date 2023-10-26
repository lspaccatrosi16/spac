package setup

import (
	"fmt"
	"os"
	"runtime"

	"github.com/lspaccatrosi16/go-cli-tools/command"
	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/spac/lib/path"
)

func Setup() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("setup command only works on linux, not %s", runtime.GOOS)
	}

	return loop()
}

func loop() error {
	manager := command.NewManager(command.ManagerConfig{Searchable: false})

	manager.Register("1", "Create Path folder ~/bin", binPath)
	manager.Register("2", "Add a directory to PATH", addCustomToPath)
	for {
		exit := manager.Tui()

		if exit {
			break
		}
		fmt.Println("\nCompleted Successfully")
	}
	return nil
}

func binPath() error {
	binDir := path.GetBin()
	err := os.MkdirAll(path.Abs(binDir), 0o755)
	if err != nil {
		return err
	}

	err = path.AddDirToPath(binDir)

	if err != nil {
		return err
	}

	return nil
}

func addCustomToPath() error {
	customDir := input.GetInput("Directory")
	err := path.AddDirToPath(customDir)
	if err != nil {
		return err
	}
	return nil
}
