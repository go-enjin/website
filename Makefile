#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

# Copyright (c) 2022  The Go-Enjin Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

-include .env

BE_LOCAL_PATH ?= ../be

APP_NAME    ?= be-website
APP_SUMMARY ?= Go-Enjin Website

DENY_DURATION ?= 60

COMMON_TAGS += htmlify
COMMON_TAGS += papertrail
COMMON_TAGS += header_proxy
COMMON_TAGS += basic_auth
COMMON_TAGS += driver_kvs_gocache memory
COMMON_TAGS += page_pql
COMMON_TAGS += page_robots
COMMON_TAGS += driver_fs_embed
COMMON_TAGS += fs_theme fs_menu fs_content fs_public fs_locale
COMMON_TAGS += driver_fts_bleve
COMMON_TAGS += page_sitemap
COMMON_TAGS += page_query
COMMON_TAGS += page_search

BUILD_TAGS     = production embeds $(COMMON_TAGS)
DEV_BUILD_TAGS = locals $(COMMON_TAGS)

# Custom go.mod locals
GOPKG_KEYS = BET SET GOXT DJHT TIF

# Basic Enjin Theme
BET_GO_PACKAGE = github.com/go-enjin/basic-enjin-theme
BET_LOCAL_PATH = ../basic-enjin-theme

# Semantic Enjin Theme
SET_GO_PACKAGE = github.com/go-enjin/semantic-enjin-theme
SET_LOCAL_PATH = ../semantic-enjin-theme

# Go-Enjin gotext package
GOXT_GO_PACKAGE = github.com/go-enjin/golang-org-x-text
GOXT_LOCAL_PATH = ../golang-org-x-text

# Go-Enjin times package
DJHT_GO_PACKAGE = github.com/go-enjin/github-com-djherbis-times
DJHT_LOCAL_PATH = ../github-com-djherbis-times

# Go-Enjin thisip package
TIF_GO_PACKAGE = github.com/go-enjin/website-thisip-fyi
TIF_LOCAL_PATH = ../website-thisip-fyi

include ./Enjin.mk

LANGUAGES = en,ja
LOCALES_CATALOG ?= /dev/null

gen-theme-locales:
	@echo "# generating go-enjin theme locales"
	@${CMD} enjenv be-update-locales \
		-lang=${LANGUAGES} \
		-out=./themes/go-enjin/locales \
		./themes/go-enjin/layouts \
		./content

gen-locales: BE_PKG_LIST=$(shell enjenv be-pkg-list)
gen-locales: gen-theme-locales
	@echo "# generating locales"
	@${CMD} \
		GOFLAGS="-tags=all" \
		gotext -srclang=en update \
			-lang=${LANGUAGES} \
			-out=${LOCALES_CATALOG} \
				${BE_PKG_LIST} \
				github.com/go-enjin/website
	@if [ -d locales ]; then \
		find locales -type f -name "*.gotext.json" -print0 | xargs -n 1 -0 sha256sum; \
	else \
		echo "# error: locales directory not found" 1>&2; \
		false; \
	fi
