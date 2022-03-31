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
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type cache struct {
	Day   string   `json:"day"`
	Names []string `json:"names"`
}

func LoadCache(today time.Time) ([]string, error) {
	path, err := CacheFile()
	if err != nil {
		return nil, fmt.Errorf("get cache file path: %w", err)
	}
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer file.Close()
	dec := json.NewDecoder(file)
	var cache cache
	if err := dec.Decode(&cache); err != nil {
		return nil, err
	}
	if cache.Day != today.Format("2006-01-02") {
		// Cache is out of date
		return nil, nil
	}
	return cache.Names, nil
}

func SaveCache(today time.Time, names []string) error {
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

func CacheFile() (string, error) {
	dir, err := CacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "latest.json"), nil
}

func CacheDir() (string, error) {
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
