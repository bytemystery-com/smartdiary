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

package calendarlayout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type CalendarLayout struct{}

func (l *CalendarLayout) getRowsAndCols(n int) (int, int) {
	if n < 7 {
		return 1, n
	}
	rows := n / 7
	if rows*7 != n {
		n++
	}
	return rows, 7
}

func (l *CalendarLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects) - 5
	if n <= 0 {
		return
	}
	rows, cols := l.getRowsAndCols(n)
	index := 0
	w := size.Width / float32(cols)
	h := size.Height / float32(rows)
	s := fyne.Min(w, h)
	si := fyne.NewSize(s, s)

	btnW := fyne.Max(objects[0].MinSize().Width, objects[1].MinSize().Width) * 1.2
	btnH := objects[0].MinSize().Height
	pad := theme.Padding()
	btnSize := fyne.NewSize(btnW, btnH)
	objects[0].Resize(btnSize)
	objects[1].Resize(btnSize)
	objects[2].Resize(btnSize)
	objects[3].Resize(btnSize)
	objects[4].Resize(fyne.NewSize(s*7-4*btnW-2*pad, objects[4].MinSize().Height))
	objects[0].Move(fyne.NewPos(0, 0))
	objects[1].Move(fyne.NewPos(btnW+pad, 0))
	objects[2].Move(fyne.NewPos(s*7-2*btnW-pad, 0))
	objects[3].Move(fyne.NewPos(s*7-btnW, 0))
	objects[4].Move(fyne.NewPos(2*btnW+pad, 0))

	index = 5
	for row := range rows {
		for col := range cols {
			objects[index].Resize(si)
			objects[index].Move(fyne.NewPos(float32(col)*s, btnH+float32(row)*s))
			index++
		}
	}
}

func (l *CalendarLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	n := len(objects)
	if n < 6 {
		return fyne.NewSize(theme.InnerPadding(), theme.InnerPadding())
	}
	rows, cols := l.getRowsAndCols(n)
	s := objects[5].MinSize()
	return fyne.NewSize(float32(cols)*s.Width, float32(rows)*s.Height+objects[0].MinSize().Height)
}
