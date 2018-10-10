package Xand0

import (
	"io"
	"net/http"
	"time"
	"mind/core/framework/skill"
	"mind/core/framework/log"

	"mind/core/framework/drivers/distance"
	"mind/core/framework/drivers/hexabody"

	"github.com/gocv.io/x/gocv"
	"github.com/gorilla/mux"
)

const (
	DISTANCE_TO_WALL   = 100 // millimeters
	MOVE_HEAD_DURATION = 500 // milliseconds
	WALK_SPEED         = 0.3 // cm per second
	SENSE_INTERVAL     = 250 // four times per second
)

type Xand0 struct {
	skill.Base
}

func NewSkill() skill.Interface {
	// Use this method to create a new skill.
	gocv.NewMat()
	return &Xand0{}
}

func (d *Xand0) OnStart() {
	// Use this method to do something when this skill is starting.
	log.Debug.Println("OnStart()")

	// Start hexabody
	err := hexabody.Start()
	if err != nil {
		log.Error.Println("Hexabody start err:", err)
		return
	}

	// Start the distance
	err = distance.Start()
	if err != nil {
		log.Error.Println("Distance start err:", err)
		return
	}
	if !distance.Available() {
		log.Error.Println("Distance sensor is not available")
	}

	// Start the server API
	go d.StartAPI()

	// Execute the sequence of operations
	d.executeSeqOfOperations()
}

func (d *Xand0) OnClose() {
	// Use this method to do something when this skill is closing.
	log.Debug.Println("OnClose()")
	distance.Close()
	hexabody.Close()
}

func (d *Xand0) OnConnect() {
	// Use this method to do something when the remote connected.
	log.Debug.Println("OnConnect()")
}

func (d *Xand0) OnDisconnect() {
	// Use this method to do something when the remote disconnected.
	log.Debug.Println("OnDisconnect()")
}

func (d *Xand0) OnRecvJSON(data []byte) {
	// Use this method to do something when skill receive json data from remote client.
}

func (d *Xand0) OnRecvString(data string) {
	// Use this method to do something when skill receive string from remote client.
	switch data {
		case "ReExec":
			d.executeSeqOfOperations()
	}
}

func (d *Xand0) executeSeqOfOperations() {
	// Move head to front of robot
	d.moveHeadTo0()

	// Start walking until wall is reached
	d.walk()
}

func (d *Xand0) moveHeadTo0() {
	hexabody.MoveHead(0, MOVE_HEAD_DURATION)
}

func (d *Xand0) getDistance() float64 {
	distanceVal, err := distance.Value()
	if err != nil {
		log.Error.Println(err)
	}

	log.Debug.Println("Distance: ", distanceVal)
	return distanceVal
}

func (d *Xand0) walk() {
	hexabody.WalkContinuously(0, WALK_SPEED)
	log.Debug.Println("walk()")
	for {
		if d.getDistance() <= DISTANCE_TO_WALL {
			hexabody.StopWalkingContinuously()
			break
		}
		time.Sleep(SENSE_INTERVAL * time.Millisecond)
	}
}

func (d *Xand0) StartAPI() {
	r := mux.NewRouter()
	r.HandleFunc("/", d.IndexHandler)
	log.Error.Println(http.ListenAndServe(":8000", r))
}

func (d *Xand0) IndexHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	d.executeSeqOfOperations()

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}
