package scache

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/lspaccatrosi16/go-cli-tools/cache"
	"github.com/lspaccatrosi16/go-cli-tools/config"
)

var configDir, _ = config.GetConfigPath("spac")

const CacheExpiry = 60 * 60 * 24

func GetCachedFile(name string) (*cache.CacheItem, error) {
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

func WriteCachedFile(name string, item *cache.CacheItem) error {
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

func RemoveCachedFile(name string) error {
	path := filepath.Join(configDir, name)
	return os.Remove(path)
}
