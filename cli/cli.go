package cli

import (
	"github.com/milansav/Castle/util"
	"os"
)

type CompilerSettings struct {
	Files  []string
	outdir string
}

func ParseArguments() CompilerSettings {
	argv := os.Args
	argc := len(argv)

	settings := CompilerSettings{
		Files: make([]string, 0),
	}

	for index := 0; index < argc; index++ {
		element := argv[index]

		if len(element) < 1 {
			os.Exit(1)
		}

		if util.GetRune(element, 0) == '-' {

			switch util.GetRune(element, 1) {
			case 'c':
				if index+1 >= argc {
					os.Exit(1)
				}
				fileName := argv[index+1]
				index++
				settings.Files = append(settings.Files, fileName)
			case 'd':
				if index+1 >= argc {
					os.Exit(1)
				}
				settings.outdir = argv[index+1]
				index++
			}
		}
	}

	return settings
}
