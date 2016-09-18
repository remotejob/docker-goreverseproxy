package handlers

import (
	// "bytes"

	"encoding/xml"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/remotejob/docker-goreverseproxy/dbhandler"
	"github.com/remotejob/kaukotyoeu/domains"
	mgo "gopkg.in/mgo.v2"
)

var mongodbuser string
var mongodbpass string
var themes string
var locale string
var mainroute string

var resultXML []byte

func init() {

	log.Println("Sitemap init user", os.Getenv("SECRET_USERNAME"))
	mongodbuser = os.Getenv("SECRET_USERNAME")
	mongodbpass = os.Getenv("SECRET_PASSWORD")
	themes = "job"
	locale = "fi_FI"
	mainroute = "blogi"

}

//CheckServeSitemap CheckServeSitemap
func CheckServeSitemap(w http.ResponseWriter, r *http.Request) {

	sitefull := r.Host
	site := strings.Split(sitefull, ":")[0]

	if site == "localhost" {

		site = "www.kaukotyo.eu"

	}
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:     []string{"mymongo-controller"},
		Timeout:   60 * time.Second,
		Database:  "admin",
		Username:  mongodbuser,
		Password:  mongodbpass,
		Mechanism: "SCRAM-SHA-1",
	}

	dbsession, err := mgo.DialWithInfo(mongoDBDialInfo)

	if err != nil {
		panic(err)
	}
	defer dbsession.Close()

	allsitemaplinks := dbhandler.GetAllSitemaplinks(*dbsession, site)

	// fmt.Println("CheckServeSitemap")

	docList := new(domains.Pages)
	docList.XmlNS = "http://www.sitemaps.org/schemas/sitemap/0.9"

	for _, sitemaplink := range allsitemaplinks {

		if sitemaplink.Site == site {

			doc := new(domains.Page)
			doc.Loc = "http://" + site + "/" + themes + "/" + locale + "/" + mainroute + "/" + sitemaplink.Stitle + ".html"
			doc.Lastmod = sitemaplink.Updated.Format(time.RFC3339)
			doc.Changefreq = "monthly"
			docList.Pages = append(docList.Pages, doc)
			// fmt.Println(site, sitemaplink.Stitle)
		}

	}

	resultXML, err = xml.MarshalIndent(docList, "", "  ")
	if err != nil {
		log.Println(err.Error())
	}

	// fmt.Println(bytes.NewBuffer(resultXML).String())

	// fmt.Println(bytes.NewBuffer(sitemap).String())

	w.Header().Add("Content-type", "application/xml")
	w.Write(resultXML)

}
