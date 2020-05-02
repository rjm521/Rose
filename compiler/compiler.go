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

/**
compile the program in the supervisor sandbox
return the executable file path
 */

func Compile(compileConfig config.CompileConfig, srcPath string, outputDir string) (string, error)  {
	exePath := filepath.Join(outputDir, compileConfig.ExeName)
	replacements := map[string]string{
		//"{src_path}": srcPath,
		"{exe_dir}":  outputDir,
		"{exe_path}": exePath,
	}
	command := compileConfig.CompileCommand.FillWith(replacements)

	compilerOUt := filepath.Join(outputDir, "compiler.out")

	args := strings.Split(command, " ")

	//parse args

	result, err := supervisor.SupervisorRun(supervisor.Config{
		MaxCpuTime:       compileConfig.MaxCpuTime,
		MaxRealTime:      compileConfig.MaxRealTime,
		MaxMemory:        compileConfig.MaxMemory,
		MaxStack:         128 * 1024 * 1024,
		MaxProcessNumber: supervisor.UNLIMITED,
		MaxOutPutSize:    supervisor.UNLIMITED,
		SupervisorExePath: outputDir,
		ExePath:          args[0],
		InputPath:        srcPath,
		OutputPath:       compilerOUt,
		ErrorPath:        compilerOUt,
		LogPath:          config.COMPILER_LOG_PATH,
		Args:             args[1:],
		Env:              []string{"PATH=" + os.Getenv("PATH")},
		SecCompRuleName:  "none",
		Uid:              config.COMPILER_USER_UID,
		Gid:              config.COMPILER_GROUP_UID,
	})
	if err != nil {
		return "", err
	}

	if result.Result != supervisor.SUCCESS {
		_, err := os.Stat(compilerOUt)
		var errOut string
		if err == nil {
			errByte, _ := ioutil.ReadFile(compilerOUt)
			errOut = string(errByte[:])
			log.Println(errOut)
			os.Remove(compilerOUt)
		} else {
			errOut = fmt.Sprintf("Compiler Runtime Error, info %#v", result)

		}
		return exePath, errors.New(errOut)
	}

	return exePath, nil
}