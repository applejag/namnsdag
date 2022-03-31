<!--
SPDX-FileCopyrightText: 2022 Kalle Fagerberg

SPDX-License-Identifier: CC-BY-4.0
-->

# namnsdag

[![REUSE status](https://api.reuse.software/badge/github.com/jilleJr/namnsdag)](https://api.reuse.software/info/github.com/jilleJr/namnsdag)

```console
$ namnsdag
Fetching names from https://www.dagensnamnsdag.nu/
=== Today's names: Erla*, Essy*, Ester, Kenji*, Lenore*, Scilla*

$ namnsdag --help
Simple CLI for fetching the list of names to celebrate today.

When run, it will query https://www.dagensnamnsdag.nu/ to obtain today's names,
and cache the results inside ~/.cache/namnsdag/

Usage:
  namnsdag [flags]

Flags:
  -h, --help   help for namnsdag
      --no-cache   Skips loading from cache.
      --no-fetch   Skips fetching via HTTP.
```

## Install

Requires Go 1.18 or higher.

```sh
go install github.com/jilleJr/namnsdag@latest
```

## License

This project is primarily licensed under the GNU General Public License v3.0 or
later (GPL-3.0-or-later).

Copyright &copy; Kalle Fagerberg
