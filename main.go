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
	"strconv"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/go-enjin/be/features/fs/content"
	"github.com/go-enjin/be/features/pages/metrics"
	"github.com/go-enjin/be/features/srv/factories/nonces"
	"github.com/go-enjin/be/features/srv/factories/tokens"
	defaultTheme "github.com/go-enjin/default-enjin-theme"
	"github.com/go-enjin/golang-org-x-text/message"

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
	wwwContent feature.Feature
	wwwPublic  feature.Feature
	wwwMenu    feature.Feature
	wwwLocales feature.Feature

	wwwDevFeatures feature.Features

	enjaContent feature.Feature
	enjaPublic  feature.Feature
	enjaMenu    feature.Feature
	enjaEditor  feature.Feature
	enjaLocales feature.Feature

	siteThemes  feature.Feature
	siteLocales feature.Feature

	hotReload bool

	enjaEnDomain string
	enjaJaDomain string
)

func init() {
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
		SetEnjinTextFn(func(printer *message.Printer) (text feature.EnjinText) {
			name := printer.Sprintf("Go-Enjin")
			return feature.EnjinText{
				Name:            name,
				TagLine:         printer.Sprintf("Done is the enjin of more."),
				CopyrightName:   name,
				CopyrightYear:   "2022-" + strconv.Itoa(time.Now().Year()),
				CopyrightNotice: printer.Sprintf("All rights reserved"),
			}
		}).
		SiteCopyrightName("Go-Enjin").
		SiteCopyrightNotice("All rights reserved").
		SiteSupportedLanguages(language.English, language.Japanese).
		SiteLanguageDisplayNames(map[language.Tag]string{
			language.English: "EN",
			language.French:  "FR",
		}).
		Set("SiteTitleReversed", true).
		Set("SiteTitleSeparator", " | ").
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo")
	return eb
}

func features(eb feature.Builder, l feature.ServiceListener) feature.Builder {
	return eb.
		AddPreset(defaults.New().SetListener(l).Make()).
		AddFeature(gocache.NewTagged(gNoncesKvsFeatureWWW).
			AddMemoryCache(gNoncesKvsCacheWWW).
			Make()).
		AddFeature(nonces.New().
			SetKeyValueCache(gNoncesKvsFeatureWWW, gNoncesKvsCacheWWW).
			Make()).
		AddFeature(gocache.NewTagged(gNoncesKvsFeatureENJA).
			AddMemoryCache(gNoncesKvsCacheENJA).
			Make()).
		AddFeature(tokens.New().
			SetKeyValueCache(gNoncesKvsFeatureENJA, gNoncesKvsCacheENJA).
			Make()).
		AddFeature().
		AddFeature(metrics.New().
			AddArchetype("blog").
			Make()).
		AddFeature(sitemap.New().Make()).
		AddFeature(robots.New().
			AddSitemap("/sitemap.xml").
			AddRuleGroup(robots.NewRuleGroup().
				AddUserAgent("*").
				AddAllowed("/").
				Make()).
			Make()).
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
		AddFeature(bleve.NewTagged(gBleveFtsFeatureWWW).Make()).
		AddFeature(gocache.NewTagged(gPagesPqlKvsFeatureWWW).
			AddMemoryCache(gPagesPqlKvsCacheWWW).
			Make()).
		AddFeature(pql.NewTagged(gPagesPqlFeatureWWW).
			SetKeyValueCache(gPagesPqlKvsFeatureWWW, gPagesPqlKvsCacheWWW).
			Make())
	features(www, nil).
		AddFeature(siteThemes).
		AddFeature(wwwLocales).
		AddFeature(siteLocales).
		AddFeature(wwwDevFeatures...).
		AddFeature(wwwMenu).
		AddFeature(wwwPublic).
		AddFeature(wwwContent)

	setup(enja).SiteTag("ENJA").
		SiteDefaultLanguage(language.English).
		SiteLanguageMode(lang.NewDomainMode().
			Set(language.English, enjaEnDomain).
			Set(language.Japanese, enjaJaDomain).
			Make(),
		).
		AddFeature(bleve.NewTagged(gBleveFtsFeatureENJA).Make()).
		AddFeature(gocache.NewTagged(gPagesPqlKvsFeatureENJA).
			AddMemoryCache(gPagesPqlKvsCacheENJA).
			Make()).
		AddFeature(pql.NewTagged(gPagesPqlFeatureENJA).
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
	features(enja, nil).
		AddFeature(siteThemes).
		AddFeature(enjaLocales).
		AddFeature(siteLocales).
		AddFeature(enjaMenu).
		AddFeature(enjaPublic).
		AddFeature(enjaContent).
		AddFeature(enjaEditor)

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
		feature.NewAction(content.Tag.Kebab(), "view", "page"),
		feature.NewAction(gSiteContentTag.Kebab(), "view", "page"),
		feature.NewAction(feature.EnjinTag.Kebab(), "view", "page"),
		feature.NewAction("rw-pages", "view", "page"),
		feature.NewAction("ro-pages", "view", "page"),
	}
)

const (
	gSiteContentTag feature.Tag = "www-content"
	gSiteLocalesTag feature.Tag = "site-locales"
	gSiteThemesTag  feature.Tag = "site-themes"

	gNoncesKvsFeatureWWW  = "nonces-kvs-feature-www"
	gNoncesKvsCacheWWW    = "nonces-kvs-cache-www"
	gNoncesKvsFeatureENJA = "nonces-kvs-feature-enja"
	gNoncesKvsCacheENJA   = "nonces-kvs-cache-enja"

	gPagesPqlKvsFeatureWWW = "pages-pql-kvs-feature-www"
	gPagesPqlKvsCacheWWW   = "pages-pql-kvs-cache-www"

	gSiteKvsFeature = "site-kvs-feature"
	gSiteKvsCache   = "site-kvs-cache"

	gSiteUsersKvsFeature         = "site-users-kvs-feature"
	gSiteUsersKvsCache           = "site-users-kvs-cache"
	gSiteAuthEmailKvsFeature     = "site-auth-email-kvs-feature"
	gSiteAuthEmailTokenKvsCache  = "site-auth-email-token-kvs-cache"
	gSiteAuthEmailBackupKvsCache = "site-auth-email-backup-kvs-cache"

	gAdminKvsCache                = "admin-kvs-cache"
	gAdminAuthEmailTokenKvsCache  = "admin-auth-email-token-kvs-cache"
	gAdminAuthEmailBackupKvsCache = "admin-auth-email-backup-kvs-cache"

	gPagesPqlKvsFeatureENJA = "pages-pql-kvs-feature-enja"
	gPagesPqlKvsCacheENJA   = "pages-pql-kvs-cache-enja"
	gPagesPqlFeatureWWW     = "pages-pql-www"
	gPagesPqlFeatureENJA    = "pages-pql-enja"

	gBleveFtsFeatureWWW  = "bleve-fts-www"
	gBleveFtsFeatureENJA = "bleve-fts-enja"

	main500tmpl = `500 - {{ _ "Internal Server Error" }}`
	main404tmpl = `404 - {{ _ "Not Found" }}`
	main204tmpl = `+++
layout = "none"
+++
204 - {{ _ "No Content" }}`
)