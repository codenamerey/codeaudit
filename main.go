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

// TODO: Update readme to include the new flags and options and python support. Remove depreciated notes
// TODO: Create more in depth mock files for testing in mock_directory in tests directory
// TODO: Add unit tests for all checks and utils.
// TODO: Instead of using a bunch of if statements and setting defaults, require the user provides atleast the file type, the check type and the root directory
// TODO: Add a config flag for allowing the user to specify if the report should be save in cwd or in the downloads folder. Will need to update logic to allow for this. Currently only saves in downloads folder
// TODO: Allow users to specify a config file that contains necessary options for the checks

func main() {
	downloadsFolderPath := "./"
	userHomeDir, userHomeErr := os.UserHomeDir()
	shouldGenerateFailureReport := true
	if userHomeErr != nil {
		println("Error occurred finding home directory: r", userHomeErr.Error())
	}
	fmt.Printf("home dir %s", userHomeDir)
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
	root_directory := "./"
	defaultCharacterLimit := 100
	defaultVariableNamingConvention := "camel"
	defaultIndentationSpaces := 2
	defaultSemiColonUsage := true

	terminalArgs := os.Args[1:]
	println("terminal arguments received:", terminalArgs)
	if len(terminalArgs) == 0 {
		fileType = ".js"
		check = "all"
	}

	if len(terminalArgs) >= 1 {
		if !utils.IsConfigOption(terminalArgs[0]) {
			check = terminalArgs[0]
		}
		fileType = ".js"
	}

	if len(terminalArgs) >= 2 {

		if !utils.IsConfigOption(terminalArgs[1]) {
			if !utils.IsConfigOption(terminalArgs[0]) {
				check = terminalArgs[0]
			}
			isFlag, err := regexp.MatchString("-", terminalArgs[1])
			if err != nil {
				fmt.Println("err 2")

				panic(err)
			}

			if isFlag {
				flagValue := strings.Split(terminalArgs[1], "-")[1]
				options := strings.Split(flagValue, "")
				switch options[0] {
				case "j":
					fileType = ".js"
				case "p":
					fileType = ".py"
				case "t":
					fileType = ".ts"
				case "g":
					fileType = ".go"
				default:
					fileType = ".js"
				}

			} else {
				fmt.Println("err 3")
				panic("Unknown value for flag found")
			}
		}
	}

	if len(terminalArgs) >= 3 {
		if !utils.IsConfigOption(terminalArgs[2]) {
			root_directory = terminalArgs[2]
			fmt.Printf("path found in terminal %s\n", root_directory)
		}
	}

	for a := 0; a < len(terminalArgs); a++ {
		arg := terminalArgs[a]
		if utils.IsConfigOption(arg) {
			configTuple := strings.Split(arg, "=")
			option := configTuple[0]
			value := configTuple[1]

			switch option {
			case "i":
				newIndentation, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("Input %s wasn't able to be converted to a number to use as new number of indentation. Used default value of 2\n", value)
				} else {
					defaultIndentationSpaces = newIndentation
				}

			case "c":
				newCharacterLineLimit, err := strconv.Atoi(value)
				if err != nil {
					fmt.Printf("Input %s wasn't able to be converted to a number to use as new character limit. Used default value of 100\n", value)
				} else {
					defaultCharacterLimit = newCharacterLineLimit
				}
			case "n":
				isValidConventions := utils.IsValidNamingConvention(value)
				if !isValidConventions {
					fmt.Printf("%s isn't a valid or supported variable naming convention. See docs\n", value)
				} else {
					defaultVariableNamingConvention = value
				}
			case "r":
				willGenerateReport, err := strconv.ParseBool(value)
				if err != nil {
					fmt.Printf("%s isn't a valid boolean value\n", value)
				} else {
					shouldGenerateFailureReport = willGenerateReport
				}
			case "s":
				newSemiColonUsage, err := strconv.ParseBool(value)
				if err != nil {
					fmt.Printf("%s isn't a valid boolean value\n", value)
				} else {
					defaultSemiColonUsage = newSemiColonUsage
				}
			}
		}
	}
	var report models.ConsistencyReport
	checksRan := true

	switch check {
	case "all":
		report = checks.PerformAllChecks(root_directory, fileType, defaultVariableNamingConvention, defaultIndentationSpaces, defaultCharacterLimit, defaultSemiColonUsage)
	case "indent":
		report = checks.PerformIndentationChecks(root_directory, fileType, defaultIndentationSpaces)
	case "naming":
		report = checks.PerformVariableNamingChecks(root_directory, fileType, defaultVariableNamingConvention)
	case "char":
		report = checks.PerformCharacterCountChecks(root_directory, fileType, defaultCharacterLimit)
	case "semi":
		report = checks.PerformSemiColonChecks(root_directory, fileType, defaultSemiColonUsage)
	default:
		checksRan = false
		println("unexpected command please try again\n")
	}

	if checksRan {
		fmt.Printf("Number of issues found: %d\n", len(report.IssuesFound))
		if shouldGenerateFailureReport {
			reportCreated := utils.GenerateFailureReport(report, downloadsFolderPath)
			if !reportCreated {
				println("Error occurred while generating report")
			}
		}
	}
}
