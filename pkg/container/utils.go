package container

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/weikeit/mydocker/util"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func SendInitCommand(cmds []string, writePipe *os.File) {
	msgStr := strings.Join(cmds, " ")
	log.Debugf("runCommand sends user-defined command: %s", msgStr)
	writePipe.WriteString(msgStr)
	writePipe.Close()
}

func GetContainer(uuid string) (*Container, error) {
	configFile := path.Join(ContainersDir, uuid, ConfigName)
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("the configFile %s doesn't exist: %v", configFile, err)
	}

	contents, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read configFile %s: %v", configFile, err)
	}

	var c Container
	if err := json.Unmarshal(contents, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

func GetAllContainers() ([]*Container, error) {
	exist, _ := util.FileOrDirExists(ContainersDir)
	if ! exist {
		if err := os.MkdirAll(ContainersDir, 0777); err != nil {
			return nil, fmt.Errorf("failed to mkdir %s: %v", ContainersDir, err)
		}
		return nil, nil
	}

	containerDirs, err := ioutil.ReadDir(ContainersDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir %s: %v", ContainersDir, err)
	}

	var containerArray []*Container
	for _, containerDir := range containerDirs {
		uuid := containerDir.Name()
		c, err := GetContainer(uuid)
		if err != nil {
			log.Errorf("failed to get the info of container %s: %v", uuid, err)
			continue
		}
		containerArray = append(containerArray, c)
	}

	return containerArray, nil
}

func GetContainerByNameOrUuid(identifier string) (*Container, error) {
	allContainers, err := GetAllContainers()
	if err != nil {
		return nil, err
	}

	var c *Container
	for _, container := range allContainers {
		if identifier == container.Name || identifier == container.Uuid {
			c = container
			break
		}
	}

	if c == nil {
		return nil, fmt.Errorf("no such container: %s", identifier)
	}

	return c, nil
}
