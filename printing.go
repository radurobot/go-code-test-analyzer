package main

import (
	"bytes"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func printCounts(fileCounts map[string]*FileAnalysis) {
	total := aggregateCounts(fileCounts)
	table := setupTable()

	for filename, analysis := range fileCounts {
		percentage := calculatePercentage(analysis.TestFunctions, analysis.Functions)
		colorizedFilename, colorizedFunctions, colorizedTestFunctions, colorizedPercentage := getColorizedStrings(filename, analysis.Functions, analysis.TestFunctions, percentage)
		row := []string{colorizedFilename, colorizedFunctions, colorizedTestFunctions, colorizedPercentage}
		table.Append(row)
	}

	appendTotalRow(table, total)
	table.Render()
	printPercentageOfTestFunctions(total.Functions, total.TestFunctions)
}

func setupTable() *tablewriter.Table {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Filename", "Functions", "Test Functions", "Percentage"})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
		tablewriter.Colors{tablewriter.Bold, tablewriter.FgWhiteColor},
	)
	table.SetColumnColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{},
		tablewriter.Colors{},
		tablewriter.Colors{},
	)
	table.SetBorder(true)
	return table
}

func calculatePercentage(numerator, denominator int) float64 {
	return float64(numerator) / float64(denominator) * 100
}

func getColorizedStrings(filename string, functions, testFunctions int, percentage float64) (string, string, string, string) {
	var colorizedFilename, colorizedFunctions, colorizedTestFunctions, colorizedPercentage string

	switch {
	case percentage < 20:
		colorizedFilename = color.RedString(filename)
		colorizedFunctions = color.RedString(fmt.Sprintf("%d", functions))
		colorizedTestFunctions = color.RedString(fmt.Sprintf("%d", testFunctions))
		colorizedPercentage = color.RedString(fmt.Sprintf("%.2f%%", percentage))
	case percentage < 60:
		colorizedFilename = color.YellowString(filename)
		colorizedFunctions = color.YellowString(fmt.Sprintf("%d", functions))
		colorizedTestFunctions = color.YellowString(fmt.Sprintf("%d", testFunctions))
		colorizedPercentage = color.YellowString(fmt.Sprintf("%.2f%%", percentage))
	default:
		colorizedFilename = color.GreenString(filename)
		colorizedFunctions = color.GreenString(fmt.Sprintf("%d", functions))
		colorizedTestFunctions = color.GreenString(fmt.Sprintf("%d", testFunctions))
		colorizedPercentage = color.GreenString(fmt.Sprintf("%.2f%%", percentage))
	}

	return colorizedFilename, colorizedFunctions, colorizedTestFunctions, colorizedPercentage
}

func appendTotalRow(table *tablewriter.Table, total *FileAnalysis) {
	table.Append([]string{"", "", "", ""})
	table.Append([]string{"Total", color.YellowString(fmt.Sprintf("%d", total.Functions)), color.HiYellowString(fmt.Sprintf("%d", total.TestFunctions))})
}

func printPercentageOfTestFunctions(functions, testFunctions int) {
	percentage := calculatePercentage(testFunctions, functions)
	var colorizedPercentage string
	switch {
	case percentage < 20:
		colorizedPercentage = color.RedString(fmt.Sprintf("%.2f%%", percentage))
	case percentage < 60:
		colorizedPercentage = color.YellowString(fmt.Sprintf("%.2f%%", percentage))
	default:
		colorizedPercentage = color.GreenString(fmt.Sprintf("%.2f%%", percentage))
	}
	fmt.Println("Percentage of test functions: ", colorizedPercentage)
}

func convertToMarkdown(fileCounts map[string]*FileAnalysis) string {
	var buf bytes.Buffer
	total := aggregateCounts(fileCounts)
	w := tabwriter.NewWriter(&buf, 0, 0, 1, ' ', tabwriter.TabIndent)

	fmt.Fprintln(w, "| Filename | Functions | Test Functions |")
	fmt.Fprintln(w, "| --- | --- | --- |")

	for filename, analysis := range fileCounts {
		fmt.Fprintf(w, "| %s | %d | %d |\n", filename, analysis.Functions, analysis.TestFunctions)
	}
	fmt.Fprintf(w, "| **Total** | **%d** | **%d** |\n", total.Functions, total.TestFunctions)

	fmt.Fprintln(w, "Percentage of test functions: ", float64(total.TestFunctions)/float64(total.Functions)*100)

	w.Flush()

	return buf.String()
}
