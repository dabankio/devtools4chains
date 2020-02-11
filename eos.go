// Copyright (c) [2019] [dabank.io]
// [devtools4chains] is licensed under the Mulan PSL v1.
// You can use this software according to the terms and conditions of the Mulan PSL v1.
// You may obtain a copy of Mulan PSL v1 at:
//     http://license.coscl.org.cn/MulanPSL
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
// PURPOSE.
// See the Mulan PSL v1 for more details.

package devtools4chains

// EOSKeyPair .
type EOSKeyPair struct {
	Privk, Pubk string
}

// some dev data
const (
	// ref: https://developers.eos.io/welcome/latest/getting-started/development-environment/create-development-wallet/#step-6-import-the-development-key
	EOSDefaultSYSUser = "eosio"
	EOSDefaultDevKey  = "5KQwrPbwdL6PhXujxW37FSSQZ1JiwsST4cqQzDeyXtP79zkvFD3"
)

// some vars
var (
	// got from:  cleos create key --to-console
	EOSPresetKeyPairs = []EOSKeyPair{
		{
			Privk: "5J1V7CAjsrdZBoqzr8XpyfkpzQGxhZc5oFVKCpdPBU95eS1n5q1",
			Pubk:  "EOS7RxMK9cYn3NFhbnmPtoHCbwukGX8TKhk9GHVmK6spXXUfkCKdw",
		},
		{
			Privk: "5KfcVU6z4JcQKNKyMTL32ganLeiWh3fthXtZS8eWf3ZG1Q2qXy9",
			Pubk:  "EOS7QuCmkGDWmv6HvKazZ4x8CJ6QnMuXn4XvpBEqxwCyFyVb4BP6t",
		},
		{
			Privk: "5J1AfFTb91ghQoYGgzURos4VBaEyyEWcKigShddcN4XndZmQiUQ",
			Pubk:  "EOS6PHbW8hqJDZjx5vEb1wfN5dmUTkt3RvQo73ZAxwea4TK9bMx8i",
		},
	}
)
