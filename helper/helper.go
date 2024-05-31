package helper

import (
	"time"

	"github.com/Igor-Basilio/lexis/async"
	CONST "github.com/Igor-Basilio/lexis/constant"
	rl "github.com/gen2brain/raylib-go/raylib"
)

// Helper function that returns the amount
// of digits in the typed manner.
// Ex : 43 -> 2, 1 -> 1, 453 -> 3, ...
// Breaks if number >= 1 Million
func GetAmountOfDigits(number int) int {

	if number >= 0 && number <= 9 {
		return 1
	} else if number >= 10 && number <= 99 {
		return 2
	} else if number >= 100 && number <= 999 {
		return 3
	} else if number >= 1000 && number <= 9999 {
		return 4
	}

	return CONST.DEFAULT_NUMBER_OFFSET

}

// Helper function to reduce size of checking
// for errors on function returns.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Function that draws text on the screen for a specified duration given a position
// a condition, if cond = true then it'll always be drawn.
func DrawTextForSpecifiedTime(text string, dur time.Duration, pos rl.Vector2,
	cond *bool, f rl.Font, color rl.Color, spacing float32) {

	if *cond {
		async.EVENT_TIMER = dur
		async.Event_Ticker.Reset(async.EVENT_TIMER)
		CONST.DRAW_COND = true
		CONST.DRAW_TEXT = text
		*cond = false
	}

	if CONST.DRAW_COND {
		rl.DrawTextEx(f, CONST.DRAW_TEXT,
			pos, float32(f.BaseSize), spacing, color)
	}

}

// Runs the function passed as an argument until the condition is true
// if the condition is true nothing is done.
func DrawFunctionUntilCond(f func(), cond bool) {

	if cond {

		f()

	}

}

// Gets an integer and returns the size of "n" characters on the
// current font
/* func GetChars( n int32 ) int32 {

	h_off :=

} */

