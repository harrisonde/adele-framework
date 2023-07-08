package adel

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

func (a *Adel) ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	// TO DO:
	// Move this to the .env and adel app config
	// Limits the max size of the body to be read.
	maxBytes := 1048576 // 1 mb limit for JSON payload
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only have a single json value")
	}

	return nil

}

func (a *Adel) WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (a *Adel) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

func (a *Adel) DownloadFile(w http.ResponseWriter, r *http.Request, pathToFile, fileName string) error {
	fp := path.Join(pathToFile, fileName)

	// clean path up
	fileToServe := filepath.Clean(fp)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)
	return nil
}

func (a *Adel) Error404(w http.ResponseWriter, r *http.Request) {
	a.ErrorStatus(w, http.StatusNotFound)
}

func (a *Adel) Error500(w http.ResponseWriter, r *http.Request) {
	a.ErrorStatus(w, http.StatusInternalServerError)
}

func (a *Adel) ErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	a.ErrorStatus(w, http.StatusUnauthorized)
}

func (a *Adel) ErrorForbidden(w http.ResponseWriter, r *http.Request) {
	a.ErrorStatus(w, http.StatusForbidden)
}

func (a *Adel) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
