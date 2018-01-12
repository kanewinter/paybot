package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        //"os"
    )

    func Unmarshal(data []byte, v interface{}) error

    func check(e error) {
        if e != nil {
            panic(e)
        }
    }

    func main() {

    datafile, err := ioutil.ReadFile("payconfig.dat")
    check(err)
    fmt.Print(string(datafile))

    var jsondata interface{}
    if err := json.Unmarshal(datafile, &jsondata) {
            log.Fatal(err) }
    fmt.Print(interface{}(jsondata))


    }