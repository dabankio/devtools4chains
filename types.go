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

//kill 时调用的函数
type killHook func() error

// DataDirOption .
type DataDirOption struct {
	NewTmpDir bool //创建并使用一个新的临时目录

	//下面的选项在不使用临时目录时无效
	TmpDirPrefix               string //如果创建临时目录，使用该值作为目录前缀
	NotRemoveTmpDirWhenKilling bool   //在kill函数中不移除临时目录
}
