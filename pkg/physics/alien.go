package physics

import (
	"math"
	"math/rand"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
)

// Alien represents a mysterious alien face that drifts through the star field
type Alien struct {
	X, Y          float32   // current position
	VX, VY        float32   // drift velocity
	Size          float32   // size of the alien face
	Bounds        fyne.Size // movement bounds
	IsActive      bool      // whether the alien is active
	// Visual components
	Image         *canvas.Image     // Alien face image
	ImageContainer *fyne.Container // Container for the image
	// Drift behavior
	DriftTimer    int     // frames until direction change
	DriftDuration int     // frames between direction changes
	Alpha         float32 // transparency (0.0 to 1.0)
	// Mysterious behavior
	PhaseOffset   float32 // for subtle floating motion
	FloatAmplitude float32 // how much it bobs up and down
}

// NewAlien creates a new alien entity that drifts through space
func NewAlien(x, y, size float32) *Alien {
	alien := &Alien{
		X:             x,
		Y:             y,
		VX:            (rand.Float32() - 0.5) * 0.8, // Slow random drift
		VY:            (rand.Float32() - 0.5) * 0.8,
		Size:          size,
		Bounds:        fyne.NewSize(800, 600),
		IsActive:      true,
		DriftTimer:    rand.Intn(300) + 180, // 3-8 seconds at 60fps
		DriftDuration: 300,                   // 5 seconds default
		Alpha:         0.7,                   // Semi-transparent
		PhaseOffset:   rand.Float32() * 2 * math.Pi,
		FloatAmplitude: 2.0, // Subtle floating motion
	}

	// Create alien face image placeholder (will be replaced in NewAlienFromFile)
	alien.Image = canvas.NewImageFromResource(nil)
	alien.Image.FillMode = canvas.ImageFillOriginal
	alien.Image.ScaleMode = canvas.ImageScaleSmooth
	alien.Image.Resize(fyne.NewSize(size, size))

	// Create container
	alien.ImageContainer = container.NewWithoutLayout(alien.Image)

	// Set initial position
	alien.UpdatePosition()

	// Alien created successfully

	return alien
}

// NewAlienFromResource creates an alien with a specific image resource
func NewAlienFromResource(x, y, size float32, resource fyne.Resource) *Alien {
	// For now, just create a simple circle alien
	return NewAlien(x, y, size)
}

// NewAlienFromFile creates an alien with an image from file
func NewAlienFromFile(x, y, size float32, filename string) *Alien {
	alien := NewAlien(x, y, size)

		// Load image from file using absolute path
	cwd, _ := os.Getwd()
	fullPath := filepath.Join(cwd, filename)
	resource := storage.NewFileURI(fullPath)
	alien.Image = canvas.NewImageFromURI(resource)

	// Fallback to human.png if alien.png doesn't exist
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		humanPath := filepath.Join(cwd, "human.png")
		resource = storage.NewFileURI(humanPath)
		alien.Image = canvas.NewImageFromURI(resource)
	}

	alien.Image.FillMode = canvas.ImageFillOriginal
	alien.Image.ScaleMode = canvas.ImageScaleSmooth
	alien.Image.Resize(fyne.NewSize(size, size))

	// Update container with new image
	alien.ImageContainer = container.NewWithoutLayout(alien.Image)

	return alien
}

// Update handles the alien's drift behavior
func (a *Alien) Update() {
	if !a.IsActive {
		return
	}

	// Alien drifting peacefully through space

	// Update drift timer
	a.DriftTimer--

	// Change direction randomly when timer expires
	if a.DriftTimer <= 0 {
		a.changeDirection()
		a.DriftTimer = rand.Intn(300) + 180 // 3-8 seconds
	}

	// Apply drift movement
	a.X += a.VX
	a.Y += a.VY

	// Add subtle floating motion
	a.PhaseOffset += 0.02 // Slow phase increment

	// Wrap around screen edges for mysterious appearances
	a.wrapAroundScreen()

	// Update visual position with floating effect
	a.UpdatePosition()
}

// changeDirection randomly changes the alien's drift direction
func (a *Alien) changeDirection() {
	// Generate new random drift velocity (very slow)
	maxSpeed := float32(0.8)
	a.VX = (rand.Float32() - 0.5) * maxSpeed
	a.VY = (rand.Float32() - 0.5) * maxSpeed

	// Sometimes pause (no movement)
	if rand.Float32() < 0.2 { // 20% chance to pause
		a.VX = 0
		a.VY = 0
	}
}

// wrapAroundScreen makes the alien wrap around screen edges
func (a *Alien) wrapAroundScreen() {
	margin := a.Size

	// Wrap horizontally
	if a.X < -margin {
		a.X = a.Bounds.Width + margin
	} else if a.X > a.Bounds.Width + margin {
		a.X = -margin
	}

	// Wrap vertically
	if a.Y < -margin {
		a.Y = a.Bounds.Height + margin
	} else if a.Y > a.Bounds.Height + margin {
		a.Y = -margin
	}
}

// UpdatePosition updates the visual position of the alien
func (a *Alien) UpdatePosition() {
	if !a.IsActive {
		return
	}

	// Add floating effect
	floatOffset := float32(math.Sin(float64(a.PhaseOffset))) * a.FloatAmplitude
	displayY := a.Y + floatOffset

				// Center the image on the alien's position
	baseX := a.X - a.Size/2
	baseY := displayY - a.Size/2

	// Update image position
	a.Image.Move(fyne.NewPos(baseX, baseY))

	// Keep consistent size
	a.Image.Resize(fyne.NewSize(a.Size, a.Size))
}

// SetAlpha sets the transparency of the alien (0.0 = invisible, 1.0 = opaque)
func (a *Alien) SetAlpha(alpha float32) {
	a.Alpha = alpha
	// Note: Fyne doesn't have direct alpha support for images
	// This could be implemented with custom rendering if needed
}

// Hide makes the alien invisible
func (a *Alien) Hide() {
	a.IsActive = false
	a.ImageContainer.Hide()
}

// Show makes the alien visible
func (a *Alien) Show() {
	a.IsActive = true
	a.ImageContainer.Show()
}

// GetVisualComponents returns the alien's visual components for UI management
func (a *Alien) GetVisualComponents() []fyne.CanvasObject {
	return []fyne.CanvasObject{a.ImageContainer}
}

// SetBounds updates the movement bounds for the alien
func (a *Alien) SetBounds(bounds fyne.Size) {
	a.Bounds = bounds
}

// Respawn moves the alien to a random edge position for mysterious re-entry
func (a *Alien) Respawn() {
	margin := a.Size

	// Choose random edge (0=top, 1=right, 2=bottom, 3=left)
	edge := rand.Intn(4)

	switch edge {
	case 0: // Top edge
		a.X = rand.Float32() * a.Bounds.Width
		a.Y = -margin
	case 1: // Right edge
		a.X = a.Bounds.Width + margin
		a.Y = rand.Float32() * a.Bounds.Height
	case 2: // Bottom edge
		a.X = rand.Float32() * a.Bounds.Width
		a.Y = a.Bounds.Height + margin
	case 3: // Left edge
		a.X = -margin
		a.Y = rand.Float32() * a.Bounds.Height
	}

	// Set new random drift direction toward screen
	a.changeDirection()

	// Ensure movement toward screen center
	centerX := a.Bounds.Width / 2
	centerY := a.Bounds.Height / 2

	if a.X < centerX {
		if a.VX < 0 {
			a.VX = -a.VX // Move right
		}
	} else {
		if a.VX > 0 {
			a.VX = -a.VX // Move left
		}
	}

	if a.Y < centerY {
		if a.VY < 0 {
			a.VY = -a.VY // Move down
		}
	} else {
		if a.VY > 0 {
			a.VY = -a.VY // Move up
		}
	}
}