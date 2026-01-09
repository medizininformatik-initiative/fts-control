package utils

import (
	"fmt"
	"log/slog"
	"time"
)

// Dividing Lines

// Large

func DivdlnL() {
	fmt.Println("==================================================================")
}

// Small

func DivdlnS() {
	fmt.Println("------------------------------------------------------------------")
}

// Logo

func GenerateLogo() string {
	const Reset = "\033[0m"
	const Cyan = "\033[96m"
	logo := `
███████╗████████╗███████╗												
██╔════╝╚══██╔══╝██╔════╝` + Cyan + `           ██╗   ██╗` + Reset + `
█████╗     ██║   ███████╗` + Cyan + ` ██████╗████████╗██║` + Reset + `
██╔══╝     ██║   ╚════██║` + Cyan + `██╔════╝╚══██╔══╝██║` + Reset + `
██║        ██║   ███████║` + Cyan + `╚██████╗   ██║   █████╗` + Reset + `
╚═╝        ╚═╝   ╚══════╝` + Cyan + ` ╚═════╝   ╚═╝   ╚════╝` + Reset + `
       
`

	return logo
}

func GenerateLogoV2() string {
	const Reset = "\033[0m"
	const Cyan = "\033[96m"
	logo := `
███████` + Cyan + `╗` + Reset + `████████` + Cyan + `╗` + Reset + `███████` + Cyan + `╗` + Reset + `
██` + Cyan + `╔════╝╚══` + Reset + `██` + Cyan + `╔══╝` + Reset + `██` + Cyan + `╔════╝` + Reset + `` + Cyan + `           ██` + Reset + `╗` + Cyan + `   ██` + Reset + `╗
█████` + Cyan + `╗` + Reset + `     ██` + Cyan + `║` + Reset + `   ███████` + Cyan + `╗` + Reset + `` + Cyan + ` ██████` + Reset + `╗` + Cyan + `████████` + Reset + `╗` + Cyan + `██` + Reset + `║
██` + Cyan + `╔══╝` + Reset + `     ██` + Cyan + `║   ╚════` + Reset + `██` + Reset + `` + Cyan + `║` + Cyan + `██` + Reset + `╔════╝╚══` + Cyan + `██` + Reset + `╔══╝` + Cyan + `██` + Reset + `║
██` + Cyan + `║` + Reset + `        ██` + Cyan + `║` + Reset + `   ███████` + Reset + `` + Cyan + `║` + Cyan + `` + Reset + `╚` + Cyan + `██████` + Reset + `╗` + Cyan + `   ██` + Reset + `║` + Cyan + `   █████` + Reset + `╗
` + Cyan + `╚═╝        ╚═╝   ╚══════╝` + Cyan + `` + Reset + ` ╚═════╝   ╚═╝   ╚════╝
`

	return logo
}

// PrintProcessStatus prints formatted status information for a Process.
func PrintProcessStatus(p Process) {
	fmt.Printf("ProcessID: %s\n", p.ProcessID)
	fmt.Printf("Phase: %s\n", p.Phase)
	fmt.Printf("CreatedAt: %s\n", p.CreatedAt.ToTime().Local().Format(time.RFC1123))
	if p.FinishedAt != nil {
		fmt.Printf("FinishedAt: %s\n", p.FinishedAt.ToTime().Local().Format(time.RFC1123))
	} else {
		fmt.Println("FinishedAt: The process is still running")
	}
	fmt.Printf("TotalPatients: %d\n", p.TotalPatients)
	fmt.Printf("TotalBundles: %d\n", p.TotalBundles)
	fmt.Printf("DeidentifiedBundles: %d\n", p.DeidentifiedBundles)
	fmt.Printf("SentBundles: %d\n", p.SentBundles)
	fmt.Printf("SkippedBundles: %d\n", p.SkippedBundles)
}

// LogHTTPError logs an HTTP error appropriately and returns true if an error was handled.
func LogHTTPError(err error, context string) bool {
	if err == nil {
		return false
	}

	if httpErr, ok := IsHTTPError(err); ok {
		slog.Error("API request failed", "context", context, "status", httpErr.Status, "details", httpErr.Body)
	} else {
		slog.Error("Request failed", "context", context, "error", err)
	}
	return true
}

// PrintHTTPError prints an HTTP error in a user-friendly format and returns true if an error was handled.
// Use this when you want to display error details to the user (not just log them).
func PrintHTTPError(err error) bool {
	if err == nil {
		return false
	}

	DivdlnL()
	if httpErr, ok := IsHTTPError(err); ok {
		fmt.Printf("Request failed! Response Status: %s\n", httpErr.Status)
		if httpErr.Body != "" {
			fmt.Printf("Server response body:\n%s\n", httpErr.Body)
		}
	} else {
		slog.Error("Request failed", "error", err)
	}
	DivdlnL()
	return true
}
