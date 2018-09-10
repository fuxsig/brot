// Copyright 2018 Espen Reich. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package brot

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/fuxsig/brot/di"
)

type UploadHandler struct {
	Files     *FileDB    `brot:"files"`
	Parameter string     `brot:"parameter"`
	Data      *DataLayer `brot:"data"`
}

func (th *UploadHandler) HandlerFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if Debug {
			Debug.Printf("UploadHandler begins")
			defer Debug.Printf("UploadHandler ends")
		}
		switch r.Method {
		case "GET":
			values := th.Data.Values(r)
			if key, ok := values.Get("key"); ok {
				fi, err := th.Files.Info(key)
				if err != nil {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				var (
					xMinPar, xMaxPar, yMinPar, yMaxPar, maskPar, srcPar       string
					xMinFlag, xMaxFlag, yMinFlag, yMaxFlag, maskFlag, srcFlag bool
				)
				xMinPar, xMinFlag = values.Get("xMin")
				xMaxPar, xMaxFlag = values.Get("xMax")
				yMinPar, yMinFlag = values.Get("yMin")
				yMaxPar, yMaxFlag = values.Get("yMax")
				maskPar, maskFlag = values.Get("mask")
				srcPar, srcFlag = values.Get("src")

				if maskFlag {
					var mfi *FileInfo
					mfi, err = th.Files.Info(maskPar)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					var mImg image.Image
					mImg, err = th.Files.Image(mfi)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}

					mBounds := mImg.Bounds()
					width := mBounds.Dx()
					height := mBounds.Dy()
					maske := image.NewRGBA(mBounds)
					black := color.RGBA{0, 0, 0, 0}
					white := color.RGBA{0, 0, 0, 255}
					for y := 0; y < height; y++ {
						for x := 0; x < width; x++ {
							red, green, blue, _ := mImg.At(x, y).RGBA()
							if red == green && green == blue && blue > 0 {
								maske.SetRGBA(x, y, white)
							} else {
								maske.SetRGBA(x, y, black)
							}
						}
					}

					var img image.Image
					img, err = th.Files.Image(fi)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					bounds := img.Bounds()

					dest := image.NewRGBA(bounds)
					if srcFlag {

						var sfi *FileInfo
						sfi, err = th.Files.Info(srcPar)
						if err != nil {
							w.WriteHeader(http.StatusNotFound)
							return
						}
						var sImg image.Image
						sImg, err = th.Files.Image(sfi)
						if err != nil {
							w.WriteHeader(http.StatusNotFound)
							return
						}
						draw.Draw(dest, bounds, sImg, image.ZP, draw.Src)
					}

					//draw.Draw(dest, dest.Bounds(), img, image.ZP, draw.Src)
					draw.DrawMask(dest, bounds, img, image.ZP, maske, image.ZP, draw.Over)
					writeImage(w, dest, fi.ContentType)

				} else if xMinFlag || xMaxFlag || yMinFlag || yMaxFlag {
					var img image.Image
					img, err = th.Files.Image(fi)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					bounds := img.Bounds()
					width := bounds.Dx()
					height := bounds.Dy()
					xMin := 0
					if xMinFlag {
						if xMin, err = strconv.Atoi(xMinPar); err != nil || xMin < 0 {
							xMin = 0
						}
					}
					xMax := width
					if xMaxFlag {
						if xMax, err = strconv.Atoi(xMaxPar); err != nil || xMax > width {
							xMax = width
						}
					}
					yMin := 0
					if yMinFlag {
						if yMin, err = strconv.Atoi(yMinPar); err != nil || yMin < 0 {
							yMin = 0
						}
					}
					yMax := height
					if yMaxFlag {
						if yMax, err = strconv.Atoi(yMaxPar); err != nil || yMax > height {
							yMax = height
						}
					}
					cropped := cropImage(xMin, xMax, yMin, yMax, img)
					writeImage(w, cropped, fi.ContentType)
				} else {
					f, err := os.Open(fi.Path)
					if err != nil {
						w.WriteHeader(http.StatusNotFound)
						return
					}
					defer f.Close()

					h := w.Header()
					h.Set("Content-Type", fi.ContentType)
					h.Set("Content-Length", string(fi.Size))
					_, err = io.Copy(w, f)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		case "POST":
			file, handle, err := r.FormFile(th.Parameter)
			if err != nil {
				r.ParseForm()
				if val, ok := r.Form["image"]; ok {
					str := val[0]
					b64data := str[strings.IndexByte(str, ',')+1:]
					reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(b64data))

					if fi, err := th.Files.Store(reader, "camera.png", "image/png", ".png"); err == nil {
						JsonResponse(w, http.StatusCreated, `{"status":"ok","message":"bravo","key":"%s"}`, fi.Key)

					} else {
						InternalError(w, err)
					}
				} else {
					if Debug {
						Debug.Printf("UploadHandler error: %s", err.Error())
					}
					JsonResponse(w, http.StatusBadRequest, `{"status":"error","message":"Reason: %s"}`, err.Error())
					return
				}
				return

			}
			defer file.Close()

			mimeType := handle.Header.Get("Content-Type")
			var extension string
			switch mimeType {
			case "image/jpeg":
				extension = ".jpeg"
			case "image/png":
				extension = ".png"
			case "image/gif":
				extension = ".gif"
			default:
				JsonResponse(w, http.StatusBadRequest, `{"status":"error","message":"unsupported content-type"}`)
				return
			}
			if fi, err := th.Files.Store(file, handle.Filename, mimeType, extension); err == nil {
				JsonResponse(w, http.StatusCreated, `{"status":"ok","message":"bravo","key":"%s"}`, fi.Key)
			} else {
				InternalError(w, err)
			}

		}
	})
}

func writeImage(w http.ResponseWriter, img image.Image, contentType string) {
	var (
		buffer bytes.Buffer
		err    error
	)
	switch contentType {
	case "image/jpeg":
		err = jpeg.Encode(&buffer, img, nil)
	case "image/png":
		err = png.Encode(&buffer, img)
	case "image/gif":
		err = gif.Encode(&buffer, img, nil)
	default:
		err = fmt.Errorf("Unsupported format %s", contentType)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	header := w.Header()
	header.Set("Content-Type", contentType)
	header.Set("Content-Length", string(buffer.Len()))
	io.Copy(w, &buffer)
}

func cropImage(xMin, xMax, yMin, yMax int, img image.Image) image.Image {
	dest := image.NewRGBA(image.Rect(xMin, yMin, xMax, yMax))
	sp := image.Point{xMin, yMin}
	draw.Draw(dest, dest.Bounds(), img, sp, draw.Src)
	return dest
}

func InternalError(w http.ResponseWriter, err error) {
	if Debug {
		Debug.Printf("UploadHandler error: %s", err.Error())
	}
	JsonResponse(w, http.StatusInternalServerError, `{"status":"error","message":"Reason: %s"}`, err.Error())
}

func JsonResponse(w http.ResponseWriter, code int, message string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	fmt.Fprintf(w, message, args...)
}

var _ ProvidesHandler = (*UploadHandler)(nil)
var _ = di.GlobalScope.Declare((*UploadHandler)(nil))
