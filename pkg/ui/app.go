package ui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/atyronesmith/bouncing-balls/pkg/physics"
)

// App represents the main application
type App struct {
	fyneApp         fyne.App
	window          fyne.Window
	balls           []*physics.Ball
	human           *physics.Human
	dragon          *physics.Dragon
	currentBounds   fyne.Size
	animationTicker *time.Ticker
	content         *fyne.Container // Main content container for dynamic elements
}

// NewApp creates a new application instance
func NewApp() *App {
	return &App{
		fyneApp:       app.New(),
		currentBounds: fyne.NewSize(800, 600),
	}
}

// updateBounds updates the bounds for all physics objects
func (a *App) updateBounds(newSize fyne.Size) {
	a.currentBounds = newSize

	// Update bounds for all balls
	for _, ball := range a.balls {
		ball.Bounds = newSize
	}

	// Update bounds for human
	if a.human != nil {
		a.human.Bounds = newSize
	}

	// Update bounds for dragon
	if a.dragon != nil {
		a.dragon.Bounds = newSize
	}
}

// startAnimation starts the animation loop - simplified version
func (a *App) startAnimation() {
	// Start the animation automatically for all balls
	for _, ball := range a.balls {
		ball.IsAnimated = true
	}

	// Simple animation timer - updates all balls 60 times per second
	a.animationTicker = time.NewTicker(time.Millisecond * 16) // ~60 FPS
	go func() {
		defer a.animationTicker.Stop()
		for {
			select {
			case <-a.animationTicker.C:
				// Update all ball positions (wall bouncing)
				for _, ball := range a.balls {
					ball.Update()
				}

				// Update eyeball positions with human tracking (if human is active)
				if a.human != nil && a.human.IsActive {
					for _, ball := range a.balls {
						ball.UpdatePositionWithHuman(a.human.X, a.human.Y)
					}
				} else {
					// If no human, use default positioning
					for _, ball := range a.balls {
						ball.UpdatePosition()
					}
				}

				// Check for ball-to-ball collisions
				for i := 0; i < len(a.balls); i++ {
					for j := i + 1; j < len(a.balls); j++ {
						if a.balls[i].CheckCollision(a.balls[j]) {
							// Store explosion state before collision
							wasExploding1 := a.balls[i].IsExploding
							wasExploding2 := a.balls[j].IsExploding

							a.balls[i].HandleCollision(a.balls[j])

							// Add explosion particles to UI if explosion just started
							if !wasExploding1 && a.balls[i].IsExploding {
								for _, particle := range a.balls[i].GetExplosionParticles() {
									if particle != nil {
										a.content.Add(particle)
									}
								}
							}
							if !wasExploding2 && a.balls[j].IsExploding {
								for _, particle := range a.balls[j].GetExplosionParticles() {
									if particle != nil {
										a.content.Add(particle)
									}
								}
							}
						}
					}
				}

				// Update human
				if a.human != nil {
					if a.human.IsActive {
						// Store bullets before update for UI management
						bulletsBeforeUpdate := a.human.GetBulletVisuals()

						a.human.Update(a.balls)

						// Add new bullets to UI
						bulletsAfterUpdate := a.human.GetBulletVisuals()
						for _, bullet := range bulletsAfterUpdate {
							found := false
							for _, oldBullet := range bulletsBeforeUpdate {
								if bullet == oldBullet {
									found = true
									break
								}
							}
							if !found {
								a.content.Add(bullet)
							}
						}

						// Check ball-human collisions
						if a.human.CheckCollisionWithBalls(a.balls) {
							// Store previous explosion state
							wasExploding := a.human.IsExploding
							a.human.Explode()

							// If explosion just started, add particles to UI
							if !wasExploding && a.human.IsExploding {
								for _, particle := range a.human.ExplosionParticles {
									if particle != nil {
										a.content.Add(particle)
									}
								}
							}
						}
					}

					// Always update explosion state (handles respawn timer and animation)
					if a.human.IsExploding {
						// Store explosion particles before update (for cleanup)
						particlesBeforeUpdate := a.human.ExplosionParticles
						wasExploding := a.human.IsExploding

						a.human.UpdateExplosion()

						// If explosion just ended (respawn happened), clean up particles from UI
						if wasExploding && !a.human.IsExploding {
							// Remove explosion particles from UI
							for _, particle := range particlesBeforeUpdate {
								if particle != nil {
									a.content.Remove(particle)
								}
							}
						}
					}
				}

				// Update dragon if active
				if a.dragon != nil && a.dragon.IsActive {
					a.dragon.Update(a.balls)
					a.dragon.UpdatePosition()
				}
			}
		}
	}()
}

// Run starts the application
func (a *App) Run() {
	a.fyneApp.SetIcon(nil)

	// Create a fixed size window
	a.window = a.fyneApp.NewWindow("Bouncing Balls - Fixed Size (800x600)")
	a.window.Resize(a.currentBounds)
	a.window.CenterOnScreen()
	a.window.SetFixedSize(true) // Make window non-resizable

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

	// Create the human figure
	a.human = physics.NewHuman(400, 300, 35)

	// Create the dragon
	a.dragon = physics.NewDragon(200, 200, 40)

	// Update bounds for all objects
	a.updateBounds(a.currentBounds)

	// Create instruction label
	label := widget.NewLabel("⚡ Bouncing Balls - Fixed Size Window (Non-Resizable)")
	label.Alignment = fyne.TextAlignCenter

	// Create UI controls
	controls := a.createControls()

	// Create the main container
	a.content = container.NewWithoutLayout()

	// Add ball trails to container
	for _, ball := range a.balls {
		for _, trail := range ball.Trail {
			a.content.Add(trail)
		}
	}

	// Add balls (eyeball background, iris, pupils, and bloodshot veins)
	a.content.Add(ball1.Circle) // White eyeball background
	// Add bloodshot veins for ball1
	for _, vein := range ball1.BloodVeins {
		a.content.Add(vein)
	}
	a.content.Add(ball1.Iris)   // Colored iris
	a.content.Add(ball1.Pupil)  // Black pupil

	a.content.Add(ball2.Circle) // White eyeball background
	// Add bloodshot veins for ball2
	for _, vein := range ball2.BloodVeins {
		a.content.Add(vein)
	}
	a.content.Add(ball2.Iris)   // Colored iris
	a.content.Add(ball2.Pupil)  // Black pupil

	a.content.Add(ball3.Circle) // White eyeball background
	// Add bloodshot veins for ball3
	for _, vein := range ball3.BloodVeins {
		a.content.Add(vein)
	}
	a.content.Add(ball3.Iris)   // Colored iris
	a.content.Add(ball3.Pupil)  // Black pupil

	// Add ball text labels
	a.content.Add(ball1.Text)
	a.content.Add(ball2.Text)
	a.content.Add(ball3.Text)

	// Add human figure PNG image
	a.content.Add(a.human.ImageContainer)

	// Add dragon figure components
	dragonComponents := a.dragon.GetVisualComponents()
	for _, component := range dragonComponents {
		a.content.Add(component)
	}

	// Create the full layout
	fullContent := container.NewVBox(
		label,
		a.content,
		controls,
	)

	// Set the content
	a.window.SetContent(fullContent)

	// Start the animation
	a.startAnimation()

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

	resetButton := widget.NewButton("🔄 Reset All", func() {
		a.resetAll()
	})

	quitButton := widget.NewButton("❌ Quit", func() {
		a.fyneApp.Quit()
	})

	// Create a horizontal container for buttons
	return container.NewHBox(
		startButton,
		stopButton,
		colorButton,
		resetButton,
		quitButton,
	)
}

// resetAll resets all objects to their initial state
func (a *App) resetAll() {
	// Reset ball positions and velocities
	a.balls[0].X = 100
	a.balls[0].Y = 100
	a.balls[0].VX = 3.5
	a.balls[0].VY = 2.8

	a.balls[1].X = 300
	a.balls[1].Y = 200
	a.balls[1].VX = -2.8
	a.balls[1].VY = 4.2

	a.balls[2].X = 500
	a.balls[2].Y = 150
	a.balls[2].VX = -4.1
	a.balls[2].VY = -3.3

	// Update ball visual positions
	for _, ball := range a.balls {
		ball.Circle.Move(fyne.NewPos(ball.X-ball.Radius, ball.Y-ball.Radius))
	}

	// Reset human
	a.human.X = 400
	a.human.Y = 300
	a.human.IsExploding = false
	a.human.IsActive = true
	a.human.RespawnTimer = 0
	a.human.Rotation = 0 // Reset rotation
	// Show human image container
	a.human.ImageContainer.Show()
	a.human.UpdatePosition()

	// Reset dragon
	a.dragon.X = 200
	a.dragon.Y = 200
	a.dragon.VX = 0
	a.dragon.VY = 0
	a.dragon.IsActive = true
	a.dragon.Show()
	a.dragon.UpdatePosition()
}
