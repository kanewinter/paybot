package main


    import (
        "fmt"
        //"io/ioutil"
        "os/exec"
        "encoding/json"
        "bufio"
        "log"
        "os"
        "strings"
        "strconv"
        "bytes"
	    "time"
	    "math"
	    "net/http"
	    "html/template"
	    "github.com/spf13/viper"
    )

    var adminpay float64
    var customerpay float64
    var balance float64
    var payments= []*Payee{}
    var paycommand bytes.Buffer
    var result bytes.Buffer
    var payabort bool = false
    var adminwallet string
    var coin string
    var coincli string
    var payoutacct string
    var collateral float64
    var adminpercentage float64
    var err error

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

    func Round(val float64, roundOn float64, places int ) (newVal float64) {
	    var round float64
	    pow := math.Pow(10, float64(places))
	    digit := pow * val
	    _, div := math.Modf(digit)
	    if div >= roundOn {
	    	round = math.Ceil(digit)
	    } else {
	    	round = math.Floor(digit)
	    }
	    newVal = round / pow
	    return
    }

    func custdata() {
        // Open file and create scanner on top of it
        file, err := os.Open("customerdata.dat")
        if err != nil {
		    fmt.Println("error customerdata ", err.Error)
		    log.Fatal(err)
        }
        scanner := bufio.NewScanner(file)

        for scanner.Scan() {
            temp := strings.Split(scanner.Text(), " ")

            payees := new(Payee)
            payees.Wallet= temp[0]

            tempshare, err:= strconv.ParseInt(temp[1], 10, 64)
                if err == nil {
                }

            payees.Share = float64(tempshare)
            payees.Pay = float64((payees.Share / collateral) * customerpay)
            payees.Pay = Round(payees.Pay, .5, 2)
            payments = append(payments, payees)

        }
    }

    func getconfig() {
        viper.SetConfigName("payconfig")
        viper.AddConfigPath(".")
        err := viper.ReadInConfig()
        if err != nil {
            fmt.Println("Config file not found...")
            payabort = true
          } else {
            coin = viper.GetString("config.coin")
            coincli = viper.GetString("config.cli")
            payoutacct = viper.GetString("config.payoutacct")
            collateral = viper.GetFloat64("config.collateral")
            adminwallet = viper.GetString("config.adminwallet")
            adminpercentage = viper.GetFloat64("config.adminpercentage")

            fmt.Printf("\n Config found:\n coin = %s\n", coin)
            fmt.Printf(" coin-cli = %j\n", coincli)
            fmt.Printf(" payoutacct = %i\n", payoutacct)
            fmt.Printf(" collateral = %h\n", collateral)
            fmt.Printf(" adminwallet = %g\n", adminwallet)
            fmt.Printf(" adminpercentage = %f\n", adminpercentage)
                }
    }

    func getbalance() (float64) {
        fmt.Println("Getting Balance...")
        balancecmd := "getbalance"
        cmd := exec.Command(coincli, balancecmd)
        var out bytes.Buffer
        cmd.Stdout = &out
        err := cmd.Run()
        if err != nil {
		    fmt.Println("exec error ", err.Error, out.String())
        	log.Fatal(err)
        }
       	//result.WriteString(out.String())
        fmt.Println(out.String())
        tmp := strings.TrimSuffix(out.String(), "\n")
        things, err := strconv.ParseFloat(tmp, 64)
                if err != nil {
        		    fmt.Println("exec error ", err.Error, tmp)
                	log.Fatal(err)
                }
        things = float64(things - collateral - 1)
        things = Round(things, .5, 1)

        if things < 20 {
        fmt.Println("Balance too low: ", things)
        result.WriteString("Balance too low: ")
        result.WriteString(strconv.FormatFloat(things, 'f', -1, 64))
        result.WriteString("\n")
        payabort= true
        }

        return things
    }


    func createcommand() {

        paycommand.WriteString("sendmany ")
        fmt.Fprintf(&paycommand, "\"")
        paycommand.WriteString(payoutacct)
        fmt.Fprintf(&paycommand, "\" \"{\\\"")
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

        var MJ_APIKEY_PUBLIC string= "cb9872db45f62a1e4b67ded1736d85a1"
        var MJ_APIKEY_PRIVATE string= "b211992104c42942713d8c4cacad7ad2"

        type Recipient struct {
            Email string `json:"Email"`
        }

        type Payload struct {
        	FromEmail  string `json:"FromEmail"`
	        FromName   string `json:"FromName"`
	        Subject    string `json:"Subject"`
	        TextPart   string `json:"Text-part"`
	        HTMLPart   string `json:"Html-part"`
	        Recipients []Recipient `json:"Recipients"`
        }

        //emaillist:= Recipient{"admin@kane.ventures"}
        emaillist:= Recipient{"kane4ventures@gmail.com"}

        s := ""
        buf := bytes.NewBufferString(s)

        t, _ := template.ParseFiles("email.html")
        t.Execute(buf, payments)

        data := Payload{
            FromEmail: "paybot@kane.ventures",
            FromName: "Paybot",
            Subject: "Payout Report",
            TextPart: result.String(),
            HTMLPart: buf.String(),
            //HTMLPart: result.String(),
            Recipients: []Recipient {emaillist},
        }

        payloadBytes, err := json.Marshal(data)
        if err != nil {
	        // handle err
	        fmt.Println("ERROR")
        }
        body := bytes.NewReader(payloadBytes)

        req, err := http.NewRequest("POST", "https://api.mailjet.com/v3/send", body)
        if err != nil {
	        // handle err
	        fmt.Println("ERROR")
        }
        req.SetBasicAuth(os.ExpandEnv(MJ_APIKEY_PUBLIC), os.ExpandEnv(MJ_APIKEY_PRIVATE))
        req.Header.Set("Content-Type", "application/json")

        resp, err := http.DefaultClient.Do(req)
        if err != nil {
	        // handle err
	        fmt.Println("ERROR")
        }
        defer resp.Body.Close()

        fmt.Println(result.String())

    }

    func main() {

        getconfig()
        fmt.Println("")
        balance = getbalance()

        adminpay = float64(balance * adminpercentage)
        customerpay = float64(balance - adminpay)
        adminpay = Round(adminpay, .5, 2)

        custdata()

        createcommand()

        var checkpayments float64
        for k := range payments {
            checkpayments= checkpayments + payments[k].Pay
            }
        if checkpayments > customerpay {
            log.Fatal(checkpayments)
	        fmt.Println(checkpayments, customerpay)
            payabort= true
        }

        if (checkpayments + adminpay) > balance {
            log.Fatal(balance)
		    fmt.Println(checkpayments, customerpay, balance)
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
        result.WriteString(coincli)
        result.WriteString(" ")
	    result.WriteString(paycommand.String())
	    result.WriteString("\n")

        var paycmd string = paycommand.String()

//////////////DEBUG MODE SWITCH set to true for testing comment out to get real
        payabort = true
/////////////////////////

        if payabort != true {
            cmd := exec.Command(coincli, paycmd)
        	var out bytes.Buffer
        	cmd.Stdout = &out
        	err := cmd.Run()
        	if err != nil {
			fmt.Println("exec error ", err.Error, out.String())
        	}
        	result.WriteString(out.String())
        }

        if (payabort == true) || (err != nil) {
            result.WriteString("Payout Aborted or failed")
            result.WriteString("payabort variable is: ")
            result.WriteString(strconv.FormatBool(payabort))
            result.WriteString("\n")
            result.WriteString("err variable is: ")
            //result.WriteString(err.Error())
            result.WriteString("\n")
         }

         notification()

}
