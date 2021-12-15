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
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type UserInput struct {
	Code       string `json:"code" is_required:"true"`
	StdInput   string `json:"std_input" is_required:"true"`
	Language   string `json:"language" is_required:"true"`
	MaxCpuTime int    `json:"max_cpu_time" is_required:"true"`
	MaxMemory  int    `json:"max_memory" is_required:"true"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleCode(c *gin.Context) {
	// upgrade http to websocket connection
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		sendToClient(ws, gin.H{
			"code": -1,
			"data": err.Error()})
		return
	}
	defer ws.Close()
	handleConn(ws)
}

func handleConn(ws *websocket.Conn) {

	req := &UserInput{}
	for {

		// read msg from connection
		_, msg, err := ws.ReadMessage()
		if err != nil {
			return
		}

		// validate json
		err = validator.ValidateJson(msg, &req)
		if err != nil {
			sendToClient(ws, gin.H{
				"code": -1,
				"data": err.Error()})
			return
		}

		// validate language
		conf := config.LanguageCompileConfig{}
		switch req.Language {
		case "c":
			conf = config.CompileC
		case "c11":
			conf = config.CompileC11
		case "cpp":
			conf = config.CompileCpp
		case "cpp14":
			conf = config.CompileCpp14
		case "cpp17":
			conf = config.CompileCpp17
		case "java":
			conf = config.CompileJava
		case "py2":
			conf = config.CompilePython2
		case "py3":
			conf = config.CompilePython3
		case "php":
			conf = config.CompilePhp7
		case "javascript":
			conf = config.CompileJsc
		case "golang":
			conf = config.CompileGo
		case "csharp":
			conf = config.CompileCsharp
		case "ruby":
			conf = config.CompileRuby
		case "rust":
			conf = config.CompileRust
		case "haskell":
			conf = config.CompileHaskell
		case "pascal":
			conf = config.CompilePascal
		case "plaintext":
			conf = config.CompilePlainText
		case "basic":
			conf = config.CompileBasic
		default:
			sendToClient(ws, gin.H{
				"code": -1,
				"data": "invaild language support",
			})
			return
		}

		// create dir
		space, err := allocateSpace(req, conf)
		// TODO: defer freeSpace(space)
		if err != nil {
			sendToClient(ws, gin.H{
				"code": -1,
				"data": err.Error(),
			})
			return
		}

		// update status
		sendToClient(ws, gin.H{
			"code":   0,
			"status": "pending",
		})

		// compile code
		exePath, err := compiler.Compile(conf.CompileConfig, space)
		if err != nil {
			sendToClient(ws, gin.H{
				"code":            0,
				"status":          "COMPILE_ERROR",
				"compilation_log": err.Error(),
			})
			freeSpace(space)
			continue
		}

		// update status
		sendToClient(ws, gin.H{
			"code":   0,
			"status": "running",
		})

		log.Println("user exe file:", exePath)
		log.Println("user space:", space)

		conf.RunConfig.Command.FillWith(map[string]string{
			"{exe_path}":   exePath,
			"{exe_dir}":    space,
			"{max_memory}": strconv.Itoa(10240),
		})

		fmt.Println("server layer: ", conf.RunConfig.Command)
		rc := client.RunClient{
			Runconf:       conf.RunConfig,
			ExePath:       exePath,
			MaxCpuTime:    5000,
			MaxMemory:     128 * 1024 * 1024,
			SubmissionDir: space,
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

		// update status
		err = sendToClient(ws, gin.H{
			"code":   0,
			"status": resultStatus,
			"stdout": runResult.Output,
			"time":   runResult.CpuTime,
		})

		freeSpace(space)
		if err != nil {
			break
		}
	}
}

func allocateSpace(input *UserInput, c config.LanguageCompileConfig) (string, error) {

	// create directory
	id := uuid.NewV4()
	dir := filepath.Join(config.SUBMISSION_DIR, id.String())
	err := os.Mkdir(dir, 0777)
	if err != nil {
		return "", fmt.Errorf("could not create dir %v", err)
	}

	// create srouce code file
	srcPath := filepath.Join(dir, c.SrcName)
	srcFile, err := os.Create(srcPath)
	if err != nil {
		return "", fmt.Errorf("could not create srouce code file %v", err)
	}
	_, err = srcFile.WriteString(input.Code)
	if err != nil {
		return "", fmt.Errorf("could not write srouce code content %v", err)
	}

	// create stdin file
	stdInputPath := filepath.Join(dir, c.RunConfig.StdinputPath)
	stdInputFile, err := os.Create(stdInputPath)
	if err != nil {
		return "", fmt.Errorf("could not create stdin file %v", err)
	}
	_, err = stdInputFile.WriteString(input.StdInput)
	if err != nil {
		return "", fmt.Errorf("could not write stdin content %v", err)
	}

	// create stdout file
	stdOutputPath := filepath.Join(dir, c.RunConfig.StdoutputPath)
	_, err = os.Create(stdOutputPath)
	if err != nil {
		return "", fmt.Errorf("could not create stdout file %v", err)
	}

	return dir, nil
}

func sendToClient(conn *websocket.Conn, v interface{}) error {
	msg, err := json.Marshal(v)
	if err != nil {
		return err
	}
	err = conn.WriteMessage(1, msg)
	if err != nil {
		return err
	}
	return nil
}

func freeSpace(submissionDirPath string) {
	os.RemoveAll(submissionDirPath)
}
