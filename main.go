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

	"github.com/go-enjin/be/features/outputs/htmlify"
	"github.com/go-enjin/be/features/pages/permalink"
	"github.com/go-enjin/be/features/pages/robots"
	"github.com/go-enjin/be/features/pages/search"
	"github.com/go-enjin/be/features/pages/sitemap"
	"github.com/go-enjin/be/pkg/lang"
	"github.com/go-enjin/golang-org-x-text/language"

	semantic "github.com/go-enjin/semantic-enjin-theme"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/features/log/papertrail"
	"github.com/go-enjin/be/features/pages/formats"
	"github.com/go-enjin/be/features/requests/headers/proxy"
	"github.com/go-enjin/be/pkg/feature"
)

var fLocales feature.Feature
var fContent feature.Feature
var fPublic feature.Feature
var fMenu feature.Feature
var hotReload bool

func main() {
	enjin := be.New().
		AddFeature(formats.New().Defaults().Make()).
		AddTheme(semantic.SemanticEnjinTheme()).
		AddTheme(goEnjinTheme()).
		SetTheme("go-enjin").
		SiteTag("BE").
		SiteName("Go-Enjin").
		SiteTagLine("Done is the enjin of more.").
		SiteCopyrightName("Go-Enjin").
		SiteCopyrightNotice("© 2022 All rights reserved").
		SiteDefaultLanguage(language.English).
		// SiteDefaultLanguage(language.Japanese).
		// SiteLanguageMode(lang.NewQueryMode().SetQueryParameter("lang").Make()).
		SiteLanguageMode(lang.NewPathMode().Make()).
		// SiteLanguageMode(lang.NewDomainMode().
		// 	Set(language.English, "http://en.localhost:3335").
		// 	Set(language.Japanese, "http://ja.localhost:3335").
		// 	Make(),
		// ).
		SiteLanguageDisplayNames(map[language.Tag]string{
			language.English: "EN",
		}).
		Set("SiteTitleReversed", true).
		Set("SiteTitleSeparator", " | ").
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo").
		// Set("SiteLoadingEffect", "true").
		AddFeature(papertrail.Make()).
		AddFeature(sitemap.New().Make()).
		AddFeature(robots.New().
			// SiteRobotsHeader("noindex").
			// SiteRobotsMetaTag("none").
			AddSitemap("/sitemap.xml").
			AddRuleGroup(robots.NewRuleGroup().
				AddUserAgent("*").AddAllowed("/").Make(),
			).Make()).
		AddFeature(proxy.New().Enable().Make()).
		AddFeature(permalink.New().Make()).
		AddFeature(search.New().Make()).
		AddFeature(fMenu).
		AddFeature(fPublic).
		AddFeature(fLocales).
		AddFeature(fContent).
		AddFeature(htmlify.New().Make()).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
		HotReload(hotReload).
		Build()
	if err := enjin.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}