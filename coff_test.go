package winres

import (
	"io/ioutil"
	"testing"
)

func Test_writeObject_WriteErrFileHeader(t *testing.T) {
	bw := &badWriter{19}
	rs := &ResourceSet{}

	err := writeObject(bw, rs, ArchAMD64)
	if err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}
}

func Test_writeObject_WriteErrSectionHeader(t *testing.T) {
	bw := &badWriter{58}
	rs := &ResourceSet{}

	err := writeObject(bw, rs, ArchAMD64)
	if err == nil || err.Error() != errWrite+" -2" {
		t.Fail()
	}
}

func Test_writeObject_WriteErrReloc(t *testing.T) {
	bw := &badWriter{165}
	rs := &ResourceSet{}
	rs.Set(RT_RCDATA, ID(1), 0, make([]byte, 1))

	err := writeObject(bw, rs, ArchAMD64)
	if err == nil || err.Error() != errWrite+" -1" {
		t.Fail()
	}
}

func Test_writeObject_WriteErrSymbol(t *testing.T) {
	bw := &badWriter{182}
	rs := &ResourceSet{}
	rs.Set(RT_RCDATA, ID(1), 0, make([]byte, 1))

	err := writeObject(bw, rs, ArchAMD64)
	if err == nil || err.Error() != errWrite+" -2" {
		t.Fail()
	}
}

func Test_writeObject_WriteErrStringTable(t *testing.T) {
	bw := &badWriter{185}
	rs := &ResourceSet{}
	rs.Set(RT_RCDATA, ID(1), 0, make([]byte, 1))

	err := writeObject(bw, rs, ArchAMD64)
	if err == nil || err.Error() != errWrite+" -3" {
		t.Fail()
	}
}

func Test_writeObject_WriteErrSection(t *testing.T) {
	bw := &badWriter{61}
	rs := &ResourceSet{}
	rs.Set(RT_RCDATA, ID(1), 0, make([]byte, 1))

	err := writeObject(bw, rs, ArchAMD64)
	if err == nil || err.Error() != errWrite+" -15" {
		t.Fail()
	}
}

func Test_writeRelocTable_UnknownArch(t *testing.T) {
	err := writeRelocTable(ioutil.Discard, 1, "*", []int{1, 2, 3, 4})
	if err == nil || err.Error() != errUnknownArch {
		t.Fail()
	}
}
