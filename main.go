package main

import (
	"errors"
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
	inputFile := flag.String("i", "", "Desired zip file")
	fileName := flag.String("n", "", "Desired zip file name")
	emptyFile := flag.Bool("empty", false, "Create a empty archive")
	archiveType := flag.String("t", "zip", "Archive type (zip/tar/gzip)")
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
	printStatus("Out File: " + *outFileName)
	printStatus("Given File Name: " + *fileName)
	print("\n")

	switch *archiveType {
	case "zip":
		zipWriter, err := zippo.NewZipFile(*outFileName)
		fatal(err)
		printStatus("Forging zip file...")
		fatal(zippo.AddFileToZIP(zipWriter, *inputFile, *fileName))
		zipWriter.Flush()
		zipWriter.Close()
	case "tar":
		tarWriter, err := zippo.NewTARFile(*outFileName)
		fatal(err)
		printStatus("Forging tar file...")
		fatal(zippo.AddFileToTAR(tarWriter, *inputFile, *fileName))
		defer tarWriter.Flush()
		defer tarWriter.Close()
	case "gzip":
		tarWriter, gzipWriter, err := zippo.NewTARGZFile(*outFileName)
		fatal(err)
		printStatus("Forging tar.gz file...")
		fatal(zippo.AddFileToTAR(tarWriter, *inputFile, *fileName))
		defer gzipWriter.Flush()
		defer gzipWriter.Close()
		defer tarWriter.Flush()
		defer tarWriter.Close()
	default:
		fatal(errors.New("unknown archive type"))
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
	banner, _ := b64.StdEncoding.DecodeString("ICAgICAgICAgICAgICAgICwuflwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgLC1gICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICBcICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICBcICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAgICAgICAgICBcICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICBcICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgIF8uLS0tLS0tLS5cICAgICAgIFwgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgIAogICAgICAob3wgbyBvIG8gfCBcICAgIC4tYCAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgX198fG9fb19vX298X2FkLWBgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgIHxgYGBgYGBgYGBgYGBgYHwKICAgIHwgICAgIFpJUFBPICAgIHwgIAogICAgfCAgIOKZoCDimaAg4pmgIOKZoCDimaAgIHwgCiAgICB8ICAgICDimaAg4pmgIOKZoCAgICB8CiAgICB8ICAgICAgIOKZoCAgICAgIHwKICAgIHxfX19fX19fX19fX19fX3wKPT09PT09PT1AZWdlYmxjPT09PT09PT09PQo=")
	white.Print(string(banner))
	print("\n")
}
