package physics

import (
	"image/color"
	"math"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// Ball represents a ball with position and velocity
type Ball struct {
	X, Y       float32 // current position
	VX, VY     float32 // velocity
	Radius     float32 // ball radius
	Circle     *canvas.Circle
	Text       *canvas.Text // AI LLM name label
	LLMName    string       // AI LLM name
	Bounds     fyne.Size    // animation bounds
	IsAnimated bool         // whether animation is running
	// Particle trail system
	Trail      []*canvas.Circle
	TrailIndex int
	// Jiggle effect for jello-like bouncing
	JiggleAmplitude float32 // Current jiggle strength
	JigglePhase     float32 // Current phase of jiggle oscillation
	JiggleDecay     float32 // How fast jiggle fades
	OriginalRadius  float32 // Original radius before jiggle
	// Explosion effects for ball collisions
	ExplosionParticles []*canvas.Circle
	ExplosionTimer     int  // frames for explosion animation
	IsExploding        bool // whether ball is currently exploding
}

// AI LLM names to choose from
var llmNames = []string{
	"GPT-4",
	"Claude",
	"Gemini",
	"LLaMA",
	"PaLM",
	"Bard",
	"ChatGPT",
	"Codex",
	"Alpaca",
	"Vicuna",
	"Mistral",
	"Llama2",
}

// getRandomLLMName returns a random AI LLM name
func getRandomLLMName() string {
	return llmNames[rand.Intn(len(llmNames))]
}

// updateTextSize calculates and sets the appropriate font size for the text to fit inside the ball
func (b *Ball) updateTextSize() {
	if b.Text == nil {
		return
	}

	// Simple font size based on ball radius
	// Larger balls get larger text, smaller balls get smaller text
	fontSize := b.Radius * 0.4 // Scale factor based on radius

	// Ensure reasonable bounds
	if fontSize < 8 {
		fontSize = 8
	} else if fontSize > 20 {
		fontSize = 20
	}
	// fontSize is within bounds, no change needed

	// Apply the font size
	b.Text.TextSize = fontSize
}

// NewBall creates a new bouncing ball
func NewBall() *Ball {
	ball := &Ball{
		X:       100,
		Y:       100,
		VX:      3.5, // horizontal velocity
		VY:      2.8, // vertical velocity
		Radius:  30,
		Bounds:  fyne.NewSize(800, 600),
		LLMName: getRandomLLMName(),
		// Initialize jiggle properties
		JiggleAmplitude: 0.0,
		JigglePhase:     0.0,
		JiggleDecay:     0.88, // Decay rate for jiggle amplitude
		OriginalRadius:  30,
		// Initialize explosion properties
		ExplosionParticles: nil,
		ExplosionTimer:     0,
		IsExploding:        false,
	}

	// Create the visual circle with no stroke
	ball.Circle = &canvas.Circle{
		FillColor:   color.RGBA{R: 100, G: 150, B: 255, A: 255}, // Light blue
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 0},         // No stroke
		StrokeWidth: 0,                                          // No stroke width
	}

	// Create the text label for the AI LLM name
	ball.Text = &canvas.Text{
		Text:      ball.LLMName,
		Color:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White text
		Alignment: fyne.TextAlignCenter,
		TextStyle: fyne.TextStyle{Bold: true},
		TextSize:  12, // Initial size, will be adjusted
	}

	// Set initial size and position
	ball.Circle.Resize(fyne.NewSize(ball.Radius*2, ball.Radius*2))

	// Set font size to fit inside ball
	ball.updateTextSize()

	ball.updatePosition()

	// Initialize particle trail
	ball.initializeTrail()

	return ball
}

// initializeTrail creates the particle trail for the ball
func (b *Ball) initializeTrail() {
	// Clean up existing trail particles first
	if b.Trail != nil {
		for _, trail := range b.Trail {
			if trail != nil {
				trail.Hide()
			}
		}
	}

	trailLength := 10 // Number of trail particles
	b.Trail = make([]*canvas.Circle, trailLength)
	b.TrailIndex = 0 // Reset trail index

	for i := 0; i < trailLength; i++ {
		trail := &canvas.Circle{
			FillColor:   color.RGBA{R: 255, G: 255, B: 255, A: uint8(255 - i*20)}, // Fading trail
			StrokeColor: color.RGBA{R: 255, G: 255, B: 255, A: 0},
			StrokeWidth: 0,
		}
		size := b.Radius * 0.3 * (1.0 - float32(i)*0.1) // Decreasing size
		trail.Resize(fyne.NewSize(size*2, size*2))
		trail.Move(fyne.NewPos(b.X-size, b.Y-size))
		b.Trail[i] = trail
	}
}

// updateTrail updates the particle trail positions
func (b *Ball) updateTrail() {
	if len(b.Trail) == 0 {
		return
	}

	// Move current trail position to the ball's old position
	currentTrail := b.Trail[b.TrailIndex]
	if currentTrail != nil {
		size := b.Radius * 0.3 * (1.0 - float32(b.TrailIndex)*0.1)
		currentTrail.Move(fyne.NewPos(b.X-size, b.Y-size))

		// Update trail color based on ball color
		ballColor := b.Circle.FillColor.(color.RGBA)
		alpha := uint8(255 - b.TrailIndex*20)
		currentTrail.FillColor = color.RGBA{
			R: ballColor.R,
			G: ballColor.G,
			B: ballColor.B,
			A: alpha,
		}
		currentTrail.Refresh()
	}

	// Move to next trail index
	b.TrailIndex = (b.TrailIndex + 1) % len(b.Trail)
}

// updatePosition updates the visual position of the circle and text
func (b *Ball) updatePosition() {
	// Apply jiggle effect to radius
	currentRadius := b.OriginalRadius
	if b.JiggleAmplitude > 0.01 { // Only apply jiggle if amplitude is significant
		// Create jiggle using sine wave with phase
		jiggleOffset := float32(math.Sin(float64(b.JigglePhase))) * b.JiggleAmplitude
		currentRadius = b.OriginalRadius + jiggleOffset

		// Update jiggle phase and decay amplitude
		b.JigglePhase += 0.3               // Reduced to 0.3 for more controlled oscillation
		b.JiggleAmplitude *= b.JiggleDecay // Use the decay rate from the struct

		// Stop jiggling when amplitude is very small
		if b.JiggleAmplitude < 0.01 {
			b.JiggleAmplitude = 0.0
			currentRadius = b.OriginalRadius
		}
	}

	// Update visual radius and position
	b.Radius = currentRadius
	b.Circle.Resize(fyne.NewSize(currentRadius*2, currentRadius*2))
	b.Circle.Move(fyne.NewPos(b.X-currentRadius, b.Y-currentRadius))

	// Update text position to center on ball
	if b.Text != nil {
		textSize := b.Text.MinSize()
		// Position text so its center aligns with ball center
		textX := b.X - textSize.Width/2
		textY := b.Y - textSize.Height/2
		b.Text.Move(fyne.NewPos(textX, textY))
		// Set text size to match its content
		b.Text.Resize(textSize)
	}

	// Update trail
	b.updateTrail()
}

// Update calculates the next position and handles wall bouncing
func (b *Ball) Update() {
	if !b.IsAnimated {
		return
	}

	// Update position
	b.X += b.VX
	b.Y += b.VY

	// Bounce off walls
	// Left and right walls
	if b.X-b.Radius <= 0 || b.X+b.Radius >= b.Bounds.Width {
		b.VX = -b.VX
		// Trigger jiggle effect based on impact velocity
		impactIntensity := float32(math.Abs(float64(b.VX))) / 8.0 // Increased to 8.0 for gentler effect
		b.triggerJiggle(impactIntensity)

		// Keep ball within bounds
		if b.X-b.Radius < 0 {
			b.X = b.Radius
		} else if b.X+b.Radius > b.Bounds.Width {
			b.X = b.Bounds.Width - b.Radius
		}
	}

	// Top and bottom walls
	if b.Y-b.Radius <= 50 || b.Y+b.Radius >= b.Bounds.Height-50 { // Account for button area
		b.VY = -b.VY
		// Trigger jiggle effect based on impact velocity
		impactIntensity := float32(math.Abs(float64(b.VY))) / 8.0 // Increased to 8.0 for gentler effect
		b.triggerJiggle(impactIntensity)

		// Keep ball within bounds
		if b.Y-b.Radius < 50 {
			b.Y = 50 + b.Radius
		} else if b.Y+b.Radius > b.Bounds.Height-50 {
			b.Y = b.Bounds.Height - 50 - b.Radius
		}
	}

	b.updatePosition()

	// Update explosion effects
	if b.IsExploding {
		b.UpdateExplosion()
	}
}

// CheckCollision checks if this ball collides with another ball
func (b *Ball) CheckCollision(other *Ball) bool {
	if b == other {
		return false
	}

	// Calculate distance between centers
	dx := b.X - other.X
	dy := b.Y - other.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Check if collision occurs (distance < sum of radii)
	return distance < (b.Radius + other.Radius)
}

// GetMass returns the mass of the ball based on its area (π * r²)
func (b *Ball) GetMass() float32 {
	return float32(math.Pi) * b.Radius * b.Radius
}

// HandleCollision handles elastic collision response between two balls with different masses
func (b *Ball) HandleCollision(other *Ball) {
	if b == other {
		return
	}

	// Calculate distance and collision normal
	dx := b.X - other.X
	dy := b.Y - other.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance == 0 {
		// Prevent division by zero - separate balls
		dx = 1
		distance = 1
	}

	// Normalize collision vector
	nx := dx / distance
	ny := dy / distance

	// Separate balls to prevent overlap based on mass ratio
	overlap := (b.Radius + other.Radius) - distance
	m1 := b.GetMass()
	m2 := other.GetMass()
	totalMass := m1 + m2

	// Heavier balls move less during separation
	separationRatio1 := m2 / totalMass
	separationRatio2 := m1 / totalMass

	b.X += nx * overlap * separationRatio1
	b.Y += ny * overlap * separationRatio1
	other.X -= nx * overlap * separationRatio2
	other.Y -= ny * overlap * separationRatio2

	// Calculate velocity components along collision normal
	v1n := b.VX*nx + b.VY*ny         // Ball 1 velocity along normal
	v2n := other.VX*nx + other.VY*ny // Ball 2 velocity along normal

	// Calculate velocity components perpendicular to collision normal
	v1p_x := b.VX - v1n*nx
	v1p_y := b.VY - v1n*ny
	v2p_x := other.VX - v2n*nx
	v2p_y := other.VY - v2n*ny

	// Do not resolve if velocities are separating
	if v1n-v2n > 0 {
		return
	}

	// Apply elastic collision formulas for 1D collision along normal
	// v1_new = ((m1-m2)/(m1+m2)) * v1_old + ((2*m2)/(m1+m2)) * v2_old
	// v2_new = ((m2-m1)/(m1+m2)) * v2_old + ((2*m1)/(m1+m2)) * v1_old
	v1n_new := ((m1-m2)/(m1+m2))*v1n + ((2*m2)/(m1+m2))*v2n
	v2n_new := ((m2-m1)/(m1+m2))*v2n + ((2*m1)/(m1+m2))*v1n

	// Reconstruct final velocities (normal + perpendicular components)
	b.VX = v1n_new*nx + v1p_x
	b.VY = v1n_new*ny + v1p_y
	other.VX = v2n_new*nx + v2p_x
	other.VY = v2n_new*ny + v2p_y

	// Optional: Add slight energy damping for more realistic behavior
	dampening := float32(0.98) // Less damping to preserve more energy
	b.VX *= dampening
	b.VY *= dampening
	other.VX *= dampening
	other.VY *= dampening

	// Trigger jiggle effects for both balls based on collision intensity
	collisionIntensity := float32(math.Sqrt(float64(v1n*v1n+v2n*v2n))) / 12.0 // Increased to 12.0 for gentler collision jiggle
	b.triggerJiggle(collisionIntensity)
	other.triggerJiggle(collisionIntensity)

	// Trigger explosions for both balls
	b.triggerExplosion()
	other.triggerExplosion()

	// Reduce ball sizes by 20%
	b.shrinkBall(0.8) // 0.8 = reduce to 80% of current size (20% reduction)
	other.shrinkBall(0.8)
}

// ChangeColor cycles through different colors
func (b *Ball) ChangeColor() {
	switch b.Circle.FillColor {
	case color.RGBA{R: 100, G: 150, B: 255, A: 255}: // Light blue
		b.Circle.FillColor = color.RGBA{R: 255, G: 100, B: 100, A: 255} // Light red
	case color.RGBA{R: 255, G: 100, B: 100, A: 255}: // Light red
		b.Circle.FillColor = color.RGBA{R: 100, G: 255, B: 100, A: 255} // Light green
	case color.RGBA{R: 100, G: 255, B: 100, A: 255}: // Light green
		b.Circle.FillColor = color.RGBA{R: 255, G: 255, B: 100, A: 255} // Light yellow
	case color.RGBA{R: 255, G: 255, B: 100, A: 255}: // Light yellow
		b.Circle.FillColor = color.RGBA{R: 255, G: 100, B: 255, A: 255} // Pink
	default:
		b.Circle.FillColor = color.RGBA{R: 100, G: 150, B: 255, A: 255} // Back to light blue
	}
	b.Circle.Refresh()
}

// NewCustomBall creates a ball with custom properties
func NewCustomBall(x, y, vx, vy, radius float32, fillColor, strokeColor color.RGBA) *Ball {
	ball := &Ball{
		X:       x,
		Y:       y,
		VX:      vx,
		VY:      vy,
		Radius:  radius,
		Bounds:  fyne.NewSize(800, 600),
		LLMName: getRandomLLMName(),
		// Initialize jiggle properties
		JiggleAmplitude: 0.0,
		JigglePhase:     0.0,
		JiggleDecay:     0.88, // Decay rate for jiggle amplitude
		OriginalRadius:  radius,
		// Initialize explosion properties
		ExplosionParticles: nil,
		ExplosionTimer:     0,
		IsExploding:        false,
	}

	// Create the visual circle with no stroke (ignoring strokeColor parameter)
	ball.Circle = &canvas.Circle{
		FillColor:   fillColor,
		StrokeColor: color.RGBA{R: 0, G: 0, B: 0, A: 0}, // No stroke
		StrokeWidth: 0,                                  // No stroke width
	}

	// Create the text label for the AI LLM name
	ball.Text = &canvas.Text{
		Text:      ball.LLMName,
		Color:     color.RGBA{R: 255, G: 255, B: 255, A: 255}, // White text
		Alignment: fyne.TextAlignCenter,
		TextStyle: fyne.TextStyle{Bold: true},
		TextSize:  12, // Initial size, will be adjusted
	}

	// Set initial size and position
	ball.Circle.Resize(fyne.NewSize(ball.Radius*2, ball.Radius*2))

	// Set font size to fit inside ball
	ball.updateTextSize()

	ball.updatePosition()

	// Initialize particle trail
	ball.initializeTrail()

	return ball
}

// triggerJiggle starts a jiggle effect (called when ball bounces)
func (b *Ball) triggerJiggle(intensity float32) {
	b.JiggleAmplitude = intensity * b.OriginalRadius * 0.8 // Reduced to 0.8 for subtle, natural jiggle
	b.JigglePhase = 0.0                                    // Reset phase
}

// triggerExplosion creates an explosion effect at the ball's location
func (b *Ball) triggerExplosion() {
	if b.IsExploding {
		return // Already exploding
	}

	b.IsExploding = true
	b.ExplosionTimer = 30 // 30 frames explosion duration (~0.5 seconds at 60 FPS)

	// Create explosion particles
	b.ExplosionParticles = make([]*canvas.Circle, 8) // 8 explosion particles per ball
	explosionColors := []color.RGBA{
		{R: 255, G: 255, B: 0, A: 255},   // Yellow
		{R: 255, G: 165, B: 0, A: 255},   // Orange
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 255, G: 255, B: 255, A: 255}, // White
	}

	for i := 0; i < 8; i++ {
		particle := &canvas.Circle{
			FillColor:   explosionColors[i%len(explosionColors)],
			StrokeColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			StrokeWidth: 1.0,
		}
		particle.Resize(fyne.NewSize(6, 6)) // Small particles
		particle.Move(fyne.NewPos(b.X-3, b.Y-3))
		b.ExplosionParticles[i] = particle
	}
}

// UpdateExplosion updates the explosion animation
func (b *Ball) UpdateExplosion() {
	if !b.IsExploding {
		return
	}

	b.ExplosionTimer--

	// Animate explosion particles
	if b.ExplosionTimer > 0 {
		explosionFrame := 30 - b.ExplosionTimer // 0 to 29
		for i, particle := range b.ExplosionParticles {
			if particle != nil {
				// Calculate particle movement in different directions
				angle := float64(i) * 2 * math.Pi / float64(len(b.ExplosionParticles))
				radius := float32(explosionFrame) * 1.5 // Smaller explosion radius than human

				newX := b.X + float32(math.Cos(angle))*radius - 3
				newY := b.Y + float32(math.Sin(angle))*radius - 3

				particle.Move(fyne.NewPos(newX, newY))

				// Fade particles
				alphaInt := 255 - int(explosionFrame*8)
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
	} else {
		// Explosion finished, clean up
		for _, particle := range b.ExplosionParticles {
			if particle != nil {
				particle.Hide()
			}
		}
		b.ExplosionParticles = nil
		b.IsExploding = false
	}
}

// shrinkBall reduces the ball size by the given factor
func (b *Ball) shrinkBall(factor float32) {
	// Calculate new size
	newRadius := b.Radius * factor
	newOriginalRadius := b.OriginalRadius * factor

	// Apply minimum size constraint (35 pixels radius)
	minRadius := float32(35.0)
	if newRadius < minRadius {
		newRadius = minRadius
		newOriginalRadius = minRadius
	}

	// Only shrink if the new size is different from current size
	if newRadius != b.Radius {
		// Hide old trail particles before resizing
		for _, trail := range b.Trail {
			if trail != nil {
				trail.Hide()
			}
		}

		// Update radius
		b.Radius = newRadius
		b.OriginalRadius = newOriginalRadius

		// Update visual size
		b.Circle.Resize(fyne.NewSize(b.Radius*2, b.Radius*2))

		// Adjust text size for new ball size
		b.updateTextSize()

		// Reinitialize trail with new size
		b.initializeTrail()
	}
}

// GetExplosionParticles returns the explosion particles for UI management
func (b *Ball) GetExplosionParticles() []*canvas.Circle {
	if b.ExplosionParticles == nil {
		return nil
	}
	return b.ExplosionParticles
}
