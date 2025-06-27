package ui

import (
	"image"
	"image/color"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/atyronesmith/bouncing-balls/pkg/physics"
)

// App represents the main application
type App struct {
	window          *app.Window
	theme           *material.Theme
	balls           []*physics.Ball
	human           *physics.Human
	dragon          *physics.Dragon
	starField       *physics.StarField
	alien           *physics.Alien
	currentBounds   image.Point
	animationTicker *time.Ticker

	// UI state
	resetButton widget.Clickable

	// Game state
	gameWidth  int
	gameHeight int
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{
		window:        new(app.Window),
		theme:         material.NewTheme(),
		currentBounds: image.Point{X: 800, Y: 600},
		gameWidth:     800,
		gameHeight:    600,
	}
}

// Run starts the application
func (a *App) Run() {
	// Initialize game objects (simplified for now)
	a.initializeGame()

	// Start animation
	a.startAnimation()

	// Configure window
	a.window.Option(app.Title("🚀 Eyeball Space Travel Simulator - Flying Through the Galaxy!"))
	a.window.Option(app.Size(unit.Dp(a.gameWidth), unit.Dp(a.gameHeight+50)))

	// Main event loop
	go func() {
		for {
			switch e := a.window.Event().(type) {
			case app.DestroyEvent:
				if a.animationTicker != nil {
					a.animationTicker.Stop()
				}
				return
			case app.FrameEvent:
				gtx := app.NewContext(&op.Ops{}, e)
				a.layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}

// initializeGame sets up all game objects (simplified)
func (a *App) initializeGame() {
	// Create three bouncing balls with slower, more controlled velocities
	ball1 := physics.NewCustomBall(
		100, 100, // position
		1.5, 1.2, // slower velocity
		30, // radius
		color.RGBA{R: 100, G: 150, B: 255, A: 255}, // Light blue fill
		color.RGBA{R: 255, G: 50, B: 50, A: 255},   // Red stroke
	)

	ball2 := physics.NewCustomBall(
		300, 200, // different starting position
		-1.2, 1.8, // slower velocity
		25, // smaller radius
		color.RGBA{R: 255, G: 100, B: 100, A: 255}, // Light red fill
		color.RGBA{R: 50, G: 255, B: 50, A: 255},   // Green stroke
	)

	ball3 := physics.NewCustomBall(
		500, 150, // different starting position
		-1.8, -1.4, // slower velocity
		35, // larger radius
		color.RGBA{R: 100, G: 255, B: 100, A: 255}, // Light green fill
		color.RGBA{R: 100, G: 50, B: 255, A: 255},  // Blue stroke
	)

	// Store all balls in a slice for easier management
	a.balls = []*physics.Ball{ball1, ball2, ball3}

	// Update bounds for all objects
	a.updateBounds()
}

// updateBounds updates the bounds for all physics objects
func (a *App) updateBounds() {
	gameArea := image.Point{X: a.gameWidth, Y: a.gameHeight}
	a.currentBounds = gameArea

	// Update bounds for all balls
	for _, ball := range a.balls {
		ball.Bounds = gameArea
	}
}

// startAnimation starts the animation loop
func (a *App) startAnimation() {
	// Start the animation automatically for all balls
	for _, ball := range a.balls {
		ball.IsAnimated = true
	}

	a.animationTicker = time.NewTicker(time.Millisecond * 16) // ~60 FPS
	go func() {
		defer a.animationTicker.Stop()
		for range a.animationTicker.C {
			// Update all ball positions (wall bouncing)
			for _, ball := range a.balls {
				ball.Update()
			}

			// Check for ball-to-ball collisions
			for i := 0; i < len(a.balls); i++ {
				for j := i + 1; j < len(a.balls); j++ {
					if a.balls[i].CheckCollision(a.balls[j]) {
						a.balls[i].HandleCollision(a.balls[j])
					}
				}
			}

			// Request window refresh
			a.window.Invalidate()
		}
	}()
}

// layout handles the main layout and rendering
func (a *App) layout(gtx layout.Context) layout.Dimensions {
	// Handle key events for human movement (TODO: implement when human is converted)
	// a.handleKeyEvents(gtx)

	return layout.Flex{
		Axis: layout.Vertical,
	}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return a.layoutControls(gtx)
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return a.layoutGame(gtx)
		}),
	)
}

// layoutControls renders the control buttons
func (a *App) layoutControls(gtx layout.Context) layout.Dimensions {
	if a.resetButton.Clicked(gtx) {
		a.resetAll()
	}

	return layout.UniformInset(unit.Dp(8)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return material.Button(a.theme, &a.resetButton, "Reset All").Layout(gtx)
	})
}

// layoutGame renders the main game area
func (a *App) layoutGame(gtx layout.Context) layout.Dimensions {
	// Set up clipping for the game area
	rect := image.Rectangle{Max: image.Point{X: a.gameWidth, Y: a.gameHeight}}
	defer clip.Rect(rect).Push(gtx.Ops).Pop()

	// Fill background with black
	paint.Fill(gtx.Ops, color.NRGBA{A: 255})

	// Render balls
	for _, ball := range a.balls {
		ball.Render(gtx.Ops)
	}

	return layout.Dimensions{Size: rect.Max}
}



// resetAll resets the entire game state
func (a *App) resetAll() {
	// Stop current animation
	if a.animationTicker != nil {
		a.animationTicker.Stop()
	}

	// Reinitialize everything
	a.initializeGame()

	// Restart animation
	a.startAnimation()
}
