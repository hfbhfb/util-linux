package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"syscall"

	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	whichnsenter string // 指出nsenter 命令位置【经过了 /proc -> /prod 改造】

	forkill string // 删除正在运行的命令

	targetip      string // 目标ip
	targetcommand string // 运行命令
)

func getNamespaceIP(pid int) (string, error) {
	// 构建 nsenter 命令，进入指定 PID 的网络命名空间
	cmd := exec.Command(whichnsenter, "--target", strconv.Itoa(pid), "--net", "sh", "-c", "ifconfig | grep 'inet ' | awk '{print $2}'")

	// 运行命令并获取输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("无法获取进程 %d 的 IP 地址: %v", pid, err)
	}
	// fmt.Println(string(output))
	// fmt.Println(pid)
	if strings.Contains(string(output), targetip) {
		return string(output), nil
	}

	return "", fmt.Errorf("进程 %v 没找到 ip", pid)
}

func getPPID(pid int) (int, error) {
	statusPath := fmt.Sprintf("/prod/%d/status", pid)
	file, err := os.Open(statusPath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PPid:") {
			fields := strings.Fields(line)
			ppid, err := strconv.Atoi(fields[1])
			if err != nil {
				return 0, err
			}
			return ppid, nil
		}
	}

	return 0, fmt.Errorf("PPid not found for pid %d", pid)
}

func commandWatchOutput(pid int, targetcommand string) {
	// 创建一个执行命令的 Cmd 对象
	// cmd := exec.Command("your-command", "arg1", "arg2")
	cmd := exec.Command("stdbuf", "-oL", whichnsenter, "--target", strconv.Itoa(pid), "--net", "sh", "-c", targetcommand)

	// 运行命令并获取输出
	// 获取命令的标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// 创建一个扫描器，逐行读取命令输出
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// 等待命令执行完毕
	if err := cmd.Wait(); err != nil {
		fmt.Println("Error waiting for command execution:", err)
	}

}

func fnForkKill(forkill string) {

	fmt.Println("need kill something:", forkill)

	// 打开 /proc 目录
	files, err := ioutil.ReadDir("/proc")
	if err != nil {
		fmt.Println("无法读取 /proc 目录:", err)
	}

	// 遍历 /proc 目录中的所有文件
	for _, file := range files {
		// 检查文件名是否是数字，即进程 ID
		if pid, err := strconv.Atoi(file.Name()); err == nil {
			// 尝试读取 /proc/<pid>/cmdline 文件来获取命令名称
			cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
			cmdline, err := ioutil.ReadFile(cmdlinePath)
			if err == nil && len(cmdline) > 0 {
				// 打印 PID 和对应的命令名称
				//fmt.Printf("PID: %d, Command: %s\n", pid, string(cmdline))
				// ppid, _ := getPPID(pid)
				if strings.Contains(string(cmdline), forkill) {

					fmt.Printf("PID: %d, Command: %s\n", pid, string(cmdline))

					// 查找进程
					process, err := os.FindProcess(pid)
					if err != nil {
						fmt.Println("无法找到进程:", err)
						return
					}

					// 向进程发送 SIGKILL 信号 (kill -9)
					err = process.Signal(syscall.SIGKILL)
					if err != nil {
						fmt.Println("无法发送信号:", err)
						return
					}

				}

			}
		}
	}

}

func runProcess() error {

	// 打开 /prod 目录
	files, err := ioutil.ReadDir("/prod")
	if err != nil {
		fmt.Println("无法读取 /prod 目录:", err)
		return err
	}

	// 遍历 /prod 目录中的所有文件
	for _, file := range files {
		// 检查文件名是否是数字，即进程 ID
		if pid, err := strconv.Atoi(file.Name()); err == nil {
			// 尝试读取 /prod/<pid>/cmdline 文件来获取命令名称
			cmdlinePath := fmt.Sprintf("/prod/%d/cmdline", pid)
			cmdline, err := ioutil.ReadFile(cmdlinePath)
			if err == nil && len(cmdline) > 0 {
				// 打印 PID 和对应的命令名称
				// fmt.Printf("PID: %d, Command: %s\n", pid, string(cmdline))
				// ppid, _ := getPPID(pid)
				if !strings.Contains(string(cmdline), "/pause") {

					// 找到对应进程匹配
					_, err := getNamespaceIP(pid)
					if err == nil {
						fmt.Printf("PID: %d, Command: %s\n", pid, string(cmdline))
						fmt.Println(pid)

						commandWatchOutput(pid, targetcommand)

						return nil
					}

				}

			}
		}
	}
	fmt.Printf("can not find target")

	return nil
}

// 获取环境变量，如果不存在则返回一个默认值
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// make all
// kubectl logs prod-opt1-0 -nhwx1166232
// kubectl exec -it prod-opt1-0 -nhwx1166232 -- bash
// f6
// kubectl get po -o wide -nhwx1166232
// kubectl exec -it nginx-helloworld-7dbb975b99-rkfhj -n hwx1166232 -- tcpdump -nnSX
// kubectl exec -it prod-opt1-0 -n hwx1166232 -- bash
// export targetip=10.0.230.176; export whichnsenter=nsenter-prod; export targetcommand='tcpdump -nnSX';proxy-prod -t 25328524
// nsenter-prod -t 1957 -n sh -c 'tcpdump -nnSX'
// stdbuf命令简介及其用途 https://developer.aliyun.com/article/1566962
// stdbuf是Linux系统中的一个实用工具命令，用于控制标准输入、输出和错误输出的缓冲模式。默认情况下，Linux中的I/O操作是缓冲的，这意味着数据在写入或读取时不会立即生效，而是先存储在内存中的缓冲区，直到满足特定条件（如缓冲区满、遇到换行符或文件结束符）才会被处理。stdbuf命令允许用户修改这种默认行为，以满足特定的数据处理和分析需求。
//
//	command: ["/bin/sh","-c"," export targetip=10.0.230.176; export whichnsenter=nsenter-prod; export targetcommand='tcpdump -nnSX ;date;sleep 2;date;sleep 2;ip a;' ; proxy-prod -t 25328524 ;while true; do  date; sleep 13; done"]
//
// export forkill=25328524 ;proxy-prod
func main() {

	// 功能，杀死正在运行的task 【超时】
	forkill = getEnv("forkill", "")
	if forkill != "" {
		fnForkKill(forkill)
		return
	}

	targetip = getEnv("targetip", "")
	fmt.Println("targetip: ", targetip)
	if targetip == "" {
		fmt.Println("targetip is empty!")
		return
	}

	whichnsenter = getEnv("whichnsenter", "nsenter") // 默认使用nsenter
	fmt.Println("whichnsenter: ", whichnsenter)

	targetcommand = getEnv("targetcommand", "")
	fmt.Println("targetcommand: ", targetcommand)
	if targetcommand == "" {
		fmt.Println("targetcommand is empty !")
		return
	}

	runProcess()

}
