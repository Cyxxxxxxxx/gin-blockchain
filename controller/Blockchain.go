package controller

import (
	"blockchain/models"
	"blockchain/services"
	"encoding/json"
	"io"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
)

//区块链是一系列已被验证的有效区块
var Blockchain []models.Block

//互斥锁
var mutex = &sync.Mutex{}

/////////////
/////////////
/////////////
/////////////
//接收到请求时写入区块链
func GetBlockchain(c *gin.Context) {
	bytes, err := json.MarshalIndent(Blockchain, "", " ")
	if err != nil {
		ResponseError(c, CodeIndentErr)
	}
	io.WriteString(c.Writer, string(bytes))
}

/////////////
/////////////
/////////////
/////////////
/////////////
/////////////
//将BPM作为有效负载输入
func WriteBlock(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	var msg models.Message

	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&msg); err != nil {
		ResponseError(c, CodeDecoderErr)
	}
	defer c.Request.Body.Close()

	mutex.Lock()
	preBlock := Blockchain[len(Blockchain)-1]
	newBlock := services.GenerateBlock(preBlock, msg.BPM)
	if services.IsBlockValid(newBlock, preBlock) {
		Blockchain = append(Blockchain, newBlock)
		spew.Dump(Blockchain)
	}
	mutex.Unlock()

	ResponseSuccess(c, CodeSuccess, newBlock)
}

func GenesisBlock() {
	t := time.Now()
	genesisBlock := models.Block{}
	genesisBlock = models.Block{0, t.String(), 0, services.CalculateHash(genesisBlock), ""}
	spew.Dump(genesisBlock)

	mutex.Lock()
	Blockchain = append(Blockchain, genesisBlock)
	mutex.Unlock()
}
