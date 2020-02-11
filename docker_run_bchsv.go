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

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

//bsv const
const (
	DockerBSVImage = "bitcoinsv/bitcoin-sv:1.0.1"
)

// DockerBSVInfo .
type DockerBSVInfo struct {
	Docker  DockerContainerInfo
	RPCPort int
	RPCUser string
	RPCPwd  string
}

/**
--name bsvreg -p 18444:18444 bitcoinsv/bitcoin-sv:1.0.1 bitcoind \
    -regtest -txindex -excessiveblocksize=0 -maxstackmemoryusageconsensus=0 \
    -server -rpcuser=usr -rpcpassword=pwd
*/

// DockerRunBSV 。
func DockerRunBSV(opt *DockerRunOptions) (KillFunc, *DockerBSVInfo, error) {
	const (
		port            = 18444
		rpcUser, rpcPwd = "usr", "pwd"
	)

	if opt == nil {
		opt = &DockerRunOptions{}
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		return nothing2do, nil, err
	}

	if opt.Image == nil {
		opt.Image = pstring(DockerBSVImage)
	}
	err = dockerIsImageExists(cli, *opt.Image)
	if err != nil {
		return nothing2do, nil, err
	}

	hostConfig := &container.HostConfig{
		AutoRemove: opt.AutoRemove,
		// "8888/tcp": [{"HostIp": "","HostPort": "8888"}]
		PortBindings:    nat.PortMap{"18444": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "18444"}}},
		PublishAllPorts: true,
		Mounts:          []mount.Mount{ //可用binds
			// {Type: "bind", Target: "/work", Source: workPath},
			// {Type: "bind", Target: "/mnt/dev/data", Source: dataPath},
			// {Type: "bind", Target: "/mnt/dev/config", Source: configPath},
		},
	}
	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		// AttachStderr: true,
		// AttachStdout: true,
		// Tty:          true,
		Image: *opt.Image,
		Cmd: []string{
			"bitcoind",
			"-regtest",
			"-txindex",
			"-excessiveblocksize=0",
			"-maxstackmemoryusageconsensus=0",
			"-server",
			"-rpcuser=" + rpcUser,
			"-rpcpassword=" + rpcPwd,
			"-rpcport=18444",
		},
		ExposedPorts: nat.PortSet{"18444": struct{}{}},
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
		}, &DockerBSVInfo{
			Docker: DockerContainerInfo{
				ListenPorts: []int{port}, //TODO fix ports
			},
			RPCUser: rpcUser,
			RPCPort: port,
			RPCPwd:  rpcPwd,
		}, nil
}
