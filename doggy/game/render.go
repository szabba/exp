package game

import (
	"image/color"

	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dkit"
)

type State interface {
	BallAt() [3]float64
}

func Draw(ctx draw2d.GraphicContext, state State, width, height float64) {
	ctx.Save()
	defer ctx.Restore()

	r := 100.0
	cx, cy := width/2, height/2

	ctx.BeginPath()
	ctx.SetStrokeColor(color.Black)
	ctx.SetFillColor(color.White)
	ctx.SetLineWidth(1)

	draw2dkit.Circle(ctx, cx, cy, r)

	ctx.FillStroke()
}
