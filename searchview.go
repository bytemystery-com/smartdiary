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
	"fmt"
	"image"
	"strings"

	"myapp/database"
	"myapp/daywidget"
	"myapp/searchlayout"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bytemystery-com/picbutton"
)

type SearchView struct {
	content    fyne.CanvasObject
	list       *widget.List
	data       []*database.EntryDataType
	scroll     *container.Scroll
	vBox       *fyne.Container
	lastSearch string
}

func NewSearchView() *SearchView {
	s := SearchView{}
	s.vBox = container.NewVBox()
	s.scroll = container.NewScroll(s.vBox)

	back := picbutton.NewPicButton(Gui.IconBackUp.StaticContent, Gui.IconBackDown.StaticContent, nil, nil, false, func() {
		ShowMainView(nil)
	}, nil)
	buttonLine := container.NewHBox(container.NewCenter(back))

	s.content = container.NewBorder(nil, buttonLine, nil, nil, s.scroll)

	return &s
}

func (s *SearchView) findSearchString(str, search string) [][2]int {
	var indices [][2]int
	from := 0
	search = strings.ToLower(search)
	str = strings.ToLower(str)
	for {
		i := strings.Index(str[from:], search)
		if i == -1 {
			break
		}
		start := from + i
		end := start + len(search)
		indices = append(indices, [2]int{start, end})
		from = end
	}
	return indices
}

func (s *SearchView) buildTextSegments(str, search string) []widget.RichTextSegment {
	searchMarks := s.findSearchString(str, search)
	start := 0
	list := make([]widget.RichTextSegment, 0, len(searchMarks)*2+1)
	for _, item := range searchMarks {
		if item[0] > start {
			txt := str[start:item[0]]
			w := widget.TextSegment{
				Text: txt,
				Style: widget.RichTextStyle{
					ColorName: theme.ColorNameForeground,
					Inline:    true,
					TextStyle: fyne.TextStyle{
						Bold:   false,
						Italic: false,
					},
				},
			}
			list = append(list, &w)
		}

		txt := str[item[0]:item[1]]
		w := widget.TextSegment{
			Text: txt,
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameError,
				Inline:    true,
				TextStyle: fyne.TextStyle{
					Bold:   true,
					Italic: true,
				},
			},
		}
		list = append(list, &w)
		start = item[1]
	}
	if start < len(str) {
		w := widget.TextSegment{
			Text: str[start:],
			Style: widget.RichTextStyle{
				ColorName: theme.ColorNameForeground,
				Inline:    true,
				TextStyle: fyne.TextStyle{
					Bold:   false,
					Italic: false,
				},
			},
		}
		list = append(list, &w)
	}
	return list
}

func (s *SearchView) updateContent() {
	s.scroll.ScrollToTop()
	s.vBox.RemoveAll()
	list := make([]fyne.CanvasObject, 0, len(s.data))
	for _, item := range s.data {
		img := canvas.NewRaster(func(w, h int) image.Image {
			return image.NewRGBA(image.Rect(0, 0, w, h))
		})
		s.createIcon(img, item)

		text1 := canvas.NewText(fmt.Sprintf("%s, %d %s %d ",
			Weekdays[item.Date.Weekday()], item.Date.Day(), Months[item.Date.Month()-1], item.Date.Year()), theme.Color(theme.ColorNameForeground))

		text2 := widget.NewRichText(s.buildTextSegments(item.Data, s.lastSearch)...)
		text2.Wrapping = fyne.TextWrapWord
		// text2.Refresh()

		btn := widget.NewButtonWithIcon("", theme.MoreHorizontalIcon(), func() {
			if item.Protected {
				doPasswordCheck(func(ok bool) {
					if ok {
						ShowEntryView(&item.Date, true)
					}
				})
			} else {
				ShowEntryView(&item.Date, true)
			}
		})
		s := btn.MinSize().Height
		img.SetMinSize(fyne.NewSize(s, s))

		// Overlay SVG (Raute)
		var mark *canvas.Image
		var protected *canvas.Image
		mark = canvas.NewImageFromResource(Gui.IconMarks[item.Mark])
		mark.FillMode = canvas.ImageFillContain
		mark.Refresh()
		if item.Mark != 0 {
			mark.Show()
		} else {
			mark.Hide()
		}
		protected = canvas.NewImageFromResource(Gui.IconProtected)
		protected.FillMode = canvas.ImageFillContain
		protected.Refresh()
		if item.Protected {
			protected.Show()
		} else {
			protected.Hide()
		}
		c := container.New(&searchlayout.SearchLayout{}, img, text1, btn, text2, mark, protected)
		list = append(list, c)
	}
	s.vBox.Objects = list
	s.vBox.Refresh()
	s.scroll.Refresh()
}

func (s *SearchView) createIcon(icon *canvas.Raster, item *database.EntryDataType) {
	icon.Generator = func(w, h int) image.Image {
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		var cat1 *database.CategoryDataType
		var cat2 *database.CategoryDataType
		cat1, _ = CategoryMap[item.CategoryId1]
		if item.CategoryId2.Valid {
			cat2, _ = CategoryMap[item.CategoryId2.Int64]
		}
		switch item.ColorMode {
		case daywidget.ColorMode_1r:
			if cat1 != nil {
				col := Gui.Theme.GetSpecialColor(cat1.Color)
				for y := range h {
					for x := range w {
						img.Set(x, y, col)
					}
				}
			}
		case daywidget.ColorMode_2rlr:
			if cat1 != nil && cat2 != nil {
				col1 := Gui.Theme.GetSpecialColor(cat1.Color)
				col2 := Gui.Theme.GetSpecialColor(cat2.Color)
				for y := range h {
					for x := range w {
						if x < w/2 {
							img.Set(x, y, col1)
						} else {
							img.Set(x, y, col2)
						}
					}
				}
			}
		case daywidget.ColorMode_2rtb:
			if cat1 != nil && cat2 != nil {
				col1 := Gui.Theme.GetSpecialColor(cat1.Color)
				col2 := Gui.Theme.GetSpecialColor(cat2.Color)
				for y := range h {
					for x := range w {
						if y < h/2 {
							img.Set(x, y, col1)
						} else {
							img.Set(x, y, col2)
						}
					}
				}
			}
		case daywidget.ColorMode_2tltrb:
			if cat1 != nil && cat2 != nil {
				col1 := Gui.Theme.GetSpecialColor(cat1.Color)
				col2 := Gui.Theme.GetSpecialColor(cat2.Color)
				m := float64(h) / float64(w)
				for y := range h {
					for x := range w {
						y1 := float64(h) - m*float64(x)
						if float64(y) < y1 {
							img.Set(x, y, col1)
						} else {
							img.Set(x, y, col2)
						}
					}
				}
			}
		case daywidget.ColorMode_2trtlb:
			if cat1 != nil && cat2 != nil {
				col1 := Gui.Theme.GetSpecialColor(cat1.Color)
				col2 := Gui.Theme.GetSpecialColor(cat2.Color)
				m := float64(h) / float64(w)
				for y := range h {
					for x := range w {
						y1 := m * float64(x)
						if float64(y) < y1 {
							img.Set(x, y, col1)
						} else {
							img.Set(x, y, col2)
						}
					}
				}
			}
		}

		return img
	}
}

func (s *SearchView) Search(search string) bool {
	setSearchView := true
	if search == "" {
		search = s.lastSearch
	}
	data, err := Database.Search("%" + search + "%")
	if err != nil {
		UIErrorHandler(err)
		return setSearchView
	}
	hasProtected := false
	for _, item := range data {
		if item.Protected {
			hasProtected = true
			break
		}
	}
	load := func() {
		s.data = data
		s.lastSearch = search
		s.updateContent()
	}

	if hasProtected {
		setSearchView = false
		s.data = s.data[:0]
		s.updateContent()

		doPasswordCheck(func(ok bool) {
			if ok {
				SetSearchView()
				load()
			}
		})
	} else {
		load()
	}
	return setSearchView
}

func (s *SearchView) GetContent() fyne.CanvasObject {
	return s.content
}

func (s *SearchView) UpdateToolBar() {
	Gui.Toolbar.Items = []widget.ToolbarItem{Gui.toolToggleThema, widget.NewToolbarSpacer(), Gui.toolInfo}
	Gui.Toolbar.Refresh()
}

func (s *SearchView) ThemeChanged() {
	for _, item := range s.vBox.Objects {
		c, ok := item.(*fyne.Container)
		if !ok {
			return
		}
		txt, ok := c.Objects[1].(*canvas.Text)
		if !ok {
			return
		}
		txt.Color = theme.Color(theme.ColorNameForeground)
		txt.Refresh()
	}
	s.content.Refresh()
}
