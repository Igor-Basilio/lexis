package text

import (
	"fmt"

	CONST "github.com/Igor-Basilio/lexis/constant"
	control "github.com/Igor-Basilio/lexis/control"
	h "github.com/Igor-Basilio/lexis/helper"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var SELECTED_NUMBER rl.Color = rl.Black

// TODO: Allow support for non-monospaced fonts,
// Draws the text inputted as an argument to Lexis on the screen as well
// as the cursor and line numbers.
// ONLY WORKS WITH MONOSPACE FONTS FOR NOW !!!
func DrawFileText(c map[int]CONST.Data, selected_color *rl.Color,
	cursor rl.Vector2, fonts *[CONST.NUMBER_OF_FONTS]rl.Font, camera *rl.Camera2D) {

	h_off := float32(fonts[CONST.SELECTED_FONT].Chars.Image.Width)
	y_off := float32(fonts[CONST.SELECTED_FONT].Chars.Image.Height)

	n_offset := h.GetAmountOfDigits(len(c))

	// Draw file text
	position := rl.Vector2{X: float32(n_offset) * h_off, Y: y_off}
	default_cursor := rl.Vector2{X: position.X + 2*h_off, Y: 2 * position.Y}

	drawSelectedSquare(rl.Vector2Add(default_cursor, cursor),
		h_off, y_off, selected_color)

	for key, value := range c {

		if value.Selected {
			SELECTED_NUMBER = rl.White
		} else {
			SELECTED_NUMBER = rl.Black
		}

		n := h.GetAmountOfDigits(key)

		rl.DrawTextEx(fonts[CONST.SELECTED_FONT], fmt.Sprint(key),
			rl.Vector2{X: h_off*float32(n_offset) - h_off*float32(n), Y: position.Y + position.Y*float32(key)},
			float32(fonts[CONST.SELECTED_FONT].BaseSize), 0, SELECTED_NUMBER)

		rl.DrawTextEx(fonts[CONST.SELECTED_FONT], value.Line,
			rl.Vector2{X: position.X + 2*h_off,
				Y: position.Y + position.Y*float32(key)},
			float32(fonts[CONST.SELECTED_FONT].BaseSize), control.Spacing, rl.Color{181, 255, 120, 255})

	}

	// Draw line that separates line numbers and file text.
	rl.DrawLine(int32(h_off*float32(n_offset+1)), int32(2*y_off),
		int32(h_off*float32(n_offset+1)),
		int32(CONST.END_POINT_POSITION), CONST.LINE_COLOR)

}

// Aka draw cursor at the position on the file text.
func drawSelectedSquare(pos rl.Vector2,
	h_off float32, y_off float32, selected_color *rl.Color) {

	if !rl.IsWindowFocused() {
		*selected_color = rl.Color{selected_color.R, selected_color.G,
			selected_color.B, 128}
	}

	rl.DrawRectangle(int32(pos.X), int32(pos.Y),
		int32(h_off), int32(y_off),
		*selected_color)

}
