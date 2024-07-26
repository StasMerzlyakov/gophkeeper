package main

import (
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func printVersion() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {

	printVersion()

	_, err := config.LoadClientConf()
	if err != nil {
		panic(err)
	}

}
