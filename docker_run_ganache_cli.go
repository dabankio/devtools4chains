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
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// DockerRunGanacheCli run ganache-cli from docker,
// require: docker started daemon
func DockerRunGanacheCli(opt *RunDockerContOptions) (func(), error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nothing2do, err
	}

	// imgs, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	// if err != nil {
	// 	return nothing2do, err
	// }
	// fmt.Println(imgs)

	// _, err = cli.ImagePull(context.Background(), "trufflesuite/ganache-cli:latest", types.ImagePullOptions{})
	// if err != nil {
	// 	return nothing2do, err
	// }

	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
		Image:        "trufflesuite/ganache-cli:latest",
	}, &container.HostConfig{}, &network.NetworkingConfig{}, "")
	if err != nil {
		return nothing2do, err
	}

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return nothing2do, err
	}
	if opt != nil && opt.Log2std {
		go followPrintContainerLog(cli, cont.ID)
	}

	// fmt.Println("试试日志打印")
	// time.Sleep(time.Second * 8)
	// fmt.Println("试试日志打印after sleep")

	return func() {
		log.Println("[info] stop and remove ganache-cli container")
		e := cli.ContainerStop(context.Background(), cont.ID, nil)
		if e != nil {
			log.Println("[Err] stop container error", e)
		}
		log.Println("[info] container stopped")

		e = cli.ContainerRemove(context.Background(), cont.ID, types.ContainerRemoveOptions{RemoveVolumes: true})
		if e != nil {
			log.Println("[Err] failed to remove container")
		}
		log.Println("[info] container removed")
	}, nil

	// cmd := exec.Command("docker", "run", "-p", "8545:8545", "trufflesuite/ganache-cli:latest")
	// cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout
	// err := cmd.Start()
	// if err != nil {
	// return nothing2do, err
	// }
	// return func() {
	// e := cmd.Process.Kill()
	// if e != nil {
	// log.Println("kill docker ganache-cli error")
	// }
	// }, nil
}