# 🚀 Eyeball Space Travel Simulator

**Flying Through the Galaxy with Realistic Astrophysics and Strategic Gameplay**

A sophisticated Go-based space simulation featuring realistic eyeball entities, stellar classification systems, and intelligent AI companions. Built using the Fyne GUI framework with advanced physics and visual effects.

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Fyne](https://img.shields.io/badge/GUI-Fyne-blue?style=for-the-badge)
![Physics](https://img.shields.io/badge/Physics-Realistic-green?style=for-the-badge)

## ✨ Features

### 👁️ Realistic Eyeball Entities
- **Anatomical Accuracy**: White sclera, colored iris, black pupil with proper proportions
- **Bloodshot Veins**: 6 semi-transparent red vessels for authentic appearance
- **Dynamic Iris Tracking**: Eyes follow the human character with 50% movement range
- **AI LLM Names**: Each eyeball displays names of popular AI models (GPT-4, Claude, Gemini, etc.)
- **Collision Physics**: Mass-based elastic collisions with jiggle effects

### ⭐ Realistic Stellar Environment
- **8 Stellar Classifications**: 
  - Red Dwarfs (76%), Orange Dwarfs (12%), Yellow Dwarfs (7.6%)
  - Blue-White Dwarfs (3%), Red Giants (0.8%), Blue Giants (0.003%)
  - White Dwarfs (0.3%), Neutron Stars (0.1%)
- **Galactic Distribution**: Non-uniform density with exponential falloff from galactic center
- **Spiral Arm Enhancement**: Mathematical modeling of galactic structure
- **Parallax Effects**: Distance-based star movement for space travel immersion
- **Advanced Twinkling**: Star-type-specific luminosity variations
- **Dynamic Regeneration**: 400 stars with seamless edge regeneration

### 🐉 Strategic Dragon Protector
- **Movement Prediction**: Tracks human velocity to anticipate direction
- **Strategic Deflection**: Deflects eyeballs opposite to human movement
- **Threat Assessment**: Prioritizes balls moving toward human within 150-pixel radius
- **Mass-Based Physics**: Dragon mass = 2x largest eyeball mass (minimum 1000 units)
- **Collision Effects**: Shrinks eyeballs to half size and reduces velocity
- **Recovery Animations**: Drift and spin cycles for realistic behavior

### 🎮 Advanced Human Character
- **Intelligent Respawn**: Grid-based algorithm finds safest position from all eyeballs
- **Bullet System**: Strategic repulsion forces push eyeballs away
- **Auto-Targeting**: Faces and shoots at closest threatening eyeball
- **Collision Avoidance**: AI-driven movement away from approaching threats
- **Explosion Effects**: Particle system with respawn timer

### 🎯 Strategic Gameplay
- **Bullet Repulsion**: Fixed physics bug - bullets now properly repel eyeballs
- **Speed Optimization**: Reduced velocities for more controlled, strategic play
- **Smart Dragon**: Clears path ahead of human movement
- **Safe Respawn**: Maximizes distance from all threats
- **Visual Feedback**: Jiggle effects, particle explosions, and trail systems

## 🛠️ Technical Implementation

### Architecture
```
cmd/bouncing-balls/     # Main application entry point
pkg/
├── physics/           # Core physics and entity logic
│   ├── ball.go       # Eyeball entities with iris tracking
│   ├── dragon.go     # Strategic AI protector
│   ├── human.go      # Player character with smart respawn
│   └── starfield.go  # Stellar classification system
└── ui/               # User interface and rendering
    └── app.go        # Main application loop and controls
```

### Key Technologies
- **Go 1.23+**: High-performance concurrent programming
- **Fyne v2**: Cross-platform GUI framework
- **Custom Physics**: Mass-based elastic collisions
- **Real-time Rendering**: 60 FPS animation with layered graphics
- **Mathematical Modeling**: Astrophysics-based stellar distribution

### Performance Features
- **Parallel Processing**: Concurrent updates for all entities
- **Efficient Collision Detection**: Optimized distance calculations
- **Memory Management**: Object pooling for particles and trails
- **Smooth Animation**: Interpolated movements and effects

## 🚀 Installation & Usage

### Prerequisites
- Go 1.23 or later
- Fyne dependencies (automatically handled by Go modules)

### Quick Start
```bash
# Clone the repository
git clone https://github.com/atyronesmith/bouncing-balls.git
cd bouncing-balls

# Build the application
go build -o bouncing-balls cmd/bouncing-balls/main.go

# Run the simulator
./bouncing-balls
```

### Controls
- **Arrow Keys**: Move human character
- **Auto-Shooting**: Character automatically targets closest eyeball
- **Mouse**: Interact with UI controls
- **Buttons**:
  - ▶️ Start All - Begin animation
  - ⏸️ Stop All - Pause simulation  
  - 🎨 Change Colors - Cycle eyeball iris colors
  - 🔄 Reset All - Return to initial state
  - ❌ Quit - Exit application

## 🎨 Vibe Coding Philosophy

This project exemplifies **"vibe coding"** - a development approach that prioritizes:

### 🌊 Flow State Development
- **Iterative Enhancement**: Features evolved organically through conversation
- **Creative Exploration**: Started as simple bouncing balls, transformed into space simulation
- **Intuitive Design**: Each addition felt natural and enhanced the overall experience
- **Emergent Complexity**: Sophisticated systems arose from simple interactions

### 🤝 Collaborative Creation
- **Human-AI Partnership**: Ideas flowed between human creativity and AI implementation
- **Real-time Feedback**: Immediate testing and refinement of each feature
- **Shared Vision**: Both participants contributed to the evolving concept
- **Iterative Improvement**: "What if we added..." led to continuous enhancement

### ✨ Key Vibe Coding Principles Demonstrated
1. **Start Simple**: Basic bouncing balls → Complex space simulation
2. **Follow Inspiration**: Each feature suggested the next logical enhancement
3. **Embrace Serendipity**: Bug discoveries led to feature improvements
4. **Maintain Playfulness**: Fun factor guided technical decisions
5. **Continuous Evolution**: No rigid plan, just organic growth

### 🎯 Results of Vibe Coding
- **Rich Feature Set**: Far exceeded initial scope
- **Cohesive Design**: All systems work harmoniously together
- **Technical Excellence**: Clean, maintainable, well-documented code
- **Engaging Experience**: Genuinely fun and visually appealing
- **Learning Journey**: Both participants discovered new possibilities

## 🎮 Gameplay Experience

The simulator creates an immersive space travel experience where you navigate through a realistic galaxy populated by watchful eyeball entities. Your dragon companion intelligently protects you by predicting your movement and clearing threats from your path. The stellar background provides authentic astronomical ambiance with proper star classifications and parallax effects.

Strategic elements include:
- **Bullet Herding**: Use repulsion shots to guide eyeballs away
- **Dragon Coordination**: Move predictably to help your protector
- **Safe Positioning**: Respawn system gives you breathing room
- **Visual Tracking**: Eyeballs watch your every move with realistic iris movement

## 🔬 Physics & Algorithms

### Collision Detection
- Circular collision detection with radius-based overlap
- Mass-based separation using momentum conservation
- Elastic collision response with energy damping

### Stellar Distribution
```go
// Galactic density function
density = baseRate * exp(-distanceFromCenter/scaleLength)
// Spiral arm enhancement
spiralBonus = 1 + spiralStrength * cos(spiralPhase)
```

### Strategic AI
- Grid-based pathfinding for optimal respawn locations
- Velocity prediction for dragon interception
- Threat assessment using distance and velocity vectors

## 🤝 Contributing

This project welcomes contributions in the spirit of vibe coding:

1. **Explore and Experiment**: Try the simulator and see what inspires you
2. **Follow Your Intuition**: If something feels like it would be cool, it probably would be
3. **Start Conversations**: Open issues to discuss ideas, not just bugs
4. **Iterate Together**: Small improvements that build on each other
5. **Maintain the Vibe**: Keep the playful, exploratory spirit alive

### Development Setup
```bash
# Fork and clone the repository
git clone https://github.com/yourusername/bouncing-balls.git
cd bouncing-balls

# Install dependencies
go mod tidy

# Run tests
go test ./...

# Start experimenting!
```

## 📜 License

This project is open source and available under the [MIT License](LICENSE).

## 🙏 Acknowledgments

- **Fyne Framework**: Excellent cross-platform GUI toolkit
- **Go Community**: Robust ecosystem and excellent tooling
- **Vibe Coding**: Philosophy of creative, collaborative development
- **Astrophysics**: Real stellar classification data for authenticity

---

*"Sometimes the best code comes not from rigid planning, but from following the vibe and seeing where it leads."*

**Built with ❤️ through vibe coding**
