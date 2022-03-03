package main

import (
    "fmt"
    "log"
    "io/ioutil"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/asaskevich/govalidator"
    "context"
    "github.com/aws/aws-lambda-go/lambda"
)


type MyEvent struct {
        Name string `json:"name"`
}

/*
func isValidEndpoint(target_ip string) bool{
    var rxIP = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}`)
    var rxDomain = regexp.MustCompile(`[a-zA-z]+[.][a-zA-Z]+[.][a-zA-Z]+|[a-zA-Z]+[.][a-zA-Z]+`)

    
    if len(target_ip) > 255 || !rxIP.MatchString(target_ip) || !rxDomain.MatchString(target_ip){
        return false
    }
    return true
}
*/

//for test purposes default page with no ip address input
func defaultPage(w http.ResponseWriter, r *http.Request){
    whoisString := fmt.Sprintf("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=at_zJSFEuLZW8F189sCNiK9kixtLkjVm&domainName=https://fishtech.group/")
    resp, err := http.Get(whoisString)
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println("hi", resp)
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalln(err)
    }
    //Convert the body to type string
    sb := string(body)
        
    //Just directly outputting the xml response from the WhoIS API back to the user.
    fmt.Fprintf(w, sb) 
    
}

func getIPInfo(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    
    /*
    govalidator is a functional version of what isValidEndpoint seeks to do
    validate whether or not the target ip address/dns hostname is in a valid format or not:
    0.0.0.0, www.____.___ or ____.___
    */
    
    if govalidator.IsIP(vars["target_ip"]) || govalidator.IsDNSName(vars["target_ip"]){
        //Utilizing the WhoIsXMLAPI found 
        //https://main.whoisxmlapi.com/ to provide a little bit of information on the desired ip/dns
        whoisString := fmt.Sprintf("https://www.whoisxmlapi.com/whoisserver/WhoisService?apiKey=at_zJSFEuLZW8F189sCNiK9kixtLkjVm&domainName=%s", vars["target_ip"])
        
        resp, err := http.Get(whoisString)
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Println("hi", resp)
        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            log.Fatalln(err)
        }
        //Convert the body to type string
        sb := string(body)
        
        //Just directly outputting the xml response from the WhoIS API back to the user.
        fmt.Fprintf(w, sb) 
    }
    
}

func handleRequests() {
    /*
    Package gorilla/mux implements a request router and dispatcher for matching incoming requests to their respective handler.
    https://github.com/gorilla/mux
    */
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/{target_ip}", getIPInfo)
    router.HandleFunc("/", defaultPage)
	log.Fatal(http.ListenAndServe(":8080", router))
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
        return fmt.Sprintf("Hello %s!", name.Name ), nil
}

func main() {
    handleRequests()
    //lambda.Start(HandleRequest)
}