package winres

import (
	"reflect"
	"testing"
)

func Test_makeManifest(t *testing.T) {
	type args struct {
		manifest AppManifest
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "empty",
			args: struct{ manifest AppManifest }{},
			want: []byte(
				// language=xml
				`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

  <assemblyIdentity type="win32" name="" version="" processorArchitecture="*"/>

  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
      <supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
      <supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"/>
    </application>
  </compatibility>

  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true</dpiAware>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">system</dpiAwareness>
    </windowsSettings>
  </application>

  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>

</assembly>
`),
		},
		{name: "full", args: struct{ manifest AppManifest }{AppManifest{
			appManifestCommon: appManifestCommon{
				Name:                              "<app.name>",
				Version:                           "1.2.3.4&",
				Description:                       "<Application Description>",
				UIAccess:                          true,
				AutoElevate:                       true,
				DisableTheming:                    true,
				DisableWindowFiltering:            true,
				HighResolutionScrollingAware:      true,
				UltraHighResolutionScrollingAware: true,
				LongPathAware:                     true,
				PrinterDriverIsolation:            true,
				GDIScaling:                        true,
				SegmentHeap:                       true,
				UseCommonControlsV6:               true,
			},
			ExecutionLevel: HighestAvailable,
			Compatibility:  WinVistaAndAbove,
			DPIAwareness:   DPIAware,
		}}, want: []byte(
			// language=xml
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

  <assemblyIdentity type="win32" name="&lt;app.name&gt;" version="1.2.3.4&amp;" processorArchitecture="*"/>
  <description>&lt;Application Description&gt;</description>

  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
      <supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
      <supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"/>
      <supportedOS Id="{e2011457-1546-43c5-a5fe-008deee3d3f0}"/>
    </application>
  </compatibility>

  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true</dpiAware>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">system</dpiAwareness>
      <autoElevate xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</autoElevate>
      <disableTheming xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</disableTheming>
      <disableWindowFiltering xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</disableWindowFiltering>
      <highResolutionScrollingAware xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</highResolutionScrollingAware>
      <printerDriverIsolation xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</printerDriverIsolation>
      <ultraHighResolutionScrollingAware xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</ultraHighResolutionScrollingAware>
      <longPathAware xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">true</longPathAware>
      <gdiScaling xmlns="http://schemas.microsoft.com/SMI/2017/WindowsSettings">true</gdiScaling>
      <heapType xmlns="http://schemas.microsoft.com/SMI/2020/WindowsSettings">SegmentHeap</heapType>
    </windowsSettings>
  </application>

  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="highestAvailable" uiAccess="true"/>
      </requestedPrivileges>
    </security>
  </trustInfo>

  <dependency>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
    </dependentAssembly>
  </dependency>

</assembly>
`)},
		{name: "win10admin", args: struct{ manifest AppManifest }{AppManifest{
			appManifestCommon: appManifestCommon{
				Name:                              "app.name",
				Version:                           "1.2.3.4",
				Description:                       "Application Description",
				UIAccess:                          false,
				AutoElevate:                       true,
				DisableTheming:                    false,
				DisableWindowFiltering:            true,
				HighResolutionScrollingAware:      false,
				UltraHighResolutionScrollingAware: true,
				LongPathAware:                     false,
				PrinterDriverIsolation:            true,
				GDIScaling:                        false,
				SegmentHeap:                       true,
			},
			ExecutionLevel: RequireAdministrator,
			Compatibility:  Win10AndAbove,
			DPIAwareness:   DPIUnaware,
		}}, want: []byte(
			// language=xml
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

  <assemblyIdentity type="win32" name="app.name" version="1.2.3.4" processorArchitecture="*"/>
  <description>Application Description</description>

  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
    </application>
  </compatibility>

  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">false</dpiAware>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">unaware</dpiAwareness>
      <autoElevate xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</autoElevate>
      <disableWindowFiltering xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</disableWindowFiltering>
      <printerDriverIsolation xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</printerDriverIsolation>
      <ultraHighResolutionScrollingAware xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</ultraHighResolutionScrollingAware>
      <heapType xmlns="http://schemas.microsoft.com/SMI/2020/WindowsSettings">SegmentHeap</heapType>
    </windowsSettings>
  </application>

  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="requireAdministrator" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>

</assembly>
`)},
		{name: "win8highest", args: struct{ manifest AppManifest }{AppManifest{
			appManifestCommon: appManifestCommon{
				Name:                              "app.name",
				Version:                           "1.2.3.4",
				Description:                       "Application Description",
				UIAccess:                          true,
				AutoElevate:                       true,
				DisableTheming:                    true,
				DisableWindowFiltering:            true,
				HighResolutionScrollingAware:      true,
				UltraHighResolutionScrollingAware: false,
				LongPathAware:                     false,
				PrinterDriverIsolation:            false,
				GDIScaling:                        false,
				SegmentHeap:                       false,
			},
			ExecutionLevel: HighestAvailable,
			Compatibility:  Win8AndAbove,
			DPIAwareness:   DPIPerMonitor,
		}}, want: []byte(
			// language=xml
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

  <assemblyIdentity type="win32" name="app.name" version="1.2.3.4" processorArchitecture="*"/>
  <description>Application Description</description>

  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
      <supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
    </application>
  </compatibility>

  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true/pm</dpiAware>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">permonitor</dpiAwareness>
      <autoElevate xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</autoElevate>
      <disableTheming xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</disableTheming>
      <disableWindowFiltering xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</disableWindowFiltering>
      <highResolutionScrollingAware xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</highResolutionScrollingAware>
    </windowsSettings>
  </application>

  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="highestAvailable" uiAccess="true"/>
      </requestedPrivileges>
    </security>
  </trustInfo>

</assembly>
`)},
		{name: "win81", args: struct{ manifest AppManifest }{AppManifest{
			appManifestCommon: appManifestCommon{
				Name:                              "app.name",
				Version:                           "1.2.3.4 α",
				Description:                       "Application™\nDescription",
				UIAccess:                          false,
				AutoElevate:                       false,
				DisableTheming:                    false,
				DisableWindowFiltering:            false,
				HighResolutionScrollingAware:      false,
				UltraHighResolutionScrollingAware: true,
				LongPathAware:                     true,
				PrinterDriverIsolation:            true,
				GDIScaling:                        true,
				SegmentHeap:                       true,
			},
			ExecutionLevel: AsInvoker,
			Compatibility:  Win81AndAbove,
			DPIAwareness:   DPIPerMonitorV2,
		}}, want: []byte(
			// language=xml
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

  <assemblyIdentity type="win32" name="app.name" version="1.2.3.4 α" processorArchitecture="*"/>
  <description>Application™
Description</description>

  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{8e0f7a12-bfb3-4fe8-b9a5-48fd50a15a9a}"/>
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
    </application>
  </compatibility>

  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware xmlns="http://schemas.microsoft.com/SMI/2005/WindowsSettings">true</dpiAware>
      <dpiAwareness xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">permonitorv2,system</dpiAwareness>
      <printerDriverIsolation xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</printerDriverIsolation>
      <ultraHighResolutionScrollingAware xmlns="http://schemas.microsoft.com/SMI/2011/WindowsSettings">true</ultraHighResolutionScrollingAware>
      <longPathAware xmlns="http://schemas.microsoft.com/SMI/2016/WindowsSettings">true</longPathAware>
      <gdiScaling xmlns="http://schemas.microsoft.com/SMI/2017/WindowsSettings">true</gdiScaling>
      <heapType xmlns="http://schemas.microsoft.com/SMI/2020/WindowsSettings">SegmentHeap</heapType>
    </windowsSettings>
  </application>

  <trustInfo xmlns="urn:schemas-microsoft-com:asm.v3">
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="asInvoker" uiAccess="false"/>
      </requestedPrivileges>
    </security>
  </trustInfo>

</assembly>
`)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeManifest(tt.args.manifest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("*** makeManifest():\n%v\n*** want:\n%v", string(got), string(tt.want))
			}
		})
	}
}

func TestMakeManifest_Bug(t *testing.T) {
	bak := manifestTemplate
	defer func() {
		recover()
		manifestTemplate = bak
	}()

	manifestTemplate = "{{.bobby}}"
	makeManifest(AppManifest{})

	t.Error("should have panicked")
}
