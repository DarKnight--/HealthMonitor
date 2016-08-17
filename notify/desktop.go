package notify

import (
	"os/exec"
	"runtime"

	"github.com/owtf/health_monitor/utils"
)

type (
	desktopNotifier interface {
		push(title string, body string, iconPath string) *exec.Cmd
		pushCritical(title string, body string, iconPath string) *exec.Cmd
	}

	// DesktopAlert hold configuration for desktop notification
	DesktopAlert struct {
		notifier    desktopNotifier
		defaultIcon string
	}

	// MessageImportance is constant type for the desktop notification messages
	MessageImportance int
)

const (
	// Normal : when notification will pop up and go after some time
	Normal MessageImportance = 0
	// Critical : when notification will not terminate until user closes it
	Critical MessageImportance = 1
)

// Push is used to send notification to the desktop.
func (n DesktopAlert) Push(title string, body string, iconPath string, urgent MessageImportance) error {
	icon := n.defaultIcon

	if iconPath != "" {
		icon = iconPath
	}

	if urgent == Critical {
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

// NewDesktopAlert constructs the struct according to the OS to send desktop notification
func NewDesktopAlert(defaultIcon string) *DesktopAlert {
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
