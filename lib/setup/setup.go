package setup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/lspaccatrosi16/go-cli-tools/input"
)

var profileFile string

func Setup() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("setup command only works on linux, not %s", runtime.GOOS)
	}

	return loop()
}

func loop() error {
	topOpts := []input.SelectOption{
		{Name: "Create PATH folder ~/bin", Value: "b"},
		{Name: "Install AUP", Value: "a"},
		{Name: "Install Scaffold", Value: "s"},
		{Name: "Exit", Value: "e"},
	}

	div := strings.Repeat("=", 20)

	for {
		fmt.Printf("%s SETUP %s\n", div, div)
		cmd, err := input.GetSearchableSelection("Setup item", topOpts)
		if err != nil {
			return err
		}

		switch cmd {
		case "b":
			err := binPath()
			if err != nil {
				return err
			}

		case "a":
			err := aup()
			if err != nil {
				return err
			}
		case "s":
			err := scaffold()
			if err != nil {
				return err
			}
		case "e":
			return nil
		}

		fmt.Println("\nCompleted Successfully")
	}
}

func getHome() (string, error) {
	homeDir, found := os.LookupEnv("HOME")
	if !found {
		return "", fmt.Errorf("$HOME is unset")
	}

	return homeDir, nil

}

func getProfile() (string, error) {
	if profileFile != "" {
		return profileFile, nil
	}

	homeDir, err := getHome()
	if err != nil {
		return "", err
	}

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

func getBinPath() (string, error) {
	homeDir, err := getHome()

	if err != nil {
		return "", err
	}

	binDir := filepath.Join(homeDir, "bin")
	return binDir, nil

}

func binPath() error {
	binDir, err := getBinPath()
	if err != nil {
		return err
	}

	err = os.MkdirAll(binDir, 0o755)
	if err != nil {
		return err
	}

	pathContent, _ := os.LookupEnv("PATH")
	pathArrs := strings.Split(pathContent, ":")

	duplicate := false

	for _, path := range pathArrs {
		if path == binDir {
			duplicate = true
		}
	}

	if !duplicate {
		profile, err := getProfile()
		if err != nil {
			return err
		}

		f, err := os.Open(profile)
		if err != nil {
			if os.IsNotExist(err) {
				f, err = os.Create(profile)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		defer f.Close()
		fmt.Fprintf(f, "\nexport PATH=\"$PATH:%s\"", binDir)
	}
	return nil
}

func aup() error {
	assetLink := "https://github.com/lspaccatrosi16/aup/releases/latest/download/aup-linux"
	resp, err := http.Get(assetLink)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("Got aup binary of size %.2f MiB\n", float32(resp.ContentLength)/(1024*1024))
	binPath, err := getBinPath()
	if err != nil {
		return err
	}

	targetPath := filepath.Join(binPath, "aup")

	f, err := os.Create(targetPath)

	if err != nil {
		return err
	}

	defer f.Close()
	io.Copy(f, resp.Body)

	err = os.Chmod(targetPath, 0o755)
	if err != nil {
		return err
	}
	return nil
}

func scaffold() error {
	cmd := exec.Command("aup", "-r", "lspaccatrosi16/scaffold", "-a", "scaffold-linux", "-b", "scaffold", "add")

	out, err := cmd.Output()
	fmt.Println(string(out))
	if err != nil {
		return err
	}

	return nil

}