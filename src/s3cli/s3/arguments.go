package s3

import (
	"errors"
	"flag"
	"os/user"
	"path/filepath"
)

func validateArguments(nonFlagArgs []string) (cmd string, src string, dst string, err error) {
	emptyStr := ""

	if len(nonFlagArgs) != 3 {
		return emptyStr, emptyStr, emptyStr, errors.New("Invalid number of arguments")
	}

	cmd, src, dst = nonFlagArgs[0], nonFlagArgs[1], nonFlagArgs[2]

	if cmd != "get" && cmd != "put" {
		return emptyStr, emptyStr, emptyStr, errors.New("Invalid command: " + cmd)
	}

	return cmd, src, dst, nil
}

func fetchConfigurationPath(args []string) (path string, nonFlagArgs []string, err error) {
	flagSet := flag.NewFlagSet("args", flag.PanicOnError)
	flagSet.StringVar(&path, "c", "", "Config file")
	flagSet.Parse(args)
	if path != "" {
		return path, flagSet.Args(), nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", []string{}, err
	}
	return filepath.Join(usr.HomeDir, ".s3cli"), flagSet.Args(), nil
}
