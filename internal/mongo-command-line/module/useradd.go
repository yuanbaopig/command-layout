package module

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	NoLogin = "/sbin/nologin"
)

func SetUserOwner(uName, path string) error {
	targetUser, err := user.Lookup(uName)
	if err != nil {
		return fmt.Errorf("error retrieving targetUser: %v", err)
	}

	// 将 UID 和 GID 转换为整数
	uid, err := strconv.Atoi(targetUser.Uid)
	if err != nil {
		return fmt.Errorf("UID conversion failed: %v", err)
	}

	gid, err := strconv.Atoi(targetUser.Gid)
	if err != nil {
		return fmt.Errorf("GID conversion failed: %v", err)
	}

	if err := RecursivePermissionSet(path, uid, gid); err != nil {
		return fmt.Errorf("data path permission set failed: %v", err)
	}

	return nil
}

// RecursivePermissionSet 递归遍历目录并设置权限和属主属组
func RecursivePermissionSet(root string, uid, gid int) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			//log.Error(err)
			//fmt.Printf("Skipping %s: %v\n", path, err)
			return err // 跳过无法访问的文件或目录
		}

		if err := os.Chown(path, uid, gid); err != nil {
			//log.Error(err)
			return err
		}

		return nil
	})
}

// AddNoLoginUser create linux process user
func AddNoLoginUser(username, shell string) error {

	cmdStr := fmt.Sprintf("useradd -r -s %s %s", shell, username)

	args := strings.Fields(cmdStr)

	runCmd := exec.Command(args[0], args[1:]...)

	//bf := bufferpool.GetBytesBuffer()
	//defer bufferpool.PutBytesBuffer(bf)

	//runCmd.Stdout = bf
	//runCmd.Stderr = bf

	//if err := runCmd.Run(); err != nil {
	if outPut, err := runCmd.CombinedOutput(); err != nil {
		//if strings.Contains(bf.String(), fmt.Sprintf("useradd: user '%s' already exists", username)) {
		if strings.Contains(string(outPut), fmt.Sprintf("useradd: user '%s' already exists", username)) {
			return nil
		}

		return fmt.Errorf("%v %s", err, outPut)
	}

	return nil
}
