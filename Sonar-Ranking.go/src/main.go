package main

import (
	"controller"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/magiconair/properties"
)

func main() {
	clearTerminal()

	prop := flag.String("p", "", "path to properties file")
	output := flag.String("o", "./Sonar.xlsx", "path to save output Excel")
	flag.Parse()

	if len(*prop) <= 1 || len(*output) <= 1 {
		log.Output(1, "Did you miss some argument?")
		log.Fatalln("try -h to fetch help")
	}

	p := properties.MustLoadFile(*prop, properties.UTF8)
	URL := p.MustGetString("URL")
	token := p.MustGetString("token")
	groups := p.MustGetString("groups")
	// var m map[string]string
	m := make(map[string]string)

	if groups != "Empty" {
		divisions := strings.Split(groups, ",")

		for _, div := range divisions {
			m[div] = p.MustGetString(div)
		}
	}

	fmt.Println("----------------------------")
	fmt.Println("PROPERTIES:")
	fmt.Println("----------------------------")
	fmt.Println("User token:", token)
	fmt.Println("SonarQube server @", URL)
	fmt.Println("----------------------------")
	fmt.Println("PROJECTS:")
	fmt.Println("----------------------------")

	for key, value := range m {
		fmt.Println(key, ":", value)
	}
	fmt.Println("----------------------------")

	c := controller.New(URL, token)

	c.MakeGeneralTop10()
	if len(m) != 0 {
		for key, value := range m {
			fmt.Println("----------------------------")
			fmt.Println(key)
			fmt.Println("----------------------------")

			c.MakeTop10ByGroup(value, key)
		}
	}
	fmt.Println("----------------------------")
	c.GenerateSecurityReport(c.FetchProjects())
	fmt.Println("----------------------------")
	c.SaveExcel(*output)
	fmt.Println("")
}

func clearTerminal() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}
