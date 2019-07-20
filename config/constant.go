package config

// 默认监听端口
const PRODUCE_PORT = "8888"

// 日志文件夹
const DEFAULT_PATH = "logs/"

// 配置文件路径
const CONFIG_PATH = "config.yml"

// redis队列键名前缀
const REDIS_LIST_KEY = "log"

// redis最大内存
const REDIS_MAX_MEMORY int64 = 1024 * 1024 * 1024 * 1024 * 1 // 1GB

// 处理端休眠间隔（秒）
const CONSUME_SLEEP_INTERVAL = 1