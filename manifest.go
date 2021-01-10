package winres

import (
	"bytes"
	"text/template"
)

type appManifestCommon struct {
	Name                              string // Application name
	Version                           string // Application version
	Description                       string // Application description
	UIAccess                          bool   // Require access to other applications' UI elements
	AutoElevate                       bool
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

// AppManifest describes an application manifest.
//
// Its zero value corresponds to the most common case.
type AppManifest struct {
	appManifestCommon
	ExecutionLevel ExecutionLevel
	Compatibility  OSCompatibility
	DPIAwareness   DPIAwareness
}

type DPIAwareness int

const (
	DPIAware DPIAwareness = iota
	DPIUnaware
	DPIPerMonitor
	DPIPerMonitorV2
)

type OSCompatibility int

const (
	Win7AndAbove OSCompatibility = iota
	Win8AndAbove
	Win81AndAbove
	Win10AndAbove
	WinVistaAndAbove
)

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

  <assemblyIdentity type="win32" name="{{.Name | html}}" version="{{.Version | html}}" processorArchitecture="*"/>
{{- if ne .Description ""}}
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
		appManifestCommon
		SupportedOS    []string
		DPIAware       string
		DPIAwareness   string
		ExecutionLevel string
	}{appManifestCommon: manifest.appManifestCommon}

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
