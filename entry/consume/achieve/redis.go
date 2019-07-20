package achieve

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"go-vrs/Model"
	"go-vrs/config"
	"go-vrs/consume"
	"go-vrs/lib/logger"
)

var logs = logger.New()

type RedisAch struct {
	Client *redis.Client
	Num    int
}

func (r *RedisAch) Collect() (consume.Data, error) {

	if r.Num == 0 {
		r.Num = 1
	}

	pip := r.Client.Pipeline()
	for i := 0; i < r.Num; i++ {
		pip.RPop(config.Cfg.Redis.ListKey)
	}

	cmds, err := pip.Exec()
	cmdsLen := len(cmds)
	if cmdsLen > 0 {
		start := cmds[0].(*redis.StringCmd)
		_, err = start.Bytes()
	}

	return cmds, err
}

func (r *RedisAch) Build(d consume.Data) (consume.Data, error) {

	cmds := d.([]redis.Cmder)
	set := make([]interface{}, 0)

	if len(cmds) != r.Num {
		fmt.Println(len(cmds))
	}

	for _, v := range cmds {
		cmd := v.(*redis.StringCmd)
		b, err := cmd.Bytes()
		if err != nil {
			r.errHandle(err, nil)
			continue
		}
		var data Model.Data
		if err := json.Unmarshal(b, &data); err != nil {
			r.errHandle(err, b)
			continue
		}
		set = append(set, data)
	}

	return set, nil
}

// 处理请求/数据解析失败
func (r *RedisAch) errHandle(err error, a interface{}) {
	if err == redis.Nil {
		return
	}
	logs.WithField("content", a).Errorln(err)
}
