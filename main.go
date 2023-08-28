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
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	defaultTheme "github.com/go-enjin/default-enjin-theme"

	"github.com/go-enjin/golang-org-x-text/language"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/drivers/fts/bleve"
	"github.com/go-enjin/be/drivers/kvs/gocache"
	"github.com/go-enjin/be/features/fs/themes"
	"github.com/go-enjin/be/features/pages/pql"
	"github.com/go-enjin/be/features/pages/robots"
	"github.com/go-enjin/be/features/pages/search"
	"github.com/go-enjin/be/features/pages/sitemap"
	"github.com/go-enjin/be/pkg/feature"
	"github.com/go-enjin/be/pkg/lang"
	"github.com/go-enjin/be/presets/defaults"
	"github.com/go-enjin/be/presets/essentials"

	"github.com/go-enjin/website-thisip-fyi/pkg/features/thisip"
)

var (
	wwwLocales feature.Feature
	wwwContent feature.Feature
	wwwPublic  feature.Feature
	wwwMenu    feature.Feature

	enjaLocales feature.Feature
	enjaContent feature.Feature
	enjaPublic  feature.Feature
	enjaMenu    feature.Feature

	fThemes            feature.Feature
	fCachePagesPqlWWW  feature.Feature
	fCachePagesPqlENJA feature.Feature

	hotReload bool

	enjaEnDomain string
	enjaJaDomain string
)

func init() {
	fCachePagesPqlWWW = gocache.NewTagged(gPagesPqlKvsFeatureWWW).AddMemoryCache(gPagesPqlKvsCacheWWW).Make()
	fCachePagesPqlENJA = gocache.NewTagged(gPagesPqlKvsFeatureENJA).AddMemoryCache(gPagesPqlKvsCacheENJA).Make()

	if v, ok := os.LookupEnv("BE_ENJA_EN_DOMAIN"); ok {
		enjaEnDomain = v
	} else {
		enjaEnDomain = "http://en.localhost:3334"
	}
	if v, ok := os.LookupEnv("BE_ENJA_JA_DOMAIN"); ok {
		enjaJaDomain = v
	} else {
		enjaJaDomain = "http://ja.localhost:3334"
	}
}

func setup(eb *be.EnjinBuilder) *be.EnjinBuilder {
	eb.SiteName("Go-Enjin").
		SiteTagLine("Done is the enjin of more.").
		SiteCopyrightName("Go-Enjin").
		SiteCopyrightNotice("© 2022 All rights reserved").
		SiteLanguageDisplayNames(map[language.Tag]string{
			language.English: "EN",
		}).
		Set("SiteTitleReversed", true).
		Set("SiteTitleSeparator", " | ").
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo").
		AddFeature(fThemes)
	return eb
}

func features(eb feature.Builder) feature.Builder {
	return eb.
		AddPreset(defaults.New().Make()).
		AddFeature(sitemap.New().Make()).
		AddFeature(robots.New().
			AddSitemap("/sitemap.xml").
			AddRuleGroup(robots.NewRuleGroup().
				AddUserAgent("*").AddAllowed("/").Make(),
			).Make()).
		AddFeature(search.New().Make()).
		AddFeature(thisip.New().Make()).
		SetPublicAccess(gPublicActions...).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
		HotReload(hotReload)
}

func main() {
	www, enja := be.New(), be.New()

	setup(www).SiteTag("WWW").
		SiteDefaultLanguage(language.English).
		SiteLanguageMode(lang.NewPathMode().Make()).
		AddFeature(bleve.NewTagged("bleve-fts-www").Make()).
		AddFeature(fCachePagesPqlWWW).
		AddFeature(pql.NewTagged("pages-pql-www").
			SetKeyValueCache(gPagesPqlKvsFeatureWWW, gPagesPqlKvsCacheWWW).
			Make())
	features(www).
		AddFeature(wwwMenu).
		AddFeature(wwwPublic).
		AddFeature(wwwContent).
		AddFeature(wwwLocales)

	setup(enja).SiteTag("ENJA").
		SiteDefaultLanguage(language.English).
		SiteLanguageMode(lang.NewDomainMode().
			Set(language.English, enjaEnDomain).
			Set(language.Japanese, enjaJaDomain).
			Make(),
		).
		AddFeature(bleve.NewTagged("bleve-fts-enja").Make()).
		AddFeature(fCachePagesPqlENJA).
		AddFeature(pql.NewTagged("pages-pql-enja").
			SetKeyValueCache(gPagesPqlKvsFeatureENJA, gPagesPqlKvsCacheENJA).
			Make()).
		AddFlags(
			&cli.StringFlag{
				Name:     "enja-en-domain",
				Category: "locale domains",
				Usage:    "enja site EN domain name (only env works)",
				EnvVars:  enja.MakeEnvKeys("ENJA_EN_DOMAIN"),
			},
			&cli.StringFlag{
				Name:     "enja-ja-domain",
				Category: "locale domains",
				Usage:    "enja site JA domain name (only env works)",
				EnvVars:  enja.MakeEnvKeys("ENJA_JA_DOMAIN"),
			},
		)
	features(enja).
		AddFeature(enjaMenu).
		AddFeature(enjaPublic).
		AddFeature(enjaContent).
		AddFeature(enjaLocales)

	enjin := be.New().
		IncludeEnjin(www, enja).
		SiteTag("MAIN").
		SiteName("main").
		SiteDefaultLanguage(language.English).
		SiteSupportedLanguages(language.English).
		AddPreset(essentials.New().Make()).
		AddFeature(themes.New().
			Include(defaultTheme.Theme()).
			SetTheme(defaultTheme.Name).
			Make()).
		AddPageFromString("204.tmpl", main204tmpl).
		AddPageFromString("404.tmpl", main404tmpl).
		AddPageFromString("500.tmpl", main500tmpl).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
		SetPublicAccess(gPublicActions...).
		Build()

	if err := enjin.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

var (
	gPublicActions = feature.Actions{
		feature.NewAction("fs-content", "view", "page"),
		feature.NewAction("enjin", "view", "page"),
	}
)

const (
	gPagesPqlKvsFeatureWWW = "pages-pql-kvs-feature-www"
	gPagesPqlKvsCacheWWW   = "pages-pql-kvs-cache-www"

	gPagesPqlKvsFeatureENJA = "pages-pql-kvs-feature-enja"
	gPagesPqlKvsCacheENJA   = "pages-pql-kvs-cache-enja"

	main500tmpl = `500 - {{ _ "Internal Server Error" }}`
	main404tmpl = `404 - {{ _ "Not Found" }}`
	main204tmpl = `+++
url = "/"
layout = "none"
+++
204 - {{ _ "No Content" }}`
)