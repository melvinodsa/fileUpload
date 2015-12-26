//Package homescreen contains the functions required for the home screen
package homescreen

import (
	"html/template"
	"log"
	"net"
	"net/http"

	"github.com/melvinodsa/fileUpload/configuration"
)

//Home generate the homepage of the aplication.
func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("User requested home view.")
	log.Println("Getting the network interface.")
	ips := []configuration.IPDetails{}
	ifaces, err := net.Interfaces()
	if err == nil {
		for _, i := range ifaces {
			log.Println("Getting the local ip address from the network interface", i.Name)
			addrs, err := i.Addrs()
			if err == nil {
				if len(addrs) > 0 {
					var ip net.IP
					switch v := addrs[0].(type) {
					case *net.IPNet:
						ip = v.IP
					case *net.IPAddr:
						ip = v.IP
					}
					log.Println("ip address of the network interface", i.Name, "is", ip)
					ips = append(ips, configuration.IPDetails{IPAddr: "http://" + ip.String() + ":9090", NWInterface: i.Name})

					if err := configuration.QRCodeGenerate("http://"+ip.String()+":9090", "resources/images/"+i.Name+".png"); nil != err {
						log.Fatal("Error in generating QR Code for", ip.String()+".", err)
					}

				}
			} else {
				log.Fatal("Error in getting the local ip addresses from the interface.", err)
			}
		}
	} else {
		log.Fatal("Error in getting the network interface.", err)
	}
	if err := configuration.SaveIPDetailsConfig(ips); nil != err {
		log.Fatal("Error in saving ip details.", err)
	}

	//Generating page from templates
	t := template.Must(template.New("header").ParseFiles("templates/snippets/homeheader.html"))
	t.Execute(w, nil)
	t = template.Must(template.New("content").ParseFiles("templates/snippets/home.html"))
	t.Execute(w, ips)
	t = template.Must(template.New("footer").ParseFiles("templates/snippets/homefooter.html"))
	t.Execute(w, nil)
}
