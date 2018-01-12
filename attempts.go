package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        //"os"
    )

    func check(e error) {
        if e != nil {
            panic(e)
        }
    }

    func main() {

    datafile, err := ioutil.ReadFile("payconfig.dat")
    check(err)
    fmt.Println(string(datafile))
    fmt.Println()

    var jsondata interface{}
    json.Unmarshal(datafile, &jsondata)
    fmt.Println(interface{}(jsondata))
    fmt.Println()

    textfile, err := ioutil.ReadFile("payconfig.dat")
        if err != nil {
            fmt.Print(err)
        }

        fmt.Println(b) // print the content as 'bytes'
        str := string(b) // convert content to a 'string'
        fmt.Println(str)

    fmt.Println()
    //payee := m.[]interface{}
    //fmt.Println(payee\n)

    }