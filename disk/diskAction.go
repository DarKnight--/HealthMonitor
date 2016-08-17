package disk

import (
	"os"

	"github.com/owtf/health_monitor/utils"
)

const (
	//DebianAPTPath is the path where downloaded apt packages are found
	DebianAPTPath = "/var/cache/apt/archives"
	// KALI is a constant with value "Kali"
	KALI = "Kali"
	// UBUNTU is a constant with value "Ubuntu"
	UBUNTU = "Ubuntu"
)

type (
	cleaner interface {
		RemovePackageManagerCache() error
		EmptyTrash() error
	}
	// BasicCleaner is a struct containing cleaner interface
	BasicCleaner struct {
		os      string
		cleaner cleaner
	}
	kali struct {
	}
	ubuntu struct {
	}
)

func (k kali) RemovePackageManagerCache() error {
	return removeDebianPackageCache()
}

func (k kali) EmptyTrash() error {
	return emptyDebianTrash()
}

func (u ubuntu) RemovePackageManagerCache() error {
	return removeDebianPackageCache()
}

func (u ubuntu) EmptyTrash() error {
	return emptyDebianTrash()
}

// NewBasicCleaner return the struct according to the OS to clean disk space
func NewBasicCleaner(OS string) *BasicCleaner {
	var basicCleaner cleaner

	switch OS {
	case KALI:
		basicCleaner = kali{}
	case UBUNTU:
		basicCleaner = ubuntu{}
	default:
		basicCleaner = nil
	}
	return &BasicCleaner{cleaner: basicCleaner, os: OS}
}

func removeDebianPackageCache() error {
	err := os.RemoveAll(DebianAPTPath)
	if err != nil {
		return err
	}
	return nil
}

func emptyDebianTrash() error {
	return os.RemoveAll(utils.GetPath(".local/share/Trash/"))
}
