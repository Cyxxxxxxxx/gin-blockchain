package tools

import "time"

// InitInMain 初始化时区
func InitInMain() {
	var cstZone = time.FixedZone("CST", 8*3600) // 东八
	time.Local = cstZone
}
