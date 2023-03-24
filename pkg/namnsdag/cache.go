// Simple CLI for fetching the list of names to celebrate today.
// <https://github.com/jilleJr/namnsdag>
//
// SPDX-FileCopyrightText: 2022 Kalle Fagerberg
//
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

package namnsdag

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Errors specific to the cache.
var (
	ErrCacheAlreadyCleared = errors.New("cache already cleared")
)

type cache struct {
	Day   string `json:"day"`
	Names []Name `json:"names"`
}

type cachev1 struct {
	Day   string   `json:"day"`
	Names []string `json:"names"`
}

// LoadCache loads the cached names from ~/.cache/namnsdag/latest.json, or the
// equivalent in other OS's cache directories (eg. %LOCALAPPDATA%).
//
// It will return nil if there is no cache or if the cache is outdated.
func LoadCache(today time.Time) ([]Name, error) {
	path, err := CacheFile()
	if err != nil {
		return nil, fmt.Errorf("get cache file path: %w", err)
	}
	fileBytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var cache cache
	if err := json.Unmarshal(fileBytes, &cache); err != nil {
		// Maybe v1 cache format
		var cachev1 cachev1
		if errv1 := json.Unmarshal(fileBytes, &cachev1); errv1 != nil {
			// If even that failed, return original error
			return nil, err
		}
		// If cachev1 succeeded, just consider it out of date
		return nil, nil
	}
	if cache.Day != today.Format("2006-01-02") {
		// Cache is out of date
		return nil, nil
	}
	return cache.Names, nil
}

// SaveCache writes the cached names to ~/.cache/namnsdag/latest.json, or the
// equivalent in other OS's cache directories (eg. %LOCALAPPDATA%).
//
// Today's year, month, and day are used to automatically detect the cache as
// outdated when loading the cached names.
func SaveCache(today time.Time, names []Name) error {
	path, err := CacheFile()
	if err != nil {
		return fmt.Errorf("get cache file path: %w", err)
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	return enc.Encode(cache{
		Day:   today.Format("2006-01-02"),
		Names: names,
	})
}

// ClearCache will remove the cached names, if any. Returns
// ErrCacheAlreadyCleared if no cache existed.
func ClearCache() error {
	path, err := CacheFile()
	if err != nil {
		return fmt.Errorf("get cache file path: %w", err)
	}
	err = os.Remove(path)
	if os.IsNotExist(err) {
		return ErrCacheAlreadyCleared
	}
	return err
}

// CacheFile returns the path to the cache file.
func CacheFile() (string, error) {
	dir, err := cacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "latest.json"), nil
}

func cacheDir() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		dir, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(dir, ".cache")
	}
	return filepath.Join(dir, "namnsdag"), nil
}
