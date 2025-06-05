package utils

import "fmt"

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
