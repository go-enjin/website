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

	"github.com/go-enjin/be/features/outputs/htmlify"
	"github.com/go-enjin/be/features/pages/permalink"
	"github.com/go-enjin/be/features/pages/robots"
	"github.com/go-enjin/be/features/pages/search"
	"github.com/go-enjin/be/features/pages/sitemap"
	"github.com/go-enjin/be/pkg/lang"
	"github.com/go-enjin/be/pkg/theme"
	"github.com/go-enjin/golang-org-x-text/language"

	semantic "github.com/go-enjin/semantic-enjin-theme"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/features/log/papertrail"
	"github.com/go-enjin/be/features/pages/formats"
	"github.com/go-enjin/be/features/requests/headers/proxy"
	"github.com/go-enjin/be/pkg/feature"
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

	hotReload bool

	enjaEnDomain string
	enjaJaDomain string
)

func init() {
	if v, ok := os.LookupEnv("BE_ENJA_EN_DOMAIN"); ok {
		enjaEnDomain = v
	} else {
		enjaEnDomain = "http://en.go-enjin.localhost:3334"
	}
	if v, ok := os.LookupEnv("BE_ENJA_JA_DOMAIN"); ok {
		enjaJaDomain = v
	} else {
		enjaJaDomain = "http://ja.go-enjin.localhost:3334"
	}
}

const (
	main500tmpl = `500 - {{ _ "Internal Server Error" }}`
	main404tmpl = `404 - {{ _ "Not Found" }}`
	main204tmpl = `+++
url = "/"
+++
204 - {{ _ "No Content" }}`
)

func setup(eb *be.EnjinBuilder) *be.EnjinBuilder {
	eb.SiteName("Go-Enjin").
		SiteTagLine("Done is the enjin of more.").
		SiteCopyrightName("Go-Enjin").
		SiteCopyrightNotice("© 2022 All rights reserved").
		AddFeature(formats.New().Defaults().Make()).
		AddTheme(semantic.SemanticEnjinTheme()).
		AddTheme(goEnjinTheme()).
		SetTheme("go-enjin").
		SiteLanguageDisplayNames(map[language.Tag]string{
			language.English: "EN",
		}).
		Set("SiteTitleReversed", true).
		Set("SiteTitleSeparator", " | ").
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo")
	// Set("SiteLoadingEffect", "true").
	return eb
}

func features(eb feature.Builder) feature.Builder {
	return eb.
		AddFeature(papertrail.Make()).
		AddFeature(sitemap.New().Make()).
		AddFeature(robots.New().
			AddSitemap("/sitemap.xml").
			AddRuleGroup(robots.NewRuleGroup().
				AddUserAgent("*").AddAllowed("/").Make(),
			).Make()).
		AddFeature(proxy.New().Enable().Make()).
		AddFeature(permalink.New().Make()).
		AddFeature(search.New().Make()).
		AddFeature(htmlify.New().Make()).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
		HotReload(hotReload)
}

func main() {
	www, enja := be.New(), be.New()

	setup(www).SiteTag("WWW").
		SiteDefaultLanguage(language.English).
		SiteLanguageMode(lang.NewPathMode().Make())
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
		AddFlags(
			&cli.StringFlag{
				Name:    "enja-en-domain",
				EnvVars: enja.MakeEnvKeys("ENJA_EN_DOMAIN"),
			},
			&cli.StringFlag{
				Name:    "enja-ja-domain",
				EnvVars: enja.MakeEnvKeys("ENJA_JA_DOMAIN"),
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
		AddTheme(theme.DefaultTheme()).
		SiteDefaultLanguage(language.English).
		SiteSupportedLanguages(language.English).
		AddPageFromString("/", main204tmpl).
		AddPageFromString("/404.tmpl", main404tmpl).
		AddPageFromString("/500.tmpl", main500tmpl).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
		Build()

	if err := enjin.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}