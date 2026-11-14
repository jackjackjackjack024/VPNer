# Changelog

**English** · [Русский](CHANGELOG.ru.md)

All notable changes are listed here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), versions follow [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.3.3] - 2026-11-14

### Added

- A virtual machine with INCY gets its own DNS server, "Cloudflare via <connection name>": encrypted DNS (DoH) reached through the VM, so the local network does not see the names of sites that open directly. It is added once INCY actually passes traffic, so a VM that is still setting up or connecting never leaves the computer without DNS. It becomes the default if the default was the system DNS; a default you chose yourself is kept. It is added once: if you delete it or choose the system DNS again on the DNS tab, it stays so.

### Fixed

- An open VPNer window kept showing old settings after they were changed elsewhere (for example, a DNS server added by VPNer itself), and saving the DNS tab could then undo the change. The window now reloads the settings whenever they are saved.

## [0.3.2] - 2026-11-14

### Fixed

- A virtual machine got no desktop and no VPN program on Debian 13, and VPNer did not say so. The setup failed on several things at once: a package that is gone in Debian 13 (policykit-1) made apt install no desktop at all; apt gave up while Debian's own automatic updates held its lock at the first boot; the autologin file went to a folder that does not exist; Xorg could not use the Hyper-V video card ("Cannot run in framebuffer mode"); and INCY's package had been replaced on GitHub, so its checksum no longer matched and the download, which could not resume, started over and over. Now packages missing from Debian are skipped instead of stopping the install, apt waits for its lock, the desktop files are written with their folders, Xorg gets a section for the Hyper-V video card, the INCY checksum is the current one and its download resumes. The setup runs as a service that retries by itself (10 minutes after a failure and whenever the VM starts) until it succeeds, and the VM card says why it failed.
- The setup time on the VM card counts from when the VM was created.

## [0.3.1] - 2026-11-12

### Added

- The card of a virtual machine shows how its first setup goes: how many minutes it has been running, which package it is installing, and that nothing has changed for a long time if it seems stuck. Before, it only said that the setup takes 5–15 minutes. Packages are shown for VMs created with this version; for older ones, the time.

### Fixed

- A virtual machine could not be created on newer Windows 11 builds: Hyper-V stopped with "Virtual hard disk files must be uncompressed and unencrypted and must not be sparse (0xC03A001A)". VPNer wrote the Debian image as a sparse file to save space; now it writes a plain file, and an image left by an earlier version is turned into one by itself the next time a VM is created, without downloading it again.
- Clicking "Connect" several times while VPNer was still starting (for example, while the routing engine was downloading) created several virtual machines. Now a connection starts only once at a time, and its card shows at once that it is waiting for routing to turn on.
- "VM window" seemed to do nothing: the window opened behind the other windows, and every click opened one more, which then asked to disconnect "another user". Now an open window of the VM is brought to the front, and a new one comes to the front once it appears.

## [0.3.0] - 2026-11-12

### Added

- A virtual machine can set up its VPN program by itself: choose INCY when adding the VM, and VPNer installs it (a pinned version checked by SHA-256), lets it create its tunnel without a password prompt, starts it with the VM in tunnel mode and connects it. The VM is named with the device name you enter, so it shows up under that name in your subscription's list of devices. A subscription link entered in VPNer is passed to the program once over VPNer's own channel to the VM and is never written to the VM's setup disk.
- Device-limit warnings. Many VPN services limit how many devices or simultaneous connections one subscription, account or key may use, so VPNer now says what can take such a place: before an action in its dialog, and right after it in a notice that stays until you close it. This covers adding and importing connections (also a subscription or a profile), exporting a profile, saving a chain and turning a connection on. A VM with a VPN program is a new device of its subscription; a corporate VPN login is one more session; a WireGuard/AmneziaWG key is one device and cannot work in two places at once; an Amnezia server gets one more user; a subscription or share-link account may count this computer, or a chain's copy that comes from another address, as one more device.

## [0.2.0] - 2026-11-12

### Added

- Chains of connections: a new connection type made of your other connections in order. Traffic enters the first one, goes through it into the next and leaves through the last, whose address the sites see — for example, reach the office VPN only through the home one. The connections in a chain keep working on their own, so rules can use each of them and the chain side by side. The first connection can be of any type (also a WSL container or a virtual machine); after it come built-in protocols. Connections for a chain can be picked from the existing ones or created right from the chain editor. If any part of a chain is down, its traffic is blocked, never sent around it. An Amnezia "full access" connection after the first hop gets its own user on your server for the chain, created when the chain is saved and removed with it, so the connection keeps working on its own at the same time.
- Importing a Clash/mihomo config turns proxies with `dialer-proxy` into chains, and their rules go to the chain. Before, such a proxy broke the engine config.

## [0.1.3] - 2026-11-02

### Fixed

- The network no longer stays blocked after Windows restarts with routing and the kill switch on. Before, VPNer was ended during shutdown before it could lift the kill switch, which lives in Windows Firewall and survives a reboot, so the PC had no network at all until VPNer was started again. Now VPNer lifts it when Windows shuts down, restarts or signs out. If that does not happen (power loss, a crash), a scheduled task restores the original firewall policy at the next startup or logon and removes itself. With autostart and routing on launch both on, the kill switch is kept until VPNer starts, as before.

## [0.1.2] - 2026-10-25

### Fixed

- Rule changes now reach connections that are already open: after saving, VPNer closes the open connections whose route is different under the new rules, and the programs open them again on the new route. Before, a browser reusing a connection kept the old route until it was closed by hand. Connections whose route did not change, and VPNer's own, stay open.
- Connections to a FakeIP address left in a program's DNS cache after an engine restart are dropped at once, so the program asks DNS again, instead of hanging until a timeout.
- Domain rules now also work for programs that resolve names on their own (DNS over HTTPS in browsers, the browser inside apps) and connect to a bare IP: the site name is taken from the connection itself (TLS/QUIC SNI, HTTP Host). Before, such connections fell through to program rules or went directly.

### Changed

- Profile export opens a window with the profile as text to copy, or to save to a file; with passwords and keys included, the window warns that the text is a secret.

## [0.1.1] - 2026-10-21

### Added

- Legal information in English and Russian (LEGAL.md): the open-source components inside VPNer and what their licenses require, the components VPNer downloads, trademarks, no warranty, limitation of liability, your responsibility and privacy with every network request VPNer makes.
- THIRD_PARTY_NOTICES.txt with the license texts of the built-in components: attached to releases, installed next to vpner.exe, included in the portable zip and built into vpner.exe (Settings → Legal, `vpner.exe licenses`).
- A short "as is" note in the installer.

## [0.1.0] - 2026-10-20

The first public beta.

### Added

- Routing of programs, sites, domains, subnets, site categories, countries, ports and protocols through several VPNs at once (the mihomo engine behind a TUN adapter), with a fail-closed default: rules leading to a VPN that is down block traffic.
- Built-in WireGuard and AmneziaWG (including 2.0), the Amnezia full access key (VPNer creates its own user on your server), VLESS/Reality, VMess, Trojan, Shadowsocks, Hysteria2, TUIC, SOCKS5/HTTP and network interfaces.
- Cisco AnyConnect / Secure Client, GlobalProtect, Fortinet, Pulse/Ivanti, F5 (OpenConnect) and OpenVPN in isolated WSL network namespaces.
- Hyper-V virtual machines for original VPN apps, created and set up automatically.
- Saved passwords (DPAPI), MFA codes, push confirmations, SAML/SSO sign-in; corporate VPNs log in one at a time; idle sessions are kept alive; logins that were not confirmed in time are retried.
- DNS: DoH, DoH3, DoT, DoQ, UDP/TCP, the system resolver, servers reached through a VPN, per-domain DNS rules, FakeIP.
- Import of keys, configs, links, subscriptions, Clash/mihomo configs and VPNer profiles; profile export.
- Leak test, live connections, log, system kill switch (Windows Firewall), emergency network restore.
- Self-healing of stalled tunnels; engine processes end together with VPNer.
- Appearance: Standard, CLI / terminal and Glass themes in light and dark, own values and variables, per-connection colors, animations and effects that can be turned off one by one.
- English and Russian interface; texts in locale files.
- Installer with language and theme choice, uninstall, portable zip.
