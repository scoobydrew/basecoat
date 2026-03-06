package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Storage is the interface for file storage operations.
// Swap the implementation to move from local FS to S3 or similar.
type Storage interface {
	// Put stores the content from r at the given key. Returns the storage key.
	Put(key string, r io.Reader, contentType string) error
	// Delete removes the file at the given key.
	Delete(key string) error
	// URL returns the publicly accessible URL for a given key.
	URL(key string) string
}

// LocalStorage stores files on the local filesystem.
type LocalStorage struct {
	root    string // absolute path to storage directory
	baseURL string // base URL prefix for serving files
}

func NewLocalStorage(root, baseURL string) (*LocalStorage, error) {
	if err := os.MkdirAll(root, 0755); err != nil {
		return nil, fmt.Errorf("create storage root: %w", err)
	}
	return &LocalStorage{root: root, baseURL: baseURL}, nil
}

func (s *LocalStorage) Put(key string, r io.Reader, _ string) error {
	dest := filepath.Join(s.root, key)
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("create dirs for %s: %w", key, err)
	}
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create file %s: %w", key, err)
	}
	defer f.Close()
	if _, err := io.Copy(f, r); err != nil {
		return fmt.Errorf("write file %s: %w", key, err)
	}
	return nil
}

func (s *LocalStorage) Delete(key string) error {
	if err := os.Remove(filepath.Join(s.root, key)); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete file %s: %w", key, err)
	}
	return nil
}

func (s *LocalStorage) URL(key string) string {
	return s.baseURL + "/" + key
}
