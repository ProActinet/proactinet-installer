package main

import (
	"fmt"
	"net"
	"runtime"
	"strings"
)

type InterfaceType int

const (
	WLAN InterfaceType = iota
	ETHERNET
)

// getInterfacePatterns returns patterns for both WLAN and Ethernet interfaces
func getInterfacePatterns() map[string][]string {
	patterns := make(map[string][]string)

	switch runtime.GOOS {
	case "linux":
		patterns["wlan"] = []string{"wlan", "wifi", "wlp", "wl"}
		patterns["ethernet"] = []string{"eth", "enp", "eno", "ens"}
	case "darwin":
		patterns["wlan"] = []string{"en"}
		patterns["ethernet"] = []string{"en"} // macOS uses en* for both
	case "windows":
		patterns["wlan"] = []string{"Wi-Fi", "Wireless"}
		patterns["ethernet"] = []string{"Ethernet", "Local Area Connection"}
	default:
		patterns["wlan"] = []string{"wlan", "wifi"}
		patterns["ethernet"] = []string{"eth", "en"}
	}

	return patterns
}

func isMatchingInterface(name string, patterns []string) bool {
	name = strings.ToLower(name)
	for _, pattern := range patterns {
		if strings.Contains(name, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func getNetworkInterfaces(interfaceType string) ([]net.Interface, error) {
	// Get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, fmt.Errorf("failed to get network interfaces: %v", err)
	}

	patterns := getInterfacePatterns()
	var matchedInterfaces []net.Interface

	for _, iface := range interfaces {
		// Special case for macOS
		if runtime.GOOS == "darwin" {
			if interfaceType == "wlan" && iface.Name == "en0" {
				matchedInterfaces = append(matchedInterfaces, iface)
				continue
			}
			if interfaceType == "ethernet" && strings.HasPrefix(iface.Name, "en") && iface.Name != "en0" {
				matchedInterfaces = append(matchedInterfaces, iface)
				continue
			}
		}

		// For other OS, check against patterns
		if isMatchingInterface(iface.Name, patterns[interfaceType]) {
			matchedInterfaces = append(matchedInterfaces, iface)
		}
	}

	return matchedInterfaces, nil
}

func getInterfaceDetails(iface net.Interface) map[string]interface{} {
	details := make(map[string]interface{})

	details["name"] = iface.Name
	details["index"] = iface.Index
	details["mtu"] = iface.MTU
	details["hardware_addr"] = iface.HardwareAddr.String()
	details["flags"] = iface.Flags.String()

	addrs, err := iface.Addrs()
	if err == nil {
		addrList := make([]string, 0)
		for _, addr := range addrs {
			addrList = append(addrList, addr.String())
		}
		details["addresses"] = addrList
	}

	return details
}

func gtInterfaceDetails(interfaceType string) interface{} {
	interfaces, err := getNetworkInterfaces(interfaceType)
	if err != nil {
		fmt.Printf("Error getting %s interfaces: %v\n", interfaceType, err)
		return nil
	}

	if len(interfaces) == 0 {
		fmt.Printf("No %s interfaces found\n", interfaceType)
		return nil
	}

	fmt.Printf("\nFound %s interfaces for %s:\n", interfaceType, runtime.GOOS)
	for _, iface := range interfaces {
		details := getInterfaceDetails(iface)
		return details["name"]
	}

	return nil
}
