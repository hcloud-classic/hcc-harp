package ifconfig

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/logger"
	"os"
	"path/filepath"
	"strings"
)

func getIfconfigScriptFilesExternal() ([]string, error) {
	var files []string

	folder := config.AdaptiveIP.IfconfigScriptFileLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// LoadExistingIfconfigScriptsExternal : Load existing ifconfig scripts for external network.
func LoadExistingIfconfigScriptsExternal() error {
	logger.Logger.Println("Loading existing ifconfig scripts for external...")

	scriptFiles, err := getIfconfigScriptFilesExternal()
	if err != nil {
		return err
	}
	if len(scriptFiles) == 1 {
		return nil
	}

	var ifconfigScriptFileName string

	for i := 0; i < len(scriptFiles); i++ {
		if scriptFiles[i] == config.AdaptiveIP.IfconfigScriptFileLocation {
			continue
		}

		ifconfigScriptFileName = scriptFiles[i][len(config.AdaptiveIP.IfconfigScriptFileLocation+"/"):]
		if !strings.Contains(ifconfigScriptFileName, ifconfigFilenamePrefix) ||
			!strings.Contains(ifconfigScriptFileName, ".sh") {
			logger.Logger.Println("Wrong ifconfig script filename: " + ifconfigScriptFileName)
			logger.Logger.Println("Filename must be as '" + ifconfigFilenamePrefix + "XXX.sh'")
			continue
		}

		err = loadIfconfigScript(scriptFiles[i])
		if err != nil {
			logger.Logger.Println(err)
		}
	}

	return nil
}

// CreateAndLoadIfconfigScriptExternal : Create and load ifconfig scripts for external network
func CreateAndLoadIfconfigScriptExternal(externelIfacename string, publicIP string, netmaskPublic string) error {
	var ifconfigExternalScriptData string
	ifconfigExternalScriptData = ifconfigReplaceString
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_IFACE_NAME", externelIfacename, -1)
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_IP", publicIP, -1)
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_NETMASK", netmaskPublic, -1)

	var ifconfigScriptData string
	ifconfigScriptData = ifconfigSHELL + ifconfigExternalScriptData
	ifconfigScriptData = strings.Replace(ifconfigScriptData, "ALIAS_STATE", "alias", -1)

	ifconfigScriptFileName := ifconfigFilenamePrefix + publicIP + ".sh"
	logger.Logger.Println("createAndLoadIfconfigScriptExternal: Creating ifconfig script file: " + ifconfigScriptFileName)
	ifconfigScriptFileLocation := config.AdaptiveIP.IfconfigScriptFileLocation + "/" + ifconfigScriptFileName

	err := fileutil.CreateDirIfNotExist(config.AdaptiveIP.IfconfigScriptFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(ifconfigScriptFileLocation, ifconfigScriptData)
	if err != nil {
		return err
	}

	logger.Logger.Println("createAndLoadIfconfigScriptExternal: Running ifconfig script file: " + ifconfigScriptFileName)
	err = loadIfconfigScript(ifconfigScriptFileLocation)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAndUnloadIfconfigScriptExternal : Delete and unload ifconfig scripts for external network
func DeleteAndUnloadIfconfigScriptExternal(externelIfacename string, publicIP string, netmaskPublic string) error {
	var ifconfigExternalScriptData string
	ifconfigExternalScriptData = ifconfigReplaceString
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_IFACE_NAME", externelIfacename, -1)
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_IP", publicIP, -1)
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_NETMASK", netmaskPublic, -1)

	var ifconfigScriptData string
	ifconfigScriptData = ifconfigSHELL + ifconfigExternalScriptData
	ifconfigScriptData = strings.Replace(ifconfigScriptData, "ALIAS_STATE", "-alias", -1)

	ifconfigScriptFileName := ifconfigFilenamePrefix + publicIP + ".sh"
	logger.Logger.Println("DeleteAndUnloadIfconfigScriptExternal: Creating ifconfig temporary script file: " + ifconfigScriptFileName)
	ifconfigScriptFileLocation := config.AdaptiveIP.IfconfigScriptFileLocation + ifconfigScriptFileName
	ifconfigScriptTemporaryFileLocation := config.AdaptiveIP.IfconfigScriptFileLocation + "/tmp/" + ifconfigScriptFileName

	err := fileutil.CreateDirIfNotExist(config.AdaptiveIP.IfconfigScriptFileLocation + "/tmp/")
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(ifconfigScriptTemporaryFileLocation, ifconfigScriptData)
	if err != nil {
		return err
	}

	logger.Logger.Println("DeleteAndUnloadIfconfigScriptExternal: Running ifconfig temporary script file: " + ifconfigScriptFileName)
	err = loadIfconfigScript(ifconfigScriptTemporaryFileLocation)
	if err != nil {
		logger.Logger.Println(err.Error())
	}

	logger.Logger.Println("DeleteAndUnloadIfconfigScriptExternal: Deleting ifconfig temporary script file: " + ifconfigScriptFileName)
	err = fileutil.DeleteFile(ifconfigScriptTemporaryFileLocation)
	if err != nil {
		logger.Logger.Println(err.Error())
	}

	err = fileutil.DeleteDir(config.AdaptiveIP.IfconfigScriptFileLocation + "/tmp/")
	if err != nil {
		logger.Logger.Println(err.Error())
	}

	logger.Logger.Println("DeleteAndUnloadIfconfigScriptExternal: Deleting ifconfig script file: " + ifconfigScriptFileName)
	err = fileutil.DeleteFile(ifconfigScriptFileLocation)
	if err != nil {
		logger.Logger.Println(err.Error())
	}

	return nil
}