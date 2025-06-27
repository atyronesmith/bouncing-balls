package physics

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
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
	// Visual components - replaced with PNG image
	Image            *canvas.Image   // PNG image with transparency
	ImageContainer   *fyne.Container // Container to hold the image
	DirectionArrow   *canvas.Line    // Arrow showing facing direction
	Rotation         float64         // Current rotation angle in radians
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
		Rotation:      0,  // Start facing right (0 radians)
	}

	// Load PNG image with transparency
	resource := storage.NewFileURI("./human.png")
	human.Image = canvas.NewImageFromURI(resource)
	human.Image.FillMode = canvas.ImageFillOriginal
	human.Image.ScaleMode = canvas.ImageScaleSmooth

	// Set the image size based on the human size parameter
	imageSize := fyne.NewSize(size*1.2, size*1.2)
	human.Image.Resize(imageSize)

	// Create direction arrow to show facing direction
	human.DirectionArrow = &canvas.Line{
		StrokeColor: color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red arrow
		StrokeWidth: 3.0,
	}

	// Create a container to hold the image and arrow
	human.ImageContainer = container.NewWithoutLayout(human.Image, human.DirectionArrow)

	// Set initial position
	human.UpdatePosition()

	return human
}

// UpdatePosition updates the visual position of the human image with direction arrow
func (h *Human) UpdatePosition() {
	if !h.IsActive {
		return
	}

	// Center the image on the human's position (no jumping around)
	imageSize := h.Image.Size()
	baseX := h.X - imageSize.Width/2
	baseY := h.Y - imageSize.Height/2

	// Always keep the image centered at the human's actual position
	h.Image.Move(fyne.NewPos(baseX, baseY))

	// Keep consistent size
	h.Image.Resize(fyne.NewSize(h.Size*1.2, h.Size*1.2))

	// Update direction arrow to show facing direction
	h.updateDirectionArrow()
}

// updateDirectionArrow updates the direction arrow to show which way the human is facing
func (h *Human) updateDirectionArrow() {
	if h.Rotation == 0 {
		// Hide arrow when not rotating
		h.DirectionArrow.Position1 = fyne.NewPos(h.X, h.Y)
		h.DirectionArrow.Position2 = fyne.NewPos(h.X, h.Y)
		return
	}

	// Calculate arrow direction based on rotation
	arrowLength := h.Size * 0.8
	endX := h.X + float32(math.Cos(h.Rotation))*arrowLength
	endY := h.Y + float32(math.Sin(h.Rotation))*arrowLength

	// Set arrow from human center to direction
	h.DirectionArrow.Position1 = fyne.NewPos(h.X, h.Y)
	h.DirectionArrow.Position2 = fyne.NewPos(endX, endY)
}

// GetFacingDirection returns a string description of which direction the human is facing
func (h *Human) GetFacingDirection() string {
	if h.Rotation == 0 {
		return "right"
	}

	degrees := h.Rotation * 180 / math.Pi

	if degrees >= -45 && degrees <= 45 {
		return "right"
	} else if degrees > 45 && degrees <= 135 {
		return "down"
	} else if degrees > 135 || degrees <= -135 {
		return "left"
	} else {
		return "up"
	}
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

// UpdateRotation calculates and updates the rotation to face the closest ball
func (h *Human) UpdateRotation(balls []*Ball) {
	if !h.IsActive {
		return
	}

	// Find the closest ball
	closestBall := h.findClosestBall(balls)
	if closestBall == nil {
		return
	}

	// Calculate angle to the closest ball
	dx := closestBall.X - h.X
	dy := closestBall.Y - h.Y
	targetAngle := math.Atan2(float64(dy), float64(dx))

	// Update rotation
	h.Rotation = targetAngle
}

// Update handles both keyboard input and AI avoidance behavior
func (h *Human) Update(balls []*Ball) {
	if !h.IsActive {
		return
	}

	// Update rotation to face closest ball
	h.UpdateRotation(balls)

	// Calculate avoidance force from all balls
	avoidX, avoidY := h.calculateAvoidance(balls)

	// Calculate centering force to stay in bounds
	centerX, centerY := h.calculateCentering()

	// Combine forces (avoidance has higher priority)
	totalForceX := avoidX*0.8 + centerX*0.2
	totalForceY := avoidY*0.8 + centerY*0.2

	// Normalize force if too strong
	forceLength := float32(math.Sqrt(float64(totalForceX*totalForceX + totalForceY*totalForceY)))
	if forceLength > h.Speed {
		totalForceX = (totalForceX / forceLength) * h.Speed
		totalForceY = (totalForceY / forceLength) * h.Speed
	}

	// Apply movement
	h.X += totalForceX
	h.Y += totalForceY

	// Keep within bounds
	h.keepWithinBounds()

	// Update visual position
	h.UpdatePosition()

	// Update pointing (now just a stub)
	h.UpdatePointing(balls)

	// Update bullets
	h.UpdateBullets()

	// Update shooting
	h.UpdateShooting(balls)

	// Check bullet collisions
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
	centerY := h.Bounds.Height / 2 // Center of game area

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
	if h.Y < margin { // Use actual bounds
		h.Y = margin
	} else if h.Y > h.Bounds.Height-margin {
		h.Y = h.Bounds.Height - margin
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

// Explode creates an explosion effect and hides the human
func (h *Human) Explode() {
	if h.IsExploding {
		return // Already exploding
	}

	h.IsExploding = true
	h.IsActive = false
	h.RespawnTimer = 180 // 3 seconds at 60 FPS
	h.Deaths++           // Increment death counter

	// Hide human components
	h.ImageContainer.Hide()
	h.DirectionArrow.Hide()

	// Create explosion particles
	numParticles := 12
	h.ExplosionParticles = make([]*canvas.Circle, numParticles)

	colors := []color.RGBA{
		{R: 255, G: 100, B: 100, A: 255}, // Red
		{R: 255, G: 200, B: 100, A: 255}, // Orange
		{R: 255, G: 255, B: 100, A: 255}, // Yellow
		{R: 100, G: 255, B: 100, A: 255}, // Green
	}

	for i := 0; i < numParticles; i++ {
		particle := &canvas.Circle{
			FillColor: colors[i%len(colors)],
		}
		particle.Resize(fyne.NewSize(8, 8))
		// Position particles at human's center initially
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

// Respawn brings the human back to life at a safe location (fallback method)
func (h *Human) Respawn() {
	// Default fallback - center position
	safeX := h.Bounds.Width / 2
	safeY := h.Bounds.Height / 2

	// Set new position
	h.X = safeX
	h.Y = safeY

	// Reset state
	h.IsExploding = false
	h.IsActive = true
	h.RespawnTimer = 0
	h.Rotation = 0 // Reset rotation

	// Show human components
	h.ImageContainer.Show()
	h.DirectionArrow.Show()

	h.UpdatePosition()
}

// RespawnWithBalls brings the human back to life at the safest location away from all balls
func (h *Human) RespawnWithBalls(balls []*Ball) {
	safeX, safeY := h.findSafestRespawnLocation(balls)

	// Set new position
	h.X = safeX
	h.Y = safeY

	// Reset state
	h.IsExploding = false
	h.IsActive = true
	h.RespawnTimer = 0
	h.Rotation = 0 // Reset rotation

	// Show human components
	h.ImageContainer.Show()
	h.DirectionArrow.Show()

	h.UpdatePosition()
}

// findSafestRespawnLocation finds the position that maximizes distance from all balls
func (h *Human) findSafestRespawnLocation(balls []*Ball) (float32, float32) {
	// Define search grid parameters
	gridSize := 20 // 20x20 grid for reasonable performance
	margin := h.Size + 10 // Margin from screen edges

	bestX := h.Bounds.Width / 2
	bestY := h.Bounds.Height / 2
	maxMinDistance := float32(0) // Maximum of minimum distances to all balls

	// Search through grid positions
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			// Calculate candidate position
			x := margin + (float32(i)/float32(gridSize-1))*(h.Bounds.Width-2*margin)
			y := 50 + margin + (float32(j)/float32(gridSize-1))*(h.Bounds.Height-100-2*margin) // Account for UI area

			// Find minimum distance to all balls from this position
			minDistanceToBalls := float32(math.Inf(1))

			for _, ball := range balls {
				if !ball.IsAnimated {
					continue
				}

				// Calculate distance to this ball
				dx := x - ball.X
				dy := y - ball.Y
				distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

				// Account for ball radius and human size for true clearance
				clearance := distance - ball.Radius - h.Size

				if clearance < minDistanceToBalls {
					minDistanceToBalls = clearance
				}
			}

			// If this position has better minimum distance, use it
			if minDistanceToBalls > maxMinDistance {
				maxMinDistance = minDistanceToBalls
				bestX = x
				bestY = y
			}
		}
	}

	return bestX, bestY
}

// UpdatePointing updates the arms to point at the closest ball
// UpdatePointing is no longer needed since we're using a PNG image
// The human image will show a static pose
func (h *Human) UpdatePointing(balls []*Ball) {
	// No longer needed with PNG image - keeping empty for compatibility
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

	// Create visual representation - large, very visible bullet
	bullet.Visual = &canvas.Circle{
		FillColor:   color.RGBA{R: 255, G: 0, B: 255, A: 255}, // Bright magenta bullet
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Black border for visibility
		StrokeWidth: 3.0,
	}
	bullet.Visual.Resize(fyne.NewSize(16, 16)) // Much larger bullet for testing
	bullet.Visual.Move(fyne.NewPos(startX-8, startY-8)) // Center the larger bullet

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
		bullet.Visual.Move(fyne.NewPos(bullet.X-8, bullet.Y-8))

		// Remove bullets that go off screen
		if bullet.X < 0 || bullet.X > h.Bounds.Width || bullet.Y < 0 || bullet.Y > h.Bounds.Height {
			bullet.IsActive = false
			bullet.Visual.Hide()
			// Remove from slice
			h.Bullets = append(h.Bullets[:i], h.Bullets[i+1:]...)
		}
	}
}

// ShootAtTarget creates bullets from the end of the direction arrow toward the target
func (h *Human) ShootAtTarget(targetX, targetY float32) {
	if !h.IsActive || h.IsExploding {
		return
	}

	// Calculate bullet spawn position from the end of the direction arrow
	var bulletX, bulletY float32

	if h.Rotation == 0 {
		// No rotation, spawn from slightly above center
		bulletX = h.X
		bulletY = h.Y - h.Size*0.3
	} else {
		// Spawn from the end of the direction arrow
		arrowLength := h.Size * 0.8
		bulletX = h.X + float32(math.Cos(h.Rotation))*arrowLength
		bulletY = h.Y + float32(math.Sin(h.Rotation))*arrowLength
	}

	// Create bullet from the arrow tip position
	bullet := NewBullet(bulletX, bulletY, targetX, targetY)
	h.Bullets = append(h.Bullets, bullet)
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

			if distance < ball.Radius+8 { // bullet radius is now 8
				// Bullet hit ball!
				bullet.IsActive = false
				bullet.Visual.Hide()

				// Apply repulsion force to the ball
				if distance > 0 {
					// Calculate repulsion direction (away from bullet impact point)
					// dx = bullet.X - ball.X, so -dx points from bullet to ball (true repulsion)
					repelX := -dx / distance // Normalized direction away from bullet
					repelY := -dy / distance

					// Apply repulsion force to ball velocity
					repulsionStrength := float32(0.8) // Adjust this to control repulsion intensity
					ball.VX += repelX * repulsionStrength
					ball.VY += repelY * repulsionStrength

					// Optional: Add slight speed dampening to prevent balls from going too fast
					maxSpeed := float32(8.0)
					currentSpeed := float32(math.Sqrt(float64(ball.VX*ball.VX + ball.VY*ball.VY)))
					if currentSpeed > maxSpeed {
						ball.VX = (ball.VX / currentSpeed) * maxSpeed
						ball.VY = (ball.VY / currentSpeed) * maxSpeed
					}

					// Trigger a subtle jiggle effect from the impact
					ball.triggerJiggle(0.3) // Smaller jiggle than wall bounces
				}

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
