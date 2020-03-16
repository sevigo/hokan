package local

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	filestore "github.com/sevigo/hokan/pkg/store/file"
	"github.com/sevigo/hokan/pkg/target/utils"
)

const TargetName = "local"

// Storage local
type localStorage struct {
	configStore       core.ConfigStore
	fileStore         core.FileStore
	bucketName        string
	targetStoragePath string
}

const bufferSize = 1024 * 1024

func DefaultConfig() *core.TargetConfig {
	return &core.TargetConfig{
		Active:      false,
		Name:        "local",
		Description: "store the files on the local disk",
		Settings: map[string]string{
			"LOCAL_STORAGE_PATH": "",
			"LOCAL_BUCKET_NAME":  "",
		},
	}
}

func New(ctx context.Context, fs core.FileStore, conf core.TargetConfig) (core.TargetStorage, error) {
	if !conf.Active {
		return nil, core.ErrTargetNotActive
	}
	// TODO: validate config
	path := filepath.Clean(conf.Settings["LOCAL_STORAGE_PATH"])
	bucket := conf.Settings["LOCAL_BUCKET_NAME"]
	return &localStorage{
		bucketName:        bucket,
		targetStoragePath: path,
		fileStore:         fs,
	}, nil
}

func (s *localStorage) Save(ctx context.Context, file *core.File) error {
	logger := log.WithFields(log.Fields{
		"target": TargetName,
		"file":   file.Path,
	})

	// TODO: this is all the same, move me
	storedFile, err := s.fileStore.Find(ctx, TargetName, file.Path)
	if errors.Is(err, filestore.ErrFileEntryNotFound) || utils.FileHasChanged(file, storedFile) {
		logger.Debug("saving file")
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
		errSave := s.save(file.Path, to)
		if err != nil {
			return errSave
		}
		return s.fileStore.Save(ctx, TargetName, file)
	}
	logger.Info("the file has not changedб шптщку")
	return nil
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
