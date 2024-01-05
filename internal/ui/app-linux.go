//go:build linux

package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	xtheme "fyne.io/x/fyne/theme"
	"github.com/godbus/dbus/v5"
)

// IsAutoStarted checks if the app is set to auto start at login.
func (app *App) IsAutoStarted() bool {
	autoStartFile := filepath.Join(os.ExpandEnv("$HOME/.config/autostart"), app.UniqueID()+".desktop")
	if _, err := os.Stat(autoStartFile); err == nil {
		return true
	}
	return false
}

// InstallAutoStart installs the app to auto start at login.
func (app *App) InstallAutoStart() error {
	desktopFile := findDesktopFile(app)
	if desktopFile == "" {
		return fmt.Errorf("The desktop file for %s was not found", app.UniqueID())
	}

	autoStartFile := filepath.Join(os.ExpandEnv("$HOME/.config/autostart"), app.UniqueID()+".desktop")

	autoStartDir := os.ExpandEnv("$HOME/.config/autostart")
	if _, err := os.Stat(autoStartDir); os.IsNotExist(err) {
		os.MkdirAll(autoStartDir, 0755)
	}
	// make a link to the desktop file
	return os.Symlink(desktopFile, autoStartFile)
}

// UninstallAutoStart uninstalls the app from auto start at login.
func (app *App) UninstallAutoStart() error {
	autoStartFile := filepath.Join(os.ExpandEnv("$HOME/.config/autostart"), app.UniqueID()+".desktop")
	return os.Remove(autoStartFile)
}

// DesktopIntegration sets the desktop integration for the app.
func (app *App) DesktopIntegration() {
	// use Adwaita theme on Linux
	app.Settings().SetTheme(xtheme.AdwaitaTheme())
	if scale := getGnomeFontSize(); scale > 0 {
		os.Setenv("FYNE_SCALE", fmt.Sprintf("%f", scale))
	}

}

func getGnomeFontSize() float32 {
	// ne noeed to scale if the env is set
	if os.Getenv("FYNE_SCALE") != "" {
		return -1
	}

	// only for Gnome now...
	if os.Getenv("XDG_CURRENT_DESKTOP") != "" {
		if !strings.Contains(
			strings.ToUpper(os.Getenv("XDG_CURRENT_DESKTOP")),
			"GNOME",
		) {
			return -1
		}
	}

	var dbusConn *dbus.Conn
	var err error
	if dbusConn, err = dbus.SessionBus(); err != nil {
		return -1
	}

	// big thanks to D-Spy developers => https://gitlab.gnome.org/GNOME/d-spy
	busname := "org.freedesktop.portal.Desktop"
	path := "/org/freedesktop/portal/desktop"
	interfaceName := "org.freedesktop.portal.Settings"
	method := interfaceName + ".ReadOne"
	namespace := "org.gnome.desktop.interface"
	property := "text-scaling-factor"

	obj := dbusConn.Object(busname, dbus.ObjectPath(path))
	resp := obj.Call(method, dbus.FlagNoAutoStart, namespace, property)

	if resp.Err != nil {
		return -1
	}

	var scale float32
	if err = resp.Store(&scale); err != nil {
		return -1
	}

	return scale * 0.92 // 0.92 is a magic number to make it look good
}

// findDesktopFile finds the desktop file for the app.
func findDesktopFile(app *App) string {
	userfile := filepath.Join(os.ExpandEnv("$HOME/.local/share/applications"), app.UniqueID()+".desktop")
	if _, err := os.Stat(userfile); err == nil {
		return userfile
	}
	systemfile := filepath.Join("/usr/share/applications", app.UniqueID()+".desktop")
	if _, err := os.Stat(systemfile); err == nil {
		return systemfile
	}
	return ""
}
