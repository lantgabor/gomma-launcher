package downloader

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Downloader struct {
	client *http.Client
	Cache  *Cache
}

func NewDownloader(client *http.Client, cache *Cache) *Downloader {
	if client == nil {
		client = http.DefaultClient
	}
	return &Downloader{client: client, Cache: cache}
}

func (d *Downloader) Download(ctx context.Context, url string, desiredHash string, dst string, onProgress func(done, total int64)) (string, error) {
	dstPath := dst
	_, filename := filepath.Split(dstPath)
	if d.Cache != nil {
		dstPath = d.Cache.Path(filename)
	}

	if fileExists(dstPath) && desiredHash != "" {
		observedHash, err := md5File(dstPath)
		if err != nil {
			return "", fmt.Errorf("failed to calculate md5 hash of %s: %w", dstPath, err)
		}

		if observedHash == desiredHash {
			if d.Cache != nil {
				err = d.Cache.Publish(filename, dst)
				if err != nil {
					return "", fmt.Errorf("failed to publish %s: %w", filename, err)
				}
			}
		}
		return observedHash, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %w", url, err)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download %s: unexpected status %s", url, resp.Status)
	}

	tmp := dstPath + ".part"
	out, err := os.Create(tmp)
	if err != nil {
		return "", fmt.Errorf("failed to create %s: %w", tmp, err)
	}
	defer os.Remove(tmp)
	defer out.Close()

	hasher := md5.New()
	progress := &progressReader{reader: resp.Body, total: resp.ContentLength, onProgress: onProgress}

	if _, err := io.Copy(io.MultiWriter(out, hasher), progress); err != nil {
		return "", fmt.Errorf("failed to write %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, dstPath); err != nil {
		return "", fmt.Errorf("failed to rename %s to %s: %w", tmp, dstPath, err)
	}

	if d.Cache != nil {
		err = d.Cache.Publish(filename, dst)
		if err != nil {
			return "", fmt.Errorf("failed to publish %s: %w", filename, err)
		}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

type progressReader struct {
	reader     io.Reader
	done       int64
	total      int64
	onProgress func(done, total int64)
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	if n > 0 {
		p.done += int64(n)
		if p.onProgress != nil {
			p.onProgress(p.done, p.total)
		}
	}
	return n, err
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Mode().IsRegular()
}

func md5File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
