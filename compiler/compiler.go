package compiler

import (
	"RunCodeServer/config"
	"RunCodeServer/supervisor"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Compile(c config.CompileConfig, space string) (string, error) {

	exePath := filepath.Join(space, c.ExeName)
	srcCodePath := filepath.Join(space, c.SrcName)
	replacements := map[string]string{
		"{src_path}": srcCodePath,
		"{exe_dir}":  space,
		"{exe_path}": exePath,
	}

	// update compile args
	command := c.CompileCommand.FillWith(replacements)

	compilerOut := filepath.Join(space, "compiler.out")

	args := strings.Split(command, " ")

	//parse args
	result, err := supervisor.SupervisorRun(supervisor.Config{
		MaxCpuTime:        c.MaxCpuTime,
		MaxRealTime:       c.MaxRealTime,
		MaxMemory:         c.MaxMemory,
		MaxStack:          128 * 1024 * 1024,
		MaxProcessNumber:  supervisor.UNLIMITED,
		MaxOutPutSize:     supervisor.UNLIMITED,
		SupervisorExePath: space,
		ExePath:           args[0],
		InputPath:         srcCodePath,
		OutputPath:        compilerOut,
		ErrorPath:         compilerOut,
		LogPath:           config.COMPILER_LOG_PATH,
		Args:              args[1:],
		Env:               []string{"PATH=" + os.Getenv("PATH"), "GOPATH=/root/go", "HOME=/root"},
		SecCompRuleName:   "none",
		Uid:               config.COMPILER_USER_UID,
		Gid:               config.COMPILER_GROUP_UID,
	})
	if err != nil {
		return "", err
	}

	if result.Result != supervisor.SUCCESS {
		_, err := os.Stat(compilerOut)
		var errOut string
		if err == nil {
			errByte, _ := ioutil.ReadFile(compilerOut)
			errOut = string(errByte[:])
			log.Println(errOut)
			os.Remove(compilerOut)
		} else {
			errOut = fmt.Sprintf("Compiler Runtime Error, info %#v", result)
		}
		return exePath, errors.New(errOut)
	}

	return exePath, nil
}
