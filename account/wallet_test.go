/*
 * Copyright (C) 2019 Zilliqa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package account

import (
	"fmt"
	provider2 "github.com/Zilliqa/gozilliqa-sdk/provider"
	"github.com/Zilliqa/gozilliqa-sdk/transaction"
	"github.com/Zilliqa/gozilliqa-sdk/util"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestWallet_SignWith(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	wallet := NewWallet()
	wallet.AddByPrivateKey("e19d05c5452598e24caad4a0d85a49146f7be089515c905ae6a19e8a578a6930")
	tx := &transaction.Transaction{
		Version:  "65535",
		Amount:   "0",
		GasPrice: "100000",
		GasLimit: "1",
	}
	provider := provider2.NewProvider("https://dev-api.zilliqa.com/")
	err := wallet.SignWith(tx, "9bfec715a6bd658fcb62b0f8cc9bfa2ade71434a", *provider)
	assert.Nil(t, err, err)
	assert.NotEmpty(t, tx.Signature)
}

func TestSendTransaction(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	wallet := NewWallet()
	wallet.AddByPrivateKey("e19d05c5452598e24caad4a0d85a49146f7be089515c905ae6a19e8a578a6930")
	provider := provider2.NewProvider("https://dev-api.zilliqa.com/")

	result, err := provider.GetMinimumGasPrice()
	assert.Nil(t, err, err)
	gasPrice := result.Result.(string)

	tx := &transaction.Transaction{
		Version:      strconv.FormatInt(int64(util.Pack(333, 1)), 10),
		SenderPubKey: "0246E7178DC8253201101E18FD6F6EB9972451D121FC57AA2A06DD5C111E58DC6A",
		ToAddr:       "4BAF5faDA8e5Db92C3d3242618c5B47133AE003C",
		Amount:       "10000000",
		GasPrice:     gasPrice,
		GasLimit:     "1",
		Code:         "",
		Data:         "",
		Priority:     false,
	}

	err2 := wallet.Sign(tx, *provider)
	assert.Nil(t, err2, err2)

	rsp, err3 := provider.CreateTransaction(tx.ToTransactionPayload())
	assert.Nil(t, err3, err3)
	assert.Nil(t, rsp.Error, rsp.Error)

	resMap := rsp.Result.(map[string]interface{})
	hash := resMap["TranID"].(string)
	fmt.Printf("hash is %s\n", hash)
	tx.Confirm(hash, 1000, 3, provider)
	assert.True(t, tx.Status == transaction.Confirmed)
}

func TestSendTransactionInsufficientAmount(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
	wallet := NewWallet()
	wallet.AddByPrivateKey("e19d05c5452598e24caad4a0d85a49146f7be089515c905ae6a19e8a578a6930")
	provider := provider2.NewProvider("https://dev-api.zilliqa.com/")

	result, err := provider.GetMinimumGasPrice()
	assert.Nil(t, err, err)
	gasPrice := result.Result.(string)

	tx := &transaction.Transaction{
		Version:      strconv.FormatInt(int64(util.Pack(333, 1)), 10),
		SenderPubKey: "0246E7178DC8253201101E18FD6F6EB9972451D121FC57AA2A06DD5C111E58DC6A",
		ToAddr:       "4BAF5faDA8e5Db92C3d3242618c5B47133AE003C",
		Amount:       "10000000000000000",
		GasPrice:     gasPrice,
		GasLimit:     "1",
		Code:         "",
		Data:         "",
		Priority:     false,
	}

	err2 := wallet.Sign(tx, *provider)
	assert.NotNil(t, err2)
	assert.Equal(t, err2.Error(), "balance is not sufficient")
}
