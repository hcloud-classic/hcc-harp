package daoext

// SubnetValid : Valid subnet
var SubnetValid int64 = 0

// SubnetValidErrorArgumentError : Some of arguments are missing
var SubnetValidErrorArgumentError int64 = 1

// SubnetValidErrorInvalidNetworkAddress : Invalid network address
var SubnetValidErrorInvalidNetworkAddress int64 = 2

// SubnetValidErrorInvalidNetmask : Invalid netmask
var SubnetValidErrorInvalidNetmask int64 = 3

// SubnetValidErrorSubnetConflict : Subnet is conflicted with one of stored in the database
var SubnetValidErrorSubnetConflict int64 = 4

// SubnetValidErrorNotPrivate : Network address is not private address
var SubnetValidErrorNotPrivate int64 = 5

// SubnetValidErrorStartIPNot1 : Start IP address is not x.x.x.1
var SubnetValidErrorStartIPNot1 int64 = 6

// SubnetValidErrorInvalidGatewayAddress : Invalid gateway address
var SubnetValidErrorInvalidGatewayAddress int64 = 7

// SubnetValidErrorGatewayNotInSubnet : Gateway IP address in not in the subnet
var SubnetValidErrorGatewayNotInSubnet int64 = 8

// SubnetValidErrorSubnetIsUsedByIface : Subnet is used by one of iface
var SubnetValidErrorSubnetIsUsedByIface int64 = 9
