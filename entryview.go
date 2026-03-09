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
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"myapp/database"
	"myapp/daywidget"
	"myapp/entrylayout"
	"myapp/menubutton"
	"myapp/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bytemystery-com/colorlabel"
)

type EntryView struct {
	content        fyne.CanvasObject
	date           time.Time
	back1d         *widget.Button
	back1w         *widget.Button
	forward1d      *widget.Button
	forward1w      *widget.Button
	label          *widget.Label
	entry          *widget.Entry
	data           *database.EntryDataType
	category1      *colorlabel.ColorLabel
	category2      *colorlabel.ColorLabel
	colorModeBtn   *menubutton.MenuButton
	colorModeIcons []*fyne.StaticResource
	markBtn        *menubutton.MenuButton
	markIcons      []*fyne.StaticResource
	protected      *widget.Check
	fromSearch     bool
}

func NewEntryView() *EntryView {
	e := EntryView{}

	e.date = time.Now()
	e.date = time.Date(e.date.Year(), e.date.Month(), 1, 0, 0, 0, 0, time.Local)
	e.back1d = widget.NewButtonWithIcon("", theme.NewThemedResource(Gui.IconBack1m), func() {
		e.doSave()
		e.checkPassword(e.date.AddDate(0, 0, -1), func(ok bool) {
			if ok {
				e.date = e.date.AddDate(0, 0, -1)
				e.updateData()
				e.SetCursorAtEnd()
			}
		})
	})
	e.back1w = widget.NewButtonWithIcon("", theme.NewThemedResource(Gui.IconBack1y), func() {
		e.doSave()
		e.checkPassword(e.date.AddDate(0, 0, -7), func(ok bool) {
			if ok {
				e.date = e.date.AddDate(0, 0, -7)
				e.updateData()
				e.SetCursorAtEnd()
			}
		})
	})
	e.forward1d = widget.NewButtonWithIcon("", theme.NewThemedResource(Gui.IconForeward1m), func() {
		e.doSave()
		e.checkPassword(e.date.AddDate(0, 0, 1), func(ok bool) {
			if ok {
				e.date = e.date.AddDate(0, 0, 1)
				e.updateData()
				e.SetCursorAtEnd()
			}
		})
	})
	e.forward1w = widget.NewButtonWithIcon("", theme.NewThemedResource(Gui.IconForeward1y), func() {
		e.doSave()
		e.checkPassword(e.date.AddDate(0, 0, 7), func(ok bool) {
			if ok {
				e.date = e.date.AddDate(0, 0, 7)
				e.updateData()
				e.SetCursorAtEnd()
			}
		})
	})
	e.label = widget.NewLabel("")
	e.label.Alignment = fyne.TextAlignCenter
	e.label.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	e.entry = widget.NewMultiLineEntry()
	e.entry.SetMinRowsVisible(1)
	e.entry.OnChanged = func(str string) {
		e.entry.SetMinRowsVisible(e.rowsNeeded(str))
	}
	e.entry.Wrapping = fyne.TextWrapWord
	top := container.New(&entrylayout.EntryLayout{}, e.back1w, e.back1d, e.forward1d, e.forward1w, e.label, e.entry)

	cancel := widget.NewButton(lang.X("cancel", "Cancel"), e.doCancel)
	ok := widget.NewButton(lang.X("save", "Save"), e.doOk)
	ok.Importance = widget.HighImportance
	wOk := util.GetDefaultTextWidth(ok.Text + "XXX")
	wCancel := util.GetDefaultTextWidth(cancel.Text + "XXX")
	w := wCancel
	if wOk > wCancel {
		w = wCancel
	}
	btnSize := fyne.NewSize(w, ok.MinSize().Height)
	okC := container.NewGridWrap(btnSize, ok)
	cancelC := container.NewGridWrap(btnSize, cancel)
	buttonLine := container.NewHBox(layout.NewSpacer(), container.NewCenter(cancelC), container.NewCenter(okC), layout.NewSpacer())
	var cat1C fyne.CanvasObject
	var cat2C fyne.CanvasObject
	e.category1, cat1C = e.createCategorySelector(false)
	e.category2, cat2C = e.createCategorySelector(true)
	e.colorModeIcons = make([]*fyne.StaticResource, 0, 5)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeSimple)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeTopBottom)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeLeftRight)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeLeftTopRightBottom)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeRightTopLeftBottom)
	e.colorModeBtn = menubutton.NewMenuButton("", e.colorModeIcons[0], func(ev *fyne.PointEvent) {
		if !e.colorModeBtn.Disabled() {
			c := container.NewVBox()
			colorModeIndex := 0
			if e.data != nil {
				colorModeIndex = daywidget.ColorModeToIndex[e.data.ColorMode]
			}
			var popup *widget.PopUp
			for index, item := range e.colorModeIcons {
				ctrl := widget.NewButtonWithIcon("", item, func() {
					popup.Hide()
					if e.data != nil {
						e.data.ColorMode = daywidget.IndexToColorMode[index]
					}
					e.colorModeBtn.SetIcon(e.colorModeIcons[index])
				})
				ctrl.Importance = widget.LowImportance
				icon := widget.NewIcon(theme.RadioButtonIcon())
				if colorModeIndex == index {
					icon.SetResource(theme.RadioButtonFillIcon())
				}
				c.Add(container.NewHBox(icon, ctrl))
			}
			popup = ShowPopUp(ev, c)
		}
	})
	e.markIcons = make([]*fyne.StaticResource, 0, len(Gui.IconMarks))
	for index, item := range Gui.IconMarks {
		if index > 0 {
			e.markIcons = append(e.markIcons, item)
		} else {
			r := theme.VisibilityOffIcon()
			res := fyne.StaticResource{
				StaticName:    r.Name(),
				StaticContent: r.Content(),
			}
			e.markIcons = append(e.markIcons, &res)
		}
	}
	e.markBtn = menubutton.NewMenuButton("", e.markIcons[0], func(ev *fyne.PointEvent) {
		c := container.NewVBox()
		markIndex := 0
		if e.data != nil {
			markIndex = e.data.Mark
		}
		var popup *widget.PopUp
		var hBox *fyne.Container = nil
		for index, item := range e.markIcons {
			if index%2 == 0 {
				if hBox != nil {
					c.Add(hBox)
				}
				hBox = container.NewHBox()
			}
			ctrl := widget.NewButtonWithIcon("", theme.NewThemedResource(item), func() {
				popup.Hide()
				if e.data != nil {
					e.data.Mark = index
				}
				e.markBtn.SetIcon(theme.NewThemedResource(e.markIcons[index]))
			})
			ctrl.Importance = widget.LowImportance
			icon := widget.NewIcon(theme.RadioButtonIcon())
			if markIndex == index {
				icon.SetResource(theme.RadioButtonFillIcon())
			}
			hBox.Add(container.NewHBox(icon, ctrl))
			if index%2 == 0 {
				hBox.Add(util.NewHFiller(3))
			}
		}
		if hBox != nil {
			c.Add(hBox)
		}
		popup = ShowPopUp(ev, c)
	})
	e.protected = widget.NewCheck(lang.X("entry.protected", "Protected"), func(b bool) {
		if e.data != nil {
			e.data.Protected = b
		}
	})
	top2 := container.NewVBox(top, container.NewHBox(layout.NewSpacer(), cat1C, cat2C, layout.NewSpacer()))
	top3 := container.NewHBox(widget.NewLabel(lang.X("entry.mark", "Mark")), e.markBtn, layout.NewSpacer(), e.colorModeBtn, layout.NewSpacer(), e.protected)
	if Gui.IsDesktop {
		e.content = container.NewVBox(top2, top3, layout.NewSpacer(), buttonLine)
	} else {
		e.content = container.NewVBox(top2, top3, util.NewVFiller(util.GetDefaultTextHeight("X")*.5), buttonLine)
	}

	return &e
}

func (e *EntryView) checkPassword(date time.Time, doIt func(ok bool)) {
	data, err := Database.GetEntryByDate(date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			if doIt != nil {
				doIt(true)
			}
			return
		}
		UIErrorHandler(err)
		if doIt != nil {
			doIt(false)
		}
		return
	}
	if !data.Protected {
		if doIt != nil {
			doIt(true)
		}
		return
	}
	doPasswordCheck(doIt)
}

func (e *EntryView) rowsNeeded(str string) int {
	n := strings.Count(str, "\n")
	lines := strings.Split(str, "\n")
	w := e.entry.Size().Width - 2*theme.InnerPadding()
	if w > 0 {
		for _, line := range lines {
			x := util.GetDefaultTextWidth(line) / w
			_ = x
			m := math.Ceil(float64(util.GetDefaultTextWidth(line) / w))
			if m > 1 {
				n += int(m) - 1
			}
		}
	}

	return n + 1
}

func (e *EntryView) createCategorySelector(hasNull bool) (*colorlabel.ColorLabel, fyne.CanvasObject) {
	category, content := CreateCategorySelector(hasNull, func(id int64) {
		if hasNull {
			if id == -1 {
				e.resetCategory2UI()
				e.data.CategoryId2.Valid = false
				e.data.ColorMode = daywidget.ColorMode_1r
				e.colorModeBtn.Disable()
			} else {
				e.data.CategoryId2.Valid = true
				e.data.CategoryId2.Int64 = id
				e.colorModeBtn.Enable()
			}
		} else {
			if e.data != nil {
				e.data.CategoryId1 = id
			}
		}
	}, func(id int64) bool {
		if hasNull {
			if e.data != nil && ((e.data.CategoryId2.Valid && e.data.CategoryId2.Int64 == id) ||
				(!e.data.CategoryId2.Valid && id == -1)) {
				return true
			}
		} else if e.data != nil && e.data.CategoryId1 == id {
			return true
		}
		return false
	}, false, false)

	return category, content
}

func (e *EntryView) doOk() {
	e.doSave()
	if e.fromSearch {
		ShowSearch("")
	} else {
		ShowMainView(nil)
	}
}

func (e *EntryView) doSave() {
	if e.data != nil {
		e.data.Data = e.entry.Text
		e.data.Date = e.date
		err := Database.InsertOrUpdateEntry(e.data)
		if err != nil {
			UIErrorHandler(err)
		}
	}
}

func (e *EntryView) doCancel() {
	if e.fromSearch {
		ShowSearch("")
	} else {
		ShowMainView(nil)
	}
}

func (e *EntryView) resetCategory2UI() {
	e.category2.SetText(lang.X("category_not_set", "not defined"))
	e.category2.SetTextColor(Gui.Theme.GetSpecialColor("nocategory_text_color"))
	e.category2.SetBackgroundColor(theme.ColorNameBackground)
	e.colorModeBtn.SetIcon(e.colorModeIcons[0])
	e.colorModeBtn.Disable()
}

func (e *EntryView) updateData() {
	var err error

	e.label.SetText(fmt.Sprintf("%s, %d.%d.%d", Weekdays[e.date.Weekday()], e.date.Day(), e.date.Month(), e.date.Year()))
	e.data, err = Database.GetEntryByDate(e.date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// New entry
			e.colorModeBtn.SetIcon(e.colorModeIcons[0])
			e.protected.SetChecked(Gui.Settings.NewEntryProtected)
			e.entry.SetText("")
			e.entry.SetMinRowsVisible(1)
			e.data = &database.EntryDataType{}
			e.data.Id = -1
			catData, ok := CategoryMap[Gui.Settings.NewEntryCategory]
			if !ok {
				catData = CategoryData[0]
			}
			e.data.CategoryId1 = catData.Id
			e.category1.SetText(catData.Name)
			e.category1.SetBackgroundColor(Gui.Theme.GetSpecialColor(catData.Color))
			e.markBtn.SetIcon(theme.NewThemedResource(e.markIcons[0]))
			e.data.CategoryId2.Valid = false
			e.resetCategory2UI()
		} else {
			UIErrorHandler(err)
		}
		e.UpdateRemoveEntryInToolBar()
		return
	}
	e.entry.SetText(e.data.Data)
	e.entry.SetMinRowsVisible(e.rowsNeeded(e.data.Data))
	catData1, ok := CategoryMap[e.data.CategoryId1]
	if ok {
		e.category1.SetText(catData1.Name)
		e.category1.SetBackgroundColor(Gui.Theme.GetSpecialColor(catData1.Color))
	}
	if e.data.CategoryId2.Valid {
		catData2, ok := CategoryMap[e.data.CategoryId2.Int64]
		if ok {
			e.category2.SetText(catData2.Name)
			e.category2.SetTextColor(Gui.Theme.GetSpecialColor("category_text_color"))
			e.category2.SetBackgroundColor(Gui.Theme.GetSpecialColor(catData2.Color))
			index, ok := daywidget.ColorModeToIndex[e.data.ColorMode]
			if ok {
				e.colorModeBtn.SetIcon(e.colorModeIcons[index])
			}
			e.colorModeBtn.Enable()
		}
	} else {
		e.resetCategory2UI()
	}
	e.markBtn.SetIcon(theme.NewThemedResource(e.markIcons[e.data.Mark]))
	e.protected.SetChecked(e.data.Protected)
	e.UpdateRemoveEntryInToolBar()
}

func (e *EntryView) SetDate(date *time.Time, fromSearch bool) {
	if date == nil {
		return
	}
	e.date = *date
	e.updateData()
	e.fromSearch = fromSearch
}

func (e *EntryView) SetCursorAtEnd() {
	Gui.MainWindow.Canvas().Focus(e.entry)
	n := len(e.entry.Text)
	if n > 0 {
		e.entry.CursorColumn = len(e.entry.Text)
		e.entry.Refresh()
	}
	e.entry.SetMinRowsVisible(e.rowsNeeded(e.entry.Text))
}

func (e *EntryView) RemoveEntry() {
	if e.data.Id > 0 {
		err := Database.DeleteEntry(e.data.Id)
		if err != nil {
			UIErrorHandler(err)
			return
		}
		ShowMainView(nil)
	}
}

func (e *EntryView) UpdateRemoveEntryInToolBar() {
	if e.data.Id > 0 {
		Gui.toolRemove.Enable()
	} else {
		Gui.toolRemove.Disable()
	}
}

func (e *EntryView) GetContent() fyne.CanvasObject {
	return e.content
}

func (e *EntryView) UpdateToolBar() {
	Gui.Toolbar.Items = []widget.ToolbarItem{Gui.toolToggleThema, widget.NewToolbarSeparator(), Gui.toolRemove, widget.NewToolbarSpacer(), Gui.toolInfo}
	Gui.Toolbar.Refresh()
}

func (e *EntryView) ThemeChanged() {
	e.colorModeIcons = make([]*fyne.StaticResource, 0, 5)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeSimple)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeTopBottom)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeLeftRight)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeLeftTopRightBottom)
	e.colorModeIcons = append(e.colorModeIcons, Gui.IconColorModeRightTopLeftBottom)
	if e.data != nil {
		if !e.data.CategoryId2.Valid {
			e.colorModeBtn.SetIcon(e.colorModeIcons[0])
			e.category2.SetTextColor(Gui.Theme.GetSpecialColor("nocategory_text_color"))
			e.category2.SetBackgroundColor(theme.ColorNameBackground)
		} else {
			index := daywidget.ColorModeToIndex[e.data.ColorMode]
			e.colorModeBtn.SetIcon(e.colorModeIcons[index])
		}
	}
	if e.data != nil {
		e.markBtn.SetIcon(theme.NewThemedResource(e.markIcons[e.data.Mark]))
	} else {
		e.markBtn.SetIcon(theme.NewThemedResource(e.markIcons[0]))
	}
	e.content.Refresh()
}
