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
	"os"
	"os/exec"
	"strings"
	"time"
)

//some getc related const
const (
	GetcCMD                    = "getc"
	GetcDefaultCoinbaseAddress = "0x49e826edacde7b646dae7312fe146e8fd2746565"
	GetcDefaultCoinbasePrivate = "55099a7ad52b17691edb35d8641de219f9d1684addfc72cec96e6144c984ac4e"
	GetcDefaultNetworkID       = "965"
	GetcDefaultRPCPort         = "8546"
)

// GetcDefaultArgs see the code
func GetcDefaultArgs() map[string]*string {
	return map[string]*string{
		"--dev":          nil,
		"--ipc-disable":  nil,
		"--mine":         nil,
		"--minerthreads": pstring("1"),
		"--network-id":   pstring(GetcDefaultNetworkID),
		"--rpc":          nil,
		"--rpc-port":     pstring(GetcDefaultRPCPort),
		"--port":         pstring("30000"),
		"--etherbase":    pstring(GetcDefaultCoinbaseAddress),
	}
}

// RunGetcOptions .
type RunGetcOptions struct {
	NewTmpDir             bool               //创建并使用新的临时目录作为datadir,和logdir
	TmpDirTag             string             //tag会作为临时路径的前缀
	NotRemoveTmpDirInKill bool               //kill func中不会移除临时目录，方便自行查看
	Args                  map[string]*string //k-v ,v 为nil时为flag
	NotPrint2stdout       bool               //不打印到stdout(cmd 的stdout不会指定到os.stdout)
}

// RunGetcNode run getc server,print out to stdout, require getc in the $PATH, this func is used for testing getc in local test env
// return func() to kill getc server
// usage:
// 		killGetc, err := RunGetcNode(options)
//  	defer killGetc()
func RunGetcNode(optionsPtr *RunGetcOptions) (func(), error) {
	if _, err := exec.LookPath(GetcCMD); err != nil {
		return nil, fmt.Errorf("getc may not in the path, %v", err)
	}

	killHooks := []killHook{}

	var options RunGetcOptions
	var err error

	if optionsPtr == nil {
		options = RunGetcOptions{}
	} else {
		options = *optionsPtr
	}
	if options.Args == nil {
		options.Args = map[string]*string{}
	}

	var dataDir string
	if options.NewTmpDir {
		for k, v := range options.Args {
			if k == "datadir" || k == "data-dir" {
				return nil, fmt.Errorf("datadir specified in args (%v), NewTmpDir not work", v)
			}
		}

		tag := strings.TrimLeft(options.TmpDirTag, "/")
		tmpDir := strings.TrimRight(os.TempDir(), "/")

		dataDir = tmpDir + "/" + tag + "getc_data_tmp_" + time.Now().Format(rfc3339Variant) + "/"
		err := os.MkdirAll(dataDir, 0777)
		if err != nil {
			return nil, fmt.Errorf("cannot create tmp dir: %v, err: %v", dataDir, err)
		}
		options.Args["--datadir"] = &dataDir
		options.Args["--log-dir"] = pstring(dataDir + "log")

		if options.NotRemoveTmpDirInKill {
		} else {
			killHooks = append(killHooks, func() error {
				return os.RemoveAll(dataDir)
			})
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

	cmd := exec.Command(GetcCMD, args...)
	fmt.Println("[debug] getc-node args", cmd.Args)
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
		fmt.Println("Waiting for message to kill getc")
		<-closeChan
		fmt.Println("Received message,killing getc server")

		if e := cmd.Process.Kill(); e != nil {
			fmt.Println("关闭 getc 时发生异常", e)
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
