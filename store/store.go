package store

import (
	"errors"
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
	hasExpiry bool
}

var (
	storeSET              = map[string]entry{}
	muSET    sync.RWMutex = sync.RWMutex{}
)

func init() {
	go func() {
		for {
			time.Sleep(5 * time.Second)
			cleanupExpired()
		}
	}()
}

func SETValue(key string, value string, ttlSeconds int, hasExpiry bool) {
	muSET.Lock()
	defer muSET.Unlock()

	if hasExpiry {
		storeSET[key] = entry{
			value:     value,
			expiresAt: time.Now().Add(time.Duration(ttlSeconds) * time.Second),
			hasExpiry: true,
		}
	} else {
		storeSET[key] = entry{
			value:     value,
			hasExpiry: false,
		}
	}
}

func GETValue(key string) (string, error) {
	muSET.RLock()
	entry, ok := storeSET[key]
	muSET.RUnlock()

	if !ok {
		return "", errors.New("invalid key")
	}
	if time.Now().After(entry.expiresAt) && entry.hasExpiry {
		muSET.Lock()
		defer muSET.Unlock()

		delete(storeSET, key)
		return "", errors.New("key expired")
	}

	return entry.value, nil
}

func DELValue(key string) bool {
	muSET.Lock()
	defer muSET.Unlock()

	if entry, ok := storeSET[key]; ok {
		if time.Now().After(entry.expiresAt) && entry.hasExpiry {
			delete(storeSET, key)
			return false
		}
		delete(storeSET, key)
		return true
	}

	return false
}

func Exists(key string) bool {
	muSET.RLock()
	defer muSET.RUnlock()

	entry, ok := storeSET[key]
	if time.Now().After(entry.expiresAt) && entry.hasExpiry {
		return false
	}

	return ok
}

func cleanupExpired() {
	muSET.Lock()
	defer muSET.Unlock()

	for k, e := range storeSET {
		if e.hasExpiry && time.Now().After(e.expiresAt) {
			delete(storeSET, k)
		}
	}
}
