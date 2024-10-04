package config

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/lspaccatrosi16/spac/lib/scache"
)

type configItem struct {
	AppName string
	CfgPath string
	Replace bool
}

type configList []configItem

func (c configList) Len() int {
	return len(c)
}

func (c configList) Swap(i int, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c configList) Less(i int, j int) bool {
	return c[i].AppName < c[j].AppName
}

func list() (configList, error) {
	cCache, err := scache.GetCachedFile("configlist.cache")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if cCache.IsValid(scache.CacheExpiry) {
		buf = *bytes.NewBuffer(cCache.Data)
	} else {
		resp, err := http.Get("https://raw.githubusercontent.com/lspaccatrosi16/config/master/manifest.json")

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()
		io.Copy(&buf, resp.Body)
	}

	var manifest configList
	dec := json.NewDecoder(&buf)
	err = dec.Decode(&manifest)
	if err != nil {
		return nil, err
	}
	return manifest, nil
}
