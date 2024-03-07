package utils

import (
	color "github.com/fatih/color"
)

func WarningSprintF(message string) string {
	return color.New(color.FgYellow).SprintFunc()(message)
}

func ErrorSprintF(message string) string {
	return color.New(color.FgRed).SprintFunc()(message)
}

func SuccessSprintF(message string) string {
	return color.New(color.FgGreen).SprintFunc()(message)
}

func InfoSprintF(message string) string {
	return color.New(color.FgBlue).SprintFunc()(message)
}

func DefaultSprintF(message string) string {
	return color.New(color.FgWhite).SprintFunc()(message)
}

func WarningPrintLn(message string) {
	color.New(color.FgYellow).Println(message)
}

func ErrorPrintLn(message string) {
	color.New(color.FgRed).Println(message)
}

func SuccessPrintLn(message string) {
	color.New(color.FgGreen).Println(message)
}

func InfoPrintLn(message string) {
	color.New(color.FgBlue).Println(message)
}

func DefaultPrintLn(message string) {
	color.New(color.FgWhite).Println(message)
}
