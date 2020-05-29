package provider

import (
	"encoding/json"
	"fmt"
	"github.com/ybbus/jsonrpc"
	"math/big"
	"strconv"
)

var EmptyBlock = fmt.Errorf("empty block")
var NotContract = fmt.Errorf("Address not contract address")

func ParseTxBlock(rpcResult *jsonrpc.RPCResponse) (*TxBlock, error) {
	if rpcResult.Error != nil {
		return nil, fmt.Errorf("ParseTxBlock: resp code %d, msg %s", rpcResult.Error.Code,
			rpcResult.Error.Message)
	}
	jsonResult, err := json.Marshal(rpcResult.Result)
	if err != nil {
		return nil, fmt.Errorf("ParseTxBlock: marshal rpc result, %s", err)
	}
	block := &TxBlock{}
	if err := json.Unmarshal(jsonResult, block); err != nil {
		return nil, fmt.Errorf("ParseTxBlock: unmarshal txBlock, %s", err)
	}
	return block, nil
}

func ParseTxHashArray(rpcResult *jsonrpc.RPCResponse) ([][]string, error) {
	if rpcResult.Error != nil {
		if rpcResult.Error.Message == "TxBlock has no transactions" {
			return nil, EmptyBlock
		}
		return nil, fmt.Errorf("ParseTxHashArray: resp code %d, msg %s", rpcResult.Error.Code,
			rpcResult.Error.Message)
	}
	jsonResult, err := json.Marshal(rpcResult.Result)
	if err != nil {
		return nil, fmt.Errorf("ParseTxHashArray: marshal rpc result, %s", err)
	}
	result := [][]string{}
	if err := json.Unmarshal(jsonResult, &result); err != nil {
		return nil, fmt.Errorf("ParseTxHashArray: unmarshal result, %s", err)
	}
	return result, nil
}

func ParseBlockHeight(rpcResult *jsonrpc.RPCResponse) (uint64, error) {
	if rpcResult.Error != nil {
		return 0, fmt.Errorf("ParseBlockHeight: resp code %d, msg %s", rpcResult.Error.Code,
			rpcResult.Error.Message)
	}
	if heightString, ok := rpcResult.Result.(string); ok {
		height, err := strconv.ParseUint(heightString, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("ParseBlockHeight: result %s invalid, %s", heightString, err)
		}
		return height, nil
	} else {
		return 0, fmt.Errorf("ParseBlockHeight: type unmatch")
	}
}

func ParseCreateTxResult(rpcResult *jsonrpc.RPCResponse) (info string, contract string, hash string, err error) {
	if rpcResult.Error != nil {
		return "", "", "", fmt.Errorf("ParseCreateTxResult: resp code %d, msg %s",
			rpcResult.Error.Code, rpcResult.Error.Message)
	}
	type Result struct {
		ContractAddress string
		Info            string
		TranID          string
	}
	jsonResult, err := json.Marshal(rpcResult.Result)
	if err != nil {
		err = fmt.Errorf("ParseCreateTxResult: marshal rpc result, %s", err)
		return
	}
	result := &Result{}
	if err = json.Unmarshal(jsonResult, &result); err != nil {
		err = fmt.Errorf("ParseCreateTxResult: unmarshal result, %s", err)
		return
	}
	return result.Info, result.ContractAddress, result.TranID, nil
}

func ParseBalanceResp(rpcResult *jsonrpc.RPCResponse) (*big.Int, uint64, error) {
	if rpcResult.Error != nil {
		return nil, 0, fmt.Errorf("ParseBalanceResp: resp code %d, msg %s", rpcResult.Error.Code,
			rpcResult.Error.Message)
	}
	type Balance struct {
		Balance string `json:"balance"`
		Nonce   uint64 `json:"nonce"`
	}
	jsonResult, err := json.Marshal(rpcResult.Result)
	if err != nil {
		return nil, 0, fmt.Errorf("ParseBalanceResp: marshal rpc result, %s", err)
	}
	result := &Balance{}
	if err = json.Unmarshal(jsonResult, &result); err != nil {
		return nil, 0, fmt.Errorf("ParseBalanceResp: unmarshal result, %s", err)
	}
	balance, ok := new(big.Int).SetString(result.Balance, 10)
	if !ok {
		return nil, 0, fmt.Errorf("ParseBalanceResp: balance %s invalid", result.Balance)
	}
	return balance, result.Nonce, nil
}

func ParseGetContractCode(rpcResult *jsonrpc.RPCResponse) (code string, err error) {
	if rpcResult.Error != nil {
		return "", fmt.Errorf("ParseGetContractCode: resp code %d, msg %s", rpcResult.Error.Code,
			rpcResult.Error.Message)
	}
	type Contract struct {
		Code string `json:"code"`
	}
	jsonResult, err := json.Marshal(rpcResult.Result)
	if err != nil {
		return "", fmt.Errorf("ParseGetContractCode: marshal rpc result, %s", err)
	}
	result := &Contract{}
	if err = json.Unmarshal(jsonResult, &result); err != nil {
		return "", fmt.Errorf("ParseGetContractCode: unmarshal result, %s", err)
	}
	return result.Code, nil
}

func ParseGetContractInit(rpcResult *jsonrpc.RPCResponse) ([]Value, error) {
	if rpcResult.Error != nil {
		if rpcResult.Error.Message == "Address not contract address" {
			return nil, NotContract
		}
		return nil, fmt.Errorf("ParseGetContractInit: resp code %d, msg %s", rpcResult.Error.Code,
			rpcResult.Error.Message)
	}
	jsonResult, err := json.Marshal(rpcResult.Result)
	if err != nil {
		return nil, fmt.Errorf("ParseGetContractCode: marshal rpc result, %s", err)
	}
	result := make([]Value, 0)
	if err = json.Unmarshal(jsonResult, &result); err != nil {
		return nil, fmt.Errorf("ParseGetContractCode: unmarshal result, %s", err)
	}
	return result, nil
}
