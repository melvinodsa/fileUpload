//This the main program that will be running.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/melvinodsa/fileUpload/app"
	"github.com/melvinodsa/fileUpload/configuration"
	"github.com/melvinodsa/fileUpload/modules/downloadscreen"
	"github.com/melvinodsa/fileUpload/modules/homescreen"
	"github.com/melvinodsa/fileUpload/modules/uploadscreen"
)

func home(w http.ResponseWriter, r *http.Request) {
	homescreen.Home(w, r)
}

func downloadview(w http.ResponseWriter, r *http.Request) {
	downloadscreen.DownloadScreen(w, r)
}

func download(w http.ResponseWriter, r *http.Request) {
	downloadscreen.Download(w, r)
}

func uploadview(w http.ResponseWriter, r *http.Request) {
	uploadscreen.UploadScreen(w, r)
}

func upload(w http.ResponseWriter, r *http.Request) {
	uploadscreen.Upload(w, r)
}

func main() {
	configuration.CacheConfiguration()
	if err := configuration.SaveShortCode(100000); nil != err {
		log.Fatal("Error in intializing the shortcode to 100000.", err)
	} else {
		//Comment this line to prevent app from starting...
		go app.StartFileUpload("http://localhost:9090")
		f, err := os.OpenFile("./logs/fileUpload.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("error opening log file: ", err)
		}
		defer f.Close()
		log.SetOutput(f)
		log.Println("The application has started.")
		http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
		//Handling all the url requests.
		http.HandleFunc("/", home)
		http.HandleFunc("/downloadscreen", downloadview)
		http.HandleFunc("/download", download)
		http.HandleFunc("/uploadscreen", uploadview)
		http.HandleFunc("/upload", upload)

		err = http.ListenAndServe(":9090", nil) // setting listening port
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}
