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
        "bytes"
    )

    var collateral float64
    var adminpercentage float64
    var adminpay float64
    var customerpay float64
    var balance float64
    var payments= []*Payee{}
    var paycommand bytes.Buffer
    var payoutacct string

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
        file, err := os.Open("customerdata.dat")
        if err != nil {
            log.Fatal(err)
        }
        scanner := bufio.NewScanner(file)

        for scanner.Scan() {
            //fmt.Println("Line:", scanner.Text())

            temp := strings.Split(scanner.Text(), " ")

            payees := new(Payee)
            payees.Wallet= temp[0]

            tempshare, err:= strconv.ParseInt(temp[1], 10, 64)
                if err == nil {
                    //fmt.Println(tempshare)
                }

            payees.Share= float64(tempshare)
            payees.Pay= float64((payees.Share / collateral) * customerpay)
            fmt.Println(payees.Wallet, payees.Share, payees.Pay)
            payments = append(payments, payees)

        }
	//for k := range payments {
	//fmt.Println(payments[k].Wallet, payments[k].Share, payments[k].Pay)
	//}
    }

    func createcommand() {
        for k := range payments {
      	    //fmt.Println(payments[k].Wallet, payments[k].Pay)
      	    tempwallet:= string(payments[k].Wallet)
      	    paycommand.WriteString(tempwallet)
      	    paycommand.WriteString("\":")
      	    temppay := strconv.FormatFloat(payments[k].Pay, 'f', -1, 64)
      	    paycommand.WriteString(temppay)

      	    if len(payments) == k {
      	    paycommand.WriteString(",\"")
      	    }
      	}
      	paycommand.WriteString("}\"")
        fmt.Println(paycommand.String())
    }



    func main() {


        datafile, err := ioutil.ReadFile("payconfig.dat")
        check(err)
        //fmt.Println(string(datafile))

        var jsondata interface{}
        json.Unmarshal(datafile, &jsondata)
        //fmt.Println(interface{}(jsondata))
        //fmt.Println()


        var balance= 37.5
        var payoutacct= "BP&C Payout" //jsondata.payoutacct
        paycommand.WriteString("sendmany ")
        paycommand.WriteString(payoutacct)
        paycommand.WriteString("" "{\"")


        collateral= 1000 //jsondata.collateral
        // balance= balance()
        //adminpercentage= jsondata.adminpercentage
        var adminpercentage= 0.1
        var adminpay float64 = float64(balance * adminpercentage)
        customerpay = float64(balance - adminpay)

        parse()
        createcommand()


    }
