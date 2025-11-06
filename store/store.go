package store

import (
	"errors"
	"fmt"
	"sync"
)

var (
	storeSET map[string]string = map[string]string{}
	muSET    sync.RWMutex      = sync.RWMutex{}
)

func SETValue(key string, value string) {
	muSET.Lock()
	defer muSET.Unlock()

	fmt.Println("key", key, "Vaalue", value)

	storeSET[key] = value
}

func GETValue(key string) (string, error) {
	muSET.RLock()
	defer muSET.RUnlock()

	value, ok := storeSET[key]
	if !ok {
		return "", errors.New("invalid key")
	}

	return value, nil
}

func DELValue(key string) bool {
	muSET.Lock()
	defer muSET.Unlock()

	if _, ok := storeSET[key]; ok {
		delete(storeSET, key)
		return true
	}

	return false
}

func Exists(key string) bool {
	muSET.RLock()
	defer muSET.RUnlock()

	_, ok := storeSET[key]

	return ok
}
