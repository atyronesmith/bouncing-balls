package effects

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/atyronesmith/bouncing-balls/pkg/physics"
)

// Lightning represents a lightning effect between two points
type Lightning struct {
	Lines     []*canvas.Line
	StartTime int64
	Duration  int64 // in milliseconds
}

// NewLightning creates a lightning bolt effect between two balls
func NewLightning(ball1, ball2 *physics.Ball) *Lightning {
	lightning := &Lightning{
		StartTime: time.Now().UnixMilli(),
		Duration:  300, // 300ms lightning effect
	}

	// Create multiple jagged lines for lightning effect
	numSegments := 8
	lightning.Lines = make([]*canvas.Line, numSegments)

	for i := 0; i < numSegments; i++ {
		// Calculate intermediate points with random jaggedness
		t := float32(i) / float32(numSegments-1)

		// Linear interpolation between ball centers
		x := ball1.X + (ball2.X-ball1.X)*t
		y := ball1.Y + (ball2.Y-ball1.Y)*t

		// Add random jaggedness
		if i > 0 && i < numSegments-1 {
			x += (float32(math.Sin(float64(i)*0.5)) * 20)
			y += (float32(math.Cos(float64(i)*0.7)) * 15)
		}

		line := &canvas.Line{
			StrokeColor: color.RGBA{R: 255, G: 255, B: 0, A: 255}, // Bright yellow
			StrokeWidth: 3.0,
		}

		if i == 0 {
			line.Position1 = fyne.NewPos(ball1.X, ball1.Y)
		} else {
			line.Position1 = lightning.Lines[i-1].Position2
		}

		if i == numSegments-1 {
			line.Position2 = fyne.NewPos(ball2.X, ball2.Y)
		} else {
			line.Position2 = fyne.NewPos(x, y)
		}

		lightning.Lines[i] = line
	}

	return lightning
}

// Update updates lightning animation and returns true if still active
func (l *Lightning) Update() bool {
	elapsed := time.Now().UnixMilli() - l.StartTime
	if elapsed > l.Duration {
		// Hide lightning lines
		for _, line := range l.Lines {
			if line != nil {
				line.Hide()
			}
		}
		return false // Lightning finished
	}

	// Animate lightning with flickering effect
	alpha := uint8(255 * (1.0 - float64(elapsed)/float64(l.Duration)))
	flicker := math.Sin(float64(elapsed)*0.05) > 0.5

	for _, line := range l.Lines {
		if line != nil {
			if flicker {
				line.StrokeColor = color.RGBA{R: 255, G: 255, B: 0, A: alpha}
			} else {
				line.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: alpha}
			}
			line.Refresh()
		}
	}

	return true // Lightning still active
}
