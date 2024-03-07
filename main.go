package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"codeAudit/checks"
	"codeAudit/models"
	"codeAudit/utils"
)

func main() {
	downloadsFolderPath := "./"
	userHomeDir, userHomeErr := os.UserHomeDir()
	shouldGenerateFailureReport := true
	if userHomeErr != nil {
		utils.WarningPrintLn(fmt.Sprintf("Error occurred finding home directory: %s", userHomeErr.Error()))
	}
	system := utils.GetOs()
	if system == "darwin" {
		downloadsFolderPath = fmt.Sprintf("%s/Downloads", userHomeDir)
	}

	if system == "windows" {
		downloadsFolderPath = fmt.Sprintf("%s/Downloads", userHomeDir)
	}

	if system == "linux" {
		downloadsFolderPath = fmt.Sprintf("/home%s/Downloads", userHomeDir)
	}
	var fileType string
	var check string
	var verboseMode bool = false
	root_directory := "./"
	defaultCharacterLimit := 100
	defaultVariableNamingConvention := "camel"
	defaultIndentationSpaces := 2
	defaultSemiColonUsage := true

	terminalArgs := os.Args[1:]
	if len(terminalArgs) == 0 {
		fileType = ".js"
		check = "all"
	}

	if len(terminalArgs) >= 1 {
		if !utils.IsConfigFlag(terminalArgs[0]) {
			check = terminalArgs[0]
		}
		fileType = ".js"
	}

	if len(terminalArgs) >= 2 {

		if !utils.IsConfigFlag(terminalArgs[1]) {
			if !utils.IsConfigFlag(terminalArgs[0]) {
				check = terminalArgs[0]
			}
			isFlag, err := regexp.MatchString("-", terminalArgs[1])
			if err != nil {
				utils.ErrorPrintLn("Unexpected error occurred while checking for flag")
				return
			}

			if isFlag {
				flagValue := strings.Split(terminalArgs[1], "-")[1]
				options := strings.Split(flagValue, "")
				switch options[0] {
				case "j":
					fileType = ".js"
				case "p":
					fileType = ".py"
				// case "t":
				// 	fileType = ".ts"
				// case "g":
				// 	fileType = ".go"
				default:
					utils.WarningPrintLn("Invalid file type provided. Using default file type (.js)")
					fileType = ".js"
				}
			}
		}
	}

	if len(terminalArgs) >= 3 {
		if !utils.IsConfigFlag(terminalArgs[2]) {
			root_directory = terminalArgs[2]
		}
	}

	for a := 0; a < len(terminalArgs); a++ {
		arg := terminalArgs[a]
		if utils.IsConfigFlag(arg) {
			configTuple := strings.Split(arg, "=")
			option := strings.Split(configTuple[0], "--")[1]
			value := configTuple[1]

			switch option {
			case "i":
				newIndentation, err := strconv.Atoi(value)
				if err != nil {
					utils.WarningPrintLn(fmt.Sprintf("Input %s wasn't able to be converted to a number to use as new number of indentation. Used default value of 2\n", value))
				} else {
					defaultIndentationSpaces = newIndentation
				}

			case "c":
				newCharacterLineLimit, err := strconv.Atoi(value)
				if err != nil {
					utils.WarningPrintLn(fmt.Sprintf("Input %s wasn't able to be converted to a number to use as new character limit. Used default value of 100", value))
				} else {
					defaultCharacterLimit = newCharacterLineLimit
				}
			case "n":
				isValidConventions := utils.IsValidNamingConvention(value)
				if !isValidConventions {
					utils.WarningPrintLn(fmt.Sprintf("%s isn't a valid or supported variable naming convention. See docs", value))
				} else {
					defaultVariableNamingConvention = value
				}
			case "r":
				willGenerateReport, err := strconv.ParseBool(value)
				if err != nil {
					utils.WarningPrintLn(fmt.Sprintf("%s isn't a valid boolean value", value))
				} else {
					shouldGenerateFailureReport = willGenerateReport
				}

			case "v":
				verbose, err := strconv.ParseBool(value)
				if err != nil {
					utils.WarningPrintLn(fmt.Sprintf("%s isn't a valid boolean value", value))
				} else {
					verboseMode = verbose
				}

			case "s":
				newSemiColonUsage, err := strconv.ParseBool(value)
				if err != nil {
					utils.WarningPrintLn(fmt.Sprintf("%s isn't a valid boolean value", value))
				} else {
					defaultSemiColonUsage = newSemiColonUsage
				}
			}
		}
	}

	if verboseMode {
		utils.WarningPrintLn("\nVerbose mode enabled\n")
	}
	utils.DefaultPrintLn(fmt.Sprintf("Directory testing: %s\n", root_directory))

	var report models.ConsistencyReport
	var foundFiles []string
	checksRan := true

	switch check {
	case "all":
		result, files := checks.PerformAllChecks(root_directory, fileType, defaultVariableNamingConvention, defaultIndentationSpaces, defaultCharacterLimit, defaultSemiColonUsage)
		report = result
		foundFiles = files
	case "indent":
		result, files := checks.PerformIndentationChecks(root_directory, fileType, defaultIndentationSpaces)
		report = result
		foundFiles = files
	case "naming":
		result, files := checks.PerformVariableNamingChecks(root_directory, fileType, defaultVariableNamingConvention)
		report = result
		foundFiles = files
	case "char":
		result, files := checks.PerformCharacterCountChecks(root_directory, fileType, defaultCharacterLimit)
		report = result
		foundFiles = files
	case "semi":
		result, files := checks.PerformSemiColonChecks(root_directory, fileType, defaultSemiColonUsage)
		report = result
		foundFiles = files
	default:
		checksRan = false
		utils.ErrorPrintLn("unexpected command please try again")
	}

	if len(foundFiles) == 0 {
		utils.ErrorPrintLn("No files found in the specified directory matching the file type provided.")
		return
	}

	if verboseMode {
		utils.WarningPrintLn("\nConfigurations Values used\n")
		utils.DefaultPrintLn(fmt.Sprintf("Line Character Limit: %d", defaultCharacterLimit))
		utils.DefaultPrintLn(fmt.Sprintf("Variable Naming Convention: %s", defaultVariableNamingConvention))
		utils.DefaultPrintLn(fmt.Sprintf("Number of Indentation Spaces: %d", defaultIndentationSpaces))
		utils.DefaultPrintLn(fmt.Sprintf("Allow Lines Ending In Semi-Colon: %t\n", defaultSemiColonUsage))
	}

	if checksRan {
		utils.DefaultPrintLn(fmt.Sprintf("Number of files tested: %d", len(foundFiles)))
		utils.DefaultPrintLn(fmt.Sprintf("Number of issues found: %d\n", len(report.IssuesFound)))
		if shouldGenerateFailureReport {
			reportCreated := utils.GenerateFailureReport(report, downloadsFolderPath)
			if !reportCreated {
				utils.WarningPrintLn("Error occurred while generating report")
			} else {
				utils.SuccessPrintLn("Report generated successfully\n")
				utils.DefaultPrintLn(fmt.Sprintf("Full Report can be found at %s/CodeAuditReport.csv\n", downloadsFolderPath))
			}
		}
	}
}
