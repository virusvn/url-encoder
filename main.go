package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"bytes"
	"github.com/gorilla/mux"
	"html/template"
	"strconv"
	 "strings"
	"time"
)

var (
	githash    string
	buildstamp string
)

func main() {
	port := flag.Int("port", 8125, "Port")
	flag.Parse()
	tmpl := template.Must(template.ParseFiles("index.go.html"))

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		vals := r.URL.Query()
		errorMsgs, ok := vals["errorMsg"] // Note type, not ID. ID wasn't specified anywhere.
		var errorMsg string
		if ok {
			if len(errorMsgs) >= 1 {
				errorMsg = errorMsgs[0] // The first `?type=model`
			}
		}
		tmpl.Execute(w, map[string]interface{}{
			"errorMsg": errorMsg,
		})
		return

	})
	fmt.Println("Git Commit Hash: " + githash)
	fmt.Println("UTC Build Time: " + buildstamp)
	r.HandleFunc("/upload", upload)
	fmt.Println("Server serve at: ")
	fmt.Println("http://localhost:" + strconv.Itoa(*port))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), r))
}

func upload(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	if r.Method == http.MethodPost {
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()
		data, err1 := ioutil.ReadAll(file)
		if err1 != nil {
			errorMsg := "Can't get uploaded file"
			http.Redirect(w, r, "/?errorMsg="+errorMsg, 301)
		}
		contentType := http.DetectContentType(data)
		log.Println("contentType", contentType)
		if !strings.Contains(contentType, "text/plain") {
			errorMsg := "File type is not supported. Only 「txt」is supported."
			http.Redirect(w, r, "/?errorMsg="+errorMsg, 301)
			return
		}

		scanner := bufio.NewScanner(bytes.NewReader(data))
		// Not actually needed since it’s a default split function.
		scanner.Split(bufio.ScanLines)

		var ret []string
		for scanner.Scan() {
			line := scanner.Text()
			line = urlencode(line)
			ret = append(ret, line)
		}
		fResult, err := ioutil.TempFile("/tmp", "encoded")
		checkError("Can't create tmp file", err)
		defer os.Remove(fResult.Name())
		
		for _,item := range ret {
			fResult.WriteString(item + "\n")
		}

		// tell the browser the returned content should be downloaded
		t := time.Now()
		name := "encoded.csv"
		w.Header().Set("Content-Disposition", "attachment; filename="+name)
		http.ServeContent(w, r, name, t, fResult)
		elapsed := time.Since(start)
		log.Printf("Took %s", elapsed)
		return
	}
}
func UrlEncoded(str string) (string, error) {
    u, err := url.Parse(str)
    if err != nil {
        return "", err
    }
    return u.String(), nil
}
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}


func urlencode(s string) (result string){
	for _, c := range(s) {
	  if c <= 0x7f { // single byte 
		result += fmt.Sprintf("%%%X", c)
	  } else if c > 0x1fffff {// quaternary byte
		result += fmt.Sprintf("%%%X%%%X%%%X%%%X",
		  0xf0 + ((c & 0x1c0000) >> 18),
		  0x80 + ((c & 0x3f000) >> 12),
		  0x80 + ((c & 0xfc0) >> 6),
		  0x80 + (c & 0x3f),
		)
	  } else if c > 0x7ff { // triple byte
		result += fmt.Sprintf("%%%X%%%X%%%X",
		  0xe0 + ((c & 0xf000) >> 12),
		  0x80 + ((c & 0xfc0) >> 6),
		  0x80 + (c & 0x3f),
		)
	  } else { // double byte
		result += fmt.Sprintf("%%%X%%%X",
		  0xc0 + ((c & 0x7c0) >> 6),
		  0x80 + (c & 0x3f),
		)
	  }
	}
  
	return result
  }