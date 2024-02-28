package dispatcher

type processor struct {
	cache       *buffer
	handler     func(interface{})
	inCh, outCh chan interface{}
}

func NewProcessor(size int) *processor {
	return &processor{
		cache: NewBuffer(size),
		outCh: make(chan interface{}),
		inCh:  make(chan interface{}),
	}
}
func (p *processor) register(handler func(interface{})) {
	p.handler = handler
}

func (p *processor) write() {
	for {
		select {
		case data := <-p.inCh:
			p.cache.WriteOne(data)
		}
	}
}

func (p *processor) read() {
	for {
		data := p.cache.ReadOne()
		if data != nil {
			p.outCh <- data
		}
	}

}

func (p *processor) run() {
	go p.read()
	go p.write()
	for {
		select {
		case data := <-p.outCh:
			if data != nil {
				p.handler(data)
			}
		}
	}
}

func (p *processor) stop() {
	close(p.inCh)
	close(p.outCh)
}
