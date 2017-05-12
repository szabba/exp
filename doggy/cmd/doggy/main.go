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

	err := glfw.Init()
	if err != nil {
		return nil, err
	}

	var gw GameWindow

	gw.Window, err = glfw.CreateWindow(800, 600, "doggy", nil, nil)
	if err != nil {
		return nil, err
	}

	return &gw, nil
}

func (gw *GameWindow) CleanUp() {
	if gw != nil {
		// other cleanup
	}
	glfw.Terminate()
}

func (gw *GameWindow) Render() {}
