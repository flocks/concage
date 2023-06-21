package utils

import (
	"github.com/gotk3/gotk3/gdk"
	"log"
)

type Image struct {
	PixelsBuffer *gdk.Pixbuf
	ImageName    string
	Selected     bool
}

type State struct {
	Margin    uint8
	BgColor   string
	Direction bool
	Scale     int // to define
	Images    []Image
}

func ShiftNextImage(state State) State {
	index := -1
	for i := 0; i < len(state.Images); i++ {
		if state.Images[i].Selected {
			index = i
			break
		}
	}
	if index >= 0 && index < len(state.Images)-1 {
		tmp := state.Images[index+1]
		state.Images[index+1] = state.Images[index]
		state.Images[index] = tmp
	}
	return state
}

func ShiftPreviousImage(state State) State {
	index := -1
	for i := 0; i < len(state.Images); i++ {
		if state.Images[i].Selected {
			index = i
			break
		}
	}
	if index > 0 {
		tmp := state.Images[index-1]
		state.Images[index-1] = state.Images[index]
		state.Images[index] = tmp
	}
	return state
}
func RemoveAllNextImage(state State) State {
	images := []Image{}
	indexCurrent := 0
	for i := 0; i < len(state.Images); i++ {
		if state.Images[i].Selected == false {
			images = append(images, state.Images[i])
		} else {
			indexCurrent = i
			break
		}
	}
	if indexCurrent < 2 {
		return state
	}

	if indexCurrent < len(images) {
		images[indexCurrent].Selected = true
	} else {
		images[indexCurrent-1].Selected = true
	}
	state.Images = images
	return state

}

func RemoveCurrentImage(state State) State {
	images := []Image{}
	indexCurrent := 0
	for i := 0; i < len(state.Images); i++ {
		if state.Images[i].Selected == false {
			images = append(images, state.Images[i])
		} else {
			indexCurrent = i
		}
	}
	if len(images) == 0 {
		log.Fatal("No more images")
	}
	if indexCurrent < len(images) {
		images[indexCurrent].Selected = true
	} else {
		images[indexCurrent-1].Selected = true
	}
	state.Images = images
	return state
}

func ToggleNextImage(state State) State {
	for i := 0; i < len(state.Images); i++ {
		if state.Images[i].Selected {
			if (i + 1) < len(state.Images) {
				state.Images[i].Selected = false
				state.Images[i+1].Selected = true
			}
			break
		}
	}

	return state
}

func TogglePreviousImage(state State) State {
	prevIndex := -1
	for i := 0; i < len(state.Images); i++ {
		if state.Images[i].Selected {
			state.Images[i].Selected = false
			prevIndex = i
			break
		}
	}
	if prevIndex > 0 {
		prevIndex = prevIndex - 1
	}
	state.Images[prevIndex].Selected = true

	return state
}

func ToggleDirection(state State) State {
	state.Direction = !state.Direction
	return state
}

func GenerateConvertCommand(state State) []string {
	commands := []string{"convert"}
	for _, image := range state.Images {
		commands = append(commands, image.ImageName)
	}
	if state.Direction {
		commands = append(commands, "-append")
	} else {
		commands = append(commands, "+append")
	}
	commands = append(commands, "output.png")

	return commands
}
