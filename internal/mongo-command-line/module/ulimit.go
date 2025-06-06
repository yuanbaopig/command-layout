package module

import (
	"fmt"
	"syscall"
)

func ULimit() error {

	/*
		[mongod@db0 ~]$ ulimit -a
		core file size          (blocks, -c) 0
		data seg size           (kbytes, -d) unlimited
		scheduling priority             (-e) 0
		file size               (blocks, -f) unlimited
		pending signals                 (-i) 14989
		max locked memory       (kbytes, -l) unlimited
		max memory size         (kbytes, -m) unlimited
		open files                      (-n) 655350
		pipe size            (512 bytes, -p) 8
		POSIX message queues     (bytes, -q) 819200
		real-time priority              (-r) 0
		stack size              (kbytes, -s) 8192
		cpu time               (seconds, -t) unlimited
		max user processes              (-u) 327675
		virtual memory          (kbytes, -v) unlimited
		file locks                      (-x) unlimited
	*/

	// 设置open files限制
	var limit syscall.Rlimit
	// 设置最大打开文件数
	limit.Max = 655350
	limit.Cur = 655350

	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &limit); err != nil {
		fmt.Printf("设置文件描述符限制失败: %v\n", err)
		return err
	}

	// 设置最大用户进程数
	limit.Cur = 327675
	limit.Max = 327675
	if err := syscall.Setrlimit(syscall.PRIO_PROCESS, &limit); err != nil {
		fmt.Printf("设置用户进程限制失败: %v\n", err)
		return err
	}

	// 设置堆栈大小限制
	limit.Cur = 8192 * 1024 // 8192 KB
	limit.Max = 8192 * 1024
	if err := syscall.Setrlimit(syscall.RLIMIT_STACK, &limit); err != nil {
		fmt.Printf("设置堆栈大小限制失败: %v\n", err)
	} else {
		fmt.Println("堆栈大小限制设置成功")
	}
	return nil
}
