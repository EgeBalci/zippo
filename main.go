package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	zippo "zippo/lib"

	b64 "encoding/base64"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

func main() {

	banner()
	inputFile := flag.String("i", "", "Desired zip file name")
	fileName := flag.String("n", "", "Desired zip file name")
	emptyFile := flag.Bool("empty", false, "Create a empty file")
	tarMode := flag.Bool("tar", false, "Create a TAR file instead")
	targzMode := flag.Bool("targz", false, "Create a TAR file instead")
	outFileName := flag.String("o", "", "Output zip file")
	flag.Parse()

	if len(os.Args) < 2 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *emptyFile {
		*inputFile = tempFile()
		defer removeFile(*inputFile)
	}
	printStatus("Input File: " + *inputFile)
	printStatus("Target File: " + *outFileName)
	printStatus("File Name: " + *fileName)
	print("\n")
	printStatus("Creating a new zippo...")

	if *tarMode {
		tarWriter, err := zippo.NewTARFile(*outFileName)
		fatal(err)
		printStatus("Adding files...")
		fatal(zippo.AddFileToTAR(tarWriter, *inputFile, *fileName))
		defer tarWriter.Flush()
		defer tarWriter.Close()
	} else if *targzMode {
		tarWriter, gzipWriter, err := zippo.NewTARGZFile(*outFileName)
		fatal(err)
		printStatus("Adding files...")
		fatal(zippo.AddFileToTAR(tarWriter, *inputFile, *fileName))
		defer gzipWriter.Flush()
		defer gzipWriter.Close()
		defer tarWriter.Flush()
		defer tarWriter.Close()
	} else {
		zipWriter, err := zippo.NewZipFile(*outFileName)
		fatal(err)
		printStatus("Adding files...")
		fatal(zippo.AddFileToZIP(zipWriter, *inputFile, *fileName))
		zipWriter.Flush()
		zipWriter.Close()
	}

	printSuccess("Done !")
}

func tempFile() string {
	tmpfile, err := ioutil.TempFile("", "."+uuid.New().String())
	fatal(err)
	return tmpfile.Name()
}

func removeFile(name string) {
	err := os.Remove(name)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			log.Fatalf("[%s] ERROR: %s\n", strings.ToUpper(strings.Split(details.Name(), ".")[1]), err)
		} else {
			log.Fatalf("[UNKOWN] ERROR: %s\n", err)
		}
	}
}

func printFail(str string) {
	red := color.New(color.FgRed).Add(color.Bold)
	red.Print("[-] ")
	fmt.Println(str)
}

func printStatus(str string) {
	yellow := color.New(color.FgYellow).Add(color.Bold)
	yellow.Print("[*] ")
	fmt.Println(str)
}

func printSuccess(str string) {
	green := color.New(color.FgGreen).Add(color.Bold)
	green.Print("[+] ")
	fmt.Println(str)
}

func banner() {
	white := color.New(color.FgWhite).Add(color.Bold)
	fmt.Println("") // Line feed
	banner, _ := b64.StdEncoding.DecodeString("ICAgICAgICAgICAgICAgICwuflwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgLC1gICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICBcICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICBcICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICBcICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICBcICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgIF8uLS0tLS0tLS5cICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAob3wgbyBvIG8gfCBcICAgIC4tYCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgX198fG9fb19vX298X2FkLWBgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgIHxgYGBgYGBgYGBgYGBgYHwKICAgIHwgICAgIFpJUFBPICAgIHwgIAogICAgfCAgIOKZoCDimaAg4pmgIOKZoCDimaAgIHwgCiAgICB8ICAgICDimaAg4pmgIOKZoCAgICB8CiAgICB8ICAgICAgIOKZoCAgICAgIHwKICAgIHxfX19fX19fX19fX19fX3wKPT09PT09PT1FR0UtQkFMQ0k9PT09PT09PT0K")
	white.Print(string(banner))
	print("\n")
}
