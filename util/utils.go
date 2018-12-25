package util

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
)

func PrintExeFile(pid int) {
	output, err := exec.Command("/bin/readlink",
		fmt.Sprintf("/proc/%d/exe", pid)).Output()
	if err != nil {
		log.Warnf("failed to readlink the /proc/%d/exe", pid)
	}
	log.Debugf("the executable file of pid [%d] is %s", pid,
		strings.Trim(string(output), "\n"))
}

func FileOrDirExists(fileOrDir string) (bool, error) {
	_, err := os.Stat(fileOrDir)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

func EnSureFileExists(fileName string) error {
	dir, _ := path.Split(fileName)
	exist, err := FileOrDirExists(dir)
	if err != nil {
		return err
	}
	if !exist {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0644)
	defer file.Close()
	return err
}

func DirIsMounted(dir string) bool {
	_, err := exec.Command("sh", "-c", fmt.Sprintf("mount | grep -qw %s", dir)).Output()
	return err == nil
}

func GetEnvsByPid(pid int) ([]string, error) {
	envFile := fmt.Sprintf("/proc/%d/environ", pid)
	envsBytes, err := ioutil.ReadFile(envFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read envfile %s: %v", envFile, err)
	}
	return strings.Split(string(envsBytes), "\u0000"), nil
}

func RandomName(n uint32) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, n)
	for i := range b {
		if i%5 == 4 {
			b[i] = '_'
			continue
		}
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
