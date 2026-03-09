// Copyright (c) 2026 Reiner Pröls
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//
// SPDX-License-Identifier: MIT
//
// Author: Reiner Pröls

package main

import "golang.org/x/crypto/bcrypt"

type ExportModeType int

const (
	EXPORTMODE_NEW ExportModeType = iota
	EXPORTMODE_OVERWRITE
	EXPORTMODE_ASK
)

const (
	PREF_THEMEVARIANT_KEY               = "theme"
	PREF_THEMEVARIANT_VALUE             = -1
	PREF_FIRSTDAYOFWEEK_KEY             = "firstdayofweek"
	PREF_FIRSTDAYOFWEEK_VALUE           = 1
	PREF_SPECIALWEEKDAY_0_KEY           = "specialweekday_0"
	PREF_SPECIALWEEKDAY_0_VALUE         = 0
	PREF_SPECIALWEEKDAY_1_KEY           = "specialweekday_1"
	PREF_SPECIALWEEKDAY_1_VALUE         = 6
	PREF_SPECIALWEEKDAY_2_KEY           = "specialweekday_2"
	PREF_SPECIALWEEKDAY_2_VALUE         = 2
	PREF_NEWENTRY_CATEGORY_KEY          = "newentrycategory"
	PREF_NEWENTRY_CATEGORY_VALUE        = -1
	PREF_NEWENTRY_PROTECTED_KEY         = "newentryprotected"
	PREF_NEWENTRY_PROTECTED_VALUE       = false
	PREF_PASSWORD_KEY                   = "password"
	PREF_PASSWORD_VALUE                 = ""
	PREF_PASSWORD_EXPIRE_KEY            = "passwordexpire"
	PREF_PASSWORD_EXPIRE_VALUE          = 5 // min
	PREF_UPDATE_LAST_CHECK_KEY          = "lastupdatecheck"
	PREF_UPDATE_LAST_CHECK_VALUE        = 0
	PREF_UPDATE_CHECK_INTERVAL_KEY      = "updatecheckinterval"
	PREF_UPDATE_CHECK_INTERVAL_VALUE    = 48
	PREF_UPDATE_CHECK_AUTO_KEY          = "autoupdatecheck"
	PREF_UPDATE_CHECK_AUTO_VALUE        = true
	PREF_EXPORTMODE_KEY                 = "exportmode"
	PREF_EXPORMODE_VALUE                = EXPORTMODE_ASK
	PREF_MAXSEARCHRESULTS_KEY           = "maxsearchresults"
	PREF_MAXSEARCHRESULTS_VALUE_MOBILE  = 50
	PREF_MAXSEARCHRESULTS_VALUE_DESKTOP = 150
)

type Preferences struct {
	ThemeVariant        int
	FirstDayOfWeek      int
	SpecialWeekDay1     int
	SpecialWeekDay2     int
	SpecialWeekDay3     int
	NewEntryCategory    int64
	NewEntryProtected   bool
	PasswordHash        []byte
	PasswordExpire      int
	LastUpdatecheck     int64
	UpdateCheckInterval int
	AutoUpdateCheck     bool
	ExportMode          ExportModeType
	MaxSearchResults    int
}

func init() {
}

func NewPreferences() *Preferences {
	hash, _ := bcrypt.GenerateFromPassword([]byte(PREF_PASSWORD_VALUE), bcrypt.DefaultCost)
	p := &Preferences{
		ThemeVariant:        Gui.App.Preferences().IntWithFallback(PREF_THEMEVARIANT_KEY, PREF_THEMEVARIANT_VALUE),
		FirstDayOfWeek:      Gui.App.Preferences().IntWithFallback(PREF_FIRSTDAYOFWEEK_KEY, PREF_FIRSTDAYOFWEEK_VALUE),
		SpecialWeekDay1:     Gui.App.Preferences().IntWithFallback(PREF_SPECIALWEEKDAY_0_KEY, PREF_SPECIALWEEKDAY_0_VALUE),
		SpecialWeekDay2:     Gui.App.Preferences().IntWithFallback(PREF_SPECIALWEEKDAY_1_KEY, PREF_SPECIALWEEKDAY_1_VALUE),
		SpecialWeekDay3:     Gui.App.Preferences().IntWithFallback(PREF_SPECIALWEEKDAY_2_KEY, PREF_SPECIALWEEKDAY_2_VALUE),
		NewEntryCategory:    int64(Gui.App.Preferences().IntWithFallback(PREF_NEWENTRY_CATEGORY_KEY, PREF_NEWENTRY_CATEGORY_VALUE)),
		NewEntryProtected:   Gui.App.Preferences().BoolWithFallback(PREF_NEWENTRY_PROTECTED_KEY, PREF_NEWENTRY_PROTECTED_VALUE),
		PasswordHash:        []byte(Gui.App.Preferences().StringWithFallback(PREF_PASSWORD_KEY, string(hash))),
		PasswordExpire:      Gui.App.Preferences().IntWithFallback(PREF_PASSWORD_EXPIRE_KEY, PREF_PASSWORD_EXPIRE_VALUE),
		LastUpdatecheck:     100 * int64(Gui.App.Preferences().IntWithFallback(PREF_UPDATE_LAST_CHECK_KEY, PREF_UPDATE_LAST_CHECK_VALUE)),
		UpdateCheckInterval: Gui.App.Preferences().IntWithFallback(PREF_UPDATE_CHECK_INTERVAL_KEY, PREF_UPDATE_CHECK_INTERVAL_VALUE),
		AutoUpdateCheck:     Gui.App.Preferences().BoolWithFallback(PREF_UPDATE_CHECK_AUTO_KEY, PREF_UPDATE_CHECK_AUTO_VALUE),
		ExportMode:          ExportModeType(Gui.App.Preferences().IntWithFallback(PREF_EXPORTMODE_KEY, int(PREF_EXPORMODE_VALUE))),
	}
	if Gui.IsDesktop {
		p.MaxSearchResults = Gui.App.Preferences().IntWithFallback(PREF_MAXSEARCHRESULTS_KEY, PREF_MAXSEARCHRESULTS_VALUE_DESKTOP)
	} else {
		p.MaxSearchResults = Gui.App.Preferences().IntWithFallback(PREF_MAXSEARCHRESULTS_KEY, PREF_MAXSEARCHRESULTS_VALUE_MOBILE)
	}
	return p
}

func (p *Preferences) Store() {
	pref := Gui.App.Preferences()
	pref.SetInt(PREF_THEMEVARIANT_KEY, p.ThemeVariant)
	pref.SetInt(PREF_FIRSTDAYOFWEEK_KEY, p.FirstDayOfWeek)
	pref.SetInt(PREF_SPECIALWEEKDAY_0_KEY, p.SpecialWeekDay1)
	pref.SetInt(PREF_SPECIALWEEKDAY_1_KEY, p.SpecialWeekDay2)
	pref.SetInt(PREF_SPECIALWEEKDAY_2_KEY, p.SpecialWeekDay3)
	pref.SetInt(PREF_NEWENTRY_CATEGORY_KEY, int(p.NewEntryCategory))
	pref.SetBool(PREF_NEWENTRY_PROTECTED_KEY, p.NewEntryProtected)
	pref.SetString(PREF_PASSWORD_KEY, string(p.PasswordHash))
	pref.SetInt(PREF_PASSWORD_EXPIRE_KEY, p.PasswordExpire)
	pref.SetInt(PREF_UPDATE_LAST_CHECK_KEY, int(p.LastUpdatecheck/100))
	pref.SetInt(PREF_UPDATE_CHECK_INTERVAL_KEY, p.UpdateCheckInterval)
	pref.SetBool(PREF_UPDATE_CHECK_AUTO_KEY, p.AutoUpdateCheck)
	pref.SetInt(PREF_EXPORTMODE_KEY, int(p.ExportMode))
	pref.SetInt(PREF_MAXSEARCHRESULTS_KEY, p.MaxSearchResults)
}

func (p *Preferences) SavePasswd(pass []byte) {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		return
	}
	p.PasswordHash = hash
	p.Store()
}

func (p *Preferences) CheckPasswd(pass []byte) error {
	return bcrypt.CompareHashAndPassword(p.PasswordHash, []byte(pass))
}
