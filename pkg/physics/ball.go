package physics

import (
	"image"
	"image/color"
	"math"
	"math/rand"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Ball represents a ball with position and velocity
type Ball struct {
	X, Y       float32 // current position
	VX, VY     float32 // velocity
	Radius     float32 // ball radius
	LLMName    string  // AI LLM name
	Bounds     image.Point // animation bounds
	IsAnimated bool        // whether animation is running

	// Visual properties
	EyeballColor  color.NRGBA
	IrisColor     color.NRGBA
	PupilColor    color.NRGBA
	TextColor     color.NRGBA

	// Particle trail system
	TrailPositions []f32.Point
	TrailIndex     int

	// Jiggle effect for jello-like bouncing
	JiggleAmplitude float32 // Current jiggle strength
	JigglePhase     float32 // Current phase of jiggle oscillation
	JiggleDecay     float32 // How fast jiggle fades
	OriginalRadius  float32 // Original radius before jiggle

	// Explosion effects for ball collisions
	ExplosionParticles []f32.Point
	ExplosionTimer     int  // frames for explosion animation
	IsExploding        bool // whether ball is currently exploding

	// Human tracking for iris movement
	IrisOffsetX, IrisOffsetY float32
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

// getTextColorForLLM returns a bright, contrasting color for each LLM name
func getTextColorForLLM(llmName string) color.NRGBA {
	// Create distinctive colors for different AI models
	switch llmName {
	case "GPT-4":
		return color.NRGBA{R: 0, G: 255, B: 255, A: 220}   // Bright cyan
	case "Claude":
		return color.NRGBA{R: 255, G: 255, B: 0, A: 220}   // Bright yellow
	case "Gemini":
		return color.NRGBA{R: 255, G: 100, B: 255, A: 220} // Bright magenta
	case "LLaMA":
		return color.NRGBA{R: 100, G: 255, B: 100, A: 220} // Bright green
	case "PaLM":
		return color.NRGBA{R: 255, G: 150, B: 0, A: 220}   // Bright orange
	case "Bard":
		return color.NRGBA{R: 150, G: 255, B: 255, A: 220} // Light cyan
	case "ChatGPT":
		return color.NRGBA{R: 255, G: 255, B: 150, A: 220} // Light yellow
	case "Codex":
		return color.NRGBA{R: 255, G: 150, B: 255, A: 220} // Light magenta
	case "Alpaca":
		return color.NRGBA{R: 150, G: 255, B: 150, A: 220} // Light green
	case "Vicuna":
		return color.NRGBA{R: 255, G: 200, B: 100, A: 220} // Light orange
	case "Mistral":
		return color.NRGBA{R: 100, G: 200, B: 255, A: 220} // Light blue
	case "Llama2":
		return color.NRGBA{R: 255, G: 100, B: 150, A: 220} // Pink
	default:
		return color.NRGBA{R: 255, G: 255, B: 255, A: 220} // White as fallback
	}
}

// NewBall creates a new bouncing ball that looks like an eyeball
func NewBall() *Ball {
	llmName := getRandomLLMName()
	ball := &Ball{
		X:       100,
		Y:       100,
		VX:      3.5, // horizontal velocity
		VY:      2.8, // vertical velocity
		Radius:  30,
		Bounds:  image.Point{X: 800, Y: 600},
		LLMName: llmName,
		// Initialize jiggle properties
		JiggleAmplitude: 0.0,
		JigglePhase:     0.0,
		JiggleDecay:     0.88, // Decay rate for jiggle amplitude
		OriginalRadius:  30,
		// Initialize explosion properties
		ExplosionParticles: nil,
		ExplosionTimer:     0,
		IsExploding:        false,
		// Visual properties
		EyeballColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // White eyeball
		IrisColor:    color.NRGBA{R: 100, G: 150, B: 255, A: 255}, // Blue iris
		PupilColor:   color.NRGBA{R: 0, G: 0, B: 0, A: 255},       // Black pupil
		TextColor:    getTextColorForLLM(llmName),
	}

	// Initialize particle trail
	ball.initializeTrail()

	return ball
}

// NewCustomBall creates a ball with custom properties that looks like an eyeball
func NewCustomBall(x, y, vx, vy, radius float32, fillColor, strokeColor color.RGBA) *Ball {
	llmName := getRandomLLMName()
	ball := &Ball{
		X:       x,
		Y:       y,
		VX:      vx,
		VY:      vy,
		Radius:  radius,
		Bounds:  image.Point{X: 800, Y: 600},
		LLMName: llmName,
		// Initialize jiggle properties
		JiggleAmplitude: 0.0,
		JigglePhase:     0.0,
		JiggleDecay:     0.88, // Decay rate for jiggle amplitude
		OriginalRadius:  radius,
		// Initialize explosion properties
		ExplosionParticles: nil,
		ExplosionTimer:     0,
		IsExploding:        false,
		// Visual properties
		EyeballColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, // White eyeball
		IrisColor:    color.NRGBA{R: fillColor.R, G: fillColor.G, B: fillColor.B, A: fillColor.A},
		PupilColor:   color.NRGBA{R: 0, G: 0, B: 0, A: 255}, // Black pupil
		TextColor:    getTextColorForLLM(llmName),
	}

	// Initialize particle trail
	ball.initializeTrail()

	return ball
}

// initializeTrail creates the particle trail for the ball
func (b *Ball) initializeTrail() {
	trailLength := 10 // Number of trail particles
	b.TrailPositions = make([]f32.Point, trailLength)
	b.TrailIndex = 0 // Reset trail index

	for i := 0; i < trailLength; i++ {
		b.TrailPositions[i] = f32.Point{X: b.X, Y: b.Y}
	}
}

// updateTrail updates the particle trail positions
func (b *Ball) updateTrail() {
	if len(b.TrailPositions) == 0 {
		return
	}

	// Move current trail position to the ball's old position
	b.TrailPositions[b.TrailIndex] = f32.Point{X: b.X, Y: b.Y}

	// Move to next trail index
	b.TrailIndex = (b.TrailIndex + 1) % len(b.TrailPositions)
}

// UpdatePosition updates the visual position of the eyeball components
func (b *Ball) UpdatePosition() {
	b.UpdatePositionWithHuman(0, 0) // Default position when no human tracking
}

// UpdatePositionWithHuman updates the visual position of the eyeball components with human tracking
func (b *Ball) UpdatePositionWithHuman(humanX, humanY float32) {
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

	// Update visual radius
	b.Radius = currentRadius

	// Calculate direction to human for iris tracking
	b.IrisOffsetX, b.IrisOffsetY = 0, 0
	if humanX != 0 || humanY != 0 { // If human position is provided
		dx := humanX - b.X
		dy := humanY - b.Y
		distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		if distance > 0 {
			// Normalize direction and apply offset (iris can move within the eyeball)
			maxOffset := currentRadius * 0.3 // Maximum iris offset from center
			b.IrisOffsetX = (dx / distance) * maxOffset
			b.IrisOffsetY = (dy / distance) * maxOffset
		}
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

	// Handle bouncing off boundaries
	if b.X-b.Radius <= 0 || b.X+b.Radius >= float32(b.Bounds.X) {
		b.VX = -b.VX
		b.triggerJiggle(0.8) // Reduced jiggle intensity for more natural effect

		// Keep ball within bounds
		if b.X-b.Radius <= 0 {
			b.X = b.Radius
		}
		if b.X+b.Radius >= float32(b.Bounds.X) {
			b.X = float32(b.Bounds.X) - b.Radius
		}
	}

	if b.Y-b.Radius <= 0 || b.Y+b.Radius >= float32(b.Bounds.Y) {
		b.VY = -b.VY
		b.triggerJiggle(0.8) // Reduced jiggle intensity for more natural effect

		// Keep ball within bounds
		if b.Y-b.Radius <= 0 {
			b.Y = b.Radius
		}
		if b.Y+b.Radius >= float32(b.Bounds.Y) {
			b.Y = float32(b.Bounds.Y) - b.Radius
		}
	}

	// Update visual position
	b.UpdatePosition()

	// Update explosion if active
	b.UpdateExplosion()
}

// CheckCollision checks if this ball collides with another ball
func (b *Ball) CheckCollision(other *Ball) bool {
	if !b.IsAnimated || !other.IsAnimated {
		return false
	}

	dx := b.X - other.X
	dy := b.Y - other.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	return distance < (b.Radius + other.Radius)
}

// GetMass returns the ball's mass (based on volume/radius)
func (b *Ball) GetMass() float32 {
	return b.Radius * b.Radius // Simplified mass calculation
}

// HandleCollision handles collision physics between two balls
func (b *Ball) HandleCollision(other *Ball) {
	// Calculate collision normal
	dx := b.X - other.X
	dy := b.Y - other.Y
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if distance == 0 {
		return // Avoid division by zero
	}

	// Normalize collision vector
	nx := dx / distance
	ny := dy / distance

	// Calculate relative velocity
	dvx := b.VX - other.VX
	dvy := b.VY - other.VY

	// Calculate relative velocity in collision normal direction
	dvn := dvx*nx + dvy*ny

	// Do not resolve if velocities are separating
	if dvn > 0 {
		return
	}

	// Calculate collision impulse
	impulse := 2 * dvn / (b.GetMass() + other.GetMass())

	// Update velocities
	b.VX -= impulse * other.GetMass() * nx
	b.VY -= impulse * other.GetMass() * ny
	other.VX += impulse * b.GetMass() * nx
	other.VY += impulse * b.GetMass() * ny

	// Separate overlapping balls
	overlap := (b.Radius + other.Radius) - distance
	if overlap > 0 {
		separationX := nx * overlap * 0.5
		separationY := ny * overlap * 0.5
		b.X += separationX
		b.Y += separationY
		other.X -= separationX
		other.Y -= separationY
	}

	// Apply dampening to prevent infinite bouncing
	dampening := float32(0.95)
	b.VX *= dampening
	b.VY *= dampening
	other.VX *= dampening
	other.VY *= dampening

	// Trigger jiggle effects for both balls based on collision intensity
	collisionIntensity := float32(math.Sqrt(float64(dvn*dvn))) / 12.0 // Increased to 12.0 for gentler collision jiggle
	b.triggerJiggle(collisionIntensity)
	other.triggerJiggle(collisionIntensity)

	// Trigger explosions for both balls
	b.triggerExplosion()
	other.triggerExplosion()

	// Reduce ball sizes by 20%
	b.shrinkBall(0.8) // 0.8 = reduce to 80% of current size (20% reduction)
	other.shrinkBall(0.8)
}

// ChangeColor cycles through different iris colors for the eyeball
func (b *Ball) ChangeColor() {
	switch b.IrisColor {
	case (color.NRGBA{R: 100, G: 150, B: 255, A: 255}): // Blue iris
		b.IrisColor = color.NRGBA{R: 100, G: 255, B: 100, A: 255} // Green iris
	case (color.NRGBA{R: 100, G: 255, B: 100, A: 255}): // Green iris
		b.IrisColor = color.NRGBA{R: 139, G: 69, B: 19, A: 255} // Brown iris
	case (color.NRGBA{R: 139, G: 69, B: 19, A: 255}): // Brown iris
		b.IrisColor = color.NRGBA{R: 128, G: 128, B: 128, A: 255} // Gray iris
	case (color.NRGBA{R: 128, G: 128, B: 128, A: 255}): // Gray iris
		b.IrisColor = color.NRGBA{R: 255, G: 140, B: 0, A: 255} // Orange iris
	default:
		b.IrisColor = color.NRGBA{R: 100, G: 150, B: 255, A: 255} // Back to blue iris
	}
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
	b.ExplosionParticles = make([]f32.Point, 8) // 8 explosion particles per ball
	for i := 0; i < 8; i++ {
		b.ExplosionParticles[i] = f32.Point{X: b.X, Y: b.Y}
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
		for i := range b.ExplosionParticles {
			// Calculate particle movement in different directions
			angle := float64(i) * 2 * math.Pi / float64(len(b.ExplosionParticles))
			radius := float32(explosionFrame) * 1.5 // Smaller explosion radius than human

			b.ExplosionParticles[i] = f32.Point{
				X: b.X + float32(math.Cos(angle))*radius,
				Y: b.Y + float32(math.Sin(angle))*radius,
			}
		}
	} else {
		// Explosion finished, clean up
		b.ExplosionParticles = nil
		b.IsExploding = false
	}
}

// shrinkBall reduces the eyeball size by the given factor
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
		// Update radius
		b.Radius = newRadius
		b.OriginalRadius = newOriginalRadius

		// Reinitialize trail with new size
		b.initializeTrail()
	}
}

// Render draws the ball using Gio operations
func (b *Ball) Render(ops *op.Ops) {
	// Render trail particles first (behind the ball)
	b.renderTrail(ops)

	// Render explosion particles if exploding
	if b.IsExploding {
		b.renderExplosion(ops)
	}

	// Render main eyeball
	b.renderEyeball(ops)

	// Render bloodshot veins
	b.renderBloodVeins(ops)

	// Render text label below the ball
	b.renderText(ops)
}

// renderEyeball renders the main eyeball components
func (b *Ball) renderEyeball(ops *op.Ops) {
	// Draw white eyeball background
	eyeballRect := image.Rectangle{
		Min: image.Point{X: int(b.X - b.Radius), Y: int(b.Y - b.Radius)},
		Max: image.Point{X: int(b.X + b.Radius), Y: int(b.Y + b.Radius)},
	}
	defer clip.Ellipse(eyeballRect).Push(ops).Pop()
	paint.Fill(ops, b.EyeballColor)

	// Draw iris (with human tracking offset)
	irisRadius := b.Radius * 0.6
	irisX := b.X + b.IrisOffsetX
	irisY := b.Y + b.IrisOffsetY
	irisRect := image.Rectangle{
		Min: image.Point{X: int(irisX - irisRadius), Y: int(irisY - irisRadius)},
		Max: image.Point{X: int(irisX + irisRadius), Y: int(irisY + irisRadius)},
	}
	defer clip.Ellipse(irisRect).Push(ops).Pop()
	paint.Fill(ops, b.IrisColor)

	// Draw pupil (follows iris)
	pupilRadius := b.Radius * 0.3
	pupilRect := image.Rectangle{
		Min: image.Point{X: int(irisX - pupilRadius), Y: int(irisY - pupilRadius)},
		Max: image.Point{X: int(irisX + pupilRadius), Y: int(irisY + pupilRadius)},
	}
	defer clip.Ellipse(pupilRect).Push(ops).Pop()
	paint.Fill(ops, b.PupilColor)
}

// renderTrail renders the particle trail
func (b *Ball) renderTrail(ops *op.Ops) {
	for i, pos := range b.TrailPositions {
		alpha := uint8(255 - i*20) // Fading trail
		if alpha < 50 {
			continue
		}

		size := b.Radius * 0.3 * (1.0 - float32(i)*0.1) // Decreasing size
		trailColor := color.NRGBA{R: b.EyeballColor.R, G: b.EyeballColor.G, B: b.EyeballColor.B, A: alpha}

		trailRect := image.Rectangle{
			Min: image.Point{X: int(pos.X - size), Y: int(pos.Y - size)},
			Max: image.Point{X: int(pos.X + size), Y: int(pos.Y + size)},
		}
		defer clip.Ellipse(trailRect).Push(ops).Pop()
		paint.Fill(ops, trailColor)
	}
}

// renderExplosion renders explosion particles
func (b *Ball) renderExplosion(ops *op.Ops) {
	explosionColors := []color.NRGBA{
		{R: 255, G: 255, B: 0, A: 255},   // Yellow
		{R: 255, G: 165, B: 0, A: 255},   // Orange
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 255, G: 255, B: 255, A: 255}, // White
	}

	explosionFrame := 30 - b.ExplosionTimer
	alpha := uint8(255 - explosionFrame*8)
	if alpha < 10 {
		return
	}

	for i, pos := range b.ExplosionParticles {
		particleColor := explosionColors[i%len(explosionColors)]
		particleColor.A = alpha

		particleRect := image.Rectangle{
			Min: image.Point{X: int(pos.X - 3), Y: int(pos.Y - 3)},
			Max: image.Point{X: int(pos.X + 3), Y: int(pos.Y + 3)},
		}
		defer clip.Ellipse(particleRect).Push(ops).Pop()
		paint.Fill(ops, particleColor)
	}
}

// renderBloodVeins renders bloodshot veins around the eyeball
func (b *Ball) renderBloodVeins(ops *op.Ops) {
	veinColor := color.NRGBA{R: 200, G: 50, B: 50, A: 180} // Semi-transparent red

	// Draw 6 blood vessels radiating from different points
	for i := 0; i < 6; i++ {
		angle := float64(i) * math.Pi / 3.0 // 60 degree intervals

		// Start point (closer to edge of eyeball)
		startRadius := b.Radius * 0.7
		startX := b.X + float32(math.Cos(angle))*startRadius
		startY := b.Y + float32(math.Sin(angle))*startRadius

		// End point (towards center but not reaching iris)
		endRadius := b.Radius * 0.3
		endX := b.X + float32(math.Cos(angle+0.5))*endRadius // Slight curve
		endY := b.Y + float32(math.Sin(angle+0.5))*endRadius

		// Create a thin line for the vein
		veinRect := image.Rectangle{
			Min: image.Point{X: int(startX - 1), Y: int(startY - 1)},
			Max: image.Point{X: int(endX + 1), Y: int(endY + 1)},
		}
		defer clip.Rect(veinRect).Push(ops).Pop()
		paint.Fill(ops, veinColor)
	}
}

// renderText renders the AI LLM name text below the ball
func (b *Ball) renderText(ops *op.Ops) {
	// TODO: Implement text rendering with Gio
	// This is more complex and would require text shaping
	// For now, we'll skip text rendering and add it later
}
