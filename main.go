package main

import (
	"flag"
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
	spring     harmonica.Spring
	position   float64
	velocity   float64
	target     float64
	cycle      int
	localizer  *Localizer
}

func NewTerminalGym(localizer *Localizer) *TerminalGym {
	return &TerminalGym{
		spring:    harmonica.NewSpring(harmonica.FPS(fps), angularFreq, dampingRatio),
		position:  0.0,
		velocity:  0.0,
		target:    0.0,
		cycle:     0,
		localizer: localizer,
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
		return tg.localizer.T("squeeze_instruction")
	case 2, 3:
		return tg.localizer.T("lift_instruction")
	}
	return ""
}

func (tg *TerminalGym) render() {
	tg.clearScreen()
	
	// Title
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("                    " + tg.localizer.T("title"))
	fmt.Println("              " + tg.localizer.T("subtitle"))
	fmt.Println(strings.Repeat("=", 60) + "\n")
	
	// Instructions
	instruction := tg.getInstructions()
	padding := (60 - len(instruction)) / 2
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s\n\n", strings.Repeat(" ", padding), instruction)
	
	// Animation area
	fmt.Println("\n" + strings.Repeat(" ", 25) + tg.localizer.T("watch_follow"))
	fmt.Println()
	
	// Render the animated butt
	tg.renderButt()
	
	// Exercise counter and tips
	fmt.Printf("\n\n%s" + tg.localizer.Tf("rep_counter", tg.cycle/2+1) + "\n", strings.Repeat(" ", 25))
	
	// Tips
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println(tg.localizer.T("tips_header"))
	fmt.Println(tg.localizer.T("tip_follow_rhythm"))
	fmt.Println(tg.localizer.T("tip_squeeze"))
	fmt.Println(tg.localizer.T("tip_lift"))
	fmt.Println(tg.localizer.T("tip_core"))
	fmt.Println(tg.localizer.T("tip_exit"))
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
			fmt.Println("\n" + tg.localizer.T("workout_complete"))
			fmt.Println(tg.localizer.T("keep_work") + "\n")
			return
		case <-ticker.C:
			tg.update()
			tg.render()
		}
	}
}

func main() {
	// Parse command line arguments
	lang := flag.String("lang", "en", "Language (en/zh)")
	help := flag.Bool("help", false, "Show help")
	flag.Parse()
	
	// Initialize localizer
	localizer, err := NewLocalizer(*lang)
	if err != nil {
		fmt.Printf("Error initializing localizer: %v\n", err)
		fmt.Println("Falling back to English...")
		localizer, _ = NewLocalizer("en")
	}
	
	if *help {
		fmt.Println(localizer.T("language_help"))
		return
	}
	
	// Hide cursor for better animation experience
	fmt.Print("\033[?25l")
	defer fmt.Print("\033[?25h") // Show cursor on exit
	
	gym := NewTerminalGym(localizer)
	
	// Welcome message
	fmt.Print("\033[H\033[2J")
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("                 " + localizer.T("welcome_title"))
	fmt.Println("                    " + localizer.T("welcome_subtitle"))
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n" + localizer.T("starting_countdown"))
	fmt.Println(localizer.T("prepare_message"))
	
	// Countdown
	for i := 3; i > 0; i-- {
		time.Sleep(time.Second)
		fmt.Printf("\r" + localizer.Tf("starting_in", i))
	}
	fmt.Println("\n\n" + localizer.T("lets_begin"))
	time.Sleep(time.Second)
	
	gym.run()
}