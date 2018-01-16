package main


    import (
        "fmt"
        "io/ioutil"
        "os/exec"
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
    var payabort bool = false

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

    func createcommand(adminpay float64) {

        paycommand.WriteString(adminwallet)
        paycommand.WriteString("\\\":")
	    paycommand.WriteString(strconv.FormatFloat(adminpay, 'f', -1, 64))
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

    func notification() {
fmt.Println("Sending Email")

   var apikey="cb9872db45f62a1e4b67ded1736d85a1:b211992104c42942713d8c4cacad7ad2"


 var mailcommand string= fmt.Sprintf("-x /usr/bin/curl") //-s -X POST --user %s https://api.mailjet.com/v3/send -H 'Content-Type: application/json' -d '{ \"FromEmail\":\"paybot@kane.ventures\", \"FromName\":\"Mailjet Pilot\", \"Subject\":\"Test4\", \"Text-part\":\"Dear passenger, welcome to Mailjet! May the delivery force be with you!\", \"Html-part\":\"<h3>Dear passenger, welcome to Mailjet!</h3><br />May the delivery force be with you!\", \"Recipients\":[ { \"Email\": \"kanewinter@gmail.com\" } ] }'", apikey)

mailcmd := []string{mailcommand}

fmt.Println(apikey)
	//fmt.Println(mailcmd)

	   cmd := exec.Command("/usr/bin/curl", )
           	var out bytes.Buffer
           	cmd.Stdout = &out
           	err := cmd.Run()
		fmt.Println(cmd.Output)
		fmt.Println(cmd.Stdout)
           	if err != nil {
           		fmt.Println(err.Error)
           	}
           	result.WriteString(out.String())

           	var mailcommand string= fmt.Sprintf("-x /usr/bin/curl") //-s -X POST --user %s https://api.mailjet.com/v3/send -H 'Content-Type: application/json' -d '{ \"FromEmail\":\"paybot@kane.ventures\", \"FromName\":\"Mailjet Pilot\", \"Subject\":\"Test4\", \"Text-part\":\"Dear passenger, welcome to Mailjet! May the delivery force be with you!\", \"Html-part\":\"<h3>Dear passenger, welcome to Mailjet!</h3><br />May the delivery force be with you!\", \"Recipients\":[ { \"Email\": \"kanewinter@gmail.com\" } ] }'", apikey)
           	resp, err := http.Post("https://api.mailjet.com/v3/send", "application/json", &buf)


MJ_APIKEY_PUBLIC:= cb9872db45f62a1e4b67ded1736d85a1
MJ_APIKEY_PRIVATE:= b211992104c42942713d8c4cacad7ad2

type Payload struct {
	FromEmail  string `json:"FromEmail"`
	FromName   string `json:"FromName"`
	Subject    string `json:"Subject"`
	TextPart   string `json:"Text-part"`
	HTMLPart   string `json:"Html-part"`
	Recipients []struct {
		Email string `json:"Email"`
	} `json:"Recipients"`
}

data := Payload{
FromEmail: "paybot@kane.ventures",
FromName: "Paybot",
Subject: "test5",
TextPart: "lots ot text can go here",
HTMLPart: "lots of html",
Recipients: Recipients{
        Email: "kanewinter@gmail.com",
        },
    }
    }

c := &Configuration{
        Val: "test",
        Proxy: Proxy{
            Address: "addr",
            Port:    "port",
        },
    }

payloadBytes, err := json.Marshal(data)
if err != nil {
	// handle err
}
body := bytes.NewReader(payloadBytes)

req, err := http.NewRequest("POST", "https://api.mailjet.com/v3/send", body)
if err != nil {
	// handle err
}
req.SetBasicAuth(os.ExpandEnv("$MJ_APIKEY_PUBLIC"), os.ExpandEnv("$MJ_APIKEY_PRIVATE"))
req.Header.Set("Content-Type", "application/json")

resp, err := http.DefaultClient.Do(req)
if err != nil {
	// handle err
}
defer resp.Body.Close()



fmt.Println("mail sent?")

    }

    func init() {
        http.HandleFunc("/", handler)
    }
    func handler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello, world!")
    }

    func main() {


        datafile, err := ioutil.ReadFile("payconfig.dat")
        check(err)
        //fmt.Println(string(datafile))

        var jsondata interface{}
        json.Unmarshal(datafile, &jsondata)
        //fmt.Println(interface{}(jsondata))
        //fmt.Println()

fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")


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
        createcommand(adminpay)

fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")

        var checkpayments float64
        for k := range payments {
            checkpayments= checkpayments + payments[k].Pay
            }
        if checkpayments > customerpay {
            log.Fatal(checkpayments)
            payabort= true
        }

        if (checkpayments + adminpay) > balance {
                    log.Fatal(balance)
                    payabort= true
        }



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

fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")


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
	    result.WriteString("\n")

        fmt.Println(result.String())

        var paycmd string = paycommand.String()
	    fmt.Println(paycmd)

fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")

        if payabort != true {
            cmd := exec.Command("gobyte-cli", "paycmd")
        	var out bytes.Buffer
        	cmd.Stdout = &out
        	err := cmd.Run()
        	if err != nil {
        		//log.Fatal(err)
        	}
        	result.WriteString(out.String())
        }
fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")
        if (payabort == true) || (err != nil) {
            result.WriteString("Payout Aborted or failed")
            result.WriteString("payabort variable is: ")
            result.WriteString(strconv.FormatBool(payabort))
            result.WriteString("\n")
            result.WriteString("err variable is: ")
            result.WriteString(err.Error())
            result.WriteString("\n")
         }

fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")
fmt.Println("*************")

         notification()

}
