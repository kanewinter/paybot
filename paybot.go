package main


    import (
        "fmt"
        //"io/ioutil"
        "os/exec"
        "encoding/json"
        "bufio"
        //"log"
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
    var mnwallet string
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

    type Bresult struct {
        Balance  float64
        Received float64
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

    func Truncate(some float64) float64 {
        return float64(int(some * 100)) / 100
    }


    func custdata() {
        // Open file and create scanner on top of it
        file, err := os.Open("customerdata.dat")
        if err != nil {
		    fmt.Println("error customerdata ", err.Error)
		    payabort = true
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
            payees.Pay = Truncate(payees.Pay)
            //payees.Pay = Round(payees.Pay, .5, 2)
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
            mnwallet = viper.GetString("config.mnwallet")
            adminpercentage = viper.GetFloat64("config.adminpercentage")

            fmt.Printf("\n Config found:\n coin = %s\n", coin)
            fmt.Printf(" coin-cli = %j\n", coincli)
            fmt.Printf(" payoutacct = %i\n", payoutacct)
            fmt.Printf(" collateral = %h\n", collateral)
            fmt.Printf(" adminwallet = %g\n", adminwallet)
            fmt.Println(" mnwallet", mnwallet)
            fmt.Printf(" adminpercentage = %f\n", adminpercentage)
                }
    }

    func getbalance() (float64) {
        fmt.Println("Getting Balance...")
        var balancecmd string = "getaddressbalance"
        t := []string{`{"addresses":["`, mnwallet, `"]}`}
        var list string = strings.Join(t, "")

        //////WORKS!!
        ////var list string = `{"addresses":["AZDgBUM6kcSTyqxH2Q4ig3G54xjpvYcynE"]}`
        //////////////////

        cmd := exec.Command(coincli, balancecmd, list)
        out, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Println("exec error ", err.Error, out)
        	payabort = true
        }

        s := string(out[:])
        var bresults Bresult
        outbyte := []byte(s)

        err = json.Unmarshal(outbyte, &bresults)
        if err != nil {
               fmt.Println("error:", err)
        }

        var tbalance float64 = bresults.Balance / 100000000
        fmt.Println("Curent RAW Balance: ", tbalance)
        result.WriteString("Curent RAW Balance: ")
	result.WriteString(strconv.FormatFloat(tbalance, 'f', -1, 64))
        result.WriteString("\n")
        result.WriteString(s)
        tbalance= float64(tbalance - collateral - 2)
        tbalance = Truncate(tbalance)

        if tbalance < 20 {
            fmt.Println("Balance too low: ", tbalance)
            result.WriteString("Balance too low: ")
            result.WriteString(strconv.FormatFloat(tbalance, 'f', -1, 64))
            result.WriteString("\n")
            payabort= true
        }

        return tbalance
    }

    func createcommand() (string) {
        p := []string{`{"`, adminwallet, `":`, strconv.FormatFloat(adminpay, 'f', -1, 64), `,"`}
        var p2 string = strings.Join(p, "")
        fmt.Println(p2)

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
            p := []string{p2, tempwallet, `":`, strconv.FormatFloat(payments[k].Pay, 'f', -1, 64)}
            p2 = strings.Join(p, "")
            fmt.Println(p2)


      	    if (k+1) < len(payments) {
      	    paycommand.WriteString(",\\\"")
	    p := []string{p2, `,"`}
      	    p2 = strings.Join(p, "")
      	    }
      	}
      	paycommand.WriteString("}\"")
	p = []string{p2, `}`}
      	p2 = strings.Join(p, "")
	fmt.Println(p2)
	return p2
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

        fmt.Println("")
        fmt.Println("")
        fmt.Println("##########################")
        fmt.Println("Command Line Run Results")
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

        paycmd := createcommand()

        var checkpayments float64
        for k := range payments {
            checkpayments= checkpayments + payments[k].Pay
            }
        if checkpayments > customerpay {
            payabort = true
	        fmt.Println(checkpayments, customerpay)
            payabort= true
        }

        if (checkpayments + adminpay) > balance {
            payabort = true
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
	    result.WriteString(paycmd)
	    result.WriteString("\n")

        //var paycmd string = paycommand.String()
	

//////////////DEBUG MODE SWITCH set to true for testing comment out to get real
//        payabort = true
/////////////////////////


        //fmt.Println(paycmd)

        if payabort != true {
	fmt.Println("######PAY COMMAND")
	fmt.Println(paycmd)
	fmt.Println("/n")
cmd := exec.Command(coincli, `sendmany`, payoutacct, paycmd)
       //     cmd := exec.Command(coincli, `sendmany`, paycmd)
        	out, err := cmd.CombinedOutput()
        	if err != nil {
        	    e := string(out[:])
                fmt.Println("exec error ", err.Error, e)
        	}
	fmt.Println(cmd)
	e := string(out[:])
	fmt.Println(e)

        //	tmp := strings.TrimSuffix(out.String(), "\n")
            //fmt.Println(cmd.Stdout)
            //fmt.Print(string(out.Bytes()))
        	//fmt.Println(out.String())
        	//result.WriteString(out.String())
        	//fmt.Println(tmp)
            //result.WriteString(tmp)

        }

        if (payabort == true) || (err != nil) {
            result.WriteString("Payout Aborted or failed")
            result.WriteString("\n")
            result.WriteString("payabort variable is: ")
            result.WriteString(strconv.FormatBool(payabort))
            result.WriteString("\n")
            result.WriteString("err variable is: ")
            //result.WriteString(err.Error())
            result.WriteString("\n")
         }

         notification()


}
