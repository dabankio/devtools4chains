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
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// RunDockerContOptions run docker container options
type RunDockerContOptions struct {
	Log2std bool
}

// DockerContainerInfo .
type DockerContainerInfo struct {
	ListenPorts []int
}

func followPrintContainerLog(cli *client.Client, id string) {
	reader, err := cli.ContainerLogs(context.Background(), id, types.ContainerLogsOptions{
		ShowStdout: true, ShowStderr: true, Follow: true, Timestamps: false,
	})
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func dockerIsImageExists(cli *client.Client, name string) error {
	summaries, err := cli.ImageList(context.Background(), types.ImageListOptions{All: true})
	if err != nil {
		return err
	}
	for _, x := range summaries {
		// fmt.Printf("%#v\n", x)
		for _, tag := range x.RepoTags {
			if tag == name {
				return nil
			}
		}
	}
	return fmt.Errorf("image [%s] not exists, pull it manually ==> docker pull %s", name, name)
}
