package install

import (
	"fmt"
	"runtime"
	"sort"

	"github.com/lspaccatrosi16/go-cli-tools/command"
	"github.com/lspaccatrosi16/go-cli-tools/input"
	"github.com/lspaccatrosi16/spac/lib/install/specific"
	"github.com/lspaccatrosi16/spac/lib/install/sudo"
)

type aupPkg struct {
	Repokey  string
	Artifact string
	Binary   string
}

type pkg struct {
	APT     string
	AUP     *aupPkg
	DNF     string
	FLATPAK string
}

type installable struct {
	Name        string
	Description string
	PkgName     *pkg
	Special     bool
	Config      func() error
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

	manager := command.NewManager(command.ManagerConfig{Searchable: true})

	for _, item := range *list {
		if item.Special {
			manager.Register(item.Name, item.Description, wrapTarget(item, handleSpecial))
		} else {
			manager.Register(item.Name, item.Description, wrapTarget(item, handleNormal))
		}
	}

	for {
		exit := manager.Tui()
		if exit {
			return nil
		}

	}
}

func wrapTarget(target installable, f func(installable) error) func() error {
	return func() error {
		err := f(target)
		if err != nil {
			return err
		}

		fmt.Printf("%s was installed successfully\n", target.Name)

		if target.Config != nil {
			err = target.Config()

			if err != nil {
				return err
			}
			fmt.Printf("%s was configured successfully\n", target.Name)
		}
		fmt.Printf("You may need to open a new terminal, or refresh your profile for the changes to take effect")
		return nil
	}
}

func handleNormal(target installable) error {
	pkgManOpt := []input.SelectOption{}

	if target.PkgName.APT != "" {
		pkgManOpt = append(pkgManOpt, input.SelectOption{Name: "apt", Value: "a"})
	}

	if target.PkgName.AUP != nil {
		pkgManOpt = append(pkgManOpt, input.SelectOption{Name: "aup", Value: "l"})
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
	case "l":
		return installAup(target)
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

func installAup(target installable) error {
	aup := target.PkgName.AUP
	cmd := fmt.Sprintf("aup -r %s -a %s -b %s add", aup.Repokey, aup.Artifact, aup.Binary)

	return sudo.RunNonSudo(cmd)
}
