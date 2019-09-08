package main

import (
	"context"
	"flag"
	"fmt"
	runtimevar "github.com/aaronland/gocloud-runtimevar-string"
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
