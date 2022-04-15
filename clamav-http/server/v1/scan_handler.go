package v1alpha

import (
	"encoding/json"
	"net/http"

	"github.com/dutchcoders/go-clamd"
	log "github.com/sirupsen/logrus"
)

type ScanHandler struct {
	Address      string
	Max_file_mem int64
	Logger       *log.Logger
}

//Scanned file model
type ScanFileResponse struct {
	Name  string `json:"name"`
	Pass  bool   `json:"isPassed"`
	Found string `json:"found"`
}

//Response model
type ScanResponse struct {
	ErrorCode        TypeError          `json:"errorCode"`
	ListFiles        []ScanFileResponse `json:"listFiles"`
	ErrorDescription string             `json:"errorDescription"`
}

type TypeError string

const (
	SCAN_OK               TypeError = "SCAN_OK"
	MULTIPART_PARSE_ERROR TypeError = "MULTIPART_PARSE_ERROR"
	NO_FILE_SENT          TypeError = "NO_FILE_SENT"
	FILE_ERROR            TypeError = "FILE_ERROR"
	SCAN_FILE_ERROR       TypeError = "SCAN_FILE_ERROR"
	SCAN_RESPONSE_ERROR   TypeError = "SCAN_RESPONSE_ERROR"
	SCAN_PARSE_ERROR      TypeError = "SCAN_PARSE_ERROR"
)

//Handle request on api
func (sh *ScanHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//Initialized response
	scanResponse := ScanResponse{
		ErrorCode:        SCAN_OK,
		ListFiles:        make([]ScanFileResponse, 0),
		ErrorDescription: "",
	}

	//Multipart parsing
	err := r.ParseMultipartForm(sh.Max_file_mem * 1024 * 1024)

	//Parsing error
	if err != nil {
		scanResponse.ErrorCode = MULTIPART_PARSE_ERROR
		scanResponse.ErrorDescription = err.Error()
		//Write response
		writeResponse(w, scanResponse, true, http.StatusOK)
		//Log
		sh.Logger.Errorf("scan error %d: %s", MULTIPART_PARSE_ERROR, err.Error())
		return
	}

	//Get file from multipart
	files := r.MultipartForm.File["file"]

	//No File sent
	if len(files) == 0 {
		scanResponse.ErrorCode = NO_FILE_SENT
		//Write response
		writeResponse(w, scanResponse, true, http.StatusOK)
		//Log
		sh.Logger.Errorf("scan error %d: %s", NO_FILE_SENT, "No file sent")
		return
	}

	//Open file
	f, err := files[0].Open()

	//Error on file opening
	if err != nil {
		scanResponse.ErrorCode = FILE_ERROR
		scanResponse.ErrorDescription = err.Error()
		//Write response
		writeResponse(w, scanResponse, true, http.StatusOK)
		//Log
		sh.Logger.Errorf("scan error %d: %s", FILE_ERROR, err.Error())
		return
	}

	defer f.Close()

	//Send file to clamd
	c := clamd.NewClamd(sh.Address)
	response, err := c.ScanStream(f, make(chan bool))

	//Error on send file
	if err != nil {
		scanResponse.ErrorCode = SCAN_FILE_ERROR
		scanResponse.ErrorDescription = err.Error()
		//Write response
		writeResponse(w, scanResponse, true, http.StatusOK)
		//Log
		sh.Logger.Errorf("scan error %d: %s", SCAN_FILE_ERROR, err.Error())
		return
	}

	//Get response from clamd
	result := <-response

	if result.Status == "OK" {
		scanResponse.ListFiles = []ScanFileResponse{{
			Name: files[0].Filename,
			Pass: true,
		}}
		sh.Logger.Infof("Scanning %v: RES_OK", files[0].Filename)
	} else if result.Status == "FOUND" {
		scanResponse.ListFiles = []ScanFileResponse{{
			Name:  files[0].Filename,
			Pass:  false,
			Found: result.Raw,
		}}
		sh.Logger.Infof("Scanning %v: RES_FOUND", files[0].Filename)
	} else if result.Status == "ERROR" {
		scanResponse.ErrorCode = SCAN_RESPONSE_ERROR
		scanResponse.ErrorDescription = result.Raw
		scanResponse.ListFiles = []ScanFileResponse{{
			Name: files[0].Filename,
			Pass: false,
		}}
		sh.Logger.Infof("Scanning %v: RES_ERROR", files[0].Filename)
	} else if result.Status == "PARSE_ERROR" {
		scanResponse.ErrorCode = SCAN_PARSE_ERROR
		scanResponse.ErrorDescription = result.Raw
		scanResponse.ListFiles = []ScanFileResponse{{
			Name: files[0].Filename,
			Pass: false,
		}}
		sh.Logger.Infof("Scanning %v: RES_PARSE_ERROR", files[0].Filename)
	} else {
		//Write response
		writeResponse(w, scanResponse, false, http.StatusNotImplemented)
		return
	}

	//Write response
	writeResponse(w, scanResponse, true, http.StatusOK)

	return
}

//Write response
func writeResponse(w http.ResponseWriter, scanResponse ScanResponse, isParseBody bool, httpStatusCode int) {
	//Verify if sending body is needed
	if isParseBody {
		//Marshall body
		body, err := json.Marshal(scanResponse)

		//Error marshall body
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Infof("Parsing response error: %v", err.Error())
			return
		}

		//if no error send body
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(httpStatusCode)
		w.Write(body)
	} else {
		w.WriteHeader(httpStatusCode)
	}

	return
}
