package devtools4chains

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-pg/pg/v9"
	"go.uber.org/zap"
)

// some docker const
const (
	DockerPGPassword = "pwd"
	DockerPGUser     = "postgres"
	DockerPGDatabase = "postgres"
)

// DockerPGInfo 运行的docker pg容器信息
type DockerPGInfo struct {
	Port          int
	User          string
	Password      string
	Database      string
	ContainerName string
}

// RunTestInDockerPG .
func RunTestInDockerPG(testFn func(info DockerPGInfo)) error {
	stopPG, info, err := DockerRunPG(DockerRunOptions{
		AutoRemove: true,
	})
	if err != nil {
		return err
	}
	defer stopPG()
	testFn(info)
	return nil
}

// MustRunPG 启动一个docker pg供测试使用，自动注册停止pg 函数
// ops[0]将生效（如果有）
// 无法创建数据库时 t.Fatal
func MustRunPG(t *testing.T, ops ...DockerRunOptions) DockerPGInfo {
	if len(ops) == 0 {
		ops = []DockerRunOptions{{AutoRemove: true}}
	}
	stopPG, info, err := DockerRunPG(ops[0])
	if err != nil {
		t.Fatal("run docker pg failed", err)
	} else {
		t.Cleanup(stopPG)
	}
	return info
}

// DockerRunPG run 一个pg container,构造一个空的pg database,通常用以执行单元测试，确保空库
// EXPOSE 端口是随机的，这样可以进行并行的测试
// 可以选择 autoRemove container
// 返回：kill函数（停止pg容器）,pgInfo 数据库信息, error
// 使用：
// stopPG, info, err := testtool.DockerRunPG(testtool.DockerRunOptions{})
// if err != nil {...}
// defer stopPG()
//
// 建议：AutoRemove=true，除非你需要保留容器观察
// 为了方便调试创建的容器会给一个名字  pg_{端口}_{时间}，例如：pg_59587_20200204T134858
// 启动后可以通过这些方式psql:
//    psql -h 127.0.0.1 -p {port} -U postgres -W
//    docker exec -it pg_{port}_{containerName} psql
func DockerRunPG(opt DockerRunOptions) (func(), DockerPGInfo, error) {
	info := DockerPGInfo{
		User:     DockerPGUser,
		Password: DockerPGPassword,
		Database: DockerPGDatabase,
	}
	cli, err := client.NewEnvClient()
	if err != nil {
		return nothing2do, info, err
	}

	if opt.Image == nil {
		_o := "postgres:latest"
		opt.Image = &_o
	}
	if err = dockerIsImageExists(cli, *opt.Image); err != nil {
		return nothing2do, info, err
	}

	idlePort, err := GetIdlePort()
	if err != nil {
		return nothing2do, info, err
	}
	info.Port = idlePort
	info.ContainerName = fmt.Sprintf("pg_%d_%s", info.Port, time.Now().Format("20060102T150405"))
	cont, err := cli.ContainerCreate(context.Background(), &container.Config{
		// AttachStderr: true,
		// AttachStdout: true,
		// Tty:          true,
		Env:          []string{fmt.Sprintf("POSTGRES_PASSWORD=%s", DockerPGPassword)},
		Image:        *opt.Image,
		ExposedPorts: nat.PortSet{"5432": struct{}{}},
	}, &container.HostConfig{
		PortBindings:    nat.PortMap{"5432": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: strconv.Itoa(info.Port)}}},
		PublishAllPorts: true,
		AutoRemove:      opt.AutoRemove,
	}, &network.NetworkingConfig{}, info.ContainerName)
	if err != nil {
		return nothing2do, info, err
	}

	err = cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	if err != nil {
		return nothing2do, info, err
	}
	zap.S().Infof("container [%s] started (listen: %d, name: %s, id: %s)\n", *opt.Image, info.Port, info.ContainerName, cont.ID[:12])

	for { //启动后不能立刻就能连接，要等一会儿，测试连接，通了再返回
		time.Sleep(time.Second)
		zap.S().Info("ping docker pg")
		connected := func() bool {
			db := pg.Connect(&pg.Options{
				Database: info.Database,
				User:     info.User,
				Password: info.Password,
				Addr:     fmt.Sprintf("127.0.0.1:%d", info.Port),
			})
			defer db.Close()
			_, err := db.Exec("SELECT 1")
			return err == nil
		}()
		if connected {
			break
		}
	}

	return func() {
		zap.S().Infof("[info] stop container: %s (name: %s ,autoRemove: %v)\n", *opt.Image, info.ContainerName, opt.AutoRemove)
		if e := cli.ContainerStop(context.Background(), cont.ID, nil); e != nil {
			zap.S().Warn("stop container error", e)
		}
	}, info, nil
}

// ToPGOption convert to *pg.Options
func (info DockerPGInfo) ToPGOption() *pg.Options {
	return &pg.Options{
		ApplicationName: "unit testing",
		Database:        info.Database,
		User:            info.User,
		Password:        info.Password,
		Addr:            fmt.Sprintf("127.0.0.1:%d", info.Port),
	}
}
