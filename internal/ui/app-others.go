//go:build !linux

package ui

import "fmt"

// IsAutoStarted checks if the app is set to auto start at login.
func (app *App) IsAutoStarted() bool {
	return false
}

// InstallAutoStart installs the app to auto start at login.
func (app *App) InstallAutoStart() error {
	return fmt.Errorf("not implemented")
}

// UninstallAutoStart uninstalls the app from auto start at login.
func (app *App) UninstallAutoStart() error {
	return fmt.Errorf("not implemented")
}

// DesktopIntegration sets the desktop integration for the app.
func (app *App) DesktopIntegration() {
	// no-op, windows and macos don't need this
}
