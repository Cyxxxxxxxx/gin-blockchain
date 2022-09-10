package services

import "blockchain/models"

//校验区块
//通过当上一个区块Index+1是否等于当前区块Index;
//上一区块Hash是否等于当前区块记录的上一区块的Prehash;
//计算当前区块Hash是否等于当前区块记录的Hash值
func IsBlockValid(newBlock, oldBlock models.Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}
	if oldBlock.Hash != newBlock.PreHash {
		return false
	}
	if CalculateHash(newBlock) != newBlock.Hash {
		return false
	}
	return true
}
