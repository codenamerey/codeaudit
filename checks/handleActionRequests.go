package checks

import (
	"codeAudit/models"
	"codeAudit/utils"
	"fmt"
)

func logSyntaxConsistencyScore(score int) {
	if score < 65 {
		utils.ErrorPrintLn(fmt.Sprintf("\nTotal Syntax Consistency Score - - - >  %d%%", score))
	} else if score < 80 {
		utils.WarningPrintLn(fmt.Sprintf("\nTotal Syntax Consistency Score - - - >  %d%%", score))
	} else {
		utils.SuccessPrintLn(fmt.Sprintf("\nTotal Syntax Consistency Score - - - >  %d%%", score))
	}
}

func PerformAllChecks(root_directory string, fileType string, name_convention string, indentation int, char_count int, allow_semicolons bool, preferDouble bool) (models.ConsistencyReport, []string) {
	files := utils.GetCorrectFiles(root_directory, fileType)
	if len(files) == 0 {
		return models.ConsistencyReport{}, files
	}
	result1 := MakeNamingConventionChecks(files, name_convention, fileType)
	result2 := MakeIndentionChecks(files, indentation)
	result3 := MakeCharacterCountChecks(files, char_count)
	result4 := MakeSemiColonChecks(files, allow_semicolons)
	result5 := MakeQuotationStyleChecks(files, preferDouble, fileType)

	results := []models.CompleteCheckResult{result1, result2, result3, result4, result5}
	checksMade := []string{"Variable Name Casing", "Indentation", "Character Count Limit Per Line", "Semi-Colon Usage", "Quoutation Style"}
	issues := []models.IssueData{}
	totalScore := (result1.FinalConsistencyScore + result2.FinalConsistencyScore + result3.FinalConsistencyScore + result4.FinalConsistencyScore + result5.FinalConsistencyScore) / 5

	issues = append(issues, result1.IssuesFound...)
	issues = append(issues, result2.IssuesFound...)
	issues = append(issues, result3.IssuesFound...)
	issues = append(issues, result4.IssuesFound...)
	issues = append(issues, result5.IssuesFound...)

	fullReport := models.ConsistencyReport{IssuesFound: issues, CheckResults: results, Checks: checksMade, CodeBaseConsistencyScore: totalScore}
	for r := 0; r < len(results); r++ {
		resultEntry := results[r]
		utils.InfoPrintLn(fmt.Sprintf("Check \"%s\" score: %d%%", resultEntry.CheckType, resultEntry.FinalConsistencyScore))
	}

	logSyntaxConsistencyScore(fullReport.CodeBaseConsistencyScore)
	return fullReport, files

}

func PerformSemiColonChecks(root_directory string, fileType string, allow_semicolons bool) (models.ConsistencyReport, []string) {
	files := utils.GetCorrectFiles(root_directory, fileType)
	if len(files) == 0 {
		return models.ConsistencyReport{}, files
	}
	result1 := MakeSemiColonChecks(files, allow_semicolons)

	results := []models.CompleteCheckResult{result1}
	checksMade := []string{"Semi-Colon Usage"}
	issues := []models.IssueData{}
	issues = append(issues, result1.IssuesFound...)
	totalScore := result1.FinalConsistencyScore

	fullReport := models.ConsistencyReport{IssuesFound: issues, CheckResults: results, Checks: checksMade, CodeBaseConsistencyScore: totalScore}
	for r := 0; r < len(results); r++ {
		resultEntry := results[r]
		utils.InfoPrintLn(fmt.Sprintf("Check \"%s\" score: %d%%", resultEntry.CheckType, resultEntry.FinalConsistencyScore))
	}
	logSyntaxConsistencyScore(fullReport.CodeBaseConsistencyScore)
	return fullReport, files

}

func PerformCharacterCountChecks(root_directory string, fileType string, char_count int) (models.ConsistencyReport, []string) {
	files := utils.GetCorrectFiles(root_directory, fileType)
	if len(files) == 0 {
		return models.ConsistencyReport{}, files
	}
	result1 := MakeCharacterCountChecks(files, char_count)

	results := []models.CompleteCheckResult{result1}
	checksMade := []string{"Character Count Limit Per Line"}
	issues := []models.IssueData{}
	issues = append(issues, result1.IssuesFound...)
	totalScore := result1.FinalConsistencyScore

	fullReport := models.ConsistencyReport{IssuesFound: issues, CheckResults: results, Checks: checksMade, CodeBaseConsistencyScore: totalScore}
	for r := 0; r < len(results); r++ {
		resultEntry := results[r]
		utils.InfoPrintLn(fmt.Sprintf("Check \"%s\" score: %d%%", resultEntry.CheckType, resultEntry.FinalConsistencyScore))
	}
	logSyntaxConsistencyScore(fullReport.CodeBaseConsistencyScore)
	return fullReport, files
}

func PerformVariableNamingChecks(root_directory string, fileType string, name_convention string) (models.ConsistencyReport, []string) {
	files := utils.GetCorrectFiles(root_directory, fileType)
	if len(files) == 0 {
		return models.ConsistencyReport{}, files
	}
	result1 := MakeNamingConventionChecks(files, name_convention, fileType)

	results := []models.CompleteCheckResult{result1}
	checksMade := []string{"Variable Name Casing"}
	issues := []models.IssueData{}
	issues = append(issues, result1.IssuesFound...)
	totalScore := result1.FinalConsistencyScore

	fullReport := models.ConsistencyReport{IssuesFound: issues, CheckResults: results, Checks: checksMade, CodeBaseConsistencyScore: totalScore}
	for r := 0; r < len(results); r++ {
		resultEntry := results[r]
		utils.InfoPrintLn(fmt.Sprintf("Check \"%s\" score: %d%%", resultEntry.CheckType, resultEntry.FinalConsistencyScore))
	}
	logSyntaxConsistencyScore(fullReport.CodeBaseConsistencyScore)
	return fullReport, files
}

func PerformIndentationChecks(root_directory string, fileType string, indentation int) (models.ConsistencyReport, []string) {
	files := utils.GetCorrectFiles(root_directory, fileType)
	if len(files) == 0 {
		return models.ConsistencyReport{}, files
	}
	result1 := MakeIndentionChecks(files, indentation)

	results := []models.CompleteCheckResult{result1}
	checksMade := []string{"Indentation"}
	issues := []models.IssueData{}
	issues = append(issues, result1.IssuesFound...)
	totalScore := result1.FinalConsistencyScore

	fullReport := models.ConsistencyReport{IssuesFound: issues, CheckResults: results, Checks: checksMade, CodeBaseConsistencyScore: totalScore}
	for r := 0; r < len(results); r++ {
		resultEntry := results[r]
		utils.InfoPrintLn(fmt.Sprintf("Check \"%s\" score: %d%%", resultEntry.CheckType, resultEntry.FinalConsistencyScore))
	}
	logSyntaxConsistencyScore(fullReport.CodeBaseConsistencyScore)
	return fullReport, files
}

func PerformQuotationStyleChecks(root_directory string, fileType string, preferDouble bool) (models.ConsistencyReport, []string) {
	files := utils.GetCorrectFiles(root_directory, fileType)
	if len(files) == 0 {
		return models.ConsistencyReport{}, files
	}
	result1 := MakeQuotationStyleChecks(files, preferDouble, fileType)

	results := []models.CompleteCheckResult{result1}
	checksMade := []string{"Quotation Style"}
	issues := []models.IssueData{}
	issues = append(issues, result1.IssuesFound...)
	totalScore := result1.FinalConsistencyScore

	fullReport := models.ConsistencyReport{IssuesFound: issues, CheckResults: results, Checks: checksMade, CodeBaseConsistencyScore: totalScore}
	for r := 0; r < len(results); r++ {
		resultEntry := results[r]
		utils.InfoPrintLn(fmt.Sprintf("Check \"%s\" score: %d%%", resultEntry.CheckType, resultEntry.FinalConsistencyScore))
	}
	logSyntaxConsistencyScore(fullReport.CodeBaseConsistencyScore)
	return fullReport, files
}
