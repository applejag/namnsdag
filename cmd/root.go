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

// Package cmd is the command-line definitions of all commands.
package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jilleJr/namnsdag/v3/pkg/namnsdag"
	"github.com/spf13/cobra"
)

var (
	colorPrefix = color.New(color.FgHiBlack)
	colorText   = color.New(color.FgYellow)
	colorStatus = color.New(color.FgHiBlack, color.Italic)
	colorError  = color.New(color.FgRed)

	colorNameOfficial         = color.New(color.FgHiCyan)
	colorNameUnofficial       = color.New(color.FgCyan, color.Italic)
	colorNameUnofficialSymbol = color.New(color.FgMagenta, color.Italic)
	colorNameDelimiter        = color.New(color.FgHiBlack)
	colorNameNone             = color.New(color.FgRed, color.Italic)

	rootFlags = struct {
		noFetch      bool
		noCache      bool
		noUnofficial bool
	}{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "namnsdag [YYYY-MM-DD]",
	Short: "Simple CLI for fetching the list of names to celebrate today",
	Long: `Simple CLI for fetching the list of names to celebrate today.

When run, it will query https://www.dagensnamnsdag.nu/ to obtain today's names,
and cache the results inside ~/.cache/namnsdag/`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		day := time.Now()
		if len(args) == 1 {
			var err error
			day, err = time.Parse(time.DateOnly, args[0])
			if err != nil {
				return fmt.Errorf("parse argument: %w", err)
			}
		}
		namesPerDay, err := loadOrFetchNames()
		if err != nil {
			if namesPerDay != nil {
				colorStatus.Println("Found cached names, but they might be outdated.")
				writeNames(namesForToday(namesPerDay, day), day)
			}
			writeError(err)
			os.Exit(1)
			return nil
		}
		writeNames(namesForToday(namesPerDay, day), day)
		return nil
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func writeError(err error) {
	colorPrefix.Print("Error: ")
	colorError.Println(err)
}

func namesForToday(namesPerDay map[namnsdag.DoM][]namnsdag.Name, today time.Time) []namnsdag.Name {
	dom := namnsdag.NewDoMFromTime(today)
	names := namesPerDay[dom]
	if rootFlags.noUnofficial {
		names = filterOnlyOfficial(names)
	}
	return names
}

func writeNames(names []namnsdag.Name, day time.Time) {
	prefix := "Today's names"
	if !sameDate(day, time.Now()) {
		prefix = fmt.Sprintf("Names for %s", day.Format(time.DateOnly))
	}

	if len(names) == 0 {
		writeColored(fmt.Sprintf("%s: %s", prefix, colorNameNone.Sprint("no names found for today")))
		return
	}
	writeColored(fmt.Sprintf("%s: %s", prefix, joinNames(names)))
}

func sameDate(a, b time.Time) bool {
	y1, m1, d1 := a.Date()
	y2, m2, d2 := b.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func writeColored(text string) {
	var sb strings.Builder
	colorPrefix.Fprint(&sb, "===")
	sb.WriteByte(' ')
	colorText.Fprint(&sb, text)
	fmt.Println(sb.String())
}

func joinNames(names []namnsdag.Name) string {
	var sb strings.Builder
	for i, name := range names {
		if i > 0 {
			colorNameDelimiter.Fprint(&sb, ", ")
		}
		if name.TypeOfName != namnsdag.TypeUnofficial {
			colorNameOfficial.Fprint(&sb, name.Name)
		} else {
			colorNameUnofficial.Fprint(&sb, name.Name)
			colorNameUnofficialSymbol.Fprint(&sb, "*")
		}
	}
	return sb.String()
}

func loadOrFetchNames() (map[namnsdag.DoM][]namnsdag.Name, error) {
	if rootFlags.noCache && rootFlags.noFetch {
		return nil, errors.New("cannot use --no-cache and --no-fetch at the same time")
	}

	var cache namnsdag.Cache

	if !rootFlags.noCache {
		c, err := namnsdag.LoadCache()
		if err != nil {
			return nil, fmt.Errorf("load cached names: %w", err)
		}
		cache = c
	}

	isCacheValid := len(cache.NamesPerDay) > 0
	if isCacheValid && rootFlags.noFetch {
		return cache.NamesPerDay, nil
	}

	isCacheOutdated := !isCacheValid || cache.UpdatedAt.Before(time.Now().Truncate(24*time.Hour))
	if isCacheOutdated && rootFlags.noFetch {
		return nil, errors.New("none or outdated cache, and skipping fetch because --no-fetch was supplied")
	}

	if !isCacheOutdated {
		return cache.NamesPerDay, nil
	}

	req := namnsdag.Request{ETag: cache.ETag}
	if !isCacheValid {
		req.ETag = ""
	}

	colorStatus.Printf("Fetching names from %s... ", namnsdag.URL)
	resp, err := namnsdag.Fetch(req)
	if errors.Is(err, namnsdag.ErrHTTPNotModified) && isCacheValid {
		colorStatus.Println("cache is up-to-date")
		return cache.NamesPerDay, nil
	}
	if err != nil {
		colorError.Println("error")
		return cache.NamesPerDay, fmt.Errorf("fetch names: %w", err)
	}
	colorStatus.Printf("fetched %d names\n", len(resp.Names))
	cache.SetNames(resp.Names)
	cache.UpdatedAt = time.Now()
	cache.ETag = resp.ETag
	if err := namnsdag.SaveCache(cache); err != nil {
		return cache.NamesPerDay, fmt.Errorf("cache names: %w", err)
	}
	return cache.NamesPerDay, nil
}

func filterOnlyOfficial(names []namnsdag.Name) []namnsdag.Name {
	var filtered []namnsdag.Name
	for _, name := range names {
		if name.TypeOfName != namnsdag.TypeUnofficial {
			filtered = append(filtered, name)
		}
	}
	return filtered
}

// Execute is the entry point for running this command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		writeError(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&rootFlags.noFetch, "no-fetch", false, "Skips fetching via HTTP.")
	rootCmd.Flags().BoolVar(&rootFlags.noCache, "no-cache", false, "Skips loading from cache.")
	rootCmd.Flags().BoolVar(&rootFlags.noUnofficial, "no-unofficial", false, `Skips showing unofficial namnsdagar, aka "Bolibompa namnsdagar".`)
}
