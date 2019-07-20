package consume

// 数据
type Data interface{}

// 结果
type Result interface{}

// 收集器
type Collection interface {
	Collect() (Data, error)
}

// 组合器
type Builder interface {
	Build(Data) (Data, error)
}

// 传输器
type Transport interface {
	Send(Data) (Result, error)
}

// 异常处理器
type ErrHandle interface {
	Catch() error
}