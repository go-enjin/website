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

COMMON_TAGS = page_search,header_proxy,papertrail,htmlify,semanticEnjinTheme,excludeDefaultTheme
BUILD_TAGS = embeds,$(COMMON_TAGS)
DEV_BUILD_TAGS = locals,$(COMMON_TAGS)
EXTRA_PKGS =

# Custom go.mod locals
GOPKG_KEYS = SET

# Semantic Enjin Theme
SET_GO_PACKAGE = github.com/go-enjin/semantic-enjin-theme
SET_LOCAL_PATH = ../semantic-enjin-theme

include ./Enjin.mk

heroku-push:
	@git push heroku trunk:main

heroku-logs:
	@heroku logs --tail

# this requires Term::ANSIColor, will error if not present,
# use `make build dev` instead
build-dev-run: build
	@( make dev 2>&1 ) | perl -p -e 'use Term::ANSIColor qw(colored);while (my $$line = <>) {print STDOUT process_line($$line)."\n";}exit(0);sub process_line {my ($$line) = @_;chomp($$line);if ($$line =~ m!^\[(\d+\-\d+\.\d+)\]\s+([A-Z]+)\s+(.+?)\s*$$!) {my ($$datestamp, $$level, $$message) = ($$1, $$2, $$3);my $$colour = "white";if ($$level eq "ERROR") {$$colour = "bold white on_red";} elsif ($$level eq "INFO") {$$colour = "green";} elsif ($$level eq "DEBUG") {$$colour = "yellow";}my $$out = "[".colored($$datestamp, "blue")."]";$$out .= " ".colored($$level, $$colour);if ($$level eq "DEBUG") {$$out .= "\t";if ($$message =~ m!^(.+?)\:(\d+)\s+\[(.+?)\]\s+(.+?)\s*$$!) {my ($$file, $$ln, $$tag, $$info) = ($$1, $$2, $$3, $$4);$$out .= colored($$file, "bright_blue");$$out .= ":".colored($$ln, "blue");$$out .= " [".colored($$tag, "bright_blue")."]";$$out .= " ".colored($$info, "bold cyan");} else {$$out .= $$message;}} elsif ($$level eq "ERROR") {$$out .= "\t".colored($$message, $$colour);} elsif ($$level eq "INFO") {$$out .= "\t".colored($$message, $$colour);} else {$$out .= "\t".$$message;}return $$out;}return $$line;}'
