package common

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	//cosmwasmcommon "github.com/airchains-network/cosmwasm-sequencer-node/common"
	"github.com/airchains-network/cosmwasm-sequencer-node/common/logs"
	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/btcutil/bech32"
	"math/big"
	"net/http"
	"os"
	"strconv"
)

func Base52Decoder(value string) string {
	decodedBytes := base58.Decode(value)
	decodedBigInt := new(big.Int).SetBytes(decodedBytes)
	return decodedBigInt.String()
}

func Bech32Decoder(value string) string {
	_, bytes, err := bech32.Decode(value)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Bech32 decoding error : %s", err.Error()))
		os.Exit(0)
	}

	decodedBigInt := new(big.Int).SetBytes(bytes)
	return decodedBigInt.String()
}

func TXHashCheck(value string) string {
	byteSlice, err := hex.DecodeString(value)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Hex decoding error : %s", err.Error()))
		os.Exit(0)
	}

	decodedBigInt := new(big.Int).SetBytes(byteSlice)
	return decodedBigInt.String()
}

func AccountBalanceCheck(walletAddress string, bloackHeight string) string {

	h, err := strconv.Atoi(bloackHeight)
	if err != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Invalid block number : %s"+err.Error()))
		os.Exit(0)
	}
	res, resErr := http.Get(
		fmt.Sprintf(
			"%s/cosmos/bank/v1beta1/balances/%s?height=%d", ExecutionClientTRPC, walletAddress, (h - 1),
		),
	)
	if resErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Resquesting in comswasm RPC : %s"+resErr.Error()))
		os.Exit(0)
	}
	defer res.Body.Close()

	var accountBalance struct {
		Balances []struct {
			Denom  string `json:"denom"`
			Amount string `json:"amount"`
		} `json:"balances"`
		Pagination struct {
			NextKey string `json:"next_key"`
			Total   string `json:"total"`
		} `json:"pagination"`
	}

	decodeError := json.NewDecoder(res.Body).Decode(&accountBalance)
	if decodeError != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Decoding in comswasm RPC : %s"+decodeError.Error()))
		os.Exit(0)
	}

	return accountBalance.Balances[0].Amount
}

func AccountNounceCheck(walletAddress string) string {
	res, resErr := http.Get(
		fmt.Sprintf(
			"%s/cosmos/auth/v1beta1/accounts/%s", ExecutionClientTRPC, walletAddress,
		),
	)
	if resErr != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Resquesting in comswasm RPC : %s"+resErr.Error()))
		os.Exit(0)
	}
	defer res.Body.Close()

	var accountNounce struct {
		Account struct {
			Type    string `json:"@type"`
			Address string `json:"address"`
			PubKey  struct {
				Type string `json:"@type"`
				Key  string `json:"key"`
			} `json:"pub_key"`
			AccountNumber string `json:"account_number"`
			Sequence      string `json:"sequence"`
		} `json:"account"`
	}

	decodeError := json.NewDecoder(res.Body).Decode(&accountNounce)
	if decodeError != nil {
		logs.LogMessage("ERROR:", fmt.Sprintf("Decoding in comswasm RPC : %s"+decodeError.Error()))
		os.Exit(0)
	}

	return accountNounce.Account.Sequence
}
