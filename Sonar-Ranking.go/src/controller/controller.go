package controller

import (
	"controller/excel"
	"fmt"
	sonarapi "lib/api"
	"lib/structs"
	"log"
	"sort"
	"strings"
)

type Controller struct {
	sonarapi *sonarapi.ApiClient
	ec       *excel.ExcelController
}

func New(URLbase string, User string) *Controller {
	controller := &Controller{
		sonarapi.New(URLbase, User),
		excel.New(),
	}
	return controller
}

func (controller *Controller) SaveExcel(path string) {
	if !controller.ec.SaveExcel(path) {
		log.Fatalln("Couldn't save Excel")
	} else {
		fmt.Println("Excel saved!")
	}
}

func (controller *Controller) FetchProjects() *structs.Projects {
	return controller.sonarapi.RetrieveProjects()

	// fmt.Println("Nº of projects: ", len(projects.Components))
	// fmt.Println("---------------------------------------------------------------------")
	// for _, project := range projects.Components {
	// 	fmt.Println("Name: ", project.Name)
	// 	fmt.Println("Id:   ", project.Id)
	// 	fmt.Println("Key:  ", project.Key)
	// 	fmt.Println("---------------------------------------------------------------------")
	// }

}

func (controller *Controller) MakeGeneralTop10() {
	fmt.Print("  > Calculating top general 10 issues... ")
	issues := controller.sonarapi.RetrieveIssues()
	if controller.ec.Top10Issues(issues) {
		fmt.Println("Done!")
		fmt.Print("  > Calculating top 10 rules... ")

		projects := controller.sonarapi.RetrieveProjects()
		rules := make([]structs.Top10Rule, 0)

		for _, project := range projects.Components {
			issues = controller.sonarapi.RetrieveIssuesByProjects(project.Key)
			checkRules(&rules, &issues.Rules)
		}

		sort.SliceStable(rules, func(i, j int) bool {
			return rules[i].Count > rules[j].Count
		})

		if controller.ec.Top10Rules(&rules) {
			fmt.Println("Done!")
		}else {
			fmt.Println("Some kind of error ocurred...")
		}
		// for _, elem := range rules {
		// 	fmt.Println(elem.Lang, elem.Name, elem.Count)
		// }

	} else {
		fmt.Println("No issues found in SonarQube")
	}
}

func (controller *Controller) MakeTop10ByGroup(projects string, div string) {
	fmt.Print("  > Calculating top 10 issues for ", div, "... ")

	issues := controller.sonarapi.RetrieveIssuesByProjects(projects)
	if controller.ec.Top10IssuesByProject(issues, div) {
		fmt.Println("Done!")
		fmt.Print("  > Calculating top 10 rules for ", div, "... ")

		rules := make([]structs.Top10Rule, 0)
		p := strings.Split(projects, ",")

		for _, elem := range p {
			issues = controller.sonarapi.RetrieveIssuesByProjects(elem)
			checkRules(&rules, &issues.Rules)
		}

		sort.SliceStable(rules, func(i, j int) bool {
			return rules[i].Count > rules[j].Count
		})

		controller.ec.Top10RulesByProject(&rules, div)
		fmt.Println("Done!")
	} else {
		fmt.Println("No issues found in", div)
	}
}

func (controller *Controller) GenerateSecurityReport(projects *structs.Projects) {
	fmt.Print("Generating OWASP Report... ")
	controller.ec.SetExcelStyleOwasp(projects)
	for index, project := range projects.Components {
		report := controller.sonarapi.RetrieveSecurityReport(project.Key)
		controller.ec.AddProjectReport(report, index, project.Name)
	}
	fmt.Println("Done!")
}

func checkRules(slice *[]structs.Top10Rule, rules *[]structs.Rule) {
	inserted := false
	for _, elem := range *rules {
		for index, rule := range *slice {
			if rule.Key == elem.Key {
				(*slice)[index].Count += 1
				inserted = true
				break
			}
		}
		if !inserted {
			*slice = append(*slice, structs.Top10Rule{elem.Key, elem.Name, elem.LangName, 1})
			inserted = false
		}
	}
}
