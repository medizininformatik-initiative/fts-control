package utils

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Dividing Lines

// Large

func DivdlnL() {
	FprintDivdlnL(os.Stdout)
}

// FprintDivdlnL writes a large dividing line to the given writer.
func FprintDivdlnL(w io.Writer) {
	fmt.Fprintln(w, "==================================================================")
}

// Small

func DivdlnS() {
	FprintDivdlnS(os.Stdout)
}

// FprintDivdlnS writes a small dividing line to the given writer.
func FprintDivdlnS(w io.Writer) {
	fmt.Fprintln(w, "------------------------------------------------------------------")
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

// PrintProcessStatus prints formatted status information for a Process to stdout.
func PrintProcessStatus(p Process) {
	FprintProcessStatus(os.Stdout, p)
}

// FprintProcessStatus writes formatted status information for a Process to the given writer.
func FprintProcessStatus(w io.Writer, p Process) {
	fmt.Fprintf(w, "ProcessID: %s\n", p.ProcessID)
	fmt.Fprintf(w, "Phase: %s\n", p.Phase)
	fmt.Fprintf(w, "CreatedAt: %s\n", p.CreatedAt.ToTime().Local().Format(time.RFC1123))
	if p.FinishedAt != nil {
		fmt.Fprintf(w, "FinishedAt: %s\n", p.FinishedAt.ToTime().Local().Format(time.RFC1123))
	} else {
		fmt.Fprintln(w, "FinishedAt: The process is still running")
	}
	fmt.Fprintf(w, "TotalPatients: %d\n", p.TotalPatients)
	fmt.Fprintf(w, "TotalBundles: %d\n", p.TotalBundles)
	fmt.Fprintf(w, "DeidentifiedBundles: %d\n", p.DeidentifiedBundles)
	fmt.Fprintf(w, "SentBundles: %d\n", p.SentBundles)
	fmt.Fprintf(w, "SkippedBundles: %d\n", p.SkippedBundles)
}
