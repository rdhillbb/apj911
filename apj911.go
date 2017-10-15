package apj911

//package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/rdhillbb/alertconfig"
	"github.com/rdhillbb/buildWatsonNumber"
	"github.com/rdhillbb/createTxT2Spch"
	"github.com/rdhillbb/slackpostmsg"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TwiML struct {
	XMLName xml.Name `xml:"Response"`

	Say  string `xml:",omitempty"`
	Play string `xml:",omitempty"`
}

//Global Variable Definitions

var ServiceConfig alertconfig.ServiceConfig
var AuthTokens alertconfig.Authtokens
var CallList []*alertconfig.CallMember

// 3 Team Members


var listSize int = 0

// Token Information
var call_count int


func C_init() {
	//Get Authentication and Account tokens
	fmt.Println("\ninit:0000 Intialization started at ", time.Now().Format(time.RFC850))
	start := time.Now()
	alertconfig.BuildCallList(&CallList, &ServiceConfig, &AuthTokens)
	fmt.Println("init:000 Initialization Elapsed time: ", time.Since(start))
	
}

func GetIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-FORWARDED-FOR"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func Fileh(w http.ResponseWriter, r *http.Request) {

	fmt.Println("\nfileh:0001:---------------------------Start---------------\n")
	start := time.Now()
	if r.Method == "GET" {
		fmt.Println("GET NOP -------------->")
	} else if r.Method == "POST" {
		fmt.Println("POST NOP  -------------->")
	}

	fmt.Println("fileh:0001:http Request", r)
	fmt.Println("fileh:0001:IP Address:", GetIP(r))
	fmt.Println("fileh:0001:Video File", ServiceConfig.Video_dir+"/"+ServiceConfig.Audio_file)
	fmt.Println("*")
	http.ServeFile(w, r, ServiceConfig.Video_dir+"/"+ServiceConfig.Audio_file)
	fmt.Println("fileh:0001:---------------------------End--------------- Elasped time:", time.Since(start))
}

func direct_sms(to string, from string, msg string) (str string) {
	// Set initial variables
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + AuthTokens.TwaccountSid + "/Messages.json"

	// Build out the data for our message
	v := url.Values{}
	v.Set("To", to)
	v.Set("From", from)
	v.Set("Body", msg)
	rb := *strings.NewReader(v.Encode())
	// Create client
	client := &http.Client{}

	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(AuthTokens.TwaccountSid, AuthTokens.TwauthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	fmt.Println(req)

	// Make request
	resp, _ := client.Do(req)
	str = "SMS Status " + to + resp.Status
	return
}

// Dial the PHone and Play the message
func direct_call(to string, from string) (str string) {
	start := time.Now()
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + AuthTokens.TwaccountSid + "/Calls.json"
	str = ""
	// Build out the data for our message
	v := url.Values{}
	v.Set("To", to)
	v.Set("From", from)
	v.Set("Url", ServiceConfig.TW_site)
	fmt.Println("\ndirect:0001:----------------------Start---------------------")
	fmt.Println("direct_call:0000 Start --> ", time.Now().Format(time.RFC850))
	rb := *strings.NewReader(v.Encode())
	// Create Client
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &rb)
	req.SetBasicAuth(AuthTokens.TwaccountSid, AuthTokens.TwauthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	fmt.Println("---req", urlStr)
	// make request
	resp, _ := client.Do(req)
	str = "direct:0001:----***** " + resp.Status
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Println("resp.StatusCode >= 200 && resp.StatusCode --", resp.StatusCode)
		var data map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		err := json.Unmarshal(bodyBytes, &data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println("call_direct:0001: good rntcode", resp.Status)
	}
	fmt.Println("direct:0001:----------------------End--------------------- Elasped time: ", time.Since(start))
	return
}

func Call(w http.ResponseWriter, r *http.Request, msg string) {
	// Let's set some initial default variables
	rntMsg := "Broadcast Message Sent to: "
	if call_count > 30 {

		w.Write([]byte("Recharge Twilio Tokeno"))
		return
	}

	for i := range CallList {
		place := CallList[i]
		start := time.Now()

		fmt.Println("to: ", place.Phone, " From", ServiceConfig.From_Numbertw)
		rntMsg += place.Name + ":" + place.Phone + " : "
		//* Mobile
		if place.Phone != "" {
			fmt.Println("Exit Call direct_call", direct_call(place.Phone, ServiceConfig.From_Numbertw))
			fmt.Println("Exit Call ", direct_sms(place.Phone, ServiceConfig.From_Numbertw, msg))
		}
		if place.Slackid != "" {
			slackpostmsg.Slkpostmsg(AuthTokens.SlackToken, place.Slackid, msg)
		}
		fmt.Println("Elapsed time:", time.Since(start))
	}
	//w.Write([]byte(rntMsg))
	return
}

func broadcastMSG(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		fmt.Println("GET-------------->")
		t, _ := template.ParseFiles("broad2field.gtpl")
		t.Execute(w, nil)
	} else if r.Method == "POST" {

		r.ParseForm()
		// logic part of log in
		fmt.Println("Company Name:", r.Form["companyName"])
		fmt.Println("Case Number:", r.Form["caseNumber"])
		fmt.Println("Case Priority:", r.Form["casePriority"])
		companyName := r.Form["companyName"][0]
		caseNumber := r.Form["caseNumber"][0]
		message := buildWatsonNumber.CrMsg4wat(companyName, caseNumber)
		start := time.Now()
		createTxT2Spch.CreateWatAudio(ServiceConfig.Video_dir, ServiceConfig.Audio_file, message, AuthTokens.WatsonToken, AuthTokens.WatsonPass)
		fmt.Println("Broadcast:000 Wattson: ", time.Since(start))
		Call(w, r, "Alert P1 "+companyName+" Case Number: "+caseNumber+" Please contact the TAC Duty Manager.")
		d, _ := template.ParseFiles("broad2field.gtpl")
		d.Execute(w, nil)
	}
}

func Twiml(w http.ResponseWriter, r *http.Request) {
	//twiml := TwiML{Say: "Test Alert you have a P1 case! Alert you have a P1 case!"}
	//twiml := TwiML{Play: "http://45.55.14.147:3000/audio"}  ServiceConfig.TW_site
	start := time.Now()
	//twiml := TwiML{Play: ServiceConfig.TW_site + "/audio"}
	twiml := TwiML{Play: "/audio"}
	fmt.Println("\ntwilm:0001:--------------------------start -----------------")
	x, err := xml.Marshal(twiml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
	fmt.Printf("twilm:0001:--------------------------End-----------------Elasped Time: ", time.Since(start), "\n")
}

func main() {
	    C_init()
	fmt.Println("Broadcast Service Starting")
	http.HandleFunc("/audio", Fileh)
	http.HandleFunc("/notifiyField", broadcastMSG)
	http.HandleFunc("/twiml", Twiml)
	//http.HandleFunc("/call", call)
	http.HandleFunc("/twiml/audio", Fileh)
	http.ListenAndServe(":3000", nil)

}
