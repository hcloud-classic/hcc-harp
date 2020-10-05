package iptables

import "os/exec"

// EnableAllRouteLocal : Enable all route of local to use NAT with iptables.
func EnableAllRouteLocal() error {
	cmd := exec.Command("sysctl", "-w", "net.ipv4.conf.all.route_localnet=1")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// EnableIPForwardV4 : Enable IPv4 IP forward to use NAT with iptables.
func EnableIPForwardV4() error {
	cmd := exec.Command("echo", "1", ">", "/proc/sys/net/ipv4/ip_forward")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}