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
	"bytemystery-com/smartdiary/database"
	"bytemystery-com/smartdiary/menubutton"
	"bytemystery-com/smartdiary/selectlayout"
	"bytemystery-com/smartdiary/util"
	"embed"
	"errors"
	"fmt"
	"image/color"
	"net/url"
	"runtime"
	"runtime/debug"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bytemystery-com/colorlabel"
)

func showInfoDialog() {
	vgo := runtime.Version()[2:]
	vfyne := ""
	os := runtime.GOOS
	arch := runtime.GOARCH
	info, _ := debug.ReadBuildInfo()
	for _, dep := range info.Deps {
		if dep.Path == "fyne.io/fyne/v2" {
			vfyne = dep.Version[1:]
		}
	}
	s := fyne.CurrentApp().Settings()
	t := Gui.Theme.GetVariant()
	thema := ""
	b := s.BuildType()
	_ = b
	switch t {
	case theme.VariantDark:
		thema = lang.X("info.thema_dark", "Dark")
	case theme.VariantLight:
		thema = lang.X("info.thema_light", "Light")
	default:
		thema = lang.X("info.thema_unknown", "Unknown")
	}

	build := ""
	switch b {
	case fyne.BuildStandard:
		build = lang.X("info.build_standard", "Standard")
	case fyne.BuildDebug:
		build = lang.X("info.build_debug", "Debug")
	case fyne.BuildRelease:
		build = lang.X("info.build_release", "Release")
	default:
		build = lang.X("info.build_unknown", "Unknown")
	}

	m := Gui.App.Metadata()
	v := fmt.Sprintf("%s (%d)", m.Version, m.Build)
	n := m.Name
	if n == "" {
		n = "Bilderwörterbuch"
	}
	tsStr := ""
	ts := m.Custom["buildts"]
	if ts != "" {
		tsStr = "Build: " + ts + "\n"
	}
	wSize := Gui.MainWindow.Canvas().Size()

	anzahl, _ := Database.GetNumberOfEntries()

	msg := fmt.Sprintf(lang.X("info.msg", "%s\nVersion: %s  \n%sAuthor: Reiner Pröls\n\nGo version: %s\n\nFyne version: %s\nBuild: %s\nThema: %s\nWindow size: %.0fx%.0f\n\nPlatform: %s\nArchitecture: %s\n\nDatabase: %s\nEntries: %d"),
		n, v, tsStr, vgo, vfyne, build, thema, wSize.Width, wSize.Height, os, arch, Gui.DatabaseFile, anzahl)
	dialog.ShowInformation(lang.X("info.title", "Info"), msg, Gui.MainWindow)
}

func loadPreferences() {
	Gui.Settings = NewPreferences()
}

func loadIcon(path, name string) *fyne.StaticResource {
	data, err := assets.ReadFile(path)
	if err != nil {
		return nil
	}
	return fyne.NewStaticResource(name, data)
}

func loadTranslations(fs embed.FS, dir string) {
	lang.AddTranslationsFS(fs, dir)
}

func loadIcons() {
	Gui.Icon = loadIcon("assets/icons/icon.png", "icon")
	Gui.App.SetIcon(Gui.Icon)

	Gui.IconBack1y = loadIcon("assets/icons/back1y.svg", "back1y")
	Gui.IconBack1m = loadIcon("assets/icons/back1m.svg", "back1m")
	Gui.IconForeward1y = loadIcon("assets/icons/foreward1y.svg", "foreward1y")
	Gui.IconForeward1m = loadIcon("assets/icons/foreward1m.svg", "foreward1m")

	Gui.IconSearchUp = loadIcon("assets/icons/search_u.png", "search_u")
	Gui.IconSearchDown = loadIcon("assets/icons/search_d.png", "search_d")

	Gui.IconBackUp = loadIcon("assets/icons/back_u.png", "back_u")
	Gui.IconBackDown = loadIcon("assets/icons/back_d.png", "back_d")

	Gui.IconExport = loadIcon("assets/icons/export.svg", "export")
	Gui.IconImport = loadIcon("assets/icons/import.svg", "import")

	Gui.IconProtected = loadIcon("assets/icons/protected.svg", "protected_svg")

	Gui.IconMarks = Gui.IconMarks[:0]
	// https://svgomg.net/
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/empty.svg", "empty_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/mark.svg", "mark_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/triangle.svg", "triangle_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/circle.svg", "circle_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/box.svg", "box_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/diamond.svg", "diamond_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/star.svg", "star_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/flash.svg", "flash_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/priority.svg", "priority_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/question.svg", "question_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/check.svg", "check_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/plus.svg", "plus_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/minus.svg", "minus_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/wave.svg", "wave_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/heart.svg", "heart_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/anchor.svg", "anchor_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/lamp.svg", "lamp_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/all.svg", "all_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/smily.svg", "smily_svg"))
	Gui.IconMarks = append(Gui.IconMarks, loadIcon("assets/icons/music.svg", "music_svg"))

	loadIconsForTheme()
}

func loadIconsForTheme() {
	dir := ""
	switch Gui.Theme.GetVariant() {
	case theme.VariantDark:
		dir = "dark"
	case theme.VariantLight:
		dir = "light"
	default:
		dir = "light"
	}
	Gui.IconColorModeSimple = loadIcon("assets/icons/"+dir+"/r.png", "colormode_rect")
	Gui.IconColorModeTopBottom = loadIcon("assets/icons/"+dir+"/tb.png", "colormode_tb")
	Gui.IconColorModeLeftRight = loadIcon("assets/icons/"+dir+"/lr.png", "colormode_lr")
	Gui.IconColorModeLeftTopRightBottom = loadIcon("assets/icons/"+dir+"/ltrb.png", "colormode_ltrb")
	Gui.IconColorModeRightTopLeftBottom = loadIcon("assets/icons/"+dir+"/rtlb.png", "colormode_rtlb")

	/*	Gui.IconImport = loadIcon("assets/icons/"+dir+"/import.png", "import")
		Gui.IconExport = loadIcon("assets/icons/"+dir+"/export.png", "export")
	*/
	_ = dir
}

func SendNotification(title, msg string) {
	fyne.Do(func() {
		n := fyne.NewNotification(title, msg)
		Gui.App.SendNotification(n)
	})
}

func doHelp() {
	u := url.URL{
		Scheme: "https",
		Host:   "github.com",
		Path:   "/bytemystery-com/smartdiary",
	}
	Gui.App.OpenURL(&u)
}

func UIErrorHandler(err error) {
	UIErrorHandlerWithMessage(err, "")
}

func UIErrorHandlerWithMessage(err error, msg string) {
	fyne.Do(func() {
		if msg != "" {
			if msg[len(msg)-1] != '\n' {
				msg += "\n"
			}
			err = errors.Join(errors.New(msg), err)
		}
		dialog.ShowError(err, Gui.MainWindow)
	})
}

func ShowExportTypeDialog(f func(overwrite bool)) {
	if Gui.IsDesktop {
		f(false)
	} else {
		switch Gui.Settings.ExportMode {
		case EXPORTMODE_NEW:
			f(false)
		case EXPORTMODE_OVERWRITE:
			f(true)
		default:
			var dia *dialog.CustomDialog
			cancel := widget.NewButton(lang.X("cancel", "Cancel"), func() {
				dia.Hide()
			})
			overwrite := widget.NewButton(lang.X("export.overwrite", "Overwrite"), func() {
				dia.Hide()
				f(true)
			})
			overwrite.Importance = widget.HighImportance
			newFile := widget.NewButton(lang.X("export.new", "New"), func() {
				dia.Hide()
				f(false)
			})
			newFile.Importance = widget.HighImportance
			c := container.NewHBox(layout.NewSpacer(), cancel, newFile, overwrite, layout.NewSpacer())
			dia = dialog.NewCustomWithoutButtons(lang.X("export.mode.title", "New file or overwrite"), c, Gui.MainWindow)
			dia.Show()
			si := Gui.MainWindow.Canvas().Size()
			var windowScale float32 = 1.0
			dia.Resize(fyne.NewSize(si.Width*windowScale, dia.MinSize().Height))
		}
	}
}

func ShowPopUp(ev *fyne.PointEvent, c fyne.CanvasObject) *widget.PopUp {
	scroll := container.NewScroll(c)
	contentSize := c.MinSize().AddWidthHeight(theme.ScrollBarSize(), theme.ScrollBarSize())
	windowSize := Gui.MainWindow.Canvas().Size()
	max := fyne.NewSize(windowSize.Width*0.8, windowSize.Height*0.8)
	size := contentSize
	if size.Width > max.Width {
		size.Width = max.Width
	}
	if size.Height > max.Height {
		size.Height = max.Height
	}
	scroll.SetMinSize(size)

	//	scroll.SetMinSize(c.MinSize())
	popup := widget.NewPopUp(scroll, Gui.MainWindow.Canvas())
	popup.Resize(size)

	pos := ev.AbsolutePosition
	if pos.X > c.MinSize().Width/2+4*theme.Padding() {
		pos.X -= c.MinSize().Width / 2
	}
	si := Gui.MainWindow.Canvas().Size()
	if pos.X+c.MinSize().Width+4*theme.Padding() > si.Width {
		pos.X = si.Width - c.MinSize().Width - 4*theme.Padding()
	}
	popup.Move(pos)
	popup.Show()
	return popup
}

func CreateCategorySelector(hasNull bool, onSelect func(int64), isSelected func(int64) bool, showAll bool, hasOrderBtn bool) (*colorlabel.ColorLabel, fyne.CanvasObject) {
	var catData database.CategoryDataType
	var categories []*database.CategoryDataType

	updateCat := func() {
		if hasNull {
			catData = database.CategoryDataType{
				Name:  lang.X("category_not_set", "not defined"),
				Color: string(theme.ColorNameBackground),
				Id:    -1,
			}
			categories = make([]*database.CategoryDataType, 0, len(CategoryData)+1)
			categories = append(categories, &catData)
			categories = append(categories, CategoryData...)
		} else {
			catData = *CategoryData[0]
			categories = CategoryData
		}
	}
	updateCat()

	category := colorlabel.NewColorLabel(catData.Name, Gui.Theme.GetSpecialColor("category_text_color"), Gui.Theme.GetSpecialColor(catData.Color), Gui.Theme.GetSpecialSize("categorie_select"))
	category.SetTextStyle(&fyne.TextStyle{
		Bold: true,
	})
	category.SetAlinment(fyne.TextAlignCenter)

	doPopup := func(*fyne.PointEvent) {}

	doPopup = func(ev *fyne.PointEvent) {
		updateCat()
		c := container.NewVBox()
		var popup *widget.PopUp
		for index, item := range categories {
			if !showAll && item.Name == "" {
				continue
			}
			ctrl := colorlabel.NewColorLabel(item.Name, Gui.Theme.GetSpecialColor("category_text_color"), Gui.Theme.GetSpecialColor(item.Color), Gui.Theme.GetSpecialSize("categorie_select_popup"))
			ctrl.SetTextStyle(&fyne.TextStyle{
				Bold: true,
			})
			if hasNull && item.Id == -1 {
				ctrl.SetTextColor(Gui.Theme.GetSpecialColor("nocategory_text_color"))
			}
			ctrl.OnTapped = func() {
				popup.Hide()
				category.SetText(item.Name)
				category.SetTextColor(Gui.Theme.GetSpecialColor("category_text_color"))
				category.SetBackgroundColor(Gui.Theme.GetSpecialColor(item.Color))
				onSelect(item.Id)
			}
			icon := widget.NewIcon(theme.RadioButtonIcon())
			var btn *menubutton.MenuButton
			if hasOrderBtn {
				btn = menubutton.NewMenuButton("", theme.MoreHorizontalIcon(), func(ev1 *fyne.PointEvent) {
					popupMenu := fyne.NewMenu("")
					if index > 0 {
						popupMenu.Items = append(popupMenu.Items, fyne.NewMenuItemWithIcon(lang.X("menu.up", "Move up"), theme.MoveUpIcon(), func() {
							CategoryData[index-1], CategoryData[index] = CategoryData[index], CategoryData[index-1]
							CategoryData[index-1].OrderIndex, CategoryData[index].OrderIndex = CategoryData[index].OrderIndex, CategoryData[index-1].OrderIndex
							CategoryData[index-1].IsEdited = true
							CategoryData[index].IsEdited = true
							popup.Hide()
							doPopup(ev)
						}))
					}
					if index < len(categories)-1 {
						popupMenu.Items = append(popupMenu.Items, fyne.NewMenuItemWithIcon(lang.X("menu.down", "Move down"), theme.MoveDownIcon(), func() {
							CategoryData[index], CategoryData[index+1] = CategoryData[index+1], CategoryData[index]
							CategoryData[index].OrderIndex, CategoryData[index+1].OrderIndex = CategoryData[index+1].OrderIndex, CategoryData[index].OrderIndex
							CategoryData[index+1].IsEdited = true
							CategoryData[index].IsEdited = true
							popup.Hide()
							doPopup(ev)
						}))
					}
					widget.ShowPopUpMenuAtPosition(popupMenu, Gui.App.Driver().CanvasForObject(btn), ev1.AbsolutePosition)
				})
			}
			if isSelected(item.Id) {
				icon.SetResource(theme.RadioButtonFillIcon())
			}
			if btn != nil {
				c.Add(container.NewBorder(nil, nil, icon, btn, ctrl))
			} else {
				c.Add(container.NewBorder(nil, nil, icon, nil, ctrl))
			}
		}
		popup = ShowPopUp(ev, c)
	}
	btn := menubutton.NewMenuButton("", theme.MoveDownIcon(), func(ev *fyne.PointEvent) {
		doPopup(ev)
	})
	category.OnTappedEx = func(ev *fyne.PointEvent) {
		doPopup(ev)
	}

	var maxWidth float32 = 0
	list := make([]string, 0, len(categories)+1)
	if hasNull {
		list = append(list, lang.X("category_not_set", "not defined"))
	}
	for _, item := range categories {
		list = append(list, item.Name)
	}
	for _, item := range list {
		w := util.GetDefaultTextWidth(item + "XX")
		if w > maxWidth {
			maxWidth = w
		}
	}
	maxWidth = fyne.Max(maxWidth, util.GetDefaultTextWidth("XXXXX"))
	maxWidth += theme.Padding() * 2

	h := fyne.Max(category.MinSize().Height, btn.MinSize().Height)
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeWidth = theme.Padding()
	border.StrokeColor = theme.Color(theme.ColorNameShadow)
	border.FillColor = color.Transparent
	ctrls := container.New(&selectlayout.SelectLayout{}, container.NewGridWrap(fyne.NewSize(maxWidth, h), category), btn, border)
	return category, ctrls
}

/*
func ShowPasswordDialog(bChange bool, doOk func(ok bool)) {
	var dia *dialog.ConfirmDialog
	pass := widget.NewPasswordEntry()
	if !bChange {
		pass.OnSubmitted = func(string) { dia.Confirm() }
	}
	pass.Validator = func(t string) error {
		return Gui.Settings.CheckPasswd([]byte(t))
	}
	pass1 := widget.NewPasswordEntry()
	pass1.Validator = func(t string) error {
		if len(t) >= 4 {
			return nil
		} else {
			return errors.New(lang.X("password.tooshort", "Password is too short"))
		}
	}

	pass2 := widget.NewPasswordEntry()
	pass2.Validator = func(t string) error {
		if t == pass1.Text {
			return nil
		} else {
			return errors.New(lang.X("password.confirm_failed", "The new password and the confirm password do not match"))
		}
	}
	pass2.AlwaysShowValidationError = true

	if !bChange {
		pass2.OnSubmitted = func(string) { dia.Confirm() }
	}
	cNew := container.NewVBox(widget.NewLabel(lang.X("password.new_password", "New password")), pass1,
		widget.NewLabel(lang.X("password.confirm_new_password", "Confirm new password")), pass2)
	c := container.NewVBox(widget.NewLabel(lang.X("password.password", "Password")), pass)
	if bChange {
		c.Add(cNew)
	}

	dia = dialog.NewCustomConfirm(lang.X("password.title", "Password"), lang.X("ok", "Ok"), lang.X("cancel", "Cancel"),
		c, func(ok bool) {
			if !ok {
				if doOk != nil {
					doOk(false)
				}
				return
			}
			err := Gui.Settings.CheckPasswd([]byte(pass.Text))
			if err != nil {
				UIErrorHandler(err)
				dia.Show()
				return
			}
			if bChange {
				err := pass1.Validate()
				if err != nil {
					UIErrorHandler(err)
					dia.Show()
					return
				}
				err = pass2.Validate()
				if err != nil {
					UIErrorHandler(err)
					dia.Show()
					return
				}
				Gui.Settings.SavePasswd([]byte(pass1.Text))
			} else {
				Gui.lastPasswordCheck = time.Now()
			}
			if doOk != nil {
				doOk(true)
			}
		}, Gui.MainWindow)
	dia.Show()
	si := Gui.MainWindow.Canvas().Size()
	dia.Resize(fyne.NewSize(si.Width, dia.MinSize().Height))
	Gui.MainWindow.Canvas().Focus(pass)
}
*/
