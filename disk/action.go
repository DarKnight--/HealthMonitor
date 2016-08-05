package disk

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"time"

	"health_monitor/setup"
	"health_monitor/utils"
)

const (
	//DebianAPTPath is the path where downloaded apt packages are found
	DebianAPTPath = "/var/cache/apt/archives"
)

var (
	kali = []byte{75, 97, 108, 105, 10}
)

func basicCleanup() {
	//TODO pause owtf
	utils.ModuleLogs(logFile, "Performing basic cleanup")
	utils.ModuleLogs(logFile, "Performing apt cache clean up")
	removeAptCache()
	utils.ModuleLogs(logFile, "Compressing owtf proxy-cache: /tmp/owtf/proxy-cache")
	compressFolder("/tmp/owtf/proxy-cache", "/tmp/owtf/proxy-cache"+time.Now().Format(time.Stamp)+".tar.gz")
	os.RemoveAll(utils.GetPath(".w3af/tmp/"))
}

func removeAptCache() {
	cmd := exec.Command("sudo", "apt-get", "clean")
	err := cmd.Run()
	if err != nil {
		utils.ModuleError(logFile, "Unable to clean apt-get", err.Error())
	}
	if reflect.DeepEqual(setup.OSVarient, kali) {
		err := os.RemoveAll(DebianAPTPath)
		if err != nil {
			utils.ModuleError(logFile, "Unable to free apt cache", err.Error())
		}
		utils.ModuleLogs(logFile, "Deleted apt cache successfully.")
	}
}

func emptyTrash() {
	os.RemoveAll(utils.GetPath(".local/share/Trash/"))
}

func dirSizeMB(path string) int {
	sizes := make(chan int64)
	readSize := func(path string, file os.FileInfo, err error) error {
		if err != nil || file == nil {
			return nil // Ignore errors
		}
		if !file.IsDir() {
			sizes <- file.Size()
		}
		return nil
	}

	go func() {
		filepath.Walk(path, readSize)
		close(sizes)
	}()

	size := int64(0)
	for s := range sizes {
		size += s
	}

	sizeMB := int(float64(size) / 1024.0 / 1024.0)

	return sizeMB
}

func compressFolder(basePath string, outFName string) error {
	outFile, err := os.Create(outFName)
	if err != nil {
		utils.ModuleError(logFile, "Unable to open file for compressing", err.Error())
		return err
	}
	defer outFile.Close()

	gzipWriter := gzip.NewWriter(outFile)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	walkFunction := func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsDir() {
			return nil
		}

		newPath := path[len(basePath)+1:]
		if len(newPath) == 0 {
			return nil
		}
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		defer fr.Close()

		if h, err := tar.FileInfoHeader(info, newPath); err != nil {
			log.Fatalln(err)
		} else {
			h.Name = newPath
			if err = tarWriter.WriteHeader(h); err != nil {
				utils.ModuleError(logFile, "Unable to add file headers", err.Error())
			}
		}
		if _, err := io.Copy(tarWriter, fr); err != nil {
			utils.ModuleError(logFile, "Unable to write to tar file", err.Error())
		}
		return nil
	}
	if err = filepath.Walk(basePath, walkFunction); err != nil {
		return err
	}
	return nil
}
