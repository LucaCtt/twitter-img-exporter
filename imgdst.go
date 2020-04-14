package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ImgDst represents a destination to which images can be written.
type ImgDst interface {
	Write(imgs []Img) error
}

// DirImgDst is an ImgDst that writes images to a local directory.
type DirImgDst struct {
	Dir string
}

// Write writes the images to a local directory.
func (d *DirImgDst) Write(imgs []Img) error {
	for i := 0; i < len(imgs); i++ {
		img := imgs[i]
		defer img.Content.Close()

		// Create dir
		if err := os.MkdirAll(d.Dir, os.ModePerm); err != nil {
			return fmt.Errorf("Create directory \"%s\" failed: %w", d.Dir, err)
		}

		file := filepath.Join(d.Dir, img.Name)

		// Create file
		out, err := os.Create(file)
		if err != nil {
			return fmt.Errorf("Create file \"%s\" failed: %w", file, err)
		}

		// Write contents of image to file
		if _, err = io.Copy(out, img.Content); err != nil {
			return fmt.Errorf("Write to file \"%s\" failed: %w", file, err)
		}
	}

	return nil
}
