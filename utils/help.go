package utils

import "fmt"

func PrintHelp() {
	fmt.Println("Usage:")
	fmt.Println("  webp2p [--no-stun] <localPort> me")
	fmt.Println("  webp2p [--no-stun] <localPort> connect")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --no-stun   Skip public IP discovery via STUN")
}
