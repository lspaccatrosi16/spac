package install

import (
	"fmt"
	"runtime"

	"github.com/lspaccatrosi16/go-cli-tools/input"
)

type pkg struct {
	APT string
	DNF string
}

type installable struct {
	Name    string
	PkgName *pkg
	Special bool
}

func Install() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("install command only works on linux, not %s", runtime.GOOS)
	}
	return loop()
}

func list() *[]installable {
	list := []installable{
		{"nvm", nil, true},
		{"gcc", &pkg{"gcc", "gcc"}, false},
	}
	return &list
}

func loop() error {
	list := list()
	options := []input.SelectOption{}
	for _, item := range *list {
		options = append(options, input.SelectOption{Name: item.Name, Value: item.Name})
	}

	for {
		_, idx, err := input.GetSearchableSelectionIdx("Select Package", options)
		if err != nil {
			return err
		}

		selected := (*list)[idx]

		if selected.Special {
			handleSpecial(selected)
		} else {
			handleNormal(selected)
		}

	}

}

func handleNormal(target installable) error {
	pkgManOpt := []input.SelectOption{{Name: "apt", Value: "a"}, {Name: "dnf", Value: "d"}}
	selected, err := input.GetSelection("Package Manager", pkgManOpt)
	if err != nil {
		return err
	}

	switch selected {

	case "a":
		return installApt(target)
	case "d":
		return installDnf(target)
	}

	return nil
}

func handleSpecial(target installable) error {

	return nil

}

func installDnf(target installable) error {
	return nil

}

func installApt(target installable) error {
	return nil

}
