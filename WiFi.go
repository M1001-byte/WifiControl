package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getNetworks(iface string) []string {
	var outpout []string
	var essid string
	var security string
	var signal string
	command := "iwctl"
	// scan
	cmdScan := exec.Command(command, "station", iface, "scan")
	_, err := cmdScan.Output()
	if err != nil {
		println(err.Error())
	}
	args := []string{"station", iface, "get-networks"}

	cmd := exec.Command(command, args...)

	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		fmt.Printf("%s: %v\n", red("ERROR"), err)
		return nil
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := string(scanner.Text())
		if strings.Contains(text, "*") {
			textSlice := strings.Split(text, " ")
			// como elimina los espacios en blanco,tambien los hace cuano la red tiene uno xd
			textSlice = removeBlankStrings(textSlice)

			securityI := indexOf("open", textSlice)
			if securityI == -1 {
				securityI = indexOf("psk", textSlice)
			}

			if len(textSlice) >= 6 {
				// current wifi network
				essid, security, signal = strings.Join(textSlice[3:securityI], " "), textSlice[securityI], textSlice[securityI+1]
			} else {
				essid, security, signal = strings.Join(textSlice[0:securityI], " "), textSlice[securityI], textSlice[securityI+1]
			}
			compactstr := fmt.Sprintf("%s:%s:%s", essid, security, signal)
			outpout = append(outpout, compactstr)

		}
	}

	return outpout
}

func connectWifi(iface string, essid string, password string) {
	cmd := exec.Command("iwctl", "--passphrase", password, "station", iface, "connect", essid)
	outpout, err := cmd.Output()
	if err != nil {
		str := string(outpout)
		if strings.Contains(str, "failed") {
			fmt.Printf("%s: password is incorrect.\n", red("ERROR"))
		} else {
			fmt.Printf("%s: %s.\n", red("ERROR"), outpout)
		}
	} else {
		fmt.Printf("%s: %s", green("Connected succesfull to"), yellow(essid))
		os.Exit(0)
	}

}
