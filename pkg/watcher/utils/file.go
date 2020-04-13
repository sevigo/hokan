package utils

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	"github.com/sevigo/hokan/pkg/core"
)

func FileChecksumInfo(path string) (string, *core.FileInfo, error) {
	f, erro := os.Open(path)
	if erro != nil {
		return "", nil, erro
	}
	defer f.Close()
	info, errs := f.Stat()
	if errs != nil {
		return "", nil, errs
	}

	if info.IsDir() {
		return "", nil, fmt.Errorf("not a file")
	}

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", nil, err
	}

	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, &core.FileInfo{info}, nil
}
