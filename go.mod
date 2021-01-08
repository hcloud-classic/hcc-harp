module hcc/harp

go 1.13

require (
	github.com/Terry-Mao/goconf v0.0.0-20161115082538-13cb73d70c44
	github.com/apparentlymart/go-cidr v1.0.1
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gojp/goreportcard v0.0.0-20201106142952-232d912e513e // indirect
	github.com/golang/protobuf v1.4.3
	github.com/hcloud-classic/hcc_errors v1.1.0
	github.com/hcloud-classic/pb v0.0.0
	github.com/mactsouk/go v0.0.0-20180603081621-6a282087f7bd // indirect
	github.com/mdlayher/arp v0.0.0-20191213142603-f72070a231fc
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/opencontainers/selinux v1.6.0
	golang.org/x/lint v0.0.0-20201208152925-83fdc39ff7b5 // indirect
	golang.org/x/tools v0.0.0-20210107193943-4ed967dd8eff // indirect
	google.golang.org/appengine v1.6.5 // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0 // indirect
)

replace github.com/hcloud-classic/pb => ../pb
