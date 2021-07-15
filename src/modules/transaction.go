package modules

import (
	"context"
	"strings"

	models "ethapi/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
)

// Retrive a transaction by hash. Returns an error if not found
func GetTx(c *ethclient.Client, hash common.Hash) (*models.Transaction, *models.Error) {
	tx, isPending, err := c.TransactionByHash(context.Background(), hash)
	if err != nil {
		if strings.ToLower(err.Error()) == "not found" {
			return nil, &models.Error{Code: 404, Msg: "Transaction not found"}
		}
		log.Error(err)
		return nil, &models.Error{Code: 500, Msg: err.Error()}
	}

	// Marshal the response to our transaction model and return
	return &models.Transaction{
		Hash:     tx.Hash().String(),
		To:       tx.To().String(),
		Value:    tx.Value().String(),
		Nonce:    tx.Nonce(),
		Pending:  isPending,
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice().Uint64(),
	}, nil
}
