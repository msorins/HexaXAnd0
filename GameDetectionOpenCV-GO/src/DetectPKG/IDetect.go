package DetectPKG

import "gocv.io/x/gocv"

type IDetect interface {
	Detect() (gocv.Mat, [3][3]int)
}
