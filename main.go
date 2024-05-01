package main

import (
	_ "embed"
	"flag"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var (
	//go:embed sb.ico
	icon       []byte
	otoContext *oto.Context
	content    map[string][]byte
	columns    int
)

func play(input []byte) {
	ap := otoContext.NewPlayer()
	if _, err := ap.Write(input); err != nil {
		log.Panicf("failed writing to player: %v", err)
	}
	if err := ap.Close(); err != nil {
		log.Panicf("failed to close player: %v", err)
	}
}

func gui() {
	a := app.New()
	w := a.NewWindow("go-soundboard")
	w.SetIcon(fyne.NewStaticResource("icon", icon))
	ng := container.NewGridWithColumns(columns)
	keys := make([]string, 0, len(content))
	for k := range content {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		sl := key
		nb := content[key]
		butt := widget.NewButton(sl, func() {
			go play(nb)
		})
		ng.Add(butt)
	}
	w.SetContent(ng)
	w.ShowAndRun()
}

func init() {
	var otoErr error
	otoContext, otoErr = oto.NewContext(48000, 2, 2, 8192)
	if otoErr != nil {
		log.Panicf("Error creating oto.NewContext %v", otoErr)
	}
}

func main() {
	var path string
	flag.IntVar(&columns, "c", 3, "columns")
	flag.StringVar(&path, "d", "audio/", "dir of 48khz mp3 files")
	flag.Parse()
	if path[len(path)-1:] != "/" {
		path = path + "/"
	}
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}
	content = map[string][]byte{}
	for _, entry := range dir {
		en := entry.Name()
		if !strings.Contains(en, ".mp3") {
			continue
		}
		file, err := os.Open(path + en)
		if err != nil {
			log.Panicf("Error opening file %s: %v", en, err)
		}
		decodedFile, err := mp3.NewDecoder(file)
		if err != nil {
			log.Panicf("Error decoding file %s: %v", en, err)
		}
		if content[strings.Replace(en, ".mp3", "", 1)], err = io.ReadAll(decodedFile); err != nil {
			log.Panicf("Error reading decodedFile %s: %v", en, err)
		}
	}
	gui()
}
