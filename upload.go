package adele

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gabriel-vasile/mimetype"
	"github.com/harrisonde/adele/filesystem"
)

func (a *Adele) UploadFile(r *http.Request, destination, field string, fs filesystem.FS) error {
	fileName, err := a.getFileToUpload(r, field)

	if err != nil {
		a.ErrorLog.Println(err)
		return err
	}

	if fs != nil {
		err = fs.Put(fileName, destination)
		if err != nil {
			a.ErrorLog.Println(err)
			return err
		}
	} else {
		err = os.Rename(fileName, fmt.Sprintf("%s/%s", destination, path.Base(fileName)))
		if err != nil {
			a.ErrorLog.Println(err)
		}
		return err
	}

	defer func() {
		_ = os.Remove(fileName)
	}()

	return nil
}

func (a *Adele) getFileToUpload(r *http.Request, fieldName string) (string, error) {
	_ = r.ParseMultipartForm(a.config.uploads.maxUploadSize)

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}

	defer file.Close()

	mimeType, err := mimetype.DetectReader(file)
	if err != nil {
		return "", err
	}

	// Mimetype put pointier to a new position, we need
	// to back to zero
	_, err = file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	if !inSlice(a.config.uploads.allowedMimeTypes, mimeType.String()) {
		return "", errors.New("invalid file type uploaded")
	}

	dst, err := os.Create(fmt.Sprintf("./tmp/%s", header.Filename))
	if err != nil {
		return "", err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("./tmp/%s", header.Filename), nil
}

func inSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
