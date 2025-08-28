# 🏋️ Terminal Gym - Exercise & Meditation Guide 🧘

A fun and interactive terminal application that guides you through various exercises with animated ASCII art! Choose between strength training and meditation exercises.

## Features

- 🏋️ **Multiple Exercise Types**: Choose between strength training and meditation
- 🍑 **Buttock Lifting**: Animated ASCII art that contracts and expands for muscle training
- 🧘 **Deep Breathing Meditation**: Guided 4-7-8 breathing technique with lung visualization
- 🎬 **Animated ASCII Art**: Physics-based smooth animations using Harmonica library
- 📊 **Progress Tracking**: Exercise rep counter and breath cycle counter
- 🎯 **Visual Cues**: Real-time guidance to help you follow the exercise rhythm
- ⌨️ **Simple Controls**: Easy keyboard navigation and Ctrl+C to exit
- 🌐 **Multi-language Support**: English and Chinese localization
- 🔄 **Easy Language Switching**: Command-line language selection

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

## How to Use

### Exercise Selection
1. **Run the application** and choose your exercise type
2. **Option 1: Buttock Lifting** - Strength training exercise
3. **Option 2: Deep Breathing Meditation** - Relaxation and mindfulness

### Buttock Lifting Exercise
1. **Stand up** and get into position
2. **Watch the animation** - the butt will contract and expand
3. **Follow the rhythm**:
   - When the butt contracts → Squeeze your glutes
   - When the butt expands → Lift your buttocks
4. **Keep your core engaged** throughout the exercise

### Deep Breathing Meditation
1. **Sit or lie down** in a comfortable position
2. **Watch the lung animation** - it will expand and contract
3. **Follow the 4-7-8 breathing technique**:
   - Inhale for 4 seconds (lung expands)
   - Hold for 7 seconds (lung stays expanded)
   - Exhale for 8 seconds (lung contracts)
   - Pause for 2 seconds (brief rest)
4. **Focus on your breath** and let go of thoughts

### General Controls
- Press **Ctrl+C** when you're done with any exercise

## Project Structure

```
terminal-gym/
├── main.go           # Main application with exercise selection and implementations
├── i18n.go          # Internationalization support
├── locales/         # Language files
│   ├── en.json      # English translations
│   └── zh.json      # Chinese translations
├── go.mod           # Go module dependencies
└── README.md        # This file
```

## Exercise Types

### 🏋️ Strength Training
- **Buttock Lifting**: Interactive glute strengthening exercise with real-time visual feedback

### 🧘 Meditation & Wellness  
- **Deep Breathing**: 4-7-8 breathing technique for stress relief and mindfulness

## Dependencies

- [Harmonica](https://github.com/charmbracelet/harmonica) - Physics-based animation library

## Contributing

Feel free to contribute improvements, new exercises, or better ASCII art!

## License

This project is open source. Exercise responsibly! 💪
