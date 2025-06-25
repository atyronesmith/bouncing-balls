package physics

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Human represents a human that avoids the balls
type Human struct {
	X, Y         float32   // current position
	Size         float32   // size (similar to ball radius)
	Speed        float32   // movement speed
	Bounds       fyne.Size // movement bounds
	IsActive     bool      // whether the human is active
	IsExploding  bool      // whether the human is currently exploding
	RespawnTimer int       // frames until respawn
	Deaths       int       // death counter
	// Keyboard control state
	KeyUp        bool      // up arrow key pressed
	KeyDown      bool      // down arrow key pressed
	KeyLeft      bool      // left arrow key pressed
	KeyRight     bool      // right arrow key pressed
	// Visual components
	Head     *canvas.Circle
	Body     *canvas.Rectangle
	LeftArm  *canvas.Rectangle
	RightArm *canvas.Rectangle
	LeftLeg  *canvas.Rectangle
	RightLeg *canvas.Rectangle
	// Explosion particles
	ExplosionParticles []*canvas.Circle
}

// NewHuman creates a new human figure
func NewHuman(x, y, size float32) *Human {
	human := &Human{
		X:        x,
		Y:        y,
		Size:     size,
		Speed:    4.5, // Increased from 2.0 to 4.5 for much faster movement
		Bounds:   fyne.NewSize(800, 600),
		IsActive: true,
	}

	// Create human figure components
	humanColor := color.RGBA{R: 255, G: 220, B: 177, A: 255}   // Skin color
	clothingColor := color.RGBA{R: 70, G: 130, B: 180, A: 255} // Blue clothing

	// Head (circle)
	human.Head = &canvas.Circle{
		FillColor:   humanColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 2.0,
	}
	human.Head.Resize(fyne.NewSize(size*0.4, size*0.4))

	// Body (rectangle)
	human.Body = &canvas.Rectangle{
		FillColor:   clothingColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	human.Body.Resize(fyne.NewSize(size*0.3, size*0.6))

	// Arms (rectangles)
	human.LeftArm = &canvas.Rectangle{
		FillColor:   humanColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	human.LeftArm.Resize(fyne.NewSize(size*0.15, size*0.4))

	human.RightArm = &canvas.Rectangle{
		FillColor:   humanColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	human.RightArm.Resize(fyne.NewSize(size*0.15, size*0.4))

	// Legs (rectangles)
	human.LeftLeg = &canvas.Rectangle{
		FillColor:   clothingColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	human.LeftLeg.Resize(fyne.NewSize(size*0.12, size*0.5))

	human.RightLeg = &canvas.Rectangle{
		FillColor:   clothingColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	human.RightLeg.Resize(fyne.NewSize(size*0.12, size*0.5))

	// Set initial position
	human.UpdatePosition()

	return human
}

// UpdatePosition updates the visual position of all human components
func (h *Human) UpdatePosition() {
	if !h.IsActive {
		return
	}

	// Head position (top center)
	h.Head.Move(fyne.NewPos(h.X-h.Size*0.2, h.Y-h.Size*0.6))

	// Body position (center)
	h.Body.Move(fyne.NewPos(h.X-h.Size*0.15, h.Y-h.Size*0.2))

	// Arms positions (sides of body)
	h.LeftArm.Move(fyne.NewPos(h.X-h.Size*0.3, h.Y-h.Size*0.15))
	h.RightArm.Move(fyne.NewPos(h.X+h.Size*0.15, h.Y-h.Size*0.15))

	// Legs positions (bottom of body)
	h.LeftLeg.Move(fyne.NewPos(h.X-h.Size*0.12, h.Y+h.Size*0.1))
	h.RightLeg.Move(fyne.NewPos(h.X+h.Size*0.06, h.Y+h.Size*0.1))
}

// Update handles both keyboard input and AI avoidance behavior
func (h *Human) Update(balls []*Ball) {
	if !h.IsActive {
		return
	}

	// Check if any keyboard keys are pressed
	keyboardActive := h.KeyUp || h.KeyDown || h.KeyLeft || h.KeyRight

	var moveX, moveY float32

	if keyboardActive {
		// Keyboard control mode - player has partial control
		if h.KeyUp {
			moveY -= h.Speed
		}
		if h.KeyDown {
			moveY += h.Speed
		}
		if h.KeyLeft {
			moveX -= h.Speed
		}
		if h.KeyRight {
			moveX += h.Speed
		}

		// Still apply AI avoidance as a safety override when in extreme danger
		avoidanceX, avoidanceY := h.calculateAvoidance(balls)

		// If there's significant danger, blend keyboard input with avoidance
		if avoidanceX != 0 || avoidanceY != 0 {
			// Reduce keyboard influence and add avoidance influence
			moveX = moveX*0.3 + avoidanceX*0.7
			moveY = moveY*0.3 + avoidanceY*0.7
		}
	} else {
		// Pure AI avoidance mode (original behavior)
		moveX, moveY = h.calculateAvoidance(balls)
	}

	// Apply movement
	h.X += moveX
	h.Y += moveY

	// Keep human within bounds
	h.keepWithinBounds()

	// Update visual position
	h.UpdatePosition()
}

// calculateAvoidance calculates AI avoidance movement (extracted from original AvoidBalls method)
func (h *Human) calculateAvoidance(balls []*Ball) (float32, float32) {
	var totalAvoidanceX, totalAvoidanceY float32
	dangerCount := 0
	maxDanger := float32(0)

	for _, ball := range balls {
		if !ball.IsAnimated {
			continue
		}

		// Calculate distance to ball
		dx := h.X - ball.X
		dy := h.Y - ball.Y
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		// Much larger safety margin for earlier detection
		dangerDistance := h.Size + ball.Radius + 120 // Increased from 50 to 120

		if distance < dangerDistance {
			// Calculate multiple future positions for better prediction
			for t := float32(5); t <= 20; t += 5 { // Check 5, 10, 15, 20 frames ahead
				futureBallX := ball.X + ball.VX*t
				futureBallY := ball.Y + ball.VY*t

				// Check if ball will be close to human
				futureDx := h.X - futureBallX
				futureDy := h.Y - futureBallY
				futureDistance := float32(math.Sqrt(float64(futureDx*futureDx + futureDy*futureDy)))

				if futureDistance < distance*0.8 { // Ball is approaching
					// Calculate avoidance direction (away from ball)
					if distance > 0 {
						// More aggressive avoidance strength
						avoidanceStrength := (dangerDistance - distance) / dangerDistance
						avoidanceStrength = avoidanceStrength * avoidanceStrength * 2.0 // Quadratic increase

						normalizedDx := dx / distance
						normalizedDy := dy / distance

						totalAvoidanceX += normalizedDx * avoidanceStrength
						totalAvoidanceY += normalizedDy * avoidanceStrength
						dangerCount++

						if avoidanceStrength > maxDanger {
							maxDanger = avoidanceStrength
						}
					}
					break // Found threat, move to next ball
				}
			}
		}
	}

	// Calculate avoidance movement
	if dangerCount > 0 {
		// Normalize avoidance vector
		avgAvoidanceX := totalAvoidanceX / float32(dangerCount)
		avgAvoidanceY := totalAvoidanceY / float32(dangerCount)

		// Apply panic speed boost when in extreme danger
		speedMultiplier := h.Speed
		if maxDanger > 0.7 {
			speedMultiplier *= 2.0 // Double speed in extreme danger
		} else if maxDanger > 0.4 {
			speedMultiplier *= 1.5 // 50% speed boost in moderate danger
		}

		return avgAvoidanceX * speedMultiplier, avgAvoidanceY * speedMultiplier
	}

	return 0, 0
}

// keepWithinBounds ensures human stays within window boundaries
func (h *Human) keepWithinBounds() {
	margin := h.Size * 0.5 // Reduced margin for more movement space
	if h.X < margin {
		h.X = margin
	} else if h.X > h.Bounds.Width-margin {
		h.X = h.Bounds.Width - margin
	}
	if h.Y < 50+margin { // Account for button area
		h.Y = 50 + margin
	} else if h.Y > h.Bounds.Height-50-margin {
		h.Y = h.Bounds.Height - 50 - margin
	}
}

// CheckCollisionWithBalls checks if the human collides with any ball
func (h *Human) CheckCollisionWithBalls(balls []*Ball) bool {
	if !h.IsActive || h.IsExploding {
		return false
	}

	for _, ball := range balls {
		if !ball.IsAnimated {
			continue
		}

		// Calculate distance between human center and ball center
		dx := h.X - ball.X
		dy := h.Y - ball.Y
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		// Check if collision occurs (human size is used as radius)
		collisionDistance := h.Size*0.6 + ball.Radius // Slightly smaller than full size for fairness
		if distance < collisionDistance {
			return true
		}
	}
	return false
}

// Explode creates an explosion effect and starts respawn timer
func (h *Human) Explode() {
	if h.IsExploding {
		return
	}

	h.IsExploding = true
	h.Deaths++
	h.RespawnTimer = 180 // 3 seconds at 60 FPS

	// Hide human components
	h.Head.Hide()
	h.Body.Hide()
	h.LeftArm.Hide()
	h.RightArm.Hide()
	h.LeftLeg.Hide()
	h.RightLeg.Hide()

	// Create explosion particles
	h.ExplosionParticles = make([]*canvas.Circle, 12) // 12 explosion particles
	explosionColors := []color.RGBA{
		{R: 255, G: 255, B: 0, A: 255},   // Yellow
		{R: 255, G: 165, B: 0, A: 255},   // Orange
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 255, G: 255, B: 255, A: 255}, // White
	}

	for i := 0; i < 12; i++ {
		particle := &canvas.Circle{
			FillColor:   explosionColors[i%len(explosionColors)],
			StrokeColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			StrokeWidth: 1.0,
		}
		particle.Resize(fyne.NewSize(8, 8))
		particle.Move(fyne.NewPos(h.X-4, h.Y-4))
		h.ExplosionParticles[i] = particle
	}
}

// UpdateExplosion updates the explosion animation
func (h *Human) UpdateExplosion() {
	if !h.IsExploding {
		return
	}

	h.RespawnTimer--

	// Animate explosion particles
	if h.RespawnTimer > 120 { // First 1 second - explosion expanding
		explosionFrame := 180 - h.RespawnTimer
		for i, particle := range h.ExplosionParticles {
			if particle != nil {
				// Calculate particle movement in different directions
				angle := float64(i) * 2 * math.Pi / float64(len(h.ExplosionParticles))
				radius := float32(explosionFrame) * 2.0

				newX := h.X + float32(math.Cos(angle))*radius - 4
				newY := h.Y + float32(math.Sin(angle))*radius - 4

				particle.Move(fyne.NewPos(newX, newY))

				// Fade particles
				alpha := 255 - uint8(explosionFrame*4)
				if alpha > 255 {
					alpha = 0
				}
				particle.FillColor = color.RGBA{
					R: particle.FillColor.(color.RGBA).R,
					G: particle.FillColor.(color.RGBA).G,
					B: particle.FillColor.(color.RGBA).B,
					A: alpha,
				}
				particle.Refresh()
			}
		}
	} else if h.RespawnTimer == 120 {
		// Hide explosion particles
		for _, particle := range h.ExplosionParticles {
			if particle != nil {
				particle.Hide()
			}
		}
	}

	// Respawn human
	if h.RespawnTimer <= 0 {
		h.Respawn()
	}
}

// Respawn creates a new human at a safe location
func (h *Human) Respawn() {
	h.IsExploding = false
	h.IsActive = true

	// Find a safe spawn location (away from balls)
	safeLocations := []fyne.Position{
		{X: 100, Y: 150},
		{X: 700, Y: 150},
		{X: 400, Y: 500},
		{X: 200, Y: 300},
		{X: 600, Y: 300},
	}

	// Choose the safest location
	h.X = safeLocations[h.Deaths%len(safeLocations)].X
	h.Y = safeLocations[h.Deaths%len(safeLocations)].Y

	// Show human components
	h.Head.Show()
	h.Body.Show()
	h.LeftArm.Show()
	h.RightArm.Show()
	h.LeftLeg.Show()
	h.RightLeg.Show()

	h.UpdatePosition()

	// Clear explosion particles
	h.ExplosionParticles = nil
}
