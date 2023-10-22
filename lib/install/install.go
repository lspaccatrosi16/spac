package install

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/spac/lib/install/specific"
	"github.com/lspaccatrosi16/spac/lib/install/sudo"
)

type pkg struct {
	APT     string
	DNF     string
	FLATPAK string
}

type installable struct {
	Name    string
	PkgName *pkg
	Special bool
	Config  func() error
}

type installList []installable

func (b installList) Len() int {
	return len(b)
}

func (b installList) Swap(i int, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b installList) Less(i int, j int) bool {
	return b[i].Name < b[j].Name
}

func Install() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("install command only works on linux, not %s", runtime.GOOS)
	}
	return loop()
}

func loop() error {
	list := list()
	sort.Sort(list)

	options := []input.SelectOption{}
	for _, item := range *list {
		options = append(options, input.SelectOption{Name: item.Name, Value: item.Name})
	}

	options = append(options, input.SelectOption{Name: "Back", Value: "e"})

	for {
		c, idx, err := input.GetSearchableSelectionIdx("Select Package", options)
		if err != nil {
			return err
		}

		if c == "e" {
			return nil
		}

		selected := (*list)[idx]

		if selected.Special {
			err = handleSpecial(selected)
		} else {
			err = handleNormal(selected)
		}

		if err != nil {
			return err
		}

		fmt.Printf("%s was installed successfully\n", selected.Name)

		if selected.Config != nil {
			err = (selected.Config)()

			if err != nil {
				return err
			}

			fmt.Printf("%s was configured successfully\n", selected.Name)
		}

	}

}

func handleNormal(target installable) error {
	pkgManOpt := []input.SelectOption{}

	if target.PkgName.APT != "" {
		pkgManOpt = append(pkgManOpt, input.SelectOption{Name: "apt", Value: "a"})
	}

	if target.PkgName.DNF != "" {
		pkgManOpt = append(pkgManOpt, input.SelectOption{Name: "dnf", Value: "d"})
	}

	if target.PkgName.FLATPAK != "" {
		pkgManOpt = append(pkgManOpt, input.SelectOption{Name: "flatpak", Value: "f"})
	}

	selected, err := input.GetSelection("Package Manager", pkgManOpt)
	if err != nil {
		return err
	}

	switch selected {

	case "a":
		return installApt(target)
	case "d":
		return installDnf(target)
	case "f":
		return installFlatpak(target)
	}

	return nil
}

func handleSpecial(target installable) error {
	exec, err := specific.FindSpecific(target.Name)

	if err != nil {
		return err
	}

	return (*exec)()
}

func installDnf(target installable) error {
	cmd := fmt.Sprintf("dnf install %s", target.PkgName.DNF)
	return sudo.RunSudo(cmd)
}

func installApt(target installable) error {
	cmd := fmt.Sprintf("apt-get install %s", target.PkgName.APT)
	return sudo.RunSudo(cmd)
}

func installFlatpak(target installable) error {
	cmd := fmt.Sprintf("flatpak install flathub %s", target.PkgName.FLATPAK)
	return sudo.RunSudo(cmd)
}
