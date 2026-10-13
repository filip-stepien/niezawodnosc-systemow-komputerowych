package export

import (
	"bytes"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCSVUsesSemicolonsAndBOM(t *testing.T) {
	var buffer bytes.Buffer
	records := [][]string{{"t_i", "F*(t_i)"}, {"1", "0,0857"}, {}, {"n(0)", "35"}}
	if err := WriteCSV(&buffer, records); err != nil {
		t.Fatal(err)
	}
	want := utf8BOM + "t_i;F*(t_i)\r\n1;0,0857\r\n\r\nn(0);35\r\n"
	if got := buffer.String(); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestSavePNGWritesDecodableImage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "chart.png")
	if err := SavePNG(path, image.NewRGBA(image.Rect(0, 0, 4, 3))); err != nil {
		t.Fatal(err)
	}
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		t.Fatal(err)
	}
	if size := img.Bounds().Size(); size.X != 4 || size.Y != 3 {
		t.Errorf("decoded size = %v, want 4x3", size)
	}
}

func TestSavePNGReportsUnwritablePath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "chart.png")
	if err := SavePNG(path, image.NewRGBA(image.Rect(0, 0, 1, 1))); err == nil {
		t.Error("want an error for a path in a missing directory")
	}
}
