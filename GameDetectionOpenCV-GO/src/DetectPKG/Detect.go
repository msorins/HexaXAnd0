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

func Max(a int, b int) int {
	if a > b {
		return a
	}

	return b
}

func (this *Detect) filterMaxDepthApproach(contours [][]image.Point) ([][]image.Point, []image.Point) {
	var graf [1000][]int
	var grafStrict [1000][]int
	visited := [1000]bool{}
	maxdepth := [1000]int{}
	var rects []image.Rectangle

	// Form Min Area Rectangles from conotours
	for i := 0; i < len(contours); i++ {
		visited[i] = false
		maxdepth[i] = -1
		rects = append(rects, gocv.BoundingRect(contours[i]))
	}

	// Form the graf & visited
	for i := 0; i < len(contours); i++  {
		for j := 0; j < len(contours); j++ {
			if i == j {
				continue
			}

			if rects[j].In( rects[i] ) {
				graf[i] = append(graf[i], j)
				visited[j] = true
			}
		}
	}

	// Form the max depth vector
	for i := 0; i < len(contours); i++ {
		used := [1000]bool{}
		maxdepth[i] = Max(0, maxdepth[i])
		used[i] = true
		q := []int{i}

		for len(q) > 0 {
			crt := q[0] // .front
			q = q[1:] // .pop

			for j := 0; j < len(graf[crt]); j += 1 {
				nxt := graf[crt][j]
				if used[nxt] == false && maxdepth[nxt] < maxdepth[crt] + 1 {
					maxdepth[nxt] = maxdepth[crt] + 1
					used[nxt] = true
				}
			}
		}
	}

	for i := 0; i < len(contours); i++  {
		for j := 0; j < len(graf[i]); j++ {
			if maxdepth[ graf[i][j] ] == maxdepth[i] + 1 {
				grafStrict[i] = append(grafStrict[i], graf[i][j])
			}
		}
	}

	for i := 0; i < len(contours); i++ {
		if len(grafStrict[i]) == 9 {
			miniContours := [][]image.Point{}
			for j := 0; j < 9; j++ {
				miniContours = append(miniContours, contours[ grafStrict[i][j] ])
			}

			return miniContours, contours[ i ]
		}
	}

	return [][]image.Point{}, []image.Point{}
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

	//this.img.ConvertTo(&this.img, gocv.MatTypeCV8S)
	//this.imgGray.ConvertTo(&this.imgGray, gocv.MatTypeCV8S)
	used := [2000][2000]bool{}

	for j := 0; j < 3; j++ {
		for i := 0; i < 3; i++ {
			dx := []int{-1,1,0,0}
			dy := []int{0,0,-1,1}

			game[i][j] = 0
			q := []image.Point{ image.Point{rects[index].Center.Y, rects[index].Center.X } }
			for len(q) > 0 {
				front := q[0]
				q = q[1:]

				b := float64(this.img.GetVeciAt(front.X, front.Y)[0])
				g := float64(this.img.GetVeciAt(front.X, front.Y)[1])
				r := float64(this.img.GetVeciAt(front.X, front.Y)[2])

				if( this.imgGray.GetUCharAt(front.X, front.Y) != 0 ) {
					avgBlack := (math.Abs(b - g) + math.Abs(b - r) + math.Abs(g - r)) / 3
					if avgBlack < 35 {
						continue
					}

					if r >= g && r >= b {
						game[i][j] = 1
						break
					}
					if b >= g && b >= r {
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

	//finalImageWindow := gocv.NewWindow("ABC")
	//for i:= 0; i < len(approxedContours); i++ {
	//	gocv.DrawContours(&this.img, approxedContours, i, color.RGBA{255, 255, 0, 255}, 1,)
	//
	//	rect := gocv.BoundingRect(approxedContours[i])
	//	gocv.DrawKeyPoints(this.img, []gocv.KeyPoint{ gocv.KeyPoint{float64(rect.Min.X), float64(rect.Min.Y), 10, 0, 0, 0, 0}, gocv.KeyPoint{float64(rect.Max.X), float64(rect.Max.Y), 10, 0, 0, 0, 0} }, &this.img, color.RGBA{0, 255, 0, 255}, 0)
	//
	//	finalImageWindow.IMShow(this.img)
	//	finalImageWindow.WaitKey(1500)
	//}


	// Filter all the extra contours (that do not have exactly 9 kids)
	this.contours, this.rootContour = this.filterMaxDepthApproach(approxedContours)

	// Get the parsed baoard
	this.board = this.computeBoardState()

	// Draw the result contours on the img
	if len(this.contours) != 0 {
		// Draw root contour
		if len(this.rootContour) != 0 {
			rc := [][]image.Point{this.rootContour}
			gocv.DrawContours(&this.img, rc, 0, color.RGBA{255, 0, 0, 100}, 3)
		}

		// Draw the other contours
		gocv.DrawContours(&this.img, this.contours, -1, color.RGBA{0, 0, 255, 100}, 2)
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


