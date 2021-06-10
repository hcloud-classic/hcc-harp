package iptablesext

// NeededTablesForHarp : iptables needed tables for harp
var NeededTablesForHarp = []string{"filter", "nat"}

// NatChains : chain names used in NAT
var NatChains = []string{"POSTROUTING", "PREROUTING"}

// HarpChainNamePrefix : iptables chain name prefix for harp
var HarpChainNamePrefix = "HARP_"

// HarpAdaptiveIPInputDropChainName : iptables chain name for dropping inputs destined to AdaptiveIP
var HarpAdaptiveIPInputDropChainName = HarpChainNamePrefix + "ADAPTIVE_IP_INPUT_DROP"
