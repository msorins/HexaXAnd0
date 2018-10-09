package main

import (
	"DetectPKG"
	"fmt"
	"gocv.io/x/gocv"
	"strconv"
)

func computeFromFile(path string) (gocv.Mat, [3][3]int) {
	img := gocv.NewMat()
	img = gocv.IMRead(path, gocv.IMReadColor)  // To read from file

	var dtct DetectPKG.IDetect
	dtct = DetectPKG.DetectBuilder(img)
	return dtct.Detect()
}

func main() {
	//// File mode
	////img, board := computeFromFile("/Users/so/Desktop/XAnd0/GameDetectionOpenCV-GO/workingimg41.png")
	////img, board := computeFromFile("/Users/so/Desktop/XAnd0/GameDetectionOpenCV-GO/workingimgv27.png")
	////img, board := computeFromFile("/Users/so/Desktop/XAnd0/GameDetectionOpenCV-GO/xand0img0.png")
	//
	//// Show the board
	//fmt.Println(board)
	//
	//// Show the image
	//finalImageWindow := gocv.NewWindow("FinalImage")
	//finalImageWindow.IMShow(img)
	//finalImageWindow.WaitKey(200000)



	//Webcam mode
	webcam, _ := gocv.VideoCaptureDevice(0)
	img := gocv.NewMat()
	counter := 0
	for {
		//Read from webcam
		webcam.Read(&img)
		gocv.IMWrite("workingimgv2" + strconv.Itoa(counter) + ".png", img)
		counter++

		//Process the image
		var dtct DetectPKG.IDetect
		dtct = DetectPKG.DetectBuilder(img)
		img, board := dtct.Detect()

		// Show the board
		fmt.Println(board)

	//	// Show the image
		finalImageWindow := gocv.NewWindow("FinalImage")
		finalImageWindow.IMShow(img)
		finalImageWindow.WaitKey(2000)
	}
}

//img = gocv.IMRead("/Users/so/Desktop/HexaXand0/GameDetectionOpenCV-GO/img.jpeg", gocv.IMReadColor)  // To read from file