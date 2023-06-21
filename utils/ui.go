package utils

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func drawHeader() (*gtk.Box, error) {
	header, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	title, _ := gtk.LabelNew("")
	title.SetMarkup("<b>Concage</b>")
	desc, _ := gtk.LabelNew("Frontend application to imagemagick append command")

	header.Add(title)
	header.Add(desc)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)
	spacingBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 20)

	usage, _ := gtk.LabelNew("")
	usage.SetMarkup("<span><i>h,j,k,l</i>to move, <i>shift+h,j,k,l</i> to shift image, <i>d,D</i> delete, <i>v</i> change direction</span>")

	box.Add(header)
	box.Add(usage)
	box.Add(spacingBox)
	return box, nil
}

func Render(window *gtk.ApplicationWindow, state State) {
	cleanWindow(window)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 40)
	box.SetMarginTop(30)
	header, _ := drawHeader()

	// TODO we need to handle auto scroll after image selection
	scrolledWindow, _ := gtk.ScrolledWindowNew(nil, nil)
	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_NEVER)

	box.Add(header)
	stateView, _ := drawState(state)
	stateView.SetHAlign(gtk.ALIGN_CENTER)
	stateView.SetVAlign(gtk.ALIGN_CENTER)
	scrolledWindow.Add(stateView)
	box.Add(scrolledWindow)
	window.Add(box)
	window.ShowAll()
}

func drawState(state State) (*gtk.Box, error) {
	direction := gtk.ORIENTATION_HORIZONTAL
	if state.Direction {
		direction = gtk.ORIENTATION_VERTICAL
	}
	box, _ := gtk.BoxNew(direction, int(state.Margin))
	for i := 0; i < len(state.Images); i++ {
		boxImage, _ := drawImageWithBorder(state.Images[i])
		box.Add(boxImage)
	}
	return box, nil
}

func cleanWindow(window *gtk.ApplicationWindow) {
	children := window.GetChildren()
	children.Foreach(
		func(item interface{}) {
			if widget, ok := item.(*gtk.Widget); ok {
				window.Remove(widget)
			}
		})
}

func resizeImage(image *gdk.Pixbuf) (*gdk.Pixbuf, error) {
	originalWidth := image.GetWidth()
	originalHeight := image.GetHeight()
	aspectRatio := float64(originalWidth) / float64(originalHeight)

	desiredWidth := 300
	desiredHeight := 300
	var newWidth, newHeight int
	if float64(desiredWidth)/float64(desiredHeight) > aspectRatio {
		newWidth = desiredHeight * originalWidth / originalHeight
		newHeight = desiredHeight
	} else {
		newWidth = desiredWidth
		newHeight = desiredWidth * originalHeight / originalWidth
	}

	// Resize the original Pixbuf while maintaining aspect ratio
	resizedPixbuf, err := image.ScaleSimple(newWidth, newHeight, 1)
	return resizedPixbuf, err
}

func drawImageWithBorder(image Image) (*gtk.Box, error) {
	resized, _ := resizeImage(image.PixelsBuffer)
	width := resized.GetWidth()
	height := resized.GetHeight()

	borderSize := 5
	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	overlay, _ := gtk.DrawingAreaNew()
	overlay.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		if image.Selected {
			cr.SetSourceRGB(1.0, 0, 0)
		} else {
			cr.SetSourceRGB(0.878, 0.878, 0.878)
		}
		cr.Rectangle(0, 0, float64(width+borderSize), float64(height+borderSize))
		cr.Fill()

		cr.Rectangle(0, 0, float64(width), float64(height))

		surface, _ := gdk.CairoSurfaceCreateFromPixbuf(resized, 1, nil)
		cr.SetSourceSurface(surface, float64(borderSize/2), float64(borderSize/2))
		cr.Paint()
	})
	overlay.SetSizeRequest(width+borderSize, height+borderSize)
	box.Add(overlay)

	return box, nil
}
