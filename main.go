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

// Package main is the entrypoint for this module. It launches the CLI from the
// cmd package.
package main

import "github.com/jilleJr/namnsdag/v3/cmd"

func main() {
	cmd.Execute()
}
