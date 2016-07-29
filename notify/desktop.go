package notify

import (
	"os/exec"
	"runtime"
)

const (
	NORMAL   = 0
	CRITICAL = 1
)

type desktopNotifier interface {
	push(title string, body string, iconPath string) *exec.Cmd
	pushCritical(title string, body string, iconPath string) *exec.Cmd
}

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

type osxDesktopAlert struct {
	appName string
}

func (o osxDesktopAlert) push(title string, body string, iconPath string) *exec.Cmd {
	return exec.Command("growlnotify", "-n", o.appName, "--image", iconPath, "-m", title)
}

func (o osxDesktopAlert) pushCritical(title string, body string, iconPath string) *exec.Cmd {
	return exec.Command("notify-send", "-i", iconPath, title, body, "--sticky", "-p", "2")
}

func DesktopAlertBuilder(appName string, defaultIcon string) *DesktopAlert {

	var notifier desktopNotifier

	switch runtime.GOOS {
	case "darwin":
		notifier = osxDesktopAlert{appName: appName}

	case "linux":
		notifier = linuxDesktopAlert{}

	// Windows support not added yet
	case "windows":
		notifier = nil

	default:
		notifier = nil

	}

	return &DesktopAlert{notifier: notifier, defaultIcon: defaultIcon}
}

// see growlnotify for darwin
func CheckDesktopAlertSupport() bool {
	switch runtime.GOOS {
	case "darwin", "linux":
		return checkNotifySend()
	// Add support for windows in future
	case "windows":
		return false
	default:
		return false
	}
}

func checkNotifySend() bool {
	command := exec.Command("notify-send", "--help")
	if command.Run() != nil {
		return false
	}
	return true
}
