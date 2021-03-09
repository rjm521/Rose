package client

import (
	"RunCodeServer/config"
	"RunCodeServer/supervisor"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type RunResult struct {
	CpuTime  int
	Result   supervisor.ResultCode
	RealTime int
	Memory   int
	Signal   int
	Error    supervisor.ErrorCode
	Output   string
}

type RunClient struct {
	Runconf       config.RunConfig
	ExePath       string
	MaxCpuTime    int
	MaxMemory     int
	SubmissionDir string
}

func (rc *RunClient) RunCode() (runResult RunResult) {
	commands := strings.Split(string(rc.Runconf.Command), " ")
	userOutputPath := filepath.Join(rc.SubmissionDir, rc.Runconf.StdoutputPath)
	userInputPath := filepath.Join(rc.SubmissionDir, rc.Runconf.StdinputPath)
	result, err := supervisor.SupervisorRun(supervisor.Config{
		MaxCpuTime:        rc.MaxCpuTime,
		MaxRealTime:       rc.MaxCpuTime * 3,
		MaxMemory:         supervisor.UNLIMITED,
		MaxStack:          128 * 1024 * 1024,
		MaxProcessNumber:  supervisor.UNLIMITED,
		MaxOutPutSize:     16 * 1024 * 1024,
		SupervisorExePath: rc.SubmissionDir,
		ExePath:           commands[0],
		InputPath:         userInputPath,
		OutputPath:        userOutputPath,
		ErrorPath:         userOutputPath,
		LogPath:           config.RUN_CODE_LOG_PATH,
		Args:              commands[1:],
		Env:               append(rc.Runconf.Env, "PATH="+os.Getenv("PATH")),
		SecCompRuleName:   rc.Runconf.SeccompRule,
		Uid:               config.RUN_USER_UID,
		Gid:               config.RUN_GROUP_UID,
	})

	fmt.Println("============== debug run exctuable  begin =====================")
	for _, v := range commands {
		fmt.Println(v)
	}
	fmt.Println("============== debug run exctuable end =====================")
	if err != nil {
		return RunResult{
			CpuTime:  result.CpuTime,
			Result:   -1,
			RealTime: result.RealTime,
			Memory:   result.Memory,
			Signal:   result.Signal,
			Error:    result.Error,
			Output:   "something wrong!",
		}
	}
	// fmt.Println(result)

	OutputContext, err := ioutil.ReadFile(userOutputPath)
	log.Println("用户的文件输出结果：----->" + string(OutputContext))

	if err != nil {
		OutputContext = []byte("something wrong!")
	}
	if result.Result != 0 {
		OutputContext = []byte("")
	}
	return RunResult{
		CpuTime: result.CpuTime,
		Result:  result.Result,
		Memory:  result.Memory,
		Signal:  result.Signal,
		Error:   result.Error,
		Output:  string(OutputContext),
	}

}
