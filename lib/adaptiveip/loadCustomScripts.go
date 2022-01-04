package adaptiveip

import (
	"hcc/harp/lib/cmd"
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"os"
	"path/filepath"
)

func loadCustomScripts() error {
	var files []string

	folder := config.AdaptiveIP.CustomScriptsLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if path != config.AdaptiveIP.CustomScriptsLocation {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	for _, file := range files {
		logger.Logger.Println("Running " + file + "...")
		err = cmd.RunScript(file)
		if err != nil {
			logger.Logger.Println("LoadCustomScripts(): Error occurred while running " + file)
			logger.Logger.Println("########## Error message start ##########")
			logger.Logger.Println(err.Error())
			logger.Logger.Println("########## Error message end ##########")
		}
	}

	return nil
}
