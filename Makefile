# SPDX-FileCopyrightText: 2022 Kalle Fagerberg
#
# SPDX-License-Identifier: CC0-1.0

namnsdag:
	go build .

.PHONY: install
install:
	go install

.PHONY: check
check:
	go test ./...

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: deps
deps: node_modules
	go install github.com/mgechev/revive@latest
	go install golang.org/x/tools/cmd/goimports@latest
	python3 -m pip install --upgrade --user reuse

node_modules:
	npm install

.PHONY: lint
lint: lint-md lint-go lint-license

.PHONY: lint-fix
lint-fix: lint-md-fix lint-go-fix

.PHONY: lint-md
lint-md: node_modules
	npx remark . .github

.PHONY: lint-md-fix
lint-md-fix: node_modules
	npx remark . .github -o

.PHONY: lint-go
lint-go:
	@echo goimports -d '**/*.go'
	@goimports -d $(shell git ls-files "*.go")
	revive -formatter stylish -config revive.toml ./...

.PHONY: lint-go-fix
lint-go-fix:
	@echo goimports -d -w '**/*.go'
	@goimports -d -w $(shell git ls-files "*.go")

.PHONY: lint-license
lint-license:
	reuse lint
