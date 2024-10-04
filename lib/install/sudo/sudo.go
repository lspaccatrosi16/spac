package sudo

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunSudo(str string) error {
	return RunNonSudo(fmt.Sprintf("sudo %s", str))
}

func RunNonSudo(str string) error {
	splitted := strings.Split(str, " ")
	cmdStr := splitted[0]
	args := splitted[1:]

	cmd := exec.Command(cmdStr, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		return err
	}
	return nil

}
