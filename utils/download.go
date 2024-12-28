package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ProgressReader struct {
	Reader     io.Reader
	Total      int64
	Downloaded int64
	OnProgress func(downloaded, total int64)
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Downloaded += int64(n)
	if pr.OnProgress != nil {
		pr.OnProgress(pr.Downloaded, pr.Total)
	}
	return n, err
}

func DownloadFile(url, destPath string) error {
	// Create request with headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Add headers to mimic browser behavior
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	// Create client with timeout
	client := &http.Client{
		Timeout: 30 * time.Minute,
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error downloading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Create the file
	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer out.Close()

	// Setup progress tracking
	progressReader := &ProgressReader{
		Reader: resp.Body,
		Total:  resp.ContentLength,
		OnProgress: func(downloaded, total int64) {
			if total > 0 {
				progress := float64(downloaded) / float64(total) * 100
				fmt.Printf("\rDownloading... %.1f%%", progress)
			} else {
				fmt.Printf("\rDownloading... %d bytes", downloaded)
			}
		},
	}

	// Copy the content
	_, err = io.Copy(out, progressReader)
	fmt.Println() // New line after progress
	if err != nil {
		return fmt.Errorf("error saving file: %v", err)
	}

	return nil
}