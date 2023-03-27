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
	"encoding"
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

// Cache is the model representing the cached data.
type Cache struct {
	ETag        string         `json:"etag"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	NamesPerDay map[DoM][]Name `json:"namesPerDay"`
}

// SetNames replaces the names of the map.
func (c *Cache) SetNames(names []Name) {
	c.NamesPerDay = nil
	c.AddNames(names)
}

// AddNames adds names to the map of names, on their appropriate dates.
func (c *Cache) AddNames(names []Name) {
	if c.NamesPerDay == nil {
		c.NamesPerDay = make(map[DoM][]Name, len(names))
	}
	for _, name := range names {
		dom := NewDoM(name.Month, name.Day)
		c.NamesPerDay[dom] = append(c.NamesPerDay[dom], name)
	}
}

// DoM (Day-of-Month) represents a day in a month, no matter what year.
type DoM struct {
	Day   int
	Month time.Month
}

var _ encoding.TextMarshaler = DoM{}
var _ encoding.TextUnmarshaler = &DoM{}

// MarshalText implements [encoding.TextMarshaler]
func (d DoM) MarshalText() (text []byte, err error) {
	return fmt.Appendf(nil, "%02d-%02d", d.Month, d.Day), nil
}

// UnmarshalText implements [encoding.TextUnmarshaler]
func (d *DoM) UnmarshalText(text []byte) error {
	_, err := fmt.Sscanf(string(text), "%02d-%02d", &d.Month, &d.Day)
	return err
}

// String implements [fmt.Stringer]
func (d DoM) String() string {
	b, _ := d.MarshalText()
	return string(b)
}

// NewDoMFromTime creates a new [DoM] based on the month and day in the
// given time. The year, as well as any hours, minutes, seconds, milliseconds,
// and time zone is ignored.
func NewDoMFromTime(t time.Time) DoM {
	_, month, day := t.Date()
	return NewDoM(month, day)
}

// NewDoM creates a new [DoM] based on the month and the day.
func NewDoM(month time.Month, day int) DoM {
	return DoM{
		Day:   day,
		Month: month,
	}
}

// LoadCache loads the cached names from ~/.cache/namnsdag/latest.json, or the
// equivalent in other OS's cache directories (eg. %LOCALAPPDATA%).
//
// It will return nil if there is no cache or if the cache is outdated.
func LoadCache() (Cache, error) {
	path, err := CacheFile()
	if err != nil {
		return Cache{}, fmt.Errorf("get cache file path: %w", err)
	}
	fileBytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Cache{}, nil
	} else if err != nil {
		return Cache{}, err
	}
	var cache Cache
	if err := json.Unmarshal(fileBytes, &cache); err != nil {
		return Cache{}, err
	}
	return cache, nil
}

// SaveCache writes the cached names to ~/.cache/namnsdag/latest.json, or the
// equivalent in other OS's cache directories (eg. %LOCALAPPDATA%).
//
// Today's year, month, and day are used to automatically detect the cache as
// outdated when loading the cached names.
func SaveCache(cache Cache) error {
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

	if cache.UpdatedAt == (time.Time{}) {
		cache.UpdatedAt = time.Now()
	}

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(cache)
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
	return filepath.Join(dir, "cache@v3.json"), nil
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
