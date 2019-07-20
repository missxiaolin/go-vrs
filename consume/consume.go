package consume


// 数据调度器
type Dispatch struct {
	c Collection
	b Builder
	t Transport
	e ErrHandle
}

func NewDispatch(opts ...Options) *Dispatch {

	d := &Dispatch{}
	for _, opt := range opts {
		opt(d)
	}
	return d
}

func (d *Dispatch) Run() (Result, error) {

	var data Data
	var err error
	data, err = d.c.Collect()
	if err != nil {
		return nil, err
	}

	if d.b != nil {
		data, err = d.b.Build(data)
		if err != nil {
			return nil, err
		}
	}
	return d.t.Send(data)
}

