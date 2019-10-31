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
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
)

func tShouldNil(t *testing.T, v interface{}, args ...interface{}) {
	if v != nil {
		debug.PrintStack()
		t.Fatalf("[test assert] should nil, but got: %v, %v", v, args)
	}
}

func tShouldTrue(t *testing.T, b bool, args ...interface{}) {
	if !b {
		debug.PrintStack()
		t.Fatalf("[test assert] should true, args: %v", args)
	}
}

func tShouldNotZero(t *testing.T, v interface{}, args ...interface{}) {
	if reflect.ValueOf(v).IsZero() {
		debug.PrintStack()
		t.Fatalf("[test assert] should not [zero value], %v", args)
	}
}

func tShouldNotContains(t *testing.T, v, containV string) {
	if strings.Contains(v, containV) {
		debug.PrintStack()
		t.Fatalf("[test assert] [%s] should not contains [%s]", v, containV)
	}
}
