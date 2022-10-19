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

	semantic "github.com/go-enjin/semantic-enjin-theme"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/features/log/papertrail"
	"github.com/go-enjin/be/features/outputs/htmlify"
	"github.com/go-enjin/be/features/pages/formats"
	"github.com/go-enjin/be/features/pages/search"
	"github.com/go-enjin/be/features/requests/headers/proxy"
	"github.com/go-enjin/be/pkg/feature"
)

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
		Set("CopyrightName", "Go-Enjin").
		Set("CopyrightNotice", "© 2022 All rights reserved").
		Set("SiteTag", "Go-Enjin").
		Set("SiteName", "Go-Enjin").
		Set("SiteTagLine", "Done is the enjin of more.").
		// Set("SiteLoadingEffect", "true").
		Set("SiteLogoUrl", "/media/go-enjin-logo.png").
		Set("SiteLogoAlt", "Go-Enjin logo").
		AddFeature(fMenu).
		AddFeature(fContent).
		AddFeature(fPublic).
		AddFeature(search.New().Make()).
		AddFeature(proxy.New().Make()).
		AddFeature(papertrail.Make()).
		AddFeature(htmlify.New().Make()).
		SetStatusPage(404, "/404").
		HotReload(hotReload).
		Build()
	if err := enjin.Run(os.Args); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}