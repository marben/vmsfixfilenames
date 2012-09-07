package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var lowerCase = false
var upperCase = false
var force = false

func usage() {
	fmt.Println("This tool renames all the files in a directory that have the vms's file.ext;number pattern to a name without it..")
	fmt.Println("Usage: vmsfixfilenames <dir>")
	fmt.Println("Options: -l or -u to convert filenames to lower/upper case")
	fmt.Println("\t -f - force rename even if the filename is not vms style")
}

var vmsRegexp = regexp.MustCompile("^(.*\\..*);(\\d+)$")

func vmsFilename(filename string) bool {
	if vmsRegexp.MatchString(filename) {
		return true
	}

	return false
}

// converts vmsFilename to normal filename (removes tailing ;<num>)
// if normal filename specified, process it without changing extension
// if toupper/tolowe specified, converts the filename accordingly in all cases
func vmsFixFilename(vmsFilename string) (string, error) {
	var output string
	matches := vmsRegexp.FindStringSubmatch(vmsFilename)
	if matches == nil {
		if force {
			output = vmsFilename
		} else {
			return "", errors.New("Filename is not vms style..")
		}
	} else {
		output = matches[1]
	}

	if upperCase {
		output = strings.ToUpper(output)
	}

	if lowerCase {
		output = strings.ToLower(output)
	}

	return output, nil
}

func main() {
	flag.Usage = usage
	flag.BoolVar(&lowerCase, "l", false, "Convert filenames to lowercase")
	flag.BoolVar(&upperCase, "u", false, "Convert filenames to uppercase")
	flag.BoolVar(&force, "f", false, "Convert to lower/uppercase even if the file is not vms style")

	flag.Parse()

	if upperCase && lowerCase {
		fmt.Fprintln(os.Stderr, "Cannot use lower case flag together with uppercase flag")
		os.Exit(1)
	}

	args := flag.Args()
	if len(args) != 1 {
		usage()
		os.Exit(1)
	}

	dirname := args[0]

	dir, err := os.Open(dirname)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer dir.Close()

	files, err_readdir := dir.Readdir(0)
	if err_readdir != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

    newFiles := make(map[string]string)
	for _, file := range files {
		if vmsFilename(file.Name()) || force {
			newFilename, err := vmsFixFilename(file.Name())
			if err != nil {
				panic("panic!!")
				os.Exit(1)
			}

			fmt.Printf("%s -> %s\n", file.Name(), newFilename)
            newFiles[file.Name()] = newFilename
		}
	}

    if len(newFiles) == 0 {
        fmt.Println("No files to process. Exitting..")
        os.Exit(0)
    }


	fmt.Println("Do you want to continue? [Y/n]")
	var reply string
	fmt.Scanf("%s", &reply)
	reply = strings.ToLower(reply)
	if reply != "y" && reply != "yes" {
		fmt.Println("Quitting.")
		os.Exit(0)
	}
	for file, newFile := range newFiles {

        err_rename := os.Rename(filepath.Join(dirname, file), filepath.Join(dirname, newFile))
        if err_rename != nil {
            fmt.Fprint(os.Stderr, err_rename)
        }
	}
}
