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

package tapcontainer

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

var (
	_ fyne.Tappable          = (*TapContainer)(nil)
	_ fyne.SecondaryTappable = (*TapContainer)(nil)
)

type TapContainer struct {
	widget.BaseWidget
	content           fyne.CanvasObject
	OnTappedSecondary func(*fyne.PointEvent)
	OnTapped          func(*fyne.PointEvent)
}

func NewTapContainer(c fyne.CanvasObject) *TapContainer {
	l := &TapContainer{content: c}
	l.ExtendBaseWidget(l)
	return l
}

func (l *TapContainer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(l.content)
}

// Event Handler
func (l *TapContainer) Tapped(ev *fyne.PointEvent) {
	if l.OnTapped != nil {
		l.OnTapped(ev)
	}
}

func (l *TapContainer) TappedSecondary(ev *fyne.PointEvent) {
	if l.OnTappedSecondary != nil {
		l.OnTappedSecondary(ev)
	}
}
