package nmcli

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"strconv"
)

const NMCLIBIN string = "nmcli"

func run(cmd *exec.Cmd) error {
	var stderr bytes.Buffer

	cmd.Stderr = &stderr

	if err := cmd.Run() ; err != nil {
		return fmt.Errorf(stderr.String())
	}

	return nil
}

func DeviceShow(name string) (*Device, error) {
	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, "device", "show", name)

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^GENERAL\.STATE:\s+([0-9\.]+)\s\(.+\)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 {
			device := Device{
				Name:  strings.TrimSpace(name),
				State:  strings.TrimSpace(match[1]),
			}

			return &device, nil
		}
	}

	return nil, nil
}

func DeviceStatus(dtype string) (*Device, error) {
	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, "device", "status")

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^(.*)\s+([a-z\-]+)\s+(connected|disconnected|unmanaged)\s+(.*)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 && match[2] == dtype {
			device := Device{
				Name: strings.TrimSpace(match[1]),
				Type: strings.TrimSpace(match[2]),
				State:  strings.TrimSpace(match[3]),
				Connection: strings.TrimSpace(match[4]),
			}

			return &device, nil
		}
	}

	return nil, nil
}

func ConnectionShow(ctype, name string) (*Connection, error) {
	var connection Connection

	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, "connection", "show", "--active")

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^(.*)\s+([a-z0-9]{8}\-[a-z0-9]{4}\-[a-z0-9]{4}\-[a-z0-9]{4}\-[a-z0-9]{12})\s+([a-z\-]+)\s+(.*)$`)

	active := false
	for _, l := range result {
		match := re.FindStringSubmatch(l)

        if len(match) > 0 && strings.TrimSpace(match[3]) == ctype {
			if len(name) > 0 && strings.TrimSpace(match[1]) != name {
				continue
			}

			connection = Connection{
				Name: strings.TrimSpace(match[1]),
				Uuid: strings.TrimSpace(match[2]),
				Type: strings.TrimSpace(match[3]),
				Device: strings.TrimSpace(match[4]),
			}
			active = true
		}
	}

	if !active {
		return nil, fmt.Errorf("no active connection")
	}

	return &connection, nil
}

func ConnectionAdd(name, ctype, ifname string, autoconnect bool) error {
	ac := "no"
	if autoconnect {
		ac = "yes"
	}

	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return err
	}
	cmd := exec.Command(execPath,
		"connection", "add",
		"type", ctype,
		"ifname", ifname,
		"con-name", name,
		"autoconnect", ac,
		"ssid", name,
	)

	if err := run(cmd); err != nil {
		return err
	}


	return nil
}

func ConnectionDelete(name string) error {
	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return err
	}
	cmd := exec.Command(execPath, "connection", "delete", name)

	if err := run(cmd); err != nil {
		return err
	}

	return nil
}

func ConnectionModify(name, option, value string) error {
	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return err
	}
	cmd := exec.Command(execPath, "connection", "modify", name, option, value)

	if err := run(cmd); err != nil {
		return err
	}

	return nil
}

func ConnectionUp(name string) error {
	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return err
	}
	cmd := exec.Command(execPath, "connection", "up", name)

	if err := run(cmd); err != nil {
		return err
	}

	return nil
}

func ConnectionDown(name string) error {
	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return err
	}
	cmd := exec.Command(execPath, "connection", "down", name)

	if err := run(cmd); err != nil {
		return err
	}

	return nil
}

func GetConnectionDhcpDns(uuid string) ([]string, error) {
	var dns []string

	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, "c", "show", uuid)

	o, err := cmd.Output()
	if err != nil {
		return dns, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^IP4\.DNS\[1\]:\s+([0-9\.]+)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 {
			dns = append(dns, match[1])
		}
	}

	return dns, nil
}

func WifiList() ([]Wifi, error) {
	var wifis []Wifi

	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, "device", "wifi", "list")

	o, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	result := strings.Split(string(o), "\n")

	re := regexp.MustCompile(`^(\*?)\s+([A-Z0-9\:]{17})\s+(.+)\s+([A-Z][a-z]+)\s+([0-9]{1,3})\s+([0-9]+\s[a-zA-Z\/]+)\s+([0-9]{1,2})\s+(\**)\s+([A-Z0-9\.\s]+)$`)

	for _, l := range result {
		match := re.FindStringSubmatch(l)

		if len(match) > 0 {
			inUse := false
			if match[1] == "*" {
				inUse = true
			}

			channel, _ := strconv.Atoi(match[5])
			signal, _ := strconv.Atoi(match[7])

			wifi := Wifi{
				InUse: inUse,
				Bssid: strings.TrimSpace(match[2]),
				Ssid: strings.TrimSpace(match[3]),
				Mode: strings.TrimSpace(match[4]),
				Chan: channel,
				Rate: strings.TrimSpace(match[6]),
				Signal: signal,
				Bars: strings.TrimSpace(match[8]),
				Security: strings.TrimSpace(match[9]),
			}

			wifis = append(wifis, wifi)
		}
	}

	return wifis, nil
}

func WifiConnect(ssid, password string) error {
	execPath, err := exec.LookPath(NMCLIBIN)
	if err != nil {
		return err
	}
	cmd := exec.Command(execPath,
        "device", "wifi",
        "connect", ssid,
        "password", password,
    )

	if err := run(cmd); err != nil {
		return err
	}

	return nil
}
