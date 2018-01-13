package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        "bufio"
        "log"
        "os"
        "strings"
        "strconv"
    )

    var collateral float64
    var adminpercentage float64
    var adminpay float64
    var customerpay float64
    var balance float64
    var payments= []*Payee{}

    func check(e error) {
        if e != nil {
            panic(e)
        }
    }

    type Payee struct {
        Wallet     string
        Share     float64
        Pay       float64
    }

    func parse() {
        // Open file and create scanner on top of it
        file, err := os.Open("test.txt")
        if err != nil {
            log.Fatal(err)
        }
        scanner := bufio.NewScanner(file)
        i:= 0

        for scanner.Scan() {
            fmt.Println("Line:", scanner.Text())

            temp := strings.Split(scanner.Text(), " ")

            payees := new(Payee)
            payees.Wallet= temp[0]

            tempshare, err:= strconv.ParseInt(temp[1], 10, 64)
                if err == nil {
                    fmt.Println(tempshare)
                }

            payees.Share= float64(tempshare)
            payees.Pay= float64((payees.Share / collateral) * customerpay)
            fmt.Println(payees.Wallet, payees.Share, payees.Pay)
            payments = append(payments, payees)

            i= i+1

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

        var balance= 37.5



        collateral= 1000 //jsondata.collateral
        // balance= balance()
        //adminpercentage= jsondata.adminpercentage
        var adminpercentage= 0.1
        var adminpay float64 = float64(balance * adminpercentage)
        customerpay = float64(balance - adminpay)

        parse()


    }