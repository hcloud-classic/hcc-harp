[![pipeline status](http://210.207.104.150:8100/iitp-sds/harp/badges/master/pipeline.svg)](http://210.207.104.150:8100/iitp-sds/harp/pipelines)
[![coverage report](http://210.207.104.150:8100/iitp-sds/harp/badges/master/coverage.svg)](http://210.207.104.150:8100/iitp-sds/harp/commits/master)
[![go report](http://210.207.104.150:8100/iitp-sds/hcloud-badge/raw/feature/dev/hcloud-badge_harp.svg)](http://210.207.104.150:8100/iitp-sds/hcloud-badge/raw/feature/dev/goreport_harp)


## Harp

- Features
  - Subnet Management
  - Master Node DHCP Management
  - Configure network interfaces of Master Node
  - Allocate Adaptive IP (Allocate public IP address to private IP address.)

- Supported OS
  - FreeBSD
  - Linux

- Pre-required
  - For FreeBSD
    - PF firewall enabled with kernel module loaded.
  - For Linux
    - SELinux disabled (If enabled, iptables will not work correctly.)
    - iptables installed with NAT kernel module loaded (iptable_nat, nf_nat)
  - For both FreeBSD and Linux
    - Golang installed
    - 2 network interfaces for use an external network and internal networks.

<br>

- How to build
  - FreeBSD
    - Just run `gmake` command.
  - Linux
    - Just run `make` command.

<br>

- How to run
  - FreeBSD
    1. Copy `harp.conf_FreeBSD` to `/etc/hcc/harp/harp.conf`
    2. Change your settings in `harp.conf`
    3. Run `harp` binary.

  - Linux
    1. Copy `harp.conf_Linux` to `/etc/hcc/harp/harp.conf`
    2. Change your settings in `harp.conf`
    3. Run `harp` binary.

  <br>

#### Adaptive IP

- How it allocates public IP addresses.
  1. Get server's UUID and public IP address from user.
  2. Check if provided server's UUID is already used in Adaptive IP.
  3. Get the private subnet information related with sever.
  4. Get the first IP address. (End with x.x.x.1. This is Leader Node's IP address.)
  5. Create NAT firewall rules.
  6. Server's nodes are now connect to Internet and can connect from external network by provided public IP address.

  <br>

- How it listing available public IP addresses.
  1. Make a array from start IP address to end IP address that configured as Adaptive IP range in `harp.conf`.
  2. Check from start IP address to end IP. First, check if the IP address is configure to external network interface.
  3. Second, send ARP request. If it received ARP reply, then the IP address is duplicate with someone.
  4. Show the available IP addresses except in 2 and 3 .

<br>

#### Additional infos

- See configuration comments in `.go` files located in `./lib/config/`.
