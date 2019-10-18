package floatingip

import (
	"hcc/harp/lib/logger"
)

// CreateFloatingIP : Create floating IP for server.
func CreateFloatingIP() error {
	var err error

	logger.Logger.Println("Create FloatingIP")

	if err != nil {
		logger.Logger.Println(err)
		return err
	}

	return nil
}
