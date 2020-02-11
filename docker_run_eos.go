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
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/pkg/errors"
)

const (
	defaultEOSNodeImage = "eostudio/eos:v2.0.0"
)

/*
docker run --name nodeos -d -p 8888:8888 -v ~/eos:/work \
-v ~/eos/data:/mnt/dev/data -v ~/eos/config:/mnt/dev/config eostudio/eos:v2.0.0 \
/bin/bash -c "nodeos -e -p eosio --plugin eosio::producer_plugin \
--plugin eosio::history_plugin --plugin eosio::chain_api_plugin \
--plugin eosio::history_api_plugin --plugin eosio::http_plugin \
-d /mnt/dev/data --config-dir /mnt/dev/config --http-server-address=0.0.0.0:8888 \
--access-control-allow-origin=* --contracts-console --http-validate-host=false"
*/

// DockerRunOptions .
type DockerRunOptions struct {
	AutoRemove bool
	Image      *string //默认 defaultEOSNodeImage
}

// DockerRunNodeosOptions .
type DockerRunNodeosOptions struct {
	DockerRunOptions
}

// DockerRunNodeos docker run eos nodeos
// default port:8888, else 16000+
// ~/.dockereos
// 使用固定的目录，所以不支持并发执行
func DockerRunNodeos(opt *DockerRunNodeosOptions) (KillFunc, *DockerContainerInfo, error) {
	if opt == nil {
		opt = &DockerRunNodeosOptions{}
	}
	// TODO 查找空闲的端口并绑定

	var workPath, dataPath, configPath string
	{ //目录处理
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return nothing2do, nil, err
		}
		volumeDir := filepath.Join(userHomeDir, ".dockereosdata")
		if _, err = os.Stat(volumeDir); err == nil { //dir exists, remove it
			err = os.RemoveAll(volumeDir)
			if err != nil {
				return nothing2do, nil, errors.Wrap(err, "remove data dir failed")
			}
		}
		for _, subpath := range []string{"data", "config"} { //2个子目录建立
			err = os.MkdirAll(filepath.Join(volumeDir, subpath), 0755)
			if err != nil {
				return nothing2do, nil, err
			}
		}
		workPath = volumeDir
		dataPath = filepath.Join(volumeDir, "data")
		configPath = filepath.Join(volumeDir, "config")
	}

	cli, err := client.NewEnvClient()
	if err != nil {
		return nothing2do, nil, err
	}

	if opt.Image == nil {
		opt.Image = pstring(defaultEOSNodeImage)
	}
	err = dockerIsImageExists(cli, *opt.Image)
	if err != nil {
		return nothing2do, nil, err
	}

	hostConfig := &container.HostConfig{
		// "8888/tcp": [{"HostIp": "","HostPort": "8888"}]
		PortBindings:    nat.PortMap{"8888": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8888"}}},
		PublishAllPorts: true,
		Mounts: []mount.Mount{ //可用binds
			{Type: "bind", Target: "/work", Source: workPath},
			{Type: "bind", Target: "/mnt/dev/data", Source: dataPath},
			{Type: "bind", Target: "/mnt/dev/config", Source: configPath},
		},
	}
	if opt.AutoRemove {
		hostConfig.AutoRemove = true
	}
	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		// AttachStderr: true,
		// AttachStdout: true,
		// Tty:          true,
		Image: *opt.Image,
		Cmd: []string{
			"nodeos",
			"-e",
			"--producer-name=eosio",
			"--plugin=eosio::producer_plugin",
			"--plugin=eosio::history_plugin",
			"--plugin=eosio::chain_api_plugin",
			"--plugin=eosio::history_api_plugin",
			"--plugin=eosio::http_plugin",
			"--data-dir=/mnt/dev/data",
			"--config-dir=/mnt/dev/config",
			"--http-server-address=0.0.0.0:8888",
			"--access-control-allow-origin=*",
			"--contracts-console",
			"--http-validate-host=false",
		},
		ExposedPorts: nat.PortSet{"8888": struct{}{}},
	}, hostConfig, &network.NetworkingConfig{}, "")
	if err != nil {
		return nothing2do, nil, err
	}

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return nothing2do, nil, err
	}
	log.Printf("nodeos container [%s] started\n", *opt.Image)

	return func() {
			log.Printf("[info] stop container: %s (autoRemove: %v)\n", *opt.Image, opt.AutoRemove)
			if e := cli.ContainerStop(context.Background(), cont.ID, nil); e != nil {
				log.Println("[Err] stop container error", e)
			}
			log.Println("[info] container stopped")
		}, &DockerContainerInfo{
			ListenPorts: []int{8888}, //TODO fix ports
		}, nil
}
