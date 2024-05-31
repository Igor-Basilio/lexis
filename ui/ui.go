package ui

import (
	CONST "github.com/Igor-Basilio/lexis/constant"
	"github.com/Igor-Basilio/lexis/control"
	h "github.com/Igor-Basilio/lexis/helper"
	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawMainUI(sw int32, sh int32, camera *rl.Camera2D,
	fonts *[CONST.NUMBER_OF_FONTS]rl.Font) {

	h_off := fonts[CONST.SELECTED_FONT].Chars.Image.Width
	y_off := fonts[CONST.SELECTED_FONT].Chars.Image.Height

	rl.DrawRectangle(0, 0, int32(float32(sw)/camera.Zoom)+
		int32(CONST.SCROLLED_COUNT_H)*h_off+int32(control.Spacing)*h_off,
		2*y_off, rl.Color{38, 37, 33, 255})

	// Draw lexis logo
	size := int32(len("Lexis"))
	rl.DrawText("Lexis", int32(float32((sw/2)-(size/2)*h_off)/camera.Zoom),
		y_off/2, fonts[CONST.SELECTED_FONT].BaseSize, rl.Color{181, 255, 120, 255})

	// Search Box
	h.DrawFunctionUntilCond(func() {

		rl.DrawRectangle(0, int32(rl.GetScreenHeight())-int32(1.5*float32(y_off)),
			25*h_off, int32(1.5*float32(y_off)), rl.Color{38, 37, 33, 255})

		rl.DrawTextPro(fonts[CONST.SELECTED_FONT], "S",
			rl.Vector2{X: float32(h_off / 2),
				Y: float32(rl.GetScreenHeight() - int(1.8*float32(y_off)))},
			rl.Vector2{}, 0, float32(2*fonts[CONST.SELECTED_FONT].BaseSize), control.Spacing,
			rl.Color{181, 255, 120, 255})

		/* if b, r := control.IsAnyKeyPressed(); b {

			CONST.SEARCH_BOX_TEXT += string(r)
			cursor.X += float32(h_off) + control.Spacing
			async.Flashing_Ticker.Reset(async.FLASHING_TIMER)
			*sc = rl.Yellow

		} */

	},
		CONST.DRAW_SEARCH_BOX)

}
