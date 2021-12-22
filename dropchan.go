package dropchan

type DropChan struct {
	input  chan interface{}
	buffer chan interface{}
}

func New(size int, drop ...bool) *DropChan {
	r := &DropChan{}

	if len(drop) > 0 && drop[0] {
		r.input = make(chan interface{})
		r.buffer = make(chan interface{}, size)

		go r.run()
	} else {
		r.buffer = make(chan interface{}, size)
		r.input = r.buffer
	}

	return r
}

func (r *DropChan) run() {
	for v := range r.input {
		select {
		case r.buffer <- v:
		default:
			<-r.buffer
			r.buffer <- v
		}
	}

	close(r.buffer)
}

func (r *DropChan) Close() {
	close(r.input)
}

func (r *DropChan) Input() chan<- interface{} {
	return r.input
}

func (r *DropChan) Output() <-chan interface{} {
	return r.buffer
}

func (r *DropChan) Len() int {
	return len(r.buffer)
}

func (r *DropChan) Cap() int {
	return cap(r.buffer)
}
