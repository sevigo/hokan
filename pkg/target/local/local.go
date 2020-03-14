package local

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

const TargetName = "local"

// Storage local
type localStorage struct {
	bucketName        string
	targetStoragePath string
}

const bufferSize = 1024 * 1024

func New(ctx context.Context, fs core.FileStore) (core.TargetStorage, error) {
	path := `C:\backup`
	bucket := "osaka" //TODO: must be from the config
	return &localStorage{
		bucketName:        bucket,
		targetStoragePath: path,
	}, nil
}

func (s *localStorage) Save(ctx context.Context, file *core.File) error {
	volume := filepath.VolumeName(file.Path)
	base := volume + string(os.PathSeparator)
	relFilePath, err := filepath.Rel(base, file.Path)
	if err != nil {
		return err
	}
	// on windows volume will be 'C:', so we just remove :
	// on all other systems it will be empty
	if volume != "" {
		volume = strings.TrimSuffix(volume, ":")
	}
	to := filepath.Join(s.targetStoragePath, s.bucketName, volume, relFilePath)
	return s.save(file.Path, to)
}

func (s *localStorage) List(context.Context) ([]*core.File, error) {
	log.WithField("target", TargetName).Print("List")
	return nil, nil
}

func (s *localStorage) Find(ctx context.Context, q string) (*core.File, error) {
	log.WithField("target", TargetName).Print("Find")
	return nil, nil
}

func (s *localStorage) Delete(ctx context.Context, file *core.File) error {
	log.WithField("target", TargetName).Print("Delete")
	return nil
}

func (s *localStorage) Ping(ctx context.Context) error {
	return nil
}

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
