package controller

import (
	log "github.com/Sirupsen/logrus"
	"github.com/adammck/hexapod"
	"github.com/adammck/hexapod/math3d"
	"github.com/adammck/sixaxis"
	"io"
	"time"
)

const (
	moveSpeed = 100.0
	rotSpeed  = 20.0
)

type Controller struct {
	sa *sixaxis.SA
}

func New(r io.Reader) *Controller {
	return &Controller{
		sa: sixaxis.New(r),
	}
}

func (c *Controller) Boot() error {
	go c.sa.Run()
	return nil
}

func (c *Controller) Tick(now time.Time, state *hexapod.State) error {

	p := math3d.Pose{
		Position: math3d.Vector3{
			X: (float64(c.sa.LeftStick.X) / 127.0) * moveSpeed,
			Z: (float64(-c.sa.LeftStick.Y) / 127.0) * moveSpeed,
		},
		Heading: (float64(c.sa.RightStick.X) / 127.0) * rotSpeed,
	}

	y := state.Target.Position.Y
	state.Target = state.Pose.Add(p)
	state.Target.Position.Y = y

	// At any time, pressing start shuts down the hex.
	if c.sa.Start && !state.Shutdown {
		log.Warn("Pressed START, shutting down")
		state.Shutdown = true
	}

	return nil
}
