package services

import (
	"blockchain/models"
	"crypto/sha256"
	"encoding/hex"
)

//计算给定数据的SHA256散列值
func CalculateHash(block models.Block) string {
	record := string(block.Index) + block.Timestamp + string(block.BPM) + block.PreHash
	hash := sha256.New()
	hash.Write([]byte(record))
	Hash := hash.Sum(nil)
	return hex.EncodeToString(Hash)
}
