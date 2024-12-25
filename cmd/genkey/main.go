// This is the main package for initializing and running the server.
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/andygeiss/cloud-native-utils/security"
)

func main() {
	key := security.GenerateKey()
	fmt.Printf("%s\n", hex.EncodeToString(key[:]))
}
