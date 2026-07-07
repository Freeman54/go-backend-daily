package consumerbackpressure

import "errors"

type Snapshot struct {
	InFlight int
	MaxInFly int
	Lag      int
}

type Controller struct {
	pauseLag    int
	resumeLag   int
	pauseInFly  int
	resumeInFly int
	paused      bool
}

func New(pauseLag, resumeLag, pauseInFly, resumeInFly int) (*Controller, error) {
	if pauseLag < 0 || resumeLag < 0 || pauseInFly <= 0 || resumeInFly < 0 {
		return nil, errors.New("thresholds must be non-negative and pauseInFly must be positive")
	}
	if resumeLag > pauseLag {
		return nil, errors.New("resumeLag must be less than or equal to pauseLag")
	}
	if resumeInFly > pauseInFly {
		return nil, errors.New("resumeInFly must be less than or equal to pauseInFly")
	}

	return &Controller{
		pauseLag:    pauseLag,
		resumeLag:   resumeLag,
		pauseInFly:  pauseInFly,
		resumeInFly: resumeInFly,
	}, nil
}

func (c *Controller) Decide(snapshot Snapshot) bool {
	if snapshot.MaxInFly < 0 || snapshot.InFlight < 0 || snapshot.Lag < 0 {
		return c.paused
	}

	pressureLag := snapshot.Lag >= c.pauseLag
	pressureInFly := snapshot.InFlight >= c.pauseInFly
	recoverLag := snapshot.Lag <= c.resumeLag
	recoverInFly := snapshot.InFlight <= c.resumeInFly

	if !c.paused && (pressureLag || pressureInFly) {
		c.paused = true
		return true
	}

	if c.paused && recoverLag && recoverInFly {
		c.paused = false
		return false
	}

	return c.paused
}

func (c *Controller) Paused() bool {
	return c.paused
}
