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
		want string
	}{
		{
			name: "empty",
			args: struct{ manifest AppManifest }{},
			// language=manifest
			want: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

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
`},
		{
			name: "full",
			args: struct{ manifest AppManifest }{AppManifest{
				Identity: AssemblyIdentity{
					Name:    "<app.name>",
					Version: [4]uint16{1, 2, 3, 4},
				},
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
				ExecutionLevel:                    HighestAvailable,
				Compatibility:                     WinVistaAndAbove,
				DPIAwareness:                      DPIAware,
			}},
			// language=manifest
			want: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

  <assemblyIdentity type="win32" name="&lt;app.name&gt;" version="1.2.3.4" processorArchitecture="*"/>
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
`},
		{
			name: "win10admin",
			args: struct{ manifest AppManifest }{AppManifest{
				Identity: AssemblyIdentity{
					// No name, no identity (empty name is forbidden)
					Version: [4]uint16{1, 2, 3, 4},
				},
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
				ExecutionLevel:                    RequireAdministrator,
				Compatibility:                     Win10AndAbove,
				DPIAwareness:                      DPIUnaware,
			}},
			// language=manifest
			want: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">
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
`},
		{
			name: "win8highest",
			args: struct{ manifest AppManifest }{AppManifest{
				Identity: AssemblyIdentity{
					Name: "app.name",
					// No version -> 0.0.0.0
				},
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
				ExecutionLevel:                    HighestAvailable,
				Compatibility:                     Win8AndAbove,
				DPIAwareness:                      DPIPerMonitor,
			}},
			// language=manifest
			want: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

  <assemblyIdentity type="win32" name="app.name" version="0.0.0.0" processorArchitecture="*"/>
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
`},
		{
			name: "win81",
			args: struct{ manifest AppManifest }{AppManifest{
				Identity: AssemblyIdentity{
					Name:    "app.name",
					Version: [4]uint16{0xFFFF, 65535, 0xFFFF, 65535},
				},
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
				ExecutionLevel:                    AsInvoker,
				Compatibility:                     Win81AndAbove,
				DPIAwareness:                      DPIPerMonitorV2,
			}},
			// language=manifest
			want: `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly xmlns="urn:schemas-microsoft-com:asm.v1" manifestVersion="1.0">

  <assemblyIdentity type="win32" name="app.name" version="65535.65535.65535.65535" processorArchitecture="*"/>
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
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeManifest(tt.args.manifest); string(got) != tt.want {
				t.Errorf("*** makeManifest():\n%v###\n*** want:\n%v###", string(got), tt.want)
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

func TestAppManifestFromXML(t *testing.T) {
	tests := []struct {
		name    string
		xml     string
		want    AppManifest
		wantErr bool
	}{
		{
			name: "zero", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly></assembly>`,
			want:    AppManifest{DPIAwareness: DPIUnaware},
			wantErr: false,
		},
		{
			name: "longVersion", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly>
  <assemblyIdentity type="win32" name=" a.name" version="-1.2.3.4.5.6.7" processorArchitecture="x86"/>
</assembly>`,
			want: AppManifest{
				Identity: AssemblyIdentity{
					Name:    " a.name",
					Version: [4]uint16{0, 2, 3, 4},
				},
				DPIAwareness: DPIUnaware,
			},
			wantErr: false,
		},
		{
			name: "shortVersion", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly>
  <assemblyIdentity type="win32" name="💨" version="42.5💨" processorArchitecture="x86"/>
</assembly>`,
			want: AppManifest{
				Identity: AssemblyIdentity{
					Name:    "💨",
					Version: [4]uint16{42},
				},
				DPIAwareness: DPIUnaware,
			},
			wantErr: false,
		},
		{
			name: "dpiAware", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly>
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAware> TrUe </dpiAware>
    </windowsSettings>
  </application>
</assembly>`,
			want:    AppManifest{},
			wantErr: false,
		},
		{
			name: "manifest1", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly>
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAwareness> something,system, </dpiAwareness>

      <autoElevate>true</autoElevate>
      <disableTheming> tRue</disableTheming>
      <disableWindowFiltering>true</disableWindowFiltering>
      
      <printerDriverIsolation>false</printerDriverIsolation>
      <ultraHighResolutionScrollingAware> err</ultraHighResolutionScrollingAware>
      <longPathAware> yes</longPathAware>
      <gdiScaling>true</gdiScaling>
      <heapType attr="x">segmentHeap</heapType>
    </windowsSettings>
    <windowsSettings>
      <highResolutionScrollingAware>true</highResolutionScrollingAware>
    </windowsSettings>
  <dependency>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="5.6.6.6" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
    </dependentAssembly>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="6.0.0.0" processorArchitecture="*" publicKeyToken="06595b64144ccf1df" language="*"/>
    </dependentAssembly>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls6" version="6.0.0.0" processorArchitecture="*" publicKeyToken="6595b64144ccf1df" language="*"/>
    </dependentAssembly>
  </dependency>
  </application>
  <description>This is a 
  description</description>
</assembly>`,
			want: AppManifest{
				Description:                  "This is a \n  description",
				AutoElevate:                  true,
				DisableTheming:               true,
				DisableWindowFiltering:       true,
				GDIScaling:                   true,
				HighResolutionScrollingAware: true,
				SegmentHeap:                  true,
			},
			wantErr: false,
		},
		{
			name: "manifest2", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly>
  <trustInfo>
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="require administrator" uiAccess=" True "/>
      </requestedPrivileges>
    </security>
  </trustInfo>
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAwareness>system</dpiAwareness>
      <ultraHighResolutionScrollingAware>true</ultraHighResolutionScrollingAware>
      <printerDriverIsolation>true</printerDriverIsolation>
      <longPathAware >false</longPathAware>
      <longPathAware >true</longPathAware>
      <gdiScaling>true</gdiScaling>
    </windowsSettings>

  </application>
  <dependency>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="a" version="5.6.6.6" processorArchitecture="*" publicKeyToken="42" language="*"/>
    </dependentAssembly>
    <dependentAssembly>
      <assemblyIdentity type="win32" name="Microsoft.Windows.Common-Controls" version="  6.0.0.0 " processorArchitecture="*" publicKeyToken=" 6595B64144CCF1DF " language="*"/>
    </dependentAssembly>
  </dependency>
</assembly>`,
			want: AppManifest{
				UIAccess:                          true,
				PrinterDriverIsolation:            true,
				UltraHighResolutionScrollingAware: true,
				LongPathAware:                     true,
				GDIScaling:                        true,
				UseCommonControlsV6:               true,
			},
			wantErr: false,
		},
		{
			name: "admin", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly>
  <trustInfo>
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="RequireAdministrator" uiAccess="yes"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
  <application xmlns="urn:schemas-microsoft-com:asm.v3">
    <windowsSettings>
      <dpiAwareness>system</dpiAwareness>
    </windowsSettings>
  </application>
</assembly>`,
			want:    AppManifest{ExecutionLevel: RequireAdministrator},
			wantErr: false,
		},
		{
			name: "highest", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly>
  <trustInfo>
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="HIGHESTAVAILABLE" uiAccess="no thanks"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="etzetezteztzertz"/>
      <supportedOS Id=" {8E0F7A12-BFB3-4FE8-B9A5-48FD50A15A9A} "/>
    </application>
  </compatibility>
  <application>
    <windowsSettings>
      <dpiAware>true</dpiAware>
    </windowsSettings>
  </application>
</assembly>`,
			want: AppManifest{
				ExecutionLevel: HighestAvailable,
				Compatibility:  Win10AndAbove,
			},
			wantErr: false,
		},
		{
			name: "parseError", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?><assembly></assemble>`,
			want:    AppManifest{},
			wantErr: true,
		},
		{
			name: "os", xml: // language=manifest
			`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<assembly>
  <trustInfo>
    <security>
      <requestedPrivileges>
        <requestedExecutionLevel level="*" uiAccess="no thanks"/>
      </requestedPrivileges>
    </security>
  </trustInfo>
  <compatibility xmlns="urn:schemas-microsoft-com:compatibility.v1">
    <application>
      <supportedOS Id="{1f676c76-80e1-4239-95bb-83d0f6d0da78}"/>
      <supportedOS Id=" {8E0F7A12-BFB3-4FE8-B9A5-48FD50A15A9A} "/>
      <supportedOS Id="{e2011457-1546-43c5-a5fe-008deee3d3f0}"/>
      <supportedOS Id="{4a2f28e3-53b9-4441-ba9c-d69d4a4a6e38}"/>
      <supportedOS Id="{35138b9a-5d96-4fbd-8e2d-a2440225f93a}"/>
    </application>
  </compatibility>
  <application>
    <windowsSettings>
      <dpiAware>true</dpiAware>
    </windowsSettings>
  </application>
</assembly>`,
			want:    AppManifest{Compatibility: WinVistaAndAbove},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AppManifestFromXML([]byte(tt.xml))
			if (err != nil) != tt.wantErr {
				t.Errorf("AppManifestFromXML() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppManifestFromXML() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readDPIAwareness(t *testing.T) {
	type args struct {
		dpiAware     string
		dpiAwareness string
	}
	tests := []struct {
		name string
		args args
		want DPIAwareness
	}{
		{name: "", args: args{dpiAware: "", dpiAwareness: "true, systeM , permonitor"}, want: DPIAware},
		{name: "", args: args{dpiAware: "true", dpiAwareness: "perMonitorv2,system"}, want: DPIPerMonitorV2},
		{name: "", args: args{dpiAware: "true", dpiAwareness: ""}, want: DPIAware},
		{name: "", args: args{dpiAware: " true/PM ", dpiAwareness: ""}, want: DPIPerMonitor},
		{name: "", args: args{dpiAware: "false", dpiAwareness: "system"}, want: DPIAware},
		{name: "", args: args{dpiAware: "true / PM", dpiAwareness: "per monitor"}, want: DPIUnaware},
		{name: "", args: args{dpiAware: "true", dpiAwareness: "perMonitorv2,system"}, want: DPIPerMonitorV2},
		{name: "", args: args{dpiAware: "true", dpiAwareness: "PerMonitor"}, want: DPIPerMonitor},
		{name: "", args: args{dpiAware: "true", dpiAwareness: " unaWarE"}, want: DPIUnaware},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readDPIAwareness(tt.args.dpiAware, tt.args.dpiAwareness); got != tt.want {
				t.Errorf("readDPIAwareness() = %v, want %v", got, tt.want)
			}
		})
	}
}
