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

package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"bytemystery-com/smartdiary/daywidget"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	_ "modernc.org/sqlite"
)

const (
	CURRENT_VERSION              = 1
	DEF_NUMBER_OF_CATEGORIES     = 15
	DEF_NUMBER_OF_SEARCH_RESULTS = 25
)

type Db struct {
	sql     *sql.DB
	prepStm map[string]*sql.Stmt
}

type CategoryDataType struct {
	Id         int64     `json:"id"`
	Name       string    `json:"n"`
	Color      string    `json:"c"`
	OrderIndex int       `json:"o"`
	Timestamp  time.Time `json:"t"`
	IsEdited   bool      `json:"-"`
}

type EntryDataSimpleType struct {
	Id          int64                   `json:"id"`
	CategoryId1 int64                   `json:"c1"`
	CategoryId2 sql.NullInt64           `json:"c2"`
	ColorMode   daywidget.ColorModeType `json:"cm"`
	Mark        int                     `json:"m"`
	Protected   bool                    `json:"p"`
}

type EntryDataType struct {
	EntryDataSimpleType `json:"e"`
	Date                time.Time `json:"d"`
	Data                string    `json:"n"`
	Timestamp           time.Time `json:"t"`
}

func GetDBFile(name string) (string, error) {
	app := fyne.CurrentApp()
	root := app.Storage().RootURI()
	dbUri, err := storage.Child(root, name+".db")
	if err != nil {
		return "", err
	}
	return filepath.FromSlash(dbUri.Path()), nil
}

func NewDb() *Db {
	return &Db{
		prepStm: make(map[string]*sql.Stmt, 10),
	}
}

func (d *Db) create() error {
	// Config
	_, err := d.sql.Exec(`CREATE TABLE config (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT DEFAULT "", value TEXT DEFAULT "", ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		return err
	}
	_, err = d.sql.Exec(`CREATE TRIGGER config_set_ts_update AFTER UPDATE ON config FOR EACH ROW BEGIN UPDATE config SET ts = CURRENT_TIMESTAMP WHERE id=OLD.id; END`)
	if err != nil {
		return err
	}
	_, err = d.sql.Exec(fmt.Sprintf(`INSERT INTO config (name, value) VALUES ('version', '%d')`, CURRENT_VERSION))
	if err != nil {
		return err
	}

	// Category
	_, err = d.sql.Exec(`CREATE TABLE category (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT DEFAULT '' NOT NULL, color TEXT NOT NULL, orderIndex INTEGER NOT NULL, ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP)`)
	if err != nil {
		return err
	}
	_, err = d.sql.Exec(`CREATE INDEX idx_category_orderindex ON category(orderIndex ASC)`)
	if err != nil {
		return err
	}
	_, err = d.sql.Exec(`CREATE TRIGGER category_set_ts_update AFTER UPDATE ON category FOR EACH ROW BEGIN UPDATE category SET ts = CURRENT_TIMESTAMP WHERE id=OLD.id; END`)
	if err != nil {
		return err
	}
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 1', '#FFFF8CFF', 100)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 2', '#99CCFFFF', 200)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 3', '#FFC2C2FF', 300)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 4', '#D4FFB5FF', 400)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 5', '#CCFFFFFF', 500)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 6', '#D6B9F4FF', 600)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 7', '#E3A86FFF', 700)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 8', '#C5C5C5FF', 800)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 9', '#92C392FF', 900)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 10', '#FFC41AFF', 1000)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 11', '#fe6d6eFF', 1100)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 12', '#FFFFFFFF', 1200)`)
	_, err = d.sql.Exec(`INSERT INTO category (name, color, orderIndex) VALUES ('Category 13', '#ff7aebFF', 1300)`)
	if err != nil {
		return err
	}

	// Entry
	_, err = d.sql.Exec(`CREATE TABLE entry (id INTEGER PRIMARY KEY AUTOINCREMENT, date TEXT NOT NULL, categoryId1 INTEGER NOT NULL, categoryId2 INTEGER DEFAULT NULL, data TEXT DEFAULT '' NOT NULL, colorMode INTEGER NOT NULL DEFAULT 0, mark INTEGER NOT NULL DEFAULT 0, protected BOOL NOT NULL DEFAULT FALSE, ts DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (categoryId1) REFERENCES category(id), FOREIGN KEY (categoryId2) REFERENCES category(id))`)
	if err != nil {
		return err
	}
	_, err = d.sql.Exec(`CREATE UNIQUE INDEX idx_entry_date ON entry(date ASC)`)
	if err != nil {
		return err
	}
	/*
		_, err = d.sql.Exec(`CREATE INDEX idx_entry_data_categoryId1 ON entry(data ASC, categoryId1)`)
		if err != nil {
			return err
		}
	*/
	_, err = d.sql.Exec(`CREATE TRIGGER entry_set_ts_update AFTER UPDATE ON entry FOR EACH ROW BEGIN UPDATE entry SET ts = CURRENT_TIMESTAMP WHERE id=OLD.id; END`)
	if err != nil {
		return err
	}

	return nil
}

func (d *Db) update(oldVer int) error {
	return nil
}

func (d *Db) Open(file string) error {
	d.Close()
	b, err := sql.Open("sqlite", fmt.Sprintf("file:%s?pragma=foreign_keys(1)", file))
	if err != nil {
		return err
	}
	d.sql = b
	v := ""
	q, err := d.sql.Query(`SELECT value from config WHERE name = 'version'`)
	if err != nil {
		return d.create()
	}
	defer q.Close()
	if !q.Next() {
		return errors.New("no version information")
	}
	err = q.Scan(&v)
	if err != nil {
		return err
	}
	ver, err := strconv.Atoi(v)
	if err != nil {
		return err
	}
	if ver == CURRENT_VERSION {
		return nil
	} else {
		return d.update(ver)
	}
}

func (d *Db) IsOpen() bool {
	return d.sql != nil
}

func (d *Db) Close() {
	if d.sql != nil {
		for _, item := range d.prepStm {
			item.Close()
		}
		clear(d.prepStm)
		d.sql.Close()
		d.sql = nil
	}
}

func (d *Db) UpdateCategory(data *CategoryDataType) error {
	if !d.IsOpen() {
		return errors.New("database is not open")
	}
	var ok bool
	var p *sql.Stmt
	var err error
	p, ok = d.prepStm["update_category"]
	if !ok {
		p, err = d.sql.Prepare(`UPDATE category SET name=?, color=?, orderIndex=? WHERE id=?`)
		if err != nil {
			return err
		}
		d.prepStm["update_category"] = p
	}
	_, err = p.Exec(data.Name, data.Color, data.OrderIndex, data.Id)
	return err
}

func (d *Db) GetNumberOfCategories() (int, error) {
	if !d.IsOpen() {
		return 0, errors.New("database is not open")
	}
	r := d.sql.QueryRow(`SELECT COUNT(*) FROM category`)
	count := 0
	err := r.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (d *Db) GetNumberOfEntries() (int, error) {
	if !d.IsOpen() {
		return 0, errors.New("database is not open")
	}
	r := d.sql.QueryRow(`SELECT COUNT(*) FROM entry`)
	count := 0
	err := r.Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}

func (d *Db) GetCategories() ([]*CategoryDataType, error) {
	if !d.IsOpen() {
		return nil, errors.New("database not open")
	}
	q, err := d.sql.Query(`SELECT id, name, color, orderIndex, ts FROM category ORDER BY orderIndex ASC`)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	list := make([]*CategoryDataType, 0, DEF_NUMBER_OF_CATEGORIES)
	for q.Next() {
		data := CategoryDataType{}
		err := q.Scan(&data.Id, &data.Name, &data.Color, &data.OrderIndex, &data.Timestamp)
		if err == nil {
			list = append(list, &data)
		}
	}
	return list, nil
}

func (d *Db) InsertOrUpdateEntry(data *EntryDataType) error {
	if !d.IsOpen() {
		return errors.New("database not open")
	}
	var ok bool
	var p *sql.Stmt
	var err error
	p, ok = d.prepStm["upsert_by_date"]
	if !ok {
		p, err = d.sql.Prepare(`INSERT INTO entry (categoryId1,categoryId2,date,data,colorMode,mark,protected) VALUES (?,?,?,?,?,?, ?) ON CONFLICT(date) DO UPDATE SET
			categoryId1=excluded.categoryId1, categoryId2=excluded.categoryId2, data=excluded.data, colorMode=excluded.colorMode, mark=excluded.mark, protected=excluded.protected`)
		if err != nil {
			return err
		}
		d.prepStm["upsert_by_date"] = p
	}
	colorMode := daywidget.ColorModeToIndex[data.ColorMode]

	r, err := p.Exec(data.CategoryId1, data.CategoryId2, data.Date.Format(time.DateOnly), data.Data, colorMode, data.Mark, data.Protected)
	if err != nil {
		return err
	}
	i, err := r.LastInsertId()
	if err == nil {
		data.Id = i
	}
	return nil
}

func (d *Db) GetEntryByDate(date time.Time) (*EntryDataType, error) {
	if !d.IsOpen() {
		return nil, errors.New("database not open")
	}
	var ok bool
	var p *sql.Stmt
	var err error
	p, ok = d.prepStm["entry_by_date"]
	if !ok {
		p, err = d.sql.Prepare(`SELECT id, data, categoryId1, categoryId2, colorMode, mark, protected, ts FROM entry WHERE date=?`)
		if err != nil {
			return nil, err
		}
		d.prepStm["entry_by_date"] = p
	}
	r := p.QueryRow(date.Format(time.DateOnly))
	e := EntryDataType{}
	var colorMode int
	err = r.Scan(&e.Id, &e.Data, &e.CategoryId1, &e.CategoryId2, &colorMode, &e.Mark, &e.Protected, &e.Timestamp)
	if err != nil {
		return nil, err
	}
	e.Date = date
	cm, ok := daywidget.IndexToColorMode[colorMode]
	if ok {
		e.ColorMode = cm
	} else {
		e.ColorMode = daywidget.ColorMode_1r
	}
	return &e, nil
}

func (d *Db) GetSimpleEntryByDate(date time.Time) (*EntryDataSimpleType, error) {
	if !d.IsOpen() {
		return nil, errors.New("database not open")
	}
	var ok bool
	var p *sql.Stmt
	var err error
	p, ok = d.prepStm["simple_entry_by_date"]
	if !ok {
		p, err = d.sql.Prepare(`SELECT id, categoryId1, categoryId2, colorMode, mark, protected FROM entry WHERE date=?`)
		if err != nil {
			return nil, err
		}
		d.prepStm["simple_entry_by_date"] = p
	}
	r := p.QueryRow(date.Format(time.DateOnly))
	e := EntryDataSimpleType{}
	var colorMode int
	err = r.Scan(&e.Id, &e.CategoryId1, &e.CategoryId2, &colorMode, &e.Mark, &e.Protected)
	if err != nil {
		return nil, err
	}
	cm, ok := daywidget.IndexToColorMode[colorMode]
	if ok {
		e.ColorMode = cm
	} else {
		e.ColorMode = daywidget.ColorMode_1r
	}
	return &e, nil
}

func (d *Db) Search(text string) ([]*EntryDataType, error) {
	if !d.IsOpen() {
		return nil, errors.New("database not open")
	}
	var ok bool
	var p *sql.Stmt
	var err error
	p, ok = d.prepStm["search"]
	if !ok {
		p, err = d.sql.Prepare(`SELECT id, date, data, categoryId1, categoryId2, colorMode, mark, protected, ts FROM entry WHERE data LIKE ? COLLATE NOCASE ORDER BY date DESC`)
		if err != nil {
			return nil, err
		}
		d.prepStm["search"] = p
	}
	r, err := p.Query(text)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	list := make([]*EntryDataType, 0, DEF_NUMBER_OF_SEARCH_RESULTS)
	for r.Next() {
		e := EntryDataType{}
		colorMode := 0
		date := ""
		err := r.Scan(&e.Id, &date, &e.Data, &e.CategoryId1, &e.CategoryId2, &colorMode, &e.Mark, &e.Protected, &e.Timestamp)
		if err != nil {
			return nil, err
		}
		d, err := time.Parse(time.DateOnly, date)
		if err != nil {
			return nil, err
		}
		e.Date = d
		cm, ok := daywidget.IndexToColorMode[colorMode]
		if ok {
			e.ColorMode = cm
		} else {
			e.ColorMode = daywidget.ColorMode_1r
		}
		list = append(list, &e)
	}
	return list, nil
}

func (d *Db) DeleteEntry(id int64) error {
	if !d.IsOpen() {
		return errors.New("database not open")
	}
	var ok bool
	var p *sql.Stmt
	var err error
	p, ok = d.prepStm["delete_entry"]
	if !ok {
		p, err = d.sql.Prepare(`DELETE FROM entry WHERE id=?`)
		if err != nil {
			return err
		}
		d.prepStm["delete_entry"] = p
	}
	_, err = p.Exec(id)
	return err
}

func (d *Db) GetLastWrite() (time.Time, error) {
	if !d.IsOpen() {
		return time.Time{}, errors.New("database not open")
	}
	r := d.sql.QueryRow(`SELECT MAX(ts) FROM (SELECT ts FROM category UNION ALL SELECT ts FROM entry)`)
	var str string
	err := r.Scan(&str)
	if err != nil {
		return time.Time{}, err
	}
	t, err := time.Parse("2006-01-02 15:04:05", str)
	return t, err
}

type AllData struct {
	Categories []*CategoryDataType `json:"categories"`
	Entries    []*EntryDataType    `json:"entries"`
}

func (d *Db) ExportJson() (string, error) {
	if !d.IsOpen() {
		return "", errors.New("database not open")
	}
	var data AllData
	ca, err := d.GetCategories()
	if err != nil {
		return "", err
	}
	data.Categories = ca

	q, err := d.sql.Query(`SELECT id, categoryId1, categoryId2, colorMode, date, data, mark, protected, ts FROM entry ORDER BY date ASC`)
	if err != nil {
		return "", err
	}
	defer q.Close()

	data.Entries = make([]*EntryDataType, 0, 100)
	var colorMode int
	for q.Next() {
		entry := EntryDataType{}
		date := ""
		err := q.Scan(&entry.Id, &entry.CategoryId1, &entry.CategoryId2, &colorMode, &date, &entry.Data, &entry.Mark, &entry.Protected, &entry.Timestamp)
		if err != nil {
			return "", err
		}
		da, err := time.Parse(time.DateOnly, date)
		if err != nil {
			return "", err
		}
		entry.Date = da
		cm, ok := daywidget.IndexToColorMode[colorMode]
		if ok {
			entry.ColorMode = cm
		} else {
			entry.ColorMode = daywidget.ColorMode_1r
		}
		data.Entries = append(data.Entries, &entry)
	}
	json, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

func (d *Db) ImportJson(str string) error {
	if !d.IsOpen() {
		return errors.New("database not open")
	}
	var data AllData
	err := json.Unmarshal([]byte(str), &data)
	if err != nil {
		return err
	}

	tx, _ := d.sql.Begin()

	q, err := tx.Query(`DELETE FROM category`)
	if err != nil {
		tx.Rollback()
		return err
	}
	q.Close()

	q, err = tx.Query(`DELETE FROM entry`)
	if err != nil {
		tx.Rollback()
		return err
	}
	q.Close()

	p, err := tx.Prepare("INSERT INTO category (id, name, color, orderIndex, ts) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer p.Close()
	for _, item := range data.Categories {
		_, err = p.Exec(item.Id, item.Name, item.Color, item.OrderIndex, item.Timestamp)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	p, err = tx.Prepare("INSERT INTO entry (id, data, categoryId1, categoryId2, date, colorMode, mark, protected, ts) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer p.Close()

	for _, item := range data.Entries {
		colorMode := daywidget.ColorModeToIndex[item.ColorMode]
		_, err = p.Exec(item.Id, item.Data, item.CategoryId1, item.CategoryId2, item.Date.Format(time.DateOnly), colorMode, item.Mark, item.Protected, item.Timestamp)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

// Import from old SmartDiary
/*
func (d *Db) ImportFromCSV(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.Comma = ';'   // Trennzeichen ändern
	r.Comment = '#' // Kommentarzeilen
	r.TrimLeadingSpace = true

	for {
		record, err := r.Read()
		if err != nil {
			break
		}
		entry := EntryDataType{}
		entry.Date, err = time.Parse(time.DateOnly, record[0])
		if err != nil {
			return err
		}
		entry.Data = record[1]
		id, err := strconv.Atoi(record[2])
		if err != nil {
			return err
		}
		entry.CategoryId1 = int64(id + 1)

		protect, err := strconv.Atoi(record[3])
		if err != nil {
			return err
		}
		entry.Protected = (protect != 0)

		mark, err := strconv.Atoi(record[4])
		if err != nil {
			return err
		}
		entry.Mark = mark
		err = d.InsertOrUpdateEntry(&entry)
		if err != nil {
			return err
		}
	}
	return nil
}
*/
