package handlers

import (
	"encoding/json"
	"math/big"
	"net/http"
	"strconv"

	models "ethapi/models"
	modules "ethapi/modules"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockHandler ethereum client instance
type BlockHandler struct {
	Client *ethclient.Client
}

func (c *BlockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		b   *models.Block
		err *models.Error
		n   *big.Int
		ok  bool
	)

	// Get the query parameters from the url request
	hash := r.URL.Query().Get("hash")
	num := r.URL.Query().Get("number")
	lt := r.URL.Query().Get("latest")
	latest, _ := strconv.ParseBool(lt)

	// Set response header to json
	w.Header().Set("Content-Type", "application/json")

	// Check query parameters
	if hash == "" && num == "" && lt == "" {
		json.NewEncoder(w).Encode(&models.Error{
			Code: 400,
			Msg:  "Malformed request, must specify only one of latest, hash or number",
		})
		return
	}
	if len(num) > 0 {
		n = new(big.Int)
		n, ok = n.SetString(num, 10)
		if !ok {
			json.NewEncoder(w).Encode(&models.Error{
				Code: 400,
				Msg:  "Malformed request, number must be a valid base 10 number",
			})
			return
		}
	}

	if latest {
		// Get the latest block
		b, err = modules.GetBlockLatest(c.Client)
	} else {
		if len(hash) > 0 {
			// Get a block by hash
			h := common.HexToHash(hash)
			b, err = modules.GetBlockByHash(c.Client, &h)
		}
		if n != nil {
			// Get a block by num
			b, err = modules.GetBlockByNumber(c.Client, n)
		}
	}

	if err != nil {
		switch err.Code {
		// we rewrite Msg to avoid possible leaks of sensitive information in the response
		case 500:
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
	json.NewEncoder(w).Encode(b)
}
