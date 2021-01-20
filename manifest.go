package winres

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

// AppManifest describes an application manifest.
//
// Its zero value corresponds to the most common case.
type AppManifest struct {
	Identity                          AssemblyIdentity
	Description                       string
	Compatibility                     SupportedOS
	ExecutionLevel                    ExecutionLevel
	UIAccess                          bool // Require access to other applications' UI elements
	AutoElevate                       bool
	DPIAwareness                      DPIAwareness
	DisableTheming                    bool
	DisableWindowFiltering            bool
	HighResolutionScrollingAware      bool
	UltraHighResolutionScrollingAware bool
	LongPathAware                     bool
	PrinterDriverIsolation            bool
	GDIScaling                        bool
	SegmentHeap                       bool
	UseCommonControlsV6               bool // Application requires Common Controls V6 (V5 remains the default)
}

// AssemblyIdentity defines the side-by-side assembly identity of the executable.
//
// It should not be needed unless another assembly depends on this one.
//
// If the Name field is empty, the <assemblyIdentity> element will be omitted.
type AssemblyIdentity struct {
	Name    string
	Version [4]uint16
}

// DPIAwareness is an enumeration which corresponds to the <dpiAware> and the <dpiAwareness> elements.
//
// When it is set to DPIPerMonitorV2, it will fallback to DPIAware if the OS does not support it.
//
// DPIPerMonitor would not scale windows on secondary monitors.
type DPIAwareness int

const (
	DPIAware DPIAwareness = iota
	DPIUnaware
	DPIPerMonitor
	DPIPerMonitorV2
)

// SupportedOS is an enumeration that provides a simplified way to fill the
// compatibility element in an application manifest, by only setting a minimum OS.
//
// Its zero value is Win7AndAbove, which matches Go's requirements.
//
// https://github.com/golang/go/wiki/MinimumRequirements#windows
type SupportedOS int

const (
	WinVistaAndAbove SupportedOS = iota - 1
	Win7AndAbove
	Win8AndAbove
	Win81AndAbove
	Win10AndAbove
)

// ExecutionLevel is used in an AppManifest to set the required execution level.
type ExecutionLevel int

const (
	AsInvoker ExecutionLevel = iota
	HighestAvailable
	RequireAdministrator
)

const (
	osWin10    = "{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"
	osWin81    = "{1f676c76-80e1-4239-95bb-83d0f6d0da78}"
	osWin8     = "{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"
	osWin7     = "{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"
	osWinVista = "{e2011457-1546-43c5-a5fe-008deee3d3f0}"
)

// language=GoTemplate
var manifestTemplate = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
{{- if .AssemblyName}}

  <assemblyIdentity type="win32" name="{{.AssemblyName | html}}" version="{{.AssemblyVersion}}" processorArchitecture="*"/>
{{- end}}
{{- if .Description}}
  <description>{{.Description | html}}</description>
{{- end}}

  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      {{- range $osID := .SupportedOS}}
      <supportedOS Id="{{$osID}}"/>
      {{- end}}
    </application>
  </compatibility>

  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">{{.DPIAware}}</dpiAware>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">{{.DPIAwareness}}</dpiAwareness>
      {{- if .AutoElevate}}
      <autoElevate xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</autoElevate>
      {{- end}}
      {{- if .DisableTheming}}
      <disableTheming xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</disableTheming>
      {{- end}}
      {{- if .DisableWindowFiltering}}
      <disableWindowFiltering xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</disableWindowFiltering>
      {{- end}}
      {{- if .HighResolutionScrollingAware}}
      <highResolutionScrollingAware xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</highResolutionScrollingAware>
      {{- end}}
      {{- if .PrinterDriverIsolation}}
      <printerDriverIsolation xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</printerDriverIsolation>
      {{- end}}
      {{- if .UltraHighResolutionScrollingAware}}
      <ultraHighResolutionScrollingAware xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</ultraHighResolutionScrollingAware>
      {{- end}}
      {{- if .LongPathAware}}
      <longPathAware xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">true</longPathAware>
      {{- end}}
      {{- if .GDIScaling}}
      <gdiScaling xmlns="http://schemas.microsoft.com/SMI/2017/WindowsSettings">true</gdiScaling>
      {{- end}}
      {{- if .SegmentHeap}}
      <heapType xmlns="http://schemas.microsoft.com/SMI/2020/WindowsSettings">SegmentHeap</heapType>
      {{- end}}
    </windowsSettings>
  </application>

  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="{{.ExecutionLevel}}" uiAccess="{{if .UIAccess}}true{{else}}false{{end}}"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
  {{- if .UseCommonControlsV6}}

  <dependency>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
    </dependentAssembly>
  </dependency>
  {{- end}}

</assembly>
`

func makeManifest(manifest AppManifest) []byte {
	vars := struct {
		AppManifest
		AssemblyName    string
		AssemblyVersion string
		SupportedOS     []string
		DPIAware        string
		DPIAwareness    string
		ExecutionLevel  string
	}{AppManifest: manifest}

	if manifest.Identity.Name != "" {
		vars.AssemblyName = manifest.Identity.Name
		v := manifest.Identity.Version
		vars.AssemblyVersion = fmt.Sprintf("%d.%d.%d.%d", v[0], v[1], v[2], v[3])
	}

	vars.SupportedOS = []string{
		osWin10,
		osWin81,
		osWin8,
		osWin7,
		osWinVista,
	}
	switch manifest.Compatibility {
	case Win7AndAbove:
		vars.SupportedOS = vars.SupportedOS[:4]
	case Win8AndAbove:
		vars.SupportedOS = vars.SupportedOS[:3]
	case Win81AndAbove:
		vars.SupportedOS = vars.SupportedOS[:2]
	case Win10AndAbove:
		vars.SupportedOS = vars.SupportedOS[:1]
	}

	switch manifest.ExecutionLevel {
	case RequireAdministrator:
		vars.ExecutionLevel = "requireAdministrator"
	case HighestAvailable:
		vars.ExecutionLevel = "highestAvailable"
	default:
		vars.ExecutionLevel = "asInvoker"
	}

	switch manifest.DPIAwareness {
	case DPIAware:
		vars.DPIAware = "true"
		vars.DPIAwareness = "system"
	case DPIPerMonitor:
		vars.DPIAware = "true/pm"
		vars.DPIAwareness = "permonitor"
	case DPIPerMonitorV2:
		// PerMonitorV2 fixes the scale on secondary monitors
		// If not available, the closest option seems to be System
		vars.DPIAware = "true"
		vars.DPIAwareness = "permonitorv2,system"
	case DPIUnaware:
		vars.DPIAware = "false"
		vars.DPIAwareness = "unaware"
	}

	buf := &bytes.Buffer{}
	tmpl := template.Must(template.New("manifest").Parse(manifestTemplate))
	err := tmpl.Execute(buf, vars)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

type appManifestXML struct {
	Identity struct {
		Name    string `xml:"name,attr"`
		Version string `xml:"version,attr"`
	} `xml:"assemblyIdentity"`
	Description   string `xml:"description"`
	Compatibility struct {
		Application struct {
			SupportedOS []struct {
				Id string `xml:"Id,attr"`
			} `xml:"supportedOS"`
		} `xml:"application"`
	} `xml:"compatibility"`
	Application struct {
		WindowsSettings struct {
			DPIAware                          string `xml:"dpiAware"`
			DPIAwareness                      string `xml:"dpiAwareness"`
			AutoElevate                       string `xml:"autoElevate"`
			DisableTheming                    string `xml:"disableTheming"`
			DisableWindowFiltering            string `xml:"disableWindowFiltering"`
			HighResolutionScrollingAware      string `xml:"highResolutionScrollingAware"`
			PrinterDriverIsolation            string `xml:"printerDriverIsolation"`
			UltraHighResolutionScrollingAware string `xml:"ultraHighResolutionScrollingAware"`
			LongPathAware                     string `xml:"longPathAware"`
			GDIScaling                        string `xml:"gdiScaling"`
			HeapType                          string `xml:"heapType"`
		} `xml:"windowsSettings"`
	} `xml:"application"`
	TrustInfo struct {
		Security struct {
			RequestedPrivileges struct {
				RequestedExecutionLevel struct {
					Level    string `xml:"level,attr"`
					UIAccess string `xml:"uiAccess,attr"`
				} `xml:"requestedExecutionLevel"`
			} `xml:"requestedPrivileges"`
		} `xml:"security"`
	} `xml:"trustInfo"`
	Dependency struct {
		DependentAssembly []struct {
			Identity struct {
				Name           string `xml:"name,attr"`
				Version        string `xml:"version,attr"`
				PublicKeyToken string `xml:"publicKeyToken,attr"`
			} `xml:"assemblyIdentity"`
		} `xml:"dependentAssembly"`
	} `xml:"dependency"`
}

// AppManifestFromXML makes an AppManifest from an xml manifest,
// trying to retrieve as much valid information as possible.
//
// If the xml contains other data, they are ignored.
//
// This function can only return xml syntax errors, other errors are ignored.
func AppManifestFromXML(data []byte) (AppManifest, error) {
	x := appManifestXML{}
	err := xml.Unmarshal(data, &x)
	if err != nil {
		return AppManifest{}, err
	}
	var m AppManifest

	m.Identity.Name = x.Identity.Name
	v := strings.Split(x.Identity.Version, ".")
	if len(v) > 4 {
		v = v[:4]
	}
	for i := range v {
		n, _ := strconv.ParseUint(v[i], 10, 16)
		m.Identity.Version[i] = uint16(n)
	}
	m.Description = x.Description

	m.Compatibility = Win10AndAbove + 1
	for _, os := range x.Compatibility.Application.SupportedOS {
		c := osIDToEnum(os.Id)
		if c < m.Compatibility {
			m.Compatibility = c
		}
	}
	if m.Compatibility > Win10AndAbove {
		m.Compatibility = Win7AndAbove
	}

	settings := x.Application.WindowsSettings
	m.DPIAwareness = readDPIAwareness(settings.DPIAware, settings.DPIAwareness)
	m.AutoElevate = manifestBool(settings.AutoElevate)
	m.DisableTheming = manifestBool(settings.DisableTheming)
	m.DisableWindowFiltering = manifestBool(settings.DisableWindowFiltering)
	m.HighResolutionScrollingAware = manifestBool(settings.HighResolutionScrollingAware)
	m.PrinterDriverIsolation = manifestBool(settings.PrinterDriverIsolation)
	m.UltraHighResolutionScrollingAware = manifestBool(settings.UltraHighResolutionScrollingAware)
	m.LongPathAware = manifestBool(settings.LongPathAware)
	m.GDIScaling = manifestBool(settings.GDIScaling)
	m.SegmentHeap = manifestString(settings.HeapType) == "segmentheap"

	for _, dep := range x.Dependency.DependentAssembly {
		if manifestString(dep.Identity.Name) == "microsoft.windows.common-controls" &&
			strings.HasPrefix(manifestString(dep.Identity.Version), "6.") &&
			manifestString(dep.Identity.PublicKeyToken) == "6595b64144ccf1df" {
			m.UseCommonControlsV6 = true
		}
	}

	m.UIAccess = manifestBool(x.TrustInfo.Security.RequestedPrivileges.RequestedExecutionLevel.UIAccess)
	switch manifestString(x.TrustInfo.Security.RequestedPrivileges.RequestedExecutionLevel.Level) {
	case "requireadministrator":
		m.ExecutionLevel = RequireAdministrator
	case "highestavailable":
		m.ExecutionLevel = HighestAvailable
	}

	return m, nil
}

func readDPIAwareness(dpiAware string, dpiAwareness string) DPIAwareness {
	for _, s := range strings.Split(dpiAwareness, ",") {
		switch manifestString(s) {
		case "permonitorv2":
			return DPIPerMonitorV2
		case "permonitor":
			return DPIPerMonitor
		case "system":
			return DPIAware
		case "unaware":
			return DPIUnaware
		}
	}
	switch manifestString(dpiAware) {
	case "true":
		return DPIAware
	case "true/pm":
		return DPIPerMonitor
	}
	return DPIUnaware
}

func manifestString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func manifestBool(s string) bool {
	return manifestString(s) == "true"
}

func osIDToEnum(osID string) SupportedOS {
	switch osID {
	case osWinVista:
		return WinVistaAndAbove
	case osWin7:
		return Win7AndAbove
	case osWin8:
		return Win8AndAbove
	case osWin81:
		return Win81AndAbove
	}
	return Win10AndAbove
}
