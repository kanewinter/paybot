package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        "os"
    )

    func Unmarshal(data []byte, v interface{}) error

    func check(e error) {
        if e != nil {
            panic(e)
        }
    }

    func main() {

    var datafile, err := ioutil.ReadFile("payconfig.dat")
    check(err)
    fmt.Print(string(datafile))

    }