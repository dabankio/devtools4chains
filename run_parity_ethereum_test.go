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

func TestRunParityEthereum(t *testing.T) {
	killParity, err := RunParityEthereum(&RunParityEthereumConfig{
		DataDir:   DataDirOption{NewTmpDir: true, NotRemoveTmpDirWhenKilling: true},
		ChainJSON: ParityMordorVariantChainJSON,
	})

	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer killParity()

	time.Sleep(time.Second * 10)
}
