package container

import "path"

const (
	Mydocker    = "mydocker"
	MyDockerDir = "/var/lib/mydocker"
	ConfigName  = "config.json"
	LogName     = "container.log"
)

const (
	ContainerPid = "ContainerPid"
	ContainerCmd = "ContainerCmd"
)

const (
	Creating = "creating"
	Running  = "running"
	Stopped  = "stopped"
	Exited   = "exited"
)

const (
	Stop    = "stop"
	Start   = "start"
	Restart = "restart"
	Create  = "create"
	Delete  = "delete"
)

var (
	WriteLayterDir = path.Join(MyDockerDir, "writelayer")
	ContainersDir  = path.Join(MyDockerDir, "containers")
)
