package replicastickiness

import "time"

type Replica struct {
	Name    string
	Lag     time.Duration
	Healthy bool
}

type Route struct {
	Replica    string
	SessionTTL time.Duration
}

type Router struct {
	MaxLag         time.Duration
	SessionTTL     time.Duration
	now            func() time.Time
	sessionReplica map[string]string
	sessionExpiry  map[string]time.Time
}

func NewRouter(maxLag, sessionTTL time.Duration) *Router {
	return &Router{
		MaxLag:         maxLag,
		SessionTTL:     sessionTTL,
		now:            time.Now,
		sessionReplica: map[string]string{},
		sessionExpiry:  map[string]time.Time{},
	}
}

func (r *Router) Pick(sessionID string, replicas []Replica) Route {
	now := r.now()
	if current, ok := r.sessionReplica[sessionID]; ok {
		if expiry := r.sessionExpiry[sessionID]; now.Before(expiry) && containsHealthyReplica(current, replicas, r.MaxLag) {
			return Route{Replica: current, SessionTTL: expiry.Sub(now)}
		}
	}

	for _, replica := range replicas {
		if replica.Healthy && replica.Lag <= r.MaxLag {
			r.sessionReplica[sessionID] = replica.Name
			r.sessionExpiry[sessionID] = now.Add(r.SessionTTL)
			return Route{Replica: replica.Name, SessionTTL: r.SessionTTL}
		}
	}
	return Route{}
}

func containsHealthyReplica(name string, replicas []Replica, maxLag time.Duration) bool {
	for _, replica := range replicas {
		if replica.Name == name && replica.Healthy && replica.Lag <= maxLag {
			return true
		}
	}
	return false
}
