package main

import (
	. "github.com/flocks/concage/utils"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

const MAX_IMAGES = 10

func main() {
	const appID = "org.gtk.example"
	application, err := gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	files := readArgs()
	if err != nil {
		log.Fatal("Could not create application.", err)
	}

	application.Connect("activate", func() {
		onActivate(application, files)
	})
	application.Connect("shutdown", func() {})
	os.Exit(application.Run(os.Args))
}

func drawSquare(cr *cairo.Context) {
	cr.SetSourceRGB(1.0, 0, 0)
	cr.Rectangle(0, 100, 200, 200)
	cr.Fill()
}

func readArgs() []string {
	args := os.Args[1:]
	if len(args) > 0 && len(args) < 2 {
		log.Fatal("Usage: minimum of 2 files needs to be passed")
	}
	// TODO handle when files are passed as arguments

	return getDirectoryImages()
}

func getDirectoryImages() []string {
	// Get the current working directory
	files := []fs.DirEntry{}
	dirPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't read current directory")
	}
	// Open the directory
	dir, err := os.Open(dirPath)
	if err != nil {
		log.Fatal("Couldn't open directory")
	}
	defer dir.Close()

	// Read directory contents
	entries, err := dir.ReadDir(0)
	if err != nil {
		log.Fatal("Couldn't read directory")
	}

	// Iterate over directory entries
	for _, entry := range entries {
		if !entry.IsDir() {
			fileExtension := filepath.Ext(entry.Name())
			if fileExtension == ".png" {
				files = append(files, entry)
			}
		}
	}

	sort.Slice(files, func(i, j int) bool {
		desc1, _ := files[i].Info()
		desc2, _ := files[j].Info()

		return desc2.ModTime().Before(desc1.ModTime())
	})
	filesPath := []string{}
	for i, item := range files {
		if i == MAX_IMAGES {
			break
		}
		filesPath = append(filesPath, item.Name())
	}
	return filesPath
}

func onActivate(application *gtk.Application, files []string) {
	state := State{
		Margin:  0,
		BgColor: "red",
		Scale:   1,
		Images:  []Image{},
	}
	state, _ = loadImages(files, state)
	appWindow, err := gtk.ApplicationWindowNew(application)
	if err != nil {
		log.Fatal("Could not create application window.", err)
	}

	da, _ := gtk.DrawingAreaNew()
	da.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		drawSquare(cr)
		da.SetSizeRequest(300, 300)
	})

	appWindow.Connect("configure-event", func() {
		// fmt.Print("hey")
	})
	appWindow.Connect("key-press-event", func(win *gtk.ApplicationWindow, ev *gdk.Event) {
		keyEvent := &gdk.EventKey{ev}
		key := gdk.KeyValName(keyEvent.KeyVal())
		switch key {
		case "L", "J":
			state = ShiftNextImage(state)
			Render(appWindow, state)
		case "H", "K":
			state = ShiftPreviousImage(state)
			Render(appWindow, state)
		case "l", "j":
			state = ToggleNextImage(state)
			Render(appWindow, state)
		case "h", "k":
			state = TogglePreviousImage(state)
			Render(appWindow, state)
		case "d":
			state = RemoveCurrentImage(state)
			Render(appWindow, state)
		case "D":
			state = RemoveAllNextImage(state)
			Render(appWindow, state)
		case "v":
			state = ToggleDirection(state)
			Render(appWindow, state)
		case "Return":
			runConvert(state)
		case "q":
			os.Exit(1)
		}
	})
	Render(appWindow, state)
}

func runConvert(state State) {
	args := GenerateConvertCommand(state)
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Fatal("Error while running convert", err)
	}
	os.Exit(0)

}

func loadImages(filepaths []string, state State) (State, error) {
	for i := 0; i < len(filepaths); i++ {
		pixbuf, e := gdk.PixbufNewFromFile(filepaths[i])
		if e != nil {
			log.Fatal("error")
		}

		image := Image{PixelsBuffer: pixbuf, Selected: i == 0, ImageName: filepaths[i]}
		state.Images = append(state.Images, image)
	}

	return state, nil
}
