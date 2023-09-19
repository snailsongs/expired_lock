package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

// GetCurrentProcessAndGoroutineID 获取进程ID+x协程ID 组成的二位标识字符串
func GetCurrentProcessAndGoroutineID() string {
	pid := GetCurrentProcessID()
	goroutineID := GetCurrentGoroutineID()
	return fmt.Sprintf("%d_%s", pid, goroutineID)
}

// GetCurrentGoroutineID 获取当前协程ID
func GetCurrentGoroutineID() string {
	buf := make([]byte, 128)
	buf = buf[:runtime.Stack(buf, false)]
	stackInfo := string(buf)
	return strings.TrimSpace(strings.Split(strings.Split(stackInfo, "[running]")[0], "goroutine ")[1])
}

// GetCurrentProcessID  获取当前进程ID
func GetCurrentProcessID() int {
	return os.Getpid()
}
