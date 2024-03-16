package install

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/go-cli-tools/cache"
	"github.com/lspaccatrosi16/go-cli-tools/config"
	"github.com/lspaccatrosi16/spac/lib/install/sudo"
)

var configDir, _ = config.GetConfigPath("spac")

const cacheExpiry = 60 * 60 * 24

type aupinstallable struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	RepoKey      string   `json:"repoKey"`
	ArtifactName string   `json:"artifactName"`
	BinaryName   string   `json:"binaryName"`
	Category     []string `json:"category"`
}

type packageinstallable struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Apt         string   `json:"apt"`
	Dnf         string   `json:"dnf"`
	Flatpak     string   `json:"flatpak"`
	Postinstall string   `json:"postinstall"`
	Category    []string `json:"category"`
}

func specificList() installList {
	return installList{
		{"nvm", "NodeJS version manager", nil, true, nil, []string{"CLI Tool", "Development"}},
		{"go", "Golang compiler", nil, true, nil, []string{"CLI Tool", "Development", "Compiler"}},
		{"starship", "Shell prompt", nil, true, nil, []string{}},
		{"aup", "package manager", nil, true, nil, []string{"CLI Tool", "Package Manager"}},
	}
}

func getCachedFile(name string) (*cache.CacheItem, error) {
	item := new(cache.CacheItem)
	path := filepath.Join(configDir, name)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return item, nil
		}
		return nil, err
	}

	defer f.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	err = item.Decode(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return item, nil
}

func writeCachedFile(name string, item *cache.CacheItem) error {
	path := filepath.Join(configDir, name)
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	data, err := item.Encode()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(data)
	io.Copy(f, buf)

	return nil

}

func aupPackages() (installList, error) {
	acache, err := getCachedFile("aupmanifest.cache")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if acache.IsValid(cacheExpiry) {
		buf = *bytes.NewBuffer(acache.Data)
	} else {
		resp, err := http.Get("https://raw.githubusercontent.com/lspaccatrosi16/spacmanifest/master/aupmanifest.json")
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		io.Copy(&buf, resp.Body)
		writeCachedFile("aupmanifest.cache", cache.CreateCacheItem(buf.Bytes()))
	}

	var manifest []aupinstallable

	err = json.Unmarshal(buf.Bytes(), &manifest)
	if err != nil {
		return nil, err
	}

	list := installList{}

	for _, item := range manifest {
		list = append(list, installable{
			Name:        item.Name,
			Description: item.Description,
			PkgName: &pkg{
				AUP: &aupPkg{
					Artifact: item.ArtifactName,
					Binary:   item.BinaryName,
					Repokey:  item.RepoKey,
				},
			},
			Category: item.Category,
		})

	}
	return list, nil

}

func managerPackages() (installList, error) {
	mcache, err := getCachedFile("packagemanifest.cache")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if mcache.IsValid(cacheExpiry) {
		buf = *bytes.NewBuffer(mcache.Data)
	} else {
		resp, err := http.Get("https://raw.githubusercontent.com/lspaccatrosi16/spacmanifest/master/packagemanifest.json")
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		io.Copy(&buf, resp.Body)
		writeCachedFile("packagemanifest.cache", cache.CreateCacheItem(buf.Bytes()))
	}

	var manifest []packageinstallable

	err = json.Unmarshal(buf.Bytes(), &manifest)
	if err != nil {
		return nil, err
	}

	list := installList{}

	for _, item := range manifest {
		var cfunc func() error

		if item.Postinstall != "" {
			cfunc = func() error {
				return sudo.RunSudo(item.Postinstall)
			}
		}

		list = append(list, installable{
			Name:        item.Name,
			Description: item.Description,
			PkgName: &pkg{
				APT:     item.Apt,
				DNF:     item.Dnf,
				FLATPAK: item.Dnf,
			},
			Config:   cfunc,
			Category: item.Category,
		})

	}
	return list, nil
}

var listchache *installList

func list() (*installList, error) {
	if listchache != nil {
		return listchache, nil
	}
	list := specificList()

	aupManifest, err := aupPackages()

	if err != nil {
		return nil, err
	}

	packageManifest, err := managerPackages()

	if err != nil {
		return nil, err
	}

	list = append(list, aupManifest...)
	list = append(list, packageManifest...)

	listchache = &list
	return &list, nil
}
