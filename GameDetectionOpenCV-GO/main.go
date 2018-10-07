package main

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"strconv"
)

func filterByHeightWidthRaport(contour []image.Point) bool {
	rect := gocv.MinAreaRect(contour)

	var rap float64
	if rect.Width > rect.Height {
		rap = float64(rect.Width) / float64(rect.Height)
	} else {
		rap = float64(rect.Height) / float64(rect.Width)
	}

	if rap <= 1.5 {
		return true
	}
	return false
}

func isCountourInContour(external []image.Point, internal []image.Point) bool {
	externalRect := gocv.MinAreaRect(external)
	internalRect := gocv.MinAreaRect(internal)

	return internalRect.BoundingRect.In( externalRect.BoundingRect )
}

func main() {
	webcam, _ := gocv.VideoCaptureDevice(0)
	window := gocv.NewWindow("Hello")
	img := gocv.NewMat()
	window2 := gocv.NewWindow("Window2")

	counter := 0
	for {
		webcam.Read(&img)

		//img = gocv.IMRead("/Users/so/Desktop/HexaXand0/GameDetectionOpenCV-GO/img.jpeg", gocv.IMReadColor)

		// Invert pixels
		inverted := gocv.NewMat()
		gocv.BitwiseNot(img, &inverted)
		gocv.IMWrite("img"+ strconv.Itoa(counter) +".png", img)
		counter += 1

		// Apply a threshold
		gray := gocv.NewMat()
		gocv.CvtColor(inverted, &gray, gocv.ColorRGBToGray)
		gocv.AdaptiveThreshold(gray, &gray, 10, gocv.AdaptiveThresholdMean, gocv.ThresholdBinary, 75, -40)
		for i := 0; i < 720; i++ {
			for j := 0; j < 1280; j++ {
				if gray.GetSCharAt(i,j) == 10 {
					gray.SetSCharAt(i,j, 100)
				}
			}
		}
		window2.IMShow(gray)
		window2.WaitKey(100)

		// Get the contours
		contours := gocv.FindContours(gray, gocv.RetrievalTree, gocv.ContourApproximationMode(gocv.ChainApproxSimple))
		var aproxedContours [][]image.Point

		// Approximate the contours with an e = 3%
		for i := 0; i < len(contours); i++{
			approxedContour := gocv.ApproxPolyDP(contours[i], 30, true)
			if len(approxedContour) >= 4 {
				aproxedContours = append(aproxedContours, approxedContour)
			}

		}

		// Draw the last contour (the most inside one)
		if len(aproxedContours) == 0 {
			continue
		}

		for i := 0; i < len(aproxedContours); i++ {
			overlaps := 0
			for j := 0; j < len(aproxedContours); j++ {
				if j == i {
					continue
				}

				if isCountourInContour(aproxedContours[i], aproxedContours[j]) {
					overlaps += 1
				}
			}

			if overlaps >= 9 {
				gocv.DrawContours(&img, aproxedContours, i, color.RGBA{0, 255, 0, 100}, 2)
			}
		}

		//for i := 0; i < len(aproxedContours); i++  {
		//	gocv.DrawContours(&img, aproxedContours, i, color.RGBA{255, 0, 0, 100}, 1)
		//	window.IMShow(img)
		//	window.WaitKey(200)

		gocv.DrawContours(&img, aproxedContours, -1, color.RGBA{255, 0, 0, 100}, 1)
		window.IMShow(img)
		window.WaitKey(200)
		}



}