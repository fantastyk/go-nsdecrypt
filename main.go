package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// netscaler uses a ROT cypher with a -1 rotation to encode the service name
// This may be dodgy. Need more sample cookies
/*func rot_d(r rune) rune {
	if r >= 'a' && r <= 'z' {
		// Rotate lowercase letters -1 places.
		if r >= 'm' {
			return r - 1
		} else {
			return r - 1
		}
	} else if r >= 'A' && r <= 'Z' {
		// Rotate uppercase letters -1 places.
		if r >= 'M' {
			return r - 1
		} else {
			return r - 1
		}
	}
	// Do nothing.
	return r
}*/

func parseCookie(cookie string) [][]string {
	//Regex to split the cookie into 3 groups. Service name, Backend IP, Server port
	re := regexp.MustCompile(`NSC_([a-zA-Z0-9\\\-\\\_\\\.]*)=[0-9a-f]{8}([0-9a-f]{8}).*([0-9a-f]{4})$`)
	cookiearray := re.FindAllStringSubmatch(cookie, 1)
	return cookiearray
}

func decryptServerIP(serverip int) string {

	ipkey := 0x03081e11 //int
	decodedip := serverip ^ ipkey
	// need to conver it back to int64
	ip := InttoIP4(int64(decodedip))
	return ip
}

//https://www.socketloop.com/tutorials/golang-convert-decimal-number-integer-to-ipv4-address
func InttoIP4(ipInt int64) string {

	// need to do two bit shifting and “0xff” masking
	b0 := strconv.FormatInt((ipInt>>24)&0xff, 10)
	b1 := strconv.FormatInt((ipInt>>16)&0xff, 10)
	b2 := strconv.FormatInt((ipInt>>8)&0xff, 10)
	b3 := strconv.FormatInt((ipInt & 0xff), 10)
	return b0 + "." + b1 + "." + b2 + "." + b3
}

func decryptServerPort(serverport int) string {
	portkey := 0x3630 //int-base16
	decryptedport := serverport ^ portkey
	return strconv.Itoa(decryptedport)
}

func main() {
	//cookiedough := "NSC_Qspe-xxx.bwjwb.dp.vl-IUUQ=ffffffff50effd8445525d5f4f58455e445a4a423660" //control
	cookiedough := os.Args[1]
	parsedcookie := parseCookie(cookiedough)

	//Need to convert from int64 to int for XOR'ing
	encserverip64, _ := strconv.ParseInt(parsedcookie[0][2], 16, 32)
	encserverip := int(encserverip64)
	encserverport64, _ := strconv.ParseInt(parsedcookie[0][3], 16, 32)
	encserverport := int(encserverport64)

	//Thanks Vince for the optimization
	rotMinus1 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return 'A' + (r-'A'+25)%26
		case r >= 'a' && r <= 'z':
			return 'a' + (r-'a'+25)%26
		}
		return r
	}

	fmt.Println("vServer Name:", strings.Map(rotMinus1, parsedcookie[0][1]))
	fmt.Println("VServer IP:", decryptServerIP(encserverip))
	fmt.Println("VServer Port:", decryptServerPort(encserverport))

}
