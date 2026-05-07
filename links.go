package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"github.com/fsnotify/fsnotify"
)

const linksFilename = "_config/links.json"

func refreshLinks(links *map[string]string) error {
	jsonFile, err := os.Open(linksFilename)
	if err != nil {
		return fmt.Errorf("file open: %w", err)
	}

	defer func() {
		_ = jsonFile.Close()
	}()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	err = json.Unmarshal(bytes, links)
	if err != nil {
		return fmt.Errorf("json: %w", err)
	}

	return nil
}

func watchLinks(links *map[string]string, lock *sync.RWMutex) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("watcher: %w", err)
	}
	defer func() {
		if err := watcher.Close(); err != nil {
			slog.Error("watcher close failed", "err", err)
		}
	}()

	if err = watcher.Add(linksFilename); err != nil {
		return fmt.Errorf("watch links file: %w", err)
	}

	defer slog.Error("watched ended")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return errors.New("events not ok")
			}

			if event.Has(fsnotify.Write) {
				lock.Lock()
				err = refreshLinks(links)
				lock.Unlock()

				if err != nil {
					return err
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return errors.New("errors not ok")
			}

			return err
		}
	}
}
