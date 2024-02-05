//go:build windows
// +build windows

package sys

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"unsafe"

	"github.com/hoshinonyaruko/palworld-go/config"
	 "github.com/shirou/gopsutil/process"
)

// GetPIDByPath 返回与给定路径匹配的进程的PID
func GetPIDByPath(path string) (int32, error) {
    processes, err := process.Processes()
    if err != nil {
        return 0, err
    }

    for _, p := range processes {
        exePath, err := p.Exe()
        if err == nil && exePath == path {
            return p.Pid, nil
        }
    }
    return 0, fmt.Errorf("no process found with path %s", path)
}

// WindowsRestarter implements the Restarter interface for Windows systems.
type WindowsRestarter struct{}

// NewRestarter creates a new Restarter appropriate for Windows systems.
func NewRestarter() *WindowsRestarter {
	return &WindowsRestarter{}
}

func (r *WindowsRestarter) Restart(executablePath string) error {
	// Separate the directory and the executable name
	execDir, execName := filepath.Split(executablePath)

	// Including -faststart parameter in the script that starts the executable
	scriptContent := "@echo off\n" +
		"pushd " + strconv.Quote(execDir) + "\n" +
		// Add the -faststart parameter here
		"start \"\" " + strconv.Quote(execName) + " -faststart\n" +
		"popd\n"

	scriptName := "restart.bat"
	if err := os.WriteFile(scriptName, []byte(scriptContent), 0755); err != nil {
		return err
	}

	cmd := exec.Command("cmd.exe", "/C", scriptName)

	if err := cmd.Start(); err != nil {
		return err
	}

	// The current process can now exit
	os.Exit(0)

	// This return statement will never be reached
	return nil
}

// windows
func setConsoleTitleWindows(title string) error {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return err
	}
	proc, err := kernel32.FindProc("SetConsoleTitleW")
	if err != nil {
		return err
	}
	p0, err := syscall.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	r1, _, err := proc.Call(uintptr(unsafe.Pointer(p0)))
	if r1 == 0 {
		return err
	}
	return nil
}

func KillProcess() error {
    // 获取当前工作目录
    cwd, err := os.Getwd()
    if err != nil {
        return err
    }

    // 构造PalServer的完整路径
    executablePath := filepath.Join(cwd, "PalServer-Win64-Test-Cmd.exe")

    // 获取PID
    pid, err := GetPIDByPath(executablePath)
    if err != nil {
        return err
    }

    var cmd *exec.Cmd
    if runtime.GOOS == "windows" {
        // Windows: 使用PID结束进程
        cmd = exec.Command("taskkill", "/PID", strconv.Itoa(int(pid)), "/F")
    } else {
        // 非Windows: 使用PID结束进程
        cmd = exec.Command("kill", "-9", strconv.Itoa(int(pid)))
    }

    cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
    return cmd.Run()
}

// RunViaBatch 函数接受配置，程序路径和参数数组
func RunViaBatch(config config.Config, exepath string, args []string) error {
	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// 创建批处理脚本内容
	batchScript := `@echo off
	start "" "` + exepath + `" ` + strings.Join(args, " ")

	// 指定批处理文件的路径
	batchFilePath := filepath.Join(cwd, "run_command.bat")

	// 写入批处理脚本到文件
	err = os.WriteFile(batchFilePath, []byte(batchScript), 0644)
	if err != nil {
		return err
	}

	// 执行批处理脚本
	cmd := exec.Command(batchFilePath)
	cmd.Dir = config.GamePath // 设置工作目录为游戏路径
	return cmd.Run()
}
