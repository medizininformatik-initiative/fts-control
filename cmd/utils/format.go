package utils

import (
	"fmt"
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
