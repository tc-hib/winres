package winres

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tc-hib/winres/version"
	"io/ioutil"
	"testing"
	"time"
)

func TestErrors(t *testing.T) {
	r := &ResourceSet{}
	var err error

	err = r.Set(RT_RCDATA, ID(0), 0, []byte{})
	if err == nil || err.Error() != errZeroID {
		t.Fail()
	}

	err = r.Set(ID(0), ID(1), 0, []byte{})
	if err == nil || err.Error() != errZeroID {
		t.Fail()
	}

	if r.Set(RT_RCDATA, ID(0xFFFF), 0, []byte{}) != nil {
		t.Fail()
	}

	err = r.Set(RT_RCDATA, Name(""), 0, []byte{})
	if err == nil || err.Error() != errEmptyName {
		t.Fail()
	}

	err = r.Set(Name(""), ID(1), 0, []byte{})
	if err == nil || err.Error() != errEmptyName {
		t.Fail()
	}

	if r.Set(RT_RCDATA, Name("look, i'm not a nice resource name"), 0, []byte{}) != nil {
		t.Fail()
	}

	err = r.Set(RT_RCDATA, Name("IAMNICER\x00"), 0, []byte{})
	if err == nil || err.Error() != errNameContainsNUL {
		t.Fail()
	}

	if r.Set(Name("look, i'm not a nice type name"), ID(1), 0, []byte{}) != nil {
		t.Fail()
	}

	err = r.Set(Name("IAMNICER\x00"), ID(42), 0, []byte{})
	if err == nil || err.Error() != errNameContainsNUL {
		t.Fail()
	}

	err = r.WriteObject(ioutil.Discard, "*")
	if err == nil || err.Error() != errUnknownArch {
		t.Fail()
	}
}

func TestEmpty(t *testing.T) {
	rs := &ResourceSet{}
	checkResourceSet(t, rs, ArchI386)
}

func TestWinRes1(t *testing.T) {
	r := &ResourceSet{}

	r.Set(RT_MANIFEST, ID(1), LCIDDefault, []byte(manifest1))
	r.Set(Name("CUSTOM TYPE"), Name("CUSTOM RESOURCE"), 1033, []byte("Hello World!"))
	r.Set(Name("CUSTOM TYPE"), Name("CUSTOM RESOURCE"), 1036, []byte("Bonjour Monde !"))
	r.Set(Name("CUSTOM TYPE"), ID(42), 1033, []byte("# Hello World!"))
	r.Set(Name("CUSTOM TYPE"), ID(42), 1036, []byte("# Bonjour Monde !"))
	r.Set(RT_RCDATA, ID(1), 1033, []byte("## Hello World!"))
	r.Set(RT_RCDATA, ID(1), 1036, []byte("## Bonjour Monde !"))
	r.Set(RT_RCDATA, ID(42), 1033, []byte("### Hello World!"))
	r.Set(RT_RCDATA, ID(42), 1036, []byte("### Bonjour Monde !"))
	icon1 := loadICOFile(t, "icon.ico")
	cursor := loadCURFile(t, "cursor.cur")
	icon2 := loadPNGFileAsIcon(t, "cur-64x128.png", nil)
	icon3 := loadPNGFileAsIcon(t, "cur-32x64.png", []int{48, 16})
	icon4 := loadPNGFileAsIcon(t, "img.png", []int{128})
	r.SetIconTranslation(ID(1), 0, icon1)
	r.SetIconTranslation(ID(1), 1033, icon2)
	r.SetIconTranslation(ID(1), 1036, icon3)
	r.SetIconTranslation(Name("SUPERB ICON"), 0, icon4)
	r.SetIcon(ID(2), icon2)
	r.SetCursor(ID(1), cursor)
	v := version.Info{
		ProductVersion: [4]uint16{5, 6, 7, 8},
	}
	v.Set(1036, "Custom Info", "Very important information")
	v.Set(1036, version.ProductName, "A test for winres")
	v.Set(1036, version.ProductVersion, "0.0.0.0-αlpha-")
	v.Set(1036, version.CompanyName, "Test Corporation ltd")
	v.Flags.SpecialBuild = true
	v.FileVersion = [4]uint16{4, 42, 424, 4242}
	v.Timestamp = time.Date(1979, 7, 3, 0, 15, 0, 0, time.UTC)
	r.SetVersionInfo(v)

	checkResourceSet(t, r, ArchAMD64)
}

func TestWinRes2(t *testing.T) {
	r := &ResourceSet{}

	r.Set(RT_MANIFEST, ID(1), LCIDDefault, []byte(manifest1))

	j, _ := ioutil.ReadFile("testdata/vi.json")
	v := version.Info{}
	json.Unmarshal(j, &v)
	r.SetVersionInfo(v)

	r.SetIcon(ID(1), loadICOFile(t, "icon.ico"))

	checkResourceSet(t, r, ArchI386)
}

func TestWinRes3(t *testing.T) {
	rs := &ResourceSet{}
	rs.SetCursor(ID(1), loadCURFile(t, "cursor.cur"))
	checkResourceSet(t, rs, ArchARM)
}

func TestWinRes4(t *testing.T) {
	rs := &ResourceSet{}
	rs.SetCursor(ID(1), loadPNGFileAsCursor(t, "cur-32x64.png", 10, 7))
	rs.SetIcon(ID(1), loadPNGFileAsIcon(t, "cur-32x64.png", []int{1, 7, 11, 15, 22, 255, 256}))
	checkResourceSet(t, rs, ArchARM64)
}

func TestResourceSet_Count(t *testing.T) {
	rs := &ResourceSet{}
	rs.SetManifest(AppManifest{})
	rs.SetManifest(AppManifest{appManifestCommon: appManifestCommon{Name: "Hello"}})
	rs.Set(RT_RCDATA, ID(42), 0x40C, make([]byte, 8))
	rs.Set(RT_RCDATA, ID(42), 0x40C, make([]byte, 9))
	rs.Set(RT_RCDATA, Name("Data"), 0x40C, make([]byte, 6))
	rs.Set(RT_RCDATA, ID(42), 0x409, make([]byte, 7))
	rs.Set(RT_VERSION, ID(1), 0x409, make([]byte, 9))
	rs.Set(RT_CURSOR, ID(42), 0x409, make([]byte, 5))
	rs.Set(Name("1"), ID(1), 0x409, make([]byte, 1))
	if rs.Count() != 7 {
		t.Fail()
	}
}

func TestResourceSet_SetManifest(t *testing.T) {
	rs := &ResourceSet{}
	rs.SetManifest(AppManifest{})
	checkResourceSet(t, rs, ArchARM64)
}

func TestResourceSet_SetVersionInfo(t *testing.T) {
	rs := &ResourceSet{}
	vi := version.Info{}
	vi.FileVersion = [4]uint16{1, 2, 3, 4}
	vi.ProductVersion = [4]uint16{1, 2, 3, 4}
	vi.Set(0x0409, version.ProductName, "Good product")
	vi.Set(0x040C, version.ProductName, "Bon produit")
	rs.SetVersionInfo(vi)
	checkResourceSet(t, rs, ArchAMD64)
}

func TestResourceSet_Walk(t *testing.T) {
	rs := ResourceSet{}
	b := &bytes.Buffer{}

	walker := func(typeID, resID Identifier, langID uint16, data []byte) bool {
		fmt.Fprintf(b, "%T(%v) -> %T(%v) -> 0x%04X -> [%d]byte\n", typeID, typeID, resID, resID, langID, len(data))
		return resID != Name("STOP")
	}

	rs.Walk(walker)
	if b.String() != "" {
		t.Fail()
	}

	rs.Set(RT_RCDATA, ID(42), 0x40C, make([]byte, 8))
	rs.Set(RT_RCDATA, Name("Data"), 0x40C, make([]byte, 6))
	rs.Set(RT_RCDATA, ID(42), 0x409, make([]byte, 7))
	rs.Set(RT_VERSION, ID(1), 0x409, make([]byte, 9))
	rs.Set(RT_CURSOR, ID(42), 0x409, make([]byte, 5))
	rs.Set(Name("1"), ID(1), 0x409, make([]byte, 1))
	rs.Set(Name("1"), ID(2), 0x409, make([]byte, 2))
	rs.Set(Name("Hi"), ID(2), 0x409, make([]byte, 3))
	rs.Set(Name("hey"), ID(2), 0x409, make([]byte, 4))
	rs.Set(ID(99), Name("STOP"), 0x409, make([]byte, 4))
	rs.Set(ID(99), Name("TOO FAR"), 0x409, make([]byte, 4))
	rs.Walk(walker)
	expected := `winres.Name(1) -> winres.ID(1) -> 0x0409 -> [1]byte
winres.Name(1) -> winres.ID(2) -> 0x0409 -> [2]byte
winres.Name(Hi) -> winres.ID(2) -> 0x0409 -> [3]byte
winres.Name(hey) -> winres.ID(2) -> 0x0409 -> [4]byte
winres.ID(1) -> winres.ID(42) -> 0x0409 -> [5]byte
winres.ID(10) -> winres.Name(Data) -> 0x040C -> [6]byte
winres.ID(10) -> winres.ID(42) -> 0x0409 -> [7]byte
winres.ID(10) -> winres.ID(42) -> 0x040C -> [8]byte
winres.ID(16) -> winres.ID(1) -> 0x0409 -> [9]byte
winres.ID(99) -> winres.Name(STOP) -> 0x0409 -> [4]byte
`
	if b.String() != expected {
		t.Fail()
	}
}

func TestResourceSet_WalkType(t *testing.T) {
	rs := ResourceSet{}
	b := &bytes.Buffer{}

	walker := func(resID Identifier, langID uint16, data []byte) bool {
		fmt.Fprintf(b, "%T(%v) -> 0x%04X -> [%d]byte\n", resID, resID, langID, len(data))
		return resID != ID(999)
	}

	rs.WalkType(RT_RCDATA, walker)
	if b.String() != "" {
		t.Fail()
	}

	rs.Set(RT_RCDATA, ID(42), 0x401, make([]byte, 8))
	rs.Set(RT_RCDATA, Name("Data"), 0x402, make([]byte, 6))
	rs.Set(RT_RCDATA, ID(42), 0x403, make([]byte, 7))
	rs.Set(RT_RCDATA, ID(999), 0x404, make([]byte, 4))
	rs.Set(RT_RCDATA, ID(1000), 0x405, make([]byte, 4))
	rs.Set(RT_VERSION, ID(1), 0x409, make([]byte, 9))
	rs.Set(RT_CURSOR, ID(42), 0x409, make([]byte, 5))
	rs.Set(Name("1"), ID(1), 0x409, make([]byte, 1))
	rs.Set(Name("1"), ID(2), 0x409, make([]byte, 2))
	rs.Set(Name("Hi"), ID(2), 0x409, make([]byte, 3))
	rs.Set(Name("hey"), ID(2), 0x409, make([]byte, 4))
	rs.WalkType(RT_RCDATA, walker)
	expected := `winres.Name(Data) -> 0x0402 -> [6]byte
winres.ID(42) -> 0x0401 -> [8]byte
winres.ID(42) -> 0x0403 -> [7]byte
winres.ID(999) -> 0x0404 -> [4]byte
`

	if b.String() != expected {
		t.Fail()
	}
}

//language=xml
const manifest1 = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

<assemblyIdentity
	version="1.0.0.0"
	processorArchitecture="*"
	name="An App"
	type="win32"
/>

<application xmlns="urn:schemas-microsoft-com:asm.v3">
	<windowsSettings>
		<dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true/PM</dpiAware>
	</windowsSettings>
</application>

<trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
	<security>
		<requestedPrivileges>
			<requestedExecutionLevel
				level="asInvoker"
				uiAccess="false"
			/>
		</requestedPrivileges>
	</security>
</trustInfo>

<compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
	<application>
		<supportedOS Id="{e2011457-1546-43c5-a5fe-008deee3d3f0}"/>
		<supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"/>
		<supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
		<supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
		<supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
	</application>
</compatibility>

</assembly>
`