/*
 * @Author: Vincent Yang
 * @Date: 2025-01-22 18:14:46
 * @LastEditors: Vincent Yang
 * @LastEditTime: 2025-01-22 18:32:38
 * @FilePath: /imgzip/imgzip.go
 * @Telegram: https://t.me/missuo
 * @GitHub: https://github.com/missuo
 *
 * Copyright © 2025 by Vincent, All Rights Reserved.
 */

package main

import (
	"fmt"
	_ "image/jpeg"
	"log"
	"math"
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

	// Get original file size
	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Fatal(err)
	}
	originalSize := inputInfo.Size() / 1024 // Convert to KB

	// Compress the image
	err = compressImage(inputPath, outputPath, quality)
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
	fmt.Printf("Original size: %d KB\n", originalSize)
	fmt.Printf("Compressed size: %d KB\n", compressedSize)
	fmt.Printf("Compression ratio: %.2f%%\n", float64(compressedSize)/float64(originalSize)*100)
	fmt.Printf("Output file: %s\n", outputPath)
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

	// Quality-based initial scale calculation
	// Higher quality means less aggressive scaling
	qualityScale := float64(quality) / 100.0
	initialScale := math.Max(0.1, qualityScale) // Ensure minimum scale is 0.1

	// Calculate initial dimensions based on quality
	newWidth := int(float64(width) * initialScale)
	newHeight := int(float64(height) * initialScale)

	// First resize based on quality
	resized := imaging.Resize(src, newWidth, newHeight, imaging.Lanczos)

	// Create temporary file for compression testing
	tempFile, err := os.CreateTemp("", "compress_*.jpg")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Try to save with current dimensions and quality
	err = imaging.Save(resized, tempFile.Name(), imaging.JPEGQuality(quality))
	if err != nil {
		return fmt.Errorf("failed to save compressed image: %v", err)
	}

	// Check if size meets target
	fileInfo, err := tempFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	currentSize := fileInfo.Size()
	if currentSize > targetSize {
		// If still too large, apply additional scaling while maintaining aspect ratio
		scale := math.Sqrt(float64(targetSize) / float64(currentSize))
		finalWidth := int(float64(newWidth) * scale)
		finalHeight := int(float64(newHeight) * scale)
		resized = imaging.Resize(resized, finalWidth, finalHeight, imaging.Lanczos)
	}

	// Save final image
	err = imaging.Save(resized, outputPath, imaging.JPEGQuality(quality))
	if err != nil {
		return fmt.Errorf("failed to save final image: %v", err)
	}

	return nil
}
