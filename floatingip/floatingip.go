package floatingip

import (
	"hcc/harp/logger"
)

// CreateFloatingIP : Create floating IP for server.
func CreateFloatingIP() error {
	var err error

	logger.Logger.Println("Create FloatingIP")

	if err != nil {
		logger.Logger.Println(err.Error())
		return err
	}

	return nil
}
