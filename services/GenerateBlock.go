package services

import (
	"blockchain/models"
	"time"
)

//生成区块
func GenerateBlock(oldBlock models.Block, BPM int) models.Block {
	var newBlock models.Block
	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.BPM = BPM
	newBlock.PreHash = oldBlock.Hash
	newBlock.Hash = CalculateHash(newBlock)

	return newBlock
}
