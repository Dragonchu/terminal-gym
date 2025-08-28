package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

// Exercise interface
type Exercise interface {
	GetName() string
	GetCategory() string
	GetDescription() string
	Render() 
	Update()
	GetInstructions() string
	GetTips() []string
	IsComplete() bool
	Reset()
	GetCounter() string
}

// ButtockExercise represents the buttock lifting exercise
type ButtockExercise struct {
	// Base exercise properties
	Name        string
	Category    string
	Description string
	Cycle       int
	FrameCount  int64
	Localizer   *Localizer
	
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
}

// Enhanced ASCII art for different butt states with more detail
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

func NewButtockExercise(localizer *Localizer) *ButtockExercise {
	return &ButtockExercise{
		Name:        "Buttock Lifting",
		Category:    "Strength",
		Description: "Buttock lifting exercise with animated guidance",
		Cycle:       0,
		FrameCount:  0,
		Localizer:   localizer,
		
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
	}
}

func (be *ButtockExercise) GetName() string {
	return be.Name
}

func (be *ButtockExercise) GetCategory() string {
	return be.Category
}

func (be *ButtockExercise) GetDescription() string {
	return be.Description
}

func (be *ButtockExercise) renderButt() {
	// Calculate the base animation state using main spring
	normalizedPos := (be.mainPosition + animationRange) / (2 * animationRange)
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
	leftOffset := int(be.leftPosition * 0.3)   // Subtle left adjustment
	rightOffset := int(be.rightPosition * 0.3) // Subtle right adjustment
	
	// Calculate breathing micro-movement
	breathOffset := int(be.breathPosition * 0.5)
	
	// Calculate muscle tension effect
	tensionIntensity := (be.tensionPosition + animationRange) / (2 * animationRange)
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
		if abs(be.leftPosition-be.rightPosition) > 1.0 {
			if be.leftPosition > be.rightPosition {
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

func (be *ButtockExercise) Render() {
	be.renderButt()
}

func (be *ButtockExercise) Update() {
	be.FrameCount++
	
	// Update main spring physics
	be.mainPosition, be.mainVelocity = be.mainSpring.Update(be.mainPosition, be.mainVelocity, be.mainTarget)
	
	// Update left cheek spring with slight delay and variation
	leftTargetVariation := be.mainTarget + sin(float64(be.FrameCount)*0.02)*0.5
	be.leftPosition, be.leftVelocity = be.leftSpring.Update(be.leftPosition, be.leftVelocity, leftTargetVariation)
	
	// Update right cheek spring with different delay and variation
	rightTargetVariation := be.mainTarget + sin(float64(be.FrameCount)*0.018)*0.4
	be.rightPosition, be.rightVelocity = be.rightSpring.Update(be.rightPosition, be.rightVelocity, rightTargetVariation)
	
	// Update breathing spring with slow oscillation
	breathingCycle := sin(float64(be.FrameCount) * 0.01) * 2.0
	be.breathPosition, be.breathVelocity = be.breathSpring.Update(be.breathPosition, be.breathVelocity, breathingCycle)
	
	// Update tension spring - follows main target but with different characteristics
	tensionTarget := be.mainTarget * 1.2 // Slightly more intense
	be.tensionPosition, be.tensionVelocity = be.tensionSpring.Update(be.tensionPosition, be.tensionVelocity, tensionTarget)
	
	// Check if we need to change target (cycle between contract and expand)
	if be.hasReachedMainTarget() {
		be.Cycle++
		if be.Cycle%2 == 0 {
			be.mainTarget = -animationRange // Contract
		} else {
			be.mainTarget = animationRange  // Expand
		}
	}
}

func (be *ButtockExercise) hasReachedMainTarget() bool {
	threshold := 0.5
	return abs(be.mainPosition-be.mainTarget) < threshold && abs(be.mainVelocity) < threshold
}

func (be *ButtockExercise) GetInstructions() string {
	phase := be.Cycle % 4
	switch phase {
	case 0, 1:
		return be.Localizer.T("squeeze_instruction")
	case 2, 3:
		return be.Localizer.T("lift_instruction")
	}
	return ""
}

func (be *ButtockExercise) GetTips() []string {
	return []string{
		be.Localizer.T("tip_follow_rhythm"),
		be.Localizer.T("tip_squeeze"),
		be.Localizer.T("tip_lift"),
		be.Localizer.T("tip_core"),
		be.Localizer.T("tip_exit"),
	}
}

func (be *ButtockExercise) GetCounter() string {
	return be.Localizer.Tf("rep_counter", be.Cycle/2+1)
}

func (be *ButtockExercise) IsComplete() bool {
	return false // This exercise runs indefinitely until user exits
}

func (be *ButtockExercise) Reset() {
	be.Cycle = 0
	be.FrameCount = 0
	be.mainPosition = 0.0
	be.mainVelocity = 0.0
	be.mainTarget = -animationRange
	be.leftPosition = 0.0
	be.leftVelocity = 0.0
	be.rightPosition = 0.0
	be.rightVelocity = 0.0
	be.breathPosition = 0.0
	be.breathVelocity = 0.0
	be.tensionPosition = 0.0
	be.tensionVelocity = 0.0
}

// MeditationExercise represents a deep breathing meditation exercise
type MeditationExercise struct {
	// Base exercise properties
	Name        string
	Category    string
	Description string
	Cycle       int
	FrameCount  int64
	Localizer   *Localizer
	
	// Breathing animation spring
	breathSpring   harmonica.Spring
	breathPosition float64
	breathVelocity float64
	breathTarget   float64
	
	// Lung expansion spring
	lungSpring     harmonica.Spring
	lungPosition   float64
	lungVelocity   float64
	lungTarget     float64
	
	// Heart rate spring for calming effect
	heartSpring    harmonica.Spring
	heartPosition  float64
	heartVelocity  float64
	heartTarget    float64
	
	// Meditation state
	isInhaling     bool
	breathCycles   int
	phase          string // "inhale", "hold", "exhale", "pause"
	phaseTimer     int
	phaseDuration  int
}

// ASCII art for different breathing states
var breathingStates = [][]string{
	// State 0: Exhaled - lungs contracted
	{
		`           ╭─────╮           `,
		`         ╱         ╲         `,
		`       ╱    ╭───╮    ╲       `,
		`      ╱    ╱  ○  ╲    ╲      `,
		`     ╱    ╱       ╲    ╲     `,
		`    ╱    ╱    ♡    ╲    ╲    `,
		`   ╱    ╱           ╲    ╲   `,
		`  ╱    ╱             ╲    ╲  `,
		` ╱____╱               ╲____╲ `,
		`╱_____________________╲`,
	},
	
	// State 1: Slight inhale - lungs beginning to expand
	{
		`          ╭───────╮          `,
		`        ╱           ╲        `,
		`      ╱    ╭─────╮    ╲      `,
		`     ╱    ╱   ○   ╲    ╲     `,
		`    ╱    ╱         ╲    ╲    `,
		`   ╱    ╱     ♡     ╲    ╲   `,
		`  ╱    ╱             ╲    ╲  `,
		` ╱    ╱               ╲    ╲ `,
		`╱____╱                 ╲____╲`,
		`╱_______________________╲`,
	},
	
	// State 2: Medium inhale - lungs expanding
	{
		`         ╭─────────╮         `,
		`       ╱             ╲       `,
		`     ╱    ╭───────╮    ╲     `,
		`    ╱    ╱    ○    ╲    ╲    `,
		`   ╱    ╱           ╲    ╲   `,
		`  ╱    ╱      ♡      ╲    ╲  `,
		` ╱    ╱               ╲    ╲ `,
		`╱    ╱                 ╲    ╲`,
		`╲____╱                 ╲____╱`,
		`╲_________________________╱`,
	},
	
	// State 3: Full inhale - lungs fully expanded
	{
		`        ╭───────────╮        `,
		`      ╱               ╲      `,
		`    ╱    ╭─────────╮    ╲    `,
		`   ╱    ╱     ○     ╲    ╲   `,
		`  ╱    ╱             ╲    ╲  `,
		` ╱    ╱       ♡       ╲    ╲ `,
		`╱    ╱                 ╲    ╲`,
		`╲    ╱                 ╲    ╱`,
		`╲____╱                 ╲____╱`,
		`╲___________________________╱`,
	},
	
	// State 4: Maximum inhale - peak expansion
	{
		`       ╭─────────────╮       `,
		`     ╱                 ╲     `,
		`   ╱    ╭───────────╮    ╲   `,
		`  ╱    ╱      ○      ╲    ╲  `,
		` ╱    ╱               ╲    ╲ `,
		`╱    ╱        ♡        ╲    ╲`,
		`╲    ╱                 ╲    ╱`,
		`╲   ╱                   ╲   ╱`,
		`╲___╱                   ╲___╱`,
		`╲_____________________________╱`,
	},
}

func NewMeditationExercise(localizer *Localizer) *MeditationExercise {
	return &MeditationExercise{
		Name:        "Deep Breathing Meditation",
		Category:    "Meditation", 
		Description: "Guided deep breathing exercise for relaxation and mindfulness",
		Cycle:       0,
		FrameCount:  0,
		Localizer:   localizer,
		
		// Breathing spring - slow, smooth breathing rhythm
		breathSpring:  harmonica.NewSpring(harmonica.FPS(fps), 0.8, 0.9),
		breathPosition: 0.0,
		breathVelocity: 0.0,
		breathTarget:   0.0,
		
		// Lung expansion spring - follows breathing but with slight delay
		lungSpring:    harmonica.NewSpring(harmonica.FPS(fps), 1.0, 0.8),
		lungPosition:  0.0,
		lungVelocity:  0.0,
		lungTarget:    0.0,
		
		// Heart rate spring - very slow, calming rhythm
		heartSpring:   harmonica.NewSpring(harmonica.FPS(fps), 0.5, 0.95),
		heartPosition: 0.0,
		heartVelocity: 0.0,
		heartTarget:   0.0,
		
		isInhaling:    true,
		breathCycles:  0,
		phase:         "inhale",
		phaseTimer:    0,
		phaseDuration: 120, // 4 seconds at 30fps
	}
}

func (me *MeditationExercise) GetName() string {
	return me.Name
}

func (me *MeditationExercise) GetCategory() string {
	return me.Category
}

func (me *MeditationExercise) GetDescription() string {
	return me.Description
}

func (me *MeditationExercise) renderBreathing() {
	// Calculate the base animation state using breath spring
	normalizedPos := (me.breathPosition + animationRange) / (2 * animationRange)
	if normalizedPos < 0 {
		normalizedPos = 0
	}
	if normalizedPos > 1 {
		normalizedPos = 1
	}
	
	// Base state selection
	baseStateIndex := int(normalizedPos * float64(len(breathingStates)-1))
	if baseStateIndex >= len(breathingStates) {
		baseStateIndex = len(breathingStates) - 1
	}
	
	// Calculate lung expansion effect
	lungOffset := int(me.lungPosition * 0.2)
	
	// Calculate heart rate effect for subtle pulsing
	heartOffset := int(me.heartPosition * 0.1)
	
	// Dynamic padding for breathing effect
	basePadding := 10
	dynamicPadding := basePadding + lungOffset + heartOffset
	if dynamicPadding < 5 {
		dynamicPadding = 5
	}
	if dynamicPadding > 20 {
		dynamicPadding = 20
	}
	
	// Render each line of the breathing animation
	breathLines := breathingStates[baseStateIndex]
	for i, line := range breathLines {
		padding := strings.Repeat(" ", dynamicPadding)
		
		// Add subtle heart beat effect to the heart symbol line
		if strings.Contains(line, "♡") {
			if me.heartPosition > 3.0 {
				line = strings.ReplaceAll(line, "♡", "💖") // Stronger heart beat
			} else if me.heartPosition > 1.0 {
				line = strings.ReplaceAll(line, "♡", "💗") // Medium heart beat
			}
		}
		
		// Add breathing indicators
		if i == 0 && me.phase == "inhale" {
			line += "  ↑ " + me.Localizer.T("inhaling")
		} else if i == 0 && me.phase == "exhale" {
			line += "  ↓ " + me.Localizer.T("exhaling")
		} else if i == 0 && me.phase == "hold" {
			line += "  ⏸ " + me.Localizer.T("holding")
		} else if i == 0 && me.phase == "pause" {
			line += "  ⏹ " + me.Localizer.T("pausing")
		}
		
		fmt.Printf("%s%s\n", padding, line)
	}
}

func (me *MeditationExercise) Render() {
	me.renderBreathing()
}

func (me *MeditationExercise) Update() {
	me.FrameCount++
	me.phaseTimer++
	
	// Update breathing phases (4-7-8 breathing technique)
	switch me.phase {
	case "inhale":
		if me.phaseTimer >= 120 { // 4 seconds
			me.phase = "hold"
			me.phaseTimer = 0
		}
		me.breathTarget = animationRange
		me.lungTarget = animationRange * 0.8
		
	case "hold":
		if me.phaseTimer >= 210 { // 7 seconds
			me.phase = "exhale" 
			me.phaseTimer = 0
		}
		me.breathTarget = animationRange
		me.lungTarget = animationRange * 0.8
		
	case "exhale":
		if me.phaseTimer >= 240 { // 8 seconds
			me.phase = "pause"
			me.phaseTimer = 0
			me.breathCycles++
		}
		me.breathTarget = -animationRange
		me.lungTarget = -animationRange * 0.6
		
	case "pause":
		if me.phaseTimer >= 60 { // 2 seconds
			me.phase = "inhale"
			me.phaseTimer = 0
		}
		me.breathTarget = -animationRange
		me.lungTarget = -animationRange * 0.6
	}
	
	// Update spring physics
	me.breathPosition, me.breathVelocity = me.breathSpring.Update(me.breathPosition, me.breathVelocity, me.breathTarget)
	me.lungPosition, me.lungVelocity = me.lungSpring.Update(me.lungPosition, me.lungVelocity, me.lungTarget)
	
	// Heart rate follows a slow, calming rhythm
	heartTarget := sin(float64(me.FrameCount) * 0.005) * 4.0
	me.heartPosition, me.heartVelocity = me.heartSpring.Update(me.heartPosition, me.heartVelocity, heartTarget)
}

func (me *MeditationExercise) GetInstructions() string {
	switch me.phase {
	case "inhale":
		return me.Localizer.T("breathe_in_instruction")
	case "hold":
		return me.Localizer.T("hold_breath_instruction")
	case "exhale":
		return me.Localizer.T("breathe_out_instruction")
	case "pause":
		return me.Localizer.T("pause_instruction")
	}
	return ""
}

func (me *MeditationExercise) GetTips() []string {
	return []string{
		me.Localizer.T("tip_breathe_478"),
		me.Localizer.T("tip_inhale"),
		me.Localizer.T("tip_hold"),
		me.Localizer.T("tip_exhale"),
		me.Localizer.T("tip_pause"),
		me.Localizer.T("tip_focus"),
		me.Localizer.T("tip_exit"),
	}
}

func (me *MeditationExercise) GetCounter() string {
	return me.Localizer.Tf("breath_counter", me.breathCycles)
}

func (me *MeditationExercise) IsComplete() bool {
	return false // Meditation runs indefinitely until user exits
}

func (me *MeditationExercise) Reset() {
	me.Cycle = 0
	me.FrameCount = 0
	me.breathPosition = 0.0
	me.breathVelocity = 0.0
	me.breathTarget = -animationRange
	me.lungPosition = 0.0
	me.lungVelocity = 0.0
	me.lungTarget = -animationRange * 0.6
	me.heartPosition = 0.0
	me.heartVelocity = 0.0
	me.heartTarget = 0.0
	me.isInhaling = true
	me.breathCycles = 0
	me.phase = "inhale"
	me.phaseTimer = 0
	me.phaseDuration = 120
}

// TerminalGym manages the overall application
type TerminalGym struct {
	currentExercise Exercise
	localizer      *Localizer
}

func NewTerminalGym(localizer *Localizer) *TerminalGym {
	return &TerminalGym{
		localizer: localizer,
	}
}

func (tg *TerminalGym) clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func (tg *TerminalGym) selectExercise() {
	tg.clearScreen()
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("                 " + tg.localizer.T("welcome_title"))
	fmt.Println("                    " + tg.localizer.T("welcome_subtitle"))
	fmt.Println(strings.Repeat("=", 60) + "\n")
	
	fmt.Println(tg.localizer.T("exercise_selection"))
	fmt.Println(tg.localizer.T("exercise_buttock"))
	fmt.Println(tg.localizer.T("exercise_meditation"))
	fmt.Print("\n" + tg.localizer.T("enter_choice"))
	
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			continue
		}
		
		choice := strings.TrimSpace(input)
		choiceNum, err := strconv.Atoi(choice)
		if err != nil || choiceNum < 1 || choiceNum > 2 {
			fmt.Print(tg.localizer.T("invalid_choice") + "\n" + tg.localizer.T("enter_choice"))
			continue
		}
		
		switch choiceNum {
		case 1:
			tg.currentExercise = NewButtockExercise(tg.localizer)
		case 2:
			tg.currentExercise = NewMeditationExercise(tg.localizer)
		}
		break
	}
}

func (tg *TerminalGym) render() {
	tg.clearScreen()
	
	// Title
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("                    " + tg.localizer.T("title"))
	fmt.Println("              " + tg.localizer.T("subtitle"))
	fmt.Println(strings.Repeat("=", 60) + "\n")
	
	// Instructions
	instruction := tg.currentExercise.GetInstructions()
	padding := (60 - len(instruction)) / 2
	if padding < 0 {
		padding = 0
	}
	fmt.Printf("%s%s\n\n", strings.Repeat(" ", padding), instruction)
	
	// Animation area
	fmt.Println("\n" + strings.Repeat(" ", 25) + tg.localizer.T("watch_follow"))
	fmt.Println()
	
	// Render the current exercise
	tg.currentExercise.Render()
	
	// Exercise counter and tips
	fmt.Printf("\n\n%s%s\n", strings.Repeat(" ", 25), tg.currentExercise.GetCounter())
	
	// Tips
	fmt.Println("\n" + strings.Repeat("-", 60))
	fmt.Println(tg.localizer.T("tips_header"))
	for _, tip := range tg.currentExercise.GetTips() {
		fmt.Println(tip)
	}
	fmt.Println(strings.Repeat("-", 60))
}

func (tg *TerminalGym) run() {
	// Set up signal handling for graceful exit
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	// Initialize the current exercise
	tg.currentExercise.Reset()
	
	// Animation loop
	ticker := time.NewTicker(time.Second / fps)
	defer ticker.Stop()
	
	for {
		select {
		case <-c:
			tg.clearScreen()
			if tg.currentExercise.GetCategory() == "Meditation" {
				fmt.Println("\n" + tg.localizer.T("meditation_complete"))
			} else {
				fmt.Println("\n" + tg.localizer.T("workout_complete"))
			}
			fmt.Println(tg.localizer.T("keep_work") + "\n")
			return
		case <-ticker.C:
			tg.currentExercise.Update()
			tg.render()
		}
	}
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
	
	// Exercise selection
	gym.selectExercise()
	
	// Preparation phase
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