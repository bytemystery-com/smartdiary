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
	"errors"
	"time"

	"bytemystery-com/smartdiary/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type PaswordView struct {
	content fyne.CanvasObject
	fOk     func(bool)
	pass    *widget.Entry
	pass1   *widget.Entry
	pass2   *widget.Entry
	bChange bool
	cNew    *fyne.Container
	icon    *canvas.Image
}

func NewPasswordView() *PaswordView {
	p := PaswordView{}
	p.pass = widget.NewPasswordEntry()
	p.pass.OnSubmitted = func(string) { p.Confirm() }
	p.pass.Validator = func(t string) error {
		return Gui.Settings.CheckPasswd([]byte(t))
	}
	p.pass1 = widget.NewPasswordEntry()
	p.pass1.Validator = func(t string) error {
		if len(t) >= 4 {
			return nil
		} else {
			return errors.New(lang.X("password.tooshort", "Password is too short"))
		}
	}

	p.pass2 = widget.NewPasswordEntry()
	p.pass2.Validator = func(t string) error {
		if t == p.pass1.Text {
			return nil
		} else {
			return errors.New(lang.X("password.confirm_failed", "The new password and the confirm password do not match"))
		}
	}
	p.pass2.AlwaysShowValidationError = true
	p.pass2.OnSubmitted = func(string) { p.Confirm() }
	p.cNew = container.NewVBox(widget.NewLabel(lang.X("password.new_password", "New password")), p.pass1,
		widget.NewLabel(lang.X("password.confirm_new_password", "Confirm new password")), p.pass2)
	p.icon = canvas.NewImageFromResource(Gui.Icon)
	p.icon.FillMode = canvas.ImageFillContain
	si := Gui.Theme.GetSpecialSize("login_icon_size")

	p.icon.SetMinSize(fyne.NewSize(si, si))
	p.icon.Refresh()
	var c *fyne.Container
	if Gui.IsDesktop {
		s1 := Gui.Theme.GetSpecialSize("login_space_logo_label")
		c = container.NewVBox(container.NewCenter(p.icon), util.NewVFiller(s1), widget.NewLabel(lang.X("password.password", "Password")), p.pass)
	} else {
		c = container.NewVBox(container.NewCenter(p.icon), widget.NewLabel(lang.X("password.password", "Password")), p.pass)
	}
	c.Add(p.cNew)

	ok := widget.NewButton(lang.X("ok", "Ok"), func() { p.Confirm() })
	ok.Importance = widget.HighImportance
	cancel := widget.NewButton(lang.X("cancel", "Cancel"), func() { p.Cancel() })

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
		p.content = container.NewBorder(nil, buttonLine, nil, nil, c)
	} else {
		c.Add(util.NewVFiller(Gui.Theme.GetSpecialSize("login_space_field_ok")))
		c.Add(buttonLine)
		p.content = c
	}

	return &p
}

func (p *PaswordView) SetMode(bChange bool, fOk func(ok bool)) {
	if bChange {
		p.pass.OnSubmitted = nil
		p.pass2.OnSubmitted = func(string) { p.Confirm() }
		p.cNew.Show()
		if Gui.IsDesktop {
			p.icon.Show()
		} else {
			p.icon.Hide()
		}
	} else {
		p.pass.OnSubmitted = func(string) { p.Confirm() }
		p.pass2.OnSubmitted = nil
		p.cNew.Hide()
		p.icon.Show()
	}
	p.pass.SetText("")
	p.pass1.SetText("")
	p.pass2.SetText("")
	p.fOk = fOk
}

func (p *PaswordView) Confirm() {
	err := Gui.Settings.CheckPasswd([]byte(p.pass.Text))
	if err != nil {
		UIErrorHandler(err)
		Gui.MainWindow.Canvas().Focus(p.pass)
		return
	}
	if p.bChange {
		err := p.pass1.Validate()
		if err != nil {
			UIErrorHandler(err)
			Gui.MainWindow.Canvas().Focus(p.pass1)
			return
		}
		err = p.pass2.Validate()
		if err != nil {
			UIErrorHandler(err)
			Gui.MainWindow.Canvas().Focus(p.pass2)
			return
		}
		Gui.Settings.SavePasswd([]byte(p.pass1.Text))
	} else {
		Gui.lastPasswordCheck = time.Now()
	}
	p.pass.SetText("")
	p.pass1.SetText("")
	p.pass2.SetText("")
	if p.fOk != nil {
		p.fOk(true)
	}
}

func (p *PaswordView) Cancel() {
	p.pass.SetText("")
	p.pass1.SetText("")
	p.pass2.SetText("")
	if p.fOk != nil {
		p.fOk(false)
	}
}

func (p *PaswordView) GetContent() fyne.CanvasObject {
	return p.content
}

func (p *PaswordView) UpdateToolBar() {
	Gui.Toolbar.Items = []widget.ToolbarItem{Gui.toolToggleThema, widget.NewToolbarSpacer(), Gui.toolInfo}
	Gui.Toolbar.Refresh()
}

func (p *PaswordView) ThemeChanged() {
	p.content.Refresh()
}
