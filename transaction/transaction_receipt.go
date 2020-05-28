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
package transaction

type TransactionReceipt struct {
	Accepted      bool             `json:"accepted"`
	Success       bool             `json:"success"`
	CumulativeGas string           `json:"cumulative_gas"`
	EpochNum      string           `json:"epoch_num"`
	EventLogs     []interface{}    `json:"event_logs"`
	Transitions   []*InnerTransfer `json:"transitions"`
}

type InnerTransfer struct {
	Addr  string `json:"addr"`
	Depth uint   `json:"depth"`
	Msg   *Msg   `json:"msg"`
}

type Msg struct {
	Amount    string        `json:"_amount"`
	Recipient string        `json:"_recipient"`
	Tag       string        `json:"_tag"`
	Params    []interface{} `json:"params"`
}
