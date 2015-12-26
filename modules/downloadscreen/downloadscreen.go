//Package downloadscreen contains the functions required for the download screen
package downloadscreen

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/melvinodsa/fileUpload/configuration"
)

//DownloadScreen generate the homepage of the aplication.
func DownloadScreen(w http.ResponseWriter, r *http.Request) {
	log.Println("User requested download view.")

	//Generating page from templates
	t := template.Must(template.New("header").ParseFiles("templates/snippets/homeheader.html"))
	t.Execute(w, nil)
	t = template.Must(template.New("content").ParseFiles("templates/snippets/download.html"))
	t.Execute(w, nil)
	t = template.Must(template.New("footer").ParseFiles("templates/snippets/homefooter.html"))
	t.Execute(w, nil)
}

//Download generate the homepage of the aplication.
func Download(w http.ResponseWriter, r *http.Request) {
	log.Println("User requested file download.")
	r.ParseForm()
	download := configuration.DownloadDetails{DownloadCode: ""}
	download.DownloadCode = template.HTMLEscapeString(r.Form.Get("downloadCode"))
	log.Println("User requested the file with shortcode", download.DownloadCode+".")
	flag := false
	files, err := ioutil.ReadDir("./resources/downloads")
	if err != nil {
		log.Fatal("Error while listing the downloads directory ./resources/downloads.", err)
	} else {
		log.Println("Files in the directory ./resources/downloads")
		for _, f := range files {
			log.Println(f.Name())
			if download.DownloadCode+".zip" == f.Name() {
				log.Println("Requested file shortcode", download.DownloadCode, "has found a match.")
				flag = true
			}
		}
	}
	if flag {
		//Generating page from templates
		t := template.Must(template.New("header").ParseFiles("templates/snippets/homeheader.html"))
		t.Execute(w, nil)
		t = template.Must(template.New("content").ParseFiles("templates/snippets/downloadfile.html"))
		t.Execute(w, download.DownloadCode+".zip")
		t = template.Must(template.New("footer").ParseFiles("templates/snippets/homefooter.html"))
		t.Execute(w, nil)
	} else {
		log.Println("Requested file shortcode", download.DownloadCode, "has no match in the ./resources/downloads.")
		//Generating page from templates
		t := template.Must(template.New("header").ParseFiles("templates/snippets/homeheader.html"))
		t.Execute(w, nil)
		t = template.Must(template.New("content").ParseFiles("templates/snippets/download.html"))
		t.Execute(w, "Sorry!!! Requested shortcode is not available now")
		t = template.Must(template.New("footer").ParseFiles("templates/snippets/homefooter.html"))
		t.Execute(w, nil)
	}

}
