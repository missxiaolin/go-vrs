package pool

type Counter struct {
	pool         *Pool
	count        int
	triggerCount int
}

func (c *Counter) Run(task func()) error {

	c.count++

	// 当次数达到触发次数时
	if c.count >= TRIGGER_COUNT {
		// 重置统计
		c.count = 0

		c.pool.mux.Lock()
		// 判断协程池是否还有空位
		if c.pool.Free() > 0 {
			// 创建协程
			if err := c.pool.Submit(task); err != nil {
				return err
			}
		}
		c.pool.mux.Unlock()
	}
	return nil
}

func (c *Counter) SetTriggerCount(count int) {
	c.triggerCount = count
}

func (c *Counter) SetCount(count int) {
	c.count = count
}