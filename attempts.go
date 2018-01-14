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
	"time"
    )

    var collateral float64
    var adminpercentage float64
    var adminpay float64
    var customerpay float64
    var balance float64
    var payments= []*Payee{}
    var paycommand bytes.Buffer
    var result bytes.Buffer
    var payoutacct string
    var adminwallet string

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
            //fmt.Println(payees.Wallet, payees.Share, payees.Pay)
            payments = append(payments, payees)

        }
	//for k := range payments {
	//fmt.Println(payments[k].Wallet, payments[k].Share, payments[k].Pay)
	//}
    }

    func createcommand() {

        paycommand.WriteString(adminwallet)
        paycommand.WriteString("\\\":")
        tempadminpay := strconv.FormatFloat(payments[k].Pay, 'f', -1, 64)
        paycommand.WriteString(tempadminpay)
        paycommand.WriteString(",\\\"")

        for k := range payments {
      	    tempwallet:= string(payments[k].Wallet)
      	    paycommand.WriteString(tempwallet)
      	    paycommand.WriteString("\\\":")
      	    temppay := strconv.FormatFloat(payments[k].Pay, 'f', -1, 64)
      	    paycommand.WriteString(temppay)

      	    if (k+1) < len(payments) {
      	    paycommand.WriteString(",\\\"")
      	    }
      	}
      	paycommand.WriteString("}\"")

    }



    func main() {


        datafile, err := ioutil.ReadFile("payconfig.dat")
        check(err)
        //fmt.Println(string(datafile))

        var jsondata interface{}
        json.Unmarshal(datafile, &jsondata)
        //fmt.Println(interface{}(jsondata))
        //fmt.Println()


        var balance float64 = 122.5
        var payoutacct= "BP&C Payout" //jsondata.payoutacct
        paycommand.WriteString("sendmany ")

	    fmt.Fprintf(&paycommand, "\"")
	    paycommand.WriteString(payoutacct)
	    fmt.Fprintf(&paycommand, "\" \"{\\\"")

        collateral= 1000 //jsondata.collateral
        // balance= balance()
        //adminpercentage= jsondata.adminpercentage
        //adminwallet= jsondata.adminwallet
        adminwallet= "dfhsdfgdfgnfjdsgdfsgkmlsmkgrimnn"
        var adminpercentage= 0.1
        var adminpay float64 = float64(balance * adminpercentage)
        customerpay = float64(balance - adminpay)

        parse()
        createcommand()

        result.WriteString("Payout Report ")
	    result.WriteString(time.Now().Format(time.RFC850))
	    result.WriteString("\n")
        result.WriteString(payoutacct)
	    result.WriteString(" ")
	    result.WriteString(strconv.FormatFloat(balance, 'f', -1, 64))
	    result.WriteString("\n")
        result.WriteString("Admin Pay ")
	    result.WriteString(strconv.FormatFloat(adminpay, 'f', -1, 64))
        result.WriteString("\n")
        result.WriteString("Wallets                             Share    Payout\n")

        for k := range payments {
		result.WriteString(payments[k].Wallet)
		result.WriteString("    ")
		result.WriteString(strconv.FormatFloat(payments[k].Share, 'f', -1, 64))
		result.WriteString("      ")
		result.WriteString(strconv.FormatFloat(payments[k].Pay, 'f', -1, 64))
		result.WriteString("\n")
        }

	    result.WriteString("\n")
        result.WriteString("Pay Command to be Used \n")
	    result.WriteString(paycommand.String())

        fmt.Println(result.String())
 
    }
