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

type soundboard struct {
	otoContext *oto.Context
	content    map[string][]byte
	columns    int
}

var (
	//go:embed sb.ico
	icon []byte
)

func (sb *soundboard) play(input []byte) {
	ap := sb.otoContext.NewPlayer()
	if _, err := ap.Write(input); err != nil {
		log.Panicf("failed writing to player: %v", err)
	}
	if err := ap.Close(); err != nil {
		log.Panicf("failed to close player: %v", err)
	}
}

func (sb *soundboard) gui() {
	a := app.New()
	w := a.NewWindow("go-soundboard")
	w.SetIcon(fyne.NewStaticResource("icon", icon))
	ng := container.NewGridWithColumns(sb.columns)
	keys := make([]string, 0, len(sb.content))
	for k := range sb.content {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		sl := key
		nb := sb.content[key]
		butt := widget.NewButton(sl, func() {
			go sb.play(nb)
		})
		ng.Add(butt)
	}
	w.SetContent(ng)
	w.ShowAndRun()
}

func main() {
	sb := soundboard{}
	var otoErr error
	sb.otoContext, otoErr = oto.NewContext(48000, 2, 2, 8192)
	if otoErr != nil {
		log.Panicf("Error creating oto.NewContext %v", otoErr)
	}
	var path string
	flag.IntVar(&sb.columns, "c", 3, "columns")
	flag.StringVar(&path, "d", "audio/", "dir of 48khz mp3 files")
	flag.Parse()
	if path[len(path)-1:] != "/" {
		path = path + "/"
	}
	dir, err := os.ReadDir(path)
	if err != nil {
		return
	}
	sb.content = map[string][]byte{}
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
		if sb.content[strings.Replace(en, ".mp3", "", 1)], err = io.ReadAll(decodedFile); err != nil {
			log.Panicf("Error reading decodedFile %s: %v", en, err)
		}
	}
	sb.gui()
}
