#!/usr/bin/make --no-print-directory --jobs=1 --environment-overrides -f

# Copyright (c) 2023  The Go-Enjin Authors
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

COMMON_TAGS += driver_db_gorm driver_fs_db_gorm sqlite
COMMON_TAGS += driver_kvs_gocache memory memshard
COMMON_TAGS += driver_fs_embed
COMMON_TAGS += driver_fts_bleve
COMMON_TAGS += user_auth_basic
COMMON_TAGS += user_base_htenv
COMMON_TAGS += papertrail
COMMON_TAGS += page_pql
COMMON_TAGS += page_robots
COMMON_TAGS += page_sitemap
COMMON_TAGS += page_search
COMMON_TAGS += page_metrics
COMMON_TAGS += fs_theme fs_menu fs_content fs_public fs_locale

ADD_TAGS_DEFAULTS := true

BUILD_TAGS     = production embeds $(COMMON_TAGS)
DEV_BUILD_TAGS = locals editor semantic_enjin_editor fs_email driver_email_gomail driver_email_fakemail $(COMMON_TAGS)

# Custom go.mod locals
GOPKG_KEYS += _DEFAULT_THEME
GOPKG_KEYS += _SEMANTIC_THEME
GOPKG_KEYS += _TIMES
GOPKG_KEYS += TIF

# Go-Enjin thisip package
TIF_GO_PACKAGE = github.com/go-enjin/website-thisip-fyi
TIF_LOCAL_PATH = ../website-thisip-fyi

AUTO_CORELIBS_KEYS := true

#
#: Language settings
#

LANGUAGES += ja

MAKE_LOCALES := true
# MAKE_THEME_LOCALES := true
# MAKE_SOURCE_LOCALES := true
# MAKE_CONTENT_LOCALES := false

include ./Enjin.mk
