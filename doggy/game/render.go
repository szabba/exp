package game

import (
	"image/color"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
)

type State interface {
	BallAt() [3]float64
}

func Draw(ctx draw2d.GraphicContext, state State) {
	ctx.SetFontData(draw2d.FontData{
		Name:   "luxi",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleBold | draw2d.FontStyleItalic})

	ctx.BeginPath()
	draw2dkit.RoundedRectangle(ctx, 200, 200, 600, 600, 100, 100)

	ctx.SetFillColor(color.Black)
	ctx.Fill()
}
