package replicalagguard

import "time"

type Replica struct {
	Name    string
	Lag     time.Duration
	Healthy bool
}

type Decision struct {
	Target string
	Reason string
}

func Choose(readAfterWrite bool, maxLag time.Duration, replicas []Replica) Decision {
	if readAfterWrite {
		return Decision{Target: "primary", Reason: "read-after-write"}
	}

	best := Replica{}
	found := false
	for _, replica := range replicas {
		if !replica.Healthy || replica.Lag > maxLag {
			continue
		}
		if !found || replica.Lag < best.Lag {
			best = replica
			found = true
		}
	}

	if found {
		return Decision{Target: best.Name, Reason: "replica-ok"}
	}
	return Decision{Target: "primary", Reason: "replica-lag"}
}
