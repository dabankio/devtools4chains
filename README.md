# devtools4chains
Golang 公链开发工具。Dev tools for chains.Tools for integrating block chains into golang dev env.

典型的，开发DAPP 钱包等区块链生态中，在单元测试环境中，需要创建一个某条链的干净环境，测试完成后清理环境。

举个例子：我正在比特币钱包，那么开始前我需要部署一个bitcoind调试节点，节点数据需要是全新的，测试完成后关闭这个节点（必要的话还要清理数据），devtools4chains完成这项工作。

## 特性
- 启动某个程序，返回关闭函数
- 可定制启动参数
- 可选的创建临时目录
- 可选的跟踪打印日志或者将日志输出到某个文件
- 支持的程序
  - [ ] getc ,需要$PATH下有`getc`
  - [ ] bitcoind
  - [ ] omnicored
  - [ ] geth
  - [ ] parity, 需要$PATH下有`parity`
  - [ ] bigbang, 需要$PATH下有`bigbang`
  - [ ] ganache-cli, 需要$PATH下有`ganache-cli`

## 使用
`go get github.com/dabankio/devtools4chains`



# LICENSE

软件基于 [木兰宽松许可证](https://license.coscl.org.cn/MulanPSL/) 发行
