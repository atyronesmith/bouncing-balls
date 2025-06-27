package physics

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Dragon represents a dragon that protects the human by following them and deflecting balls
type Dragon struct {
	X, Y          float32   // current position
	VX, VY        float32   // velocity for bouncing and movement
	Size          float32   // size of the dragon
	Mass          float32   // mass for collision physics (twice the largest ball)
	Speed         float32   // movement speed
	FollowDistance float32  // preferred distance to maintain from human
	ProtectRadius  float32  // radius within which dragon will intercept balls
	Bounds        fyne.Size // movement bounds
	IsActive      bool      // whether the dragon is active
	// Human movement tracking for strategic deflection
	LastHumanX    float32 // previous human X position
	LastHumanY    float32 // previous human Y position
	HumanVX       float32 // calculated human velocity X
	HumanVY       float32 // calculated human velocity Y
	// Collision and drift state
	IsDrifting    bool // whether dragon is currently drifting after collision
	DriftTimer    int  // frames remaining in drift mode
	DriftDuration int  // total frames to drift
	// Spinning animation state
	IsSpinning bool    // whether dragon is spinning while resuming follow
	SpinAngle  float32 // current spin angle
	SpinCount  int     // number of spins completed
	SpinTarget int     // target number of spins (4)
	// Visual components
	Head      *canvas.Circle
	Body      *canvas.Rectangle
	Tail      *canvas.Rectangle
	LeftWing  *canvas.Rectangle
	RightWing *canvas.Rectangle
	LeftEye   *canvas.Circle
	RightEye  *canvas.Circle
	// Animation state
	WingFlap       float32 // wing flapping animation
	FlameParticles []*canvas.Circle
	FlameTimer     int
}

// NewDragon creates a new dragon figure that protects the human
func NewDragon(x, y, size float32) *Dragon {
	dragon := &Dragon{
		X:             x,
		Y:             y,
		VX:            0,
		VY:            0,
		Size:          size,
		Mass:          70.0,  // Will be updated to twice the largest ball mass
		Speed:         2.0,   // Slower, more controlled movement
		FollowDistance: 80.0, // Preferred distance from human
		ProtectRadius:  150.0, // Will intercept balls within this radius of human
		Bounds:        fyne.NewSize(800, 600),
		IsActive:      true,
		// Initialize human tracking
		LastHumanX:    x, // Initialize with dragon's starting position
		LastHumanY:    y,
		HumanVX:       0,
		HumanVY:       0,
		DriftDuration: 60, // Shorter drift duration for more responsive protection
		SpinTarget:    2,  // Fewer spins for faster recovery
	}

	// Dragon colors
	dragonColor := color.RGBA{R: 150, G: 50, B: 200, A: 255} // Purple dragon
	wingColor := color.RGBA{R: 100, G: 30, B: 150, A: 255}   // Darker purple wings
	eyeColor := color.RGBA{R: 255, G: 255, B: 0, A: 255}     // Yellow eyes
	flameColor := color.RGBA{R: 255, G: 100, B: 50, A: 255}  // Orange flames

	// Head (circle)
	dragon.Head = &canvas.Circle{
		FillColor:   dragonColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 2.0,
	}
	dragon.Head.Resize(fyne.NewSize(size*0.5, size*0.5))

	// Body (ellipse)
	dragon.Body = &canvas.Rectangle{
		FillColor:   dragonColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 2.0,
	}
	dragon.Body.Resize(fyne.NewSize(size*0.8, size*0.4))

	// Tail (rectangle)
	dragon.Tail = &canvas.Rectangle{
		FillColor:   dragonColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	dragon.Tail.Resize(fyne.NewSize(size*0.6, size*0.2))

	// Wings (rectangles)
	dragon.LeftWing = &canvas.Rectangle{
		FillColor:   wingColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	dragon.LeftWing.Resize(fyne.NewSize(size*0.4, size*0.6))

	dragon.RightWing = &canvas.Rectangle{
		FillColor:   wingColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	dragon.RightWing.Resize(fyne.NewSize(size*0.4, size*0.6))

	// Eyes (small circles)
	dragon.LeftEye = &canvas.Circle{
		FillColor:   eyeColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	dragon.LeftEye.Resize(fyne.NewSize(size*0.1, size*0.1))

	dragon.RightEye = &canvas.Circle{
		FillColor:   eyeColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		StrokeWidth: 1.0,
	}
	dragon.RightEye.Resize(fyne.NewSize(size*0.1, size*0.1))

	// Initialize flame particles
	dragon.FlameParticles = make([]*canvas.Circle, 8)
	for i := range dragon.FlameParticles {
		flame := &canvas.Circle{
			FillColor:   flameColor,
			StrokeColor: color.RGBA{R: 255, G: 50, B: 0, A: 255},
			StrokeWidth: 1.0,
		}
		flameSize := size * 0.15 * (1.0 - float32(i)*0.1)
		flame.Resize(fyne.NewSize(flameSize, flameSize))
		dragon.FlameParticles[i] = flame
	}

	// Set initial position
	dragon.UpdatePosition()

	return dragon
}

// FindLargestBall finds the ball with the largest radius to calculate dragon mass
func (d *Dragon) FindLargestBall(balls []*Ball) *Ball {
	if len(balls) == 0 {
		return nil
	}

	var largestBall *Ball
	maxRadius := float32(0)

	for _, ball := range balls {
		if ball.Radius > maxRadius {
			maxRadius = ball.Radius
			largestBall = ball
		}
	}

	return largestBall
}

// UpdateMass updates dragon mass to be twice the largest ball's mass
func (d *Dragon) UpdateMass(balls []*Ball) {
	largestBall := d.FindLargestBall(balls)
	if largestBall != nil {
		// Calculate ball mass based on area (π * r²) and assume unit density
		ballMass := math.Pi * float64(largestBall.Radius) * float64(largestBall.Radius)
		d.Mass = float32(ballMass) * 2.0 // Dragon has twice the mass of largest ball

		// Ensure minimum mass for dragon effectiveness
		minMass := float32(1000.0)
		if d.Mass < minMass {
			d.Mass = minMass
		}
	}
}

// updateHumanVelocity tracks human movement to predict movement direction
func (d *Dragon) updateHumanVelocity(human *Human) {
	if human == nil || !human.IsActive {
		d.HumanVX = 0
		d.HumanVY = 0
		return
	}

	// Calculate human velocity from position change
	d.HumanVX = human.X - d.LastHumanX
	d.HumanVY = human.Y - d.LastHumanY

	// Store current position for next frame
	d.LastHumanX = human.X
	d.LastHumanY = human.Y
}

// FindThreateningBalls finds balls that are moving toward the human and within protect radius
func (d *Dragon) FindThreateningBalls(balls []*Ball, human *Human) []*Ball {
	if human == nil || !human.IsActive {
		return nil
	}

	var threateningBalls []*Ball

	for _, ball := range balls {
		if !ball.IsAnimated {
			continue
		}

		// Calculate distance from ball to human
		dx := human.X - ball.X
		dy := human.Y - ball.Y
		distanceToHuman := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		// Only consider balls within protect radius
		if distanceToHuman > d.ProtectRadius {
			continue
		}

		// Check if ball is moving toward human
		// Calculate dot product of ball velocity and direction to human
		dotProduct := ball.VX*dx + ball.VY*dy
		if dotProduct > 0 { // Ball is moving toward human
			threateningBalls = append(threateningBalls, ball)
		}
	}

	return threateningBalls
}

// FindClosestThreat finds the closest threatening ball to intercept
func (d *Dragon) FindClosestThreat(threateningBalls []*Ball) *Ball {
	if len(threateningBalls) == 0 {
		return nil
	}

	var closestBall *Ball
	minDistance := float32(math.Inf(1))

	for _, ball := range threateningBalls {
		dx := ball.X - d.X
		dy := ball.Y - d.Y
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		if distance < minDistance {
			minDistance = distance
			closestBall = ball
		}
	}

	return closestBall
}

// CheckCollisionWithBalls checks if dragon collides with any ball and handles deflection
func (d *Dragon) CheckCollisionWithBalls(balls []*Ball) *Ball {
	if !d.IsActive {
		return nil
	}

	for _, ball := range balls {
		if !ball.IsAnimated {
			continue
		}

		// Calculate distance between dragon center and ball center
		dx := d.X - ball.X
		dy := d.Y - ball.Y
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		// Check if collision occurs (dragon size/2 + ball radius)
		collisionDistance := d.Size*0.4 + ball.Radius
		if distance < collisionDistance {
			return ball
		}
	}

	return nil
}

// HandleBallCollision handles collision with a ball using mass-based physics to deflect it
func (d *Dragon) HandleBallCollision(ball *Ball, human *Human) {
	// Calculate collision direction
	dx := d.X - ball.X
	dy := d.Y - ball.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance > 0 && human != nil {
		// Normalize collision direction
		normalizedDx := dx / distance
		normalizedDy := dy / distance

		// Calculate ball mass (assuming unit density)
		ballMass := math.Pi * float64(ball.Radius) * float64(ball.Radius)

				// Calculate collision impulse
		relativeVX := d.VX - ball.VX
		relativeVY := d.VY - ball.VY
		relativeSpeed := relativeVX*normalizedDx + relativeVY*normalizedDy

		// Collision impulse magnitude
		impulse := 2.0 * relativeSpeed / (1.0 + float32(ballMass)/d.Mass)

		// Apply impulse to ball (dragon deflects ball strategically)
		// Calculate direction opposite to human's movement for strategic deflection
		var deflectDx, deflectDy float32

		// Check if human is moving (velocity magnitude > threshold)
		humanSpeed := float32(math.Sqrt(float64(d.HumanVX*d.HumanVX + d.HumanVY*d.HumanVY)))

		if humanSpeed > 0.5 { // Human is moving significantly
			// Deflect ball in opposite direction of human movement
			deflectDx = -d.HumanVX / humanSpeed // Opposite of human movement X
			deflectDy = -d.HumanVY / humanSpeed // Opposite of human movement Y
		} else {
			// Human is stationary or moving slowly, deflect away from human
			humanDx := ball.X - human.X
			humanDy := ball.Y - human.Y
			humanDistance := float32(math.Sqrt(float64(humanDx*humanDx + humanDy*humanDy)))

			if humanDistance > 0 {
				deflectDx = humanDx / humanDistance
				deflectDy = humanDy / humanDistance
			}
		}

		// Combine collision direction with strategic deflection direction
		finalDx := (normalizedDx + deflectDx) * 0.5
		finalDy := (normalizedDy + deflectDy) * 0.5
		finalDistance := float32(math.Sqrt(float64(finalDx*finalDx + finalDy*finalDy)))

		if finalDistance > 0 {
			finalDx /= finalDistance
			finalDy /= finalDistance

			// Apply deflection to ball (reduced strength for more control)
			deflectionStrength := impulse * 0.8 // Moderate deflection
			ball.VX += finalDx * deflectionStrength
			ball.VY += finalDy * deflectionStrength
		}

		// Shrink the ball by half and trigger jiggle effect
		ball.shrinkBall(0.5) // Shrink to half size (including mass)
		ball.triggerJiggle(0.5) // Add satisfying jiggle effect

		// Reduce ball's velocity by half to make it less threatening
		ball.VX *= 0.5 // Half the X velocity
		ball.VY *= 0.5 // Half the Y velocity

		// Dragon bounces back slightly (less due to higher mass)
		dragonBounce := impulse * float32(ballMass) / d.Mass * 0.3
		d.VX -= normalizedDx * dragonBounce
		d.VY -= normalizedDy * dragonBounce

		// Enter brief drift mode
		d.IsDrifting = true
		d.DriftTimer = d.DriftDuration / 2 // Shorter drift for responsiveness
		d.IsSpinning = false
		d.SpinAngle = 0
		d.SpinCount = 0
	}
}

// Update handles all dragon behavior including following human and protecting them
func (d *Dragon) Update(balls []*Ball, human *Human) {
	if !d.IsActive {
		return
	}

	// Update human movement tracking for strategic deflection
	d.updateHumanVelocity(human)

	// Update mass based on current largest ball
	d.UpdateMass(balls)

	// Check for collisions with balls (only when not already drifting)
	if !d.IsDrifting {
		if collidedBall := d.CheckCollisionWithBalls(balls); collidedBall != nil {
			d.HandleBallCollision(collidedBall, human)
		}
	}

	// Handle different behavior states
	if d.IsDrifting {
		d.updateDrifting()
	} else if d.IsSpinning {
		d.updateSpinning()
	} else {
		d.updateProtecting(balls, human)
	}

	// Apply movement
	d.X += d.VX
	d.Y += d.VY

	// Keep dragon within bounds
	d.keepWithinBounds()

	// Update animations
	d.updateAnimations()
}

// updateDrifting handles drifting behavior after collision
func (d *Dragon) updateDrifting() {
	d.DriftTimer--

	// Apply friction to slow down drifting
	d.VX *= 0.95 // Less friction for more responsive recovery
	d.VY *= 0.95

	// Check if drift time is over
	if d.DriftTimer <= 0 {
		d.IsDrifting = false
		d.VX = 0
		d.VY = 0
		// Start brief spinning animation before resuming protection
		d.IsSpinning = true
		d.SpinAngle = 0
		d.SpinCount = 0
	}
}

// updateSpinning handles spinning animation when resuming protection
func (d *Dragon) updateSpinning() {
	// Faster spin speed for quicker recovery
	spinSpeed := float32(2 * math.Pi / 10)
	d.SpinAngle += spinSpeed

	// Check if completed a full rotation
	if d.SpinAngle >= 2*math.Pi {
		d.SpinAngle -= 2 * math.Pi
		d.SpinCount++

		// Check if completed target number of spins
		if d.SpinCount >= d.SpinTarget {
			d.IsSpinning = false
			d.SpinAngle = 0
			d.SpinCount = 0
		}
	}
}

// updateProtecting handles protective behavior - following human and intercepting threats
func (d *Dragon) updateProtecting(balls []*Ball, human *Human) {
	if human == nil || !human.IsActive {
		d.VX = 0
		d.VY = 0
		return
	}

	// Find threatening balls
	threateningBalls := d.FindThreateningBalls(balls, human)
	closestThreat := d.FindClosestThreat(threateningBalls)

	if closestThreat != nil {
		// Priority: Intercept the closest threat
		d.interceptBall(closestThreat, human)
	} else {
		// No immediate threats: Follow the human at preferred distance
		d.followHuman(human)
	}
}

// interceptBall moves dragon to intercept a threatening ball
func (d *Dragon) interceptBall(ball *Ball, human *Human) {
	// Calculate interception point
	// Predict where the ball will be when dragon reaches it
	dx := ball.X - d.X
	dy := ball.Y - d.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance > 0 {
		// Simple interception: move toward ball's current position aggressively
		normalizedDx := dx / distance
		normalizedDy := dy / distance

		// Move at controlled speed toward threat
		d.VX = normalizedDx * d.Speed * 1.2 // 1.2x speed when intercepting (was 1.5x)
		d.VY = normalizedDy * d.Speed * 1.2
	}
}

// followHuman makes dragon follow the human at preferred distance
func (d *Dragon) followHuman(human *Human) {
	// Calculate distance to human
	dx := human.X - d.X
	dy := human.Y - d.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance > 0 {
		normalizedDx := dx / distance
		normalizedDy := dy / distance

		if distance > d.FollowDistance+20 {
			// Too far: Move closer to human
			moveSpeed := (distance - d.FollowDistance) * 0.1 // Gradual approach
			if moveSpeed > d.Speed {
				moveSpeed = d.Speed
			}
			d.VX = normalizedDx * moveSpeed
			d.VY = normalizedDy * moveSpeed
		} else if distance < d.FollowDistance-20 {
			// Too close: Back away slightly
			d.VX = -normalizedDx * d.Speed * 0.3
			d.VY = -normalizedDy * d.Speed * 0.3
		} else {
			// Good distance: Maintain position with slight drift
			d.VX *= 0.9
			d.VY *= 0.9
		}
	}
}

// keepWithinBounds ensures dragon stays within window boundaries
func (d *Dragon) keepWithinBounds() {
	margin := d.Size * 0.5
	if d.X < margin {
		d.X = margin
		if d.VX < 0 {
			d.VX = -d.VX * 0.5 // Bounce off left wall
		}
	} else if d.X > d.Bounds.Width-margin {
		d.X = d.Bounds.Width - margin
		if d.VX > 0 {
			d.VX = -d.VX * 0.5 // Bounce off right wall
		}
	}
	if d.Y < 50+margin {
		d.Y = 50 + margin
		if d.VY < 0 {
			d.VY = -d.VY * 0.5 // Bounce off top wall
		}
	} else if d.Y > d.Bounds.Height-50-margin {
		d.Y = d.Bounds.Height - 50 - margin
		if d.VY > 0 {
			d.VY = -d.VY * 0.5 // Bounce off bottom wall
		}
	}
}

// updateAnimations handles wing flapping and flame effects
func (d *Dragon) updateAnimations() {
	// Update wing flap animation (faster when spinning)
	flapSpeed := float32(0.3)
	if d.IsSpinning {
		flapSpeed = 0.8 // Faster wing flapping during spin
	}
	d.WingFlap += flapSpeed
	if d.WingFlap > 2*math.Pi {
		d.WingFlap = 0
	}

	// Update flame animation
	d.FlameTimer++
	if d.FlameTimer > 60 {
		d.FlameTimer = 0
	}
}

// UpdatePosition updates the visual position of all dragon components
func (d *Dragon) UpdatePosition() {
	if !d.IsActive {
		return
	}

	// Wing flap effect
	wingOffset := float32(math.Sin(float64(d.WingFlap))) * 5

	// Spinning effect - rotate all components around dragon center
	var spinCos, spinSin float32 = 1, 0
	if d.IsSpinning {
		spinCos = float32(math.Cos(float64(d.SpinAngle)))
		spinSin = float32(math.Sin(float64(d.SpinAngle)))
	}

	// Helper function to apply spin rotation to a position offset
	applySpinRotation := func(offsetX, offsetY float32) (float32, float32) {
		if !d.IsSpinning {
			return offsetX, offsetY
		}
		rotatedX := offsetX*spinCos - offsetY*spinSin
		rotatedY := offsetX*spinSin + offsetY*spinCos
		return rotatedX, rotatedY
	}

	// Head position (front center)
	headOffsetX, headOffsetY := applySpinRotation(-d.Size*0.25, -d.Size*0.25)
	d.Head.Move(fyne.NewPos(d.X+headOffsetX, d.Y+headOffsetY))

	// Body position (center)
	bodyOffsetX, bodyOffsetY := applySpinRotation(-d.Size*0.4, -d.Size*0.2)
	d.Body.Move(fyne.NewPos(d.X+bodyOffsetX, d.Y+bodyOffsetY))

	// Tail position (behind body)
	tailOffsetX, tailOffsetY := applySpinRotation(-d.Size*0.8, -d.Size*0.1)
	d.Tail.Move(fyne.NewPos(d.X+tailOffsetX, d.Y+tailOffsetY))

	// Wings positions (animated flapping + spinning)
	leftWingOffsetX, leftWingOffsetY := applySpinRotation(-d.Size*0.6, -d.Size*0.5+wingOffset)
	d.LeftWing.Move(fyne.NewPos(d.X+leftWingOffsetX, d.Y+leftWingOffsetY))

	rightWingOffsetX, rightWingOffsetY := applySpinRotation(d.Size*0.2, -d.Size*0.5-wingOffset)
	d.RightWing.Move(fyne.NewPos(d.X+rightWingOffsetX, d.Y+rightWingOffsetY))

	// Eyes positions (on head)
	leftEyeOffsetX, leftEyeOffsetY := applySpinRotation(-d.Size*0.15, -d.Size*0.3)
	d.LeftEye.Move(fyne.NewPos(d.X+leftEyeOffsetX, d.Y+leftEyeOffsetY))

	rightEyeOffsetX, rightEyeOffsetY := applySpinRotation(-d.Size*0.05, -d.Size*0.3)
	d.RightEye.Move(fyne.NewPos(d.X+rightEyeOffsetX, d.Y+rightEyeOffsetY))

	// Update flame particles (breath effect + spinning)
	for i, flame := range d.FlameParticles {
		if flame != nil {
			// Base flame position
			baseFlameX := -d.Size*0.5 - float32(i)*8
			baseFlameY := float32(math.Sin(float64(d.FlameTimer)*0.1+float64(i)*0.5)) * 3

			// Apply spin rotation to flame position
			flameOffsetX, flameOffsetY := applySpinRotation(baseFlameX, baseFlameY)

			flameSize := d.Size * 0.15 * (1.0 - float32(i)*0.1)
			flame.Move(fyne.NewPos(d.X+flameOffsetX-flameSize/2, d.Y+flameOffsetY-flameSize/2))

			// Animate flame color (more intense during spinning)
			alpha := uint8(200 - i*20)
			red := uint8(255)
			green := uint8(100 + i*10)
			if d.IsSpinning {
				green = uint8(150 + i*10) // More intense flames during spin
			}
			flame.FillColor = color.RGBA{R: red, G: green, B: 50, A: alpha}
			flame.Refresh()
		}
	}
}

// GetVisualComponents returns all visual components for adding to container
func (d *Dragon) GetVisualComponents() []fyne.CanvasObject {
	components := []fyne.CanvasObject{
		d.Tail,     // Draw tail first (behind)
		d.LeftWing, // Wings behind body
		d.RightWing,
		d.Body,    // Body in middle
		d.Head,    // Head on top
		d.LeftEye, // Eyes on top of head
		d.RightEye,
	}

	// Add flame particles
	for _, flame := range d.FlameParticles {
		if flame != nil {
			components = append(components, flame)
		}
	}

	return components
}

// Hide hides all dragon components
func (d *Dragon) Hide() {
	d.Head.Hide()
	d.Body.Hide()
	d.Tail.Hide()
	d.LeftWing.Hide()
	d.RightWing.Hide()
	d.LeftEye.Hide()
	d.RightEye.Hide()
	for _, flame := range d.FlameParticles {
		if flame != nil {
			flame.Hide()
		}
	}
}

// Show shows all dragon components
func (d *Dragon) Show() {
	d.Head.Show()
	d.Body.Show()
	d.Tail.Show()
	d.LeftWing.Show()
	d.RightWing.Show()
	d.LeftEye.Show()
	d.RightEye.Show()
	for _, flame := range d.FlameParticles {
		if flame != nil {
			flame.Show()
		}
	}
}
