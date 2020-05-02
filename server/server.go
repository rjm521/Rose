package server

import (
	"RunCodeServer/client"
	"RunCodeServer/compiler"
	"RunCodeServer/config"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	validator "github.com/matthewgao/gojsonvalidator"
	uuid "github.com/satori/go.uuid"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

type UserInput struct {
	Code string `json:"code" is_required:"true"`
	StdInput string `json:"std_input" is_required:"true"`
	Language string `json:"language" is_required:"true"`
	MaxCpuTime int `json:"max_cpu_time" is_required:"true"`
	MaxMemory int `json:"max_memory" is_required:"true"`
}
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
func Ping(c *gin.Context)  {
	if !checkToken(c.GetHeader("Token")) {
		c.JSON(http.StatusOK,gin.H{
			"code": -1,
			"data": "wrong token",
		})
		return
	}
	hostname, err  := os.Hostname()
	if err != nil {
		hostname = ""
	}
	cpuPercent, _ := cpu.Percent(0, false)
	vmem, _ := mem.VirtualMemory()
	c.JSON(http.StatusOK, gin.H{
		"server_version": "0.1.0",
		"hostname":       hostname,
		"cpu_core":       runtime.NumCPU(),
		"cpu": cpuPercent,
		"memory": vmem.UsedPercent,
		"action": "pong",
	})
}
func HandleCode(c *gin.Context) {

	// Token验证模块，暂时不想要
	//if !checkToken(c.GetHeader("Token")) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"code": -1,
	//		"data": "wrong token",
	//	})
	//	return
	//}

	var input UserInput
	var err error

	//改变协议，获得websocket连接
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	for  {

		//读消息
		_, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}

		//判断消息格式，消息格式不正确返回结束
		err = validator.ValidateJson(msg,&input)
		if err != nil {
			err = sendToClient(ws,gin.H{
				"code": -1,
				"data": err.Error(),})
			if err != nil {
				break
			}
			break
		}

		//为提交创建目录,有异常断开连接
		subbmissionId := uuid.NewV4()
		var submissionDir string
		if submissionDir , err = initSubmissionEnv(subbmissionId.String()); err != nil {
			err = sendToClient(ws, gin.H{
				"code": -1,
				"data": err.Error(),
			})
			break
		}
		defer freeSubmissionEnv(submissionDir)

		var conf config.LanguageCompileConfig
		switch input.Language {
		case "c":
				conf = config.CompileC
		case "cpp":
				conf = config.CompileCpp
		case "java":
			conf = config.CompileJava
		case "py2":
			conf = config.CompilePython2
		default:
			err = sendToClient(ws, gin.H{
			"code": -1,
			"data":"invaild language support",
			})
			if err != nil {
				break
			}
		}

		//为提交代码和stdin创建文件
		srcPath := filepath.Join(submissionDir, conf.SrcName)
		var srcFile *os.File
		if srcFile, err = os.Create(srcPath); err != nil {
			err = sendToClient(ws,gin.H{
				"code": -1,
				"data": err.Error(),
			})
			if err != nil {
				break
			}
			break
		}
		defer srcFile.Close()
		stdInPath := filepath.Join(submissionDir, conf.RunConfig.StdinputPath)
		var stdInFile *os.File
		if stdInFile, err = os.Create(stdInPath); err != nil {
			err = sendToClient(ws, gin.H{
				"code": -1,
				"data": err.Error(),
			})
			break
		}
	_, err =os.Create(filepath.Join(submissionDir,conf.RunConfig.StdoutputPath))
		if err != nil {
			err = sendToClient(ws, gin.H{
				"code": -1,
				"data": err.Error(),
			})
			break
		}

		//defer stdInFile.Close()

		_, err = srcFile.WriteString(input.Code)
		if err != nil {
			err = sendToClient(ws,gin.H{
				"code": -1,
				"data": err.Error(),
			})
			if err != nil {
				break
			}
			break
		}
		_, err = stdInFile.WriteString(input.StdInput)
		if err != nil {
			err = sendToClient(ws,gin.H{
				"code": -1,
				"data": err.Error(),
			})
			if err != nil {
				break
			}
			break
		}

		 err = sendToClient(ws, gin.H{
			"code": 0,
			"status": "pending",
		})

		 if err != nil {
		 	break
		 }

		//编译代码中
		var exePath string
		exePath, err = compiler.Compile(conf.CompileConfig, srcPath, submissionDir)
		//代码编译错误
		if err != nil {
			err = sendToClient(ws, gin.H{
				"code": 0,
				"status": "COMPILE_ERROR",
				"compilation_log": err.Error(),
			})
			if err != nil {
				break
			}
			os.RemoveAll(submissionDir)
			continue
		}

		//编译代码完成更新状态
		err = sendToClient(ws, gin.H{
			"code": 0,
			"status": "running",
		})
		if err != nil {
			break
		}

		log.Println("用户可执行文件路径：--->" + exePath)
		log.Println("用户目录：--->" + submissionDir)
		conf.RunConfig.Command.FillWith(map[string]string{
			"{exe_path}":   exePath,
			"{exe_dir}":    submissionDir,
			"{max_memory}": strconv.Itoa(input.MaxMemory / 1024),
		})

		rc := client.RunClient{
			Runconf:       conf.RunConfig,
			ExePath:       exePath,
			MaxCpuTime:    input.MaxCpuTime,
			MaxMemory:     input.MaxMemory,
			SubmissionDir: submissionDir,
		}
		runResult := rc.RunCode()
		var resultStatus string
		switch runResult.Result {
		case 0:
			resultStatus = "FINISHED"
		case 1, 2:
			resultStatus = "TIME_LIMIT_EXCEEDED"
		case 3:
			resultStatus = "MEMORY_LIMIT_EXCEEDED"
		case 4:
			resultStatus = "RUNTIME_ERROR"
		default:
			resultStatus = "SYSTEM_ERROR"
		}
		fmt.Println("input cpu time: ", input.MaxCpuTime)
		fmt.Println("status: ", resultStatus)
		fmt.Println("cpu time: ", runResult.CpuTime)
		fmt.Println("real time: ", runResult.RealTime)

		//更新代码执行情况
		err = sendToClient(ws, gin.H{
			"code":   0,
			"status": resultStatus,
			"stdout": runResult.Output,
			"time":   runResult.CpuTime,
		})
		os.RemoveAll(submissionDir)
		if err != nil {
			break
		}
	}
}

func sendToClient(conn *websocket.Conn, v interface{}) error  {
	msg, err := json.Marshal(v);
	if err != nil {
		return err
	}
	err = conn.WriteMessage(1, msg)
	if err != nil {
		return err
	}
	return nil
}

func freeSubmissionEnv(submissionDirPath string) error {
	return os.RemoveAll(submissionDirPath)
}

func initSubmissionEnv(subbmissionId string) (string, error) {
	subbmissionDirPath := filepath.Join(config.SUBMISSION_DIR, subbmissionId)
	if err := os.Mkdir(subbmissionDirPath, 0777); err != nil {
		return "", err

	}
	return subbmissionDirPath, nil
}
func checkToken(token string)  bool {
	localToken := os.Getenv("RPC_TOKEN")
	if token != localToken {
		return false
	}
	return true
}

