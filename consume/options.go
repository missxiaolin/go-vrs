package consume

type Options func (d *Dispatch)

func NewCollection(c Collection) Options {
	return func(d *Dispatch) {
		d.c = c
	}
}

func NewBuilder(b Builder) Options {
	return func(d *Dispatch) {
		d.b = b
	}
}

func NewTransport(t Transport) Options {
	return func(d *Dispatch) {
		d.t = t
	}
}

func NewErrHandle(e ErrHandle) Options {
	return func(d *Dispatch) {
		d.e = e
	}
}