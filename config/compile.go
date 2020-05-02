package config

import (
	"RunCodeServer/supervisor"
	"strings"
)

type CommandStr string

func (c *CommandStr) FillWith(param map[string]string) string  {
	for k,v := range param {
		*c = CommandStr(strings.Replace(string(*c), k, v, -1))
	}
	return string(*c)
}

type CompileConfig struct {
	SrcName string //source code file name
	ExeName string // executable file name (after compile)
	MaxCpuTime int
	MaxRealTime int
	MaxMemory int
	CompileCommand CommandStr
}

type RunConfig struct {
	Command CommandStr
	StdinputPath string
	StdoutputPath string
	SeccompRule string
	Env []string
}

type LanguageCompileConfig struct {
	CompileConfig
	RunConfig
}

var DefaultEnv = []string{
	"LANG=en_US.UTF-8",
	"LANGUAGE=en_US:en",
	"LC_ALL=en_US.UTF-8",
}

var CompileC = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "main.c",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      supervisor.UNLIMITED,
		CompileCommand: "/usr/bin/gcc -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c99 main.c -lm -o main",
	},
	RunConfig:     RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "c_cpp",
		Env:           DefaultEnv,
	},
}

var CompileCpp = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "main.cpp",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      supervisor.UNLIMITED,
		CompileCommand: "/usr/bin/g++ -DONLINE_JUDGE -O2 -w -fmax-errors=3 -std=c++11 main.cpp -lm -o main",
	},
	RunConfig:     RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "c_cpp",
		Env:           DefaultEnv,
	},
}

var CompileJava = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "Main.java",
		ExeName:        "Main",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      supervisor.UNLIMITED,
		CompileCommand: "/usr/bin/javac Main.java -d Main -encoding UTF8",
	},
	RunConfig:     RunConfig{
		Command:     "/usr/bin/java -cp {exe_dir} -Xss1M -XX:MaxPermSize=16M -XX:PermSize=8M -Xms16M -Xmx{max_memory}k -Djava.security.manager -Dfile.encoding=UTF-8 -Djava.security.policy==/etc/java_policy -Djava.awt.headless=true Main",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "none",
		Env:            append(DefaultEnv, "MALLOC_ARENA_MAX=1"),
	},
}

var CompilePython2 = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "solution.py",
		ExeName:        "solution",
		MaxCpuTime:     3000,
		MaxRealTime:    5000,
		MaxMemory:      supervisor.UNLIMITED,
		CompileCommand: "/usr/bin/python -m py_compile solution.py",
	},
	RunConfig:     RunConfig{
		Command:       "/usr/bin/python {exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "general",
		Env:           DefaultEnv,
	},
}