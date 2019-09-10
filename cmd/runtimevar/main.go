package main

import (
	"context"
	"flag"
	"fmt"
	runtimevar "github.com/aaronland/gocloud-runtimevar-string"
	_ "gocloud.dev/runtimevar/blobvar"
	_ "gocloud.dev/runtimevar/constantvar"
	_ "gocloud.dev/runtimevar/filevar"
	"log"
)

func main() {

	url := flag.String("url", "", "...")

	flag.Parse()

	ctx := context.Background()

	s, err := runtimevar.OpenString(ctx, *url)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s)
}
