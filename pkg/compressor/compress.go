package compressor

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func Zip(sourceDir, zipFileName string) error {
	// Read the files from the source directory
	files, err := ioutil.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	zipFile, err := os.Create(filepath.Join(sourceDir, zipFileName))
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if !strings.Contains(file.Name(), ".png") {
			continue
		}
		filePath := filepath.Join(sourceDir, file.Name())
		if err := addFileToZip(zipWriter, filePath); err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filePath string) error {
	// Open the file to be added to the zip archive
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a zip file header
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = filepath.Base(filePath)
	header.Method = zip.Deflate

	// Create a writer for the file in the zip archive
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	// Copy the file content to the zip archive
	_, err = io.Copy(writer, file)
	return err
}
