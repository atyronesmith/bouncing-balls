package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/atyronesmith/bouncing-balls/pkg/effects"
	"github.com/atyronesmith/bouncing-balls/pkg/physics"
)

// App represents the main application
type App struct {
	fyneApp fyne.App
	window  fyne.Window
	balls   []*physics.Ball
	human   *physics.Human
	dragon  *physics.Dragon
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{
		fyneApp: app.New(),
	}
}

// Run starts the application
func (a *App) Run() {
	a.fyneApp.SetIcon(nil)

	// Create a larger window for bouncing animation
	a.window = a.fyneApp.NewWindow("Dragon Chases Balls + Human Keyboard Control - macOS")
	a.window.Resize(fyne.NewSize(800, 600))
	a.window.CenterOnScreen()

	// Create three bouncing balls with different properties
	ball1 := physics.NewCustomBall(
		100, 100, // position
		3.5, 2.8, // velocity
		30, // radius
		color.RGBA{R: 100, G: 150, B: 255, A: 255}, // Light blue fill
		color.RGBA{R: 255, G: 50, B: 50, A: 255},   // Red stroke
	)

	ball2 := physics.NewCustomBall(
		300, 200, // different starting position
		-2.8, 4.2, // different velocity (negative x for opposite direction)
		25, // smaller radius
		color.RGBA{R: 255, G: 100, B: 100, A: 255}, // Light red fill
		color.RGBA{R: 50, G: 255, B: 50, A: 255},   // Green stroke
	)

	ball3 := physics.NewCustomBall(
		500, 150, // different starting position
		-4.1, -3.3, // different velocity (both negative)
		35, // larger radius
		color.RGBA{R: 100, G: 255, B: 100, A: 255}, // Light green fill
		color.RGBA{R: 100, G: 50, B: 255, A: 255},  // Blue stroke
	)

	// Store all balls in a slice for easier management
	a.balls = []*physics.Ball{ball1, ball2, ball3}

	// Create the human figure (same size as largest ball - ball3 has radius 35)
	a.human = physics.NewHuman(400, 300, 35) // Center of screen, size of largest ball

	// Create the dragon (slightly larger than human, positioned differently)
	a.dragon = physics.NewDragon(200, 200, 40) // Upper left area, slightly larger than human

	// Lightning system
	var activeLightning []*effects.Lightning
	frameCount := 0

	// Create instruction label
	label := widget.NewLabel("⚡🌟🐉⌨️ Lightning + Trails + Dragon Chasing + Human Keyboard Control (Arrow Keys)!")
	label.Alignment = fyne.TextAlignCenter

	// Create UI controls
	controls := a.createControls()

	// Create the main container with trails, circles, and human
	content := container.NewWithoutLayout()

	// Add ball trails to container
	for _, ball := range a.balls {
		for _, trail := range ball.Trail {
			content.Add(trail)
		}
	}

	// Add balls
	content.Add(ball1.Circle) // Blue circle
	content.Add(ball2.Circle) // Red circle
	content.Add(ball3.Circle) // Green circle

	// Add ball text labels
	content.Add(ball1.Text) // Blue ball text
	content.Add(ball2.Text) // Red ball text
	content.Add(ball3.Text) // Green ball text

	// Add human figure components
	content.Add(a.human.Head)
	content.Add(a.human.Body)
	content.Add(a.human.LeftArm)
	content.Add(a.human.RightArm)
	content.Add(a.human.LeftLeg)
	content.Add(a.human.RightLeg)

	// Add dragon figure components
	dragonComponents := a.dragon.GetVisualComponents()
	for _, component := range dragonComponents {
		content.Add(component)
	}

	// Create the full layout
	fullContent := container.NewVBox(
		label,
		content,
		controls,
	)

	// Set the content
	a.window.SetContent(fullContent)

		// Set up keyboard event handling for continuous key presses
	// We'll use a map to track key press timestamps for auto-release
	keyTimestamps := make(map[fyne.KeyName]int)
	keyTimeout := 3 // frames before auto-release (at 60fps = ~50ms)

	a.window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyUp, fyne.KeyDown, fyne.KeyLeft, fyne.KeyRight:
			keyTimestamps[key.Name] = frameCount
		}
	})

	// Update human key states based on recent key presses
	updateKeyStates := func() {
		a.human.KeyUp = (frameCount - keyTimestamps[fyne.KeyUp]) < keyTimeout
		a.human.KeyDown = (frameCount - keyTimestamps[fyne.KeyDown]) < keyTimeout
		a.human.KeyLeft = (frameCount - keyTimestamps[fyne.KeyLeft]) < keyTimeout
		a.human.KeyRight = (frameCount - keyTimestamps[fyne.KeyRight]) < keyTimeout
	}

	// Start the animation automatically for all balls
	for _, ball := range a.balls {
		ball.IsAnimated = true
	}

	// Animation timer - updates all balls and human 60 times per second for smooth animation
	go func() {
		ticker := time.NewTicker(time.Millisecond * 16) // ~60 FPS
		defer ticker.Stop()
		for range ticker.C {
			frameCount++

			// Update keyboard states
			updateKeyStates()

			// Update all ball positions (wall bouncing)
			for _, ball := range a.balls {
				ball.Update()
			}

			// Check for ball-to-ball collisions and create lightning
			for i := 0; i < len(a.balls); i++ {
				for j := i + 1; j < len(a.balls); j++ {
					if a.balls[i].CheckCollision(a.balls[j]) {
						a.balls[i].HandleCollision(a.balls[j])

						// Create lightning effect on collision
						lightning := effects.NewLightning(a.balls[i], a.balls[j])
						activeLightning = append(activeLightning, lightning)

						// Add lightning lines to visual container
						for _, line := range lightning.Lines {
							content.Add(line)
						}
					}
				}
			}

			// Update lightning effects
			var stillActiveLightning []*effects.Lightning
			for _, lightning := range activeLightning {
				if lightning.Update() {
					stillActiveLightning = append(stillActiveLightning, lightning)
				}
			}
			activeLightning = stillActiveLightning

			// Update human behavior
			if a.human.IsExploding {
				a.human.UpdateExplosion()
			} else {
				// Update human behavior (includes keyboard control and AI avoidance)
				a.human.Update(a.balls)

				// Check for collisions with balls
				if a.human.CheckCollisionWithBalls(a.balls) {
					a.human.Explode()

					// Add explosion particles to the visual container
					for _, particle := range a.human.ExplosionParticles {
						if particle != nil {
							content.Add(particle)
						}
					}
				}
			}

			// Update dragon behavior (handles collision, drifting, spinning, and chasing)
			a.dragon.Update(a.balls)
			a.dragon.UpdatePosition()
		}
	}()

	// Show and run the application
	a.window.ShowAndRun()
}

// createControls creates the UI control buttons
func (a *App) createControls() *fyne.Container {
	// Create animation control buttons
	startButton := widget.NewButton("▶️ Start All", func() {
		for _, ball := range a.balls {
			ball.IsAnimated = true
		}
	})

	stopButton := widget.NewButton("⏸️ Stop All", func() {
		for _, ball := range a.balls {
			ball.IsAnimated = false
		}
	})

	colorButton := widget.NewButton("🎨 Change Colors", func() {
		for _, ball := range a.balls {
			ball.ChangeColor()
		}
	})

	speedUpButton := widget.NewButton("⚡ Speed Up", func() {
		for _, ball := range a.balls {
			ball.VX *= 1.2
			ball.VY *= 1.2
		}
	})

	slowDownButton := widget.NewButton("🐌 Slow Down", func() {
		for _, ball := range a.balls {
			ball.VX *= 0.8
			ball.VY *= 0.8
		}
	})

	massInfoButton := widget.NewButton("⚖️ Show Masses", func() {
		// Display mass information (for demonstration)
		// In a real app, you might use a dialog or status bar
		println("Ball Masses:")
		println("Blue ball (radius 30):", a.balls[0].GetMass())
		println("Red ball (radius 25):", a.balls[1].GetMass())
		println("Green ball (radius 35):", a.balls[2].GetMass())
	})

	humanButton := widget.NewButton("🏃 Toggle Human", func() {
		if a.human.IsExploding {
			return // Don't allow toggle during explosion
		}
		a.human.IsActive = !a.human.IsActive
		if !a.human.IsActive {
			// Hide human components when inactive
			a.human.Head.Hide()
			a.human.Body.Hide()
			a.human.LeftArm.Hide()
			a.human.RightArm.Hide()
			a.human.LeftLeg.Hide()
			a.human.RightLeg.Hide()
		} else {
			// Show human components when active
			a.human.Head.Show()
			a.human.Body.Show()
			a.human.LeftArm.Show()
			a.human.RightArm.Show()
			a.human.LeftLeg.Show()
			a.human.RightLeg.Show()
		}
	})

	deathCountButton := widget.NewButton("💀 Death Count", func() {
		println("Human Deaths:", a.human.Deaths)
	})

	dragonButton := widget.NewButton("🐉 Toggle Dragon", func() {
		a.dragon.IsActive = !a.dragon.IsActive
		if !a.dragon.IsActive {
			a.dragon.Hide()
		} else {
			a.dragon.Show()
		}
	})

	resetButton := widget.NewButton("🔄 Reset All", func() {
		a.resetAll()
	})

	quitButton := widget.NewButton("❌ Quit", func() {
		a.fyneApp.Quit()
	})

	// Create button container with more buttons
	return container.NewHBox(
		startButton,
		stopButton,
		colorButton,
		speedUpButton,
		slowDownButton,
		massInfoButton,
		humanButton,
		deathCountButton,
		dragonButton,
		resetButton,
		quitButton,
	)
}

// resetAll resets all balls and human to initial state
func (a *App) resetAll() {
	// Reset ball1
	a.balls[0].X, a.balls[0].Y = 100, 100
	a.balls[0].VX, a.balls[0].VY = 3.5, 2.8
	a.balls[0].Circle.Move(fyne.NewPos(a.balls[0].X-a.balls[0].Radius, a.balls[0].Y-a.balls[0].Radius))

	// Reset ball2
	a.balls[1].X, a.balls[1].Y = 300, 200
	a.balls[1].VX, a.balls[1].VY = -2.8, 4.2
	a.balls[1].Circle.Move(fyne.NewPos(a.balls[1].X-a.balls[1].Radius, a.balls[1].Y-a.balls[1].Radius))

	// Reset ball3
	a.balls[2].X, a.balls[2].Y = 500, 150
	a.balls[2].VX, a.balls[2].VY = -4.1, -3.3
	a.balls[2].Circle.Move(fyne.NewPos(a.balls[2].X-a.balls[2].Radius, a.balls[2].Y-a.balls[2].Radius))

	// Reset human
	a.human.X, a.human.Y = 400, 300
	a.human.IsExploding = false
	a.human.IsActive = true
	a.human.RespawnTimer = 0
	a.human.Deaths = 0
	a.human.KeyUp = false
	a.human.KeyDown = false
	a.human.KeyLeft = false
	a.human.KeyRight = false

	// Show human components
	a.human.Head.Show()
	a.human.Body.Show()
	a.human.LeftArm.Show()
	a.human.RightArm.Show()
	a.human.LeftLeg.Show()
	a.human.RightLeg.Show()
	a.human.UpdatePosition()

	// Clear any explosion particles
	if a.human.ExplosionParticles != nil {
		for _, particle := range a.human.ExplosionParticles {
			if particle != nil {
				particle.Hide()
			}
		}
		a.human.ExplosionParticles = nil
	}

	// Reset dragon
	a.dragon.X, a.dragon.Y = 200, 200
	a.dragon.VX, a.dragon.VY = 0, 0
	a.dragon.IsActive = true
	a.dragon.IsDrifting = false
	a.dragon.DriftTimer = 0
	a.dragon.IsSpinning = false
	a.dragon.SpinAngle = 0
	a.dragon.SpinCount = 0
	a.dragon.Show()
	a.dragon.UpdatePosition()
}
