package modules

import (
	"context"
	"math/big"
	"strings"

	log "github.com/sirupsen/logrus"

	models "ethapi/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Retrive a block using the provided ethclient. Returns the block or an error if not found
func getBlock(c *ethclient.Client, hash *common.Hash, num *big.Int) (*models.Block, *models.Error) {
	// Get the header first. HeaderByNumber returns the last block if number is nil
	header, err := c.HeaderByNumber(context.Background(), num)
	if err != nil {
		if strings.ToLower(err.Error()) == "not found" {
			return nil, &models.Error{Code: 404, Msg: "Block not found"}
		}
		log.Error(err)
		return nil, &models.Error{Code: 500, Msg: err.Error()}
	}
	blockNumber := big.NewInt(header.Number.Int64())
	// Get the block now by number
	b, err := c.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Error(err)
		return nil, &models.Error{Code: 500, Msg: err.Error()}
	}

	// Marshal the response to our block model
	bModel := &models.Block{
		BlockNum:   b.Number().Int64(),
		Hash:       b.Hash().String(),
		Difficulty: b.Difficulty().Uint64(),
		Timestamp:  b.Time(),
		TxCount:    len(b.Transactions()),
		Txs:        []models.Transaction{},
	}

	// Enumerate the transactions in the block
	for _, tx := range b.Transactions() {
		var to string
		if tx.To() != nil {
			to = tx.To().String()
		}
		bModel.Txs = append(bModel.Txs, models.Transaction{
			Hash:     tx.Hash().String(),
			To:       to,
			Value:    tx.Value().String(),
			Nonce:    tx.Nonce(),
			Gas:      tx.Gas(),
			GasPrice: tx.GasPrice().Uint64(),
		})
	}

	return bModel, nil
}

// Retrieve a block by its hash or an error if not found
func GetBlockByHash(c *ethclient.Client, hash *common.Hash) (*models.Block, *models.Error) {
	return getBlock(c, hash, nil)
}

// Retrive a block by its number or an error if not found
func GetBlockByNumber(c *ethclient.Client, num *big.Int) (*models.Block, *models.Error) {
	return getBlock(c, nil, num)
}

// Retrive the latest block or an error if not found
func GetBlockLatest(c *ethclient.Client) (*models.Block, *models.Error) {
	return getBlock(c, nil, nil)
}
