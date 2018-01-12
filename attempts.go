package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        "bufio"
        "log"
        "os"
    )

    func check(e error) {
        if e != nil {
            panic(e)
        }
    }

func parse() {
    // Open file and create scanner on top of it
    file, err := os.Open("test.txt")
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)

    for scanner.Scan {
        fmt.Println("First line found:", scanner.Text())

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

    fmt.Println("print bytes")
    fmt.Println(textfile) // print the content as 'bytes'
    fmt.Println()

        fmt.Println("string")
    str := string(textfile) // convert content to a 'string'
    fmt.Println(str)

    fmt.Println()

    parse()





    }