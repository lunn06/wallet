package semaphore

// Semaphore implements semaphore pattern
type Semaphore struct {
	semaCh chan struct{}
}

func New(maxReq int) *Semaphore {
	return &Semaphore{
		semaCh: make(chan struct{}, maxReq),
	}
}

// Acquire take place in semaphore channel
// if channel full method will sleep
// until channel will be Release
func (s *Semaphore) Acquire() {
	s.semaCh <- struct{}{}
}

// Release place in semaphore channel
func (s *Semaphore) Release() {
	<-s.semaCh
}
