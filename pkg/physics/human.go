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
	// Eyeball components for bullets
	Eyeball  *canvas.Circle  // White eyeball
	Iris     *canvas.Circle  // Colored iris
	Pupil    *canvas.Circle  // Black pupil
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
	// Rotation field for facing direction calculations
	Rotation       float64           // Current rotation angle in radians
	// Firing circle system
	FiringCircle   *canvas.Circle    // Transparent circle around human for bullet origin
	FiringEye      *canvas.Circle    // Outer eyeball (white) that shows when firing
	FiringIris     *canvas.Circle    // Colored iris (red/orange) inside the eye
	FiringPupil    *canvas.Circle    // Black pupil in center of eye
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

	// Create firing circle system
	human.FiringRadius = size * 1.5 // Circle radius around human

	// Transparent firing circle with only edge visible
	human.FiringCircle = &canvas.Circle{
		FillColor:   color.RGBA{R: 0, G: 0, B: 0, A: 0},           // Completely transparent interior
		StrokeColor: color.RGBA{R: 0, G: 150, B: 255, A: 3},       // 99% transparent blue edge (almost invisible)
		StrokeWidth: 1.0,
	}
	human.FiringCircle.Resize(fyne.NewSize(human.FiringRadius*2, human.FiringRadius*2))
	human.FiringCircle.Move(fyne.NewPos(x-human.FiringRadius, y-human.FiringRadius))

	// Firing effect (highlights a portion of the circle edge when shooting)
	human.FiringEye = &canvas.Circle{
		FillColor:   color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White eyeball
		StrokeColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},     // Red outline for menacing look
		StrokeWidth: 3.0,
	}
	eyeballSize := human.FiringRadius * 0.267  // Smaller eye size (1/3 of original)
	human.FiringEye.Resize(fyne.NewSize(eyeballSize, eyeballSize))
	human.FiringEye.Move(fyne.NewPos(x-human.FiringRadius, y-human.FiringRadius))

	// Firing iris - make it more vibrant
	human.FiringIris = &canvas.Circle{
		FillColor:   color.RGBA{R: 255, G: 50, B: 0, A: 255},     // Bright orange-red iris
		StrokeColor: color.RGBA{R: 150, G: 0, B: 0, A: 255},      // Dark red outline
		StrokeWidth: 2.0,
	}
	irisSize := eyeballSize * 0.7 // Larger iris proportion
	human.FiringIris.Resize(fyne.NewSize(irisSize, irisSize))
	human.FiringIris.Move(fyne.NewPos(x-human.FiringRadius, y-human.FiringRadius))

	// Firing pupil - make it more prominent
	human.FiringPupil = &canvas.Circle{
		FillColor:   color.RGBA{R: 0, G: 0, B: 0, A: 255},         // Black pupil
		StrokeColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},   // White outline for contrast
		StrokeWidth: 2.0,
	}
	pupilSize := eyeballSize * 0.35 // Larger pupil for more intensity
	human.FiringPupil.Resize(fyne.NewSize(pupilSize, pupilSize))
	human.FiringPupil.Move(fyne.NewPos(x-human.FiringRadius, y-human.FiringRadius))

	// Hide eye components initially
	human.FiringEye.Hide()
	human.FiringIris.Hide()
	human.FiringPupil.Hide()

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

	// Update firing circle position
	h.updateFiringCircle()
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
	h.FiringCircle.Hide()
	h.FiringEye.Hide()
	h.FiringIris.Hide()
	h.FiringPupil.Hide()

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
	h.FiringCircle.Show()
	h.FiringEye.Show()
	h.FiringIris.Show()
	h.FiringPupil.Show()

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
	h.FiringCircle.Show()
	h.FiringEye.Show()
	h.FiringIris.Show()
	h.FiringPupil.Show()

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
func (h *Human) UpdatePointing(_ []*Ball) {
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

	// Create eyeball bullet components
	bulletSize := float32(20) // Size of the bullet eyeball

	// White eyeball (outer layer)
	bullet.Eyeball = &canvas.Circle{
		FillColor:   color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White eyeball
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},       // Black outline
		StrokeWidth: 2.0,
	}
	bullet.Eyeball.Resize(fyne.NewSize(bulletSize, bulletSize))
	bullet.Eyeball.Move(fyne.NewPos(startX-bulletSize/2, startY-bulletSize/2))

	// Colored iris (middle layer)
	irisSize := bulletSize * 0.7
	bullet.Iris = &canvas.Circle{
		FillColor:   color.RGBA{R: 0, G: 255, B: 255, A: 255},   // Cyan iris for visibility
		StrokeColor: color.RGBA{R: 0, G: 150, B: 150, A: 255},   // Darker cyan outline
		StrokeWidth: 1.0,
	}
	bullet.Iris.Resize(fyne.NewSize(irisSize, irisSize))
	bullet.Iris.Move(fyne.NewPos(startX-irisSize/2, startY-irisSize/2))

	// Black pupil (inner layer)
	pupilSize := bulletSize * 0.35
	bullet.Pupil = &canvas.Circle{
		FillColor:   color.RGBA{R: 0, G: 0, B: 0, A: 255},       // Black pupil
		StrokeColor: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White outline
		StrokeWidth: 1.0,
	}
	bullet.Pupil.Resize(fyne.NewSize(pupilSize, pupilSize))
	bullet.Pupil.Move(fyne.NewPos(startX-pupilSize/2, startY-pupilSize/2))

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
		bulletSize := float32(20) // Size of the bullet eyeball (match NewBullet)
		irisSize := bulletSize * 0.7
		pupilSize := bulletSize * 0.35

		bullet.X += bullet.VX
		bullet.Y += bullet.VY
		bullet.Eyeball.Move(fyne.NewPos(bullet.X-bulletSize/2, bullet.Y-bulletSize/2))
		bullet.Iris.Move(fyne.NewPos(bullet.X-irisSize/2, bullet.Y-irisSize/2))
		bullet.Pupil.Move(fyne.NewPos(bullet.X-pupilSize/2, bullet.Y-pupilSize/2))

		// Remove bullets that go off screen
		if bullet.X < 0 || bullet.X > h.Bounds.Width || bullet.Y < 0 || bullet.Y > h.Bounds.Height {
			bullet.IsActive = false
			bullet.Eyeball.Hide()
			bullet.Iris.Hide()
			bullet.Pupil.Hide()
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

			if distance < ball.Radius+10 { // bullet radius is now 10 (bulletSize/2 = 20/2 = 10)
				// Bullet hit ball!
				bullet.IsActive = false
				bullet.Eyeball.Hide()
				bullet.Iris.Hide()
				bullet.Pupil.Hide()

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
			visuals = append(visuals, bullet.Eyeball)
			visuals = append(visuals, bullet.Iris)
			visuals = append(visuals, bullet.Pupil)
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

		// Position firing effect to appear as a highlighted segment on the main circle edge
		// Position effect center on the main circle's edge at the firing angle
		effectCenterX := h.X + float32(math.Cos(float64(h.FiringAngle))) * h.FiringRadius
		effectCenterY := h.Y + float32(math.Sin(float64(h.FiringAngle))) * h.FiringRadius

		// Position eyeball at the calculated position (use consistent sizing)
		eyeballSize := h.FiringRadius * 0.267   // Smaller eye size (1/3 of original)
		irisSize := eyeballSize * 0.7         // Larger iris proportion
		pupilSize := eyeballSize * 0.35       // Larger pupil for more intensity

		eyeballRadius := eyeballSize / 2      // Convert size to radius for positioning
		irisRadius := irisSize / 2
		pupilRadius := pupilSize / 2

		// Position eyeball
		h.FiringEye.Move(fyne.NewPos(effectCenterX-eyeballRadius, effectCenterY-eyeballRadius))
		h.FiringEye.Show()

		// Position iris centered within eyeball
		h.FiringIris.Move(fyne.NewPos(effectCenterX-irisRadius, effectCenterY-irisRadius))
		h.FiringIris.Show()

		// Position pupil centered within iris
		h.FiringPupil.Move(fyne.NewPos(effectCenterX-pupilRadius, effectCenterY-pupilRadius))
		h.FiringPupil.Show()
	} else {
		h.FiringEye.Hide()
		h.FiringIris.Hide()
		h.FiringPupil.Hide()
	}
}
