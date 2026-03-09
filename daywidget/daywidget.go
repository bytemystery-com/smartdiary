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

package daywidget

import (
	"image/color"

	"bytemystery-com/smartdiary/mytheme"
	"bytemystery-com/smartdiary/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const OVERLAY_SIZE = 0.3

type ColorModeType int

const (
	ColorMode_1r ColorModeType = iota
	ColorMode_2rtb
	ColorMode_2rlr
	ColorMode_2tltrb
	ColorMode_2trtlb
)

var (
	IndexToColorMode map[int]ColorModeType
	ColorModeToIndex map[ColorModeType]int
)

func init() {
	IndexToColorMode = map[int]ColorModeType{0: ColorMode_1r, 1: ColorMode_2rtb, 2: ColorMode_2rlr, 3: ColorMode_2tltrb, 4: ColorMode_2trtlb}
	ColorModeToIndex = map[ColorModeType]int{ColorMode_1r: 0, ColorMode_2rtb: 1, ColorMode_2rlr: 2, ColorMode_2tltrb: 3, ColorMode_2trtlb: 4}
}

type DayWidget struct {
	widget.BaseWidget
	text       string
	txtColor   string
	backColor1 string
	backColor2 string
	bold       bool
	protected  *fyne.StaticResource
	mark       *fyne.StaticResource
	colorMode  ColorModeType
	OnTapped   func(*fyne.PointEvent)
}

type DayWidgetRenderer struct {
	w         *DayWidget
	text      *canvas.Text
	rect1     *canvas.Rectangle
	rect2     *canvas.Rectangle
	tri       *canvas.Raster
	protected *canvas.Image
	mark      *canvas.Image
}

var (
	_ fyne.Widget   = (*DayWidget)(nil)
	_ fyne.Tappable = (*DayWidget)(nil)
	/*
		_ fyne.SecondaryTappable = (*PicButton)(nil)
		_ desktop.Mouseable      = (*PicButton)(nil)
		_ desktop.Hoverable      = (*PicButton)(nil)
		_ desktop.Cursorable     = (*PicButton)(nil)
		_ fyne.Focusable         = (*PicButton)(nil)
	*/

	_ fyne.WidgetRenderer = (*DayWidgetRenderer)(nil)
)

func NewDayWidget(t string, txtColor, backColor1, backColor2 string, colorMode ColorModeType, bold bool, mark *fyne.StaticResource, protected *fyne.StaticResource) *DayWidget {
	d := DayWidget{
		text:       t,
		txtColor:   txtColor,
		backColor1: backColor1,
		backColor2: backColor2,
		bold:       bold,
		colorMode:  colorMode,
		protected:  protected,
		mark:       mark,
	}
	d.ExtendBaseWidget(&d)
	return &d
}

func (d *DayWidget) SetText(text string) {
	d.text = text
	d.Refresh()
}

func (d *DayWidget) GetOverlayScaleFactor() float32 {
	t := fyne.CurrentApp().Settings().Theme()
	myt, ok := t.(mytheme.MyThemeInterface)
	if ok {
		return myt.GetSpecialSize("overlay_size_scale")
	} else {
		return OVERLAY_SIZE
	}
}

func (d *DayWidget) GetTextColor() color.Color {
	t := fyne.CurrentApp().Settings().Theme()
	myt, ok := t.(mytheme.MyThemeInterface)
	if ok {
		return myt.GetSpecialColor(d.txtColor)
	} else {
		return t.Color(fyne.ThemeColorName(d.text), fyne.CurrentApp().Settings().ThemeVariant())
	}
}

func (d *DayWidget) GetSecondColor() color.Color {
	t := fyne.CurrentApp().Settings().Theme()
	myt, ok := t.(mytheme.MyThemeInterface)
	if ok {
		return myt.GetSpecialColor(d.backColor2)
	} else {
		return t.Color(fyne.ThemeColorName(d.text), fyne.CurrentApp().Settings().ThemeVariant())
	}
}

func (d *DayWidget) GetFirstColor() color.Color {
	t := fyne.CurrentApp().Settings().Theme()
	myt, ok := t.(mytheme.MyThemeInterface)
	if ok {
		return myt.GetSpecialColor(d.backColor1)
	} else {
		return t.Color(fyne.ThemeColorName(d.text), fyne.CurrentApp().Settings().ThemeVariant())
	}
}

// Widget interface
func (d *DayWidget) CreateRenderer() fyne.WidgetRenderer {
	r := DayWidgetRenderer{
		w: d,
	}
	r.text = canvas.NewText(d.text, d.GetTextColor())
	r.text.Alignment = fyne.TextAlignCenter
	switch d.colorMode {
	case ColorMode_1r:
		r.rect1 = canvas.NewRectangle(d.GetFirstColor())
	case ColorMode_2rtb, ColorMode_2rlr:
		r.rect1 = canvas.NewRectangle(d.GetFirstColor())
		r.rect2 = canvas.NewRectangle(d.GetSecondColor())
	case ColorMode_2tltrb, ColorMode_2trtlb:
		r.tri = canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
			switch r.w.colorMode {
			case ColorMode_2tltrb:
				m := float64(h) / float64(w)
				y1 := float64(h) - float64(x)*m
				if float64(y) < y1 {
					return r.w.GetFirstColor()
				} else {
					return r.w.GetSecondColor()
				}
			case ColorMode_2trtlb:
				m := float64(h) / float64(w)
				y1 := float64(x) * m
				if float64(y) < y1 {
					return r.w.GetFirstColor()
				} else {
					return r.w.GetSecondColor()
				}
			}
			return r.w.GetFirstColor()
		})
	}
	if r.w.protected != nil {
		r.protected = canvas.NewImageFromResource(r.w.protected)
		r.protected.FillMode = canvas.ImageFillContain

	} else {
		r.protected = nil
	}
	if r.w.mark != nil {
		r.mark = canvas.NewImageFromResource(r.w.mark)
		r.mark.FillMode = canvas.ImageFillContain

	} else {
		r.mark = nil
	}
	return &r
}

// Tapped Interface
func (d *DayWidget) Tapped(e *fyne.PointEvent) {
	if d.OnTapped != nil {
		d.OnTapped(e)
	}
}

// WidgetRenderer interface
func (r *DayWidgetRenderer) Layout(size fyne.Size) {
	switch r.w.colorMode {
	case ColorMode_1r:
		r.rect1.Resize(size.AddWidthHeight(.5, .5))
		r.rect1.Move(fyne.NewPos(0, 0))
	case ColorMode_2rlr:
		size2 := fyne.NewSize(size.Width/2, size.Height)
		r.rect1.Resize(size2.AddWidthHeight(.5, .5))
		r.rect2.Resize(size2.AddWidthHeight(.5, .5))
		r.rect1.Move(fyne.NewPos(0, 0))
		r.rect2.Move(fyne.NewPos(size2.Width, 0))
	case ColorMode_2rtb:
		size2 := fyne.NewSize(size.Width, size.Height/2)
		r.rect1.Resize(size2.AddWidthHeight(.5, .5))
		r.rect2.Resize(size2.AddWidthHeight(.5, .5))
		r.rect1.Move(fyne.NewPos(0, 0))
		r.rect2.Move(fyne.NewPos(0, size.Height/2))
	case ColorMode_2tltrb, ColorMode_2trtlb:
		r.tri.Resize(size.AddWidthHeight(.5, .5))
		r.tri.Move(fyne.NewPos(0, 0))
	}
	r.text.Resize(r.text.MinSize())
	r.text.Move(fyne.NewPos(size.Width/2-r.text.MinSize().Width/2, size.Height/2-r.text.MinSize().Height/2-theme.InnerPadding()/42))

	s := r.w.GetOverlayScaleFactor()
	ps := fyne.NewSize(size.Width*s, size.Height*s)

	if r.protected != nil {
		r.protected.Resize(ps)
		r.protected.Move(fyne.NewPos(size.Width-ps.Width-2, size.Height-ps.Height-2))
	}
	if r.mark != nil {
		r.mark.Resize(ps)
		r.mark.Move(fyne.NewPos(2, 2))
	}
}

// WidgetRenderer interface
func (r *DayWidgetRenderer) MinSize() fyne.Size {
	s := util.GetDefaultTextSize("XX")
	pad := theme.InnerPadding() * 2
	return s.AddWidthHeight(pad, pad)
}

// WidgetRenderer interface
func (r *DayWidgetRenderer) Refresh() {
	r.text.Text = r.w.text
	r.text.TextStyle = fyne.TextStyle{
		Bold: r.w.bold,
	}
	r.text.Color = r.w.GetTextColor()
	r.text.Refresh()
	if r.protected != nil {
		r.protected.Refresh()
	}
	if r.mark != nil {
		r.mark.Refresh()
	}
	switch r.w.colorMode {
	case ColorMode_1r:
		r.rect1.FillColor = r.w.GetFirstColor()
		r.rect1.Refresh()
	case ColorMode_2rlr:
		r.rect1.FillColor = r.w.GetFirstColor()
		r.rect2.FillColor = r.w.GetSecondColor()
		r.rect1.Refresh()
		r.rect2.Refresh()
	case ColorMode_2tltrb, ColorMode_2trtlb:
		r.tri.Refresh()
	}
}

// WidgetRenderer interface
func (r *DayWidgetRenderer) Destroy() {
}

// WidgetRenderer interface
func (r *DayWidgetRenderer) Objects() []fyne.CanvasObject {
	list := make([]fyne.CanvasObject, 0, 5)
	switch r.w.colorMode {
	case ColorMode_1r:
		list = append(list, r.rect1, r.text)
	case ColorMode_2rlr, ColorMode_2rtb:
		list = append(list, r.rect1, r.rect2, r.text)
	case ColorMode_2tltrb, ColorMode_2trtlb:
		list = append(list, r.tri, r.text)
	}
	if r.protected != nil {
		list = append(list, r.protected)
	}
	if r.mark != nil {
		list = append(list, r.mark)
	}
	return list
}
