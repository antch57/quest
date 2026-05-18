package jambase

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/antch57/quest/internal/store"
)

const cacheFile = "cache.json"
const cacheDir = "jamz"

var ErrorCacheNotFound = errors.New("cache does not exist")

func saveCache(key, value string) error {
	path, err := store.StorePath(cacheDir, cacheFile)
	if err != nil {
		return err
	}

	cacheData := map[string]string{
		key: value,
	}
	cacheBytes, err := json.Marshal(cacheData)
	if err != nil {
		return err
	}
	existingCache, err := loadCache()
	if err != nil {
		if !errors.Is(err, ErrorCacheNotFound) {
			return err
		}
	}

	for k, v := range existingCache {
		if k != key {
			cacheData[k] = v.(string)
		}
	}

	if err := os.WriteFile(path, cacheBytes, 0o600); err != nil {
		return err
	}
	return nil
}

func loadCache() (map[string]interface{}, error) {
	path, err := store.StorePath(cacheDir, cacheFile)
	if err != nil {
		return nil, err
	}

	store, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, ErrorCacheNotFound
	}
	if err != nil {
		return nil, err
	}
	var cacheData map[string]interface{}
	if err := json.Unmarshal(store, &cacheData); err != nil {
		return nil, err
	}
	return cacheData, nil

}

func checkCache(key string) (string, error) {
	path, err := store.StorePath(cacheDir, cacheFile)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", ErrorCacheNotFound
	}

	cacheData, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var cacheMap map[string]interface{}
	if err := json.Unmarshal(cacheData, &cacheMap); err != nil {
		return "", err
	}
	for k, v := range cacheMap {
		if k == key {
			return v.(string), nil
		}
	}
	return "", ErrorCacheNotFound
}
