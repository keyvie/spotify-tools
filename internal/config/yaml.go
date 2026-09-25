package config

import (
	"os"
	"path/filepath"
	"errors"
	"context"
	"time"

	"gopkg.in/yaml.v3"
	"github.com/gofrs/flock"
)

var fileLock *flock.Flock

func createDirIfNotExists() error {
	dir := filepath.Dir(Path)
	if dir == "" || dir == "." {
		return nil
	}
	return os.MkdirAll(dir, 0700)
}

func Lock() error {
	if fileLock == nil {
		if err := createDirIfNotExists(); err != nil {
			return err
		}
		fileLock = flock.New(Path + ".lock")
	} 
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	locked, err := fileLock.TryLockContext(ctx, 100*time.Millisecond)
	if err != nil {
		return err
	}

	if !locked {
		return errors.New("config file is locked")
	}
	return nil
}

func Unlock() error {
	if fileLock == nil {
		return nil
	}
	return fileLock.Unlock()
}

func Locked() bool {
	return fileLock != nil && fileLock.Locked()
}

func Load() (*Config, error) {
	if err := createDirIfNotExists(); err != nil {
		return nil, err
	}
	
	var fileData, err = os.ReadFile(Path)
	if errors.Is(err, os.ErrNotExist) {
		defaultConfig := &Config{}
		err = Save(defaultConfig)
		if err != nil {
			return nil, err
		}
		fileData, err = os.ReadFile(Path)
	}
	if err != nil {
		return nil, err
	}

	var config Config

	err = yaml.Unmarshal(fileData, &config) 
	if err != nil {
		return nil, err
	}

	edited := false
	if config.Secrets.RedirectURI == "" {
		config.Secrets.RedirectURI = "localhost:36793/callback"
		edited = true
	}
	if edited {
		if err := Save(&config); err != nil {
			return nil, err
		}
	}

	return &config, nil
}


func Save(config *Config) error {
	fileData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	if err := createDirIfNotExists(); err != nil {
		return err
	}

	return os.WriteFile(Path, fileData, 0600)
}

