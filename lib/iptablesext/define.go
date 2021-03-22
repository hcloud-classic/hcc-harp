package iptablesext

// NeededTablesForHarp : iptables needed tables for harp
var NeededTablesForHarp = []string{"filter", "nat"}

// NeededKernelModulesForHarp : Needed kernel modules for using iptables in harp
var NeededKernelModulesForHarp = "iptable_filter, iptable_nat, nf_nat"

// NatChains : chain names used in NAT
var NatChains = []string{"POSTROUTING", "PREROUTING"}

// HarpChainNamePrefix : iptables chain name prefix for harp
var HarpChainNamePrefix = "HARP_"