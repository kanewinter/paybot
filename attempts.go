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

    var jsondata interface{}
    json.Unmarshal(datafile, &jsondata)
    fmt.Print(interface{}(jsondata))

    m := jsondata.(map[string]interface{})
    fmt.Println(m)
    //payee := m.[]interface{}
    //fmt.Println(payee)

    }