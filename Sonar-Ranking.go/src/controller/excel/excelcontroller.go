package excel

import (
	"lib/structs"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

const (
	PROBLEMS int = 4
	RULES    int = 18
)

type ExcelController struct {
	excelFile *excelize.File
}

func New() *ExcelController {
	ec := &ExcelController{
		excelize.NewFile(),
	}
	return ec
}

func (ec *ExcelController) SaveExcel(path string) bool {
	err := ec.excelFile.SaveAs(path)
	if err != nil {
		return false
	}
	return true
}

func (ec *ExcelController) Top10Issues(issues *structs.Issues) bool {
	if len(issues.Rules) != 0 {
		ec.excelFile.SetSheetName("Sheet1", "Top10General")
		setExcelStyleTop10(ec, "Top10General")
		setExcelValuesTop10(ec, "Top10General", issues)
		return true
	} else {
		return false
	}
}

func (ec *ExcelController) Top10Rules(rules *[]structs.Top10Rule) bool {
	if len(*rules) != 0 {
		for i := 0; i < len(*rules); i++ {

			ec.excelFile.SetCellValue("Top10General", getCellLevel("B", i, RULES), (*rules)[i].Lang)
			ec.excelFile.SetCellValue("Top10General", getCellLevel("C", i, RULES), (*rules)[i].Name)
			ec.excelFile.SetCellValue("Top10General", getCellLevel("D", i, RULES), (*rules)[i].Count)

			if i >= 9 {
				break
			}
		}

		return true
	} else {
		return false
	}
}

func (ec *ExcelController) Top10IssuesByProject(issues *structs.Issues, div string) bool {
	if len(issues.Rules) != 0 {
		ec.excelFile.NewSheet(div)
		setExcelStyleTop10(ec, div)
		setExcelValuesTop10(ec, div, issues)
		return true
	} else {
		return false
	}
}

func (ec *ExcelController) Top10RulesByProject(rules *[]structs.Top10Rule, div string) bool {
	if len(*rules) != 0 {
		for i := 0; i < len(*rules); i++ {

			ec.excelFile.SetCellValue(div, getCellLevel("B", i, RULES), (*rules)[i].Lang)
			ec.excelFile.SetCellValue(div, getCellLevel("C", i, RULES), (*rules)[i].Name)
			ec.excelFile.SetCellValue(div, getCellLevel("D", i, RULES), (*rules)[i].Count)

			if i >= 9 {
				break
			}
		}
		return true
	} else {
		return false
	}
}

func (ec *ExcelController) SetExcelStyleOwasp(projects *structs.Projects) {
	name := "OWASP"
	ec.excelFile.NewSheet(name)
	ec.excelFile.MergeCell(name, "C2", "Q2")
	titleStyle, _ := ec.excelFile.NewStyle(`{"fill":{"type":"pattern","color":["#696969"],"pattern":1}, "font":{"bold":true, "color":"#FFFFFF"}, "alignment":{"horizontal":"center"}, "border":[{"type":"left","color":"000000","style":3},{"type":"top","color":"000000","style":3},{"type":"right","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "C2", "Q2", titleStyle)
	ec.excelFile.SetCellValue(name, "C2", "OWASP Report")

	ec.excelFile.SetColWidth(name, "E", "P", 15)
	ec.excelFile.SetColWidth(name, "C", "C", 20)
	subtitleStyle, _ := ec.excelFile.NewStyle(`{"fill":{"type":"pattern","color":["#BEBEBE"],"pattern":1}, "font":{"bold":true}, "border":[{"type":"left","color":"000000","style":3},{"type":"top","color":"000000","style":3},{"type":"bottom","color":"000000","style":3},{"type":"right","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "C3", "Q3", subtitleStyle)
	ec.excelFile.SetCellValue(name, "C3", "Proyectos")
	ec.excelFile.SetCellValue(name, "D3", "Analisis")
	ec.excelFile.SetCellValue(name, "E3", "A1 - Injection")
	ec.excelFile.SetCellValue(name, "F3", "A2 - Broken Authentication")
	ec.excelFile.SetCellValue(name, "G3", "A3 - Sensitive Data Exposure")
	ec.excelFile.SetCellValue(name, "H3", "A4 - XML External Entities (XXE)")
	ec.excelFile.SetCellValue(name, "I3", "A5 - Broken Access Control")
	ec.excelFile.SetCellValue(name, "J3", "A6 - Security Misconfiguration")
	ec.excelFile.SetCellValue(name, "K3", "A7 - Cross-Site Scripting (XSS)")
	ec.excelFile.SetCellValue(name, "L3", "A8 - Insecure Deserialization")
	ec.excelFile.SetCellValue(name, "M3", "A9 - Using Components with Known Vulnerabilities")
	ec.excelFile.SetCellValue(name, "N3", "A10 - Insufficient Logging & Monitoring")
	ec.excelFile.SetCellValue(name, "O3", "Not OWASP")
	ec.excelFile.SetCellValue(name, "P3", "Affected %")
	ec.excelFile.SetCellValue(name, "Q3", "TOTAL")

	leftStyle, _ := ec.excelFile.NewStyle(`{"fill":{"type":"pattern","color":["#EBEBEB"],"pattern":1}, "border":[{"type":"left","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "C4", getCellLevel("C", len(projects.Components)-1, PROBLEMS), leftStyle)

	rightStyle, _ := ec.excelFile.NewStyle(`{"border":[{"type":"right","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "Q4", getCellLevel("Q", len(projects.Components)-1, PROBLEMS), rightStyle)

	bottomleftStyle, _ := ec.excelFile.NewStyle(`{"fill":{"type":"pattern","color":["#EBEBEB"],"pattern":1}, "border":[{"type":"bottom","color":"000000","style":3},{"type":"left","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, getCellLevel("C", len(projects.Components)-1, PROBLEMS), getCellLevel("C", len(projects.Components)-1, PROBLEMS), bottomleftStyle)

	bottomStyle, _ := ec.excelFile.NewStyle(`{"border":[{"type":"bottom","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, getCellLevel("D", len(projects.Components)-1, PROBLEMS), getCellLevel("Q", len(projects.Components)-1, PROBLEMS), bottomStyle)

	bottomrightStyle, _ := ec.excelFile.NewStyle(`{"border":[{"type":"bottom","color":"000000","style":3},{"type":"right","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, getCellLevel("Q", len(projects.Components)-1, PROBLEMS), getCellLevel("Q", len(projects.Components)-1, PROBLEMS), bottomrightStyle)
}

func (ec *ExcelController) AddProjectReport(report *structs.SecurityReport, index int, project string) {
	letters := []string{"E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}
	total := 0
	affected := 0
	ec.excelFile.SetCellValue("OWASP", getCellLevel("C", index, PROBLEMS), project)

	for i, elem := range report.Categories {
		ec.excelFile.SetCellValue("OWASP", getCellLevel(letters[i], index, PROBLEMS), elem.Vulnerabilities)
		if elem.Vulnerabilities != 0 && elem.Category != "unknown" {
			affected += 10
		}
		total += elem.Vulnerabilities
	}
	ec.excelFile.SetCellValue("OWASP", getCellLevel("Q", index, PROBLEMS), total)
	ec.excelFile.SetCellValue("OWASP", getCellLevel("P", index, PROBLEMS), strconv.Itoa(affected)+"%")
}

func setExcelStyleTop10(ec *ExcelController, name string) {
	subtitleStyle, _ := ec.excelFile.NewStyle(`{"fill":{"type":"pattern","color":["#BEBEBE"],"pattern":1}, "font":{"bold":true}, "border":[{"type":"left","color":"000000","style":3},{"type":"top","color":"000000","style":3},{"type":"bottom","color":"000000","style":3},{"type":"right","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "B3", "D3", subtitleStyle)
	_ = ec.excelFile.SetCellStyle(name, "B17", "D17", subtitleStyle)

	ec.excelFile.SetColWidth(name, "D", "D", 10)
	ec.excelFile.SetColWidth(name, "C", "C", 70)
	ec.excelFile.SetColWidth(name, "B", "B", 15)
	ec.excelFile.MergeCell(name, "B2", "D2")
	ec.excelFile.MergeCell(name, "B16", "D16")
	titleStyle, _ := ec.excelFile.NewStyle(`{"fill":{"type":"pattern","color":["#696969"],"pattern":1}, "font":{"bold":true, "color":"#FFFFFF"}, "alignment":{"horizontal":"center"}, "border":[{"type":"left","color":"000000","style":3},{"type":"top","color":"000000","style":3},{"type":"right","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "B2", "D2", titleStyle)
	_ = ec.excelFile.SetCellStyle(name, "B16", "D16", titleStyle)
	ec.excelFile.SetCellValue(name, "B2", "Top10 General Problemas")
	ec.excelFile.SetCellValue(name, "B16", "Top10 General Reglas")

	leftStyle, _ := ec.excelFile.NewStyle(`{"fill":{"type":"pattern","color":["#EBEBEB"],"pattern":1}, "border":[{"type":"left","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "B4", "B13", leftStyle)
	_ = ec.excelFile.SetCellStyle(name, "B18", "B27", leftStyle)

	rightStyle, _ := ec.excelFile.NewStyle(`{"border":[{"type":"right","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "D4", "D13", rightStyle)
	_ = ec.excelFile.SetCellStyle(name, "D18", "D27", rightStyle)

	bottomleftStyle, _ := ec.excelFile.NewStyle(`{"fill":{"type":"pattern","color":["#EBEBEB"],"pattern":1}, "border":[{"type":"bottom","color":"000000","style":3},{"type":"left","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "B13", "B13", bottomleftStyle)
	_ = ec.excelFile.SetCellStyle(name, "B27", "B27", bottomleftStyle)

	bottomStyle, _ := ec.excelFile.NewStyle(`{"border":[{"type":"bottom","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "C13", "C13", bottomStyle)
	_ = ec.excelFile.SetCellStyle(name, "C27", "C27", bottomStyle)

	bottomrightStyle, _ := ec.excelFile.NewStyle(`{"border":[{"type":"bottom","color":"000000","style":3},{"type":"right","color":"000000","style":3}]}`)
	_ = ec.excelFile.SetCellStyle(name, "D13", "D13", bottomrightStyle)
	_ = ec.excelFile.SetCellStyle(name, "D27", "D27", bottomrightStyle)

	ec.excelFile.SetCellValue(name, "B3", "Lenguaje")
	ec.excelFile.SetCellValue(name, "C3", "Nombre")
	ec.excelFile.SetCellValue(name, "D3", "Problemas")
	ec.excelFile.SetCellValue(name, "B17", "Lenguaje")
	ec.excelFile.SetCellValue(name, "C17", "Nombre")
	ec.excelFile.SetCellValue(name, "D17", "Reglas")

	ec.excelFile.MergeCell(name, "F2", "G2")
	_ = ec.excelFile.SetCellStyle(name, "F2", "G2", titleStyle)
	ec.excelFile.SetCellValue(name, "F2", "Severities")

	_ = ec.excelFile.SetCellStyle(name, "F3", "G3", subtitleStyle)
	ec.excelFile.SetCellValue(name, "F3", "Severity")
	ec.excelFile.SetCellValue(name, "G3", "Count")

	_ = ec.excelFile.SetCellStyle(name, "F4", "F7", leftStyle)
	_ = ec.excelFile.SetCellStyle(name, "F8", "F8", bottomleftStyle)
	_ = ec.excelFile.SetCellStyle(name, "G4", "G7", rightStyle)
	_ = ec.excelFile.SetCellStyle(name, "G8", "G8", bottomrightStyle)

	ec.excelFile.SetColWidth(name, "I", "I", 20)
	ec.excelFile.MergeCell(name, "I2", "J2")
	_ = ec.excelFile.SetCellStyle(name, "I2", "J2", titleStyle)
	ec.excelFile.SetCellValue(name, "I2", "Types")

	_ = ec.excelFile.SetCellStyle(name, "I3", "J3", subtitleStyle)
	ec.excelFile.SetCellValue(name, "I3", "Type")
	ec.excelFile.SetCellValue(name, "J3", "Count")

	_ = ec.excelFile.SetCellStyle(name, "I4", "I6", leftStyle)
	_ = ec.excelFile.SetCellStyle(name, "I7", "I7", bottomleftStyle)
	_ = ec.excelFile.SetCellStyle(name, "J4", "J6", rightStyle)
	_ = ec.excelFile.SetCellStyle(name, "J7", "J7", bottomrightStyle)
}

func setExcelValuesTop10(ec *ExcelController, sheet string, issues *structs.Issues) {

	for index, issue := range issues.Facets[2].Values {
		name, lang := searchRuleDefinition(&issues.Rules, issue.Val)

		ec.excelFile.SetCellValue(sheet, getCellLevel("B", index, PROBLEMS), lang)
		ec.excelFile.SetCellValue(sheet, getCellLevel("C", index, PROBLEMS), name)
		ec.excelFile.SetCellValue(sheet, getCellLevel("D", index, PROBLEMS), issue.Count)

		if index == 9 {
			break
		}
	}

	for index, severity := range issues.Facets[0].Values {
		ec.excelFile.SetCellValue(sheet, getCellLevel("F", index, PROBLEMS), severity.Val)
		ec.excelFile.SetCellValue(sheet, getCellLevel("G", index, PROBLEMS), severity.Count)
	}

	for index, tipo := range issues.Facets[1].Values {
		ec.excelFile.SetCellValue(sheet, getCellLevel("I", index, PROBLEMS), tipo.Val)
		ec.excelFile.SetCellValue(sheet, getCellLevel("J", index, PROBLEMS), tipo.Count)
	}

	ec.excelFile.AddChart(sheet, "M2", "{\"type\":\"pie\",\"series\":[{\"name\":\""+sheet+"!$F$2\",\"categories\":\""+sheet+"!$F$4:$F$8\",\"values\":\""+sheet+"!$G$4:$G$8\"}],\"format\":{\"x_scale\":1.0,\"y_scale\":1.0,\"x_offset\":15,\"y_offset\":10,\"print_obj\":true,\"lock_aspect_ratio\":false,\"locked\":false},\"legend\":{\"position\":\"bottom\",\"show_legend_key\":false},\"title\":{\"name\":\"Severities\"},\"plotarea\":{\"show_bubble_size\":true,\"show_cat_name\":false,\"show_leader_lines\":false,\"show_percent\":true,\"show_series_name\":false,\"show_val\":false},\"show_blanks_as\":\"gap\"}")
	ec.excelFile.AddChart(sheet, "F11", "{\"type\":\"pie\",\"series\":[{\"name\":\""+sheet+"!$I$2\",\"categories\":\""+sheet+"!$I$4:$I$7\",\"values\":\""+sheet+"!$J$4:$J$7\"}],\"format\":{\"x_scale\":1.0,\"y_scale\":1.0,\"x_offset\":15,\"y_offset\":10,\"print_obj\":true,\"lock_aspect_ratio\":false,\"locked\":false},\"legend\":{\"position\":\"bottom\",\"show_legend_key\":false},\"title\":{\"name\":\"Type\"},\"plotarea\":{\"show_bubble_size\":true,\"show_cat_name\":false,\"show_leader_lines\":false,\"show_percent\":true,\"show_series_name\":false,\"show_val\":false},\"show_blanks_as\":\"gap\"}")
}

func searchRuleDefinition(rules *[]structs.Rule, searching string) (string, string) {
	var name string
	var lang string
	for _, rule := range *rules {
		if rule.Key == searching {
			name = rule.Name
			lang = rule.LangName
			break
		}
	}
	return name, lang
}

func getCellLevel(column string, i int, tipo int) string {
	return strings.Join([]string{column, strconv.Itoa(i + tipo)}, "")
}
