package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/charmbracelet/harmonica"
)

const (
	// Animation parameters
	fps            = 30
	angularFreq    = 4.0
	dampingRatio   = 0.3
	animationRange = 8.0
)

// ASCII art for different butt states
var buttStates = []string{
	// Contracted state (smaller)
	`    ( . Y . )    `,
	
	// Slightly expanded
	`   (  . Y .  )   `,
	
	// Medium expansion
	`  (   . Y .   )  `,
	
	// Fully expanded
	` (    . Y .    ) `,
	
	// Maximum expansion
	`(     . Y .     )`,
}

type TerminalGym struct {
	spring   harmonica.Spring
	position float64
	velocity float64
	target   float64
	cycle    int
}

func NewTerminalGym() *TerminalGym {
	return &TerminalGym{
		spring:   harmonica.NewSpring(harmonica.FPS(fps), angularFreq, dampingRatio),
		position: 0.0,
		velocity: 0.0,
		target:   0.0,
		cycle:    0,
	}
}

func (tg *TerminalGym) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (tg *TerminalGym) renderButt() {
	// Map position to butt state
	normalizedPos := (tg.position + animationRange) / (2 * animationRange)
	if normalizedPos < 0 {
		normalizedPos = 0
	}
	if normalizedPos > 1 {
		normalizedPos = 1
	}
	
	stateIndex := int(normalizedPos * float64(len(buttStates)-1))
	if stateIndex >= len(buttStates) {
		stateIndex = len(buttStates) - 1
	}
	
	// Center the butt on screen
	padding := strings.Repeat(" ", 20)
	fmt.Printf("%s%s\n", padding, buttStates[stateIndex])
}

func (tg *TerminalGym) getInstructions() string {
	phase := tg.cycle % 4
	switch phase {
	case 0, 1:
		return "🏋️  SQUEEZE YOUR GLUTES! Contract those muscles! 🏋️"
	case 2, 3:
		return "🚀  LIFT YOUR BUTTOCKS! Push up and engage! 🚀"
	}
	return ""
}

func (tg *TerminalGym) render() {
	tg.clearScreen()
	
	// Title
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("                    🍑 TERMINAL GYM 🍑")
	fmt.Println("              Buttock Lifting Exercise Guide")
	fmt.Println(strings.Repeat("=", 60) + "\n")
	
	// Instructions
	instruction := tg.getInstructions()
	padding := (60 - len(instruction)) / 2
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s\n\n", strings.Repeat(" ", padding), instruction)
	
	// Animation area
	fmt.Println("\n" + strings.Repeat(" ", 25) + "👀 WATCH AND FOLLOW 👀")
	fmt.Println()
	
	// Render the animated butt
	tg.renderButt()
	
	// Exercise counter and tips
	fmt.Printf("\n\n%sRep: %d\n", strings.Repeat(" ", 25), tg.cycle/2+1)
	
	// Tips
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println("💡 Tips:")
	fmt.Println("   • Follow the animation rhythm")
	fmt.Println("   • Squeeze when the butt contracts")
	fmt.Println("   • Lift when the butt expands")
	fmt.Println("   • Keep your core engaged")
	fmt.Println("   • Press Ctrl+C to exit")
	fmt.Println(strings.Repeat("-", 60))
}

func (tg *TerminalGym) update() {
	// Update spring physics
	tg.position, tg.velocity = tg.spring.Update(tg.position, tg.velocity, tg.target)
	
	// Check if we need to change target (cycle between contract and expand)
	if tg.hasReachedTarget() {
		tg.cycle++
		if tg.cycle%2 == 0 {
			tg.target = -animationRange // Contract
		} else {
			tg.target = animationRange  // Expand
		}
	}
}

func (tg *TerminalGym) hasReachedTarget() bool {
	threshold := 0.5
	return abs(tg.position-tg.target) < threshold && abs(tg.velocity) < threshold
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func (tg *TerminalGym) run() {
	// Set up signal handling for graceful exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	// Start with contraction
	tg.target = -animationRange
	
	// Animation loop
	ticker := time.NewTicker(time.Second / fps)
	defer ticker.Stop()
	
	for {
		select {
		case <-c:
			tg.clearScreen()
			fmt.Println("\n🎉 Great workout! Your glutes thank you! 🎉")
			fmt.Println("💪 Keep up the good work! 💪\n")
			return
		case <-ticker.C:
			tg.update()
			tg.render()
		}
	}
}

func main() {
	// Hide cursor for better animation experience
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit
	
	gym := NewTerminalGym()
	
	// Welcome message
	fmt.Print("\033[H\033[2J")
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("                 🏋️  WELCOME TO TERMINAL GYM! 🏋️")
	fmt.Println("                    Get Ready to Lift! 🍑")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n⏰ Starting in 3 seconds... Get into position!")
	fmt.Println("🧘 Stand up, engage your core, and prepare your glutes!")
	
	// Countdown
	for i := 3; i > 0; i-- {
		time.Sleep(time.Second)
		fmt.Printf("\r🚀 Starting in %d... ", i)
	}
	fmt.Println("\n\n🎬 Let's begin!")
	time.Sleep(time.Second)
	
	gym.run()
}