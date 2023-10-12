package cvdocker

import (
	_ "github.com/clickvisual/clickvisual/api/internal/pkg/cvdocker/containerd"
	_ "github.com/clickvisual/clickvisual/api/internal/pkg/cvdocker/docker"
	"github.com/clickvisual/clickvisual/api/internal/pkg/cvdocker/manager"
	"github.com/clickvisual/clickvisual/api/internal/pkg/utils"
	"github.com/gotomicro/ego/core/elog"
)

const containerdSock = "/run/containerd/containerd.sock"
const dockerSock = "/var/run/docker.sock"

// Component 组件
type Component struct {
	config       *manager.Config
	containerMap map[string]*manager.DockerInfo
	clientType   string // docker, containerd
}

func NewContainer() *Component {
	obj := &Component{
		config: manager.DefaultConfig(),
	}
	isExistDockerSock, err := utils.PathExist(dockerSock)
	if err != nil {
		elog.Panic("docker sock fail", elog.FieldErr(err))
	}

	isExistContainerdSock, err := utils.PathExist(containerdSock)
	if err != nil {
		elog.Panic("containerd sock fail", elog.FieldErr(err))
	}

	if !isExistDockerSock && !isExistContainerdSock {
		elog.Panic("docker.sock and containerd.sock is empty", elog.FieldValue("we need "+containerdSock+" or "+dockerSock))
	}
	if isExistContainerdSock {
		obj.config.ClientSocket = containerdSock
		obj.clientType = "containerd"
	} else if isExistDockerSock {
		obj.config.ClientSocket = dockerSock
		obj.clientType = "docker"
	}
	return obj
}

func (c *Component) GetActiveContainers() (containerMap map[string]*manager.DockerInfo) {
	var err error
	obj := manager.Get(c.clientType)
	obj.Run(c.config)
	containerMap, err = obj.GetAllDockerInfo()
	if err != nil {
		elog.Panic("containerdSockObj fetchAll fail", elog.FieldErr(err))
	}
	return
}
