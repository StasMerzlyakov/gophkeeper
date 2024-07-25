package main

import "github.com/StasMerzlyakov/gophkeeper/internal/config"

func main() {
	_, err := config.LoadClientConf()
	if err != nil {
		panic(err)
	}

}
