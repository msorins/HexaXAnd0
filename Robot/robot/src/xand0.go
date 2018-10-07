package Xand0

import (
	"mind/core/framework/skill"
)

type Xand0 struct {
	skill.Base
}

func NewSkill() skill.Interface {
	
	// Use this method to create a new skill.

	return &Xand0{}
}

func (d *Xand0) OnStart() {
	// Use this method to do something when this skill is starting.
}

func (d *Xand0) OnClose() {
	// Use this method to do something when this skill is closing.
}

func (d *Xand0) OnConnect() {
	// Use this method to do something when the remote connected.
}

func (d *Xand0) OnDisconnect() {
	// Use this method to do something when the remote disconnected.
}

func (d *Xand0) OnRecvJSON(data []byte) {
	// Use this method to do something when skill receive json data from remote client.
}

func (d *Xand0) OnRecvString(data string) {
	// Use this method to do something when skill receive string from remote client.
}
