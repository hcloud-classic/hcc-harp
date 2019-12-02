package adaptiveip

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/fileutil"
	"hcc/harp/lib/logger"
	"os"
	"os/exec"
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

func getIfconfigScriptFilesInternal() ([]string, error) {
	var files []string

	folder := config.DHCPD.IfconfigScriptFileLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func loadIfconfigScript(filepath string) error {
	logger.Logger.Println("Loading ifconfig script file: " + filepath)

	cmd := exec.Command("csh", filepath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func loadExistingIfconfigScriptsExternal() error {
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

func loadExistingIfconfigScriptsInternal() error {
	logger.Logger.Println("Loading existing ifconfig scripts for internal...")

	scriptFiles, err := getIfconfigScriptFilesInternal()
	if err != nil {
		return err
	}
	if len(scriptFiles) == 1 {
		return nil
	}

	var ifconfigScriptFileName string

	for i := 0; i < len(scriptFiles); i++ {
		if scriptFiles[i] == config.DHCPD.IfconfigScriptFileLocation {
			continue
		}

		ifconfigScriptFileName = scriptFiles[i][len(config.DHCPD.IfconfigScriptFileLocation+"/"):]
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

// CreateAndLoadIfconfigScriptInternal : Create and load ifconfig script for internal interface
func CreateAndLoadIfconfigScriptInternal(internelIfacename string, privateGatewayIP string, netmaskPrivate string) error {
	var ifconfigInternalScriptData string
	ifconfigInternalScriptData = ifconfigReplaceString
	ifconfigInternalScriptData = strings.Replace(ifconfigInternalScriptData, "IFCONFIG_IFACE_NAME", internelIfacename, -1)
	ifconfigInternalScriptData = strings.Replace(ifconfigInternalScriptData, "IFCONFIG_IP", privateGatewayIP, -1)
	ifconfigInternalScriptData = strings.Replace(ifconfigInternalScriptData, "IFCONFIG_NETMASK", netmaskPrivate, -1)

	var ifconfigScriptData string
	ifconfigScriptData = ifconfigSHELL + ifconfigInternalScriptData
	ifconfigScriptData = strings.Replace(ifconfigScriptData, "ALIAS_STATE", "alias", -1)

	ifconfigScriptFileName := ifconfigFilenamePrefix + privateGatewayIP + ".sh"
	logger.Logger.Println("CreateAndLoadIfconfigScriptInternal: Creating ifconfig script file: " + ifconfigScriptFileName)
	ifconfigScriptFileLocation := config.DHCPD.IfconfigScriptFileLocation + "/" + ifconfigScriptFileName

	err := logger.CreateDirIfNotExist(config.DHCPD.IfconfigScriptFileLocation)
	if err != nil {
		return err
	}

	err = fileutil.WriteFile(ifconfigScriptFileLocation, ifconfigScriptData)
	if err != nil {
		return err
	}

	logger.Logger.Println("CreateAndLoadIfconfigScriptInternal: Running ifconfig script file: " + ifconfigScriptFileName)
	err = loadIfconfigScript(ifconfigScriptFileLocation)
	if err != nil {
		return err
	}

	return nil
}

func createAndLoadIfconfigScriptExternal(externelIfacename string, publicIP string, netmaskPublic string) error {
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

	err := logger.CreateDirIfNotExist(config.AdaptiveIP.IfconfigScriptFileLocation)
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
