package main

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
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

func filterFirstTwoContourLevels(contours [][]image.Point) [][]image.Point {
	var graf [1000][]int
	var visited[1000]bool
	var depth[1000]int
	var rects []gocv.RotatedRect

	// Form Min Area Rectangles from conotours
	for i := 0; i < len(contours); i++ {
		visited[i] = false
		depth[i] = -1
		rects = append(rects, gocv.MinAreaRect(contours[i]))
	}

	// Form the graf & visited
	for i := 0; i < len(contours); i++  {
		for j := 0; j < len(contours); j++ {
			if i == j {
				continue
			}

			if rects[j].BoundingRect.In( rects[i].BoundingRect ) {
				graf[i] = append(graf[i], j)
				visited[j] = true
			}
		}
	}

	// Form the depth vector
	for i := 0; i < len(contours); i++ {
		// If we found a node at root level
		if visited[i] == false {
			depth[i] = 0
			q := []int{i}
			kids := 0

			for len(q) > 0 {
				crt := q[0] // .front
				q = q[1:] // .pop

				for j := 0; j < len(graf[crt]); j += 1 {
					if depth[ graf[crt][j] ] < depth[crt] + 1 {
						depth[ graf[crt][j] ] = depth[crt] + 1
						kids += 1
						q = append(q, graf[crt][j])
					}
				}
			}

			// Also filter contours that do not have any kids
			if kids == 0 {
				depth[i] = 999
			}
		}
	}

	// Filter all the contours with depth >= 2
	finalContours := [][]image.Point{}
	for i := 0; i < len(contours); i++ {
		if depth[i] < 2 {
			finalContours = append(finalContours, contours[i])
		}
	}

	return finalContours
}

func filterContoursThatDoNotHaveNineKids(contours [][]image.Point) ([][]image.Point, int) {
	var graf [1000][]int
	var visited[1000]bool
	var depth[1000]int
	var rects []gocv.RotatedRect
	theOne := -1

	// Form Min Area Rectangles from conotours
	for i := 0; i < len(contours); i++ {
		visited[i] = false
		depth[i] = -1
		rects = append(rects, gocv.MinAreaRect(contours[i]))
	}

	// Form the graf & visited
	for i := 0; i < len(contours); i++  {
		for j := 0; j < len(contours); j++ {
			if i == j {
				continue
			}

			if rects[j].BoundingRect.In( rects[i].BoundingRect ) {
				graf[i] = append(graf[i], j)
				visited[j] = true
			}
		}
	}

	// Form the depth vector
	for i := 0; i < len(contours); i++ {
		// If we found a node at root level
		if visited[i] == false {
			depth[i] = 0
			q := []int{i}
			kids := 0

			for len(q) > 0 {
				crt := q[0] // .front
				q = q[1:]   // .pop

				for j := 0; j < len(graf[crt]); j += 1 {
					if depth[ graf[crt][j] ] < depth[crt]+1 {
						depth[ graf[crt][j] ] = depth[crt] + 1
						kids += 1
						q = append(q, graf[crt][j])
					}
				}
			}

			// Also filter contours that do not have any kids
			if kids != 9 {
				depth[i] = 999
			}
			if kids == 9 {
				theOne = i
			}
		}
	}

	// Filter all the contours with depth >= 2
	finalContours := [][]image.Point{}
	for i := 0; i < len(contours); i++ {
		if depth[i] < 2 {
			finalContours = append(finalContours, contours[i])
		}
		if i == theOne {
			theOne = len(finalContours) - 1
		}
	}

	return finalContours, theOne
}

func processImage(img gocv.Mat) gocv.Mat {
	// Declare visualisation windows
	threshold := gocv.NewWindow("Threshold")

	// Invert pixels
	inverted := gocv.NewMat()
	gocv.BitwiseNot(img, &inverted)

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

	// Show the threshold window
	threshold.IMShow(gray)
	threshold.WaitKey(100)

	// Get the contours
	contours := gocv.FindContours(gray, gocv.RetrievalTree, gocv.ContourApproximationMode(gocv.ChainApproxSimple))
	var approxedContours [][]image.Point

	// Approximate the contours with an e = 3%
	for i := 0; i < len(contours); i++{
		approxedContour := gocv.ApproxPolyDP(contours[i], 50, true)

		// Only contours that have between 4 and 6 edges will remain
		if len(approxedContour) == 4 {
			approxedContours = append(approxedContours, approxedContour)
		}
	}

	// If no contours => bye
	if len(approxedContours) == 0 {
		return img
	}

	// Filter all the contours that are very nested
	var theOne int
	approxedContours  = filterFirstTwoContourLevels(approxedContours)
	approxedContours, theOne = filterContoursThatDoNotHaveNineKids(approxedContours)

	if len(approxedContours) != 0 {
		if theOne != -1 {
			gocv.DrawContours(&img, approxedContours, theOne, color.RGBA{0, 255, 0, 100}, 2)
		}
		gocv.DrawContours(&img, approxedContours, -1, color.RGBA{0, 0, 255, 100}, 1)
	}

	return img
}

func main() {
	webcam, _ := gocv.VideoCaptureDevice(0)
	finalImageWindow := gocv.NewWindow("FinalImage")
	img := gocv.NewMat()

	//img = gocv.IMRead("/Users/so/Desktop/XAnd0/GameDetectionOpenCV-GO/probimg32.png", gocv.IMReadColor)  // To read from file
	//img = processImage(img)
	//finalImageWindow.IMShow(img)
	//finalImageWindow.WaitKey(200000)

	//counter := 0
	for {
		//Read from webcam
		webcam.Read(&img)
		//gocv.IMWrite("probimg" + strconv.Itoa(counter) + ".png", img)
		//counter++

		//Process the image
		img = processImage(img)

		//Show the result
		finalImageWindow.IMShow(img)
		finalImageWindow.WaitKey(200)
		}
}

//img = gocv.IMRead("/Users/so/Desktop/HexaXand0/GameDetectionOpenCV-GO/img.jpeg", gocv.IMReadColor)  // To read from file