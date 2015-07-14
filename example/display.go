package main
// Display.go contains all the code to display to the screen the Mover's registered.
// Goroutines will be setup for two tasks, building the image, and displaying that image in the window.
// Use defined var's in this file to do basic setup, that can be changed later to allow config out of this file.
// Mover -> Goroutine 1 -> Creates Image -> On signal draws image and erases image for next frame

import (
	"fmt"
	"reflect"
	//"bufio"
	"log"
	"os"
	//"github.com/skelterjohn/geom"
	//"github.com/skelterjohn/go.uik/layouts"
	"github.com/llgcode/draw2d"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/Arrow/generals"
	"image"
	"image/color"
	//"image/png"
	"math"
	//"time"
)

const (
	heading          = 25
	border           = 5
)

var (
	// The path to the font used to draw text.
	fontPath = "resource/font/luximr.ttf"

	// The color of the text.
	fg = color.Black//xgraphics.BGRA{B: 0xff, G: 0xff, R: 0xff, A: 0xff}

	// The size of the text.
	size = 20.0
)

func SetupDisplay() (ctrlChan chan interface{}) {
	ctrlChan = make(chan interface{})
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	fontReader, err := os.Open(fontPath)
	if err != nil {
		log.Fatal(err)
	}

	// Now parse the font.
	font, err := xgraphics.ParseFont(fontReader)
	if err != nil {
		log.Fatal(err)
	}

	img, gc := initGc(width + border * 2, height + heading + border * 2)
	ximg := xgraphics.NewConvert(X, img)
	wid := ximg.XShowExtra("Regiment", true)

	circle := new(draw2d.PathStorage)

	ctr := 0
	fps := 0
	go func() {
		for {
			c := <-ctrlChan
			switch c := c.(type) {
			case *generals.Regiment:
				c.Step()
				rp := c.Position()
				gc.SetFillColor(image.Black)
				gc.SetStrokeColor(image.Black)
				
				for i := 0; i < len(rp.RowPath); i++ {
					x, y := rp.RowPath[i].X, rp.RowPath[i].Y
					circle = new(draw2d.PathStorage)
					circle.ArcTo(x + border, y + border + heading, 1, 1, 0, 2*math.Pi)
					gc.FillStroke(circle)
				}
				
				for i := 0; i < len(rp.ColPath); i++ {
					x, y := rp.ColPath[i].X, rp.ColPath[i].Y
					circle = new(draw2d.PathStorage)
					circle.ArcTo(x + border, y + border + heading, 1, 1, 0, 2*math.Pi)
					gc.FillStroke(circle)
				}
				
			case *Frame:
				for xp := border; xp < (width + border); xp++ {
					for yp := border; yp < (height + border + heading); yp++ {
						ximg.Set(xp, yp, img.At(xp, yp))
					}
				}
				_, _, err = ximg.Text(10, 0, fg, size, font, fmt.Sprintf("FPS: %v", fps))
				if err != nil {
					log.Fatal(err)
				}
				ximg.XDraw()
				ximg.XPaint(wid.Id)
				ctr++

				gc.SetFillColor(image.White)
				// fill the background 
				gc.Clear()
			case *Timing:
				fps = ctr
				_, _, err = ximg.Text(10, 0, fg, size, font, fmt.Sprintf("FPS: %v", fps))
				ctr = 0
			default:
				fmt.Println(reflect.TypeOf(c))
			}
		}
	}()

	go func() {
		xevent.Main(X)
	}()
	return ctrlChan
}

func initGc(w, h int) (image.Image, draw2d.GraphicContext) {
	i := image.NewRGBA(image.Rect(0, 0, w, h))
	gc := draw2d.NewGraphicContext(i)

	gc.SetStrokeColor(image.Black)
	gc.SetFillColor(image.White)
	// fill the background 
	gc.Clear()

	return i, gc
}
/*
func saveToPngFile(TestName string, m image.Image) {
	filePath := folder + TestName + ".png"
	f, err := os.Create(filePath)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer f.Close()
	b := bufio.NewWriter(f)
	err = png.Encode(b, m)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Wrote %s OK.\n", filePath)
}*/
