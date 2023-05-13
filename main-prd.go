//go:build production

// Copyright (c) 2022  The Go-Enjin Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"embed"

	semantic "github.com/go-enjin/semantic-enjin-theme"

	"github.com/go-enjin/be/features/fs/content"
	"github.com/go-enjin/be/features/fs/locale"
	"github.com/go-enjin/be/features/fs/menu"
	"github.com/go-enjin/be/features/fs/public"
)

//go:embed content/www/**
var contentFsWWW embed.FS

//go:embed content/enja/**
var contentFsENJA embed.FS

//go:embed public/**
var publicFs embed.FS

//go:embed menus/www/**
var menuFsWWW embed.FS

//go:embed menus/enja/**
var menuFsENJA embed.FS

//go:embed themes/**
//go:embed themes/*/layouts/_default/*
var themeFs embed.FS

//go:embed locales/*
var localesFs embed.FS

func init() {
	wwwMenu = menu.New().
		MountEmbedPath("/", "menus", menuFsWWW).
		Make()
	wwwPublic = public.New().
		MountEmbedPath("/", "public", publicFs).
		Make()
	wwwPublic = content.New().
		MountEmbedPath("/", "content", contentFsWWW).
		AddToIndexProviders("pages-pql").
		SetKeyValueCache(gFsContentKvsFeatureWWW, gFsContentKvsCacheWWW).
		Make()
	wwwLocales = locale.New().Include("locales", localesFs).Make()

	enjaMenu = menu.New().
		MountEmbedPath("/", "menus", menuFsWWW).
		Make()
	enjaPublic = public.New().
		MountEmbedPath("/", "public", publicFs).
		Make()
	enjaContent = content.New().
		MountEmbedPath("/", "content", contentFsWWW).
		AddToIndexProviders("pages-pql").
		SetKeyValueCache(gFsContentKvsFeatureWWW, gFsContentKvsCacheWWW).
		Make()
	enjaLocales = locale.New().Include("locales", localesFs).Make()

	fThemes = beFsTheme.New().
		AddTheme(semantic.SemanticEnjinTheme()).
		EmbedTheme("themes/go-enjin", themeFs).
		SetTheme("go-enjin").
		Make()

	hotReload = false
}
