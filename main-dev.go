//go:build !production

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
	"time"

	"github.com/go-enjin/be"
	"github.com/go-enjin/be/drivers/db/gorm"
	"github.com/go-enjin/be/drivers/email/fakemail"
	"github.com/go-enjin/be/drivers/email/gomail"
	"github.com/go-enjin/be/drivers/kvs/gocache"
	"github.com/go-enjin/be/features/fs/content"
	"github.com/go-enjin/be/features/fs/email"
	"github.com/go-enjin/be/features/fs/locale"
	"github.com/go-enjin/be/features/fs/menu"
	"github.com/go-enjin/be/features/fs/public"
	"github.com/go-enjin/be/features/fs/site-users"
	"github.com/go-enjin/be/features/fs/themes"
	"github.com/go-enjin/be/features/site"
	"github.com/go-enjin/be/features/site/auth"
	"github.com/go-enjin/be/features/site/auth/otp/app_totp"
	"github.com/go-enjin/be/features/site/auth/otp/backup_codes"
	"github.com/go-enjin/be/features/site/auth/otp/email_hotp"
	"github.com/go-enjin/be/features/site/auth/provider/email_backup"
	"github.com/go-enjin/be/features/site/auth/provider/email_token"
	"github.com/go-enjin/be/features/site/dashboard"
	"github.com/go-enjin/be/features/site/fs-editor"
	enjinInfo "github.com/go-enjin/be/features/site/fs-editor/enjin-info"
	"github.com/go-enjin/be/features/site/fs-editor/locales"
	"github.com/go-enjin/be/features/site/fs-editor/menus"
	"github.com/go-enjin/be/features/site/fs-editor/pages"
	themesEditor "github.com/go-enjin/be/features/site/fs-editor/themes"
	"github.com/go-enjin/be/features/site/fs-editor/unsafe"
	"github.com/go-enjin/be/features/site/profile"
	"github.com/go-enjin/be/features/site/settings"
	"github.com/go-enjin/be/features/site/user-manager"
	"github.com/go-enjin/be/features/user/base/ipenv"
	"github.com/go-enjin/be/pkg/cli/env"
	"github.com/go-enjin/be/pkg/feature"
	semantic "github.com/go-enjin/semantic-enjin-theme"
)

func init() {
	// locals environment, early startup debug logging
	// log.Config.LogLevel = log.LevelDebug
	// log.Config.Apply()

	//siteListener = ngrokio.New().Make()

	wwwMenu = menu.New().
		MountLocalPath("/", "menus/www").
		Make()
	wwwPublic = public.New().
		MountLocalPath("/", "public").
		Make()
	wwwContent = content.NewTagged(gSiteContentTag).
		MountLocalPath("/", "content/www").
		AddToIndexProviders(gPagesPqlFeatureWWW).
		AddToSearchProviders(gBleveFtsFeatureWWW).
		Make()
	wwwLocales = locale.NewTagged(gSiteLocalesTag).
		MountLocalPath("/", "locales").
		Make()

	enjaMenu = menu.New().
		MountLocalPath("/", "menus/enja").
		Make()
	enjaPublic = public.New().
		MountLocalPath("/", "public").
		Make()
	enjaContent = content.NewTagged(gSiteContentTag).
		MountLocalPath("/", "content/enja").
		AddToIndexProviders(gPagesPqlFeatureENJA).
		AddToSearchProviders(gBleveFtsFeatureENJA).
		Make()
	enjaLocales = locale.NewTagged(gSiteLocalesTag).
		MountLocalPath("/", "locales").
		Make()

	siteThemes = themes.NewTagged(gSiteThemesTag).
		Include(semantic.Theme()).
		Include(semantic.SiteTheme()).
		LocalTheme("themes/go-enjin").
		SetTheme("go-enjin").
		Make()
	siteLocales = be.MakeLocalLocales()

	userPermissions := feature.Actions{
		feature.NewAction(feature.EnjinTag.Kebab(), "view", "page"),
		feature.NewAction(auth.Tag.Kebab(), "access", "feature"),
		feature.NewAction(settings.Tag.Kebab(), "access", "feature"),
		feature.NewAction(dashboard.Tag.Kebab(), "access", "feature"),
		feature.NewAction(profile.Tag.Kebab(), "access", "feature"),
		feature.NewAction(profile.Tag.Kebab(), "view-own", "profile"),
		feature.NewAction("www-users", "view-own", "user"),
		feature.NewAction("www-users", "update-own", "user"),
		feature.NewAction("www-users", "delete-own", "user"),
	}
	for _, tag := range []feature.Tag{
		editor.Tag, enjinInfo.Tag, locales.Tag, menus.Tag, pages.Tag,
	} {
		kebab := tag.Kebab()
		userPermissions = append(userPermissions,
			feature.NewAction(kebab, "access", "feature"),
			feature.NewAction(kebab, "view", "file-browser"),
			feature.NewAction(kebab, "view", "file-editor"),
			feature.NewAction(kebab, "edit", "file-editor"),
			feature.NewAction(kebab, "create", "file-editor"),
			feature.NewAction(kebab, "delete", "file-editor"),
		)
	}

	adminPermissions := append(userPermissions,
		feature.NewAction("admin-auth", "access", "feature"),
		feature.NewAction("admin-auth", "reset-own", "multi-factors"),
		feature.NewAction("admin-auth", "reset-other", "multi-factors"),
		feature.NewAction("admin-settings", "access", "feature"),
		feature.NewAction("admin-dashboard", "access", "feature"),
		feature.NewAction("admin-profile", "view-own", "profile"),
		feature.NewAction("admin-profile", "view-other", "profile"),
		feature.NewAction("www-users", "create", "user"),
		feature.NewAction("www-users", "view-own", "user"),
		feature.NewAction("www-users", "update-own", "user"),
		feature.NewAction("www-users", "delete-own", "user"),
		feature.NewAction("www-users", "view-other", "user"),
		feature.NewAction("www-users", "update-other", "user"),
		feature.NewAction("www-users", "delete-other", "user"),
		feature.NewAction("admin-user-manager", "access", "feature"),
		feature.NewAction("admin-user-manager", "create", "user"),
		feature.NewAction("admin-user-manager", "update", "user"),
		feature.NewAction("admin-user-manager", "delete", "user"),
	)
	for _, tag := range []feature.Tag{
		//auth.Tag, editor.Tag, enjinInfo.Tag, unsafe.Tag, locales.Tag, menus.Tag, pages.Tag, themesEditor.Tag,
		"admin-auth", "admin-settings", "admin-dashboard", "admin-profile", "admin-user-manager", "admin-editor",
		"admin-enjin-info", "admin-unsafe", "admin-locales", "admin-menus", "admin-pages", "admin-themes",
	} {
		kebab := tag.Kebab()
		adminPermissions = append(adminPermissions,
			feature.NewAction(kebab, "access", "feature"),
			feature.NewAction(kebab, "view", "file-browser"),
			feature.NewAction(kebab, "view", "file-editor"),
			feature.NewAction(kebab, "edit", "file-editor"),
			feature.NewAction(kebab, "create", "file-editor"),
			feature.NewAction(kebab, "delete", "file-editor"),
		)
	}

	emailAccount := env.Get("BE_DEV_USE_EMAIL_ACCOUNT", "fakemail")
	switch emailAccount {
	case "fakemail":
		wwwDevFeatures = append(wwwDevFeatures, fakemail.New().
			AddAccount("fakemail").
			Make())
	case "default":
		wwwDevFeatures = append(wwwDevFeatures, gomail.NewTagged("gomail").
			AddAccount("default", gomail.SmtpConfig{}).
			Make())
	default:
		panic(fmt.Sprintf("error: BE_DEV_USE_EMAIL_ACCOUNT value must be one of: \"default\" or \"fakemail\""))
	}

	wwwDevFeatures = append(wwwDevFeatures,
		gocache.NewTagged(gSiteAuthEmailKvsFeature).
			AddMemoryCache(gSiteAuthEmailTokenKvsCache).
			AddMemoryCache(gSiteAuthEmailBackupKvsCache).
			AddMemoryCache(gAdminAuthEmailTokenKvsCache).
			AddMemoryCache(gAdminAuthEmailBackupKvsCache).
			Make(),
		gocache.NewTagged(gSiteUsersKvsFeature).
			AddMemoryCache(gSiteUsersKvsCache).
			Make(),
		gocache.NewTagged(gSiteKvsFeature).
			AddMemoryCache(gSiteKvsCache).
			AddMemoryCache(gAdminKvsCache).
			Make(),
		gorm.NewTagged("gorm-db").
			AddConnection("locales-rw").
			AddConnection("www-users").
			AddConnection("rw-pages-db").
			Make(),
		content.NewTagged("enja-content").
			MountLocalPath("/", "content/enja").
			Make(),
		content.NewTagged("rw-pages").
			MountGormDBPath("/", "content", "rw-pages-db").
			AddToIndexProviders(gPagesPqlFeatureWWW).
			AddToSearchProviders(gBleveFtsFeatureWWW).
			Make(),
		ipenv.New().Make(),
		email.New().
			MountLocalPath("/", "emails").
			Make(),
		site_users.NewTagged("www-users").
			MountGormDBPath("/", "", "www-users").
			InitGroup("users", userPermissions...).
			InitGroup("admins", adminPermissions...).
			SetKeyValueCache(gSiteUsersKvsFeature, gSiteUsersKvsCache).
			Make(),
		site.NewTagged("www-site").
			SetSitePath("/site").
			SetSiteTheme(semantic.SiteName).
			SetSiteUsers("www-users").
			SetKeyValueCache(gSiteKvsFeature, gSiteKvsCache).
			UseSiteRootFeature(
				dashboard.New().
					SetSiteFeatureKey("dashboard").
					Make()).
			IncludeSiteFeatures(
				editor.New().
					Include(
						enjinInfo.New().Make(),
						pages.New().
							AddContentFileSystems(gSiteContentTag, "rw-pages", "enja-content").
							Make(),
						menus.New().
							AddMenuFileSystems(menu.NewTagged("enja-menus").
								MountLocalPath("/", "menus/enja").
								Make()).
							Make(),
						locales.New().Make(),
					).
					Make(),
				profile.New().
					SetSiteFeatureKey("profile").
					SetSiteLocaleEnabled(true).
					SetSelfProfilePageEnabled(true).
					SetOtherProfilePagesEnabled(true).
					DefaultProfileImageNames().
					Make(),
				settings.New().
					SetSiteFeatureKey("settings").
					Make(),
				auth.New().
					SetSessionDuration(time.Hour*24).
					SetVerifiedDuration(time.Minute*10).
					SetNumRequiredFactors(0).
					SetJwtCookieName("www-site-cookie").
					SetRoutePaths("/sign-in", "/sign-out", "/challenge").
					IncludeMultiFactors(
						app_totp.New().Make(),
						email_hotp.New().
							SetEmailAccount(emailAccount).
							SetEmailProvider(email.Tag).
							Make(),
						backup_codes.New().Make(),
					).
					IncludeProviders(
						email_token.New().
							SetKeyValueCache(gSiteAuthEmailKvsFeature, gSiteAuthEmailTokenKvsCache).
							SetEmailAccount(emailAccount).
							SetEmailProvider(email.Tag).
							Make(),
						email_backup.New().
							SetKeyValueCache(gSiteAuthEmailKvsFeature, gSiteAuthEmailBackupKvsCache).
							SetEmailAccount(emailAccount).
							SetEmailProvider(email.Tag).
							Make(),
					).
					Make(),
			).
			Make(),
		site.NewTagged("www-admin").
			SetSitePath("/admin").
			SetSiteTheme(semantic.SiteName).
			SetSiteUsers("www-users").
			SetKeyValueCache(gSiteKvsFeature, gAdminKvsCache).
			UseSiteRootFeature(dashboard.NewTagged("admin-dashboard").
				SetSiteFeatureKey("dashboard").
				Make()).
			IncludeSiteFeatures(
				editor.NewTagged("admin-editor").
					Include(
						enjinInfo.NewTagged("admin-enjin-info").Make(),
						pages.NewTagged("admin-pages").
							AddContentFileSystems(gSiteContentTag, "rw-pages", "enja-content").
							Make(),
						menus.NewTagged("admin-menus").
							AddMenuFileSystems(menu.NewTagged("admin-enja-menus").
								MountLocalPath("/", "menus/enja").
								Make()).
							Make(),
						locales.NewTagged("admin-locales").Make(),
						themesEditor.NewTagged("admin-themes").Make(),
						unsafe.NewTagged("admin-unsafe").Make(),
					).
					Make(),
				profile.NewTagged("admin-profile").
					SetSiteFeatureKey("profile").
					SetSiteLocaleEnabled(true).
					DefaultProfileImageNames().
					Make(),
				settings.NewTagged("admin-settings").
					SetSiteFeatureKey("settings").
					Make(),
				user_manager.NewTagged("admin-user-manager").
					SetSiteFeatureKey("users").
					Make(),
				auth.NewTagged("admin-auth").
					SetSessionDuration(time.Hour).
					SetVerifiedDuration(time.Minute*10).
					SetNumRequiredFactors(2).
					SetJwtCookieName("www-admin-cookie").
					SetRoutePaths("/sign-in", "/sign-out", "/challenge").
					IncludeMultiFactors(
						app_totp.NewTagged("admin-totp").Make(),
						email_hotp.NewTagged("admin-hotp").
							SetEmailAccount(emailAccount).
							SetEmailProvider(email.Tag).
							Make(),
						backup_codes.NewTagged("admin-backup-codes").Make(),
					).
					IncludeProviders(
						email_token.NewTagged("admin-email-token").
							SetKeyValueCache(gSiteAuthEmailKvsFeature, gAdminAuthEmailTokenKvsCache).
							SetEmailAccount(emailAccount).
							SetEmailProvider(email.Tag).
							Make(),
						email_backup.NewTagged("admin-email-backup").
							SetKeyValueCache(gSiteAuthEmailKvsFeature, gAdminAuthEmailBackupKvsCache).
							SetEmailAccount(emailAccount).
							SetEmailProvider(email.Tag).
							Make(),
					).
					Make(),
			).
			Make(),
	)

	hotReload = true
}