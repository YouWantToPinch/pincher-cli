package cli

import "github.com/charmbracelet/lipgloss"

// Styles holds the values for the pincher-cli lipgloss styles.
type styles struct {
	Orange lipgloss.Style
	Green  lipgloss.Style
	White  lipgloss.Style
}

func (s *styles) Init() {
	s.Orange = lipgloss.NewStyle().Foreground(lipgloss.Color("#F79269"))
	s.White = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF"))
	s.Green = lipgloss.NewStyle().Foreground(lipgloss.Color("#4FD6BE"))
}
