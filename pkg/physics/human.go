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
	// Visual components - drawn programmatically
	Head           *canvas.Circle    // Head (circle)
	Body           *canvas.Rectangle // Body (rectangle)
	LeftEye        *canvas.Circle    // Left eye
	RightEye       *canvas.Circle    // Right eye
	LeftPupil      *canvas.Circle    // Left pupil (tracks closest ball)
	RightPupil     *canvas.Circle    // Right pupil (tracks closest ball)
	LeftArm        *canvas.Rectangle // Left arm
	RightArm       *canvas.Rectangle // Right arm
	LeftLeg        *canvas.Rectangle // Left leg
	RightLeg       *canvas.Rectangle // Right leg
	DirectionArrow *canvas.Line      // Arrow showing facing direction (deprecated)
	Rotation       float64           // Current rotation angle in radians
	// Firing circle system
	FiringCircle   *canvas.Circle    // Transparent circle around human for bullet origin
	FiringEffect   *canvas.Circle    // Highlighted arc/segment that shows when firing
	FiringAngle    float32           // Current angle where bullets are fired from
	FiringEffectTimer int            // Timer for showing firing effect
	FiringRadius   float32           // Radius of the firing circle
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

	// Create head (circle)
	human.Head = &canvas.Circle{
		FillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White head
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Black outline
		StrokeWidth: 2.0,
	}
	human.Head.Resize(fyne.NewSize(size*0.8, size*0.8))
	human.Head.Move(fyne.NewPos(x-size*0.4, y-size*0.4))

	// Create body (rectangle)
	human.Body = &canvas.Rectangle{
		FillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White body
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Black outline
		StrokeWidth: 2.0,
	}
	human.Body.Resize(fyne.NewSize(size*0.6, size*0.4))
	human.Body.Move(fyne.NewPos(x-size*0.3, y-size*0.2))

	// Create eyes
	human.LeftEye = &canvas.Circle{
		FillColor: color.RGBA{R: 0, G: 0, B: 0, A: 255}, // Black eye
		StrokeColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		StrokeWidth: 2.0,
	}
	human.LeftEye.Resize(fyne.NewSize(size*0.1, size*0.1))
	human.LeftEye.Move(fyne.NewPos(x-size*0.35, y-size*0.3))

	human.RightEye = &canvas.Circle{
		FillColor: color.RGBA{R: 0, G: 0, B: 0, A: 255}, // Black eye
		StrokeColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		StrokeWidth: 2.0,
	}
	human.RightEye.Resize(fyne.NewSize(size*0.1, size*0.1))
	human.RightEye.Move(fyne.NewPos(x+size*0.35, y-size*0.3))

	// Create pupils
	human.LeftPupil = &canvas.Circle{
		FillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White pupil
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 2.0,
	}
	human.LeftPupil.Resize(fyne.NewSize(size*0.05, size*0.05))
	human.LeftPupil.Move(fyne.NewPos(x-size*0.35, y-size*0.3))

	human.RightPupil = &canvas.Circle{
		FillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White pupil
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 2.0,
	}
	human.RightPupil.Resize(fyne.NewSize(size*0.05, size*0.05))
	human.RightPupil.Move(fyne.NewPos(x+size*0.35, y-size*0.3))

	// Create arms
	human.LeftArm = &canvas.Rectangle{
		FillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White arm
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Black outline
		StrokeWidth: 2.0,
	}
	human.LeftArm.Resize(fyne.NewSize(size*0.2, size*0.2))
	human.LeftArm.Move(fyne.NewPos(x-size*0.3, y-size*0.2))

	human.RightArm = &canvas.Rectangle{
		FillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White arm
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Black outline
		StrokeWidth: 2.0,
	}
	human.RightArm.Resize(fyne.NewSize(size*0.2, size*0.2))
	human.RightArm.Move(fyne.NewPos(x+size*0.3, y-size*0.2))

	// Create legs
	human.LeftLeg = &canvas.Rectangle{
		FillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White leg
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Black outline
		StrokeWidth: 2.0,
	}
	human.LeftLeg.Resize(fyne.NewSize(size*0.2, size*0.4))
	human.LeftLeg.Move(fyne.NewPos(x-size*0.3, y+size*0.2))

	human.RightLeg = &canvas.Rectangle{
		FillColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White leg
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},     // Black outline
		StrokeWidth: 2.0,
	}
	human.RightLeg.Resize(fyne.NewSize(size*0.2, size*0.4))
	human.RightLeg.Move(fyne.NewPos(x+size*0.3, y+size*0.2))

	// Create direction arrow to show facing direction
	human.DirectionArrow = &canvas.Line{
		StrokeColor: color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red arrow
		StrokeWidth: 3.0,
	}

	// Create firing circle system
	human.FiringRadius = size * 1.5 // Circle radius around human

	// Transparent firing circle
	human.FiringCircle = &canvas.Circle{
		FillColor:   color.RGBA{R: 100, G: 200, B: 255, A: 60}, // Light blue, very transparent
		StrokeColor: color.RGBA{R: 0, G: 150, B: 255, A: 120},   // Blue outline, semi-transparent
		StrokeWidth: 2.0,
	}
	human.FiringCircle.Resize(fyne.NewSize(human.FiringRadius*2, human.FiringRadius*2))
	human.FiringCircle.Move(fyne.NewPos(x-human.FiringRadius, y-human.FiringRadius))

	// Firing effect (shows briefly when shooting)
	human.FiringEffect = &canvas.Circle{
		FillColor:   color.RGBA{R: 255, G: 255, B: 100, A: 200}, // Bright yellow flash
		StrokeColor: color.RGBA{R: 255, G: 200, B: 0, A: 255},   // Orange outline
		StrokeWidth: 4.0,
	}
	human.FiringEffect.Resize(fyne.NewSize(size*0.3, size*0.3)) // Small effect
	human.FiringEffect.Hide() // Initially hidden

	// Initialize firing system
	human.FiringAngle = 0
	human.FiringEffectTimer = 0

	// Set initial position
	human.UpdatePosition()

	return human
}

// UpdatePosition updates the visual position of all human components
func (h *Human) UpdatePosition() {
	if !h.IsActive {
		return
	}

	// Update all component positions relative to human center
	// Head
	h.Head.Move(fyne.NewPos(h.X-h.Size*0.4, h.Y-h.Size*0.6))
	h.Head.Resize(fyne.NewSize(h.Size*0.8, h.Size*0.8))

	// Body
	h.Body.Move(fyne.NewPos(h.X-h.Size*0.3, h.Y-h.Size*0.2))
	h.Body.Resize(fyne.NewSize(h.Size*0.6, h.Size*0.8))

	// Eyes (fixed positions on head)
	h.LeftEye.Move(fyne.NewPos(h.X-h.Size*0.25, h.Y-h.Size*0.5))
	h.LeftEye.Resize(fyne.NewSize(h.Size*0.15, h.Size*0.15))

	h.RightEye.Move(fyne.NewPos(h.X+h.Size*0.1, h.Y-h.Size*0.5))
	h.RightEye.Resize(fyne.NewSize(h.Size*0.15, h.Size*0.15))

	// Arms
	h.LeftArm.Move(fyne.NewPos(h.X-h.Size*0.6, h.Y-h.Size*0.1))
	h.LeftArm.Resize(fyne.NewSize(h.Size*0.25, h.Size*0.35))

	h.RightArm.Move(fyne.NewPos(h.X+h.Size*0.35, h.Y-h.Size*0.1))
	h.RightArm.Resize(fyne.NewSize(h.Size*0.25, h.Size*0.35))

	// Legs
	h.LeftLeg.Move(fyne.NewPos(h.X-h.Size*0.25, h.Y+h.Size*0.3))
	h.LeftLeg.Resize(fyne.NewSize(h.Size*0.2, h.Size*0.5))

	h.RightLeg.Move(fyne.NewPos(h.X+h.Size*0.05, h.Y+h.Size*0.3))
	h.RightLeg.Resize(fyne.NewSize(h.Size*0.2, h.Size*0.5))

	// Update direction arrow to show facing direction
	h.updateDirectionArrow()

	// Update firing circle position
	h.updateFiringCircle()
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

// updateEyeTracking makes the pupils look at the closest ball
func (h *Human) updateEyeTracking() {
	if !h.IsActive {
		return
	}

	// Get left and right eye centers
	leftEyeCenterX := h.X - h.Size*0.175 // Center of left eye
	leftEyeCenterY := h.Y - h.Size*0.425

	rightEyeCenterX := h.X + h.Size*0.175 // Center of right eye
	rightEyeCenterY := h.Y - h.Size*0.425

	// Default pupil positions (center of eyes) when no balls
	leftPupilX := leftEyeCenterX
	leftPupilY := leftEyeCenterY
	rightPupilX := rightEyeCenterX
	rightPupilY := rightEyeCenterY

	// Find closest ball for eye tracking (we'll need to pass balls to this method)
	// For now, pupils stay centered - we'll update this when we modify the Update method

	// Position pupils (slightly smaller than eyes)
	pupilSize := h.Size * 0.08
	h.LeftPupil.Move(fyne.NewPos(leftPupilX-pupilSize/2, leftPupilY-pupilSize/2))
	h.LeftPupil.Resize(fyne.NewSize(pupilSize, pupilSize))

	h.RightPupil.Move(fyne.NewPos(rightPupilX-pupilSize/2, rightPupilY-pupilSize/2))
	h.RightPupil.Resize(fyne.NewSize(pupilSize, pupilSize))
}

// updateEyeTrackingWithBalls makes the pupils look at the closest ball
func (h *Human) updateEyeTrackingWithBalls(balls []*Ball) {
	if !h.IsActive {
		return
	}

	// Get left and right eye centers
	leftEyeCenterX := h.X - h.Size*0.175
	leftEyeCenterY := h.Y - h.Size*0.425

	rightEyeCenterX := h.X + h.Size*0.175
	rightEyeCenterY := h.Y - h.Size*0.425

	// Default pupil positions (center of eyes)
	leftPupilX := leftEyeCenterX
	leftPupilY := leftEyeCenterY
	rightPupilX := rightEyeCenterX
	rightPupilY := rightEyeCenterY

	// Find closest ball
	closestBall := h.findClosestBall(balls)
	if closestBall != nil {
		// Calculate direction from each eye to the closest ball
		// Left eye
		leftDx := closestBall.X - leftEyeCenterX
		leftDy := closestBall.Y - leftEyeCenterY
		leftDistance := float32(math.Sqrt(float64(leftDx*leftDx + leftDy*leftDy)))

		if leftDistance > 0 {
			// Normalize direction and scale by eye radius to keep pupil inside eye
			eyeRadius := h.Size * 0.06 // Maximum pupil movement within eye
			leftNormX := leftDx / leftDistance
			leftNormY := leftDy / leftDistance

			leftPupilX = leftEyeCenterX + leftNormX*eyeRadius
			leftPupilY = leftEyeCenterY + leftNormY*eyeRadius
		}

		// Right eye
		rightDx := closestBall.X - rightEyeCenterX
		rightDy := closestBall.Y - rightEyeCenterY
		rightDistance := float32(math.Sqrt(float64(rightDx*rightDx + rightDy*rightDy)))

		if rightDistance > 0 {
			// Normalize direction and scale by eye radius
			eyeRadius := h.Size * 0.06
			rightNormX := rightDx / rightDistance
			rightNormY := rightDy / rightDistance

			rightPupilX = rightEyeCenterX + rightNormX*eyeRadius
			rightPupilY = rightEyeCenterY + rightNormY*eyeRadius
		}
	}

	// Position pupils
	pupilSize := h.Size * 0.08
	h.LeftPupil.Move(fyne.NewPos(leftPupilX-pupilSize/2, leftPupilY-pupilSize/2))
	h.LeftPupil.Resize(fyne.NewSize(pupilSize, pupilSize))

	h.RightPupil.Move(fyne.NewPos(rightPupilX-pupilSize/2, rightPupilY-pupilSize/2))
	h.RightPupil.Resize(fyne.NewSize(pupilSize, pupilSize))
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

	// Update eye tracking to look at closest ball
	h.updateEyeTrackingWithBalls(balls)

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
	h.Head.Hide()
	h.Body.Hide()
	h.LeftEye.Hide()
	h.RightEye.Hide()
	h.LeftPupil.Hide()
	h.RightPupil.Hide()
	h.LeftArm.Hide()
	h.RightArm.Hide()
	h.LeftLeg.Hide()
	h.RightLeg.Hide()
	h.DirectionArrow.Hide()
	h.FiringCircle.Hide()
	h.FiringEffect.Hide()

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
	h.Head.Show()
	h.Body.Show()
	h.LeftEye.Show()
	h.RightEye.Show()
	h.LeftPupil.Show()
	h.RightPupil.Show()
	h.LeftArm.Show()
	h.RightArm.Show()
	h.LeftLeg.Show()
	h.RightLeg.Show()
	h.DirectionArrow.Show()
	h.FiringCircle.Show()
	h.FiringEffect.Show()

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
	h.Head.Show()
	h.Body.Show()
	h.LeftEye.Show()
	h.RightEye.Show()
	h.LeftPupil.Show()
	h.RightPupil.Show()
	h.LeftArm.Show()
	h.RightArm.Show()
	h.LeftLeg.Show()
	h.RightLeg.Show()
	h.DirectionArrow.Show()
	h.FiringCircle.Show()
	h.FiringEffect.Show()

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

// ShootAtTarget creates bullets from the firing circle edge toward the target
func (h *Human) ShootAtTarget(targetX, targetY float32) {
	if !h.IsActive || h.IsExploding {
		return
	}

	// Calculate angle to target
	dx := targetX - h.X
	dy := targetY - h.Y
	h.FiringAngle = float32(math.Atan2(float64(dy), float64(dx)))

	// Calculate bullet spawn position on the firing circle edge
	bulletX := h.X + float32(math.Cos(float64(h.FiringAngle))) * h.FiringRadius
	bulletY := h.Y + float32(math.Sin(float64(h.FiringAngle))) * h.FiringRadius

	// Create bullet from the circle edge position
	bullet := NewBullet(bulletX, bulletY, targetX, targetY)
	h.Bullets = append(h.Bullets, bullet)

	// Trigger firing effect
	h.FiringEffectTimer = 15 // Show effect for 15 frames (quarter second at 60fps)
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

// updateFiringCircle updates the position of the firing circle and effect
func (h *Human) updateFiringCircle() {
	if !h.IsActive {
		return
	}

	// Update firing circle position (centered on human)
	h.FiringCircle.Move(fyne.NewPos(h.X-h.FiringRadius, h.Y-h.FiringRadius))

	// Update firing effect timer and visibility
	if h.FiringEffectTimer > 0 {
		h.FiringEffectTimer--

		// Position firing effect at the firing angle on the circle edge
		effectX := h.X + float32(math.Cos(float64(h.FiringAngle))) * h.FiringRadius
		effectY := h.Y + float32(math.Sin(float64(h.FiringAngle))) * h.FiringRadius

		effectSize := h.Size * 0.3
		h.FiringEffect.Move(fyne.NewPos(effectX-effectSize/2, effectY-effectSize/2))
		h.FiringEffect.Show()
	} else {
		h.FiringEffect.Hide()
	}
}
