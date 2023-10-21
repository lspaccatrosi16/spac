package main

import (
	"github.com/lspaccatrosi16/spac/lib/setup"
	"github.com/lspaccatrosi16/spac/lib/types"
)

func main() {
	manager := types.NewManager()

	manager.Register("setup", "Set up my system", setup.Setup)

	manager.Gui()
}

// what do I want this to do?

// setup environment
