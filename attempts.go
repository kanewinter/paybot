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
    fmt.Print(string(datafile))
    fmt.Println()

    var jsondata interface{}
    json.Unmarshal(datafile, &jsondata)
    fmt.Print(interface{}(jsondata))
    fmt.Println()

    m := jsondata.customerdata
    fmt.Println(m)

    fmt.Println()
    //payee := m.[]interface{}
    //fmt.Println(payee\n)

    }