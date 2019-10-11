package dhcpd

var confBase = "group {\n" +
	"    use-host-decl-names on;\n" +
	"    filename \"HARP_DHCPD_PXE_FILENAME\";\n" +
	"    option domain-name \"HARP_DHCPD_DOMAIN_NAME\";\n" +
	"    min-lease-time HARP_DHCPD_MIN_LEASE_TIME;\n" +
	"    default-lease-time HARP_DHCPD_DEFAULT_LEASE_TIME;\n" +
	"    max-lease-time HARP_DHCPD_MAX_LEASE_TIME;\n" +
	"\n" +
	"    # node entries start here\n" +
	"HARP_DHCPD_NODES_ENTRIES\n" +
	"}"

var nodeEntry = "    host HARP_DHCPD_NODE_NAME { hardware ethernet HARP_DHCPD_NODE_PXE_MAC ; fixed-address HARP_DHCPD_NODE_IP ; }\n"
