package floatingip

import (
	"hcc/harp/logger"
)

func CreateFloatingip() error {
	var err error

	logger.Logger.Println("Create FloatingIP")

	if err != nil {
		logger.Logger.Println(err.Error())
		return err
	}

	return nil
}
