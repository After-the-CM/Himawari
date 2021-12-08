package scanner

import (
	"fmt"

	"Himawari/models/entity"
)

func MarkDown() string {
	var md string
	md = fmt.Sprintf("# Vulnerabilities report by Himawari🌻\n\n")
	for _, vuln := range entity.Vulnmap {
		if len(vuln.Issues) != 0 {
			md += vuln2md(vuln)
		}
	}
	return md
}

func vuln2md(vuln *entity.Vuln) string {
	var md string
	md += fmt.Sprintf("## %s %s (Severity: %s)\n\n", vuln.CWE, vuln.Name, vuln.Severity)
	md += fmt.Sprintf("### 概要\n\n")
	md += fmt.Sprintf("%s\n\n", vuln.Description)

	md += fmt.Sprintf("### 必須対策\n\n")
	md += fmt.Sprintf("%s\n\n", vuln.Mandatory)

	md += fmt.Sprintf("### 保険的対策\n\n")
	md += fmt.Sprintf("%s\n\n", vuln.Insurance)

	for _, issue := range vuln.Issues {
		md += issue2md(issue)
	}
	return md
}

func issue2md(issue entity.Issue) string {
	var md string
	md += fmt.Sprintf("#### <%s>\n\n", issue.URL)
	md += fmt.Sprintf("Name|Value\n-|-\n")

	if issue.Parameter != "" {
		md += fmt.Sprintf("Parameter|`%s`\n", issue.Parameter)
	} else {
		md += fmt.Sprintf("Parameter|\n")
	}

	if issue.Payload != "" {
		md += fmt.Sprintf("Payload|`%s`\n", issue.Payload)
	} else {
		md += fmt.Sprintf("Payload|\n")
	}

	md += fmt.Sprintf("Evidence|`%s`\n\n", issue.Evidence)
	md += fmt.Sprintf("Request\n\n")
	md += fmt.Sprintf("```http\n%s\n```\n\n", issue.Request)
	md += fmt.Sprintf("Response\n\n")
	md += fmt.Sprintf("```http\n%s\n```\n\n", issue.Response)
	return md
}
