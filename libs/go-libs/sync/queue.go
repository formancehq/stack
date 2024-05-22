package sync

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type Queue struct {
	noCopy noCopy
	sem    chan struct{}
}

func (q *Queue) Lock() {
	q.sem <- struct{}{}

}

func (q *Queue) Unlock() {
	<-q.sem
}

func NewQueue(maxConcurrency int) *Queue {
	return &Queue{
		sem: make(chan struct{}, maxConcurrency),
	}
}
