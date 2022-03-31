<!--
SPDX-FileCopyrightText: 2022 Kalle Fagerberg

SPDX-License-Identifier: CC-BY-4.0
-->

# Contributing

## How to contribute

There's multiple ways to contribute:

- Report bugs or wanted features to our issue tracker:
  <https://github.com/go-typ/typ/issues/new>

- Tackle an issue in <https://github.com/go-typ/typ/issues>. If you see one
  that's tempting, ask in the issue's thread and I'll assign you so we don't get
  multiple people working on the same thing.

## Development

### Prerequisites

[Go](https://go.dev/) v1.18 (or higher)

### Formatting

- Go formatting is performed via [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports)
  and can be installed via `make deps`

- Markdown formatting can be (but doesn't have to be) done via [Prettier](https://prettier.io/)

Make sure to regularly lint the project locally. We're relying on the Markdown
linting for its formatting.

### Linting

- Go linting is performed via [Revive](https://revive.sh)

- Markdown linting is performed via [remarklint](https://github.com/remarkjs/remark-lint)
  and requires NPM.

- Licensing linting is performed via [REUSE](https://reuse.software/) and
  requires Python v3.

You can install all the above linting rules by running `make deps`

```sh
# Lint .go & .md files, and REUSE compliance:
make lint

# Apply fixes, where possible, to .md files:
make lint-fix

# Only lint some:
make lint-md
make lint-go
make lint-license

# Only lint & fix some:
make lint-md-fix
```
