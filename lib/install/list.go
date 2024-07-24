package install

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/lspaccatrosi16/go-cli-tools/cache"
	"github.com/lspaccatrosi16/spac/lib/install/sudo"
	"github.com/lspaccatrosi16/spac/lib/scache"
)

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

func aupPackages() (installList, error) {
	acache, err := scache.GetCachedFile("aupmanifest.cache")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if acache.IsValid(scache.CacheExpiry) {
		buf = *bytes.NewBuffer(acache.Data)
	} else {
		resp, err := http.Get("https://raw.githubusercontent.com/lspaccatrosi16/spacmanifest/master/aupmanifest.json")
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		io.Copy(&buf, resp.Body)
		scache.WriteCachedFile("aupmanifest.cache", cache.CreateCacheItem(buf.Bytes()))
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
	mcache, err := scache.GetCachedFile("packagemanifest.cache")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if mcache.IsValid(scache.CacheExpiry) {
		buf = *bytes.NewBuffer(mcache.Data)
	} else {
		resp, err := http.Get("https://raw.githubusercontent.com/lspaccatrosi16/spacmanifest/master/packagemanifest.json")
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		io.Copy(&buf, resp.Body)
		scache.WriteCachedFile("packagemanifest.cache", cache.CreateCacheItem(buf.Bytes()))
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
