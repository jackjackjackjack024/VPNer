package agent

import (
	"strings"
	"time"
)

// dpkgStep returns the package dpkg worked on last, from the tail of
// /var/log/dpkg.log, so the host can show how the VM setup goes. Lines older
// than since are left out: the Debian image comes with entries of its own
// build. dpkg.log lines look like:
//
//	2026-11-12 21:40:01 status half-configured xfce4-panel:amd64 4.18.4-1
//	2026-11-12 21:40:01 install incy:amd64 <none> 3.9.0
func dpkgStep(tail string, since time.Time) string {
	lines := strings.Split(strings.TrimSpace(tail), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		f := strings.Fields(lines[i])
		if len(f) < 4 {
			continue
		}
		at, err := time.ParseInLocation("2006-01-02 15:04:05", f[0]+" "+f[1], time.Local)
		if err != nil {
			continue
		}
		if at.Before(since) {
			return ""
		}
		pkg := ""
		switch f[2] {
		case "status":
			if len(f) >= 5 {
				pkg = f[4]
			}
		case "install", "upgrade", "configure", "trigproc":
			pkg = f[3]
		}
		if pkg != "" {
			pkg, _, _ = strings.Cut(pkg, ":")
			return pkg
		}
	}
	return ""
}
