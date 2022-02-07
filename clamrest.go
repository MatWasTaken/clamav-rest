package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"clamav-rest/go-clamd"
)

var opts map[string]string
var positiveCounter int = 0

type Result struct {
	Status, Description string
}

func init() {
	log.SetOutput(ioutil.Discard)
}

func home(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, "...running...")
}

//this is where the upload happens.
func UploadHandler(res http.ResponseWriter, req *http.Request) {
	positiveCounter++
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(res, err.Error(), status)
		}
	}()
	// parse request
	// const _24K = (1 << 20) * 24
	if err = req.ParseMultipartForm(32 << 20); nil != err {
		status = http.StatusInternalServerError
		return
	}
	fmt.Println("No memory problem")
	for _, fheaders := range req.MultipartForm.File {
		for _, hdr := range fheaders {
			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); nil != err {
				status = http.StatusInternalServerError
				return
			}

			// create destination
			err = os.MkdirAll("/home/web.app/data.clamav/quarantine/"+time.Now().Format("2006-01-02")+"/", os.ModePerm)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}

			// open destination
			var outfile *os.File
			if outfile, err = os.Create("/home/web.app/data.clamav/quarantine/" + time.Now().Format("2006-01-02") + "/" + time.Now().Format("15h04") + "-" + path.Base(hdr.Filename)); nil != err {
				status = http.StatusInternalServerError
				return
			}
			// 32K buffer copy
			var written int64
			if written, err = io.Copy(outfile, infile); nil != err {
				status = http.StatusInternalServerError
				return
			}
			res.Write([]byte("uploaded file:" + path.Base(hdr.Filename) + ";length:" + strconv.Itoa(int(written)) + "\n"))
			fmt.Printf(time.Now().Format(time.RFC3339) + " Finished uploading: " + hdr.Filename + " in quarantine" + "\n")
		}
	}
}

//This is where the scan happens.
func scanHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//POST takes the uploaded file(s) and saves it to disk.
	case "POST":
		c := clamd.NewClamd(opts["CLAMD_PORT"])
		//get the multipart reader for the request.
		reader, err := r.MultipartReader()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//copy each part to destination.
		for {
			part, err := reader.NextPart()
			if err == io.EOF {
				break
			}

			//if part.FileName() is empty, skip this iteration.
			if part.FileName() == "" {
				continue
			}

			fmt.Printf(time.Now().Format(time.RFC3339) + " Started scanning: " + part.FileName() + "\n")
			var abort chan bool
			response, err := c.ScanStream(part, abort)
			for s := range response {
				w.Header().Set("Content-Type", "application/json; charset=utf-8")
				// respJson := fmt.Sprintf("{ \"Status\": \"%s\", \"Description\": \"%s\" }", s.Status, s.Description)
				// respJson := fmt.Sprintf("{ \"Statut\": ", s.Status, "; \"Description\": ", s.Description, "}")
				respJson := Result{
					Status:      s.Status,
					Description: s.Description,
				}
				jsonData, err := json.Marshal(respJson)

				switch s.Status {
				case clamd.RES_OK:
					w.WriteHeader(http.StatusOK)
				case clamd.RES_FOUND:
					w.WriteHeader(http.StatusNotAcceptable)
				case clamd.RES_ERROR:
					w.WriteHeader(http.StatusBadRequest)
				case clamd.RES_PARSE_ERROR:
					w.WriteHeader(http.StatusPreconditionFailed)
				default:
					w.WriteHeader(http.StatusNotImplemented)
				}
				fmt.Fprint(w, string(jsonData)+"\n")
				fmt.Printf("\n")
				if err != nil {
					log.Println(err)
				}
				// fmt.Fprint(w, respJson)
				fmt.Printf(time.Now().Format(time.RFC3339)+" Scan result for: %v, %v\n", part.FileName(), s)
			}
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf(time.Now().Format(time.RFC3339) + " Finished scanning: " + part.FileName() + "\n")
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func waitForClamD(port string, times int) {
	clamdTest := clamd.NewClamd(port)
	clamdTest.Ping()
	version, err := clamdTest.Version()

	if err != nil {
		if times < 30 {
			fmt.Printf("clamD not running, waiting times [%v]\n", times)
			time.Sleep(time.Second * 4)
			waitForClamD(port, times+1)
		} else {
			fmt.Printf("Error getting clamd version: %v\n", err)
			os.Exit(1)
		}
	} else {
		for version_string := range version {
			fmt.Printf("Clamd version: %#v\n", version_string.Raw)
		}
	}
}

func main() {

	const (
		PORT     = ":9000"
		SSL_PORT = ":9443"
	)

	opts = make(map[string]string)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		opts[pair[0]] = pair[1]
	}

	if opts["CLAMD_PORT"] == "" {
		opts["CLAMD_PORT"] = "tcp://127.0.0.1:3310"
	}

	fmt.Printf("Starting clamav rest bridge\n")
	fmt.Printf("Connecting to clamd on %v\n", opts["CLAMD_PORT"])
	waitForClamD(opts["CLAMD_PORT"], 1)

	fmt.Printf("Connected to clamd on %v\n", opts["CLAMD_PORT"])

	http.HandleFunc("/scan", scanHandler)
	http.HandleFunc("/upload", UploadHandler)
	http.HandleFunc("/", home)

	//Listen on port PORT
	if opts["PORT"] == "" {
		opts["PORT"] = "9000"
	}
	fmt.Printf("Listening on port " + opts["PORT"] + "\n")
	//http.ListenAndServe(":"+opts["PORT"], nil)

	l, err := net.Listen("tcp4", ":"+opts["PORT"])
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.Serve(l, nil))

	// Start the HTTPS server in a goroutine
	go http.ListenAndServeTLS(SSL_PORT, "/etc/ssl/clamav-rest/atgpedi.net.cer", "/etc/ssl/clamav-rest/atgpedi.net.key", nil)

	// Start the HTTP server
	http.ListenAndServe(PORT, nil)
}
