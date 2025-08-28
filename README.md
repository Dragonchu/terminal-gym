# 🍑 Terminal Gym - Buttock Lifting Exercise Guide

A fun and interactive terminal application that guides you through buttock lifting exercises with animated ASCII art!

## Features

- 🎬 Animated ASCII art butt that contracts and expands
- 🏋️ Real-time exercise guidance and instructions
- 💪 Physics-based smooth animations using Harmonica library
- 📊 Exercise rep counter
- 🎯 Visual cues to help you follow the exercise rhythm
- ⌨️ Simple keyboard controls (Ctrl+C to exit)
- 🌐 Multi-language support (English and Chinese)
- 🔄 Easy language switching via command-line arguments

## Installation

Make sure you have Go installed on your system, then:

```bash
git clone https://github.com/Dragonchu/terminal-gym.git
cd terminal-gym
go mod download
go build -o terminal-gym main.go i18n.go
```

## Usage

Run the terminal gym:

```bash
./terminal-gym
```

Or run directly with Go:

```bash
go run main.go i18n.go i18n.go
```

### Language Support

The application supports both English and Chinese. You can switch languages using command-line arguments:

```bash
# Run in English (default)
./terminal-gym --lang=en

# Run in Chinese
./terminal-gym --lang=zh

# Show help
./terminal-gym --help
```

## How to Exercise

1. **Stand up** and get into position
2. **Watch the animation** - the butt will contract and expand
3. **Follow the rhythm**:
   - When the butt contracts → Squeeze your glutes
   - When the butt expands → Lift your buttocks
4. **Keep your core engaged** throughout the exercise
5. Press **Ctrl+C** when you're done

## Project Structure

```
terminal-gym/
├── main.go           # Main application logic and animation
├── i18n.go          # Internationalization support
├── locales/         # Language files
│   ├── en.json      # English translations
│   └── zh.json      # Chinese translations
├── go.mod           # Go module dependencies
└── README.md        # This file
```

## Dependencies

- [Harmonica](https://github.com/charmbracelet/harmonica) - Physics-based animation library

## Contributing

Feel free to contribute improvements, new exercises, or better ASCII art!

## License

This project is open source. Exercise responsibly! 💪
