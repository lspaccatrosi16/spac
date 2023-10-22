package sudo

import (
	"os"
	"os/exec"
	"strings"
)

func RunSudo(str string) error {
	splitted := strings.Split(str, " ")
	cmd := exec.Command("sudo", splitted...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil
}
