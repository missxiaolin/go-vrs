package main

import (
	"github.com/go-redis/redis"
	"go-vrs/config"
	"go-vrs/entry/consume/achieve"
	"go-vrs/lib/logger"
	"go-vrs/lib/pool"
	"time"
)

var p *pool.Pool
var logs = logger.New()

func init ()  {
	if p == nil {
		var err error
		p, err = pool.NewPool(config.Cfg.Consume.PoolMaxCapacity)
		if err != nil {
			logs.Fatalln(err)
		}
	}
}

func main()  {
	counter := p.NewCounter(config.Cfg.Consume.TriggerCount)
	dispatch := achieve.NewDispatch()

	for {
		_, err := dispatch.Run()
		if err != nil {
			if err == redis.Nil {
				counter.SetCount(0)
				time.Sleep(time.Duration(config.Cfg.Consume.SleepInterval) * time.Second)
				continue
			}
			logs.Errorln(err)
		}
		if err := counter.Run(Loop); err != nil {
			logs.Errorln(err)
		}
	}
}

func Loop()  {
	counter := p.NewCounter(config.Cfg.Consume.TriggerCount)
	dispatch := achieve.NewDispatch()
	for {
		_, err := dispatch.Run()
		if err != nil {
			if err != redis.Nil {
				// 抓到异常，直接记录异常然后退出
				logs.Fatalln(err)
			}
			break
		}
		if err := counter.Run(Loop); err != nil {
			// 抓到异常，直接记录异常然后退出
			logs.Fatalln(err)
		}
	}
}