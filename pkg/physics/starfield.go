package physics

import (
	"image/color"
	"math"
	"math/rand"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

// StarType represents different types of stars with realistic properties
type StarType int

const (
	RedDwarf StarType = iota
	OrangeDwarf
	YellowDwarf
	BlueWhiteDwarf
	RedGiant
	BlueGiant
	WhiteDwarf
	Neutron
)

// StarClass holds properties for different star types
type StarClass struct {
	Name         string
	Color        color.RGBA
	MinSize      float32
	MaxSize      float32
	MinBrightness uint8
	MaxBrightness uint8
	Frequency    float32 // How common this star type is (0.0 to 1.0)
}

// Star represents a single star in the star field
type Star struct {
	X, Y         float32   // current position
	InitialX, InitialY float32 // initial position for regeneration
	StarType     StarType  // type of star
	Distance     float32   // simulated distance (affects brightness and parallax)
	Size         float32   // star size
	Brightness   uint8     // star brightness (alpha value)
	TwinklePhase float32   // current twinkling phase
	Visual       *canvas.Circle
}

// StarField represents a collection of moving stars with realistic distribution
type StarField struct {
	Stars       []*Star
	Bounds      fyne.Size
	GalacticCenterX float32 // Center of galaxy for density calculations
	GalacticCenterY float32
	StarClasses map[StarType]StarClass
	TravelSpeed float32     // Base speed of travel through space
	TravelAngle float32     // Direction of travel (in radians)
}

// Initialize star classification system based on real stellar populations
func getStarClasses() map[StarType]StarClass {
	return map[StarType]StarClass{
		RedDwarf: {
			Name:         "Red Dwarf (M-class)",
			Color:        color.RGBA{R: 255, G: 180, B: 120, A: 255},
			MinSize:      0.8,
			MaxSize:      1.5,
			MinBrightness: 120,
			MaxBrightness: 180,
			Frequency:    0.76, // 76% of all stars
		},
		OrangeDwarf: {
			Name:         "Orange Dwarf (K-class)",
			Color:        color.RGBA{R: 255, G: 200, B: 140, A: 255},
			MinSize:      1.2,
			MaxSize:      2.0,
			MinBrightness: 140,
			MaxBrightness: 200,
			Frequency:    0.12, // 12% of all stars
		},
		YellowDwarf: {
			Name:         "Yellow Dwarf (G-class)",
			Color:        color.RGBA{R: 255, G: 255, B: 200, A: 255},
			MinSize:      1.5,
			MaxSize:      2.5,
			MinBrightness: 160,
			MaxBrightness: 220,
			Frequency:    0.076, // 7.6% of all stars (like our Sun)
		},
		BlueWhiteDwarf: {
			Name:         "Blue-White Dwarf (A/F-class)",
			Color:        color.RGBA{R: 200, G: 220, B: 255, A: 255},
			MinSize:      1.8,
			MaxSize:      3.0,
			MinBrightness: 180,
			MaxBrightness: 240,
			Frequency:    0.03, // 3% of all stars
		},
		RedGiant: {
			Name:         "Red Giant",
			Color:        color.RGBA{R: 255, G: 140, B: 80, A: 255},
			MinSize:      3.0,
			MaxSize:      6.0,
			MinBrightness: 200,
			MaxBrightness: 255,
			Frequency:    0.008, // 0.8% of all stars
		},
		BlueGiant: {
			Name:         "Blue Giant (O/B-class)",
			Color:        color.RGBA{R: 150, G: 180, B: 255, A: 255},
			MinSize:      4.0,
			MaxSize:      8.0,
			MinBrightness: 220,
			MaxBrightness: 255,
			Frequency:    0.00003, // 0.003% of all stars (very rare)
		},
		WhiteDwarf: {
			Name:         "White Dwarf",
			Color:        color.RGBA{R: 240, G: 240, B: 255, A: 255},
			MinSize:      0.5,
			MaxSize:      1.0,
			MinBrightness: 100,
			MaxBrightness: 160,
			Frequency:    0.003, // 0.3% of all stars
		},
		Neutron: {
			Name:         "Neutron Star",
			Color:        color.RGBA{R: 180, G: 200, B: 255, A: 255},
			MinSize:      0.3,
			MaxSize:      0.8,
			MinBrightness: 80,
			MaxBrightness: 140,
			Frequency:    0.001, // 0.1% of all stars (extremely rare)
		},
	}
}

// NewStarField creates a new realistic star field with space travel effect
func NewStarField(numStars int, bounds fyne.Size) *StarField {
	starField := &StarField{
		Stars:           make([]*Star, numStars),
		Bounds:          bounds,
		GalacticCenterX: bounds.Width * 0.6,  // Offset galactic center
		GalacticCenterY: bounds.Height * 0.4,
		StarClasses:     getStarClasses(),
		TravelSpeed:     0.75, // Base speed of travel through space (slowed by half)
		TravelAngle:     math.Pi,  // Traveling horizontally to the left (stars move right)
	}

	// Create stars with realistic distribution
	for i := 0; i < numStars; i++ {
		starField.Stars[i] = starField.createRealisticStar()
	}

	return starField
}

// selectStarType chooses a star type based on realistic stellar population frequencies
func (sf *StarField) selectStarType() StarType {
	random := rand.Float32()
	cumulative := float32(0.0)

	// Use cumulative frequency distribution
	for starType, class := range sf.StarClasses {
		cumulative += class.Frequency
		if random <= cumulative {
			return starType
		}
	}

	// Fallback to most common star type
	return RedDwarf
}

// calculateGalacticDensity returns density multiplier based on distance from galactic center
func (sf *StarField) calculateGalacticDensity(x, y float32) float32 {
	// Distance from galactic center
	dx := x - sf.GalacticCenterX
	dy := y - sf.GalacticCenterY
	distance := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// Normalize distance (0.0 at center, 1.0 at edges)
	maxDistance := float32(math.Sqrt(float64(sf.Bounds.Width*sf.Bounds.Width + sf.Bounds.Height*sf.Bounds.Height)))
	normalizedDistance := distance / maxDistance

	// Create galactic disk density profile (exponential falloff)
	densityMultiplier := float32(math.Exp(-2.0 * float64(normalizedDistance)))

	// Add spiral arm enhancement
	angle := math.Atan2(float64(dy), float64(dx))
	spiralArm1 := math.Sin(angle*2 + float64(normalizedDistance)*4)
	spiralArm2 := math.Sin(angle*2 + math.Pi + float64(normalizedDistance)*4)
	spiralEnhancement := float32(1.0 + 0.3*math.Max(spiralArm1, spiralArm2))

	return densityMultiplier * spiralEnhancement
}

// createRealisticStar creates a star with realistic properties and galactic distribution
func (sf *StarField) createRealisticStar() *Star {
	// Use rejection sampling for realistic galactic distribution
	var x, y float32
	for attempts := 0; attempts < 100; attempts++ {
		x = rand.Float32() * sf.Bounds.Width
		y = rand.Float32() * sf.Bounds.Height

		density := sf.calculateGalacticDensity(x, y)
		if rand.Float32() < density {
			break // Accept this position
		}
	}

	// Select star type based on realistic frequencies
	starType := sf.selectStarType()
	starClass := sf.StarClasses[starType]

	// Generate star properties based on star class
	size := starClass.MinSize + rand.Float32()*(starClass.MaxSize-starClass.MinSize)
	brightness := starClass.MinBrightness + uint8(rand.Float32()*float32(starClass.MaxBrightness-starClass.MinBrightness))

	// Simulate distance (affects parallax and brightness)
	// Use more varied distance distribution for better parallax effect
	distance := rand.Float32()*rand.Float32()*0.9 + 0.1 // Skewed towards closer stars (0.1 to 1.0)

	// Adjust brightness based on distance (inverse square law approximation)
	distanceBrightness := float32(brightness) * (1.0 / (distance * distance + 0.5))
	if distanceBrightness > 255 {
		distanceBrightness = 255
	}
	brightness = uint8(distanceBrightness)

	// Adjust size based on distance (closer stars appear larger)
	size = size * (2.0 - distance) // Closer stars can be up to 2x larger

	star := &Star{
		X:            x,
		Y:            y,
		InitialX:     x,
		InitialY:     y,
		StarType:     starType,
		Distance:     distance,
		Size:         size,
		Brightness:   brightness,
		TwinklePhase: rand.Float32() * 2 * math.Pi,
	}

	// Create visual representation with star-type-specific color
	starColor := starClass.Color
	starColor.A = brightness

	star.Visual = &canvas.Circle{
		FillColor:   starColor,
		StrokeColor: color.RGBA{R: starColor.R, G: starColor.G, B: starColor.B, A: 0},
		StrokeWidth: 0,
	}

	star.Visual.Resize(fyne.NewSize(star.Size, star.Size))
	star.Visual.Move(fyne.NewPos(star.X-star.Size/2, star.Y-star.Size/2))

	return star
}

// Update updates all stars with space travel parallax effect
func (sf *StarField) Update() {
	for _, star := range sf.Stars {
		if star == nil {
			continue
		}

		// Calculate parallax speed based on distance
		// Closer stars (lower distance values) move faster
		parallaxMultiplier := (1.0 - star.Distance) * 3.0 + 0.5 // Range from 0.5x to 3.5x speed

		// Simple horizontal movement: we're traveling forward, so stars move from right to left
		star.X -= sf.TravelSpeed * parallaxMultiplier

		// Regenerate stars that have moved off the left edge
		margin := float32(50.0)
		if star.X < -margin {
			// Regenerate star on the right edge with random Y position
			star.X = sf.Bounds.Width + margin
			star.Y = rand.Float32() * sf.Bounds.Height

			// Generate new star properties for variety
			sf.regenerateStarProperties(star)
		}

		// Update visual position
		star.Visual.Move(fyne.NewPos(star.X-star.Size/2, star.Y-star.Size/2))

		// Advanced twinkling based on star type and atmospheric effects
		star.updateTwinkling()
	}
}

// regenerateStarProperties generates new properties for a star that has moved off screen
func (sf *StarField) regenerateStarProperties(star *Star) {
	// Generate new star properties for variety
	starType := sf.selectStarType()
	starClass := sf.StarClasses[starType]

	// Generate new properties
	size := starClass.MinSize + rand.Float32()*(starClass.MaxSize-starClass.MinSize)
	brightness := starClass.MinBrightness + uint8(rand.Float32()*float32(starClass.MaxBrightness-starClass.MinBrightness))
	distance := rand.Float32()*rand.Float32()*0.9 + 0.1 // Skewed towards closer stars

	// Adjust brightness and size based on distance
	distanceBrightness := float32(brightness) * (1.0 / (distance * distance + 0.5))
	if distanceBrightness > 255 {
		distanceBrightness = 255
	}
	brightness = uint8(distanceBrightness)
	size = size * (2.0 - distance)

	// Update star properties
	star.StarType = starType
	star.Distance = distance
	star.Size = size
	star.Brightness = brightness
	star.TwinklePhase = rand.Float32() * 2 * math.Pi

	// Update visual properties
	starColor := starClass.Color
	starColor.A = brightness
	star.Visual.FillColor = starColor
	star.Visual.Resize(fyne.NewSize(star.Size, star.Size))
}

// updateTwinkling creates realistic twinkling effects
func (s *Star) updateTwinkling() {
	// Update twinkling phase
	twinkleSpeed := 0.05 + rand.Float32()*0.03 // Vary twinkling speed
	s.TwinklePhase += twinkleSpeed

	// Different star types twinkle differently
	baseBrightness := float32(s.Brightness)

	// Calculate twinkling intensity based on star type and distance
	twinkleIntensity := float32(0.1) // Base twinkling
	if s.StarType == BlueGiant || s.StarType == RedGiant {
		twinkleIntensity *= 1.5 // Giants twinkle more
	}
	if s.Distance < 0.5 {
		twinkleIntensity *= (1.0 - s.Distance) // Closer stars twinkle more
	}

	// Apply sine wave twinkling
	twinkleEffect := float32(math.Sin(float64(s.TwinklePhase))) * twinkleIntensity
	newBrightness := baseBrightness * (1.0 + twinkleEffect)

	// Clamp brightness
	if newBrightness < 50 {
		newBrightness = 50
	} else if newBrightness > 255 {
		newBrightness = 255
	}

	// Update visual color with new brightness
	currentColor := s.Visual.FillColor.(color.RGBA)
	currentColor.A = uint8(newBrightness)
	s.Visual.FillColor = currentColor
	s.Visual.Refresh()
}

// GetVisuals returns all star visual objects for UI management
func (sf *StarField) GetVisuals() []*canvas.Circle {
	visuals := make([]*canvas.Circle, 0, len(sf.Stars))
	for _, star := range sf.Stars {
		if star != nil && star.Visual != nil {
			visuals = append(visuals, star.Visual)
		}
	}
	return visuals
}

// UpdateBounds updates the star field bounds and redistributes stars
func (sf *StarField) UpdateBounds(newBounds fyne.Size) {
	sf.Bounds = newBounds
	sf.GalacticCenterX = newBounds.Width * 0.6
	sf.GalacticCenterY = newBounds.Height * 0.4

	// Reposition stars that are now outside the new bounds
	for _, star := range sf.Stars {
		if star == nil {
			continue
		}

		if star.X > newBounds.Width {
			star.X = rand.Float32() * newBounds.Width
		}
		if star.Y > newBounds.Height {
			star.Y = rand.Float32() * newBounds.Height
		}

		// Update visual position
		star.Visual.Move(fyne.NewPos(star.X-star.Size/2, star.Y-star.Size/2))
	}
}

// SetTravelSpeed allows dynamic adjustment of travel speed
func (sf *StarField) SetTravelSpeed(speed float32) {
	sf.TravelSpeed = speed
}

// SetTravelDirection allows dynamic adjustment of travel direction
func (sf *StarField) SetTravelDirection(angle float32) {
	sf.TravelAngle = angle
}

// GetStarFieldInfo returns information about the star field composition
func (sf *StarField) GetStarFieldInfo() map[StarType]int {
	counts := make(map[StarType]int)
	for _, star := range sf.Stars {
		if star != nil {
			counts[star.StarType]++
		}
	}
	return counts
}