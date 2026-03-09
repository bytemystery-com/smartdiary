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

package selectlayout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type SelectLayout struct{}

func (l *SelectLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	n := len(objects)
	if n != 3 {
		return
	}
	pad := theme.Padding()

	h := fyne.Max(objects[0].MinSize().Height, objects[1].MinSize().Height)
	objects[0].Resize(fyne.NewSize(objects[0].MinSize().Width, h))
	objects[1].Resize(fyne.NewSize(objects[1].MinSize().Width, h))
	objects[2].Resize(fyne.NewSize(objects[0].MinSize().Width+objects[1].MinSize().Width, h))

	objects[0].Move(fyne.NewPos(pad, pad))
	objects[1].Move(fyne.NewPos(objects[0].MinSize().Width+pad, pad))
	objects[2].Move(fyne.NewPos(pad, pad))
}

func (l *SelectLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	n := len(objects)
	if n != 3 {
		return fyne.NewSize(theme.InnerPadding(), theme.InnerPadding())
	}
	pad := theme.Padding()
	h := fyne.Max(objects[0].MinSize().Height, objects[1].MinSize().Height) + 2*pad
	return fyne.NewSize(objects[0].MinSize().Width+objects[1].MinSize().Width+2*pad, h)
}
