package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
)

// COLOR
var yellow = color.New(color.FgYellow).SprintFunc()
var green = color.New(color.FgGreen).SprintFunc()
var red = color.New(color.FgRed).SprintFunc()

func parseArgs() (string, error) {
	args := os.Args
	if len(args) > 1 {
		iface := os.Args[1]
		return iface, nil
	} else {
		return "", fmt.Errorf("iface not found")
	}
}

func removeBlankStrings(slice []string) []string {
	var result []string

	for _, str := range slice {
		if str != "" {
			result = append(result, str)
		}
	}

	return result
}
func printSlice(slice []string, cursorPostion int) {
	fmt.Print("\033[H\033[2J") // clear
	fmt.Printf("\033[38;5;240m-----------------------------------------------------------------------------------\033[0m\n")
	fmt.Printf("\tUse w/a/s/d or arrows to select wifi, and (R) to reload, (Q) to quit.\n")
	fmt.Printf("\033[38;5;240m-----------------------------------------------------------------------------------\033[0m\n")
	fmt.Printf("\n %-3s %-32s %-15s %-5s \n", "", "ESSID", "SECURITY", "SIGNAL")

	for index, v := range slice {
		splitResult := strings.Split(v, ":")

		essid, signal, encrypt := splitResult[0], splitResult[1], splitResult[2]
		if index == cursorPostion {
			fmt.Printf("\n %-3s %-32s %-15s %-5s \n", ">>>", essid, signal, encrypt)

		} else {
			fmt.Printf("\n %-3s %-32s %-15s %-5s \n", "", essid, signal, encrypt)
		}
	}
	fmt.Printf("\033[38;5;240m-----------------------------------------------------------------------------------\033[0m\n")
}

func selectOpt(opt []string, iface string) {
	var selectNetwork = false
	currentPosition := len(opt)
	printSlice(opt, currentPosition)
	for {
		char, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Printf("%s", err.Error())
			break
		}
		if key == keyboard.KeyArrowUp || char == 'w' {
			currentPosition--
			if currentPosition == -1 {
				currentPosition = len(opt)
			}
		}
		if key == keyboard.KeyArrowDown || char == 's' {
			currentPosition++
			if currentPosition > len(opt) {
				currentPosition = len(opt) - 1
			}
		}
		if key == keyboard.KeyEnter {
			if currentPosition > len(opt)-1 || currentPosition == -1 {
				continue
			}
			var password string
			essidstr := opt[currentPosition]
			essid := strings.Split(essidstr, ":")[0]

			for {
				fmt.Printf("Enter password for %s:", essid)
				_, err := fmt.Scan(&password)
				if err != nil {
					println(err.Error())
				}
				if len(password) < 8 {
					fmt.Printf("\n%s: Password lenght is < 8. (%d)\n", red("ERROR"), len(password))
					continue
				} else {
					break
				}

			}

			fmt.Printf(" %s\n", password)
			connectWifi(iface, essid, password)
			selectNetwork = true
		}
		if char == 'q' {
			os.Exit(0)
		}
		if char == 'r' {
			selectOpt(getNetworks(iface), iface)
		}
		if !selectNetwork {
			printSlice(opt, currentPosition)
		}
	}
}

func indexOf(string string, slice []string) int {
	for k, v := range slice {
		if string == v {
			return k
		}
	}
	return -1 //not found.
}

func main() {
	iface, err := parseArgs()
	if err != nil {
		fmt.Printf("%s: %s {nameIface}\n", yellow("Usage"), os.Args[0])
		fmt.Printf("%s: %s\n", red("ERROR"), err.Error())
		os.Exit(1)
	}
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		keyboard.Close()
	}()

	wifiList := getNetworks(iface)
	selectOpt(wifiList, iface)

}
