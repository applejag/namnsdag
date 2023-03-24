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

// Package namnsdag contains functions to programatically retrieve today's names,
// as well as caching them.
package namnsdag

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// URL is the HTTP URL that namnsdag.Fetch will query.
const URL = "https://dagensnamnsdag.nu/"

// Name contains fields for a given name.
type Name struct {
	URL        string     `json:"url"`
	Name       string     `json:"name"`
	Day        int        `json:"day"`
	Month      time.Month `json:"month"`
	TypeOfName Type       `json:"typeOfName"`
	Gender     Gender     `json:"gender"`
}

// Type is an enum stating what kind of namnsdag-name it is.
type Type string

// Known values for [Type]. There may be other values from
// [https://dagensnamnsdag.nu], but these are the ones found so far.
const (
	TypeName    Type = "NAME"
	TypeNewName Type = "NEW_NAME"
)

// Gender is an enum stating what gender a namnsdag-name has, if any.
type Gender string

// Known values for [Gender]. There may be other values from
// [https://dagensnamnsdag.nu], but these are the ones found so far.
const (
	GenderBoth   Gender = "BOTH"
	GenderBoy    Gender = "BOY"
	GenderGirl   Gender = "GIRL"
	GenderNotSet Gender = "NOT_SET"
)

// FetchToday performs a HTTP GET request and parses the HTML response
// to extract today's names.
func FetchToday() ([]Name, error) {
	data, err := fetchTodayNextJSData()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	todaysDay := now.Day()
	todaysMonth := now.Month()

	var names []Name
	for _, name := range data.Props.PageProps.Names {
		if name.Day == todaysDay && name.Month == todaysMonth {
			names = append(names, name)
		}
	}
	sort.Slice(names, func(i, j int) bool {
		return names[i].Name < names[j].Name
	})
	return names, nil
}

type nextJSData struct {
	Props struct {
		PageProps struct {
			Names []Name `json:"names"`
		} `json:"pageProps"`
	} `json:"props"`
}

func fetchTodayNextJSData() (*nextJSData, error) {
	doc, err := fetchDocument()
	if err != nil {
		return nil, err
	}
	q := doc.Find(`script[id="__NEXT_DATA__"]`).First()
	if len(q.Nodes) == 0 {
		return nil, fmt.Errorf("no <script id='__NEXT_DATA__'> tag found")
	}
	var data nextJSData
	if err := json.Unmarshal([]byte(q.Text()), &data); err != nil {
		return nil, fmt.Errorf("parsing JSON in <script id='__NEXT_DATA__'> tag: %w", err)
	}
	return &data, nil
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
