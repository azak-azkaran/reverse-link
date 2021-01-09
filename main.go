package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

const (
	ERROR_PARSE     = "ERROR: while parsing: "
	ERROR_READ_FLAG = "ERROR: while retrieving: "
	ERROR_READ_FILE = "ERROR: while reading File: "
	ERROR_PATH_DIR  = "ERROR: while reading Path is Directory"
	ERROR_NOT_LINK  = "ERROR: file is not a link"
	FLAG_UNSET ="_*NOT_FLAG_SET*_"
)

func ReadFlags(args []string) (string, error) {
	fileFlag := pflag.NewFlagSet("file", pflag.ContinueOnError)
	fileFlag.String("file", FLAG_UNSET, "the file for which the link should be reversed")
	err := fileFlag.Parse(args)
	if err != nil {
		log.Println(ERROR_PARSE, err)
		return "", err
	}

	fileName, err := fileFlag.GetString("file")
	if err != nil {
		log.Println(ERROR_READ_FLAG, err)
		return "", err
	}

	if fileName == "" {
		err =errors.New("Empty flag")
		log.Println(ERROR_READ_FLAG, err)
		return "", err
	}

	if fileName == FLAG_UNSET  && len(args) >= 1 {
		return args[len(args)-1], nil
	}

	return fileName, nil
} 

func CheckFile(fileName string)(string,error){
	log.Println("Cheaking file", fileName)

	stat, err := os.Stat(fileName)
	if err != nil {
		log.Println(ERROR_READ_FILE, err)
		return "", err
	}

	lstat, err := os.Lstat(fileName)
	if err != nil {
		log.Println(ERROR_READ_FILE, err)
		return "", err
	}

	if stat.IsDir() || lstat.IsDir() {
		log.Println(ERROR_PATH_DIR)
		return "", errors.New(ERROR_PATH_DIR)
	}

	if lstat.Mode()&os.ModeSymlink != os.ModeSymlink {
		log.Println(ERROR_NOT_LINK)
		return "", errors.New(ERROR_NOT_LINK)
	}
	return fileName, nil
}

func ReverseFile(path string) error {
	p, err := filepath.EvalSymlinks(path)
	if err != nil {
		return err
	}

	err = os.Rename(p, path)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	log.Println(os.Args)
	fileName, err := ReadFlags(os.Args)
	if err != nil {
		panic(err.Error())
	}

	fileName, err = CheckFile(fileName)
	if err != nil {
		panic(err.Error())
	}

	err = ReverseFile(fileName)
	if err != nil {
		panic(err.Error())
	}
}
