package cmd

import "errors"

var (
	errNotFound     = errors.New("请求资源不存在，请检查url")
	errNotSupported = errors.New("该资源不支持基于http的tail, 请更改服务器配置。")
)
