package poolbackpressure

type Decision string

const (
	DecisionAdmit   Decision = "admit"
	DecisionDegrade Decision = "degrade"
	DecisionReject  Decision = "reject"
)

type Snapshot struct {
	InUse   int
	MaxOpen int
	Waiters int
}

type Controller struct {
	degradeUsageRatio float64
	rejectUsageRatio  float64
	rejectWaiters     int
}

func New(degradeUsageRatio float64, rejectUsageRatio float64, rejectWaiters int) *Controller {
	if degradeUsageRatio <= 0 {
		degradeUsageRatio = 0.7
	}
	if rejectUsageRatio <= 0 {
		rejectUsageRatio = 0.9
	}
	if rejectWaiters <= 0 {
		rejectWaiters = 32
	}

	return &Controller{
		degradeUsageRatio: degradeUsageRatio,
		rejectUsageRatio:  rejectUsageRatio,
		rejectWaiters:     rejectWaiters,
	}
}

func (c *Controller) Decide(s Snapshot) Decision {
	if s.MaxOpen <= 0 {
		return DecisionReject
	}

	usage := float64(s.InUse) / float64(s.MaxOpen)
	if usage >= c.rejectUsageRatio || s.Waiters >= c.rejectWaiters {
		return DecisionReject
	}
	if usage >= c.degradeUsageRatio {
		return DecisionDegrade
	}
	return DecisionAdmit
}
