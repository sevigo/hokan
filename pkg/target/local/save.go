package local

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func (s *localStorage) save(fromPath, toPath string) error {
	fmt.Printf("storage.local.store(): Copy file from [%s] to [%s]\n", fromPath, toPath)

	log.WithFields(log.Fields{
		"op":        "Save",
		"target":    TargetName,
		"from-file": fromPath,
		"to-file":   toPath,
	}).Info("copy file")

	from, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("cannot open file  [%s]: %v", fromPath, err)
	}
	defer from.Close()
	readBuffer := bufio.NewReader(from)

	fileStoragePath := filepath.Dir(toPath)
	if errDir := os.MkdirAll(fileStoragePath, 0744); errDir != nil {
		return fmt.Errorf("mkdirAll for path: [%s] err: %v", fileStoragePath, errDir)
	}

	to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file [%s] to write: %v", toPath, err)
	}
	defer to.Close()
	writeBuffer := bufio.NewWriter(to)

	totalWritten := 0
	buf := make([]byte, bufferSize)
	for {
		// read a chunk
		n, errRead := readBuffer.Read(buf)
		if errRead != nil && errRead != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		written, errWrite := writeBuffer.Write(buf[:n])
		if errWrite != nil {
			return errWrite
		}
		totalWritten += written
	}
	if err = writeBuffer.Flush(); err != nil {
		return fmt.Errorf("cannot write buffer: %v", err)
	}
	return nil
}
