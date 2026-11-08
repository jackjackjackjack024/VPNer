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
	incySHA256  = "1dc8f1c110068424ec657ab11cac3745ec61c8e1f6154918f9451250db26db8a"
)

// IncySubscriptionFile is the name, in the VM's shared folder, of the file
// through which VPNer hands INCY a subscription link. The session script
// reads it, deletes it and passes the link to INCY.
const IncySubscriptionFile = "vpner-incy-subscription.txt"

// incyInstall runs as root once, at the first boot.
const incyInstall = `#!/bin/sh
# Installs the INCY VPN client (pinned version, checked by SHA-256).
set -u
DEB=/var/tmp/incy.deb
for i in 1 2 3 4 5; do
  if curl -fsSL --retry 3 -o "$DEB" "` + incyURL + `" &&
     echo "` + incySHA256 + `  $DEB" | sha256sum -c - >/dev/null; then
    break
  fi
  rm -f "$DEB"
  sleep 15
done
[ -f "$DEB" ] || { echo "vpner: INCY could not be downloaded" >&2; exit 1; }
for i in 1 2 3; do
  DEBIAN_FRONTEND=noninteractive apt-get install -y "$DEB" pkexec polkitd && break
  sleep 10
done
rm -f "$DEB"
[ -x /opt/incy/bin/incy ]
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
// content) and the commands to run at the first boot.
func appFiles(app string) (files []seedFile, cmds []string) {
	switch app {
	case AppIncy:
		return []seedFile{
			{"/usr/local/sbin/vpner-incy-install", "0755", incyInstall},
			{"/etc/polkit-1/rules.d/49-vpner-incy.rules", "0644", incyPolkit},
			{"/usr/local/bin/vpner-incy", "0755", incySession},
			{"/etc/xdg/autostart/vpner-incy.desktop", "0644", incyAutostart},
		}, []string{
			`DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends curl ca-certificates python3`,
			`/usr/local/sbin/vpner-incy-install`,
		}
	}
	return nil, nil
}

type seedFile struct {
	Path, Mode, Content string
}

// KnownApp reports whether VPNer can set up the app by itself.
func KnownApp(app string) bool { return app == AppIncy }
