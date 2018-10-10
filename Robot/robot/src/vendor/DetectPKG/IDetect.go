package DetectPKG

import "github.com/gocv.io/x/gocv"

type IDetect interface {
	Detect() (gocv.Mat, [3][3]int)
}
