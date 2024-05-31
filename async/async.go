package async

import (
	"time"

	CONST "github.com/Igor-Basilio/lexis/constant"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	FLASHING_TIMER = 500 * time.Millisecond
	DEBOUNCE_TIMER = 100 * time.Millisecond
)

var (
	Debounce_Ticker   *time.Ticker
	Flashing_Ticker   *time.Ticker
	Event_Ticker      *time.Ticker
	Interrupt_Tickers chan bool     = make(chan bool)
	EVENT_TIMER       time.Duration = 1 * time.Second
)

func Async_manager(sc *rl.Color) {

	Debounce_Ticker = time.NewTicker(DEBOUNCE_TIMER)
	Flashing_Ticker = time.NewTicker(FLASHING_TIMER)
	Event_Ticker = time.NewTicker(EVENT_TIMER)

	defer Flashing_Ticker.Stop()
	defer Debounce_Ticker.Stop()

	go func() {

		for {
			select {
			case <-Interrupt_Tickers:
				return
			case <-Debounce_Ticker.C:
				CONST.DEBOUNCE_MOVEL = false
				CONST.DEBOUNCE_MOVER = false
			}
		}

	}()

	go func() {

		for {
			select {
			case <-Interrupt_Tickers:
				return
			case <-Flashing_Ticker.C:
				if *sc == rl.Yellow {
					*sc = CONST.BACKGROUND_COLOR
				} else {
					*sc = rl.Yellow
				}
			}
		}

	}()

	go func() {

		for {
			select {
			case <-Interrupt_Tickers:
				return
			case <-Event_Ticker.C:
				CONST.DRAW_COND = false
				Event_Ticker.Stop()
			}
		}

	}()

}
