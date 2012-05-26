package main

import (
    "fmt"
    "flag"
    )

func main() {
    fmt.Println("This is a tool to fix vms ;256 style end of files")

    var lowerCaseOpt *bool = flag.Bool("l", false, "Convert filenames to lowercase")
    var upperCaseOpt *bool = flag.Bool("u", false, "Convert filenames to uppercase")

    flag.Parse()
    lowerCase := *lowerCaseOpt
    upperCase := *upperCaseOpt

    if lowerCase == true {
        fmt.Println("LowerCase opt defined")
    }
    if upperCase == true {
        fmt.Println("UpperCase opt defined")
    }

}
