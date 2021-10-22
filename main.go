package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"os"
	"os/exec"
	"strings"

	"github.com/BurntSushi/freetype-go/freetype/truetype"
	"github.com/Neokil/lastpass-search/lastpasshelper"
	"github.com/Neokil/lastpass-search/searchwindow"
	"github.com/jezek/xgbutil/xgraphics"
)

type Config struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	OTP      bool   `json:"OTP"`
}

func main() {
	if len(os.Args) != 2 {
		panic("Exactly 1 Argument required")
	}

	switch os.Args[1] {
	case "update":
		fallthrough
	case "u":
		mainUpdate()
	case "search":
		fallthrough
	case "s":
		mainSearch()
	default:
		panic(fmt.Sprintf("Unknown action '%s', allowed actions are 'update' and 'search'", os.Args[1]))
	}
}

func mainUpdate() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	otpCode := ""
	if config.OTP {
		fmt.Print("Please enter your OTP-Code for LastPass: ")
		r := bufio.NewReader(os.Stdin)
		t, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		otpCode = strings.TrimSpace(t)
	}
	err = lastpasshelper.UpdateAccounts(config.Username, config.Password, otpCode)
	if err != nil {
		panic(err)
	}
	fmt.Println("Update done")
}

func mainSearch() {
	accounts, err := lastpasshelper.GetAccounts()
	if errors.Is(err, os.ErrNotExist) {
		panic("Cannot find Accounts-Cache-File. Please create one using 'lastpass-search update'.")
	}

	w := searchwindow.New()

	fg := xgraphics.BGRA{R: 0xdd, G: 0xdd, B: 0xdd, A: 0xff}
	bgFocus := xgraphics.BGRA{R: 0x4D, G: 0x4D, B: 0x4D, A: 0xff}
	fontSize := 14.0
	currentSearch := ""
	currentSelection := 0
	visibleAccountCount := 10
	visibleAccounts := getMatchingNAccounts(accounts, currentSearch, visibleAccountCount)

	w.DrawFunc = func(ximg *xgraphics.Image, font *truetype.Font) {
		_, y, _ := ximg.Text(10, 10, fg, fontSize, font, "Search: "+currentSearch)
		ximg.SubImage(image.Rect(10, y+5, w.Width-10, y+6)).(*xgraphics.Image).For(func(x, y int) xgraphics.BGRA {
			return fg
		})

		for i, a := range visibleAccounts {
			t := fmt.Sprintf("%s (%s at %s)", a.Name, a.Username, a.URL)
			if i == currentSelection {
				t = "> " + t
				ximg.SubImage(image.Rect(5, y+int(fontSize), w.Width-5, y+int(2.4*fontSize))).(*xgraphics.Image).For(func(x, y int) xgraphics.BGRA {
					return bgFocus
				})
			}
			_, y, _ = ximg.Text(10, y+int(fontSize), fg, float64(fontSize), font, t)
		}
	}
	w.KeypressFunc = func(key string) {
		fmt.Printf("key: %s\n", key)
		switch key {
		case "BackSpace":
			if len(currentSearch) > 0 {
				currentSearch = currentSearch[0 : len(currentSearch)-1]
			}
			visibleAccounts = getMatchingNAccounts(accounts, currentSearch, visibleAccountCount)
			currentSelection = 0
			return
		case "Down":
			if currentSelection >= visibleAccountCount-1 {
				currentSelection = 0
			} else {
				currentSelection++
			}
		case "Up":
			if currentSelection == 0 {
				currentSelection = visibleAccountCount - 1
			} else {
				currentSelection--
			}
		case "Return":
			err := copyToClipboard(visibleAccounts[currentSelection].Password)
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		case "Escape":
			os.Exit(0)
		}

		if len(key) == 1 || isSpecial(key) {
			currentSearch += key
			visibleAccounts = getMatchingNAccounts(accounts, currentSearch, visibleAccountCount)
			currentSelection = 0
		}
	}
	err = w.Show()
	if err != nil {
		panic(err)
	}
}

func isSpecial(s string) bool {
	if len(s) > 2 {
		return false
	}
	s = strings.ToLower(s)
	return s == "ä" || s == "ö" || s == "ü" || s == "ß"
}

func getMatchingNAccounts(accounts []lastpasshelper.Account, search string, count int) (accs []lastpasshelper.Account) {
	for _, a := range accounts {
		if accountIsMatch(a, search) {
			accs = append(accs, a)
			if len(accs) >= 10 {
				return accs
			}
		}
	}
	return accs
}

func accountIsMatch(a lastpasshelper.Account, s string) bool {
	s = strings.ToLower(s)
	if strings.Contains(strings.ToLower(a.Name), s) {
		return true
	}
	if strings.Contains(strings.ToLower(a.Username), s) {
		return true
	}
	if strings.Contains(strings.ToLower(a.URL), s) {
		return true
	}
	if strings.Contains(strings.ToLower(a.Notes), s) {
		return true
	}
	return false
}

func loadConfig() (*Config, error) {
	b, err := os.ReadFile("auth.json")
	if err != nil {
		return nil, fmt.Errorf("Failed to read auth.json: %w", err)
	}
	c := &Config{}
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal auth.json: %w", err)
	}
	return c, nil
}

func copyToClipboard(s string) error {
	cmd := exec.Command("xclip", "-in", "-selection", "clipboard")
	in, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("Failed to create StdinPipe: %w", err)
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("Failed to start command: %w", err)
	}
	_, err = in.Write([]byte(s))
	if err != nil {
		return fmt.Errorf("Failed to write to StdinPipe: %w", err)
	}
	err = in.Close()
	if err != nil {
		return fmt.Errorf("Failed to close StdinPipe: %w", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Failed to wait for command to end: %w", err)
	}
	return nil
}
