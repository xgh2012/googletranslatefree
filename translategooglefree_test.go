package translategooglefree

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	cfg := &Config{
		TimeOut: 10,
		Proxy:   "http://127.0.0.1:1080/",
	}
	result, lang, err := NewClient(cfg).Translate("这是测试数据", "zh-CN", "en")
	fmt.Println(result, lang, err)
}
