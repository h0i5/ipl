package main

// this file is to store the ascii arts for the team names and numbers
// generated from https://patorjk.com/software/taag

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const teamMI string = `
 /$$      /$$ /$$$$$$
| $$$    /$$$|_  $$_/
| $$$$  /$$$$  | $$
| $$ $$/$$ $$  | $$
| $$  $$$| $$  | $$
| $$\  $ | $$  | $$
| $$ \/  | $$ /$$$$$$
|__/     |__/|______/
`

const teamRCB string = `
 /$$$$$$$   /$$$$$$  /$$$$$$$ 
| $$__  $$ /$$__  $$| $$__  $$
| $$  \ $$| $$  \__/| $$  \ $$
| $$$$$$$/| $$      | $$$$$$$ 
| $$__  $$| $$      | $$__  $$
| $$  \ $$| $$    $$| $$  \ $$
| $$  | $$|  $$$$$$/| $$$$$$$/
|__/  |__/ \______/ |_______/ 
`

const teamCSK string = `
  /$$$$$$   /$$$$$$  /$$   /$$
 /$$__  $$ /$$__  $$| $$  /$$/
| $$  \__/| $$  \__/| $$ /$$/ 
| $$      |  $$$$$$ | $$$$$/  
| $$       \____  $$| $$  $$  
| $$    $$ /$$  \ $$| $$\  $$ 
|  $$$$$$/|  $$$$$$/| $$ \  $$
 \______/  \______/ |__/  \__/
`

const teamGT string = `
  /$$$$$$  /$$$$$$$$
 /$$__  $$|__  $$__/
| $$  \__/   | $$   
| $$ /$$$$   | $$   
| $$|_  $$   | $$   
| $$  \ $$   | $$   
|  $$$$$$/   | $$   
 \______/    |__/   
`
const teamKKR string = `
 /$$   /$$ /$$   /$$ /$$$$$$$ 
| $$  /$$/| $$  /$$/| $$__  $$
| $$ /$$/ | $$ /$$/ | $$  \ $$
| $$$$$/  | $$$$$/  | $$$$$$$/
| $$  $$  | $$  $$  | $$__  $$
| $$\  $$ | $$\  $$ | $$  \ $$
| $$ \  $$| $$ \  $$| $$  | $$
|__/  \__/|__/  \__/|__/  |__/
`
const teamLSG string = `
 /$$        /$$$$$$   /$$$$$$ 
| $$       /$$__  $$ /$$__  $$
| $$      | $$  \__/| $$  \__/
| $$      |  $$$$$$ | $$ /$$$$
| $$       \____  $$| $$|_  $$
| $$       /$$  \ $$| $$  \ $$
| $$$$$$$$|  $$$$$$/|  $$$$$$/
|________/ \______/  \______/ 
`
const teamPK string = `
 /$$$$$$$  /$$   /$$
| $$__  $$| $$  /$$/
| $$  \ $$| $$ /$$/ 
| $$$$$$$/| $$$$$/  
| $$____/ | $$  $$  
| $$      | $$\  $$ 
| $$      | $$ \  $$
|__/      |__/  \__/
`
const teamRR string = `
 /$$$$$$$  /$$$$$$$ 
| $$__  $$| $$__  $$
| $$  \ $$| $$  \ $$
| $$$$$$$/| $$$$$$$/
| $$__  $$| $$__  $$
| $$  \ $$| $$  \ $$
| $$  | $$| $$  | $$
|__/  |__/|__/  |__/
`
const teamDC string = `
 /$$$$$$$   /$$$$$$ 
| $$__  $$ /$$__  $$
| $$  \ $$| $$  \__/
| $$  | $$| $$      
| $$  | $$| $$      
| $$  | $$| $$    $$
| $$$$$$$/|  $$$$$$/
|_______/  \______/ 
`
const teamSRH string = `
  /$$$$$$  /$$$$$$$  /$$   /$$
 /$$__  $$| $$__  $$| $$  | $$
| $$  \__/| $$  \ $$| $$  | $$
|  $$$$$$ | $$$$$$$/| $$$$$$$$
 \____  $$| $$__  $$| $$__  $$
 /$$  \ $$| $$  \ $$| $$  | $$
|  $$$$$$/| $$  | $$| $$  | $$
 \______/ |__/  |__/|__/  |__/
`

const numberOne string = `
▗ 
▜ 
▟▖
`

const numberTwo string = `
▄▖
▄▌
▙▖
`

const numberThree string = `
▄▖
▄▌
▄▌
`
const numberFour string = `
▖▖
▙▌
 ▌
`
const numberFive string = `
▄▖
▙▖
▄▌
`
const numberSix string = `
▄▖
▙▖
▙▌
`
const numberSeven string = `
▄▖
 ▌
 ▌
`
const numberEight string = `
▄▖
▙▌
▙▌
`
const numberNine string = `
▄▖
▙▌
▄▌
`
const numberZero string = `
▄▖
▛▌
█▌
`

const dash string = `

▄▖

`

var teamASCII = map[string]string{
	"MI":   teamMI,
	"RCB":  teamRCB,
	"CSK":  teamCSK,
	"GT":   teamGT,
	"KKR":  teamKKR,
	"LSG":  teamLSG,
	"PK":   teamPK,
	"PBKS": teamPK,
	"RR":   teamRR,
	"DC":   teamDC,
	"SRH":  teamSRH,
}

var digitASCII = map[rune]string{
	'0': numberZero,
	'1': numberOne,
	'2': numberTwo,
	'3': numberThree,
	'4': numberFour,
	'5': numberFive,
	'6': numberSix,
	'7': numberSeven,
	'8': numberEight,
	'9': numberNine,
	'/': dash,
}

func teamArt(code string) string {
	art, ok := teamASCII[strings.ToUpper(code)]
	if !ok {
		return ""
	}
	return strings.TrimLeft(art, "\n")
}

func scoreArt(score string) string {
	var parts []string
	for _, ch := range score {
		art, ok := digitASCII[ch]
		if !ok {
			continue
		}
		parts = append(parts, strings.TrimPrefix(art, "\n"))
	}
	if len(parts) == 0 {
		return ""
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, parts...)
}
