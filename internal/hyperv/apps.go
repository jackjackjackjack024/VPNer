//go:build windows

package hyperv

// VM apps: VPN programs VPNer installs and sets up inside a new VM by itself,
// so the user does not have to do it in the VM window.

// AppIncy is the INCY VPN client (incy.cc). VPNer installs its Linux
// package, lets it set up its TUN device without a password prompt, starts
// it in TUN mode with every desktop session and passes it the subscription
// the user entered in VPNer.
const AppIncy = "incy"

// The INCY package is pinned by version and checked by SHA-256, like
// everything else VPNer downloads.
const (
	incyVersion = "3.9.0"
	incyURL     = "https://github.com/INCY-DEV/incy-platforms/releases/download/desktop-v" + incyVersion + "/incy-linux-x64.deb"
	// INCY replaced the 3.9.0 package on GitHub on 2026-11-10; this is the
	// SHA-256 of the current one (GitHub's own digest of the asset).
	incySHA256 = "088e7229a33539626cfac10b32b0a77c9a9883c3ce8eb73a9cac1a4baf096be9"
)

// IncySubscriptionFile is the name, in the VM's shared folder, of the file
// through which VPNer hands INCY a subscription link. The session script
// reads it, deletes it and passes the link to INCY.
const IncySubscriptionFile = "vpner-incy-subscription.txt"

// incyInstall runs as root from the VM setup (setupScript). It exits 2 when
// the package could not be downloaded and 3 when it could not be installed;
// the setup reports that and retries later.
//
// The package is 130 MB from GitHub, which can be slow: the download
// resumes where it stopped instead of starting over, the partial file is
// kept between attempts and between setup runs, and only a file that is
// complete but not the pinned one is thrown away.
const incyInstall = `#!/bin/sh
# Installs the INCY VPN client (pinned version, checked by SHA-256).
set -u
[ -x /opt/incy/bin/incy ] && exit 0
DEB=/var/tmp/incy.deb
SUM="` + incySHA256 + `"
good() { [ -f "$DEB" ] && echo "$SUM  $DEB" | sha256sum -c - >/dev/null 2>&1; }
i=0
while [ $i -lt 30 ] && ! good; do
  i=$((i + 1))
  curl -fsSL -C - --retry 5 --retry-delay 10 --connect-timeout 30 -o "$DEB" "` + incyURL + `"
  rc=$?
  # complete but not the pinned file, or the server refused to resume
  if ! good && { [ $rc -eq 0 ] || [ $rc -eq 22 ] || [ $rc -eq 33 ] || [ $rc -eq 36 ]; }; then
    rm -f "$DEB"
  fi
  good || sleep 10
done
good || { echo "vpner: INCY could not be downloaded" >&2; exit 2; }
for i in 1 2 3; do
  DEBIAN_FRONTEND=noninteractive apt-get -o DPkg::Lock::Timeout=1800 install -y "$DEB" pkexec polkitd && break
  sleep 15
done
[ -x /opt/incy/bin/incy ] || { echo "vpner: INCY could not be installed" >&2; exit 3; }
rm -f "$DEB"
`

// incyPolkit lets the VM user run INCY's privileged helper (it creates the
// TUN device and routes) without a password. The helper only accepts its own
// short list of commands.
const incyPolkit = `// VPNer: INCY sets up its tunnel without asking for a password.
polkit.addRule(function(action, subject) {
  if (action.id == "cc.incy.vpn.run-helper" && subject.user == "vpner") {
    return polkit.Result.YES;
  }
});
`

// incySession runs in the desktop session of the VM user.
const incySession = `#!/bin/sh
# VPNer: keeps INCY running in TUN mode and passes it the subscription that
# VPNer puts into the shared folder.
BIN=/opt/incy/bin/incy
SUB=/srv/vpner-share/` + IncySubscriptionFile + `
while [ ! -x "$BIN" ]; do sleep 5; done
CONF="${XDG_CONFIG_HOME:-$HOME/.config}/incy"
mkdir -p "$CONF"
python3 - "$CONF/preferences.json" <<'EOF'
import json, os, sys
p = sys.argv[1]
try:
    d = json.load(open(p))
except Exception:
    d = {}
d.update({"tunnelMode": "TUN", "autoConnect": True, "startMinimized": True, "minimizeToTray": True,
          "killSwitch": False, "autostart": False, "hasLaunchedBefore": True})
tmp = p + ".vpner"
with open(tmp, "w") as f:
    json.dump(d, f, indent=2)
os.replace(tmp, p)
EOF
# INCY names the device after the computer: the name chosen in VPNer.
export HOSTNAME="$(hostname)"
pgrep -u "$(id -u)" -f /opt/incy/ >/dev/null || "$BIN" >/dev/null 2>&1 &
while :; do
  if [ -s "$SUB" ]; then
    url=$(head -n 1 "$SUB" | tr -d '\r\n')
    rm -f "$SUB"
    sleep 10
    "$BIN" "incy://add/$url" >/dev/null 2>&1 &
    sleep 30
    "$BIN" "incy://connect" >/dev/null 2>&1 &
  fi
  sleep 3
done
`

const incyAutostart = `[Desktop Entry]
Type=Application
Name=VPNer: INCY
Exec=/usr/local/bin/vpner-incy
NoDisplay=true
X-GNOME-Autostart-enabled=true
`

// appFiles returns the files an app puts into the VM (path -> mode,
// content), the Debian packages it needs and the script that installs it.
func appFiles(app string) (files []seedFile, pkgs []string, install string) {
	switch app {
	case AppIncy:
		return []seedFile{
			{"/usr/local/sbin/vpner-incy-install", "0755", incyInstall},
			{"/etc/polkit-1/rules.d/49-vpner-incy.rules", "0644", incyPolkit},
			{"/usr/local/bin/vpner-incy", "0755", incySession},
			{"/etc/xdg/autostart/vpner-incy.desktop", "0644", incyAutostart},
		}, []string{"curl", "ca-certificates", "python3"}, "/usr/local/sbin/vpner-incy-install"
	}
	return nil, nil, ""
}

type seedFile struct {
	Path, Mode, Content string
}

// KnownApp reports whether VPNer can set up the app by itself.
func KnownApp(app string) bool { return app == AppIncy }
