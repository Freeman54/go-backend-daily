package lagautoscaler

type Snapshot struct {
	Lag                   int64
	InflowPerSecond       int64
	CurrentConsumers      int
	ThroughputPerConsumer int64
}

type Policy struct {
	TargetCatchupSeconds int64
	MinConsumers         int
	MaxConsumers         int
	MaxScaleStep         int
}

func Decide(policy Policy, snapshot Snapshot) int {
	if snapshot.ThroughputPerConsumer <= 0 {
		return clamp(snapshot.CurrentConsumers, policy.MinConsumers, policy.MaxConsumers)
	}

	lagDriven := ceilDiv(snapshot.Lag, policy.TargetCatchupSeconds)
	requiredRate := max64(lagDriven, snapshot.InflowPerSecond)
	desired := int(ceilDiv(requiredRate, snapshot.ThroughputPerConsumer))
	desired = clamp(desired, policy.MinConsumers, policy.MaxConsumers)

	high := snapshot.CurrentConsumers + policy.MaxScaleStep
	low := snapshot.CurrentConsumers - policy.MaxScaleStep
	if desired > high {
		return high
	}
	if desired < low {
		return low
	}
	return desired
}

func ceilDiv(numerator, denominator int64) int64 {
	if numerator <= 0 {
		return 0
	}
	if denominator <= 0 {
		return numerator
	}
	return (numerator + denominator - 1) / denominator
}

func clamp(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func max64(left, right int64) int64 {
	if left > right {
		return left
	}
	return right
}
