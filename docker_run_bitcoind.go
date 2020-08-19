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
	"context"
	"log"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// const
const (
	DockerBitcoinImage = "ruimarinho/bitcoin-core:0.18"
)

// DockerBitcoinInfo .
type DockerBitcoinInfo struct {
	Docker  DockerContainerInfo
	RPCPort int
	RPCUser string
	RPCPwd  string
}

/**
docker run -d -p 18443:18443 kylemanna/bitcoind \
    -regtest -txindex \
	-server -rpcuser=usr -rpcpassword=pwd

docker exec -it cocky_allen bitcoin-cli -rpcport=18443 getblockchaininfo

curl --data-binary '{"jsonrpc":"1.0","id":"curltext","method":"getblockchaininfo","params":[]}' -H 'content-type:application/json;' http://rpcusr:233@127.0.0.1:18443/
*/

// DockerRunBitcoin 。
func DockerRunBitcoin(opt DockerRunOptions) (KillFunc, *DockerBitcoinInfo, error) {
	const (
		rpcUser, rpcPwd = "rpcusr", "233"
	)
	idlePort, err := GetIdlePort()
	if err != nil {
		return func() {}, nil, err
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		return nothing2do, nil, err
	}

	if opt.Image == nil {
		opt.Image = pstring(DockerBitcoinImage)
	}
	err = dockerIsImageExists(cli, *opt.Image)
	if err != nil {
		return nothing2do, nil, err
	}

	hostConfig := &container.HostConfig{
		AutoRemove:      opt.AutoRemove,
		PortBindings:    nat.PortMap{"18443": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(idlePort)}}},
		PublishAllPorts: true,
		Mounts:          []mount.Mount{ //可用binds
			// {Type: "bind", Target: "/work", Source: workPath},
		},
	}
	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		Image: *opt.Image,
		Cmd: []string{
			"-regtest",
			"-txindex",
			"-server",
			"-rpcallowip=0.0.0.0/0",
			"-rpcbind=0.0.0.0", //default 127.0.0.1, it does not work in docker publish port out
			"-rpcauth=rpcusr:656f9dabc62f0eb697c801369617dc60$422d7fca742d4a59460f941dc9247c782558367edcbf1cd790b2b7ff5624fc1b", //rpcusr:233
			"-rpcport=18443",
			"-fallbackfee=0.000002",
		},
		ExposedPorts: nat.PortSet{"18443": struct{}{}},
	}, hostConfig, &network.NetworkingConfig{}, "")
	if err != nil {
		return nothing2do, nil, err
	}

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return nothing2do, nil, err
	}
	log.Printf("container [%s] started\n", *opt.Image)

	return func() {
			log.Printf("[info] stop container: %s (autoRemove: %v)\n", *opt.Image, opt.AutoRemove)
			if e := cli.ContainerStop(context.Background(), cont.ID, nil); e != nil {
				log.Println("[Err] stop container error", e)
			}
			log.Println("[info] container stopped")
		}, &DockerBitcoinInfo{
			Docker: DockerContainerInfo{
				ListenPorts: []int{idlePort}, //TODO fix ports
			},
			RPCUser: rpcUser,
			RPCPort: idlePort,
			RPCPwd:  rpcPwd,
		}, nil
}
