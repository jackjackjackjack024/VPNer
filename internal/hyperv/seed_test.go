//go:build windows

package hyperv

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

type userData struct {
	WriteFiles []struct {
		Path        string `yaml:"path"`
		Permissions string `yaml:"permissions"`
		Content     string `yaml:"content"`
	} `yaml:"write_files"`
	Runcmd [][]string `yaml:"runcmd"`
}

func readSeed(t *testing.T, app string) userData {
	dir := t.TempDir()
	err := WriteSeed(dir, SeedConfig{ID: "t", Hostname: "vpner-t", LANIP: "10.77.0.10", WANMac: "00:15:5d:00:00:01",
		LANMac: "00:15:5d:00:00:02", GatewayUser: "vpner", GatewayPass: "p", ConsolePass: "c", Desktop: true,
		Agent: []byte("agent"), App: app})
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "user-data"))
	if err != nil {
		t.Fatal(err)
	}
	var u userData
	if err := yaml.Unmarshal(b, &u); err != nil {
		t.Fatalf("user-data is not valid YAML: %v\n%s", err, b)
	}
	return u
}

// INCY's files reach the VM unchanged and it is installed after the desktop.
func TestSeedIncy(t *testing.T) {
	u := readSeed(t, AppIncy)
	want, _ := appFiles(AppIncy)
	for _, w := range want {
		found := false
		for _, f := range u.WriteFiles {
			if f.Path == w.Path {
				found = true
				if f.Content != w.Content {
					t.Errorf("%s changed on the way:\n%q\nwant\n%q", w.Path, f.Content, w.Content)
				}
				if f.Permissions != w.Mode {
					t.Errorf("%s permissions %s, want %s", w.Path, f.Permissions, w.Mode)
				}
			}
		}
		if !found {
			t.Errorf("%s is missing", w.Path)
		}
	}
	desktop, install := -1, -1
	for i, c := range u.Runcmd {
		line := strings.Join(c, " ")
		if strings.Contains(line, "lightdm restart") || strings.Contains(line, "restart lightdm") {
			desktop = i
		}
		if strings.Contains(line, "/usr/local/sbin/vpner-incy-install") {
			install = i
		}
	}
	if install < 0 || install < desktop {
		t.Errorf("INCY must be installed after the desktop: desktop %d, install %d", desktop, install)
	}
	if !strings.Contains(incyInstall, incySHA256) || !strings.Contains(incyInstall, "sha256sum -c") {
		t.Error("the package must be checked by SHA-256")
	}
}

// Without an app the seed is as before.
func TestSeedNoApp(t *testing.T) {
	u := readSeed(t, "")
	for _, f := range u.WriteFiles {
		if strings.Contains(f.Path, "incy") {
			t.Errorf("unexpected %s", f.Path)
		}
	}
	for _, c := range u.Runcmd {
		if strings.Contains(strings.Join(c, " "), "incy") {
			t.Errorf("unexpected command %v", c)
		}
	}
}

// The subscription link is never written to the setup disk.
func TestSeedHasNoSubscription(t *testing.T) {
	u := readSeed(t, AppIncy)
	for _, f := range u.WriteFiles {
		if strings.Contains(f.Content, "https://") && !strings.Contains(f.Content, incyURL) {
			t.Errorf("%s holds a link", f.Path)
		}
	}
}

// The scripts parse (sh -n only reads them); skipped without a shell.
func TestIncyScriptsParse(t *testing.T) {
	sh, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("no sh")
	}
	files, _ := appFiles(AppIncy)
	for _, f := range files {
		if !strings.HasPrefix(f.Content, "#!/bin/sh") {
			continue
		}
		p := filepath.Join(t.TempDir(), "s.sh")
		if err := os.WriteFile(p, []byte(f.Content), 0o644); err != nil {
			t.Fatal(err)
		}
		if out, err := exec.Command(sh, "-n", p).CombinedOutput(); err != nil {
			t.Errorf("%s: %v\n%s", f.Path, err, out)
		}
	}
}
