package importer

import "testing"

// "XRay native format" from the AmneziaVPN app pasted into the window.
func TestParseXRayJSON(t *testing.T) {
	js := `{"outbounds":[{"protocol":"vless","settings":{"vnext":[{"address":"203.0.113.7","port":443,
	  "users":[{"id":"11111111-2222-3333-4444-555555555555","encryption":"none","flow":"xtls-rprx-vision"}]}]},
	  "streamSettings":{"network":"tcp","security":"reality","realitySettings":{"serverName":"www.example.com","publicKey":"PUBKEY","shortId":"abcd","fingerprint":"chrome"}}}]}`
	r, err := Parse("amnezia_for_xray.json", js)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Connections) != 1 || r.Connections[0].Type != "link" {
		t.Fatalf("%+v", r)
	}
}
