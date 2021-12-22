package dropchan

type DropChan struct {
	input  chan interface{}
	buffer chan interface{}
	handle HandleFunc
}

type HandleFunc func(data interface{})

func New(size int, drop bool, handle HandleFunc) *DropChan {
	r := &DropChan{handle: handle}

	if size <= 0 {
		size = 1
	}

	if drop {
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
			d := <-r.buffer
			r.buffer <- v

			if r.handle != nil {
				r.handle(d)
			}
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
