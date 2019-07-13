package logger

import (
	"github.com/sirupsen/logrus"
	"go-vrs/config"
	"io"
	"os"
	"strings"
	"time"
)

type hook struct {
	w    io.ReadWriteCloser
	path string
}

func (h *hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// 自定义钩子执行（默认协程安全）
func (h *hook) Fire(e *logrus.Entry) error {
	// 判断log文件是否变更
	path := config.DEFAULT_PATH + time.Now().Format("20060102")
	if h.path != path {
		h.path = path
		if h.w != nil {
			if err := h.w.Close(); err != nil {
				return err
			}
			h.w = nil
		}
	}
	path = path + ".log"

	if _, err := os.Stat(path); err == nil && h.w != nil {
		return nil
	}

	// 自动创建文件
	var pathArr = strings.Split(path, "/")
	var pathLen = len(pathArr)
	if pathLen > 1 {
		dir := strings.Join(pathArr[:pathLen-1], "/")
		// 自动创建日志文件夹
		_, err := os.Stat(dir)
		if err != nil {
			mkErr := os.MkdirAll(dir, os.ModePerm)
			if mkErr != nil {
				return mkErr
			}
		}
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		return err
	}

	h.w = f
	e.Logger.Out = h.w
	return nil
}