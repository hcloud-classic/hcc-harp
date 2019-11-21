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

func getIfconfigScriptFiles() ([]string, error) {
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

func loadIfconfigScript(filepath string) error {
	logger.Logger.Println("Loading ifconfig script file: " + filepath)

	cmd := exec.Command("csh", filepath)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func loadExistingIfconfigScripts() error {
	logger.Logger.Println("Loading existing ifconfig scripts...")

	scriptFiles, err := getIfconfigScriptFiles()
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

func createAndLoadIfconfigScript(internelIfacename string, externelIfacename string,
	privateGatewayIP string, publicIP string, netmaskPrivate string,
	netmaskPublic string, isAdd bool) error {
	var ifconfigInternalScriptData string
	ifconfigInternalScriptData = ifconfigReplaceString
	ifconfigInternalScriptData = strings.Replace(ifconfigInternalScriptData, "IFCONFIG_IFACE_NAME", internelIfacename, -1)
	ifconfigInternalScriptData = strings.Replace(ifconfigInternalScriptData, "IFCONFIG_IP", privateGatewayIP, -1)
	ifconfigInternalScriptData = strings.Replace(ifconfigInternalScriptData, "IFCONFIG_NETMASK", netmaskPrivate, -1)

	var ifconfigExternalScriptData string
	ifconfigExternalScriptData = ifconfigReplaceString
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_IFACE_NAME", externelIfacename, -1)
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_IP", publicIP, -1)
	ifconfigExternalScriptData = strings.Replace(ifconfigExternalScriptData, "IFCONFIG_NETMASK", netmaskPublic, -1)

	var ifconfigScriptData string
	ifconfigScriptData = ifconfigSHELL + ifconfigInternalScriptData + ifconfigExternalScriptData

	var aliasState string
	if isAdd {
		aliasState = "alias"
	} else {
		aliasState = "-alias"
	}
	ifconfigScriptData = strings.Replace(ifconfigScriptData, "ALIAS_STATE", aliasState, -1)

	ifconfigScriptFileName := ifconfigFilenamePrefix + publicIP + ".sh"
	logger.Logger.Println("createAndLoadIfconfigScript: Creating ifconfig script file: ifconfigScriptFileName")
	ifconfigScriptFileLocation := config.AdaptiveIP.IfconfigScriptFileLocation + "/" + ifconfigScriptFileName
	err := fileutil.WriteFile(ifconfigScriptFileLocation, ifconfigScriptData)
	if err != nil {
		return err
	}

	logger.Logger.Println("createAndLoadIfconfigScript: Running ifconfig script file: " + ifconfigScriptFileName)
	err = loadIfconfigScript(ifconfigScriptFileLocation)
	if err != nil {
		return err
	}

	return nil
}
