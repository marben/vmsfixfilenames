package main

import (
    "fmt"
    "flag"
    "os"
    "regexp"
    )


func usage() {
    fmt.Println("This tool renames all the files in a directory that have the vms's file.ext;number pattern to a name without it..")
    fmt.Println("Usage: vmsfixfilenames <dir>")
}

var vmsRegExp = regexp.MustCompile("^.*\\..*;\\d+$")

func vmsPattern(filename string) bool {
    vmsRegexp  := regexp.MustCompile("^.*\\..*;\\d+$")
    if vmsRegexp.MatchString(filename) {
        return true
    }

    return false
}

func main() {
    var lowerCase bool
    var upperCase bool

    flag.Usage = usage
    flag.BoolVar(&lowerCase, "l", false, "Convert filenames to lowercase")
    flag.BoolVar(&upperCase, "u", false, "Convert filenames to uppercase")

    flag.Parse()

    if upperCase && lowerCase {
        fmt.Fprintln(os.Stderr, "Cannot user lower case flag together with uppercase flag")
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

    files, err_readdir := dir.Readdir(0)
    if err_readdir != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }

    for i, file := range files {
        var sign string
        if vmsPattern(file.Name()){
            sign = "VMS"
        } else {
            sign = ""
        }

        fmt.Println(i, " - ", file.Name(), " - ", sign)
    }
}
