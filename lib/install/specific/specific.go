package specific

import (
	"fmt"

	"github.com/lspaccatrosi16/spac/lib/install/specific/golang"
	"github.com/lspaccatrosi16/spac/lib/install/specific/script"
)

type specific struct {
	Name string
	Exec func() error
}

type specificList struct {
	Items []specific
}

func (l *specificList) Register(name string, exec func() error) {
	s := specific{
		Name: name,
		Exec: exec,
	}

	l.Items = append(l.Items, s)
}

func (l *specificList) Get(name string) *(func() error) {
	for i := 0; i < len(l.Items); i++ {
		s := l.Items[i]
		if s.Name == name {
			return &s.Exec
		}
	}

	return nil
}

func FindSpecific(name string) (*(func() error), error) {
	list := specificList{}

	list.Register("nvm", script.UseScript("https://github.com/nvm-sh/nvm/blob/master/install.sh", "nvm"))
	list.Register("starship", script.UseScript("https://starship.rs/install.sh", "starship"))

	list.Register("go", golang.Install)

	exec := list.Get(name)
	if exec == nil {
		return nil, fmt.Errorf("could not find an handler command for \"%s\"", name)
	}

	return exec, nil

}
