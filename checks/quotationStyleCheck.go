package checks

import (
	"bufio"
	"codeAudit/models"
	"codeAudit/utils"
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/dariubs/percent"
)

func shouldSkipLine(line string, fileType string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return true
	}
	switch fileType {
	case ".js":
		return strings.HasPrefix(trimmed, "//")
	case ".py":
		return strings.HasPrefix(trimmed, "#") || strings.Contains(trimmed, `'''`) || strings.Contains(trimmed, `"""`)
	}
	return false
}

func findIncorrectQuotes(line string, preferDouble bool) []int {
	var incorrectIndices []int
	inSingleQuote := false
	inDoubleQuote := false
	escaped := false
	inTemplateLiteral := false
	inRegExp := false
	inRegexCharClass := false

	for i, char := range line {
		if i > 0 && char == '/' && line[i-1] == '/' && !inSingleQuote && !inDoubleQuote && !inTemplateLiteral && !inRegExp {
			break
		}

		if escaped {
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if char == '`' && !escaped {
			inTemplateLiteral = !inTemplateLiteral
			continue
		}

		if inTemplateLiteral {
			continue
		}

		if !inSingleQuote && !inDoubleQuote && char == '/' && i > 0 {
			prevChar := line[i-1]
			if prevChar == '=' || prevChar == '(' || prevChar == ',' || prevChar == ':' ||
				prevChar == '[' || prevChar == '?' || prevChar == ';' || prevChar == '{' ||
				prevChar == '&' || prevChar == '|' || prevChar == '!' {
				inRegExp = true
				continue
			}
		}

		if inRegExp && char == '/' && !escaped {
			inRegExp = false
			continue
		}

		if inRegExp {
			continue
		}

		if inRegExp {
			if char == '[' && !inRegexCharClass {
				inRegexCharClass = true
				continue
			}
			if char == ']' && inRegexCharClass {
				inRegexCharClass = false
				continue
			}
			if char == '/' && !escaped && !inRegexCharClass {
				inRegExp = false
				continue
			}
			continue
		}

		if char == '\'' && !inDoubleQuote {
			inSingleQuote = !inSingleQuote
			if inSingleQuote && preferDouble {
				incorrectIndices = append(incorrectIndices, i)
			}
		} else if char == '"' && !inSingleQuote {
			inDoubleQuote = !inDoubleQuote
			if inDoubleQuote && !preferDouble {
				incorrectIndices = append(incorrectIndices, i)
			}
		}
	}

	return incorrectIndices
}

func lineHasCorrectQuotationStyle(line string, fileType string, preferDouble bool) (bool, []int) {
	if shouldSkipLine(line, fileType) {
		return true, []int{}
	}
	incorrectIndices := findIncorrectQuotes(line, preferDouble)
	return len(incorrectIndices) == 0, incorrectIndices
}

func runQuotationStyleCheck(path string, fileType string, preferDouble bool) models.CheckResult {
	file, fileOpenError := os.Open(path)

	if fileOpenError != nil {
		panic(fileOpenError)
	} else {
		defer file.Close()
		checksPassed := 0
		linesFailed := []int{}
		totalLines := 0
		filescanner := bufio.NewScanner(file)
		filescanner.Split(bufio.ScanLines)
		lineNumber := 0

		inMultilineString := false
		inJsBlockComment := false

		for filescanner.Scan() {
			lineNumber++
			totalLines++
			line := filescanner.Text()

			if strings.TrimSpace(line) == "" {
				checksPassed++
				continue
			}

			if fileType == ".py" {
				if strings.Contains(line, "'''") || strings.Contains(line, "\"\"\"") {
					inMultilineString = !inMultilineString
					checksPassed++
					continue
				}

				if inMultilineString {
					checksPassed++
					continue
				}
			}

			if fileType == ".js" {
				if !inJsBlockComment && strings.Contains(line, "/*") {
					inJsBlockComment = true

					if strings.Contains(line, "*/") {
						inJsBlockComment = false
					}

					checksPassed++
					continue
				}

				if inJsBlockComment {
					if strings.Contains(line, "*/") {
						inJsBlockComment = false
					}
					checksPassed++
					continue
				}
			}

			isCorrect, _ := lineHasCorrectQuotationStyle(line, fileType, preferDouble)
			if isCorrect {
				checksPassed++
			} else {
				linesFailed = append(linesFailed, lineNumber)
			}
		}
		consistency := int(math.Round(percent.PercentOf(checksPassed, totalLines)))
		return models.CheckResult{File: path, ChecksMade: totalLines, ChecksPassed: checksPassed, ConsistencyScore: consistency, LinesFailed: linesFailed}
	}
}

func MakeQuotationStyleChecks(files []string, preferDouble bool, fileType string) models.CompleteCheckResult {
	scores := []int{}
	issues := []models.IssueData{}

	quotationStyle := "double"
	if !preferDouble {
		quotationStyle = "single"
	}

	for i := 0; i < len(files); i++ {
		utils.InfoPrintLn(fmt.Sprintf("Performing quotation style checks on file %s", files[i]))
		result := runQuotationStyleCheck(files[i], fileType, preferDouble)
		scores = append(scores, result.ConsistencyScore)
		for a := 0; a < len(result.LinesFailed); a++ {
			issues = append(issues, models.IssueData{Line: result.LinesFailed[a], File: result.File, Issue: fmt.Sprintf("%s quotation marks are required but found incorrect style", quotationStyle)})
		}

	}

	totalScore := 0
	if len(scores) > 0 {
		for n := 0; n < len(scores); n++ {
			totalScore = totalScore + scores[n]
		}
		totalScore = totalScore / len(scores)
	}
	completeCheckResult := models.CompleteCheckResult{CheckType: "quotation_style", FilesChecked: files, ConsistencyScores: scores, FinalConsistencyScore: totalScore, IssuesFound: issues}
	return completeCheckResult
}
