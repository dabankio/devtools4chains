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
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/dabankio/bbrpc"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func MustRunDockerMKFDev(t *testing.T, imageName string, autoClean bool, autoRemove bool) DockerCore {
	fn, info, err := DockerRunMKFDev(imageName, autoRemove)
	if err != nil {
		t.Fatalf("docker run image failed, %v", err)
		return info
	}
	if autoClean {
		t.Cleanup(fn)
	}
	return info
}

func DockerRunMKFDev(imageName string, autoRemove bool) (func(), DockerCore, error) {
	info := DockerCore{
		MinerAddress:   "20c003rgxdn4s64r4d0dchvb87p791q4epswkn1txadgv1evjqqwk97tv",
		MinerOwnerPubk: "3bc3e5f2e5e44f1cdbc44d3bf9325c93314be123f7563b8e6a88dc6eb1a25465",
		UnlockPass:     "123",
		Conf: bbrpc.ConnConfig{
			User:       "mkf",
			Pass:       "123",
			DisableTLS: true,
		},
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		return func() {}, info, err
	}
	idlePort, err := GetIdlePort()
	if err != nil {
		return func() {}, info, err
	}
	info.RPCPort = idlePort

	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		// AttachStderr: true,
		// AttachStdout: true,
		// Tty:          true,
		Image:        imageName,
		ExposedPorts: nat.PortSet{"9550": struct{}{}},
	}, &container.HostConfig{
		PortBindings:    nat.PortMap{"9550": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(idlePort)}}},
		PublishAllPorts: true,
		AutoRemove:      autoRemove,
	}, &network.NetworkingConfig{}, "")
	if err != nil {
		return func() {}, info, err
	}
	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return func() {}, info, err
	}
	log.Println("[info] mkf container started", cont.ID)

	stopContainer := func() {
		if er := cli.ContainerStop(context.Background(), cont.ID, nil); er != nil {
			log.Println("[warn] stop bbc core dev container err", er)
		}
	}

	info.Conf.Host = fmt.Sprintf("127.0.0.1:%d", idlePort)
	info.Client, err = bbrpc.NewClient(&info.Conf)
	if err != nil {
		return stopContainer, info, err
	}
	for { //wait for booting
		time.Sleep(time.Second)
		_, e := info.Client.Getforkheight(nil)
		if e == nil {
			break
		}
		log.Println("check alive:", e)
	}
	return stopContainer, info, nil

}
