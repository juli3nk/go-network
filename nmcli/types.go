package nmcli

type Device struct {
	Name       string
	Type       string
	State      string
	Connection string
}

type Connection struct {
	Name   string
	Uuid   string
	Type   string
	Device string
}

type Wifi struct {
	InUse    bool   `json:"in_use"`
	Bssid    string `json:"bssid"`
	Ssid     string `json:"ssid"`
	Mode     string `json:"mode"`
	Chan     int    `json:"chan"`
	Rate     string `json:"rate"`
	Signal   int    `json:"signal"`
	Bars     string `json:"bars"`
	Security string `json:"security"`
}
