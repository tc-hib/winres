package winres

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const testDataDir = "testdata"

func checkResourceSet(t *testing.T, rs *ResourceSet, arch Arch) {
	buf := &bytes.Buffer{}
	if err := rs.WriteObject(buf, arch); err != nil {
		t.Fatal(err)
	}
	checkBinary(t, buf.Bytes())
}

func golden(t *testing.T) string {
	return filepath.Join(testDataDir, t.Name()+".golden")
}

func checkBinary(t *testing.T, data []byte) {
	refFile := golden(t)
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

func loadBinary(t *testing.T, filename string) []byte {
	data, err := ioutil.ReadFile(filepath.Join(testDataDir, filename))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func loadPNGFileAsIcon(t *testing.T, name string, sizes []int) *Icon {
	f, err := os.Open(filepath.Join(testDataDir, name))
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
	f, err := os.Open(filepath.Join(testDataDir, name))
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
	f, err := os.Open(filepath.Join(testDataDir, name))
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
	f, err := os.Open(filepath.Join(testDataDir, name))
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
	f, err := os.Open(filepath.Join(testDataDir, name))
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
	br     *bytes.Reader
	errPos int64
}

type badSeeker struct {
	br      *bytes.Reader
	errIter int
}

type badWriter struct {
	badLen int
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
	p, _ := r.br.Seek(0, io.SeekCurrent)
	if p <= r.errPos && r.errPos < p+int64(len(b)) {
		n, _ := r.br.Read(b[:r.errPos-p])
		return n, errors.New(errRead)
	}
	return r.br.Read(b)
}

func (r *badReader) Seek(offset int64, whence int) (int64, error) {
	return r.br.Seek(offset, whence)
}

func (r *badSeeker) Read(b []byte) (n int, err error) {
	return r.br.Read(b)
}

func (r *badSeeker) Seek(offset int64, whence int) (int64, error) {
	if r.errIter <= 0 {
		return 0, errors.New(errSeek)
	}
	r.errIter--
	return r.br.Seek(offset, whence)
}

func (r *badWriter) Write(b []byte) (n int, err error) {
	r.badLen -= len(b)
	if r.badLen <= 0 {
		return 0, &badError{errWrite, r.badLen}
	}
	return len(b), nil
}
