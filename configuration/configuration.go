//Package configuration has all the functions and structures required for the application.
package configuration

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/melvinodsa/fileUpload/modules/qr"
)

//DownloadDetails structure for storing download file details.
type DownloadDetails struct {
	//DownloadCode for the files.
	DownloadCode string
}

//UploadDetails structure for storing the upload file details.
type UploadDetails struct {
	//UploadFiles to store the files.
	UploadFiles []multipart.File
	//UploadFileNames to store the filenames.
	UploadFileNames []string
}

//IPDetails has the interface and its corresponding ip
type IPDetails struct {
	//IPAddr holds the ipaddres
	IPAddr string
	//NWInterface holds the interface name
	NWInterface string
}

//UploadFileDetails structure has IPDetails array and Shortcode
type UploadFileDetails struct {
	//ShortCode holds the shortcode for the download
	ShortCode string
	//IPDetailsConfig holds the IPDetails array
	IPDetailsConfig []IPDetails
}

//Configuration has the required configuration variables
type Configuration struct {
	//ShortCode for file uploads and downloads
	ShortCode int64
	//IPDetailsConfig array conatins the IPDetails
	IPDetailsConfig []IPDetails
}

var (
	//ConfigurationLoaded will have the configuration of the application
	ConfigurationLoaded Configuration
)

//GetShortCode method will return user credentails if available.
func GetShortCode() int64 {
	return ConfigurationLoaded.ShortCode
}

//SaveShortCode function will save the user credentials
func SaveShortCode(shortCode int64) error {
	ConfigurationLoaded.ShortCode = shortCode
	err := SaveConfiguration()
	if err != nil {
		log.Fatal("Error in while saving the shortcode. ", err)
		return err
	}
	log.Println("Shortcode has been saved successfully to the conf file.")
	return nil
}

//GetIPDetailsConfig function will return user credentails if available.
func GetIPDetailsConfig() []IPDetails {
	return ConfigurationLoaded.IPDetailsConfig
}

//SaveIPDetailsConfig function will save the user credentials
func SaveIPDetailsConfig(iPDetailsConfig []IPDetails) error {
	ConfigurationLoaded.IPDetailsConfig = iPDetailsConfig
	err := SaveConfiguration()
	if err != nil {
		log.Fatal("Error in while saving the iPDetailsConfig. ", err)
		return err
	}
	log.Println("IPDetailsConfig has been saved successfully to the conf file.")
	return nil
}

//CacheConfiguration method will cache the configuration file contents.
func CacheConfiguration() {
	os.Remove("conf.json")
	// open output file
	fo, err := os.Create("conf.json")
	if err != nil {
		log.Fatal("Error while creating conf.json file.")
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			log.Fatal("Error in closing the conf.json file.")
		}
	}()
	err = ioutil.WriteFile("conf.json", []byte("{\n}"), 0644)
	os.RemoveAll("./resources/downloads")
	if err := os.Mkdir("./resources/downloads", 0777); nil != err {
		log.Fatal("Error in creating the folder ./resources/downloads.", err)
	}
	os.RemoveAll("./logs")
	if err := os.Mkdir("./logs", 0777); nil != err {
		log.Fatal("Error in creating the folder ./logs.", err)
	}
	file, err := os.Open("conf.json")
	if err != nil {
		log.Fatal("Error in opening the conf file while caching. ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ConfigurationLoaded)
	if err != nil {
		log.Fatal("Error in reading the conf file while caching. ", err)
	}
	log.Println("Data has been cached from conf file. ")
}

//SaveConfiguration function will save the user credentials
func SaveConfiguration() error {
	file, err := os.Create("conf.json")
	if err != nil {
		log.Fatal("Error in opening the conf file while saving configuration. ", err)
		return err
	}
	fileData, jerr := json.MarshalIndent(ConfigurationLoaded, "", " ")
	if jerr != nil {
		log.Fatal("Error in parsing the json while saving configuration. ", err)
		return jerr
	}
	_, err = file.Write(fileData)
	if err != nil {
		log.Fatal("Error in writinig to conf file while saving configuration. ", err)
		return err
	}
	log.Println("Configuration has been saved successfully to the conf file.")
	return nil
}

//Zipit function will zip a folder
func Zipit(source, target string, keepHierarchy bool) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if keepHierarchy && info.IsDir() {
		baseDir = filepath.Base(source)
	}

	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !keepHierarchy && source == path {
			return nil
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		if keepHierarchy && baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

//QRCodeGenerate function will generate the QR Code for particular string
func QRCodeGenerate(data, target string) error {
	//generating the qr code for the same
	code, err := qr.Encode(data, qr.H)
	if err != nil {
		log.Fatal("Error in generating the QR code for data", data, err)
		return err
	}

	imgByte := code.PNG()
	// convert byte to image for saving to file
	img, _, _ := image.Decode(bytes.NewReader(imgByte))

	//save the imgByte to file
	out, err := os.Create(target)

	if err != nil {
		log.Fatal("Error in creating the file to save QR code for data", data, err)
		return err
	}

	err = png.Encode(out, img)

	if err != nil {
		log.Fatal("Error in ecoding the QR code for data", data, err)
		return err
	}

	// everything ok
	log.Println("QR code generated and saved to", target)
	return nil
}
