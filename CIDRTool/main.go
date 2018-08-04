package main

import (
	"fmt"

	"os"

	"github.com/zsoumya/smutils/net"
	"github.com/zsoumya/smutils/str"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode. Includes descriptive verbiage in the output.").Default("false").Short('v').Bool()

	analyzeCommand = kingpin.Command("analyze", "analyze a CIDR block")
	analyzeCIDRVal = analyzeCommand.Arg("cidrBlock", "CIDR block to anlyze").Required().String()

	containsCommand = kingpin.Command("contains", "Checks if a CIDR block contains another CIDR block or an IP address")
	containsParent  = containsCommand.Arg("cidrBlock", "Parent CIDR block").Required().String()
	containsValue   = containsCommand.Arg("value", "CIDR block or IP address to check if contained in parent").Required().String()

	spanCommand = kingpin.Command("span", "Expresses a span of two IP addresses as a series of consecutive CIDR blocks")
	spanStartIP = spanCommand.Arg("startIP", "Start IP address of the span").Required().String()
	spanEndIP   = spanCommand.Arg("endIP", "End IP address of the span").Required().String()

	splitCommand = kingpin.Command("split", "Splits a CIDR block into smaller equal sized CIDR blocks")
	splitParent  = splitCommand.Arg("cidrBlock", "CIDR block to be split into smaller blocks").Required().String()
	splitCount   = splitCommand.Arg("splitCount", "Number of smaller blocks the parent CIDR block to be split into; has to be power of 2").Required().Int()
)

func main() {
	run()
}

func run() {
	command := kingpin.Parse()
	switch command {
	case analyzeCommand.FullCommand():
		analyze(*analyzeCIDRVal, *verbose)
	case containsCommand.FullCommand():
		contains(*containsParent, *containsValue, *verbose)
	case spanCommand.FullCommand():
		span(*spanStartIP, *spanEndIP, *verbose)
	case splitCommand.FullCommand():
		split(*splitParent, uint(*splitCount), *verbose)
	default:
		printErrorln("Cannot interpret command:", command)
	}
}

func analyze(cidrVal string, verbose bool) {
	cidr, err := net.ParseCIDRv4Str(cidrVal)
	if err != nil {
		printErrorln("Invalid CIDR block:", cidrVal)
		return
	}

	padChar := "."

	if verbose {
		fmt.Println(str.PadRight("CIDR", padChar, 25), cidr)
		fmt.Println(str.PadRight("Netmask", padChar, 25), cidr.NetMask())
		fmt.Println(str.PadRight("Wildcard Bits", padChar, 25), cidr.Wildcard())
		fmt.Println(str.PadRight("Network IP", padChar, 25), cidr.NetworkIP())
		fmt.Println(str.PadRight("Broadcast IP", padChar, 25), cidr.BroadcastIP())
		fmt.Println(str.PadRight("First Usable IP", padChar, 25), cidr.FirstIP())
		fmt.Println(str.PadRight("Last Usable IP", padChar, 25), cidr.LastIP())
		fmt.Println(str.PadRight("IP Address Count", padChar, 25), cidr.Count())
		fmt.Println(str.PadRight("Usable IP Address Count", padChar, 25), cidr.UsableCount())
	} else {
		fmt.Println(cidr)
		fmt.Println(cidr.NetMask())
		fmt.Println(cidr.Wildcard())
		fmt.Println(cidr.NetworkIP())
		fmt.Println(cidr.BroadcastIP())
		fmt.Println(cidr.FirstIP())
		fmt.Println(cidr.LastIP())
		fmt.Println(cidr.Count())
		fmt.Println(cidr.UsableCount())
	}
}

func contains(cidrVal string, value string, verbose bool) {
	cidr, err := net.ParseCIDRv4Str(cidrVal)
	if err != nil {
		printErrorln("Invalid parent CIDR block:", cidrVal)
		return
	}

	contains := false
	mode := ""

	childCIDR, _ := net.ParseCIDRv4Str(value)
	if childCIDR != nil {
		mode = "CIDR block"
		contains = cidr.ContainsCIDRv4(childCIDR)
	} else {
		ip, _ := net.ParseIPv4Str(value)
		if ip != nil {
			mode = "IP address"
			contains = cidr.ContainsIPv4(ip)
		}
	}

	if mode == "" {
		printErrorln("Not an IP address or a CIDR block:", value)
		return
	}

	if verbose {
		s := "contains"
		if !contains {
			s = "does not contain"
		}

		fmt.Printf("CIDR block %v %s %s %v\n", cidr, s, mode, value)
	} else {
		fmt.Println(contains)
	}
}

func span(startIPVal string, endIPVal string, verbose bool) {
	startIP, _ := net.ParseIPv4Str(startIPVal)
	endIP, _ := net.ParseIPv4Str(endIPVal)

	if startIP == nil {
		printErrorln("Not a valid IP address:", startIPVal)
		return
	}

	if startIP == nil {
		printErrorln("Not a valid IP address:", endIPVal)
		return
	}

	cidrs := startIP.CIDRSpan(endIP)

	for i, cidr := range cidrs {
		if verbose {
			fmt.Printf("%03d: %s\n", i, cidr.Desc())
		} else {
			fmt.Println(cidr)
		}
	}
}

func split(cidrVal string, count uint, verbose bool) {
	cidr, err := net.ParseCIDRv4Str(cidrVal)
	if err != nil {
		printErrorln("Invalid CIDR block:", cidrVal)
		return
	}

	cidrs, err := cidr.Split(count)
	if err != nil {
		printErrorln(err)
		return
	}

	if verbose {
		fmt.Println("Parent CIDR:")
		fmt.Println("    ", cidr.Desc())

		fmt.Println()
		fmt.Println("Child CIDRs:")
	}

	for i, eachCIDR := range cidrs {
		if verbose {
			fmt.Printf("%03d: %s\n", i, eachCIDR.Desc())
		} else {
			fmt.Println(eachCIDR)
		}
	}
}

func printErrorln(msg ...interface{}) {
	fmt.Fprint(os.Stderr, "Error: ")
	fmt.Fprintln(os.Stderr, msg...)
}
