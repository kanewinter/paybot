package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        "os"
    )

    func Unmarshal(data []byte, v interface{}) error




    datafile, err := ioutil.ReadFile("payconfig.dat")
    check(err)
    fmt.Print(string(datafile))