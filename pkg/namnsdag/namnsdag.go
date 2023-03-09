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
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// URL is the HTTP URL that namnsdag.Fetch will query.
const URL = "https://dagensnamnsdag.nu/"

// Fetch performs a HTTP GET request and parses the HTML response to extract
// today's names.
func Fetch() ([]string, error) {
	doc, err := fetchDocument()
	if err != nil {
		return nil, err
	}
	var names []string
	doc.Find(".container p").Each(func(i int, s *goquery.Selection) {
		class, ok := s.Attr("class")
		if !ok || !strings.HasPrefix(class, "index_todaysNames") {
			return
		}
		names = append(names, getNames(s)...)
	})
	sort.Strings(names)
	return names, nil
}

func getNames(s *goquery.Selection) []string {
	var names []string
	s.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		names = append(names, s.Text())
	})
	return names
}

func fetchDocument() (*goquery.Document, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("non-2xx status code: %s", resp.Status)
	}
	return goquery.NewDocumentFromReader(resp.Body)
}
