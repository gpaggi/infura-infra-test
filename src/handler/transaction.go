package handlers

import (
	"encoding/json"
	"net/http"

	models "ethapi/models"
	modules "ethapi/modules"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// TxHandler ethereum client instance
type TxHandler struct {
	Client *ethclient.Client
}

func (c *TxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the query parameters from the url request
	hash := r.URL.Query().Get("hash")

	// Set response header to json
	w.Header().Set("Content-Type", "application/json")

	// Check query parameters
	if len(hash) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&models.Error{
			Code: 400,
			Msg:  "Malformed request, must specify hash",
		})
		return
	}

	// Get the tx by hash
	txHash := common.HexToHash(hash)
	tx, err := modules.GetTx(c.Client, txHash)

	if err != nil {
		switch err.Code {
		// we rewrite Msg to avoid possible leaks of sensitive information in the response
		case 500:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(&models.Error{
				Code: 500,
				Msg:  "Internal server error",
			})
		default:
			json.NewEncoder(w).Encode(err)
		}
		return
	}

	// finally we write out response
	json.NewEncoder(w).Encode(tx)
}
