package init

import "hcc/harp/lib/syscheck"

func syscheckInit() error {
	err := syscheck.CheckRoot()
	if err != nil {
		return err
	}

	err = syscheck.CheckArpingCommand()
	if err != nil {
		return err
	}

	return nil
}
