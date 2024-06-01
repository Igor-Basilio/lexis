package constant

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	ScrolledBottom        bool     = false
	ScrolledTop           bool     = false
	ScrolledRight         bool     = false
	ScrolledLeft          bool     = false
	END_POINT_POSITION    int      = 0
	DEBOUNCE_MOVER        bool     = false
	DEBOUNCE_MOVEL        bool     = false
	BACKGROUND_COLOR      rl.Color = rl.Color{64, 64, 64, 255}
	LINE_COLOR            rl.Color = rl.Color{255, 255, 255, 128}
	SELECTED_FONT         int      = 0
	SCROLLED_COUNT        int      = 0
	SCROLLED_COUNT_H      int      = 0
	DEFAULT_NUMBER_OFFSET int      = 5
	CURRENT_FILE          string   = ""
	DRAW_COND             bool     = false
	DRAW_TEXT             string   = ""
	DRAW_SEARCH_BOX       bool     = false
	SEARCH_BOX_TEXT       string   = ""
    NORMAL_MODE                    = true 
    CMD_RUNNING                    = false
    CMD_DIGIT                      = 0 
)

const (
	NUMBER_OF_FONTS   = 3
	DEFAULT_FONT_SIZE = int32(32)
	ROBOTO_MONO       = 2
	SPACE_MONO        = 0
	FIRA_CODE         = 1
	WRITE_GROUP_PERM  = 16  // XXX-W-XXX
	WRITE_OTHERS_PERM = 2   // XXXXXX-W-
	WRITE_OWNER_PERM  = 128 // -W-XXXXXX
)

type Data struct {
	Line     string
	Selected bool
}
