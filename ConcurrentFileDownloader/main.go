package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// One File
func DownloadFile(url, destDir string) error {
	filename := filepath.Base(url)
	filePath := filepath.Join(destDir, filename)
	out, error := os.Create(filePath)
	if error != nil {
		return error
	}
	defer out.Close()
	fmt.Println("Downloading... ", url)
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		_ = os.Remove(filePath)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		_ = os.Remove(filePath)
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		_ = os.Remove(filePath)
		return err
	}

	fmt.Printf("Downloaded %s in %s\n", filename, time.Since(start))
	return nil
}

func SequentialDownloader(urls []string, destDir string) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	start := time.Now()
	for _, url := range urls {
		if err := DownloadFile(url, destDir); err != nil {
			fmt.Println("Error downloading file:", err)
			continue
		}
	}
	fmt.Printf("Downloaded %d files in %s\n", len(urls), time.Since(start))
	return nil
}

type Result struct {
	Url      string
	FileName string
	Size     int64
	Duration time.Duration
	Error    error
}

func ConcurrentDownloader(urls []string, destDir string, maxConcurrent int) error {
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	results := make(chan Result)
	//start := time.Now()
	var wg sync.WaitGroup
	limiter := make(chan struct{}, maxConcurrent)
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			limiter <- struct{}{}
			defer func() { <-limiter }()
			start := time.Now()
			filename := filepath.Base(url)
			filepath := filepath.Join(destDir, filename)
			out, err := os.Create(filepath)
			if err != nil {
				results <- Result{Url: url, Error: err}
				return
			}
			defer out.Close()
			resp, err := http.Get(url)
			if err != nil {
				results <- Result{Url: url, Error: err}
				return
			}

			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				results <- Result{Url: url, Error: fmt.Errorf("failed to download file: %s", url)}
				return
			}
			size, err := io.Copy(out, resp.Body)
			if err != nil {
				results <- Result{Url: url, Error: err}
				return
			}
			results <- Result{Url: url, FileName: filename, Size: size, Duration: time.Since(start), Error: nil}
		}(url)

	}

	//part 1
	go func() {
		wg.Wait()
		close(results)
	}()

	//part 2
	var totalSize int64
	var errors []error
	start := time.Now()

	for result := range results {
		if result.Error != nil {
			fmt.Println("Error downloading file:", result.Error.Error())
			errors = append(errors, result.Error)
		} else {
			totalSize += result.Size
			fmt.Printf("Downloaded %s, size: %d bytes in time: %s\n", result.FileName, result.Size, result.Duration)
		}
	}

	startedSince := time.Since(start)
	fmt.Printf("Total size downloaded: %d bytes in time: %s\n", totalSize, startedSince)

	if len(errors) > 0 {
		return fmt.Errorf("Errors occurred during file downloads: %+v", errors)
	}
	return nil
}

func main() {

	urls := []string{
		"https://images.pexels.com/photos/38293401/pexels-photo-38293401.jpeg",
		"https://images.pexels.com/photos/38293401/pexels-photo-38293401.jpeg",
		"https://images.pexels.com/photos/38293401/pexels-photo-38293401.jpeg",
		"https://images.pexels.com/photos/38293401/pexels-photo-38293401.jpeg",
		"https://go.dev/images/go-logo-white.svg",
		"https://go.dev/images/go-logo-white.svg",
		"https://go.dev/images/go-logo-white.svg",
		"https://go.dev/images/go-logo-white.svg",
	}
	err := ConcurrentDownloader(urls, "./imgs", 8)
	if err != nil {
		fmt.Println("Error downloading file:", err)
	}
	fmt.Println("File Downloaded")
}
