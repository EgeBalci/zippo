package zippo

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
)

// NewTARGZFile returns a Tar.gz writer
func NewTARGZFile(gzipFileName string) (*tar.Writer, *gzip.Writer, error) {
	newTargzFile, err := os.OpenFile(gzipFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, nil, err
	}
	gzipWriter := gzip.NewWriter(newTargzFile)
	return tar.NewWriter(gzipWriter), gzipWriter, nil
}

// NewTARFile creates a new Tar file
func NewTARFile(tarFileName string) (*tar.Writer, error) {
	newTarFile, err := os.OpenFile(tarFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}
	return tar.NewWriter(newTarFile), nil
}

// NewZipFile creates a new zip file
func NewZipFile(zipFileName string) (*zip.Writer, error) {
	newZipFile, err := os.OpenFile(zipFileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0664)
	if err != nil {
		return nil, err
	}
	zipWriter := zip.NewWriter(newZipFile)
	return zipWriter, nil
}

// ReadZIP reads the zip file bytes and returns a Zippo struct
func ReadZIP(zipFileName string) (*zip.Writer, error) {
	zipFile, err := os.Open(zipFileName)
	if err != nil {
		return nil, err
	}
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	return zipWriter, nil
}

// ReadTAR reads the zip file bytes and returns a Zippo struct
func ReadTAR(tarFileName string) (*tar.Writer, error) {
	tarFile, err := os.Open(tarFileName)
	if err != nil {
		return nil, err
	}
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	return tarWriter, nil
}

// AddFileToZIP adds a new file to the ZIP archive by the given name
func AddFileToZIP(zipWriter *zip.Writer, fileName, desiredName string) error {

	rawFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	f, err := zipWriter.Create(desiredName)
	if err != nil {
		return err
	}

	// // Using FileInfoHeader() above only uses the basename of the file. If we want
	// // to preserve the folder structure we can overwrite this with the full path.
	// header.Name = desiredName

	// // Change to deflate to gain better compression
	// // see http://golang.org/pkg/archive/zip/#pkg-constants
	// header.Method = zip.Deflate
	// writer, err := zipWriter.CreateHeader(header)
	// if err != nil {
	// 	return err
	// }
	// _, err = io.Copy(writer, fileToZip)
	// return err

	_, err = f.Write(rawFile)
	return err
}

// AddFileToTAR adds a new file to the ZIP archive by the given name
func AddFileToTAR(tarWriter *tar.Writer, fileName, desiredName string) error {

	// rawFile, err := ioutil.ReadFile(fileName)
	// if err != nil {
	// 	return err
	// }

	file, err := os.Open(fileName)
	if err != nil {
		return nil
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return nil
	}

	// create a new dir/file header
	header, err := tar.FileInfoHeader(fileInfo, file.Name())
	if err != nil {
		return err
	}

	header.Name = desiredName

	// write the header
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	if _, err := io.Copy(tarWriter, file); err != nil {
		return err
	}

	file.Close()
	return nil
}

// NewZipBomb creates a ZIP containing jump file with huge size when deflated
func NewZipBomb(size int) []byte {
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	zipWriter := zip.NewWriter(buf)

	// Add some files to the archive.
	body := ""

	for i := 0; i < size*1024*1024; i++ {
		body += "A"
	}

	zipFile, _ := zipWriter.Create("readme.txt")
	zipFile.Write([]byte(body))
	zipWriter.Close()
	return buf.Bytes()
}
