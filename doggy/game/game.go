package game

import "time"

type Game struct{}

// Current position of the dog.
func (g *Game) DogAt() [3]float64 {
	var zero [3]float64
	return zero
}

// Current position of the ball.
func (g *Game) BallAt() [3]float64 {
	var zero [3]float64
	return zero
}

// Integrate progresses the simulation by dt.
func (g *Game) Integrate(dt time.Duration) {
	g.integrateFloatSeconds(float64(dt / time.Second))
}

func (g *Game) integrateFloatSeconds(dt float64) {}

// WithBall is passed a callback to be called if the avatar has the ball.
func (g *Game) WithBall(f func(BallHandle)) {}

// A BallHandle represents the actions available to the avatar when they have the ball.
type BallHandle struct {
}

// Angle returns the angle at which the ball will be thrown if done so.
func (bh *BallHandle) Angle() float64 { return 0.0 }

// Throw makes the avatar throw the ball.
func (bh *BallHandle) Throw() {}

// FakeAThrow makes the avatar fake throwing the ball.
// The avatar makes a swing, but keeps the ball in it's hand.
func (bh *BallHandle) FakeAThrow() {}
