package models

//对应区块结构体
type Block struct {
	Index     int    //这个块在整个链中的位置
	Timestamp string //块生成时的时间戳
	BPM       int    //这个块通过SHA256算法生成的散列值
	Hash      string //代表前一个块的SHA256散列值
	PreHash   string //每分钟心跳数 也就是心率
}

//传入json参数
type Message struct {
	BPM int `json:"bpm"`
}
