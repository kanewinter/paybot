package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        "os"
    )





    datafile, err := ioutil.ReadFile("payconfig.dat")
    check(err)
    fmt.Print(string(datafile))