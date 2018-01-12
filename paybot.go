package main


    import (
        "fmt"
        "io/ioutil"
        "encoding/json"
        "os"
    )

    func Unmarshal(data []byte, v interface{}) error

    func check(e error) {
        if e != nil {
            panic(e)
        }
    }

    func balance() {
    	if out, err := exec.Command(/opt/gobyte/gobyte-cli getbalance).Output(); err == nil {
    		log.Info("Balance is " out.Output)
    	}
    }

    func paycommand() {
        for i := range payoutdatawallet

        }


    datafile, err := ioutil.ReadFile("payconfig.dat")
    check(err)
    fmt.Print(string(datafile))

    var jsondata interface{}
    err := json.Unmarshal(datafile, &jsondata)

    m := jsondata.(map[string]interface{})
    payee := m.[]interface{}
    fmt.Println(payee)

    var payoutdatawallet []string
    var payoutdatashare []int
    var i int
    i := 0

    for k, v := range payee {

        payoutdatashare[i] := v
        payoutdatawallet[i] : = k
        ++i

    }



    var data map[string]interface{}

    if err := json.Unmarshal(datafil, &data); err != nil {
        panic(err)
    }
    fmt.Println(data)

    collateral := data["collateral"].(float64)
    fmt.Println(collateral)

    customerdata := data["customerdata"].([]interface{})
    array payoutwallets
    array payoutamounts
    for each customerdata
    payoutwallets[i] := customerdata[i].(string)
    payoutamounts[i] := customerdata[i].




    fmt.Println(str1)


for k, v := range m {
    switch vv := v.(type) {
    case string:
        fmt.Println(k, "is string", vv)
    case float64:
        fmt.Println(k, "is float64", vv)
    case []interface{}:
        fmt.Println(k, "is an array:")
        for i, u := range vv {
            fmt.Println(i, u)
        }
    default:
        fmt.Println(k, "is of a type I don't know how to handle")
    }
}



total=`/opt/gobyte/gobyte-cli getbalance`

admin=$( bc -l <<<"0.10*$total" )

customer=$( bc -l <<<"$total-$admin" )

echo "/opt/gobyte/gobyte-cli sendmany "kv" "{\"GfvkGVFKSCPRpdfTXRgATxvyxCATzpj3LE\":$admin,\"GfeJJFCU7qULYNXk3pJ3d56pfazYjKCDS8\":$customer}""