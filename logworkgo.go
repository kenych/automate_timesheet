package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func confirm(date string) bool {

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("add ", date, "?(y/n):")
		text, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(text)) == "y" {
			return true
		} else if strings.TrimSpace(strings.ToLower(text)) == "n" {
			return false
		}
	}

}

func validate(date string) {
	b, err := ioutil.ReadFile("data-1.txt")

	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else if strings.Contains(string(b), date) {
		// fmt.Println("already done with ", date)
		os.Exit(0)
	}

}

func template(date string) {
	b, err := ioutil.ReadFile("tpl-1.txt")
	if err != nil {
		panic(err)
	}

	r, err := regexp.Compile("startDate=")
	ioutil.WriteFile("data-1.txt", r.ReplaceAll(b, []byte("startDate="+date+"%2F18+01%3A26+PM")), 0666)
}

func sendRequest(username string, password string, url string) int {
	fmt.Println("URL:>", url)

	b, err := ioutil.ReadFile("data-1.txt")
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.SetBasicAuth(username, password)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)

	return resp.StatusCode
}

func main() {

	if len(os.Args) < 6 {
		fmt.Println("Missing args, run as \n1) 'logworkgo.go path_to_data $username $password $jira_add $jira_open 21 Jan' - for specific date or\n2) 'logworkgo.go /Users/kayanazimov/Downloads/ $username $password $jira_add $jira_open' -  for current date")
		os.Exit(1)
	}

	workingPath := os.Args[1]
	userName := os.Args[2]
	password := os.Args[3]
	jira_add := os.Args[4]
	jira_open := os.Args[5]

	var date string

	os.Chdir(workingPath)

	if len(os.Args) == 8 {
		day := os.Args[6]
		month := os.Args[7]
		date = day + "%2F" + month
	} else {
		date = time.Now().Format("2") + "%2F" + time.Now().Format("Jan")
	}

	validate(date)
	confirmed := confirm(date)
	template(date)
	if confirmed {
		if sendRequest(userName, password, jira_add) == 200 {
			cmd := exec.Command("open", jira_open+strings.ToLower(userName))
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			template("")
		}
	} else {
		fmt.Println("ignoring ", date)
	}

}
