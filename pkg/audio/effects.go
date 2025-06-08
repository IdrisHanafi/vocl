package audio

// Effect represents an audio processing effect
type Effect interface {
	Process(sample float32) float32
}

// EchoEffect implements a simple echo/delay effect
type EchoEffect struct {
	buffer    []float32
	position  int
	delayTime float32
	feedback  float32
	mix       float32
}

// NewEchoEffect creates a new echo effect with the given parameters
func NewEchoEffect(sampleRate float64, delayTimeMs float32, feedback float32, mix float32) *EchoEffect {
	// Convert delay time from milliseconds to samples
	delaySamples := int(float32(sampleRate) * delayTimeMs / 1000.0)
	return &EchoEffect{
		buffer:    make([]float32, delaySamples),
		position:  0,
		delayTime: delayTimeMs,
		feedback:  feedback,
		mix:       mix,
	}
}

// Process applies the echo effect to a single sample
func (e *EchoEffect) Process(sample float32) float32 {
	// Get the delayed sample
	delayedSample := e.buffer[e.position]

	// Mix the current sample with the delayed sample
	output := sample + delayedSample*e.mix

	// Store the current sample plus feedback in the buffer
	e.buffer[e.position] = sample + delayedSample*e.feedback

	// Move the position in the buffer
	e.position = (e.position + 1) % len(e.buffer)

	return output
}
