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

import (
	"strconv"

	"bytemystery-com/smartdiary/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bytemystery-com/colorlabel"
)

type SettingsView struct {
	content        fyne.CanvasObject
	weekStart      *widget.Select
	specialDay1    *widget.Select
	specialDay2    *widget.Select
	specialDay3    *widget.Select
	category       *colorlabel.ColorLabel
	defCategory    int64
	defProtected   *widget.Check
	passwordExpire *widget.Entry
	categoryEdit   *colorlabel.ColorLabel
	editCategory   int64
	catEdit        *widget.Entry
	exportMode     *widget.Select
	maxResults     *widget.Entry
}

func NewSettingsView() *SettingsView {
	s := SettingsView{}
	weekdays := make([]string, 0, len(Weekdays)+1)
	weekdays = append(weekdays, lang.X("weekday_no", "---"))
	weekdays = append(weekdays, Weekdays...)

	s.weekStart = widget.NewSelect(Weekdays, nil)
	s.specialDay1 = widget.NewSelect(weekdays, nil)
	s.specialDay2 = widget.NewSelect(weekdays, nil)
	s.specialDay3 = widget.NewSelect(weekdays, nil)

	var catContent fyne.CanvasObject
	s.category, catContent = CreateCategorySelector(nil, false, func(id int64) {
		s.defCategory = id
	}, func(id int64) bool {
		return id == s.defCategory
	}, false, false)
	s.defProtected = widget.NewCheck("", nil)

	changePassword := widget.NewButtonWithIcon(lang.X("settings.password.change", "Change"), theme.AccountIcon(), func() {
		ShowPasswordDialog(true, nil)
	})

	s.passwordExpire = widget.NewEntry()
	s.passwordExpire.SetPlaceHolder(lang.X("settings.password_expire.placeholder", "Duration in minutes"))
	s.passwordExpire.OnChanged = util.GetNumberFilter(s.passwordExpire, nil)

	s.catEdit = widget.NewEntry()
	s.catEdit.SetPlaceHolder(lang.X("settings.category_name.placeholder", "New name [Press ENTER]"))
	s.catEdit.OnSubmitted = func(str string) {
		cat := CategoryMap[s.editCategory]
		if cat != nil {
			cat.Name = str
			cat.IsEdited = true
		}
		s.categoryEdit.SetText(str)
	}
	var catContentEdit fyne.CanvasObject
	s.categoryEdit, catContentEdit = CreateCategorySelector(nil, false, func(id int64) {
		s.editCategory = id
		cat := CategoryMap[s.editCategory]
		if cat != nil {
			s.catEdit.SetText(cat.Name)
		}
	}, func(id int64) bool {
		return id == s.editCategory
	}, true, true)

	s.maxResults = widget.NewEntry()
	s.maxResults.SetPlaceHolder(lang.X("Max search results", "Max search results"))
	s.maxResults.OnChanged = util.GetNumberFilter(s.maxResults, nil)

	form := container.New(layout.NewFormLayout(),
		widget.NewLabel(lang.X("settings.week_start", "Start of week")), s.weekStart,
		colorlabel.NewColorLabel(lang.X("settings.week_first_special", "First special day"), Gui.Theme.GetSpecialColor("txt_color_special_0"), nil, 1.0), s.specialDay1,
		colorlabel.NewColorLabel(lang.X("settings.week_second_special", "Second special day"), Gui.Theme.GetSpecialColor("txt_color_special_1"), nil, 1.0), s.specialDay2,
		colorlabel.NewColorLabel(lang.X("settings.week_third_special", "Third special day"), Gui.Theme.GetSpecialColor("txt_color_special_2"), nil, 1.0), s.specialDay3,
		widget.NewLabel(lang.X("settings.default_category", "Default category")), catContent,
		widget.NewLabel(lang.X("settings.default_protected", "Default protected")), s.defProtected,
		widget.NewLabel(lang.X("settings.change_password", "Change password")), changePassword,
		widget.NewLabel(lang.X("settings.password_expire", "Password expire\n[min]")), s.passwordExpire,
		widget.NewLabel(lang.X("settings.max_search_results", "Max. results")), s.maxResults,
		catContentEdit, s.catEdit,
	)
	labelExportMode := widget.NewLabel(lang.X("settings.exportmode", "Export file mode"))

	s.exportMode = widget.NewSelect([]string{
		lang.X("settings.exportmode.new", "New file"),
		lang.X("settings.exportmode.overwrite", "Overwrite existing file"),
		lang.X("settings.exportmode.ask", "Ask"),
	}, nil)

	if !Gui.IsDesktop {
		form.Add(labelExportMode)
		form.Add(s.exportMode)
	}

	updateCheck := widget.NewButton(lang.X("settings.checkupdate", "Check for update"), func() {
		CheckForUpdate(false)
	})

	ok := widget.NewButton(lang.X("save", "Save"), func() { s.doSave() })
	ok.Importance = widget.HighImportance
	cancel := widget.NewButton(lang.X("cancel", "Cancel"), func() { s.doCancel() })

	wOk := util.GetDefaultTextWidth(ok.Text + "XXX")
	wCancel := util.GetDefaultTextWidth(cancel.Text + "XXX")
	w := wCancel
	if wOk > wCancel {
		w = wCancel
	}
	btnSize := fyne.NewSize(w, ok.MinSize().Height)
	okC := container.NewGridWrap(btnSize, ok)
	cancelC := container.NewGridWrap(btnSize, cancel)

	buttonLine := container.NewHBox(layout.NewSpacer(), cancelC, okC, layout.NewSpacer())
	if Gui.IsDesktop {
		c := container.NewScroll(container.NewVBox(form, util.NewVFiller(1), updateCheck))
		s.content = container.NewBorder(nil, buttonLine, nil, nil, c)
	} else {
		c := container.NewScroll(container.NewVBox(form, util.NewVFiller(1), updateCheck, util.NewVFiller(1), buttonLine))
		s.content = container.NewBorder(nil, nil, nil, nil, c)
	}
	return &s
}

func (s *SettingsView) doSave() {
	Gui.Settings.FirstDayOfWeek = s.weekStart.SelectedIndex()
	Gui.Settings.SpecialWeekDay1 = s.specialDay1.SelectedIndex() - 1
	Gui.Settings.SpecialWeekDay2 = s.specialDay2.SelectedIndex() - 1
	Gui.Settings.SpecialWeekDay3 = s.specialDay3.SelectedIndex() - 1
	Gui.Settings.NewEntryCategory = s.defCategory
	Gui.Settings.NewEntryProtected = s.defProtected.Checked
	val, err := strconv.Atoi(s.passwordExpire.Text)
	if err == nil {
		Gui.Settings.PasswordExpire = val
	}
	exportMode := EXPORTMODE_ASK
	index := s.exportMode.SelectedIndex()
	switch index {
	case 0:
		exportMode = EXPORTMODE_NEW
	case 1:
		exportMode = EXPORTMODE_OVERWRITE
	case 2:
		exportMode = EXPORTMODE_ASK
	}
	Gui.Settings.ExportMode = exportMode
	v, err := strconv.Atoi(s.maxResults.Text)
	if err == nil {
		Gui.Settings.MaxSearchResults = v
	}

	Gui.Settings.Store()

	for _, item := range CategoryData {
		if item.IsEdited {
			err := Database.UpdateCategory(item)
			if err != nil {
				UIErrorHandler(err)
			}
		}
	}
	UpdateCategoryData()
	Reload()
	RestoreBeforeSettings()
}

func (s *SettingsView) doCancel() {
	UpdateCategoryData()
	RestoreBeforeSettings()
}

func (s *SettingsView) Init() {
	s.weekStart.SetSelectedIndex(Gui.Settings.FirstDayOfWeek)
	s.specialDay1.SetSelectedIndex(Gui.Settings.SpecialWeekDay1 + 1)
	s.specialDay2.SetSelectedIndex(Gui.Settings.SpecialWeekDay2 + 1)
	s.specialDay3.SetSelectedIndex(Gui.Settings.SpecialWeekDay3 + 1)
	s.defCategory = Gui.Settings.NewEntryCategory
	cat, ok := CategoryMap[s.defCategory]
	if !ok {
		cat = CategoryData[0]
	}
	s.category.SetText(cat.Name)
	s.category.SetBackgroundColor(Gui.Theme.GetSpecialColor(cat.Color))
	s.defProtected.SetChecked(Gui.Settings.NewEntryProtected)
	s.passwordExpire.Text = strconv.Itoa(Gui.Settings.PasswordExpire)
	s.categoryEdit.SetText(CategoryData[0].Name)
	s.categoryEdit.SetBackgroundColor(Gui.Theme.GetSpecialColor(CategoryData[0].Color))
	s.editCategory = CategoryData[0].Id
	s.catEdit.SetText(CategoryData[0].Name)

	selIndex := 0
	switch Gui.Settings.ExportMode {
	case EXPORTMODE_NEW:
		selIndex = 0
	case EXPORTMODE_OVERWRITE:
		selIndex = 1
	case EXPORTMODE_ASK:
		selIndex = 2
	}
	s.exportMode.SetSelectedIndex(selIndex)
	s.maxResults.Text = strconv.Itoa(Gui.Settings.MaxSearchResults)
}

func (s *SettingsView) GetContent() fyne.CanvasObject {
	return s.content
}

func (s *SettingsView) UpdateToolBar() {
	Gui.Toolbar.Items = []widget.ToolbarItem{Gui.toolToggleThema, widget.NewToolbarSpacer(), Gui.toolInfo}
	Gui.Toolbar.Refresh()
}

func (s *SettingsView) ThemeChanged() {
	s.content.Refresh()
}
