package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

func main() {
	inputDir := "input"
	outDir := "output"

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("Failed to read input directory: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".pdf") {
			continue
		}

		inputPath := filepath.Join(inputDir, file.Name())
		outPath := filepath.Join(outDir, strings.TrimSuffix(file.Name(), ".pdf")+"_unlocked.pdf")

		fmt.Printf("Enter password for %s: ", file.Name())
		password, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Failed to read password for %s: %v\n", file.Name(), err)
			continue
		}
		password = strings.TrimSpace(password)

		err = unlockPdf(inputPath, outPath, password)
		if err != nil {
			fmt.Printf("Failed to unlock %s: %v\n", file.Name(), err)
		} else {
			fmt.Printf("Successfully unlocked %s\n", file.Name())
		}
	}
}

func unlockPdf(inputPath, outputPath, password string) error {
	conf := model.NewDefaultConfiguration()
	conf.UserPW = password
	return api.DecryptFile(inputPath, outputPath, conf)
}
