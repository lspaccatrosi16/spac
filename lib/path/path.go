package path

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/input"
)

var profileFile string

func AddStringToPath(str string) error {
	profileLoc, err := GetProfileFile()
	if err != nil {
		return err
	}

	fmt.Printf("Opening profile file at %s\n", profileLoc)

	src, err := os.Open(profileLoc)
	contents := bytes.NewBuffer(nil)

	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		io.Copy(contents, src)
		src.Close()
	}

	fmt.Fprintln(contents, str)

	dst, err := os.Create(profileLoc)

	if err != nil {
		return err
	}

	io.Copy(dst, contents)
	dst.Close()

	fmt.Printf("Added %s to PATH\n", str)

	return nil

}

func AddDirToPath(dir string) error {
	pathContent, _ := os.LookupEnv("PATH")
	dirs := strings.Split(pathContent, ":")

	absDir := Abs(dir)

	for _, d := range dirs {
		if d == absDir {
			fmt.Printf("PATH already has %s\n", absDir)
			return nil
		}
	}

	pathStr := fmt.Sprintf("export PATH=\"$PATH:%s\"", absDir)

	err := AddStringToPath(pathStr)

	if err != nil {
		return err
	}
	return nil
}

func GetProfileFile() (string, error) {
	if profileFile != "" {
		return profileFile, nil
	}
	homeDir := GetHome()

	profLocOpts := []input.SelectOption{
		{Name: "ZSH: ~./zshrc", Value: filepath.Join(homeDir, ".zshrc")},
		{Name: "BASH: ~/.bashrc", Value: filepath.Join(homeDir, ".bashrc")},
		{Name: "Other", Value: "other"},
	}

	profLoc, err := input.GetSelection("Profile File", profLocOpts)

	if err != nil {
		return "", err
	}

	if profLoc == "other" {
		profLoc = input.GetInput("Profile File Location")
	}

	profileFile = profLoc
	return profLoc, err
}

func GetHome() string {
	homeDir, found := os.LookupEnv("HOME")
	if !found {
		panic("$HOME is unset")
	}

	return homeDir

}

func GetBin() string {
	return "~/bin"
}

func Abs(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	} else if strings.HasPrefix(path, "~/") {
		path = path[2:]

		home := GetHome()
		return filepath.Join(home, path)
	} else {
		home := GetHome()
		return filepath.Join(home, path)
	}
}
