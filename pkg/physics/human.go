package physics

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Bullet represents a bullet fired by the human
type Bullet struct {
	X, Y     float32 // current position
	VX, VY   float32 // velocity
	Visual   *canvas.Circle
	IsActive bool
}

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
	KeyUp    bool // up arrow key pressed
	KeyDown  bool // down arrow key pressed
	KeyLeft  bool // left arrow key pressed
	KeyRight bool // right arrow key pressed
	// Visual components
	Head     *canvas.Circle
	Body     *canvas.Rectangle
	LeftArm  *canvas.Line
	RightArm *canvas.Line
	LeftLeg  *canvas.Rectangle
	RightLeg *canvas.Rectangle
	// Eye components
	LeftEye    *canvas.Circle // Left eye white
	RightEye   *canvas.Circle // Right eye white
	LeftPupil  *canvas.Circle // Left pupil (tracks balls)
	RightPupil *canvas.Circle // Right pupil (tracks balls)
	// Explosion particles
	ExplosionParticles []*canvas.Circle
	// Bullet system
	Bullets       []*Bullet
	ShootTimer    int // frames until next shot
	ShootCooldown int // frames between shots
}

// NewHuman creates a new human figure
func NewHuman(x, y, size float32) *Human {
	human := &Human{
		X:             x,
		Y:             y,
		Size:          size,
		Speed:         4.5, // Increased from 2.0 to 4.5 for much faster movement
		Bounds:        fyne.NewSize(800, 600),
		IsActive:      true,
		Bullets:       make([]*Bullet, 0),
		ShootTimer:    0,
		ShootCooldown: 15, // Shoot every 15 frames (4 times per second at 60 FPS)
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

	// Arms (lines)
	human.LeftArm = &canvas.Line{
		StrokeColor: humanColor,
		StrokeWidth: 4.0, // Thicker line to make arms more visible
	}
	// Set initial left arm position (from shoulder to default position)
	leftShoulderX := x - size*0.15
	leftShoulderY := y - size*0.1
	human.LeftArm.Position1 = fyne.NewPos(leftShoulderX, leftShoulderY)
	human.LeftArm.Position2 = fyne.NewPos(leftShoulderX-size*0.2, leftShoulderY+size*0.3)

	human.RightArm = &canvas.Line{
		StrokeColor: humanColor,
		StrokeWidth: 4.0, // Thicker line to make arms more visible
	}
	// Set initial right arm position (from shoulder to default position)
	rightShoulderX := x + size*0.15
	rightShoulderY := y - size*0.1
	human.RightArm.Position1 = fyne.NewPos(rightShoulderX, rightShoulderY)
	human.RightArm.Position2 = fyne.NewPos(rightShoulderX+size*0.2, rightShoulderY+size*0.3)

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

	// Eyes (circles on the head)
	eyeColor := color.RGBA{R: 255, G: 255, B: 255, A: 255} // White eyes
	pupilColor := color.RGBA{R: 0, G: 0, B: 0, A: 255}     // Black pupils

	// Left eye
	human.LeftEye = &canvas.Circle{
		FillColor:   eyeColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	human.LeftEye.Resize(fyne.NewSize(size*0.08, size*0.08))

	// Right eye
	human.RightEye = &canvas.Circle{
		FillColor:   eyeColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	human.RightEye.Resize(fyne.NewSize(size*0.08, size*0.08))

	// Left pupil
	human.LeftPupil = &canvas.Circle{
		FillColor:   pupilColor,
		StrokeColor: pupilColor,
		StrokeWidth: 0.0,
	}
	human.LeftPupil.Resize(fyne.NewSize(size*0.04, size*0.04))

	// Right pupil
	human.RightPupil = &canvas.Circle{
		FillColor:   pupilColor,
		StrokeColor: pupilColor,
		StrokeWidth: 0.0,
	}
	human.RightPupil.Resize(fyne.NewSize(size*0.04, size*0.04))

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

	// Legs positions (bottom of body)
	h.LeftLeg.Move(fyne.NewPos(h.X-h.Size*0.12, h.Y+h.Size*0.1))
	h.RightLeg.Move(fyne.NewPos(h.X+h.Size*0.06, h.Y+h.Size*0.1))

	// Eye positions (on the head)
	h.LeftEye.Move(fyne.NewPos(h.X-h.Size*0.12, h.Y-h.Size*0.55))
	h.RightEye.Move(fyne.NewPos(h.X+h.Size*0.04, h.Y-h.Size*0.55))

	// Pupil positions (centered in eyes)
	h.LeftPupil.Move(fyne.NewPos(h.X-h.Size*0.1, h.Y-h.Size*0.53))
	h.RightPupil.Move(fyne.NewPos(h.X+h.Size*0.06, h.Y-h.Size*0.53))
}

// findClosestBall finds the closest animated ball to the human
func (h *Human) findClosestBall(balls []*Ball) *Ball {
	var closestBall *Ball
	minDistance := float32(math.Inf(1))

	for _, ball := range balls {
		if !ball.IsAnimated {
			continue
		}

		// Calculate distance to ball
		dx := h.X - ball.X
		dy := h.Y - ball.Y
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		if distance < minDistance {
			minDistance = distance
			closestBall = ball
		}
	}

	return closestBall
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
		// AI mode: combine avoidance with centering behavior
		avoidanceX, avoidanceY := h.calculateAvoidance(balls)
		centeringX, centeringY := h.calculateCentering()

		if avoidanceX != 0 || avoidanceY != 0 {
			// In danger: prioritize avoidance but add slight centering influence
			moveX = avoidanceX*0.9 + centeringX*0.1
			moveY = avoidanceY*0.9 + centeringY*0.1
		} else {
			// Safe: move toward center
			moveX = centeringX
			moveY = centeringY
		}
	}

	// Apply movement
	h.X += moveX
	h.Y += moveY

	// Keep human within bounds
	h.keepWithinBounds()

	// Update visual position
	h.UpdatePosition()

	// Update arms to point at closest ball
	h.UpdatePointing(balls)

	// Update shooting system
	h.UpdateShooting(balls)

	// Update bullet positions
	h.UpdateBullets()

	// Check bullet collisions with balls
	h.CheckBulletCollisions(balls)
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

// calculateCentering calculates movement toward the center of the panel
func (h *Human) calculateCentering() (float32, float32) {
	// Calculate center of the panel
	centerX := h.Bounds.Width / 2
	centerY := (h.Bounds.Height + 50) / 2 // Account for button area (50px top, 50px bottom)

	// Calculate distance to center
	dx := centerX - h.X
	dy := centerY - h.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// If already very close to center, don't move
	if distance < 20 {
		return 0, 0
	}

	// Normalize direction and apply gentle centering force
	if distance > 0 {
		centeringStrength := h.Speed * 0.3 // Gentle centering force (30% of normal speed)
		normalizedDx := dx / distance
		normalizedDy := dy / distance

		return normalizedDx * centeringStrength, normalizedDy * centeringStrength
	}

	return 0, 0
}

// keepWithinBounds ensures human stays within window boundaries with wrap-around for horizontal edges
func (h *Human) keepWithinBounds() {
	margin := h.Size * 0.5

	// Horizontal wrap-around: if human goes off left edge, appear on right edge (and vice versa)
	if h.X < -margin {
		h.X = h.Bounds.Width + margin // Appear on right side
	} else if h.X > h.Bounds.Width+margin {
		h.X = -margin // Appear on left side
	}

	// Vertical boundaries: still clamp to prevent going off top/bottom
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
	h.LeftEye.Hide()
	h.RightEye.Hide()
	h.LeftPupil.Hide()
	h.RightPupil.Hide()

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
				alphaInt := 255 - int(explosionFrame*4)
				if alphaInt < 0 {
					alphaInt = 0
				}
				alpha := uint8(alphaInt)
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
	h.LeftEye.Show()
	h.RightEye.Show()
	h.LeftPupil.Show()
	h.RightPupil.Show()

	h.UpdatePosition()

	// Clear explosion particles
	h.ExplosionParticles = nil
}

// UpdatePointing updates the arms to point at the closest ball
func (h *Human) UpdatePointing(balls []*Ball) {
	if !h.IsActive {
		return
	}

	// Calculate shoulder positions based on current human position
	leftShoulderX := h.X - h.Size*0.15
	leftShoulderY := h.Y - h.Size*0.1
	rightShoulderX := h.X + h.Size*0.15
	rightShoulderY := h.Y - h.Size*0.1

	// Find the closest ball
	closestBall := h.findClosestBall(balls)
	if closestBall == nil {
		// No balls to track, arms should be in default position
		h.LeftArm.Position1 = fyne.NewPos(leftShoulderX, leftShoulderY)
		h.LeftArm.Position2 = fyne.NewPos(leftShoulderX-h.Size*0.2, leftShoulderY+h.Size*0.3)
		h.RightArm.Position1 = fyne.NewPos(rightShoulderX, rightShoulderY)
		h.RightArm.Position2 = fyne.NewPos(rightShoulderX+h.Size*0.2, rightShoulderY+h.Size*0.3)

		// Refresh the arms to update their visual appearance
		h.LeftArm.Refresh()
		h.RightArm.Refresh()
		return
	}

	// Calculate angle to closest ball from human center
	dx := closestBall.X - h.X
	dy := closestBall.Y - h.Y
	angle := math.Atan2(float64(dy), float64(dx))

	// Calculate arm length for pointing
	armLength := h.Size * 0.35

	// Both arms point toward the ball
	// Calculate end positions of arms pointing toward the ball
	leftArmEndX := leftShoulderX + float32(math.Cos(angle))*armLength
	leftArmEndY := leftShoulderY + float32(math.Sin(angle))*armLength
	rightArmEndX := rightShoulderX + float32(math.Cos(angle))*armLength
	rightArmEndY := rightShoulderY + float32(math.Sin(angle))*armLength

	// Update arm positions to point toward the ball
	h.LeftArm.Position1 = fyne.NewPos(leftShoulderX, leftShoulderY)
	h.LeftArm.Position2 = fyne.NewPos(leftArmEndX, leftArmEndY)
	h.RightArm.Position1 = fyne.NewPos(rightShoulderX, rightShoulderY)
	h.RightArm.Position2 = fyne.NewPos(rightArmEndX, rightArmEndY)

	// Refresh the arms to update their visual appearance
	h.LeftArm.Refresh()
	h.RightArm.Refresh()
}

// NewBullet creates a new bullet at the specified position with velocity toward target
func NewBullet(startX, startY, targetX, targetY float32) *Bullet {
	// Calculate direction vector to target
	dx := targetX - startX
	dy := targetY - startY
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Normalize direction and set bullet speed
	bulletSpeed := float32(8.0) // Fast bullet speed
	vx := (dx / distance) * bulletSpeed
	vy := (dy / distance) * bulletSpeed

	bullet := &Bullet{
		X:        startX,
		Y:        startY,
		VX:       vx,
		VY:       vy,
		IsActive: true,
	}

	// Create visual representation - small yellow circle
	bullet.Visual = &canvas.Circle{
		FillColor:   color.RGBA{R: 255, G: 255, B: 0, A: 255}, // Yellow bullet
		StrokeColor: color.RGBA{R: 255, G: 200, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	bullet.Visual.Resize(fyne.NewSize(4, 4)) // Small bullet
	bullet.Visual.Move(fyne.NewPos(startX-2, startY-2))

	return bullet
}

// UpdateBullets updates all bullet positions and removes inactive ones
func (h *Human) UpdateBullets() {
	// Update existing bullets
	for i := len(h.Bullets) - 1; i >= 0; i-- {
		bullet := h.Bullets[i]
		if !bullet.IsActive {
			continue
		}

		// Update bullet position
		bullet.X += bullet.VX
		bullet.Y += bullet.VY
		bullet.Visual.Move(fyne.NewPos(bullet.X-2, bullet.Y-2))

		// Remove bullets that go off screen
		if bullet.X < 0 || bullet.X > h.Bounds.Width || bullet.Y < 0 || bullet.Y > h.Bounds.Height {
			bullet.IsActive = false
			bullet.Visual.Hide()
			// Remove from slice
			h.Bullets = append(h.Bullets[:i], h.Bullets[i+1:]...)
		}
	}
}

// ShootAtTarget creates bullets from both arms toward the target
func (h *Human) ShootAtTarget(targetX, targetY float32) {
	if !h.IsActive || h.IsExploding {
		return
	}

	// Get arm endpoints (where bullets spawn from)
	leftArmEndX := h.LeftArm.Position2.X
	leftArmEndY := h.LeftArm.Position2.Y
	rightArmEndX := h.RightArm.Position2.X
	rightArmEndY := h.RightArm.Position2.Y

	// Create bullets from both arms
	leftBullet := NewBullet(leftArmEndX, leftArmEndY, targetX, targetY)
	rightBullet := NewBullet(rightArmEndX, rightArmEndY, targetX, targetY)

	h.Bullets = append(h.Bullets, leftBullet, rightBullet)
}

// UpdateShooting handles the shooting timer and creates bullets when ready
func (h *Human) UpdateShooting(balls []*Ball) {
	if !h.IsActive || h.IsExploding {
		return
	}

	// Decrement shoot timer
	if h.ShootTimer > 0 {
		h.ShootTimer--
		return
	}

	// Find closest ball to shoot at
	closestBall := h.findClosestBall(balls)
	if closestBall == nil {
		return
	}

	// Shoot at the closest ball
	h.ShootAtTarget(closestBall.X, closestBall.Y)

	// Reset shoot timer
	h.ShootTimer = h.ShootCooldown
}

// CheckBulletCollisions checks if any bullets hit any balls and handles the collision
func (h *Human) CheckBulletCollisions(balls []*Ball) {
	for i := len(h.Bullets) - 1; i >= 0; i-- {
		bullet := h.Bullets[i]
		if !bullet.IsActive {
			continue
		}

		for _, ball := range balls {
			if !ball.IsAnimated {
				continue
			}

			// Check collision between bullet and ball
			dx := bullet.X - ball.X
			dy := bullet.Y - ball.Y
			distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

			if distance < ball.Radius+2 { // bullet radius is 2
				// Bullet hit ball!
				bullet.IsActive = false
				bullet.Visual.Hide()

				// Make ball grow slightly
				growthAmount := float32(0.8) // Small growth per hit
				ball.Radius += growthAmount
				ball.OriginalRadius += growthAmount

				// Update ball visual size
				ball.Circle.Resize(fyne.NewSize(ball.Radius*2, ball.Radius*2))
				ball.updateTextSize() // Adjust text size for new ball size

				// Remove bullet from slice
				h.Bullets = append(h.Bullets[:i], h.Bullets[i+1:]...)
				break // Bullet can only hit one ball
			}
		}
	}
}

// GetBulletVisuals returns all bullet visual objects for UI management
func (h *Human) GetBulletVisuals() []*canvas.Circle {
	visuals := make([]*canvas.Circle, 0, len(h.Bullets))
	for _, bullet := range h.Bullets {
		if bullet.IsActive {
			visuals = append(visuals, bullet.Visual)
		}
	}
	return visuals
}
