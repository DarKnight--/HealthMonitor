package notify

import (
	"os/exec"
	"runtime"

	"health_monitor/utils"
)

const (
	// NORMAL : when notification will pop up and go after some time
	NORMAL = 0
	// CRITICAL : when notification will not terminate until user closes it
	CRITICAL = 1
)

type desktopNotifier interface {
	push(title string, body string, iconPath string) *exec.Cmd
	pushCritical(title string, body string, iconPath string) *exec.Cmd
}

// DesktopAlert hold configuration for desktop notification
type DesktopAlert struct {
	notifier    desktopNotifier
	defaultIcon string
}

// Push is used to send notification to the desktop.
func (n DesktopAlert) Push(title string, body string, iconPath string, urgent int) error {
	icon := n.defaultIcon

	if iconPath != "" {
		icon = iconPath
	}

	if urgent == CRITICAL {
		return n.notifier.pushCritical(title, body, icon).Run()
	}

	return n.notifier.push(title, body, icon).Run()

}

type linuxDesktopAlert struct{}

func (l linuxDesktopAlert) push(title string, body string, iconPath string) *exec.Cmd {
	return exec.Command("notify-send", "-i", iconPath, title, body)
}

func (l linuxDesktopAlert) pushCritical(title string, body string, iconPath string) *exec.Cmd {
	return exec.Command("notify-send", "-i", iconPath, title, body, "-u", "critical")
}

// DesktopAlertBuilder return the struct according to the OS to send desktop notification
func DesktopAlertBuilder(appName string, defaultIcon string) *DesktopAlert {
	var notifier desktopNotifier

	switch runtime.GOOS {
	case "linux":
		notifier = linuxDesktopAlert{}
	default:
		notifier = nil

	}

	return &DesktopAlert{notifier: notifier, defaultIcon: defaultIcon}
}

// CheckDesktopAlertSupport returns true if packages required to send
// desktop notification are installed.
func CheckDesktopAlertSupport() bool {
	switch runtime.GOOS {
	case "linux":
		return utils.CheckInstalledPackage("notify-send")
	default:
		return false
	}
}
