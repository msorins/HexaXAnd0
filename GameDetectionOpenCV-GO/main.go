package main

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"io"
	"math"
	"os"
	"sort"
	"strings"
)

func WriteStringToFile(filepath, s string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(s))
	if err != nil {
		return err
	}

	return nil
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

func filterContoursThatDoNotHaveNineKids(contours [][]image.Point) ([][]image.Point, []image.Point) {
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
	rootContour  := []image.Point{}
	for i := 0; i < len(contours); i++ {
		if i == theOne {
			rootContour = contours[i]
		} else {
			if depth[i] < 2 {
				finalContours = append(finalContours, contours[i])
			}
		}
	}

	return finalContours, rootContour
}

func getBoardState(contours [][]image.Point, rootContour []image.Point, img gocv.Mat, gray gocv.Mat) [3][3]int {
	//s := ""
	//for i := 0; i < 720; i++ {
	//	for j := 0; j < 1280; j++ {
	//		if( int(gray.GetUCharAt(i, j)) != 0) {
	//			s += "1"
	//		} else {
	//			s += "0"
	//		}
	//	}
	//	s += "\n"
	//}
	//WriteStringToFile("/Users/so/Desktop/XAnd0/workingimg0.txt", s)

	var rects []gocv.RotatedRect

	if len(contours) != 9 {
		return [3][3]int{}
	}

	// Form Min Area Rectangles from contours
	for i := 0; i < len(contours); i++ {
		rects = append(rects, gocv.MinAreaRect(contours[i]))
	}


	// Sort the Min Area Rectangles from contours
	sort.Slice(rects, func(i int, j int) bool {
		if math.Abs( float64(rects[i].Center.X) - float64(rects[j].Center.X) ) <= 100 {
			return rects[i].Center.Y <= rects[j].Center.Y
		}

		return rects[i].Center.X <= rects[j].Center.X
	});

	//Show drawing animation
	//finalImageWindow := gocv.NewWindow("DrawingPortionsImage")
	//for i := 0; i < len(rects); i++ {
	//	gocv.DrawContours(&img, [][]image.Point{rects[i].Contour}, 0, color.RGBA{255, 0, 0, 100}, 1)
	//	finalImageWindow.IMShow(img)
	//	finalImageWindow.WaitKey(1500)
	//}

	// Check to see the actual state of the game
	game := [3][3]int { {0, 0, 0}, {0, 0, 0}, {0, 0, 0} }
	game[0][0] = 1
	index := 0

	img.ConvertTo(&img, gocv.MatTypeCV8S)
	gray.ConvertTo(&gray, gocv.MatTypeCV8S)
	used := [2000][2000]bool{}

	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			dx := []int{-1,1,0,0}
			dy := []int{0,0,-1,1}

			q := []image.Point{ image.Point{rects[index].Center.Y, rects[index].Center.X } }
			chosen := false
			for len(q) > 0 {
				front := q[0]
				q = q[1:]

				b := float64(img.GetVeciAt(front.X, front.Y)[0])
				g := float64(img.GetVeciAt(front.X, front.Y)[1])
				r := float64(img.GetVeciAt(front.X, front.Y)[2])
				//fmt.Println(front.X, ":", front.Y)
				//fmt.Println(r, "-",g,"-",b)

				if( gray.GetUCharAt(front.X, front.Y) != 0 ) {
					if r >= 120 {
						fmt.Println("RED ", r, "->", g, "->", b)
						game[i][j] = 1
						chosen = true
						break
					}
					if b >= 120 {
						fmt.Println("BLUE",  r, "->", g, "->", b)
						game[i][j] = 2
						chosen = true
						break
					}
				}

				if gray.GetUCharAt(front.X, front.Y) == 0 {
					for h := 0; h < 4; h++ {
						next := image.Point{front.X + dx[h], front.Y + dy[h]}

						// Check if next point is NOT in current rectangle
						if next.X < 0 || next.Y < 0 || next.X >= img.Size()[0] || next.Y >= img.Size()[1] || used[next.X][next.Y] == true {
							continue
						}

						// Add it to queue
						q = append(q, next)
						used[next.X][next.Y] = true
					}
				}
			}
			if !chosen {
				fmt.Println("BLANK")
			}

			index += 1
		}
	}

	return game
}

func processImage(img gocv.Mat) (gocv.Mat, [3][3]int) {
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
		return img, [3][3]int{}
	}

	// Filter all the contours that are very nested
	var rootContour []image.Point
	approxedContours  = filterFirstTwoContourLevels(approxedContours)
	approxedContours, rootContour = filterContoursThatDoNotHaveNineKids(approxedContours)

	// Get the parsed baoard
	board := getBoardState(approxedContours, rootContour, img, gray)

	// Draw the result contours on the img
	if len(approxedContours) != 0 {
		// Draw root contour
		if len(rootContour) != 0 {
			rc := [][]image.Point{rootContour}
			gocv.DrawContours(&img, rc, 0, color.RGBA{255, 0, 0, 100}, 3)
		}

		// Draw the other contours
		gocv.DrawContours(&img, approxedContours, -1, color.RGBA{0, 0, 255, 100}, 2)
	}

	return img, board
}

func main() {
	//webcam, _ := gocv.VideoCaptureDevice(0)
	finalImageWindow := gocv.NewWindow("FinalImage")
	img := gocv.NewMat()

	var board [3][3]int
	img = gocv.IMRead("/Users/so/Desktop/XAnd0/GameDetectionOpenCV-GO/xand0img0.png", gocv.IMReadColor)  // To read from file
	img, board = processImage(img)

	// Print the board
	fmt.Println(board)

	// Show the image
	finalImageWindow.IMShow(img)
	finalImageWindow.WaitKey(200000)

	//counter := 0
	//for {
	//	//Read from webcam
	//	webcam.Read(&img)
	//	gocv.IMWrite("workingimg" + strconv.Itoa(counter) + ".png", img)
	//	counter++
	//
	//	//Process the image
	//	img = processImage(img)
	//
	//
	//	//Show the result
	//	finalImageWindow.IMShow(img)
	//	finalImageWindow.WaitKey(10000)
	//	//break
	//	}
}

//img = gocv.IMRead("/Users/so/Desktop/HexaXand0/GameDetectionOpenCV-GO/img.jpeg", gocv.IMReadColor)  // To read from file