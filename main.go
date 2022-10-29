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

	"github.com/go-enjin/golang-org-x-text/language"

	semantic "github.com/go-enjin/semantic-enjin-theme"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/features/log/papertrail"
	"github.com/go-enjin/be/features/pages/formats"
	"github.com/go-enjin/be/features/pages/search"
	"github.com/go-enjin/be/features/requests/headers/proxy"
	"github.com/go-enjin/be/pkg/feature"
)

var fMinifyHtmlify feature.Feature
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
		SiteTag("go-enjin").
		SiteName("Go-Enjin").
		SiteTagLine("Done is the enjin of more.").
		SiteCopyrightName("Go-Enjin").
		SiteCopyrightNotice("© 2022 All rights reserved").
		SiteLanguageMode("path").
		SiteDefaultLanguage(language.English).
		// SiteDefaultLanguage(language.Japanese).
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo").
		// Set("SiteLoadingEffect", "true").
		AddHtmlHeadTag("meta", map[string]string{
			"name":    "robots",
			"content": "none",
		}).
		AddFeature(fMenu).
		AddFeature(fPublic).
		AddFeature(fLocales).
		AddFeature(fContent).
		AddFeature(search.New().Make()).
		AddFeature(proxy.New().Enable().Make()).
		AddFeature(papertrail.Make()).
		AddFeature(fMinifyHtmlify).
		SetStatusPage(404, "/404").
		SetStatusPage(500, "/500").
		HotReload(hotReload).
		Build()
	if err := enjin.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}