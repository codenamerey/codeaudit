package checks

import (
	"bufio"
	"fmt"
	"os"

	"math"
	"regexp"
	"strings"

	"codeAudit/models"

	"github.com/dariubs/percent"
)

// This methods works with both python and javascript files
func RunNamingConventionCheck(filePath string, style string, file_extension string) models.CheckResult {
	var rg *regexp.Regexp
	camelCaseRegex := regexp.MustCompile("^[a-z0-9]+([A-Z].*)*$")
	snake_case_regex := regexp.MustCompile(`^[a-z0-9](\_[a-z0-9])*$`)
	pascalCaseRegex := regexp.MustCompile(`^[A-Z][a-z0-9]*([A-Z].*)*$`)

	if style == "camel" {
		rg = camelCaseRegex
	}
	if style == "snake" {
		rg = snake_case_regex
	}

	if style == "pascal" {
		rg = pascalCaseRegex
	}
	f, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)
	varibleCount := 0
	checksPassed := 0
	lines_failed := []int{}
	lineNumber := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		lineNumber = lineNumber + 1
		// this section of code checks for variable declarations only in javascript

		if file_extension == ".js" {
			line_items := strings.Split(line, " ")
			for i := 0; i < len(line_items); i++ {
				// check if line contains a variable declaration
				if line_items[i] == "const" || line_items[i] == "let" || line_items[i] == "var" {
					varibleCount = varibleCount + 1
					passesChecks := rg.MatchString(line_items[i+1])
					if passesChecks {
						checksPassed = checksPassed + 1
					} else {
						lines_failed = append(lines_failed, lineNumber)
					}
				}

				if line_items[i] == "function" {
					functionName := strings.Split(line_items[i+1], "(")[0]
					varibleCount = varibleCount + 1
					passesChecks := rg.MatchString(functionName)
					if passesChecks {
						checksPassed = checksPassed + 1
					} else {
						lines_failed = append(lines_failed, lineNumber)
					}
				}
			}
		}
		// do the same for python
		if file_extension == ".py" {
			python_function_definition_rg := regexp.MustCompile(`def\s+([[:word:]]+)\s*\((?:\s*[[:word:]]+\s*,\s*)*\w*\)\s*:`)
			python_variable_declaration_rg := regexp.MustCompile(`(\?\<\!\.)\b(?:[^\d\W]\w*\s*,\s*)*[^\d\W]\w*\s*`)

			variable_declarations := python_variable_declaration_rg.FindAllString(line, -1)
			function_definitions := python_function_definition_rg.FindAllString(line, -1)
			names_found := append(variable_declarations, function_definitions...)

			if len(names_found) > 0 {
				for i := 0; i < len(names_found); i++ {
					varibleCount = varibleCount + 1
					passesChecks := rg.MatchString(names_found[i])
					if passesChecks {
						checksPassed = checksPassed + 1
					} else {
						lines_failed = append(lines_failed, lineNumber)
					}
				}
			}

		}

	}
	consistency := int(math.Round(percent.PercentOf(checksPassed, varibleCount)))
	for a := 0; a < len(lines_failed); a++ {
	}
	result := models.CheckResult{File: filePath, ChecksMade: varibleCount, ChecksPassed: checksPassed, ConsistencyScore: consistency, LinesFailed: lines_failed}
	return result
}

func MakeNamingConventionChecks(files []string, variableNamingConvention string, file_extension string) models.CompleteCheckResult {
	scores := []int{}
	issues := []models.IssueData{}
	for i := 0; i < len(files); i++ {
		result1 := RunNamingConventionCheck(files[i], variableNamingConvention, file_extension)
		scores = append(scores, result1.ConsistencyScore)
		for n := 0; n < len(result1.LinesFailed); n++ {
			issueMessage := fmt.Sprintf("Variable or function found on line %d doesn't meet the specified naming convention", result1.LinesFailed[n])
			issue := models.IssueData{Line: result1.LinesFailed[n], File: files[i], Issue: issueMessage}
			issues = append(issues, issue)
		}

	}
	totalScore := 0
	for n := 0; n < len(scores); n++ {
		totalScore = totalScore + scores[n]
	}
	totalScore = totalScore / len(scores)
	completeCheckResult := models.CompleteCheckResult{CheckType: "variable_naming", FilesChecked: files, ConsistencyScores: scores, FinalConsistencyScore: totalScore, IssuesFound: issues}
	return completeCheckResult
}
