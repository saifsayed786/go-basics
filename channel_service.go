package libgochannel

func CreateQueue() chan<- interface{} {
	in, out := makeInfinite()
	go func() {
		var m any
		for m = range out {
			m.(ChanMod).Func(m.(ChanMod).Model)
		}
	}()
	return in
}

func makeInfinite() (chan<- interface{}, <-chan interface{}) {
	in := make(chan interface{})
	out := make(chan interface{})
	go func() {
		var inQueue []interface{}
		curVal := func() interface{} {
			//* Check wether queue has any data otherwise return nil
			if len(inQueue) == 0 {
				return nil
			}
			return inQueue[0]
		}
		outCh := func() chan interface{} {
			//* Check wether queue has any data otherwise return nil
			if len(inQueue) == 0 {
				return nil
			}
			return out
		}
		for len(inQueue) > 0 || in != nil {
			select {
			case v, ok := <-in:
				//* Enqueue
				if !ok {
					in = nil
				} else {
					inQueue = append(inQueue, v)

				}
				//* Deqeue to out and remove value from list
			case outCh() <- curVal():
				inQueue = inQueue[1:]
			}
		}

		close(out)
	}()
	return in, out
}

type ChanMod struct {
	Model any
	Func  func(any)
}
