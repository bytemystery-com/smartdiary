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
	"strconv"
	"time"

	"bytemystery-com/smartdiary/calendarlayout"
	"bytemystery-com/smartdiary/database"
	"bytemystery-com/smartdiary/daywidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bytemystery-com/picbutton"
)

type MainView struct {
	content         fyne.CanvasObject
	calenderContent *fyne.Container
	date            time.Time
	back1m          *widget.Button
	back1y          *widget.Button
	forward1m       *widget.Button
	forward1y       *widget.Button
	label           *widget.Label
	dayWidgets      []*daywidget.DayWidget
	data            *database.EntryDataSimpleType
	lastSearch      string
}

type SwipeBox struct {
	widget.BaseWidget
	label *widget.Label
}

func NewMainView() *MainView {
	m := MainView{}
	m.calenderContent = container.New(&calendarlayout.CalendarLayout{})

	m.date = time.Now()
	m.date = time.Date(m.date.Year(), m.date.Month(), 1, 0, 0, 0, 0, time.Local)
	m.back1m = widget.NewButtonWithIcon("", theme.NewThemedResource(Gui.IconBack1m), func() {
		y := m.date.Year()
		mo := m.date.Month()
		if mo > 1 {
			mo--
		} else {
			y--
			mo = 12
		}
		m.date = time.Date(y, mo, 1, 0, 0, 0, 0, time.Local)
		m.createCalendar()
	})
	m.back1y = widget.NewButtonWithIcon("", theme.NewThemedResource(Gui.IconBack1y), func() {
		m.date = m.date.AddDate(-1, 0, 0)
		m.createCalendar()
	})
	m.forward1m = widget.NewButtonWithIcon("", theme.NewThemedResource(Gui.IconForeward1m), func() {
		y := m.date.Year()
		mo := m.date.Month()
		if mo > 11 {
			mo = 1
			y++
		} else {
			mo += 1
		}
		m.date = time.Date(y, mo, 1, 0, 0, 0, 0, time.Local)
		m.createCalendar()
	})
	m.forward1y = widget.NewButtonWithIcon("", theme.NewThemedResource(Gui.IconForeward1y), func() {
		m.date = m.date.AddDate(1, 0, 0)
		m.createCalendar()
	})
	m.label = widget.NewLabel("")
	m.label.Alignment = fyne.TextAlignCenter
	m.label.TextStyle = fyne.TextStyle{
		Bold: true,
	}
	scroll := container.NewScroll(m.calenderContent)

	search := picbutton.NewPicButton(Gui.IconSearchUp.StaticContent, Gui.IconSearchDown.StaticContent, nil, nil, false, m.DoSearch, nil)
	buttonLine := container.NewHBox(layout.NewSpacer(), search, layout.NewSpacer())

	m.content = container.NewBorder(nil, buttonLine, nil, nil, scroll)

	return &m
}

func (m *MainView) DoSearch() {
	var dia *dialog.ConfirmDialog
	text := widget.NewEntry()
	text.OnSubmitted = func(string) {
		dia.Confirm()
	}
	text.SetPlaceHolder(lang.X("search.name.placeholder", "Search text"))
	c := container.NewBorder(widget.NewLabel(lang.X("search.name", "Search text")), nil, nil, nil, text)
	dia = dialog.NewCustomConfirm(lang.X("search.title", "Search entries"), lang.X("search", "Search"), lang.X("cancel", "Cancel"),
		c, func(ok bool) {
			if !ok {
				return
			}
			m.lastSearch = text.Text
			ShowSearch(text.Text)
		}, Gui.MainWindow,
	)
	dia.Show()
	si := Gui.MainWindow.Canvas().Size()
	var windowScale float32 = 1.0
	dia.Resize(fyne.NewSize(si.Width*windowScale, dia.MinSize().Height))

	Gui.MainWindow.Canvas().Focus(text)
}

func (m *MainView) createCalendar() {
	t := time.Date(m.date.Year(), m.date.Month(), 1, 0, 0, 0, 0, time.Local)
	w := t.Weekday()
	diff := Gui.Settings.FirstDayOfWeek - int(w)
	if diff < 0 {
		diff = -diff
	} else {
		diff = 7 - diff
	}
	start := t.AddDate(0, 0, -diff)
	startDay := int(start.Weekday())
	m.dayWidgets = m.dayWidgets[:0]

	for i := range 7 {
		w := (startDay + i) % 7
		txtColor := ""
		switch w {
		case Gui.Settings.SpecialWeekDay1:
			txtColor = "txt_color_special_0"
		case Gui.Settings.SpecialWeekDay2:
			txtColor = "txt_color_special_1"
		case Gui.Settings.SpecialWeekDay3:
			txtColor = "txt_color_special_2"
		default:
			txtColor = "txt_color"

		}
		m.dayWidgets = append(m.dayWidgets, daywidget.NewDayWidget(Weekdays[w], txtColor, "bg_color_header", "bg_color_header", daywidget.ColorMode_1r, true, nil, nil))
	}
	y := t.Year()
	mo := t.Month()
	if mo > 11 {
		mo = 1
		y += 1
	} else {
		mo += 1
	}
	end := time.Date(y, mo, 1, 0, 0, 0, 0, time.Local)
	days := int(end.Sub(start).Hours() / 24)
	daysTotal := (days / 7) * 7
	if daysTotal != days {
		daysTotal += 7
	}
	d := start
	now := time.Now().Local()
	for i := range daysTotal {
		w := (startDay + i) % 7
		txtColor := ""
		switch w {
		case Gui.Settings.SpecialWeekDay1:
			if m.date.Month() == d.Month() {
				txtColor = "txt_color_special_0"
			} else {
				txtColor = "txt_color_special_0_other"
			}
		case Gui.Settings.SpecialWeekDay2:
			if m.date.Month() == d.Month() {
				txtColor = "txt_color_special_1"
			} else {
				txtColor = "txt_color_special_1_other"
			}
		case Gui.Settings.SpecialWeekDay3:
			if m.date.Month() == d.Month() {
				txtColor = "txt_color_special_2"
			} else {
				txtColor = "txt_color_special_2_other"
			}
		default:
			if m.date.Month() == d.Month() {
				txtColor = "txt_color"
			} else {
				txtColor = "txt_color_other"
			}
		}
		bold := false
		if now.Year() == d.Year() && now.Month() == d.Month() && now.Day() == d.Day() {
			bold = true
		}
		var err error
		m.data, err = Database.GetSimpleEntryByDate(d)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			UIErrorHandler(err)
		}
		color1 := "bg_color_calendar"
		color2 := "bg_color_calendar"
		colorMode := daywidget.ColorMode_1r
		var protected *fyne.StaticResource = nil
		var mark *fyne.StaticResource = nil

		if m.data != nil && m.data.Mark != 0 {
			mark = Gui.IconMarks[m.data.Mark]
		}
		if m.data != nil && m.data.Protected {
			protected = Gui.IconProtected
		}
		if err == nil {
			cat, ok := CategoryMap[m.data.CategoryId1]
			if ok {
				color1 = cat.Color
			}
			colorMode = daywidget.ColorMode_1r
			if m.data.CategoryId2.Valid {
				cat, ok := CategoryMap[m.data.CategoryId2.Int64]
				if ok {
					color2 = cat.Color
					colorMode = m.data.ColorMode
				}
			}
		}
		if color1 != "bg_color_calendar" && (txtColor == "txt_color" || txtColor == "txt_color_other") {
			if m.date.Month() == d.Month() {
				txtColor = "txt_color_bl"
			} else {
				txtColor = "txt_color_bl_other"
			}
		}
		day := daywidget.NewDayWidget(strconv.Itoa(d.Day()), txtColor, color1, color2, colorMode, bold, mark, protected)
		isProtected := (protected != nil)
		date := d
		day.OnTapped = func(pe *fyne.PointEvent) {
			if isProtected {
				doPasswordCheck(func(ok bool) {
					if ok {
						ShowEntryView(&date, false)
					}
				})
			} else {
				ShowEntryView(&date, false)
			}
		}
		m.dayWidgets = append(m.dayWidgets, day)
		d = d.AddDate(0, 0, 1)
	}
	m.label.SetText(fmt.Sprintf("%s, %d", Months[m.date.Month()-1], m.date.Year()))

	list := make([]fyne.CanvasObject, 0, len(m.dayWidgets)+5)
	list = append(list, m.back1y)
	list = append(list, m.back1m)
	list = append(list, m.forward1m)
	list = append(list, m.forward1y)
	list = append(list, m.label)

	for _, item := range m.dayWidgets {
		list = append(list, item)
	}
	m.calenderContent.Objects = list
	m.calenderContent.Refresh() // calls layout()
}

func (m *MainView) SetDate(date *time.Time) {
	if date != nil {
		m.date = *date
	}
	m.createCalendar()
}

func (m *MainView) GetContent() fyne.CanvasObject {
	return m.content
}

func (m *MainView) UpdateToolBar() {
	Gui.Toolbar.Items = []widget.ToolbarItem{Gui.toolToggleThema, widget.NewToolbarSeparator(), Gui.toolSettings, widget.NewToolbarSeparator(), Gui.toolExport, Gui.toolImport, widget.NewToolbarSpacer(), Gui.toolInfo}
	Gui.Toolbar.Refresh()
}

func (m *MainView) ThemeChanged() {
	m.content.Refresh()
}
