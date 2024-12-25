// This program generates a new 256-bit key and prints it as a hex string.
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/andygeiss/cloud-native-utils/security"
)

func main() {
	// Generate a new 256-bit key.
	key := security.GenerateKey()
	// Print the key as a hex string.
	fmt.Printf("%s\n", hex.EncodeToString(key[:]))
}
