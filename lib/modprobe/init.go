package modprobe

import (
	"errors"
	"hcc/harp/lib/logger"
	"strings"
)

// LoadHarpKernelModules : Load kernel modules needed for harp
func LoadHarpKernelModules() error {
	modules := strings.Split(neededKernelModulesForHarp, ", ")

	for _, module := range modules {
		err := Load(module, "")
		logger.Logger.Println("Loading kernel module: " + module)
		if err != nil {
			return errors.New(err.Error() + " (module: " + module + ")")
		}
	}

	return nil
}
