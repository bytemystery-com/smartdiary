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

package searchlayout

import (
	"myapp/daywidget"
	"myapp/mytheme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type SearchLayout struct{}

func (l *SearchLayout) GetOverlayScaleFactor() float32 {
	t := fyne.CurrentApp().Settings().Theme()
	myt, ok := t.(mytheme.MyThemeInterface)
	if ok {
		return myt.GetSpecialSize("overlay_size_scale")
	} else {
		return daywidget.OVERLAY_SIZE
	}
}

func (l *SearchLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects)
	if n != 6 {
		return
	}
	icon := objects[0]
	date := objects[1]
	btn := objects[2]
	rich := objects[3]
	mark := objects[4]
	protected := objects[5]
	pad := theme.Padding()

	icon.Resize(icon.MinSize())
	btn.Resize(btn.MinSize())
	date.Resize(fyne.NewSize(size.Width-pad-icon.MinSize().Width-btn.MinSize().Width, date.MinSize().Height))
	rich.Resize(fyne.NewSize(size.Width, rich.MinSize().Height))

	s := l.GetOverlayScaleFactor()

	s4 := fyne.NewSize(icon.MinSize().Width*s, icon.MinSize().Height*s)
	if mark.Visible() {
		mark.Resize(s4)
	}
	if protected.Visible() {
		protected.Resize(s4)
	}
	icon.Move(fyne.NewPos(0, 0))
	y := (icon.MinSize().Height - date.MinSize().Height) / 2
	date.Move(fyne.NewPos(icon.MinSize().Width+pad, y))
	btn.Move(fyne.NewPos(size.Width-theme.ScrollBarSize()-btn.MinSize().Width, icon.MinSize().Height-btn.MinSize().Height))
	rich.Move(fyne.NewPos(-2*pad, icon.MinSize().Height-2*pad))

	if mark.Visible() {
		mark.Move(fyne.NewPos(2, 2))
	}
	if protected.Visible() {
		protected.Move(fyne.NewPos(icon.MinSize().Width-s4.Width-2, icon.MinSize().Height-s4.Height-2))
	}
}

func (l *SearchLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	n := len(objects)
	pad := theme.Padding()
	if n != 6 {
		return fyne.NewSize(pad, pad)
	}
	icon := objects[0]
	date := objects[1]
	btn := objects[2]
	rich := objects[3]
	return fyne.NewSize(icon.MinSize().Width+pad/2+date.MinSize().Width+btn.MinSize().Width+pad, icon.MinSize().Height+rich.MinSize().Height)
}
