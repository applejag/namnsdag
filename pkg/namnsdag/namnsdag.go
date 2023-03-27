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
	"errors"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	// URL is the HTTP URL of the website to find data from.
	URL = "https://dagensnamnsdag.nu/namnsdagar"

	// ErrHTTPNotModified is returned from [Fetch] when the server responded
	// with status "304 not modified", which means that the etag matched
	// and our local cache is up to date.
	ErrHTTPNotModified = errors.New("http status: 304 not modified")
)

// Name contains fields for a given name.
type Name struct {
	URL        string     `json:"url"`
	Name       string     `json:"name"`
	Day        int        `json:"day"`
	Month      time.Month `json:"month"`
	TypeOfName Type       `json:"typeOfName"`
	Gender     Gender     `json:"gender"`
}

// DoM returns this name's Day-of-Month.
func (n Name) DoM() DoM {
	return NewDoM(n.Month, n.Day)
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

// Request is the model used for a [Fetch] of names from [URL].
type Request struct {
	ETag string
}

// Response is the data received from a [Fetch] of names from [URL].
type Response struct {
	Names []Name
	ETag  string
}

// Fetch performs a HTTP GET request and parses the HTML response
// to extract all names.
func Fetch(req Request) (Response, error) {
	data, etag, err := fetchAllNextJSData(req.ETag)
	if errors.Is(err, ErrHTTPNotModified) {
		return Response{ETag: etag}, err
	}
	if err != nil {
		return Response{}, err
	}
	names := data.Props.PageProps.Names
	SortNames(names)
	return Response{
		Names: names,
		ETag:  etag,
	}, nil
}

// SortNames will sort a slice of names first by month, then by day, and finally
// by name, all in ascending order.
func SortNames(names []Name) {
	sort.Slice(names, func(i, j int) bool {
		diffMonths := names[i].Month != names[j].Month
		if diffMonths {
			return names[i].Month < names[j].Month
		}
		diffDays := names[i].Day != names[j].Day
		if diffDays {
			return names[i].Day < names[j].Day
		}
		return names[i].Name < names[j].Name
	})
}

type nextJSData struct {
	Props struct {
		PageProps struct {
			Names []Name `json:"names"`
		} `json:"pageProps"`
	} `json:"props"`
}

func fetchAllNextJSData(etag string) (*nextJSData, string, error) {
	doc, newEtag, err := fetchDocument(etag)
	if errors.Is(err, ErrHTTPNotModified) {
		return nil, etag, err
	}
	if err != nil {
		return nil, "", err
	}
	q := doc.Find(`script[id="__NEXT_DATA__"]`).First()
	if len(q.Nodes) == 0 {
		return nil, "", fmt.Errorf("no <script id='__NEXT_DATA__'> tag found")
	}
	var data nextJSData
	if err := json.Unmarshal([]byte(q.Text()), &data); err != nil {
		return nil, "", fmt.Errorf("parsing JSON in <script id='__NEXT_DATA__'> tag: %w", err)
	}
	return &data, newEtag, nil
}

func fetchDocument(etag string) (*goquery.Document, string, error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, "", err
	}
	if etag != "" {
		req.Header.Add("If-None-Match", etag)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotModified {
		return nil, "", ErrHTTPNotModified
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, "", fmt.Errorf("non-2xx status code: %s", resp.Status)
	}
	q, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("parse HTML: %w", err)
	}
	return q, resp.Header.Get("etag"), nil
}
