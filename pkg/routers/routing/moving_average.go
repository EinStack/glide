package routing

// MovingAverage represents the exponentially weighted moving average of a series of numbers
type MovingAverage struct {
	// The multiplier factor by which the previous samples decay
	decay float64
	// The current value of the average
	value float64
	// The number of samples added to this instance.
	count uint8
	// The number of samples required to start estimating average
	warmupSamples uint8
}

func NewMovingAverage(decay float64, warmupSamples uint8) *MovingAverage {
	return &MovingAverage{
		decay:         decay,
		warmupSamples: warmupSamples,
		count:         0,
		value:         0,
	}
}

// Add a value to the series and updates the moving average
func (e *MovingAverage) Add(value float64) {
	switch {
	case e.count < e.warmupSamples:
		e.count++
		e.value += value
	case e.count == e.warmupSamples:
		e.count++
		e.value = e.value / float64(e.warmupSamples)
		e.value = (value * e.decay) + (e.value * (1 - e.decay))
	default:
		e.value = (value * e.decay) + (e.value * (1 - e.decay))
	}
}

func (e *MovingAverage) WarmedUp() bool {
	return e.count > e.warmupSamples
}

// Value returns the current value of the average, or 0.0 if the series hasn't
// warmed up yet
func (e *MovingAverage) Value() float64 {
	if !e.WarmedUp() {
		return 0.0
	}

	return e.value
}

// Set sets the moving average value
func (e *MovingAverage) Set(value float64) {
	e.value = value

	if !e.WarmedUp() {
		e.count = e.warmupSamples + 1
	}
}
