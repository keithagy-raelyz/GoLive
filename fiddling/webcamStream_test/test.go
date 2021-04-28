package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"
	"net/http"

	"github.com/lazywei/go-opencv/opencv"
)

func broadcast(w http.ResponseWriter, r *http.Request) {
	webCamera := opencv.NewCameraCapture(0)

	if webCamera == nil {
		panic("Unable to open camera")
	}

	defer webCamera.Release()

	for {
		if webCamera.GrabFrame() {
			imgFrame := webCamera.RetrieveFrame(1)
			if imgFrame != nil {
				//fmt.Println(imgFrame.ImageSize())
				//fmt.Println(imgFrame.ToImage())

				// convert IplImage(Intel Image Processing Library)
				// to image.Image
				goImgFrame := imgFrame.ToImage()

				// and then convert to []byte
				// with the help of png.Encode() function

				frameBuffer := new(bytes.Buffer)
				//frameBuffer := make([]byte, imgFrame.ImageSize())
				err := png.Encode(frameBuffer, goImgFrame)

				if err != nil {
					panic(err)
				}

				// convert the buffer bytes to base64 string - use buf.Bytes() for new image
				imgBase64Str := base64.StdEncoding.EncodeToString(frameBuffer.Bytes())

				// Embed into an html without PNG file
				img2html := "<html><body><img src=\"data:image/png;base64," + imgBase64Str + "\" /></body></html>"

				w.Write([]byte(fmt.Sprintf(img2html)))

				// TODO :
				// encode frames to stream via WebRTC

				fmt.Println("Streaming....")

			}
		}
	}

}

func main() {
	fmt.Println("Broadcasting...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", broadcast)
	http.ListenAndServe(":8080", mux)
}
