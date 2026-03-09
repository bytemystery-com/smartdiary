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

package mytheme

import (
	"fmt"
	"image/color"
	"log/slog"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const Scaling = 1.3

type MyThemeInterface interface {
	GetSpecialColor(string) color.Color
	GetSpecialSize(string) float32
}

type MyTheme struct {
	base    fyne.Theme
	variant fyne.ThemeVariant
}

func (p *MyTheme) GetVariant() fyne.ThemeVariant {
	return p.variant
}

func (p *MyTheme) SetVariant(variant fyne.ThemeVariant) {
	p.variant = variant
}

func NewMyTheme(variant fyne.ThemeVariant) *MyTheme {
	return &MyTheme{
		base:    theme.DefaultTheme(),
		variant: variant,
	}
}

func (p *MyTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	if p.variant == theme.VariantDark {
		switch c {
		case theme.ColorNameBackground:
			return color.NRGBA{25, 25, 25, 255}
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 50, G: 50, B: 50, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{90, 90, 90, 255}
		case theme.ColorNameFocus:
			return color.NRGBA{R: 50, G: 145, B: 245, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 60, G: 60, B: 60, A: 255}
		case theme.ColorNamePrimary:
			return color.NRGBA{R: 87, G: 139, B: 255, A: 255}
			// return color.NRGBA{R: 41, G: 111, B: 246, A: 255}
		case theme.ColorNameSelection:
			return color.NRGBA{R: 132, G: 173, B: 255, A: 255}
		case theme.ColorNameError:
			// return color.NRGBA{{244, 67, 54, 255}
			return color.NRGBA{R: 255, G: 104, B: 82, A: 255}
		case theme.ColorNameForeground:
			return color.NRGBA{R: 230, G: 230, B: 230, A: 255}
		}
	} else {
		switch c {
		case theme.ColorNameBackground:
			return color.NRGBA{247, 247, 247, 255}
		case theme.ColorNameOverlayBackground:
			return color.NRGBA{R: 230, G: 230, B: 230, A: 255}
		case theme.ColorNameInputBackground:
			return color.NRGBA{235, 235, 235, 255}
		case theme.ColorNameFocus:
			return color.NRGBA{R: 160, G: 200, B: 242, A: 255}
		case theme.ColorNameButton:
			return color.NRGBA{R: 225, G: 225, B: 225, A: 255}
		case theme.ColorNameSelection:
			return color.NRGBA{R: 191, G: 225, B: 255, A: 255}

		}
	}
	/*
		val := p.base.Color(c, p.variant)
		if c == theme.ColorNameForeground {
			fmt.Println(val)
		}
		return val
	*/
	return p.base.Color(c, p.variant)
}

func (p *MyTheme) Font(s fyne.TextStyle) fyne.Resource {
	return p.base.Font(s)
}

func (p *MyTheme) Icon(i fyne.ThemeIconName) fyne.Resource {
	return p.base.Icon(i)
}

func (p *MyTheme) Size(s fyne.ThemeSizeName) float32 {
	val := p.base.Size(s)
	switch s {
	case theme.SizeNameSubHeadingText, theme.SizeNameHeadingText, theme.SizeNameCaptionText, theme.SizeNameText:
		return val * Scaling
	}
	return val
}

func (p *MyTheme) parseHexColor(s string) (color.NRGBA, error) {
	c := color.NRGBA{A: 255}
	if s[0] != '#' {
		return c, fmt.Errorf("invalid format")
	}

	var err error
	switch len(s) {
	case 7: // #RRGGBB
		r, _ := strconv.ParseUint(s[1:3], 16, 8)
		g, _ := strconv.ParseUint(s[3:5], 16, 8)
		b, _ := strconv.ParseUint(s[5:7], 16, 8)
		c.R = uint8(r)
		c.G = uint8(g)
		c.B = uint8(b)
	case 9: // #RRGGBBAA
		r, _ := strconv.ParseUint(s[1:3], 16, 8)
		g, _ := strconv.ParseUint(s[3:5], 16, 8)
		b, _ := strconv.ParseUint(s[5:7], 16, 8)
		a, _ := strconv.ParseUint(s[7:9], 16, 8)
		c.R = uint8(r)
		c.G = uint8(g)
		c.B = uint8(b)
		c.A = uint8(a)
	default:
		return c, fmt.Errorf("invalid length")
	}
	return c, err
}

func (p *MyTheme) GetSpecialColor(c string) color.Color {
	if p.variant == theme.VariantDark {
		switch c {
		case "txt_color":
			return color.NRGBA{247, 247, 247, 255}
		case "txt_color_other":
			return color.NRGBA{155, 155, 155, 255}

		case "txt_color_special_0":
			return color.NRGBA{255, 0, 0, 255}
		case "txt_color_special_0_other":
			return color.NRGBA{215, 80, 80, 255}

		case "txt_color_special_1":
			return color.NRGBA{17, 130, 255, 255}
		case "txt_color_special_1_other":
			return color.NRGBA{17, 97, 192, 255}

		case "txt_color_special_2":
			return color.NRGBA{0, 170, 0, 255}
		case "txt_color_special_2_other":
			return color.NRGBA{77, 135, 77, 255}

		case "bg_color_calendar":
			return color.NRGBA{50, 50, 50, 255}
		case "bg_color_header":
			return color.NRGBA{178, 178, 178, 255}

		case "category_text_color":
			return color.NRGBA{70, 70, 70, 255}
		case "nocategory_text_color":
			return color.NRGBA{190, 190, 190, 255}

		case "txt_color_bl":
			return color.NRGBA{0, 0, 0, 255}
		case "txt_color_bl_other":
			return color.NRGBA{40, 40, 40, 255}
		}
	} else {
		switch c {
		case "txt_color":
			return color.NRGBA{10, 10, 10, 255}
		case "txt_color_other":
			return color.NRGBA{120, 120, 120, 255}

		case "txt_color_special_0":
			return color.NRGBA{255, 0, 0, 255}
		case "txt_color_special_0_other":
			return color.NRGBA{202, 92, 92, 255}

		case "txt_color_special_1":
			return color.NRGBA{0, 0, 255, 255}
		case "txt_color_special_1_other":
			return color.NRGBA{94, 94, 220, 255}

		case "txt_color_special_2":
			return color.NRGBA{0, 127, 0, 255}
		case "txt_color_special_2_other":
			return color.NRGBA{95, 155, 95, 255}

		case "bg_color_calendar":
			return color.NRGBA{240, 240, 240, 255}
		case "bg_color_header":
			return color.NRGBA{150, 150, 150, 255}

		case "category_text_color":
			return color.NRGBA{70, 70, 70, 255}
		case "nocategory_text_color":
			return color.NRGBA{70, 70, 70, 255}

		case "txt_color_bl":
			return color.NRGBA{0, 0, 0, 255}
		case "txt_color_bl_other":
			return color.NRGBA{85, 85, 85, 255}
		}
	}
	col, err := p.parseHexColor(c)
	if err == nil {
		return col
	} else {
		return p.Color(fyne.ThemeColorName(c), p.variant)
	}
}

func (p *MyTheme) GetSpecialSize(s string) float32 {
	switch s {
	case "categorie_select_popup":
		return 1.3
	case "categorie_select":
		return 1.0
	case "overlay_size_scale":
		return 0.3
	case "login_icon_size":
		return 135
	case "login_space_logo_label":
		return 5
	case "login_space_field_ok":
		return 3
	default:
		slog.Error("Unknown special size", "value", s)
	}
	return 1.0
}
