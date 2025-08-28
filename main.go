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

// Enhanced ASCII art for different butt states with more detail
// Each state contains multiple lines representing a single frame
var buttStates = [][]string{
	// State 0: Fully contracted state - tight muscle definition
	{
		`    ╭─────╮    `,
		`   ╱  ╭─╮  ╲   `,
		`  ╱  ╱   ╲  ╲  `,
		` ╱  ╱  ●  ╲  ╲ `,
		`╱  ╱       ╲  ╲`,
		`╲  ╲       ╱  ╱`,
		` ╲  ╲     ╱  ╱ `,
		`  ╲  ╲___╱  ╱  `,
		`   ╲_______╱   `,
	},
	
	// State 1: Slight expansion - muscles beginning to engage
	{
		`     ╭─────╮     `,
		`   ╭─╱  ╭─╮  ╲─╮   `,
		`  ╱  ╱  ╱   ╲  ╲  ╲  `,
		` ╱  ╱  ╱  ●  ╲  ╲  ╲ `,
		`╱  ╱  ╱       ╲  ╲  ╲`,
		`╲  ╲  ╲       ╱  ╱  ╱`,
		` ╲  ╲  ╲     ╱  ╱  ╱ `,
		`  ╲  ╲  ╲___╱  ╱  ╱  `,
		`   ╲─╲_______╱─╱   `,
	},
	
	// State 2: Medium expansion - balanced muscle tone
	{
		`      ╭──────╮      `,
		`   ╭──╱   ╭─╮   ╲──╮   `,
		`  ╱   ╱   ╱   ╲   ╲   ╲  `,
		` ╱   ╱   ╱  ●  ╲   ╲   ╲ `,
		`╱   ╱   ╱       ╲   ╲   ╲`,
		`╲   ╲   ╲       ╱   ╱   ╱`,
		` ╲   ╲   ╲     ╱   ╱   ╱ `,
		`  ╲   ╲   ╲___╱   ╱   ╱  `,
		`   ╲──╲_________╱──╱   `,
	},
	
	// State 3: Full expansion - muscles fully engaged
	{
		`       ╭───────╮       `,
		`   ╭───╱    ╭─╮    ╲───╮   `,
		`  ╱    ╱    ╱   ╲    ╲    ╲  `,
		` ╱    ╱    ╱  ●  ╲    ╲    ╲ `,
		`╱    ╱    ╱       ╲    ╲    ╲`,
		`╲    ╲    ╲       ╱    ╱    ╱`,
		` ╲    ╲    ╲     ╱    ╱    ╱ `,
		`  ╲    ╲    ╲___╱    ╱    ╱  `,
		`   ╲───╲___________╱───╱   `,
	},
	
	// State 4: Maximum expansion - peak muscle activation
	{
		`        ╭────────╮        `,
		`   ╭────╱     ╭─╮     ╲────╮   `,
		`  ╱     ╱     ╱   ╲     ╲     ╲  `,
		` ╱     ╱     ╱  ●  ╲     ╲     ╲ `,
		`╱     ╱     ╱       ╲     ╲     ╲`,
		`╲     ╲     ╲       ╱     ╱     ╱`,
		` ╲     ╲     ╲     ╱     ╱     ╱ `,
		`  ╲     ╲     ╲___╱     ╱     ╱  `,
		`   ╲────╲_____________╱────╱   `,
	},
}

type TerminalGym struct {
	// Main animation spring
	mainSpring     harmonica.Spring
	mainPosition   float64
	mainVelocity   float64
	mainTarget     float64
	
	// Secondary springs for subtle effects
	leftSpring     harmonica.Spring
	leftPosition   float64
	leftVelocity   float64
	leftTarget     float64
	
	rightSpring    harmonica.Spring
	rightPosition  float64
	rightVelocity  float64
	rightTarget    float64
	
	// Breathing/micro-movement spring
	breathSpring   harmonica.Spring
	breathPosition float64
	breathVelocity float64
	breathTarget   float64
	
	// Muscle tension spring for definition
	tensionSpring  harmonica.Spring
	tensionPosition float64
	tensionVelocity float64
	tensionTarget   float64
	
	cycle          int
	localizer      *Localizer
	frameCount     int64  // For time-based effects
}

func NewTerminalGym(localizer *Localizer) *TerminalGym {
	return &TerminalGym{
		// Main spring for primary animation
		mainSpring:    harmonica.NewSpring(harmonica.FPS(fps), angularFreq, dampingRatio),
		mainPosition:  0.0,
		mainVelocity:  0.0,
		mainTarget:    0.0,
		
		// Left cheek with slightly different characteristics
		leftSpring:    harmonica.NewSpring(harmonica.FPS(fps), angularFreq*1.1, dampingRatio*0.9),
		leftPosition:  0.0,
		leftVelocity:  0.0,
		leftTarget:    0.0,
		
		// Right cheek with slightly different characteristics  
		rightSpring:   harmonica.NewSpring(harmonica.FPS(fps), angularFreq*0.9, dampingRatio*1.1),
		rightPosition: 0.0,
		rightVelocity: 0.0,
		rightTarget:   0.0,
		
		// Breathing effect - slower, more subtle
		breathSpring:  harmonica.NewSpring(harmonica.FPS(fps), 1.5, 0.8),
		breathPosition: 0.0,
		breathVelocity: 0.0,
		breathTarget:   0.0,
		
		// Muscle tension - faster response, higher damping
		tensionSpring: harmonica.NewSpring(harmonica.FPS(fps), angularFreq*2.0, dampingRatio*2.0),
		tensionPosition: 0.0,
		tensionVelocity: 0.0,
		tensionTarget:   0.0,
		
		cycle:         0,
		localizer:     localizer,
		frameCount:    0,
	}
}

func (tg *TerminalGym) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (tg *TerminalGym) renderButt() {
	// Calculate the base animation state using main spring
	normalizedPos := (tg.mainPosition + animationRange) / (2 * animationRange)
	if normalizedPos < 0 {
		normalizedPos = 0
	}
	if normalizedPos > 1 {
		normalizedPos = 1
	}
	
	// Base state selection
	baseStateIndex := int(normalizedPos * float64(len(buttStates)-1))
	if baseStateIndex >= len(buttStates) {
		baseStateIndex = len(buttStates) - 1
	}
	
	// Calculate subtle asymmetry from left/right springs
	leftOffset := int(tg.leftPosition * 0.3)   // Subtle left adjustment
	rightOffset := int(tg.rightPosition * 0.3) // Subtle right adjustment
	
	// Calculate breathing micro-movement
	breathOffset := int(tg.breathPosition * 0.5)
	
	// Calculate muscle tension effect
	tensionIntensity := (tg.tensionPosition + animationRange) / (2 * animationRange)
	if tensionIntensity < 0 {
		tensionIntensity = 0
	}
	if tensionIntensity > 1 {
		tensionIntensity = 1
	}
	
	// Dynamic padding based on breathing and asymmetry
	basePadding := 15
	dynamicPadding := basePadding + breathOffset + leftOffset - rightOffset
	if dynamicPadding < 5 {
		dynamicPadding = 5
	}
	if dynamicPadding > 25 {
		dynamicPadding = 25
	}
	
	// Render each line of the butt with subtle modifications
	buttLines := buttStates[baseStateIndex]
	for _, line := range buttLines {
		
		// Apply tension-based character substitution for more defined look
		if tensionIntensity > 0.7 {
			line = strings.ReplaceAll(line, "╱", "╱")  // Keep sharp angles
			line = strings.ReplaceAll(line, "╲", "╲")  // Keep sharp angles
		} else if tensionIntensity < 0.3 {
			// Softer look for relaxed state
			line = strings.ReplaceAll(line, "╭", "╭")
			line = strings.ReplaceAll(line, "╮", "╮")
		}
		
		// Apply asymmetric padding for left/right variation
		leftPad := dynamicPadding + (leftOffset - rightOffset)/2
		if leftPad < 0 {
			leftPad = 0
		}
		
		padding := strings.Repeat(" ", leftPad)
		
		// Add subtle rotation effect based on spring differences
		rotationEffect := ""
		if abs(tg.leftPosition-tg.rightPosition) > 1.0 {
			if tg.leftPosition > tg.rightPosition {
				rotationEffect = " ↗" // Slight tilt indicator
			} else {
				rotationEffect = " ↖" // Slight tilt indicator  
			}
		}
		
		fmt.Printf("%s%s%s\n", padding, line, rotationEffect)
	}
	
	// Add subtle muscle activation indicators
	if tensionIntensity > 0.8 {
		indicatorPadding := strings.Repeat(" ", dynamicPadding+8)
		fmt.Printf("%s💪 Peak Activation 💪\n", indicatorPadding)
	} else if tensionIntensity > 0.5 {
		indicatorPadding := strings.Repeat(" ", dynamicPadding+10)
		fmt.Printf("%s⚡ Engaged ⚡\n", indicatorPadding)
	}
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
	tg.frameCount++
	
	// Update main spring physics
	tg.mainPosition, tg.mainVelocity = tg.mainSpring.Update(tg.mainPosition, tg.mainVelocity, tg.mainTarget)
	
	// Update left cheek spring with slight delay and variation
	leftTargetVariation := tg.mainTarget + sin(float64(tg.frameCount)*0.02)*0.5
	tg.leftPosition, tg.leftVelocity = tg.leftSpring.Update(tg.leftPosition, tg.leftVelocity, leftTargetVariation)
	
	// Update right cheek spring with different delay and variation
	rightTargetVariation := tg.mainTarget + sin(float64(tg.frameCount)*0.018)*0.4
	tg.rightPosition, tg.rightVelocity = tg.rightSpring.Update(tg.rightPosition, tg.rightVelocity, rightTargetVariation)
	
	// Update breathing spring with slow oscillation
	breathingCycle := sin(float64(tg.frameCount) * 0.01) * 2.0
	tg.breathPosition, tg.breathVelocity = tg.breathSpring.Update(tg.breathPosition, tg.breathVelocity, breathingCycle)
	
	// Update tension spring - follows main target but with different characteristics
	tensionTarget := tg.mainTarget * 1.2 // Slightly more intense
	tg.tensionPosition, tg.tensionVelocity = tg.tensionSpring.Update(tg.tensionPosition, tg.tensionVelocity, tensionTarget)
	
	// Check if we need to change target (cycle between contract and expand)
	if tg.hasReachedMainTarget() {
		tg.cycle++
		if tg.cycle%2 == 0 {
			tg.mainTarget = -animationRange // Contract
		} else {
			tg.mainTarget = animationRange  // Expand
		}
	}
}

func (tg *TerminalGym) hasReachedMainTarget() bool {
	threshold := 0.5
	return abs(tg.mainPosition-tg.mainTarget) < threshold && abs(tg.mainVelocity) < threshold
}

// Add sin function for smooth oscillations
func sin(x float64) float64 {
	// Simple sine approximation for smooth oscillations
	// Using Taylor series approximation for efficiency
	x = x - float64(int(x/(2*3.14159)))*2*3.14159 // Normalize to 0-2π range
	if x < 0 {
		x += 2 * 3.14159
	}
	
	// Taylor series: sin(x) ≈ x - x³/6 + x⁵/120 - x⁷/5040
	x2 := x * x
	x3 := x2 * x
	x5 := x3 * x2
	x7 := x5 * x2
	
	return x - x3/6.0 + x5/120.0 - x7/5040.0
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
	tg.mainTarget = -animationRange
	
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