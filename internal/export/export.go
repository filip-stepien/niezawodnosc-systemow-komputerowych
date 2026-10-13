// Package export writes results to files.
package export

import (
	"encoding/csv"
	"errors"
	"image"
	"image/png"
	"io"
	"os"
)

const (
	// csvSeparator is the field separator expected by spreadsheets that use a
	// decimal comma.
	csvSeparator = ';'
	// utf8BOM makes spreadsheets read the file as UTF-8.
	utf8BOM = "\xEF\xBB\xBF"
)

// WriteCSV writes records as a semicolon-separated UTF-8 file.
func WriteCSV(w io.Writer, records [][]string) error {
	if _, err := io.WriteString(w, utf8BOM); err != nil {
		return err
	}
	writer := csv.NewWriter(w)
	writer.Comma = csvSeparator
	writer.UseCRLF = true
	return writer.WriteAll(records)
}

// SavePNG writes img to a PNG file at path.
func SavePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	return errors.Join(png.Encode(file, img), file.Close())
}
