package converter

import (
	"fmt"

	"github.com/lspaccatrosi16/go-cli-tools/command"
	"github.com/lspaccatrosi16/go-cli-tools/input"
)

func Converter() error {
	manager := command.NewManager(command.ManagerConfig{Searchable: true})

	items := list()

	for _, c := range items {
		convertFunc := func() error {
			inputOpts := optsFromArr(c.InputFormat())
			outputOpts := optsFromArr(c.OutputFormat())

			inputSel, err := input.GetSearchableSelection("Select Input Format", inputOpts)
			if err != nil {
				return err
			}

			outputSel, err := input.GetSearchableSelection("Select Output Format", outputOpts)
			if err != nil {
				return err
			}

			val := c.Convert(inputSel, outputSel)
			fmt.Println(val)
			return nil
		}

		manager.Register(c.Name(), "", convertFunc)
	}

	manager.Tui()
	return nil
}

func optsFromArr(arr []string) []input.SelectOption {
	opts := []input.SelectOption{}

	for _, i := range arr {
		opts = append(opts, input.SelectOption{
			Name:  i,
			Value: i,
		})
	}

	return opts
}
