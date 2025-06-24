# Bouncing Balls Simulator

A physics simulation featuring bouncing balls with realistic collision physics, lightning effects, particle trails, and an AI-controlled human character that tries to avoid the balls.

## Features

### Ball Physics
- **Realistic Collision Detection**: Balls bounce off walls and each other with proper physics
- **Mass-Based Collisions**: Different sized balls have different masses affecting collision outcomes
- **Elastic Collisions**: Energy is conserved during ball-to-ball collisions
- **Particle Trails**: Each ball leaves a glowing trail of particles that fade over time

### Visual Effects
- **Lightning Effects**: Ball collisions create spectacular lightning bolts with flickering animation
- **Particle Trails**: Glowing trails follow each ball with fading alpha and decreasing size
- **Smooth Animation**: 60 FPS animation for fluid motion

### Human AI Character
- **Intelligent Avoidance**: Human character uses predictive algorithms to avoid approaching balls
- **Panic Mode**: Speed increases when in extreme danger
- **Collision Detection**: Human explodes when hit by balls
- **Respawn System**: Human respawns at safe locations after being hit
- **Death Tracking**: Keeps count of how many times the human has been hit

### Interactive Controls
- **Start/Stop Animation**: Control the simulation
- **Speed Control**: Speed up or slow down the balls
- **Color Changes**: Cycle through different ball colors
- **Mass Information**: Display ball masses
- **Human Toggle**: Show/hide the human character
- **Death Counter**: View human death statistics
- **Reset Function**: Reset all elements to initial state

## Project Structure

Following Go best practices, the project is organized as follows:

```
bouncing-balls/
├── cmd/
│   └── bouncing-balls/
│       └── main.go              # Application entry point
├── pkg/
│   ├── physics/
│   │   ├── ball.go              # Ball physics and collision system
│   │   └── human.go             # Human AI and behavior system
│   ├── effects/
│   │   └── lightning.go         # Lightning effects system
│   └── ui/
│       └── app.go               # UI application and controls
├── build/                       # Build artifacts (created by make)
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── Makefile                     # Build automation
└── README.md                    # This file
```

### Key Components

1. **Physics Package** (`pkg/physics/`)
   - `ball.go`: Ball physics, collisions, particle trails
   - `human.go`: Human AI behavior, explosion effects, respawn logic

2. **Effects Package** (`pkg/effects/`)
   - `lightning.go`: Lightning bolt effects and animations

3. **UI Package** (`pkg/ui/`)
   - `app.go`: Application setup, UI controls, and main game loop

4. **Command** (`cmd/bouncing-balls/`)
   - `main.go`: Application entry point

## Requirements

- Go 1.21 or later
- Fyne GUI toolkit (automatically installed via go.mod)
- macOS (optimized for macOS but should work on other platforms)
- Make (for build automation)

## Quick Start

### Using Make (Recommended)

```bash
# Show all available commands
make help

# Quick development cycle: clean, format, vet, build, and run
make dev

# Just run the application
make run

# Build the application
make build

# Run the built binary
make run-binary
```

### Manual Installation

```bash
# Clone the repository
git clone <repository-url>
cd bouncing-balls

# Install dependencies
go mod tidy

# Run the application
go run ./cmd/bouncing-balls

# Or build and run
go build -o bouncing-balls ./cmd/bouncing-balls
./bouncing-balls
```

## Build System

The project includes a comprehensive Makefile with the following targets:

### Development Commands
- `make run` - Run the application without building
- `make build` - Build the application binary
- `make clean` - Clean build artifacts and cache
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make test` - Run all tests
- `make dev` - Quick development cycle

### Production Commands
- `make prod` - Production build with tests and packaging
- `make package` - Create distributable package
- `make build-all` - Cross-compile for multiple platforms
- `make install` - Install binary to GOPATH/bin

### Utility Commands
- `make help` - Show all available commands
- `make info` - Show project information
- `make deps` - Download dependencies
- `make tidy` - Tidy and verify modules

## Usage

1. Run the program using `make run` or `make dev`
2. Balls will start bouncing automatically
3. Use the control buttons at the bottom to interact with the simulation:
   - **▶️ Start All**: Start ball animation
   - **⏸️ Stop All**: Pause ball animation
   - **🎨 Change Colors**: Cycle through different ball colors
   - **⚡ Speed Up**: Increase ball velocities
   - **🐌 Slow Down**: Decrease ball velocities
   - **⚖️ Show Masses**: Display ball mass information in console
   - **🏃 Toggle Human**: Show/hide the human character
   - **💀 Death Count**: Display human death count in console
   - **🔄 Reset All**: Reset everything to initial state
   - **❌ Quit**: Exit the application

## Physics Details

### Ball Collisions
- Uses elastic collision formulas with momentum conservation
- Mass calculated as π × radius² (area-based mass)
- Larger balls have more mass and affect smaller balls more dramatically
- Collision separation prevents balls from sticking together

### Human AI
- **Predictive Avoidance**: Calculates ball positions 5, 10, 15, and 20 frames ahead
- **Danger Zones**: 120-unit detection radius around each ball
- **Speed Scaling**:
  - Base speed: 4.5 units/frame
  - Moderate danger: 1.5× speed boost
  - Extreme danger: 2× speed boost
- **Quadratic Avoidance**: Avoidance strength increases quadratically with proximity

### Visual Effects
- **Lightning**: 8-segment jagged bolts lasting 300ms with flickering
- **Particle Trails**: 10 particles per ball with fading alpha and decreasing size
- **Explosions**: 12 colored particles expanding radially on human death

## Performance

- Runs at 60 FPS (16ms frame time)
- Efficient collision detection algorithms
- Optimized rendering with Fyne's native graphics

## Development

The code follows Go best practices with:
- Clear package separation and modular design
- Exported functions with proper documentation
- Consistent naming conventions
- Clean imports and dependencies
- Comprehensive build system with Make

### Adding New Features

1. **Physics Features**: Add to `pkg/physics/`
2. **Visual Effects**: Add to `pkg/effects/`
3. **UI Components**: Add to `pkg/ui/`
4. **New Commands**: Add to `cmd/` directory

## Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Run all quality checks
make check
```

## Distribution

```bash
# Create a distributable package
make package

# Cross-compile for multiple platforms
make build-all

# Install system-wide
make install
```

## Future Enhancements

Potential additions could include:
- Sound effects for collisions and explosions
- Multiple human characters
- Adjustable physics parameters
- Save/load simulation states
- Ball spawning/removal during runtime
- Configuration file support
- Performance metrics and profiling