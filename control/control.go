package control

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
	"time"

	"github.com/Igor-Basilio/lexis/async"
	CONST "github.com/Igor-Basilio/lexis/constant"
	h "github.com/Igor-Basilio/lexis/helper"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var Cur_line int32 = 1
var Cur_col int32 = 0
var Spacing float32 = 0
var FIRST_KEY_PRESSED int32

func Control_manager(camera *rl.Camera2D,
	cursor *rl.Vector2, sc *rl.Color, c map[int]CONST.Data, w *bool,
	fonts *[CONST.NUMBER_OF_FONTS]rl.Font) {

	FIRST_KEY_PRESSED = rl.GetKeyPressed()
	y_off := fonts[CONST.SELECTED_FONT].Chars.Image.Height
	h_off := fonts[CONST.SELECTED_FONT].Chars.Image.Width

	if rl.IsKeyDown(rl.KeyLeftControl) &&
		(FIRST_KEY_PRESSED == 70 || FIRST_KEY_PRESSED == 102) {

		if CONST.DRAW_SEARCH_BOX {
			*cursor = rl.Vector2{X: float32(Cur_col*h_off + Cur_col*int32(Spacing)),
				Y: float32((Cur_line - 1) * y_off)}
		} else {
			*cursor = rl.Vector2{X: 0,
				Y: float32(rl.GetScreenHeight()) -
					float32(y_off) - 2*float32(y_off) - 0.25*float32(y_off)}

		}

		async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
		*sc = rl.Yellow
		CONST.DRAW_SEARCH_BOX = !CONST.DRAW_SEARCH_BOX

	}

	if !CONST.DRAW_SEARCH_BOX {

		control_zoom(camera)
		control_camera(camera, fonts /* , c */)
		arr_keyMovement(cursor, sc, c, fonts)
		textManagementFunctionalities(c, cursor, fonts, sc)
		fileFunctionalities(c, fonts)

	}

}

func fileFunctionalities(c map[int]CONST.Data,
	fonts *[CONST.NUMBER_OF_FONTS]rl.Font) {

	font_height := fonts[CONST.SELECTED_FONT].Chars.Image.Height
	var num_bytes int = 0
	var has_written = false

	if rl.IsKeyDown(rl.KeyLeftControl) &&
		(FIRST_KEY_PRESSED == 83 || FIRST_KEY_PRESSED == 115) {

		user_infile_group := false
		f_info, err := os.Stat(CONST.CURRENT_FILE)
		h.Check(err)
		m := f_info.Mode()

		user, err := user.Current()
		h.Check(err)

		u_groups, err := user.GroupIds()
		h.Check(err)

		file_sys := f_info.Sys()

		file_gid := fmt.Sprint(file_sys.(*syscall.Stat_t).Gid)
		file_uid := fmt.Sprint(file_sys.(*syscall.Stat_t).Uid)

		// TODO: Use binary search or some other search algorithm
		// to speed the search for a group id the user is in.
		for _, gid := range u_groups {
			if gid == file_gid {
				user_infile_group = true
				break
			}
		}

		if m&(CONST.WRITE_GROUP_PERM) != 0 && user_infile_group ||
			m&(CONST.WRITE_OTHERS_PERM) != 0 ||
			(file_uid == user.Uid) && m&(CONST.WRITE_OWNER_PERM) != 0 {

			f, err := os.Create(CONST.CURRENT_FILE)
			h.Check(err)
			defer f.Close()

			for counter := 0; counter <= len(c); counter++ {

				if d, ok := c[counter]; ok {
					n, err := f.WriteString(d.Line + "\n")
					h.Check(err)
					num_bytes += n
				}

			}

            err = f.Sync()
            h.Check(err)
			has_written = true

		} else {
			has_written = false
			permission := true
			h.DrawTextForSpecifiedTime(
				"Cannot write, no permission.",
				1*time.Second+500*time.Millisecond,
				rl.Vector2{X: 0, Y: float32(rl.GetScreenHeight()) - float32(font_height)},
				&permission, fonts[CONST.SELECTED_FONT], rl.Black, Spacing)
		}

	}

	h.DrawTextForSpecifiedTime(
		" Bytes Written : "+fmt.Sprint(num_bytes)+" B",
		1*time.Second+500*time.Millisecond,
		rl.Vector2{X: 0, Y: float32(rl.GetScreenHeight()) - float32(font_height)},
		&has_written, fonts[CONST.SELECTED_FONT], rl.Black, Spacing)

}

func control_zoom(camera *rl.Camera2D) {

	if rl.IsKeyDown(rl.KeyLeftControl) &&
		rl.IsKeyDown(rl.KeyMinus) {

		var zoomIncrement float32 = 0.125
		camera.Zoom += zoomIncrement

	} else if rl.IsKeyDown(rl.KeyLeftShift) &&
		rl.IsKeyDown(rl.KeyMinus) {

		var zoomIncrement float32 = 0.125
		camera.Zoom -= zoomIncrement
	}

}

/* func GetAmountOfCharsInFirstScreen(fonts *[CONST.NUMBER_OF_FONTS]rl.Font,
	c map[int]CONST.Data) int {

	var n int
	sw := rl.GetScreenWidth()

	for sw != 0 {

	}

	return n
} */

// This works based on the height of the font used and raylibs GetScreenHeight() function
func GetAmountOfLinesOnScreen(camera *rl.Camera2D, fonts *[CONST.NUMBER_OF_FONTS]rl.Font) int32 {

	y_off := float32(fonts[CONST.SELECTED_FONT].Chars.Image.Height)

	logo_size := int(2 * y_off)
	sh := rl.GetScreenHeight() - logo_size
	var counter int32 = 0

	for sh > 0 {

		sh -= int(y_off)
		counter++

	}

	return counter

}

func ScrollCameraDownOneCharacter(camera *rl.Camera2D,
	fonts *[CONST.NUMBER_OF_FONTS]rl.Font) {

	y_off := float32(fonts[CONST.SELECTED_FONT].Chars.Image.Height)

	camera.Target.Y += y_off
	CONST.END_POINT_POSITION += int(y_off)
	CONST.SCROLLED_COUNT++
	CONST.ScrolledBottom = true

}

// Scroll functionalities, size functionalities for camera2D
func control_camera(camera *rl.Camera2D,
	fonts *[CONST.NUMBER_OF_FONTS]rl.Font /* , c map[int]CONST.Data */) {

	y_off := float32(fonts[CONST.SELECTED_FONT].Chars.Image.Height)
	/* h_off := fonts[CONST.SELECTED_FONT].Chars.Image.Width */

	pred_size := int((Cur_line) * int32(y_off))
	logo_size := int(2 * y_off)

	/* 	n_off := int32(h.GetAmountOfDigits(len(c))) */

	// THIS WORKS AND I DON'T KNOW WHY ?
	// VERTICAL CAMERA CONTROL
	if pred_size > (rl.GetScreenHeight()-logo_size+CONST.SCROLLED_COUNT*int(y_off)) &&
		!CONST.ScrolledBottom {

		camera.Target.Y += y_off
		CONST.ScrolledBottom = true
		CONST.END_POINT_POSITION += int(y_off)
		CONST.SCROLLED_COUNT++

	} else if (float32(Cur_line)*y_off < camera.Target.Y) && !CONST.ScrolledTop {

		camera.Target.Y -= y_off
		CONST.ScrolledTop = true
		CONST.SCROLLED_COUNT--

	} else if Cur_line == 1 {

		camera.Target.Y = 0
		CONST.SCROLLED_COUNT = 0

		// HORIZONTAL CAMERA CONTROL
		// FUCKNIG BUG FEST MAN
	} /* else if d, ok := c[int(Cur_line)]; ok &&
		(int32(rl.GetScreenWidth())-(n_off+4)*h_off+int32(CONST.SCROLLED_COUNT_H)*h_off <
			int32(len(d.Line[:Cur_col]))*h_off+int32(Spacing)*int32(len(d.Line[:Cur_col]))) &&
		!CONST.ScrolledRight {

		camera.Target.X += float32(h_off) + Spacing
		CONST.SCROLLED_COUNT_H++
		CONST.ScrolledRight = true

	} else if _, ok := c[int(Cur_line)]; ok &&
		Cur_col == int32(CONST.SCROLLED_COUNT_H-1) {

		camera.Target.X -= float32(h_off) + Spacing
		CONST.SCROLLED_COUNT_H--
		CONST.ScrolledLeft = true

	} else if Cur_col == 0 && !CONST.ScrolledLeft {

		camera.Target.X = 0
		CONST.SCROLLED_COUNT_H = 0
		CONST.ScrolledLeft = true

	} */

}

func arr_keyMovement(cursor *rl.Vector2, sc *rl.Color, c map[int]CONST.Data,
	fonts *[CONST.NUMBER_OF_FONTS]rl.Font) {

	h_off := float32(fonts[CONST.SELECTED_FONT].Chars.Image.Width)
	y_off := float32(fonts[CONST.SELECTED_FONT].Chars.Image.Height)

	if rl.IsKeyPressed(rl.KeyRight) {

		pred_curcol := Cur_col + 1
		d, ok := c[int(Cur_line)]

		if ok {
			if int(pred_curcol) <= len(d.Line) {

				Cur_col++
				cursor.X += h_off + Spacing

				async.Debounce_Ticker.Reset(async.DEBOUNCE_TIMER)
				CONST.DEBOUNCE_MOVER = true
				async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
				*sc = rl.Yellow
				CONST.ScrolledRight = false

			}
		}

	} else if rl.IsKeyPressed(rl.KeyDown) {

		pred_curline := Cur_line + 1
		d, ok := c[int(pred_curline)]
		k, ok_k := c[int(Cur_line)]

		if ok && ok_k && Cur_line <= int32(len(c)) {

			if len(d.Line) <= len(k.Line[:Cur_col]) {

				Cur_col = int32(len(d.Line))
				*cursor = rl.Vector2{X: float32(Cur_col)*h_off + float32(Cur_col)*Spacing,
					Y: float32(Cur_line) * y_off}

			} else {
				cursor.Y += y_off
			}

			d.Selected = true
			k.Selected = false

			c[int(Cur_line)] = k
			c[int(pred_curline)] = d

			async.Debounce_Ticker.Reset(async.DEBOUNCE_TIMER)
			CONST.ScrolledBottom = false
			Cur_line++
			async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
			*sc = rl.Yellow

		}

	} else if rl.IsKeyPressed(rl.KeyUp) && Cur_line >= 2 {

		d, ok_d := c[int(Cur_line)]
		l, ok_l := c[int(Cur_line-1)]

		if ok_d && ok_l {

			if len(l.Line) <= len(d.Line[:Cur_col]) {

				Cur_col = int32(len(l.Line))
				*cursor = rl.Vector2{X: float32(Cur_col)*h_off + float32(Cur_col)*Spacing,
					Y: float32(Cur_line-2) * y_off}

			} else {
				cursor.Y -= y_off
			}

			l.Selected = true
			d.Selected = false

			c[int(Cur_line)] = d
			c[int(Cur_line-1)] = l

			async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
			*sc = rl.Yellow
			CONST.ScrolledTop = false
			Cur_line--

		}

	} else if rl.IsKeyPressed(rl.KeyLeft) && Cur_col >= 1 {

		cursor.X -= h_off + Spacing
		async.Debounce_Ticker.Reset(async.DEBOUNCE_TIMER)
		CONST.DEBOUNCE_MOVEL = true
		Cur_col--
		async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
		*sc = rl.Yellow
		CONST.ScrolledLeft = false

	}

}

func textManagementFunctionalities(c map[int]CONST.Data, cursor *rl.Vector2,
	fonts *[CONST.NUMBER_OF_FONTS]rl.Font, sc *rl.Color) {

	h_off := fonts[CONST.SELECTED_FONT].Chars.Image.Width
	y_off := fonts[CONST.SELECTED_FONT].Chars.Image.Height

	// Add newline on ENTER key.
	if rl.IsKeyPressed(rl.KeyEnter) {

		for i := len(c); i > int(Cur_line); i-- {

			str, ok := c[int(i)]
			if ok {
				c[int(i)+1] = str
			}

		}

		var d2 CONST.Data
		d, ok := c[int(Cur_line)]
		if ok {

			d2 = CONST.Data{
				Line:     d.Line[Cur_col:],
				Selected: true,
			}

			d.Line = d.Line[:Cur_col]
			d.Selected = false

		}
		c[int(Cur_line+1)] = d2
		c[int(Cur_line)] = d

		async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
		*sc = rl.Yellow
		*cursor = rl.Vector2{X: 0, Y: float32(Cur_line * y_off)}
		Cur_col = 0
		Cur_line++
		CONST.ScrolledBottom = false

		// Concatenate current line with previous on remove.
	} else if rl.IsKeyPressed(rl.KeyBackspace) {

		CONST.ScrolledTop = false

		if Cur_col == 0 {

			if Cur_line > 1 {

				d, ok := c[int(Cur_line-1)]

				// Graphical part
				if ok {
					Cur_col = int32(len(d.Line))
					*cursor = rl.Vector2{X: float32(Cur_col*h_off + Cur_col*int32(Spacing)),
						Y: float32((Cur_line - 2) * y_off)}
				}

				l, ok := c[int(Cur_line)]

				// Backend part
				if ok {
					d.Line += l.Line
					d.Selected = true
					c[int(Cur_line-1)] = d
				}

				for i := int(Cur_line); i < len(c); i++ {
					d, ok := c[int(i)+1]
					if ok {
						c[int(i)] = d
					}
				}

				async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
				*sc = rl.Yellow

				delete(c, len(c))
				Cur_line--

			}
		} else {

			d, ok := c[int(Cur_line)]

			if Cur_col <= int32(len(d.Line)) && ok {
				deleteCharAtCol(Cur_line, Cur_col-1, c)
				cursor.X -= float32(h_off) + Spacing
				Cur_col--
			}

			async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
			*sc = rl.Yellow

		}

	} else if rl.IsKeyPressed(rl.KeySpace) {

		d, ok := c[int(Cur_line)]

		if ok {

			cline_straf := d.Line[Cur_col:]
			cline_strbf := d.Line[:Cur_col] + " "
			d.Line = (cline_strbf + cline_straf)

			Cur_col++
			cursor.X += float32(h_off) + Spacing
			c[int(Cur_line)] = d

			async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
			*sc = rl.Yellow

		}

	} else if rl.IsKeyPressed(rl.KeyLeft) && !CONST.DEBOUNCE_MOVEL {

		if Cur_col == 0 && Cur_line > 1 {

			d, ok := c[int(Cur_line-1)]
			l, ok_l := c[int(Cur_line)]

			if ok && ok_l {

				d.Selected = true
				l.Selected = false

				c[int(Cur_line-1)] = d
				c[int(Cur_line)] = l

				Cur_col = int32(len(d.Line))
				*cursor = rl.Vector2{X: float32(Cur_col*h_off + Cur_col*int32(Spacing)),
					Y: float32((Cur_line - 2) * y_off)}
				Cur_line--
			}

			async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
			*sc = rl.Yellow

		}

	} else if rl.IsKeyPressed(rl.KeyRight) && !CONST.DEBOUNCE_MOVER {

		d, ok := c[int(Cur_line)]
		l, ok_l := c[int(Cur_line+1)]

		if Cur_col == int32(len(d.Line)) && (ok && ok_l) {

			d.Selected = false
			l.Selected = true

			c[int(Cur_line+1)] = l
			c[int(Cur_line)] = d

			Cur_col = 0
			*cursor = rl.Vector2{X: 0,
				Y: float32((Cur_line) * y_off)}
			Cur_line++

		}

		async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
		*sc = rl.Yellow

	} else if b, r := IsAnyKeyPressed(); b {

		d, ok := c[int(Cur_line)]

		if ok {
			d.Line = d.Line[:Cur_col] + string(r) + d.Line[Cur_col:]
			c[int(Cur_line)] = d
			Cur_col++
			cursor.X += float32(h_off) + Spacing
		}

		async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
		*sc = rl.Yellow

	} else if rl.IsKeyDown(rl.KeyLeftControl) &&
		(FIRST_KEY_PRESSED <= 57 && FIRST_KEY_PRESSED >= 48) {

		switch FIRST_KEY_PRESSED {
		case 49:
			CONST.SELECTED_FONT = CONST.SPACE_MONO
		case 50:
			CONST.SELECTED_FONT = CONST.ROBOTO_MONO
		default:
			CONST.SELECTED_FONT = CONST.FIRA_CODE
		}

		// Reset cursor since sizes changed.
		h_off := fonts[CONST.SELECTED_FONT].Chars.Image.Width
		y_off := fonts[CONST.SELECTED_FONT].Chars.Image.Height
		*cursor = rl.Vector2{
			X: float32(Cur_col)*float32(h_off) + float32(Cur_col)*Spacing,
			Y: float32(Cur_line-1) * float32(y_off)}

		async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
		*sc = rl.Yellow

	}

}

// Returns a boolean and the rune key, given a pressed key that is in the
// constant range preemptively specified.
func IsAnyKeyPressed() (bool, rune) {

	keyPressed := false
	ch := rl.GetCharPressed()

	if (ch >= 33) && (ch <= 126) {
		keyPressed = true
	}

	return keyPressed, ch

}

// Pass in a line and column and this function deletes the character at
// c[int(line)][col], implemented using go slices.
func deleteCharAtCol(line int32, col int32,
	c map[int]CONST.Data) map[int]CONST.Data {

	d, ok := c[int(line)]

	if ok {
		later_part := d.Line[col+1:]
		first_part := d.Line[:col]
		d.Line = (first_part + later_part)
		c[int(line)] = d
	}

	return c

}