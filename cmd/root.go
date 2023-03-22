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
	"github.com/jilleJr/namnsdag/v2/pkg/namnsdag"
	"github.com/spf13/cobra"
)

var (
	colorPrefix = color.New(color.FgHiBlack)
	colorText   = color.New(color.FgYellow)
	colorStatus = color.New(color.FgHiBlack, color.Italic)

	rootFlags = struct {
		noFetch      bool
		noCache      bool
		noUnofficial bool
	}{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "namnsdag",
	Short: "Simple CLI for fetching the list of names to celebrate today",
	Long: `Simple CLI for fetching the list of names to celebrate today.

When run, it will query https://www.dagensnamnsdag.nu/ to obtain today's names,
and cache the results inside ~/.cache/namnsdag/`,
	RunE: func(cmd *cobra.Command, args []string) error {
		names, err := loadOrFetchNames()
		if err != nil {
			return err
		}
		if rootFlags.noUnofficial {
			names = filterOnlyOfficial(names)
		}
		writeColored(fmt.Sprintf("Today's names: %s", strings.Join(names, ", ")))
		return nil
	},
}

func writeColored(text string) {
	var sb strings.Builder
	colorPrefix.Fprint(&sb, "===")
	sb.WriteByte(' ')
	colorText.Fprint(&sb, text)
	fmt.Println(sb.String())
}

func loadOrFetchNames() ([]string, error) {
	if rootFlags.noCache && rootFlags.noFetch {
		return nil, errors.New("cannot use --no-cache and --no-fetch at the same time")
	}

	today := time.Now()

	if !rootFlags.noCache {
		names, err := namnsdag.LoadCache(today)
		if err != nil {
			return nil, fmt.Errorf("load cached names: %w", err)
		}
		if names != nil {
			return names, nil
		}
		if rootFlags.noFetch {
			return nil, errors.New("none or outdated cache, and skipping fetch because --no-fetch was supplied")
		}
	}

	colorStatus.Println("Fetching names from " + namnsdag.URL)
	names, err := namnsdag.Fetch()
	if err != nil {
		return nil, fmt.Errorf("fetch names: %w", err)
	}
	if err := namnsdag.SaveCache(today, names); err != nil {
		return nil, fmt.Errorf("cache names: %w", err)
	}
	return names, nil
}

func filterOnlyOfficial(names []string) []string {
	var filtered []string
	for _, name := range names {
		if !strings.HasSuffix(name, "*") {
			filtered = append(filtered, name)
		}
	}
	return filtered
}

// Execute is the entry point for running this command.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&rootFlags.noFetch, "no-fetch", false, "Skips fetching via HTTP.")
	rootCmd.Flags().BoolVar(&rootFlags.noCache, "no-cache", false, "Skips loading from cache.")
	rootCmd.Flags().BoolVar(&rootFlags.noUnofficial, "no-unofficial", false, `Skips showing unofficial namnsdagar, aka "Bolibompa namnsdagar".`)
}
