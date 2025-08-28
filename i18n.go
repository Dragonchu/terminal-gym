package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Localizer handles internationalization
type Localizer struct {
	translations map[string]string
	language     string
}

// NewLocalizer creates a new localizer with the specified language
func NewLocalizer(lang string) (*Localizer, error) {
	l := &Localizer{
		language: lang,
		translations: make(map[string]string),
	}
	
	if err := l.loadTranslations(); err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}
	
	return l, nil
}

// loadTranslations loads the translation file for the current language
func (l *Localizer) loadTranslations() error {
	filename := filepath.Join("locales", l.language+".json")
	
	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Fallback to English if the language file doesn't exist
		l.language = "en"
		filename = filepath.Join("locales", "en.json")
	}
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read translation file %s: %w", filename, err)
	}
	
	if err := json.Unmarshal(data, &l.translations); err != nil {
		return fmt.Errorf("failed to parse translation file %s: %w", filename, err)
	}
	
	return nil
}

// T translates a key to the current language
func (l *Localizer) T(key string) string {
	if translation, exists := l.translations[key]; exists {
		return translation
	}
	// Return the key itself if translation is not found
	return key
}

// Tf translates a key with formatting
func (l *Localizer) Tf(key string, args ...interface{}) string {
	template := l.T(key)
	return fmt.Sprintf(template, args...)
}

// GetLanguage returns the current language
func (l *Localizer) GetLanguage() string {
	return l.language
}