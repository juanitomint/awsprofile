// Copyright 2020 Juan Ignacio Borda <juanignacioborda@gmail.com>. All rights reserved.
// Use of this source code is governed by a MIT license that can
// be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"github.com/spf13/viper"
)

func main() {
	var re = regexp.MustCompile(`(?m)\[.*\]`)
	rows := []string{}
	profiles := []string{}
	awsCredentials := os.Getenv("HOME") + "/.aws/credentials"
	content, err := ioutil.ReadFile(awsCredentials)
	if err != nil {
		log.Fatal(err)
	}
	for i, match := range re.FindAllString(string(content), -1) {
		match = strings.Replace(match, "[", "", -1)
		match = strings.Replace(match, "]", "", -1)
		profiles = append(profiles, match)
		rows = append(rows, fmt.Sprintf("[%d] %s", i, match))

	}

	// fmt.Println(">>")
	selected := drawList(rows)
	// selected := 3
	// fmt.Println("echo \"Profile selected:", rows[selected], "\";\")
	if selected >= 0 {
		iniKey := profiles[selected]
		exportConfig(awsCredentials, iniKey)
		fmt.Println("export", "AWS_PROFILE"+"="+iniKey+"\n")
	}

	os.Exit(0)

}

func exportConfig(awsCredentials string, iniKey string) {
	mmap := map[string]string{
		"aws_access_key_id":     "AWS_ACCESS_KEY_ID",
		"aws_secret_access_key": "AWS_SECRET_ACCESS_KEY",
		"region":                "AWS_DEFAULT_REGION",
	}
	viper.SetConfigType("ini")
	viper.SetConfigName("credentials")               // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("HOME") + "/.aws") // optionally look for config in the working directory
	err := viper.ReadInConfig()                      // Find and read the config file
	if err != nil {                                  // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	for key, value := range mmap {
		fmt.Println("export", strings.ToUpper(value)+"="+viper.GetString(iniKey+"."+key)+"\n")

		// os.Setenv(strings.ToUpper(value), viper.GetString(iniKey+"."+key))
	}

}

func drawList(myRows []string) int {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()
	l.Title = "Select AWS Credentials"
	l.Rows = myRows
	l.TextStyle = ui.NewStyle(ui.ColorBlue)
	l.WrapText = false
	l.SetRect(0, 0, 40, 12)

	ui.Render(l)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<Enter>":
			// defer ui.Close()
			return l.SelectedRow

		case "q", "<C-c>":
			return -1
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(l)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
