package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        "bufio"
        "log"
        "os"
        "strings"
    )

    func check(e error) {
        if e != nil {
            panic(e)
        }
    }

    type Payee struct {
        Wallet string
        Share     int
        Pay     int64
    }

func parse() {
    // Open file and create scanner on top of it
    file, err := os.Open("test.txt")
    if err != nil {
        log.Fatal(err)
    }
    scanner := bufio.NewScanner(file)



    var int i:= 0

    for scanner.Scan() {
        fmt.Println("Line:", scanner.Text())

        temp := strings.Split(scanner.Text(), " ")

        payees[i].Wallet := temp[0]
        payees[i].Share := temp[1]
        payees[i].Pay := ((payees[i].Share / collateral) * customerpay)
        fmt.Println(payees[i].Wallet payees[i].Share payees[i].Pay)

        ++i

    }
}


    func main() {

    payees := make([]Payee)



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

    var int64 balance:= 37.5



    collateral := jsondata.collateral
    balance := balance()
    adminpercentage := jsondata.adminpercentage
    adminpay:= (balance * adminpercentage)
    customerpay:= balance - adminpay



    parse()





    }