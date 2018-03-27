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
        //"math"
        "net/http"
        "html/template"
        "github.com/spf13/viper"
    )

    var payments= []*Payee{}
    var result bytes.Buffer
    var payabort bool = false
    var err error
    var info= Config{}

    func check(e error) {
        if e != nil {
            panic(e)
        }
    }

    //struct for holding payee data
    type Payee struct {
        Wallet     string
        Share     float64
        Pay       float64
    }

    type Config struct {
        Customer        string
        Coin            string
        Coincli         string
        Payoutacct      string
        Collateral      float64
        Adminwallet     string
        Mnwallet        string
        Adminpercentage float64
        Payinfo         []*Payee
        Rbalance        float64
        Adminfee        float64
        Custpayout      float64
        Cmdstring       string
        Tx              string
        Date            string
    }

    //struct for holding getaddressbalance data
    type Bresult struct {
        Balance  float64
        Received float64
    }

    //chops everything after 2 decimal places
    func Truncate(some float64) float64 {
        return float64(int(some * 100)) / 100
    }

    //parses customerdata.dat for wallet and share amounts
    func custdata() {
        // Open file and create scanner on top of it
        file, err := os.Open("customerdata.dat")
        if err != nil {
            fmt.Println("error customerdata ", err.Error)
            payabort = true
        }
        scanner := bufio.NewScanner(file)
        //scans line by line to read customerdata and store data in payee
        for scanner.Scan() {
            temp := strings.Split(scanner.Text(), " ")
            payees := new(Payee)
            payees.Wallet= temp[0]
            tempshare, err:= strconv.ParseInt(temp[1], 10, 64)
                if err == nil {
                }
            payees.Share = float64(tempshare)
            payees.Pay = float64((payees.Share / info.Collateral) * info.Custpayout)
            payees.Pay = Truncate(payees.Pay)
            //payments is an array of payee
            payments = append(payments, payees)
        }
    }

    //parses payconfig.toml to set config for this coin
    func getconfig() {
        viper.SetConfigName("payconfig")
        viper.AddConfigPath(".")
        err := viper.ReadInConfig()
        if err != nil {
            fmt.Println("Config file not found...")
            payabort = true
        } else {
            info.Coin = viper.GetString("config.coin")
            info.Coincli = viper.GetString("config.cli")
            info.Payoutacct = viper.GetString("config.payoutacct")
            info.Collateral = viper.GetFloat64("config.collateral")
            info.Adminwallet = viper.GetString("config.adminwallet")
            info.Mnwallet = viper.GetString("config.mnwallet")
            info.Adminpercentage = viper.GetFloat64("config.adminpercentage")
            fmt.Printf("\n Config found:\n coin = %s\n", info.Coin)
            fmt.Printf(" coin-cli = %j\n", info.Coincli)
            fmt.Printf(" payoutacct = %i\n", info.Payoutacct)
            fmt.Printf(" collateral = %h\n", info.Collateral)
            fmt.Printf(" adminwallet = %g\n", info.Adminwallet)
            fmt.Println(" mnwallet = ", info.Mnwallet)
            fmt.Printf(" adminpercentage = %f\n", info.Adminpercentage)
        }
    }

    func getbalance() (float64) {
        fmt.Println("Getting Balance...")

        balancecmd := "getbalance"
        cmd := exec.Command(info.Coincli, balancecmd, info.Payoutacct)
        out, err := cmd.CombinedOutput()
        if err != nil {
            fmt.Println("exec error ", err.Error, out)
            payabort = true
        }

        n := strings.TrimSuffix(string(out[:]), "\n")
        s, err := strconv.ParseFloat(n, 64)
        //balance has no decimal so this puts it in the right place, this may need to be adjusted per coin project in which case I'll make it a variable for the payconfig
        var tbalance float64 = s / 100000000

        fmt.Println("Curent RAW Balance: ", tbalance)
        result.WriteString("Curent RAW Balance: ")
        result.WriteString(strconv.FormatFloat(tbalance, 'f', -1, 64))
        result.WriteString("\n")
        result.WriteString(n)
        tbalance= float64(tbalance - info.Collateral - 0.01)
        tbalance = Truncate(tbalance)
        //if balance is less than 20 for any reason don't pay out. prevents micro payments, also might need to be adjust per project
        if tbalance < 20 {
            fmt.Println("Balance too low: ", tbalance)
            result.WriteString("Balance too low: ")
            result.WriteString(strconv.FormatFloat(tbalance, 'f', -1, 64))
            result.WriteString("\n")
            payabort= true
        }
        return tbalance
    }

    func getaddressbalance() (float64) {
        fmt.Println("Getting Balance...")

        var balancecmd string = "getaddressbalance"
        t := []string{`{"addresses":["`, info.Mnwallet, `"]}`}
        var list string = strings.Join(t, "")
        cmd := exec.Command(info.Coincli, balancecmd, list)

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
        //balance has no decimal so this puts it in the right place, this may need to be adjusted per coin project in which case I'll make it a variable for the payconfig
        var tbalance float64 = bresults.Balance / 100000000
        fmt.Println("Curent RAW Balance: ", tbalance)
        result.WriteString("Curent RAW Balance: ")
        result.WriteString(strconv.FormatFloat(tbalance, 'f', -1, 64))
        result.WriteString("\n")
        result.WriteString(s)
        tbalance= float64(tbalance - info.Collateral - 0.01)
        tbalance = Truncate(tbalance)
        //if balance is less than 20 for any reason don't pay out. prevents micro payments, also might need to be adjust per project
        if tbalance < 20 {
            fmt.Println("Balance too low: ", tbalance)
            result.WriteString("Balance too low: ")
            result.WriteString(strconv.FormatFloat(tbalance, 'f', -1, 64))
            result.WriteString("\n")
            payabort= true
        }
        return tbalance
    }



    //assembles the command string
    func createcommand() (string) {
        p := []string{`{"`, info.Adminwallet, `":`, strconv.FormatFloat(info.Adminfee, 'f', -1, 64), `,"`}
        var p2 string = strings.Join(p, "")
        for k := range payments {
            p := []string{p2, string(payments[k].Wallet), `":`, strconv.FormatFloat(payments[k].Pay, 'f', -1, 64)}
            p2 = strings.Join(p, "")
      	    if (k+1) < len(payments) {
                p := []string{p2, `,"`}
      	        p2 = strings.Join(p, "")
      	    }
      	}
        p = []string{p2, `}`}
      	p2 = strings.Join(p, "")
    	return p2
    }

    func notification() {
        if payabort != true {
            fmt.Println("Sending Email")
            //api keys for email forwarder, will have to be variables and secured so it's not on github
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
            //need to learn how to create JSON with a list of emails so I can email customers reports
            //emaillist:= Recipient{"admin@kane.ventures"}
            emaillist:= Recipient{"kane4ventures@gmail.com"}
            s := ""
            buf := bytes.NewBufferString(s)
            t, _ := template.ParseFiles("email.html")
            t.Execute(buf, info)
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
        }
        fmt.Println("")
        fmt.Println("")
        fmt.Println("##########################")
        fmt.Println("Command Line Run Results")
        fmt.Println(result.String())
    }

    func main() {

        fmt.Println("Starting...")

        getconfig() 

        fmt.Println("")
        if info.Coin == "Shekel" {
            info.Rbalance = getbalance()
        } else {
            if info.Coin == "Advance" {
                info.Rbalance = getbalance()
            }
        } else {
            info.Rbalance = getaddressbalance()
        }
        fmt.Println("")

        //adminpay is % of balance, that it deducted from balance and the rest split among cust according to share
        info.Adminfee = float64(info.Rbalance * info.Adminpercentage)
        info.Adminfee = Truncate(info.Adminfee)
        info.Custpayout = float64(info.Rbalance - info.Adminfee)

        custdata()  

        info.Payinfo = payments

        info.Cmdstring = createcommand()

        var checkpayments float64
        for k := range payments {
            checkpayments= checkpayments + payments[k].Pay
            }
        if checkpayments > info.Custpayout {
            payabort = true
            fmt.Println("checkpayments > info.Custpayout    ", checkpayments, info.Custpayout)
        }
        if (checkpayments + info.Adminfee) > info.Rbalance {
            payabort = true
            fmt.Println("payments and adminpay higher than balance     ", checkpayments, info.Custpayout, info.Rbalance)
        }

        result.WriteString("Payout Report ")
        result.WriteString(time.Now().Format(time.RFC850))
        info.Date = time.Now().Format(time.RFC850)
        result.WriteString("\n")
        result.WriteString(info.Payoutacct)
	    result.WriteString(" ")
	    result.WriteString(strconv.FormatFloat(info.Rbalance, 'f', -1, 64))
	    result.WriteString("\n")
        result.WriteString("Admin Pay ")
	    result.WriteString(strconv.FormatFloat(info.Adminfee, 'f', -1, 64))
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
        result.WriteString(info.Coincli)
        result.WriteString(" sendmany ")
	result.WriteString(info.Cmdstring)
	result.WriteString("\n")

//////////////DEBUG MODE SWITCH set to true for testing comment out to get real
                               payabort = true
/////////////////////////

        if payabort != true {
            cmd := exec.Command(info.Coincli, `sendmany`, info.Payoutacct, info.Cmdstring)
            out, err := cmd.CombinedOutput()
            if err != nil {
                e := string(out[:])
                fmt.Println("exec error ", err.Error, e)
            }
            e := string(out[:])
            result.WriteString("Paycommand Output\n")
            result.WriteString(e)
            info.Tx = e
            result.WriteString("\n")
        }

        if (payabort == true) || (err != nil) {
            result.WriteString("Payout Aborted or failed")
            result.WriteString("\n")
            result.WriteString("payabort variable is: ")
            result.WriteString(strconv.FormatBool(payabort))
            result.WriteString("\n")
        }

        notification()
}
