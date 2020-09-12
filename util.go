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
	"encoding/json"
	"net"
	"time"
)

func pstring(s string) *string { return &s }

// JSONIndent marshal indent to string
func JSONIndent(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}

// GetIdlePort 随机获取一个空闲的端口
func GetIdlePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0") //当指定的端口为0时，操作系统会自动分配一个空闲的端口
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	time.AfterFunc(3*time.Second, func() { //端口在3秒后释放，避免多个请求同时发起取得相同的端口
		l.Close()
	})
	return l.Addr().(*net.TCPAddr).Port, nil
}
