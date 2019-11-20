package adaptiveip

import (
	"hcc/harp/lib/config"
	"hcc/harp/lib/logger"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

func flushPFRules() error {
	logger.Logger.Println("Flushing pf rules...")

	cmd := exec.Command("pfctl", "-F", "all")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func loadPFRules(pfRulesFileLocation string) error {
	logger.Logger.Println("Loading pf rules from " + pfRulesFileLocation + "...")

	cmd := exec.Command("pfctl", "-f", pfRulesFileLocation)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getBinatAnchorConfigFiles() ([]string, error) {
	var files []string

	folder := config.AdaptiveIP.PFBinatConfigFileLocation
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func loadExstingBinatAnchorServersRules() error {
	logger.Logger.Println("Loading existing binat anchor servers rules...")

	configFiles, err := getBinatAnchorConfigFiles()
	if err != nil {
		return err
	}
	if len(configFiles) == 1 {
		return nil
	}

	var RoutineMAX = int(config.AdaptiveIP.ArpingRoutineMaxNum)
	if RoutineMAX == 0 {
		RoutineMAX = 5
	}
	var routineMax = RoutineMAX
	var wait sync.WaitGroup

	for i := 0; i < len(configFiles); {
		if len(configFiles)-i < RoutineMAX {
			routineMax = len(configFiles) - i
		}

		wait.Add(routineMax)

		for j := 0; j < routineMax; j++ {
			if configFiles[i] == config.AdaptiveIP.PFBinatConfigFileLocation {
				j--
				i++
				wait.Done()
				continue
			}

			go func(file string) {
				var binatanchorFileName string
				var binatanchorName string
				var err error

				binatanchorFileName = file[len(config.AdaptiveIP.PFBinatConfigFileLocation+"/"):]
				if !strings.Contains(binatanchorFileName, binatanchorFilenamePrefix) ||
					!strings.Contains(binatanchorFileName, ".conf") {
					logger.Logger.Println("Wrong binat anchor filename: " + binatanchorFileName)
					logger.Logger.Println("Filename must be as '" + binatanchorFilenamePrefix + "XXX'")
					goto RoutineDone
				}

				binatanchorName = binatanchorFileName[0 : len(binatanchorFileName)-len(".conf")]
				err = LoadPFBinatAnchorRule(binatanchorName, file)
				if err != nil {
					logger.Logger.Println(err)
				}

			RoutineDone:
				wait.Done()
			}(configFiles[i])

			i++
			if i == len(configFiles) {
				break
			}
		}

		wait.Wait()
	}

	return nil
}

// LoadPFBinatAnchorRule : Load binat anchor rules configuration file and apply it to pf firewall
func LoadPFBinatAnchorRule(binatanchorName string, binatanchorConfigFileLocation string) error {
	logger.Logger.Println("Loading binat anchor of " + binatanchorName + "...")

	IP := binatanchorName[len(binatanchorFilenamePrefix):]
	err := CheckDuplicatedIPAddress(IP)
	if err != nil {
		return err
	}

	cmd := exec.Command("pfctl", "-a", binatanchorName, "-f", binatanchorConfigFileLocation)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// RemvoePFBinatAnchorRule : Remove binat anchor rules of provided name from pf firewall
func RemvoePFBinatAnchorRule(binatanchorName string) error {
	logger.Logger.Println("Removing binat anchor rules of " + binatanchorName + "...")

	cmd := exec.Command("pfctl", "-a", binatanchorName, "-F", "all")
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// LoadHarpPFRules : Load pf rules for harp module
func LoadHarpPFRules() error {
	err := flushPFRules()
	if err != nil {
		return err
	}

	err = loadPFRules(config.AdaptiveIP.PFRulesFileLocation)
	if err != nil {
		return err
	}

	err = loadExstingBinatAnchorServersRules()
	if err != nil {
		return err
	}

	return nil
}
