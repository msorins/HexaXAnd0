package DetectPKG

import (
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"math"
	"sort"
)

type Detect struct {
	img gocv.Mat
	imgGray gocv.Mat
	contours [][]image.Point
	rootContour []image.Point
	board [3][3]int
}

func (this *Detect) filterFirstTwoContourLevels(contours [][]image.Point) [][]image.Point {
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

func (this *Detect) filterContoursThatDoNotHaveNineKids(contours [][]image.Point) ([][]image.Point, []image.Point) {
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

func (this *Detect) computeBoardState() [3][3]int {
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

	// this.contours, this.rootContour, this.img, this.imgGray

	var rects []gocv.RotatedRect

	if len(this.contours) != 9 {
		return [3][3]int{}
	}

	// Form Min Area Rectangles from contours
	for i := 0; i < len(this.contours); i++ {
		rects = append(rects, gocv.MinAreaRect(this.contours[i]))
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

	this.img.ConvertTo(&this.img, gocv.MatTypeCV8S)
	this.imgGray.ConvertTo(&this.imgGray, gocv.MatTypeCV8S)
	used := [2000][2000]bool{}

	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			dx := []int{-1,1,0,0}
			dy := []int{0,0,-1,1}

			q := []image.Point{ image.Point{rects[index].Center.Y, rects[index].Center.X } }
			for len(q) > 0 {
				front := q[0]
				q = q[1:]

				b := float64(this.img.GetVeciAt(front.X, front.Y)[0])
				//g := float64(this.img.GetVeciAt(front.X, front.Y)[1])
				r := float64(this.img.GetVeciAt(front.X, front.Y)[2])

				if( this.imgGray.GetUCharAt(front.X, front.Y) != 0 ) {
					if r >= 120 {
						game[i][j] = 1
						break
					}
					if b >= 120 {
						game[i][j] = 2
						break
					}
				}

				if this.imgGray.GetUCharAt(front.X, front.Y) == 0 {
					for h := 0; h < 4; h++ {
						next := image.Point{front.X + dx[h], front.Y + dy[h]}

						// Check if next point is NOT in current rectangle
						if next.X < 0 || next.Y < 0 || next.X >= this.img.Size()[0] || next.Y >= this.img.Size()[1] || used[next.X][next.Y] == true {
							continue
						}

						// Add it to queue
						q = append(q, next)
						used[next.X][next.Y] = true
					}
				}
			}

			index += 1
		}
	}

	return game
}

func (this *Detect) processImage() (gocv.Mat, [3][3]int) {
	// Declare visualisation windows
	threshold := gocv.NewWindow("Threshold")

	// Invert pixels
	inverted := gocv.NewMat()
	gocv.BitwiseNot(this.img, &inverted)

	// Apply a threshold
	gocv.CvtColor(inverted, &this.imgGray, gocv.ColorRGBToGray)
	gocv.AdaptiveThreshold(this.imgGray, &this.imgGray, 10, gocv.AdaptiveThresholdMean, gocv.ThresholdBinary, 75, -40)
	for i := 0; i < 720; i++ {
		for j := 0; j < 1280; j++ {
			if this.imgGray.GetSCharAt(i,j) == 10 {
				this.imgGray.SetSCharAt(i,j, 100)
			}
		}
	}

	// Show the threshold window
	threshold.IMShow(this.imgGray)
	threshold.WaitKey(100)

	// Get the contours
	contours := gocv.FindContours(this.imgGray, gocv.RetrievalTree, gocv.ContourApproximationMode(gocv.ChainApproxSimple))
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
		return this.img, [3][3]int{}
	}

	// Filter all the contours that are very nested
	this.contours  = this.filterFirstTwoContourLevels(approxedContours)
	this.contours, this.rootContour = this.filterContoursThatDoNotHaveNineKids(this.contours)

	// Get the parsed baoard
	this.board = this.computeBoardState()

	// Draw the result contours on the img
	if len(approxedContours) != 0 {
		// Draw root contour
		if len(this.rootContour) != 0 {
			rc := [][]image.Point{this.rootContour}
			gocv.DrawContours(&this.img, rc, 0, color.RGBA{255, 0, 0, 100}, 3)
		}

		// Draw the other contours
		gocv.DrawContours(&this.img, approxedContours, -1, color.RGBA{0, 0, 255, 100}, 2)
	}

	return this.img, this.board
}

func DetectBuilder(img gocv.Mat) IDetect {
	dtct := Detect{}
	dtct.img = img
	dtct.imgGray = gocv.NewMat()

	return &dtct
}

func (this *Detect) Detect() (gocv.Mat, [3][3]int) {
	return this.processImage()
}


