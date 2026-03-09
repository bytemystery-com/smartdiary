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
	"embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"myapp/database"
	"myapp/mytheme"
	"myapp/util"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type GUI struct {
	App          fyne.App
	MainWindow   fyne.Window
	Toolbar      *widget.Toolbar
	IsDesktop    bool
	Icon         *fyne.StaticResource
	FyneSettings fyne.Settings
	Settings     *Preferences
	Theme        *mytheme.MyTheme

	toolToggleThema *widget.ToolbarAction
	toolInfo        *widget.ToolbarAction
	toolSettings    *widget.ToolbarAction
	toolRemove      *widget.ToolbarAction
	toolExport      *widget.ToolbarAction
	toolImport      *widget.ToolbarAction

	IconBack1y     *fyne.StaticResource
	IconBack1m     *fyne.StaticResource
	IconForeward1y *fyne.StaticResource
	IconForeward1m *fyne.StaticResource

	IconColorModeSimple             *fyne.StaticResource
	IconColorModeTopBottom          *fyne.StaticResource
	IconColorModeLeftRight          *fyne.StaticResource
	IconColorModeLeftTopRightBottom *fyne.StaticResource
	IconColorModeRightTopLeftBottom *fyne.StaticResource

	IconSearchUp   *fyne.StaticResource
	IconSearchDown *fyne.StaticResource

	IconBackUp   *fyne.StaticResource
	IconBackDown *fyne.StaticResource

	IconProtected *fyne.StaticResource
	IconMarks     []*fyne.StaticResource

	IconExport *fyne.StaticResource
	IconImport *fyne.StaticResource

	content *fyne.Container

	DatabaseFile string

	mainView     *MainView
	entryView    *EntryView
	searchView   *SearchView
	settingsView *SettingsView
	passwordView *PaswordView

	beforeSettingsContent fyne.CanvasObject
	lastPasswordCheck     time.Time
}

//go:embed assets/*
var assets embed.FS

func init() {
	CategoryMap = make(map[int64]*database.CategoryDataType, 12)
	Gui.IconMarks = make([]*fyne.StaticResource, 0, 10)
}

var (
	Database     *database.Db
	CategoryData []*database.CategoryDataType
	CategoryMap  map[int64]*database.CategoryDataType
	Gui          = GUI{}
)

var Months = []string{}

var Weekdays = []string{}

func forceLanguage() {
	if *Flags.language == "" {
		return
	}
	// Hack. Ongoing discussion in https://github.com/fyne-io/fyne/issues/5333
	lcontent, err := assets.ReadFile("assets/lang/" + *Flags.language + ".json")
	if err != nil {
		return
	}
	lang.AddTranslationsForLocale(lcontent, lang.SystemLocale())
}

type FlagsType struct {
	language *string
}

var Flags FlagsType

func main() {
	Flags.language = flag.String("l", "", "language (en, de ....)")
	flag.Parse()

	loadTranslations(assets, "assets/lang")
	forceLanguage()
	Months = []string{
		lang.X("month_0", "January"), lang.X("month_1", "February"), lang.X("month_2", "March"),
		lang.X("month_3", "April"), lang.X("month_4", "May"), lang.X("month_5", "June"),
		lang.X("month_6", "July"), lang.X("month_7", "August"), lang.X("month_8", "September"),
		lang.X("month_9", "October"), lang.X("month_10", "November"), lang.X("month_11", "December"),
	}

	Weekdays = []string{
		lang.X("weekday_0", "Su"), lang.X("weekday_1", "Mo"), lang.X("weekday_2", "Tu"),
		lang.X("weekday_3", "We"), lang.X("weekday_4", "Th"), lang.X("weekday_5", "Fr"), lang.X("weekday_6", "Sa"),
	}

	//  go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	Gui.App = app.NewWithID("com.bytemystery.smartdiary2")
	loadPreferences()
	Gui.FyneSettings = Gui.App.Settings()
	var tv fyne.ThemeVariant
	switch Gui.Settings.ThemeVariant {
	case -1:
		tv = fyne.CurrentApp().Settings().ThemeVariant()
	case 0:
		tv = theme.VariantDark
	case 1:
		tv = theme.VariantLight
	}

	// tv = theme.VariantLight

	Gui.Theme = mytheme.NewMyTheme(tv)
	Gui.App.Settings().SetTheme(Gui.Theme)
	loadIcons()

	if _, ok := Gui.App.(desktop.App); ok {
		Gui.IsDesktop = true
	}
	Gui.MainWindow = Gui.App.NewWindow("MyApp")
	Gui.MainWindow.SetIcon(Gui.Icon)

	Gui.toolToggleThema = widget.NewToolbarAction(theme.BrokenImageIcon(), func() {
		if Gui.Theme.GetVariant() == theme.VariantDark {
			Gui.Theme.SetVariant(theme.VariantLight)
		} else {
			Gui.Theme.SetVariant(theme.VariantDark)
		}
		Gui.Settings.ThemeVariant = int(Gui.Theme.GetVariant())
		Gui.Settings.Store()
		Gui.App.Settings().SetTheme(Gui.Theme)
		updateTheme()
	})

	Gui.toolInfo = widget.NewToolbarAction(theme.InfoIcon(), func() {
		showInfoDialog()
	})
	Gui.toolSettings = widget.NewToolbarAction(theme.SettingsIcon(), func() { ShowSettingsView() })
	Gui.toolRemove = widget.NewToolbarAction(theme.DeleteIcon(), func() { RemoveEntry() })

	Gui.toolExport = widget.NewToolbarAction(theme.NewThemedResource(Gui.IconExport), func() { Export() })
	Gui.toolImport = widget.NewToolbarAction(theme.NewThemedResource(Gui.IconImport), func() { Import() })

	Gui.Toolbar = widget.NewToolbar(Gui.toolToggleThema, widget.NewToolbarSeparator(), widget.NewToolbarSpacer(), Gui.toolInfo)

	scaling := theme.Size("text") / 14.0

	Gui.content = container.NewStack(widget.NewLabel(""))
	Gui.MainWindow.SetContent(container.NewBorder(Gui.Toolbar, nil, nil, nil, Gui.content))

	Gui.MainWindow.Resize(fyne.NewSize(380*scaling, 700*scaling))
	Gui.MainWindow.CenterOnScreen()

	Database = database.NewDb()

	str, err := database.GetDBFile("smartdiary")
	if err != nil {
		slog.Error("No database path")
		return
	}
	Gui.DatabaseFile = str
	// os.Remove(Gui.DatabaseFile)
	err = Database.Open(Gui.DatabaseFile)
	if err != nil {
		slog.Error("open database error", "err", err)
		return
	}
	defer Database.Close()
	UpdateCategoryData()
	/*
		err = Database.ImportFromCSV("/home/reiner/Dropbox/private/export_entries.csv")
		if err != nil {
			panic("")
		}
	*/
	Gui.mainView = NewMainView()
	Gui.entryView = NewEntryView()
	Gui.searchView = NewSearchView()
	Gui.settingsView = NewSettingsView()
	Gui.passwordView = NewPasswordView()

	fyne.CurrentApp().Settings().AddListener(func(settings fyne.Settings) {
		updateTheme()
	})

	Gui.App.Lifecycle().SetOnExitedForeground(func() {
		// if fyne.CurrentDevice().IsMobile() {
	})
	Gui.App.Lifecycle().SetOnEnteredForeground(func() {
	})

	t := time.Now()
	ShowMainView(&t)

	Gui.MainWindow.ShowAndRun()
}

func updateTheme() {
	loadIconsForTheme()
	Gui.entryView.ThemeChanged()
	Gui.mainView.ThemeChanged()
	Gui.searchView.ThemeChanged()
	Gui.settingsView.ThemeChanged()
	Gui.passwordView.ThemeChanged()
}

func UpdateCategoryData() error {
	data, err := Database.GetCategories()
	if err != nil {
		slog.Error("database error", "err", err)
		UIErrorHandler(err)
		return err
	}
	CategoryData = data
	clear(CategoryMap)
	for _, item := range CategoryData {
		CategoryMap[item.Id] = item
	}
	return nil
}

func setContent(c fyne.CanvasObject) {
	Gui.content.RemoveAll()
	Gui.content.Add(c)
	/*
		Gui.content.Add(Gui.busy)
		Gui.busy.Hide()
	*/
	if c == Gui.mainView.GetContent() {
		Gui.mainView.UpdateToolBar()
	}
	if c == Gui.entryView.GetContent() {
		Gui.entryView.UpdateToolBar()
	}
	if c == Gui.settingsView.GetContent() {
		Gui.settingsView.UpdateToolBar()
	}
	if c == Gui.searchView.GetContent() {
		Gui.searchView.UpdateToolBar()
	}
	if c == Gui.passwordView.GetContent() {
		Gui.passwordView.UpdateToolBar()
	}
	Gui.content.Refresh()
}

func ShowMainView(date *time.Time) {
	Gui.mainView.SetDate(date)
	setContent(Gui.mainView.GetContent())
	if Gui.Settings.AutoUpdateCheck {
		CheckForUpdate(true)
	}
}

func ShowEntryView(date *time.Time, fromSearch bool) {
	Gui.entryView.SetDate(date, fromSearch)
	setContent(Gui.entryView.GetContent())
	Gui.MainWindow.Canvas().Focus(Gui.entryView.entry)
	Gui.entryView.SetCursorAtEnd()
}

func ShowSearch(search string) {
	if Gui.searchView.Search(search) {
		SetSearchView()
	}
}

func SetSearchView() {
	setContent(Gui.searchView.GetContent())
}

func ShowPasswordDialog(bChange bool, fOk func(bool)) {
	oldView := Gui.content.Objects[0]
	Gui.passwordView.SetMode(bChange, func(ok bool) {
		setContent(oldView)
		if fOk != nil {
			fOk(ok)
		}
	})
	setContent(Gui.passwordView.GetContent())
	Gui.MainWindow.Canvas().Focus(Gui.passwordView.pass)
}

func RestoreBeforeSettings() {
	if Gui.beforeSettingsContent != nil {
		if Gui.beforeSettingsContent == Gui.mainView.GetContent() {
			Gui.mainView.SetDate(nil)
		}
		setContent(Gui.beforeSettingsContent)
		Gui.beforeSettingsContent = nil
		return
	}
	ShowMainView(nil)
}

func RemoveEntry() {
	if Gui.content.Objects[0] == Gui.entryView.GetContent() {
		Gui.entryView.RemoveEntry()
	}
}

func ShowSettingsView() {
	Gui.beforeSettingsContent = Gui.content.Objects[0]
	Gui.settingsView.Init()
	setContent(Gui.settingsView.GetContent())
}

func doPasswordCheck(doOk func(ok bool)) {
	if time.Since(Gui.lastPasswordCheck) > time.Duration(Gui.Settings.PasswordExpire)*time.Minute {
		ShowPasswordDialog(false, doOk)
	} else if doOk != nil {
		doOk(true)
	}
}

func Reload() {
	Gui.entryView = NewEntryView()
	Gui.settingsView = NewSettingsView()
}

func CheckForUpdate(notify bool) {
	if notify {
		now := time.Now().Unix()
		if now-Gui.Settings.LastUpdatecheck < int64(Gui.Settings.UpdateCheckInterval)*3600 {
			return
		}
	}
	go func() {
		m := Gui.App.Metadata()
		type Version struct {
			maj   int
			min   int
			patch int
		}
		thisVersion := Version{}
		gitVersion := Version{}
		web, newVer, err := util.CheckForUpdate()
		if err != nil {
			return
		}
		n, err := fmt.Sscanf(m.Version, "%d.%d.%d", &thisVersion.maj, &thisVersion.min, &thisVersion.patch)
		if n != 3 || err != nil {
			return
		}
		n, err = fmt.Sscanf(newVer, "v%d.%d.%d", &gitVersion.maj, &gitVersion.min, &gitVersion.patch)
		if n != 3 || err != nil {
			return
		}
		if thisVersion.maj < gitVersion.maj || (thisVersion.maj == gitVersion.maj && thisVersion.min < gitVersion.min) ||
			(thisVersion.maj == gitVersion.maj && thisVersion.min == gitVersion.min && thisVersion.patch < gitVersion.patch) {
			link, err := url.Parse(web)
			if err != nil {
				return
			}
			fyne.Do(func() {
				if notify {
					SendNotification(lang.X("update.notify.title", "New version"), fmt.Sprintf(lang.X("update.notify.msg", "New version %s is available"), newVer))
					Gui.Settings.LastUpdatecheck = time.Now().Unix()
					Gui.Settings.Store()
				} else {
					msg := widget.NewHyperlinkWithStyle(fmt.Sprintf(lang.X("update.msg", "A new version %s is available !"), newVer),
						link, fyne.TextAlignCenter, fyne.TextStyle{
							Bold: true,
						})
					var dia *dialog.CustomDialog
					ok := widget.NewButton(lang.X("ok", "Ok"), func() {
						dia.Hide()
					})
					dia = dialog.NewCustomWithoutButtons(lang.X("update.title", "Update"),
						container.NewVBox(msg, util.NewVFiller(2), ok), Gui.MainWindow)
					dia.Show()
					dia.Resize(fyne.NewSize(Gui.MainWindow.Canvas().Size().Width, dia.MinSize().Height))
				}
			})
		} else {
			if !notify {
				fyne.Do(func() {
					dialog.ShowInformation(lang.X("update.title", "Update"), lang.X("update.nonew", "You are alread running the latest version."), Gui.MainWindow)
				})
			}
		}
	}()
}

func Export() {
	doPasswordCheck(func(ok bool) {
		doExport(func(err error) {
			if err != nil {
				UIErrorHandler(err)
			} else {
				SendNotification(lang.X("export.ok.title", "Export"), lang.X("export.ok.msg", "Data were successfully exported !"))
			}
		})
	})
}

func doExport(fDone func(error)) {
	w := func(writer fyne.URIWriteCloser) {
		json, err := Database.ExportJson()
		if err != nil {
			return
		}
		err = util.WriteFileToStorage(writer, []byte(json))
		if fDone != nil {
			fDone(err)
		}
	}

	ShowExportTypeDialog(func(overwrite bool) {
		var dia *dialog.FileDialog
		if overwrite {
			dia = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err != nil {
					if fDone != nil {
						fDone(err)
					}
					return
				}
				if reader == nil {
					return
				}
				defer reader.Close()
				writer, err := storage.Writer(reader.URI())
				if err != nil {
					if fDone != nil {
						fDone(err)
					}
					return
				}
				defer writer.Close()
				w(writer)
			}, Gui.MainWindow)
		} else {
			dia = dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					if fDone != nil {
						fDone(err)
					}
					return
				}
				if writer == nil {
					return
				}
				defer writer.Close()
				w(writer)
			}, Gui.MainWindow)
		}
		dia.SetView(dialog.ListView)
		if !overwrite {
			fName := path.Base(Gui.DatabaseFile)
			fName = strings.TrimSuffix(fName, path.Ext(Gui.DatabaseFile)) + ".json"
			dia.SetFileName(fName)
		}
		if Gui.IsDesktop {
			filter := storage.NewExtensionFileFilter([]string{".json"})
			dia.SetFilter(filter)
		}
		dia.Show()
		si := Gui.MainWindow.Canvas().Size()
		var windowScale float32 = 1.0
		dia.Resize(fyne.NewSize(si.Width*windowScale, si.Height*windowScale))
	})
}

func Import() {
	doPasswordCheck(func(ok bool) {
		if !ok {
			return
		}
		n, _ := Database.GetNumberOfEntries()
		dia := dialog.NewConfirm(lang.X("import.confirm.title", "Import data"),
			fmt.Sprintf(lang.X("import.confirm.msg", "Import data and drop database with %d entries ?"), n), func(ok bool) {
				if !ok {
					return
				}
				doImport(func(err error) {
					if err != nil {
						UIErrorHandler(err)
					} else {
						SendNotification(lang.X("import.ok.title", "Import"), lang.X("import.ok.msg", "Data were successfully imported !"))
						UpdateCategoryData()
						Gui.mainView.SetDate(nil)
					}
				})
			}, Gui.MainWindow)
		dia.SetConfirmImportance(widget.DangerImportance)
		dia.Show()
	})
}

func doImport(fDone func(error)) {
	dia := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			if fDone != nil {
				fDone(err)
			}
			return
		}
		if reader == nil {
			if fDone != nil {
				fDone(errors.New("reader is nil"))
			}
			return
		}
		defer reader.Close()
		data, err := io.ReadAll(reader)
		if err != nil {
			if fDone != nil {
				fDone(err)
			}
			return
		}
		err = Database.ImportJson(string(data))
		if fDone != nil {
			fDone(err)
		}
	}, Gui.MainWindow)
	dia.SetView(dialog.ListView)
	if Gui.IsDesktop {
		filter := storage.NewExtensionFileFilter([]string{".json"})
		dia.SetFilter(filter)
	}
	dia.Show()
	si := Gui.MainWindow.Canvas().Size()
	var windowScale float32 = 1.0
	dia.Resize(fyne.NewSize(si.Width*windowScale, si.Height*windowScale))
}
