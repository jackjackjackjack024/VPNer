# Changelog

All notable changes are listed here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), versions follow [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Chains of connections: a new connection type made of your other connections in order. Traffic enters the first one, goes through it into the next and leaves through the last, whose address the sites see — for example, reach the office VPN only through the home one. The connections in a chain keep working on their own, so rules can use each of them and the chain side by side. The first connection can be of any type (also a WSL container or a virtual machine); after it come built-in protocols. Connections for a chain can be picked from the existing ones or created right from the chain editor. If any part of a chain is down, its traffic is blocked, never sent around it. An Amnezia "full access" connection after the first hop gets its own user on your server for the chain, created when the chain is saved and removed with it, so the connection keeps working on its own at the same time.
- A virtual machine can set up its VPN program by itself: choose INCY when adding the VM, and VPNer installs it (a pinned version checked by SHA-256), lets it create its tunnel without a password prompt, starts it with the VM in tunnel mode and connects it. A subscription link entered in VPNer is passed to the program once over VPNer's own channel to the VM and is never written to the VM's setup disk.
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
