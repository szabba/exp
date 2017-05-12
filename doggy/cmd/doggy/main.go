package main

import (
	"log"
	"time"

	"runtime"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	runtime.LockOSThread()

	window, err := New()
	defer window.CleanUp()
	if err != nil {
		log.Fatal(err)
	}

	for !window.ShouldClose() {
		window.Render()
		glfw.PollEvents()
		time.Sleep(time.Second / 60)
	}
}

type GameWindow struct {
	*glfw.Window
}

func New() (*GameWindow, error) {

	const (
		width, height = 800, 600
		title         = "doggy"
	)

	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	var gw GameWindow

	gw.Window, err = glfw.CreateWindow(width, height, title, nil, nil)
	if err != nil {
		return nil, err
	}

	err = gl.Init()
	if err != nil {
		return nil, err
	}

	gl.ClearColor(1, 1, 1, 1)

	gw.MakeContextCurrent()
	gw.SetSizeCallback(gw.OnResize)
	gw.OnResize(gw.Window, 800, 600)

	return &gw, nil
}

func (gw *GameWindow) CleanUp() {
	if gw != nil {
		// other cleanup
	}
	glfw.Terminate()
}

func (gw *GameWindow) OnResize(_ *glfw.Window, w, h int) {
	gl.Viewport(0, 0, int32(w), int32(h))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), 0, float64(h), -1, 1)
	gl.Scalef(1, -1, 1)
	gl.Translatef(0, float32(-h), 0)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
}

func (gw *GameWindow) Render() {
	defer gw.SwapBuffers()
	defer gl.Flush()

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}
