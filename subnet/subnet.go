package subnet

import (
	"bytes"
	"fmt"
	"hcc/harp/logger"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func printCommand(cmd *exec.Cmd) {
	fmt.Printf("==> Executing: %s\n", strings.Join(cmd.Args, ""))
}

func printError(err error) {
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error : %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf("==> Output: %s\n", string(outs))
	}
}

func ReadFile() {
	data, err := ioutil.ReadFile("test.txt")
	if err != nil {
		log.Panicf("failed reading data from file: %s", err)
	}
	fmt.Printf("\nFile Content: %s", data)
}

func AppendFile() {
	file, err := os.OpenFile("test.txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	len, err := file.WriteString("\nsubnet {\n\tThe Go language was conceived in Today\n}")
	if err != nil {
		fmt.Printf("\nLength: %d bytes", len)
		fmt.Printf("\nFile Name: %s", file.Name())
	}
}

func EditFile() {
	file, err := os.OpenFile("test.txt", os.O_RDWR, 0644)
	if err != nil {
		log.Fatalf("failed opening file : %s", err)
	}
	defer file.Close()
	len, err := file.WriteAt([]byte{'S'}, 0)
	if err != nil {
		log.Fatalf("failed writing to file: %s", err)
	}
	fmt.Printf("\nLength: %d bytes", len)
	fmt.Printf("\nFile Name: %s", file.Name())
}

func UpdateSubnet() error {
	var err error

	logger.Logger.Println("Create Subnet")

	if err != nil {
		logger.Logger.Println(err.Error())
		return err
	}

	cmd := exec.Command("go", "version")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	printCommand(cmd)
	cmdErr := cmd.Run()
	printError(cmdErr)
	printOutput(cmdOutput.Bytes())

	//ReadFile()
	//EditFile()
	//ReadFile()
	AppendFile()

	return nil
}
