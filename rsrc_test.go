package winres

import (
	"bytes"
	"testing"
)

func Test_ResourceSet_write_WriteErr(t *testing.T) {
	rs := ResourceSet{}
	rs.Set(Name("NAME"), Name("NAME"), 0, make([]byte, 6))

	if _, err := rs.write(&badWriter{15}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{23}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{39}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{47}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{63}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{71}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{87}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{103}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{111}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	rs.Set(ID(1), ID(1), 0, make([]byte, 6))

	if _, err := rs.write(&badWriter{31}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}

	if _, err := rs.write(&badWriter{79}); err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}
}

func Test_dataEntry_writeData(t *testing.T) {
	expected := [][]byte{
		{},
		{1, 0, 0, 0, 0, 0, 0, 0},
		{1, 2, 0, 0, 0, 0, 0, 0},
		{1, 2, 3, 0, 0, 0, 0, 0},
		{1, 2, 3, 4, 0, 0, 0, 0},
		{1, 2, 3, 4, 5, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, 0, 0},
		{1, 2, 3, 4, 5, 6, 7, 0},
		{1, 2, 3, 4, 5, 6, 7, 8},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 0, 0, 0, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0, 0, 0, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0, 0, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 0, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 0, 0, 0},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 0, 0},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 0},
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
	}

	for i := 0; i <= 16; i++ {
		de := dataEntry{data: make([]byte, i)}
		for j := 0; j < i; j++ {
			de.data[j] = byte(j + 1)
		}
		buf := &bytes.Buffer{}

		err := de.writeData(buf)
		if err != nil || !bytes.Equal(buf.Bytes(), expected[i]) {
			t.Fail()
		}
		err = de.writeData(&badWriter{i - 1})
		if err == nil || err.Error() != errWrite+" -1" {
			t.Fail()
		}
	}
}
