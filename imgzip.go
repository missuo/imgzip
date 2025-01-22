package main

import (
	"fmt"
	_ "image/jpeg"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/disintegration/imaging"
)

const (
	defaultQuality = 80
	targetSize     = 1024 * 1024 // 1MB
)

func main() {
	// Check arguments
	args := os.Args[1:]
	if len(args) == 0 {
		log.Fatal("Usage: imgzip <image path> [quality(1-100)]")
	}

	// Get input file path
	inputPath := args[0]

	// Get compression quality
	quality := defaultQuality
	if len(args) > 1 {
		q, err := strconv.Atoi(args[1])
		if err != nil || q < 1 || q > 100 {
			log.Fatal("Quality must be a number between 1-100")
		}
		quality = q
	}

	// Verify input file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Fatalf("Input file does not exist: %s", inputPath)
	}

	// Generate output file path
	timestamp := time.Now().Format("150405") // HHMMSS timestamp
	ext := filepath.Ext(inputPath)
	filename := strings.TrimSuffix(filepath.Base(inputPath), ext)
	outputPath := fmt.Sprintf("%s_compressed_%s%s", filename, timestamp, ext)

	// Compress the image
	err := compressImage(inputPath, outputPath, quality)
	if err != nil {
		log.Fatal(err)
	}

	// Get compressed file size
	outputInfo, err := os.Stat(outputPath)
	if err != nil {
		log.Fatal(err)
	}
	compressedSize := outputInfo.Size() / 1024 // Convert to KB

	fmt.Printf("Compression complete!\n")
	fmt.Printf("Output file: %s\n", outputPath)
	fmt.Printf("File size: %d KB\n", compressedSize)
}

func compressImage(inputPath, outputPath string, quality int) error {
	// Read source image
	src, err := imaging.Open(inputPath)
	if err != nil {
		return fmt.Errorf("cannot open image: %v", err)
	}

	// Get original image dimensions
	bounds := src.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// Create temporary file for compression testing
	tempFile, err := os.CreateTemp("", "compress_*.jpg")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// First attempt: compress with specified quality
	err = imaging.Save(src, tempFile.Name(), imaging.JPEGQuality(quality))
	if err != nil {
		return fmt.Errorf("failed to save compressed image: %v", err)
	}

	// Check file size
	fileInfo, err := tempFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// If file is still too large, adjust image dimensions
	currentSize := fileInfo.Size()
	if currentSize > targetSize {
		// Calculate scale factor to reach target size
		scale := float64(targetSize) / float64(currentSize)
		newWidth := int(float64(width) * scale)
		newHeight := int(float64(height) * scale)

		// Resize image
		src = imaging.Resize(src, newWidth, newHeight, imaging.Lanczos)
	}

	// Save final image
	err = imaging.Save(src, outputPath, imaging.JPEGQuality(quality))
	if err != nil {
		return fmt.Errorf("failed to save final image: %v", err)
	}

	return nil
}
