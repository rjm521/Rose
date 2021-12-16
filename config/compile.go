package config

import (
	"RunCodeServer/supervisor"
	"strings"
)

type CommandStr string

func (c *CommandStr) FillWith(param map[string]string) string {
	for k, v := range param {
		*c = CommandStr(strings.Replace(string(*c), k, v, -1))
	}
	return string(*c)
}

type CompileConfig struct {
	SrcName        string //source code file name
	ExeName        string // executable file name (after compile)
	MaxCpuTime     int
	MaxRealTime    int
	MaxMemory      int
	CompileCommand CommandStr
}

type RunConfig struct {
	Command       CommandStr
	StdinputPath  string
	StdoutputPath string
	SeccompRule   string
	Env           []string
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
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/gcc -O2 -w -fmax-errors=3 -std=c99 main.c -lm -o main",
	},
	RunConfig: RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "c_cpp",
		Env:           DefaultEnv,
	},
}

var CompileC11 = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "main.c",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/gcc -O2 -w -fmax-errors=3 -std=c11 main.c -lm -o main",
	},
	RunConfig: RunConfig{
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
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/g++ -O2 -w -fmax-errors=3 -std=c++11 main.cpp -lm -o main",
	},
	RunConfig: RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "c_cpp",
		Env:           DefaultEnv,
	},
}

var CompileCpp14 = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "main.cpp",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/g++ -O2 -w -fmax-errors=3 -std=c++14 main.cpp -lm -o main",
	},
	RunConfig: RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "c_cpp",
		Env:           DefaultEnv,
	},
}

var CompileCpp17 = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "main.cpp",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/g++ -O2 -w -fmax-errors=3 -std=c++17 main.cpp -lm -o main",
	},
	RunConfig: RunConfig{
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
		MaxRealTime:    10000,
		MaxMemory:      supervisor.UNLIMITED,
		CompileCommand: "/usr/bin/javac {src_path} -d {exe_dir} -encoding UTF8",
	},
	RunConfig: RunConfig{
		Command:       "/usr/bin/java -cp {exe_dir} -XX:MaxRAM={max_memory}k -Djava.security.manager -Dfile.encoding=UTF-8 -Djava.security.policy==/etc/java_policy -Djava.awt.headless=true Main",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "none",
		Env:           DefaultEnv,
	},
}

var CompilePython2 = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "solution.py",
		ExeName:        "solution.pyc",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/python -m py_compile {src_path}",
	},
	RunConfig: RunConfig{
		Command:       "/usr/bin/python {exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "general",
		Env:           DefaultEnv,
	},
}

var CompilePython3 = LanguageCompileConfig{
	CompileConfig: CompileConfig{
		SrcName:        "solution.py",
		ExeName:        "__pycache__/solution.cpython-37.pyc",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/python3.7 -m py_compile {src_path}",
	},
	RunConfig: RunConfig{
		Command:       "/usr/bin/python3 {exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "general",
		Env:           append(DefaultEnv, "PYTHONIOENCODING=UTF-8", "MALLOC_ARENA_MAX=1"),
	},
}

var CompilePhp7 = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "solution.php",
		ExeName:        "solution.php",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/bin/cat {src_path}",
	},
	RunConfig{
		Command:       "/usr/bin/php -d error_reporting=0 -f {exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "none",
		Env:           DefaultEnv,
	},
}

var CompileJsc = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "solution.js",
		ExeName:        "solution.js",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/bin/cat {src_path}",
	},
	RunConfig{
		Command:       "/usr/bin/jsc {exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "none",
		Env:           DefaultEnv,
	},
}

var CompileGo = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "main.go",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      supervisor.UNLIMITED,
		CompileCommand: "/usr/local/go/bin/go build -o main main.go",
	},
	RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "none",
		Env:           append(DefaultEnv, "GODEBUG=madvdontneed=1", "GOCACHE=off"),
	},
}

var CompileCsharp = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "main.cs",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/mcs -optimize+ build -out:{main} main.cs",
	},
	RunConfig{
		Command:       "/usr/bin/mono {exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "none",
		Env:           DefaultEnv,
	},
}

var CompileRuby = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "solution.rb",
		ExeName:        "solution.rb",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/bin/cat solution.rb",
	},
	RunConfig{
		Command:       "/usr/bin/ruby {exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "none",
		Env:           DefaultEnv,
	},
}

var CompileRust = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "main.rs",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/rustc -O -o main main.rs",
	},
	RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "general",
		Env:           DefaultEnv,
	},
}

var CompileHaskell = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "main.hs",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/ghc -O -outputdir /tmp -o main main.hs",
	},
	RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "general",
		Env:           DefaultEnv,
	},
}

var CompilePascal = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "main.pas",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/bin/ghc -O -outputdir /tmp -o main main.hs",
	},
	RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "general",
		Env:           DefaultEnv,
	},
}

var CompilePlainText = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "solution.txt",
		ExeName:        "solution.txt",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/bin/cat solution.txt",
	},
	RunConfig{
		Command:       "/bin/cat {exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "general",
		Env:           DefaultEnv,
	},
}

var CompileBasic = LanguageCompileConfig{
	CompileConfig{
		SrcName:        "main.bas",
		ExeName:        "main",
		MaxCpuTime:     3000,
		MaxRealTime:    10000,
		MaxMemory:      1024 * 1024 * 1024,
		CompileCommand: "/usr/local/bin/fbc main.bas",
	},
	RunConfig{
		Command:       "{exe_path}",
		StdinputPath:  "1.in",
		StdoutputPath: "1.out",
		SeccompRule:   "general",
		Env:           DefaultEnv,
	},
}
