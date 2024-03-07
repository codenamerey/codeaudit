package checks

import (
	"bufio"
	"codeAudit/models"
	"math"
	"os"
	"strings"

	"github.com/dariubs/percent"
)

// Currently just checks if semi colons are present or not and if they are allowed or not

func lineEndsWithSemiColon(line string) bool {
	if len(line) < 1 {
		return false
	} else {
		return strings.HasSuffix(line, ";")
	}
}

func checkIfSemiColonsPresentInFile(path string, semiColonsAllowed bool) models.CheckResult {
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
		for filescanner.Scan() {
			lineNumber++
			totalLines++
			line := filescanner.Text()
			lineEndsWithSemiColon := lineEndsWithSemiColon(line)

			if semiColonsAllowed {
				checksPassed++
			} else {
				if !semiColonsAllowed && lineEndsWithSemiColon {
					linesFailed = append(linesFailed, lineNumber)
				} else {
					checksPassed++
				}
			}
		}
		consistency := int(math.Round(percent.PercentOf(checksPassed, totalLines)))
		return models.CheckResult{File: path, ChecksMade: totalLines, ChecksPassed: checksPassed, ConsistencyScore: consistency, LinesFailed: linesFailed}
	}
}

func MakeSemiColonChecks(files []string, semiColonsAllowed bool) models.CompleteCheckResult {
	scores := []int{}
	issues := []models.IssueData{}
	for i := 0; i < len(files); i++ {
		if semiColonsAllowed {
			scores = append(scores, 100)
		} else {
			result := checkIfSemiColonsPresentInFile(files[i], semiColonsAllowed)
			scores = append(scores, result.ConsistencyScore)

			for a := 0; a < len(result.LinesFailed); a++ {
				issues = append(issues, models.IssueData{Line: result.LinesFailed[a], File: result.File, Issue: "Semi-colons not allowed"})
			}
		}
	}

	totalScore := 0
	for n := 0; n < len(scores); n++ {
		totalScore = totalScore + scores[n]
	}
	totalScore = totalScore / len(scores)
	return models.CompleteCheckResult{CheckType: "semiColon", FilesChecked: files, ConsistencyScores: scores, FinalConsistencyScore: totalScore, IssuesFound: issues}
}
