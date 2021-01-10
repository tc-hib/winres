package winres

import (
	"github.com/tc-hib/winres/version"
	"io"
)

// Standard type IDs from  https://docs.microsoft.com/en-us/windows/win32/menurc/resource-types
const (
	RT_CURSOR       ID = 1
	RT_BITMAP       ID = 2
	RT_ICON         ID = 3
	RT_MENU         ID = 4
	RT_DIALOG       ID = 5
	RT_STRING       ID = 6
	RT_FONTDIR      ID = 7
	RT_FONT         ID = 8
	RT_ACCELERATOR  ID = 9
	RT_RCDATA       ID = 10
	RT_MESSAGETABLE ID = 11
	RT_GROUP_CURSOR ID = 12
	RT_GROUP_ICON   ID = 14
	RT_VERSION      ID = 16
	RT_PLUGPLAY     ID = 19
	RT_VXD          ID = 20
	RT_ANICURSOR    ID = 21
	RT_ANIICON      ID = 22
	RT_HTML         ID = 23
	RT_MANIFEST     ID = 24
)

const (
	LCIDNeutral = 0
	LCIDDefault = 0x409 // en-US is default
)

// Arch defines the target architecture.
// Its value can be used as a target suffix too: "rsrc_windows_" + string(arch) + ".syso"
type Arch string

const (
	ArchI386  Arch = "386"
	ArchAMD64 Arch = "amd64"
	ArchARM   Arch = "arm"
	ArchARM64 Arch = "arm64"
)

// ResourceSet is the main object in the package.
//
// Create an empty ResourceSet and call Set methods to add resources, then WriteObject to produce a COFF object file.
type ResourceSet struct {
	types        map[Identifier]*typeEntry
	lastIconID   uint16
	lastCursorID uint16
}

// Set adds or replaces a resource.
//
// typeID is the resource type's identifier.
// It can be either a standard type number (RT_ICON, RT_VERSION, ...) or any type name.
//
// resID is the resource's unique identifier for a given type.
// It can either be an ID starting from 1, or a Name.
//
// A resource ID can have different data depending on the user's locale.
// In this case Set can be called several times with the same resID but a different language ID.
//
// langID can be 0 (neutral), or one of those LCID:
// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-lcid
//
// Warning: the ResourceSet takes ownership of the data parameter.
// The caller should not write into it anymore after calling this method.
//
func (rs *ResourceSet) Set(typeID, resID Identifier, langID uint16, data []byte) error {
	if err := checkIdentifier(resID); err != nil {
		return err
	}
	if err := checkIdentifier(typeID); err != nil {
		return err
	}

	rs.getEntry(typeID, resID).data[ID(langID)] = &dataEntry{data}
	return nil
}

// SetVersionInfo set the VersionInfo structure.
func (rs *ResourceSet) SetVersionInfo(vi version.Info) {
	for langID, res := range vi.SplitTranslations() {
		rs.Set(RT_VERSION, ID(1), langID, res.Bytes())
	}
}

// SetManifest is a simplified way to embed a typical application manifest.
func (rs *ResourceSet) SetManifest(manifest AppManifest) {
	rs.Set(RT_MANIFEST, ID(1), LCIDDefault, makeManifest(manifest))
}

// WriteObject writes a full object file into w.
func (rs *ResourceSet) WriteObject(w io.Writer, arch Arch) error {
	return writeObject(w, rs, arch)
}

// Count returns the number of resources declared in the resource set.
func (rs *ResourceSet) Count() int {
	return rs.numDataEntries()
}

// Walk walks through the resources in same order as they will be written.
//
// It takes a callback function that takes same parameters as Set and returns a bool that should be true to continue, false to stop.
//
// If you add resources during the walk, they will be ignored.
// You should never call WriteObject during the walk.
func (rs *ResourceSet) Walk(f func(typeID, resID Identifier, langID uint16, data []byte) bool) {
	s := &state{}
	rs.order(s)
	for _, tk := range s.orderedKeys {
		te := rs.types[tk]
		for _, rk := range te.orderedKeys {
			re := te.resources[rk]
			for _, dk := range re.orderedKeys {
				if !f(tk, rk, uint16(dk), re.data[dk].data) {
					return
				}
			}
		}
	}
}

// WalkType walks through the resources or a certain type, in same order as they will be written.
//
// It takes a callback function that takes same parameters as Set and returns a bool that should be true to continue, false to stop.
//
// If you add resources during the walk, they will be ignored.
// You should never call WriteObject during the walk.
func (rs *ResourceSet) WalkType(typeID Identifier, f func(resID Identifier, langID uint16, data []byte) bool) {
	te := rs.types[typeID]
	if te == nil {
		return
	}
	te.order()
	for _, rk := range te.orderedKeys {
		re := te.resources[rk]
		for _, dk := range re.orderedKeys {
			if !f(rk, uint16(dk), re.data[dk].data) {
				return
			}
		}
	}
}

// Get returns resource data.
// Returns nil if the resource (translated resource) was not found.
func (rs *ResourceSet) Get(typeID, resID Identifier, langID uint16) []byte {
	te := rs.types[typeID]
	if te == nil {
		return nil
	}

	re := te.resources[resID]
	if re == nil {
		return nil
	}

	de := re.data[ID(langID)]
	if de == nil {
		return nil
	}

	return de.data
}

// getEntry creates, if necessary, and returns a resource entry.
func (rs *ResourceSet) getEntry(typeID Identifier, resID Identifier) *resourceEntry {
	if rs.types == nil {
		rs.types = make(map[Identifier]*typeEntry)
	}
	te := rs.types[typeID]
	if te == nil {
		te = &typeEntry{
			resources: make(map[Identifier]*resourceEntry),
		}
		rs.types[typeID] = te
	}
	if te.resources[resID] == nil {
		te.resources[resID] = &resourceEntry{
			data: make(map[ID]*dataEntry),
		}
	}
	return te.resources[resID]
}
