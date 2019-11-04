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

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

//some parity related const
const (
	ParityCMD                   = "parity"
	ParityDefaultJSONRPCPort    = "8548"
	ParityDefaultJSONRPCPortInt = 8548
)

// ParityEthereumDefaultArgs .
func ParityEthereumDefaultArgs() map[string]*string {
	return map[string]*string{
		"--no-ws":         nil,
		"--no-ipc":        nil,
		"--json-rpc-port": pstring(ParityDefaultJSONRPCPort),
	}
}

// RunParityEthereumConfig parity 以太坊/以太经典运行配置
type RunParityEthereumConfig struct {
	DataDir         DataDirOption
	Args            map[string]*string //k-v ,v 为nil时为flag
	NotPrint2stdout bool               //不打印到stdout(cmd 的stdout不会指定到os.stdout)

	// NetworkID string
	ChainJSON string //创建临时目录后，json 会写入临时目录下的parity.json里，通过命令行引用 (如果不使用临时目录，该选项不会生效)
}

// RunParityEthereum .
func RunParityEthereum(optionsP *RunParityEthereumConfig) (func(), error) {
	if _, err := exec.LookPath(ParityCMD); err != nil {
		return nil, fmt.Errorf("look path err %v", err)
	}

	killHooks := []killHook{}

	var options RunParityEthereumConfig
	var err error

	if optionsP == nil {
		options = RunParityEthereumConfig{}
	} else {
		options = *optionsP
	}
	if options.Args == nil {
		options.Args = map[string]*string{}
	}

	var dataDir string
	if options.DataDir.NewTmpDir {
		for k, v := range options.Args {
			if k == "-d" || k == "--base-path" {
				return nil, fmt.Errorf("datadir specified in args (%v), NewTmpDir not work", v)
			}
		}

		tmpDirPrefix := strings.TrimLeft(options.DataDir.TmpDirPrefix, "/")

		dataDir = strings.TrimRight(os.TempDir(), "/") + "/" + tmpDirPrefix + "parity_data_tmp_" + time.Now().Format(rfc3339Variant) + "/"
		err := os.MkdirAll(dataDir, 0777)
		if err != nil {
			return nil, fmt.Errorf("cannot create tmp dir: %v, err: %v", dataDir, err)
		}
		options.Args["-d"] = &dataDir

		if options.DataDir.NotRemoveTmpDirWhenKilling {
		} else {
			killHooks = append(killHooks, func() error {
				return os.RemoveAll(dataDir)
			})
		}

		if options.ChainJSON != "" {
			jsonFile := dataDir + "parity.json"
			if e := ioutil.WriteFile(jsonFile, []byte(options.ChainJSON), 0777); e != nil {
				return nil, fmt.Errorf("failed to create parity.json in created tmp dir, %v", e)
			}
			options.Args["--chain"] = &jsonFile
		}
	}

	args := []string{}
	for k, v := range options.Args {
		if v == nil {
			args = append(args, k)
		} else {
			args = append(args, k+"="+*v)
		}
	}

	closeChan := make(chan struct{})

	cmd := exec.Command(ParityCMD, args...)
	fmt.Println("[debug] parity-node args", cmd.Args)
	if options.NotPrint2stdout {
	} else {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		fmt.Println("Waiting for message to kill parity")
		<-closeChan
		fmt.Println("Received message,killing parity server")

		if e := cmd.Process.Kill(); e != nil {
			fmt.Println("关闭 parity 时发生异常", e)
		}
		closeChan <- struct{}{}
	}()

	// err = cmd.Wait()
	return func() {
		closeChan <- struct{}{}
		for _, hook := range killHooks {
			hook()
		}
		<-closeChan
	}, nil
}

//some chain json
const (
	//copied from https://github.com/paritytech/parity-ethereum/blob/master/ethcore/res/instant_seal.json
	// removed: account 0000000000000000000000000000000000001337
	ParityDevChainJSON = `
	{
		"name": "DevelopmentChain",
		"engine": {
			"instantSeal": {
				"params": {}
			}
		},
		"params": {
			"gasLimitBoundDivisor": "0x0400",
			"accountStartNonce": "0x0",
			"maximumExtraDataSize": "0x20",
			"minGasLimit": "0x1388",
			"networkID" : "0x11",
			"registrar" : "0x0000000000000000000000000000000000001337",
			"eip150Transition": "0x0",
			"eip160Transition": "0x0",
			"eip161abcTransition": "0x0",
			"eip161dTransition": "0x0",
			"eip155Transition": "0x0",
			"eip98Transition": "0x7fffffffffffff",
			"maxCodeSize": 24576,
			"maxCodeSizeTransition": "0x0",
			"eip140Transition": "0x0",
			"eip211Transition": "0x0",
			"eip214Transition": "0x0",
			"eip658Transition": "0x0",
			"eip145Transition": "0x0",
			"eip1014Transition": "0x0",
			"eip1052Transition": "0x0",
			"wasmActivationTransition": "0x0"
		},
		"genesis": {
			"seal": {
				"generic": "0x0"
			},
			"difficulty": "0x20000",
			"author": "0x0000000000000000000000000000000000000000",
			"timestamp": "0x00",
			"parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
			"extraData": "0x",
			"gasLimit": "0x7A1200"
		},
		"accounts": {
			"0000000000000000000000000000000000000001": { "balance": "1", "builtin": { "name": "ecrecover", "pricing": { "linear": { "base": 3000, "word": 0 } } } },
			"0000000000000000000000000000000000000002": { "balance": "1", "builtin": { "name": "sha256", "pricing": { "linear": { "base": 60, "word": 12 } } } },
			"0000000000000000000000000000000000000003": { "balance": "1", "builtin": { "name": "ripemd160", "pricing": { "linear": { "base": 600, "word": 120 } } } },
			"0000000000000000000000000000000000000004": { "balance": "1", "builtin": { "name": "identity", "pricing": { "linear": { "base": 15, "word": 3 } } } },
			"0000000000000000000000000000000000000005": { "balance": "1", "builtin": { "name": "modexp", "activate_at": 0, "pricing": { "modexp": { "divisor": 20 } } } },
			"0000000000000000000000000000000000000006": {
				"balance": "1",
				"builtin": {
					"name": "alt_bn128_add",
					"pricing": {
						"0": {
							"price": { "alt_bn128_const_operations": { "price": 500 }}
						},
						"0x7fffffffffffff": {
							"info": "EIP 1108 transition",
							"price": { "alt_bn128_const_operations": { "price": 150 }}
						}
					}
				}
			},
			"0000000000000000000000000000000000000007": {
				"balance": "1",
				"builtin": {
					"name": "alt_bn128_mul",
					"pricing": {
						"0": {
							"price": { "alt_bn128_const_operations": { "price": 40000 }}
						},
						"0x7fffffffffffff": {
							"info": "EIP 1108 transition",
							"price": { "alt_bn128_const_operations": { "price": 6000 }}
						}
					}
				}
			},
			"0000000000000000000000000000000000000008": {
				"balance": "1",
				"builtin": {
					"name": "alt_bn128_pairing",
					"pricing": {
						"0": {
							"price": { "alt_bn128_pairing": { "base": 100000, "pair": 80000 }}
						},
						"0x7fffffffffffff": {
							"info": "EIP 1108 transition",
							"price": { "alt_bn128_pairing": { "base": 45000, "pair": 34000 }}
						}
					}
				}
			},
			"00a329c0648769a73afac7f9381e08fb43dbea72": { "balance": "1606938044258990275541962092341162602522202993782792835301376" }
		}
	}
	`

	// etc单元测试用这个
	// 修改自：https://github.com/eth-classic/mordor/blob/master/parity.json
	// 去除了nodes,engine改为instantSeal,networkID改为0x77,chainID改为0x4f，增加了默认有余额的地址
	ParityMordorVariantChainJSON = `
	{
		"name":"Local Mordor Classic Testnet",
		"dataDir":"mordor_private",
		"engine":{
			"instantSeal": {
				"params": {
				}
			}
		},
		"params":{
		   "gasLimitBoundDivisor":"0x400",
		   "accountStartNonce":"0x0",
		   "maximumExtraDataSize":"0x20",
		   "minGasLimit":"0x1388",
		   "networkID":"0x77",
		   "chainID":"0x4f",
		   "eip150Transition":"0x0",
		   "eip160Transition":"0x0",
		   "eip161abcTransition":"0x0",
		   "eip161dTransition":"0x0",
		   "eip155Transition":"0x0",
		   "maxCodeSize":"0x6000",
		   "maxCodeSizeTransition":"0x0",
		   "eip140Transition":"0x0",
		   "eip211Transition":"0x0",
		   "eip214Transition":"0x0",
		   "eip658Transition":"0x0",
		   "eip145Transition": "0x498bb",
		   "eip1014Transition": "0x498bb",
		   "eip1052Transition": "0x498bb"
		},
		"genesis":{
		   "seal":{
			  "ethereum":{
				 "nonce":"0x0000000000000000",
				 "mixHash":"0x0000000000000000000000000000000000000000000000000000000000000000"
			  }
		   },
		   "difficulty":"0x20000",
		   "author":"0x0000000000000000000000000000000000000000",
		   "timestamp":"0x5d9676db",
		   "parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000",
		   "extraData":"0x70686f656e697820636869636b656e206162737572642062616e616e61",
		   "gasLimit":"0x2fefd8"
		},
		"nodes":[
		],
		"accounts":{
		   "0x0000000000000000000000000000000000000001":{
			  "builtin":{
				 "name":"ecrecover",
				 "pricing":{
					"linear":{
					   "base":3000,
					   "word":0
					}
				 }
			  }
		   },
		   "0x0000000000000000000000000000000000000002":{
			  "builtin":{
				 "name":"sha256",
				 "pricing":{
					"linear":{
					   "base":60,
					   "word":12
					}
				 }
			  }
		   },
		   "0x0000000000000000000000000000000000000003":{
			  "builtin":{
				 "name":"ripemd160",
				 "pricing":{
					"linear":{
					   "base":600,
					   "word":120
					}
				 }
			  }
		   },
		   "0x0000000000000000000000000000000000000004":{
			  "builtin":{
				 "name":"identity",
				 "pricing":{
					"linear":{
					   "base":15,
					   "word":3
					}
				 }
			  }
		   },
		   "0x0000000000000000000000000000000000000005":{
			  "builtin":{
				 "activate_at":"0x0",
				 "name":"modexp",
				 "pricing":{
					"modexp":{
					   "divisor":20
					}
				 }
			  }
		   },
		   "0x0000000000000000000000000000000000000006":{
			  "builtin":{
				 "activate_at":"0x0",
				 "name":"alt_bn128_add",
				 "eip1108_transition": "0x7fffffffffffff",
				 "pricing":{
					"alt_bn128_const_operations": {
					   "price": 500,
					   "eip1108_transition_price": 150
					}
				 }
			  }
		   },
		   "0x0000000000000000000000000000000000000007":{
			  "builtin":{
				 "activate_at":"0x0",
				 "eip1108_transition": "0x7fffffffffffff",
				 "name":"alt_bn128_mul",
				 "pricing":{
					"alt_bn128_const_operations": {
					   "price": 40000,
					   "eip1108_transition_price": 6000
					}
				 }
			  }
		   },
		   "0x0000000000000000000000000000000000000008":{
			  "builtin":{
				 "activate_at":"0x0",
				 "eip1108_transition": "0x7fffffffffffff",
				 "name":"alt_bn128_pairing",
				 "pricing":{
					"alt_bn128_pairing": {
					   "base": 100000,
					   "pair": 80000,
					   "eip1108_transition_base": 45000,
					   "eip1108_transition_pair": 34000
					}
				 }
			  }
		   },
		   "00a329c0648769a73afac7f9381e08fb43dbea72": {"balance": "1606938044258990275541962092341162602522202993782792835301376"}
		}
	 }
	`
)
