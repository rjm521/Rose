package server

import (
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"net/http"
	"os"
	"runtime"
)

// Ping will offer information about this server.
func Ping(c *gin.Context) {
	name, err := os.Hostname()
	if err != nil {
		name = ""
	}
	cpu, _ := cpu.Percent(0, false)
	vmem, _ := mem.VirtualMemory()
	c.JSON(http.StatusOK, gin.H{
		"server_version": "rose",
		"name":           name,
		"cpu_core":       runtime.NumCPU(),
		"cpu":            cpu,
		"memory":         vmem.UsedPercent,
		"action":         "pong",
	})
}
