package main

import (
	"DatabaseManage/internal/pkg/filecheck"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	mongo  = "mongod"
	mongoS = "mongos"
	mysql  = "mysqld"
	redis  = "redis-server"
	pgsql  = "postmaster"
)

var TypeList = []string{
	mongo,
	mongoS,
	mysql,
	redis,
	pgsql,
}

type Process struct {
	Pid   int   `json:"pid"`
	Ports []int `json:"ports"`
}

type ScanRequest struct {
	Status bool                 `json:"status"`
	Type   map[string][]Process `json:"type"`
}

// InodeToPort 是可以存储网络信息的映射
type InodeToPort map[string]int

func getProcessesByName(name string) ([]int, error) {
	var pids []int
	procDirs, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^\d+$`) // 正则匹配只包含数字的字符串

	for _, procDir := range procDirs {

		if !procDir.IsDir() || !re.MatchString(procDir.Name()) {
			continue
		}

		pid, _ := strconv.Atoi(procDir.Name()) // 由于使用正则表达式匹配，这里不会出错

		cmdline, err := os.ReadFile(fmt.Sprintf("/proc/%d/comm", pid))
		if err != nil {
			continue
		}

		if strings.TrimSpace(string(cmdline)) == name {
			pids = append(pids, pid)
		}

	}

	return pids, nil
}

// ReadNetworkInfo 读取/proc/net/tcp和/proc/net/tcp6文件，以构建网络信息的映射。
func ReadNetworkInfo() (InodeToPort, error) {
	// 检查所有网络类型，这里是tcp和tcp6
	networkTypes := []string{"tcp", "tcp6"}
	inodeToPort := make(InodeToPort)

	for _, network := range networkTypes {

		tcpFile := fmt.Sprintf("/proc/net/%s", network)

		// 如果文件不存在，则直跳过
		if err := filecheck.CheckFileOrDirExist(tcpFile); err != nil {
			continue
		}

		content, err := os.ReadFile(tcpFile)
		if err != nil {
			return nil, err
		}

		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			fields := strings.Fields(line)
			if len(fields) > 9 && fields[3] == "0A" { // st 字段为 "0A"（LISTEN状态）
				inode := fields[9]
				localAddress := fields[1]
				portHex := strings.Split(localAddress, ":")[1]

				port, err := strconv.ParseInt(portHex, 16, 64)
				if err != nil {
					continue // 解析错误，跳过该行
				}

				inodeToPort[inode] = int(port)
			}
		}
	}

	return inodeToPort, nil
}

func getListeningPortsByPID(pid int, inodeToPort InodeToPort) ([]int, error) {
	var ports []int

	netFiles, err := os.ReadDir(fmt.Sprintf("/proc/%d/fd", pid))
	if err != nil {
		return nil, err
	}

	// 遍历文件描述符找到匹配的inode
	for _, netFile := range netFiles {
		link, err := os.Readlink(fmt.Sprintf("/proc/%d/fd/%s", pid, netFile.Name()))
		if err != nil {
			continue
		}

		if strings.HasPrefix(link, "socket:") {
			inode := strings.Trim(link, "socket:[")
			inode = strings.TrimSuffix(inode, "]")

			if port, ok := inodeToPort[inode]; ok {
				ports = append(ports, port)
			}
		}
	}

	return removeSliceDuplicates(ports), nil
}

func removeSliceDuplicates(s []int) []int {
	seen := make(map[int]struct{}) // struct类型不占用空间
	var result []int
	for _, val := range s {
		if _, ok := seen[val]; !ok {
			seen[val] = struct{}{}
			result = append(result, val)
		}
	}
	return result

}

func main() {
	// 获取mongod进程列表

	// 一次性读取网络信息，并存储在内存中
	inodeToPort, err := ReadNetworkInfo()
	if err != nil {
		fmt.Printf("Error reading network information: %v\n", err)
		return
	}

	result := ScanRequest{
		Type: make(map[string][]Process),
	}

	for _, dbType := range TypeList {

		var dbList []Process

		PidList, err := getProcessesByName(dbType)
		if err != nil {
			fmt.Printf("Error finding %s processes: %s \n", dbType, err)
			return
		}

		for _, pid := range PidList {

			ports, err := getListeningPortsByPID(pid, inodeToPort)
			if err != nil {
				fmt.Printf("Error finding ports for PID %d: %v\n", pid, err)
				continue
			}

			if len(ports) == 0 {
				continue
			}

			process := Process{
				Pid:   pid,
				Ports: ports,
			}

			dbList = append(dbList, process)
			result.Status = true
			result.Type[dbType] = dbList

		}
	}

	s, _ := json.Marshal(result)
	fmt.Println(string(s))
}
