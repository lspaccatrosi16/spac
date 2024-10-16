package main

import (
	"github.com/lspaccatrosi16/go-cli-tools/command"
	"github.com/lspaccatrosi16/spac/lib/config"
	"github.com/lspaccatrosi16/spac/lib/converter"
	"github.com/lspaccatrosi16/spac/lib/credmanager"
	"github.com/lspaccatrosi16/spac/lib/install"
	"github.com/lspaccatrosi16/spac/lib/setup"
)

func main() {
	manager := command.NewManager(command.ManagerConfig{Searchable: true})

	manager.Register("setup", "Basic System Setup", setup.Setup)
	manager.Register("install", "Install Programs", install.Install)
	manager.Register("config", "Update App Configs", config.Config)
	manager.Register("credential", "Configure GCT Credential Manager", credmanager.Credential)
	manager.Register("converter", "Convert Quantities", converter.Converter)
	manager.Register("refresh", "Refresh Install List", install.RefreshCache)

	for {
		res := manager.Tui()
		if res {
			break
		}
	}
}
