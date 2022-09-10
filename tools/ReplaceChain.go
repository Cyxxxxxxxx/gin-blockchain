package tools

import "blockchain/models"

//最长链原则
func ReplaceChain(BlockChain, NewBlocks []models.Block) {
	if len(NewBlocks) > len(BlockChain) {
		BlockChain = NewBlocks
	}
}
