package searchwindow

import (
	"fmt"
	"image"
	"os"
	"time"

	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/jezek/xgb/xproto"
	"github.com/jezek/xgbutil"
	"github.com/jezek/xgbutil/keybind"
	"github.com/jezek/xgbutil/xevent"
	"github.com/jezek/xgbutil/xgraphics"
	"github.com/jezek/xgbutil/xwindow"
)

type Window struct {
	FontPath string
	Width    int
	Height   int
	FPS      int
	BG       xgraphics.BGRA
	Title    string

	DrawFunc     func(ximg *xgraphics.Image, font *truetype.Font)
	KeypressFunc func(key string)

	x    *xgbutil.XUtil
	font *truetype.Font
	ximg *xgraphics.Image
	wnd  *xwindow.Window
}

func New() *Window {
	return &Window{
		FontPath: "/usr/share/fonts/TTF/Fira Code Light Nerd Font Complete.ttf",
		Width:    600,
		Height:   310,
		FPS:      30,
		BG:       xgraphics.BGRA{R: 0x2E, G: 0x2E, B: 0x2E, A: 0xff},
		Title:    "Search",
	}
}

func (w *Window) Show() error {
	x, err := xgbutil.NewConn()
	if err != nil {
		return fmt.Errorf("Failed to connect to X: %w", err)
	}
	w.x = x
	err = w.loadFont()
	if err != nil {
		return fmt.Errorf("Failed to load font: %w", err)
	}

	ximg := xgraphics.New(w.x, image.Rect(0, 0, w.Width, w.Height))
	w.ximg = ximg

	wnd := ximg.XShowExtra(w.Title, true)
	w.wnd = wnd

	keybind.Initialize(w.x)
	err = w.wnd.Listen(xproto.EventMaskKeyPress, xproto.EventMaskKeyRelease)
	if err != nil {
		return fmt.Errorf("Failed to add key-event-listener: %w", err)
	}
	w.wnd.Map()
	xevent.KeyPressFun(func(xu *xgbutil.XUtil, event xevent.KeyPressEvent) {
		keystr := keybind.LookupString(xu, event.State, event.Detail)
		switch keystr {
		case "odiaeresis":
			keystr = "ö"
		case "Odiaeresis":
			keystr = "ö"
		case "adiaeresis":
			keystr = "ä"
		case "Adiaeresis":
			keystr = "Ä"
		case "udiaeresis":
			keystr = "ü"
		case "Udiaeresis":
			keystr = "Ü"
		case "ssharp":
			keystr = "ß"
		}
		if w.KeypressFunc != nil {
			fmt.Println(keystr)
			w.KeypressFunc(keystr)
		}
	}).Connect(w.x, w.wnd.Id)

	go func() {
		for {
			w.draw()
			time.Sleep(time.Second / time.Duration(w.FPS))
		}
	}()

	xevent.Main(w.x)
	return nil
}

func (w *Window) draw() {
	w.ximg.For(func(x, y int) xgraphics.BGRA {
		return w.BG
	})
	if w.DrawFunc != nil {
		w.DrawFunc(w.ximg, w.font)
	}
	w.ximg.XDraw()
	w.ximg.XPaint(w.wnd.Id)
}

func (w *Window) loadFont() error {
	r, err := os.Open(w.FontPath)
	if err != nil {
		return fmt.Errorf("Failed to open font '%s': %w", w.FontPath, err)
	}

	f, err := xgraphics.ParseFont(r)
	if err != nil {
		return fmt.Errorf("Failed to parse font: %w", err)
	}
	w.font = f

	return nil
}
