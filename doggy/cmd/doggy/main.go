package main

import (
	"log"

	"github.com/go-gl/glfw/v3.2/glfw"
)

func main() {
	window, err := New()
	if err != nil {
		log.Fatal(err)
	}

	for !window.ShouldClose() {
		window.Render()
		glfw.PollEvents()
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

func (gw *GameWindow) OnResize(_ *glfw.Window, newW, newH int) {
	log.Printf("new size is %dx%d", newW, newH)
}

func (gw *GameWindow) Render() {}
