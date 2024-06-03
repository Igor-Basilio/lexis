package main

import (
	"bufio"
	"embed"
	"log"
	"os"

	"github.com/Igor-Basilio/lexis/async"
	CONST "github.com/Igor-Basilio/lexis/constant"
	"github.com/Igor-Basilio/lexis/control"
	"github.com/Igor-Basilio/lexis/text"
	"github.com/Igor-Basilio/lexis/ui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var screenWidth int32 = 800
var screenHeight int32 = 600
var camera rl.Camera2D = rl.NewCamera2D(rl.Vector2{X: 0, Y: 0},
	rl.Vector2{X: 0, Y: 0}, 0, 1.0)
var selected_color = rl.Yellow
var cursor_position = rl.Vector2{X: 0, Y: 0}
var White bool = false

//go:embed resources/fonts/*
var fontsFS embed.FS

func main() {

	args := os.Args[1:]

	if len(args) < 1 {
		log.Printf("Usage : Lexis <file_path> <...>")
		os.Exit(0)
	}
    
	CONST.CURRENT_FILE = "./" + args[0]

	file, err := os.Open(CONST.CURRENT_FILE)
	content := make(map[int]CONST.Data)
	var counter int = 1

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Breaks if line read is longer than 64K >
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content[counter] = CONST.Data{Line: scanner.Text(),
			Selected: false}
		counter++
	}

	// TODO: change so that we save the line
	// the cursor was at after saving the file
	// Start line at Cur_line = 1 should be selected
	// for new files
	n1, ok := content[1]
	if ok {
		n1.Selected = true
		content[1] = n1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	rl.SetTargetFPS(60)
	rl.InitWindow(screenWidth, screenHeight, "Lexis")

    rl.SetExitKey(rl.KeyDelete)
	control.Spacing = 2
	CONST.END_POINT_POSITION = rl.GetScreenHeight()

	// Loading default fonts that come with Lexis into the binary :
	space_mono, err := fontsFS.ReadFile(
		"resources/fonts/space-mono/SpaceMono-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	fira_code, err := fontsFS.ReadFile(
		"resources/fonts/fira-code/static/FiraCode-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	roboto_mono, err := fontsFS.ReadFile(
		"resources/fonts/roboto-mono/static/RobotoMono-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	var fonts [CONST.NUMBER_OF_FONTS]rl.Font
	fonts[0] = rl.LoadFontFromMemory(".ttf", space_mono,
		CONST.DEFAULT_FONT_SIZE, nil)
	fonts[1] = rl.LoadFontFromMemory(".ttf", fira_code,
		CONST.DEFAULT_FONT_SIZE, nil)
	fonts[2] = rl.LoadFontFromMemory(".ttf", roboto_mono,
		CONST.DEFAULT_FONT_SIZE, nil)
	// End of loading default fonts ;

	async.Async_manager(&selected_color)
	async.Debounce_Ticker.Reset(async.DEBOUNCE_TIMER)
	async.Flashing_Ticker.Reset(async.FLASHING_TIMER)

	for !rl.WindowShouldClose() {

		rl.BeginDrawing() // Draw

		rl.BeginMode2D(camera) // 2D

		/* 	h_off := fonts[CONST.SELECTED_FONT].Chars.Image.Width
		y_off := fonts[CONST.SELECTED_FONT].Chars.Image.Height */

		rl.ClearBackground(rl.Color{64, 64, 64, 255})

		ui.DrawMainUI(screenWidth, screenHeight, &camera, &fonts)

		text.DrawFileText(content, &selected_color,
			cursor_position, &fonts, &camera)

		control.Control_manager(&camera,
			&cursor_position, &selected_color,
			content, &White, &fonts)

        control.CheckCursorInCamera(&fonts, &camera, content)

		rl.EndMode2D() // 2D

		rl.EndDrawing() // Draw

	}

	rl.CloseWindow()
	async.Interrupt_Tickers <- true

}
