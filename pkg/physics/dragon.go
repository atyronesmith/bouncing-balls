package physics

import (
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Dragon represents a dragon that chases the biggest ball but never catches it
type Dragon struct {
	X, Y          float32   // current position
	VX, VY        float32   // velocity for bouncing and drifting
	Size          float32   // size of the dragon
	Speed         float32   // movement speed
	MinDistance   float32   // minimum distance to maintain from target
	ChaseDistance float32   // distance at which dragon starts chasing
	Bounds        fyne.Size // movement bounds
	IsActive      bool      // whether the dragon is active
	// Collision and drift state
	IsDrifting    bool // whether dragon is currently drifting after collision
	DriftTimer    int  // frames remaining in drift mode
	DriftDuration int  // total frames to drift
	// Spinning animation state
	IsSpinning bool    // whether dragon is spinning while resuming chase
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

// NewDragon creates a new dragon figure
func NewDragon(x, y, size float32) *Dragon {
	dragon := &Dragon{
		X:             x,
		Y:             y,
		VX:            0,
		VY:            0,
		Size:          size,
		Speed:         3.0,   // Fast enough to chase but not too fast
		MinDistance:   60.0,  // Never gets closer than this to the target
		ChaseDistance: 400.0, // Increased from 200.0 - Starts chasing when within this distance
		Bounds:        fyne.NewSize(800, 600),
		IsActive:      true,
		DriftDuration: 120, // 2 seconds at 60 FPS
		SpinTarget:    4,   // Spin 4 times when resuming chase
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

// FindBiggestBall finds the ball with the largest radius
func (d *Dragon) FindBiggestBall(balls []*Ball) *Ball {
	if len(balls) == 0 {
		return nil
	}

	var biggestBall *Ball
	maxRadius := float32(0)

	for _, ball := range balls {
		if ball.Radius > maxRadius {
			maxRadius = ball.Radius
			biggestBall = ball
		}
	}

	return biggestBall
}

// CheckCollisionWithBalls checks if dragon collides with any ball
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
		collisionDistance := d.Size*0.4 + ball.Radius // Dragon collision radius is smaller than visual size
		if distance < collisionDistance {
			return ball
		}
	}

	return nil
}

// HandleBallCollision handles collision with a ball and starts drifting
func (d *Dragon) HandleBallCollision(ball *Ball) {
	// Calculate collision direction
	dx := d.X - ball.X
	dy := d.Y - ball.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance > 0 {
		// Normalize collision direction
		normalizedDx := dx / distance
		normalizedDy := dy / distance

		// Set drift velocity (bounce away from ball)
		bounceSpeed := d.Speed * 1.5 // Bounce away faster than normal movement
		d.VX = normalizedDx * bounceSpeed
		d.VY = normalizedDy * bounceSpeed

		// Enter drift mode
		d.IsDrifting = true
		d.DriftTimer = d.DriftDuration
		d.IsSpinning = false
		d.SpinAngle = 0
		d.SpinCount = 0
	}
}

// Update handles all dragon behavior including drifting, spinning, and chasing
func (d *Dragon) Update(balls []*Ball) {
	if !d.IsActive {
		return
	}

	// Check for collisions with balls (only when not already drifting)
	if !d.IsDrifting {
		if collidedBall := d.CheckCollisionWithBalls(balls); collidedBall != nil {
			d.HandleBallCollision(collidedBall)
		}
	}

	// Handle different behavior states
	if d.IsDrifting {
		d.updateDrifting()
	} else if d.IsSpinning {
		d.updateSpinning()
	} else {
		d.updateChasing(balls)
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
	d.VX *= 0.98
	d.VY *= 0.98

	// Check if drift time is over
	if d.DriftTimer <= 0 {
		d.IsDrifting = false
		d.VX = 0
		d.VY = 0
		// Start spinning animation before resuming chase
		d.IsSpinning = true
		d.SpinAngle = 0
		d.SpinCount = 0
	}
}

// updateSpinning handles spinning animation when resuming chase
func (d *Dragon) updateSpinning() {
	// Spin speed - complete one rotation in about 15 frames
	spinSpeed := float32(2 * math.Pi / 15)
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

// updateChasing handles normal chasing behavior
func (d *Dragon) updateChasing(balls []*Ball) {
	biggestBall := d.FindBiggestBall(balls)
	if biggestBall == nil || !biggestBall.IsAnimated {
		d.VX = 0
		d.VY = 0
		return
	}

	// Calculate distance to the biggest ball
	dx := biggestBall.X - d.X
	dy := biggestBall.Y - d.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Only chase if the ball is within chase distance
	if distance > d.ChaseDistance {
		d.VX = 0
		d.VY = 0
		return
	}

	// Calculate desired movement (maintain minimum distance)
	if distance > d.MinDistance {
		// Move towards the ball
		normalizedDx := dx / distance
		normalizedDy := dy / distance

		// Calculate target position that maintains minimum distance
		targetDistance := d.MinDistance + 10 // Add small buffer
		if distance > targetDistance {
			moveDistance := (distance - targetDistance) * 0.8 // Gradual approach
			if moveDistance > d.Speed {
				moveDistance = d.Speed
			}

			d.VX = normalizedDx * moveDistance
			d.VY = normalizedDy * moveDistance
		} else {
			d.VX = 0
			d.VY = 0
		}
	} else {
		// Too close! Back away while still facing the ball
		normalizedDx := dx / distance
		normalizedDy := dy / distance

		// Move away from the ball
		d.VX = -normalizedDx * d.Speed * 0.5
		d.VY = -normalizedDy * d.Speed * 0.5
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
