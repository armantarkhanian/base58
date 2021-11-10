# Example
```golang
package main

import (
	"fmt"

	"github.com/armantarkhanian/base58"
)

func main() {
	ed, err := base58.New(base58.Flickr, 100000000000000)
	if err != nil {
		panic(err)
	}

	var (
		userID       int64  = 42
		userIDString string = ed.Encode(userID)
		newID        int64  = ed.Decode(userIDString)
	)

	fmt.Println(userID)       // 42
	fmt.Println(userIDString) // MhQayQkd
	fmt.Println(newID)        // 42
}
```
