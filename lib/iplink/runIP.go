package iplink

import (
	"hcc/harp/lib/cmd"
)

func runIP(args string) error {
	return cmd.RunCMD("ip " + args)
}
