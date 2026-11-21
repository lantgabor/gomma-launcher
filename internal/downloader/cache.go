package downloader

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Cache struct {
	dir string
}

func NewCache(dir string) (*Cache, error) {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create cache dir at %s: %w", dir, err)
	}
	return &Cache{
		dir: dir,
	}, nil
}

func (c Cache) Path(filename string) string {
	return filepath.Join(c.dir, filename)
}

func (c *Cache) Exists(filename string) bool {
	return fileExists(c.Path(filename))
}

func (c Cache) Publish(filename string, dstPath string) error {
	if !c.Exists(filename) {
		return fmt.Errorf("cached file doesn't exists: %s", filename)
	}
	srcPath := c.Path(filename)
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open cached file %s: %w", srcPath, err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", dstPath, err)
	}

	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(dstPath)
		return fmt.Errorf("failed to copy %s to %s: %w", srcPath, dstPath, err)
	}
	if err := dst.Close(); err != nil {
		os.Remove(dstPath)
		return fmt.Errorf("failed to close %s: %w", dstPath, err)
	}
	return nil
}
