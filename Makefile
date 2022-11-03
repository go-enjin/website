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

BE_LOCAL_PATH ?= ../be

APP_NAME    ?= be-website
APP_SUMMARY ?= Go-Enjin Website

DENY_DURATION ?= 60

COMMON_TAGS = page_search,header_proxy,papertrail,semanticEnjinTheme,excludeDefaultTheme
BUILD_TAGS = embeds,minify,$(COMMON_TAGS)
DEV_BUILD_TAGS = locals,htmlify,$(COMMON_TAGS)
EXTRA_PKGS =

# Custom go.mod locals
GOPKG_KEYS = SET

# Semantic Enjin Theme
SET_GO_PACKAGE = github.com/go-enjin/semantic-enjin-theme
SET_LOCAL_PATH = ../semantic-enjin-theme

include ./Enjin.mk
