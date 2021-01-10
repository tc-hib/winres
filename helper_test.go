package winres

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func checkResourceSet(t *testing.T, rs *ResourceSet, arch Arch) {
	buf := &bytes.Buffer{}
	if err := rs.WriteObject(buf, arch); err != nil {
		t.Fatal(err)
	}
	checkBinary(t, buf.Bytes())
}

func checkBinary(t *testing.T, data []byte) {
	refFile := filepath.Join("testdata", t.Name()+".golden")
	ref, _ := ioutil.ReadFile(refFile)

	if !bytes.Equal(ref, data) {
		t.Error(t.Name() + " output is different")
		bugFile := refFile[:len(refFile)-7] + ".bug"
		err := ioutil.WriteFile(bugFile, data, 0666)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("dumped output to", bugFile)
	}
}

func loadPNGFileAsIcon(t *testing.T, name string, sizes []int) *Icon {
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	icon, err := NewIconFromResizedImage(img, sizes)
	if err != nil {
		t.Fatal(err)
	}
	return icon
}

func loadPNGFileAsCursor(t *testing.T, name string, spotX, spotY uint16) *Cursor {
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	cursor, err := NewCursorFromImages([]CursorImage{{img, HotSpot{spotX, spotY}}})
	if err != nil {
		t.Fatal(err)
	}
	return cursor
}

func loadICOFile(t *testing.T, name string) *Icon {
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	icon, err := LoadICO(f)
	if err != nil {
		t.Fatal(err)
	}
	return icon
}

func loadCURFile(t *testing.T, name string) *Cursor {
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	cursor, err := LoadCUR(f)
	if err != nil {
		t.Fatal(err)
	}
	return cursor
}

func loadImage(t *testing.T, name string) image.Image {
	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		t.Fatal(err)
	}

	return img
}

func shiftImage(img image.Image, x, y int) image.Image {
	shifted := image.NewNRGBA(image.Rectangle{
		Min: image.Point{
			X: img.Bounds().Min.X + x,
			Y: img.Bounds().Min.Y + y,
		},
		Max: image.Point{
			X: img.Bounds().Max.X + x,
			Y: img.Bounds().Max.Y + y,
		},
	})

	for srcY := img.Bounds().Min.Y; srcY < img.Bounds().Max.Y; srcY++ {
		for srcX := img.Bounds().Min.X; srcX < img.Bounds().Max.X; srcX++ {
			shifted.Set(srcX+x, srcY+y, img.At(srcX, srcY))
		}
	}

	return shifted
}

type badReader struct {
	br            *bytes.Reader
	seekCountdown int
	readCountdown int
}

type badWriter struct {
	writeCountdown int
}

type badError struct {
	f     string
	count int
}

func (e *badError) Error() string {
	return fmt.Sprintf("%s %d", e.f, e.count)
}

const (
	errRead  = "read"
	errSeek  = "seek"
	errWrite = "write"
)

func (r *badReader) Read(b []byte) (n int, err error) {
	r.readCountdown -= len(b)
	if r.readCountdown <= 0 {
		return 0, &badError{errRead, r.readCountdown}
	}
	return r.br.Read(b)
}

func (r *badReader) Seek(offset int64, whence int) (int64, error) {
	if r.seekCountdown <= 0 {
		return 0, &badError{errSeek, r.seekCountdown}
	}
	r.seekCountdown--
	return r.br.Seek(offset, whence)
}

func (r *badWriter) Write(b []byte) (n int, err error) {
	r.writeCountdown -= len(b)
	if r.writeCountdown <= 0 {
		return 0, &badError{errWrite, r.writeCountdown}
	}
	return len(b), nil
}
