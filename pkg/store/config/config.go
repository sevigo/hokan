package config

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
	"github.com/sevigo/hokan/pkg/store/db"
)

var ErrConfigNotFound = errors.New("config not found")

const bucketName = "config"
const targetPrefix = "taget"

type configStore struct {
	db *db.DB
}

func New(database *db.DB) core.ConfigStore {
	err := database.Write(bucketName, "ts", time.Now().String())
	if err != nil {
		// TODO: return error here!
		log.WithError(err).Fatalf("store.config.New(): can't write to the bucket %q", bucketName)
		return nil
	}
	return &configStore{database}
}

func (s *configStore) Save(ctx context.Context, conf *core.TargetConfig) error {
	key := fmt.Sprintf("%s:%s", targetPrefix, conf.Name)
	var value bytes.Buffer
	if err := json.NewEncoder(&value).Encode(conf); err != nil {
		return err
	}
	return s.db.Write(bucketName, key, value.String())
}

func (s *configStore) Find(ctx context.Context, name string) (*core.TargetConfig, error) {
	config := &core.TargetConfig{}
	key := fmt.Sprintf("%s:%s", targetPrefix, name)

	value, err := s.db.Read(bucketName, key)
	if err != nil {
		return nil, err
	}
	if value == nil {
		return nil, ErrConfigNotFound
	}

	errJ := json.NewDecoder(bytes.NewReader(value)).Decode(config)
	return config, errJ
}
