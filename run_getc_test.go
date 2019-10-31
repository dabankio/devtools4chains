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
	"testing"
	"time"
)

func TestRunGetcNode(t *testing.T) {
	args := GetcDefaultArgs()
	args["--rpc-port"] = pstring("38545")
	killGetc, err := RunGetcNode(&RunGetcOptions{
		NewTmpDir:             true,
		NotRemoveTmpDirInKill: true,
		NotPrint2stdout:       true,
		Args:                  args,
	})
	tShouldNil(t, err)
	defer killGetc()
	time.Sleep(time.Second * 10)
}
