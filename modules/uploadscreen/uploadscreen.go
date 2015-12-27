//Package uploadscreen contains the functions required for the download screen
package uploadscreen

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/melvinodsa/fileUpload/configuration"
)

//UploadScreen generate the homepage of the aplication.
func UploadScreen(w http.ResponseWriter, r *http.Request) {
	log.Println("User requested upload view.")

	//Generating page from templates
	t := template.Must(template.New("header").ParseFiles("templates/snippets/homeheader.html"))
	t.Execute(w, nil)
	t = template.Must(template.New("content").ParseFiles("templates/snippets/upload.html"))
	t.Execute(w, nil)
	t = template.Must(template.New("footer").ParseFiles("templates/snippets/homefooter.html"))
	t.Execute(w, nil)
}

//Upload generate the homepage of the aplication.
func Upload(w http.ResponseWriter, r *http.Request) {
	ips := configuration.UploadFileDetails{}
	shortCode := configuration.GetShortCode()
	shortCode++
	if shortCode >= 999999 {
		shortCode = 100001
	}
	if err := configuration.SaveShortCode(shortCode); nil != err {
		log.Fatal("Error in updating the shortcode to", strconv.FormatInt(shortCode, 10)+".", err)
	}

	if err := os.MkdirAll("./resources/downloads/"+strconv.FormatInt(shortCode, 10), 0777); nil != err {
		log.Fatal("Error in creating the folder", strconv.FormatInt(shortCode, 10)+".", err)
	}
	log.Println("Updated the shortcode in conf file to", strconv.FormatInt(shortCode, 10))
	log.Println("User requested file upload.")
	uploadFiles := configuration.UploadDetails{}
	log.Println("Extracting the filenames from the form.")
	r.ParseMultipartForm(32 << 20)
	for _, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			uploadFiles.UploadFileNames = append(uploadFiles.UploadFileNames, value)
			log.Println(value)
		}
	}
	log.Println("Extracting the files from the form.")
	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {
			// open uploaded
			if infile, err := hdr.Open(); nil != err {
				log.Fatal("Error in opening a file from form.", err)
			} else {
				uploadFiles.UploadFiles = append(uploadFiles.UploadFiles, infile)
			}
		}
	}

	for i := 0; i < len(uploadFiles.UploadFiles); i++ {
		// open destination
		if outfile, err := os.Create("./resources/downloads/" + strconv.FormatInt(shortCode, 10) + "/" + uploadFiles.UploadFileNames[i]); nil != err {
			log.Fatal("Error in creating a file in system.", err)
		} else {
			if _, err := io.Copy(outfile, uploadFiles.UploadFiles[i]); nil != err {
				log.Fatal("Error in writing a file in system.", err)
			}
		}
	}
	log.Println("Successfully stored all the files into the folder ./resources/downloads/" + strconv.FormatInt(shortCode, 10))
	log.Println("Compressing the folder ./resources/downloads/" + strconv.FormatInt(shortCode, 10))
	if err := configuration.Zipit("./resources/downloads/"+strconv.FormatInt(shortCode, 10), "./resources/downloads/"+strconv.FormatInt(shortCode, 10)+".zip", false); nil != err {
		log.Fatal("Error in compressing file.", err)
	}
	os.RemoveAll("./resources/downloads/" + strconv.FormatInt(shortCode, 10))
	log.Println("Successfully compressed the folder and the uploaded files are available at ./resources/downloads/" + strconv.FormatInt(shortCode, 10) + ".zip")
	ips.IPDetailsConfig = configuration.GetIPDetailsConfig()
	for j := 0; j < len(ips.IPDetailsConfig); j++ {
		if err := configuration.QRCodeGenerate(ips.IPDetailsConfig[j].IPAddr+"/download?downloadCode="+strconv.FormatInt(shortCode, 10), "resources/downloads/"+ips.IPDetailsConfig[j].NWInterface+"_"+strconv.FormatInt(shortCode, 10)+".png"); nil != err {
			log.Fatal("Error in generating QR Code for", ips.IPDetailsConfig[j].IPAddr+".", err)
		}
		ips.IPDetailsConfig[j].IPAddr = "resources/downloads/" + ips.IPDetailsConfig[j].NWInterface + "_" + strconv.FormatInt(shortCode, 10) + ".png"
	}
	ips.ShortCode = strconv.FormatInt(shortCode, 10)

	//Generating page from templates
	t := template.Must(template.New("header").ParseFiles("templates/snippets/homeheader.html"))
	t.Execute(w, nil)
	t = template.Must(template.New("content").ParseFiles("templates/snippets/uploaded.html"))
	t.Execute(w, ips)
	t = template.Must(template.New("footer").ParseFiles("templates/snippets/homefooter.html"))
	t.Execute(w, nil)

}
