package main

import (
	"flag"
	"fmt"
	"github.com/goinggo/tracelog"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/calbucci/go-htmlparser"
)

var stores []string

func extract(filename string) []string {

	buf, err := ioutil.ReadFile(filename)

	s := string(buf)

	if err != nil {
		tracelog.Errorf(err, "getListOfStores", "extractListFromXml", "Error parsing XML's XPath")
	}

	parser := htmlparser.NewParser(s)

	//<a href="/anywhere/redirect?id=318" target="_blank" rel="L.K. Bennett" data-rate="8.5" class="affiliate_link card_container">

	parser.Parse(nil, func(e *htmlparser.HtmlElement, isEmpty bool) {
		if e.TagName == "a" {
			class, _ := e.GetAttributeValue("class")
			if class == "affiliate_link card_container" {
				//stores = append(stores, e.Attributes)
				rel, _ := e.GetAttributeValue("rel")

				stores := append(stores, rel)

				fmt.Println(stores)
			}
		}
	}, nil)
	return stores
}

func download(url, filename string) {
	fmt.Println("Downloading " + url + " ...")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	io.Copy(f, resp.Body)
}

func prepairListOfStores() {
	pUrl := flag.String("url", "https://jet.com/anywhere", "URL to be processed")
	flag.Parse()
	url := *pUrl
	if url == "" {
		tracelog.Errorf(fmt.Errorf("Exception At..."), "getListOfStores", "main", "Error: empty URL!\n")
		return
	}

	filename := path.Base(url) + ".html"
	//tracelog.Info("getListOfStrores", "main", "Checking if " + filename + " exists ...")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		download(url, filename)
		tracelog.Trace("getListOfStores", "main", filename+" saved!")
	} else {
		//tracelog.Trace("getListOfStores", "main", filename + " already exists!")
	}

	extractedResult := extract(filename)

	fmt.Println(extractedResult)
	//tracelog.Info("getListOfStores", "main", "Extracted : " + extractedResult)

	insertAll(extractedResult)
}

func getListOfStoresEmails() []string {
	return getAllStoresEmails()
}
