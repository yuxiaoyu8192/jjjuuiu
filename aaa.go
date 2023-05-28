package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

// Downloader is a struct that represents a concurrent file downloader
type Downloader struct {
	url         string  // the url of the file to download
	output      string  // the output filename
	concurrency int     // the number of goroutines to use
	size        int64   // the size of the file in bytes
	ranges      [][2]int64 // the ranges of bytes to download by each goroutine
}

// NewDownloader creates a new Downloader with the given url, output and concurrency
func NewDownloader(url, output string, concurrency int) *Downloader {
	return &Downloader{
		url:         url,
		output:      output,
		concurrency: concurrency,
	}
}

// checkSupportRange checks if the server supports partial requests
func (d *Downloader) checkSupportRange() error {
	resp, err := http.Head(d.url)
	if err != nil {
		return err
	}
	if resp.StatusCode == http.StatusOK && resp.Header.Get("Accept-Ranges") == "bytes" {
		d.size = resp.ContentLength
		return nil
	}
	return fmt.Errorf("server does not support range requests")
}

// calculateRanges calculates the ranges of bytes to download by each goroutine
func (d *Downloader) calculateRanges() {
	chunkSize := d.size / int64(d.concurrency)
	for i := 0; i < d.concurrency; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1
		if i == d.concurrency-1 {
			end = d.size - 1
		}
		d.ranges = append(d.ranges, [2]int64{start, end})
	}
}

// downloadChunk downloads a chunk of the file and writes it to a temporary file
func (d *Downloader) downloadChunk(filename string, r [2]int64) error {
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", r[0], r[1]))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}
	return nil
}

// mergeFiles merges the temporary files into one output file and deletes them
func (d *Downloader) mergeFiles() error {
	outputFile, err := os.Create(d.output)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	for i := 0; i < d.concurrency; i++ {
		tempFile, err := os.Open(strconv.Itoa(i))
		if err != nil {
			return err
		}
		defer tempFile.Close()
		if _, err = io.Copy(outputFile, tempFile); err != nil {
			return err
		}
		os.Remove(strconv.Itoa(i))
	}
	return nil
}

// Download downloads the file concurrently and saves it to the output file
func (d *Downloader) Download() error {
	log.Println("Checking server support for range requests...")
	if err := d.checkSupportRange(); err != nil {
		return err
	}
	log.Printf("The size of the file is %d bytes\n", d.size)
	d.calculateRanges()
	log.Println("The ranges are:", d.ranges)

	var wg sync.WaitGroup

	for i, r := range d.ranges {
		wg.Add(1)
		go func(i int, r [2]int64) {
			defer wg.Done()
			filename := strconv.Itoa(i)
			log.Printf("Downloading %s range %v\n", filename, r)
			err := d.downloadChunk(filename, r)
			if err != nil {
				log.Printf("Error downloading %s: %v\n", filename, err)
			} else {
				log.Printf("Finished downloading %s\n", filename)
			}
			
		}(i, r)
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
			
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
		
		
	
		
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			
			

			

		
	
		
	
		
	
		
	
		
	
		
	
		
	
		
	
		
	
		
	
		
	
		
	
		

		
	

	wg.Wait()

	log.Println("Merging files...")
	err := d.mergeFiles()
	if err != nil {
		return err
	}
	log.Println("Download completed")
	return nil
	
}

func main() {

	urlFlag := flag.String("url", "", "The url of the file to download")
	outputFlag := flag.String("output", "", "The output filename")
	concurrencyFlag := flag.Int("concurrency", 10, "The number of goroutines to use")

	flag.Parse()

	if *urlFlag == "" || *outputFlag == "" {
        log.Fatal("url and output are required")
    }

	downloader := NewDownloader(*urlFlag, *outputFlag, *concurrencyFlag)

	err := downloader.Download()
	if err != nil {
        log.Fatal(err)
    }
}

