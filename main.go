package main

import (
    "fmt"
    "flag"
    "os"
    "regexp"
    "errors"
    "strings"
    "path/filepath"
    )


func usage() {
    fmt.Println("This tool renames all the files in a directory that have the vms's file.ext;number pattern to a name without it..")
    fmt.Println("Usage: vmsfixfilenames <dir>")
    fmt.Println("Options: -l or -u to convert filenames to lower/upper case")
}

var vmsRegexp = regexp.MustCompile("^(.*\\..*);(\\d+)$")

func vmsFilename(filename string) bool {
    if vmsRegexp.MatchString(filename) {
        return true
    }

    return false
}

func vmsFixFilename(vmsFilename string) (string, error) {
    matches := vmsRegexp.FindStringSubmatch(vmsFilename)
    if matches == nil {
        return "", errors.New("Filename is not vms style..")
    }
    return matches[1], nil
}

func main() {
    var lowerCase bool
    var upperCase bool

    flag.Usage = usage
    flag.BoolVar(&lowerCase, "l", false, "Convert filenames to lowercase")
    flag.BoolVar(&upperCase, "u", false, "Convert filenames to uppercase")

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

    for _, file := range files {
        if vmsFilename(file.Name()){
            newFilename, err := vmsFixFilename(file.Name())
            if err != nil {
                panic("panic!!")
                os.Exit(1)
            }
            if upperCase {
                newFilename = strings.ToUpper(newFilename)
            }
            if lowerCase {
                newFilename = strings.ToLower(newFilename)
            }

            fmt.Printf("%s -> %s\n", file.Name(), newFilename)

            err_rename := os.Rename(filepath.Join(dirname, file.Name()), filepath.Join(dirname, newFilename))
            if err_rename != nil {
                fmt.Fprint(os.Stderr, err_rename)
            }
        }
    }
}
