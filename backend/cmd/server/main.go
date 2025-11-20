package main

import "github.com/joho/godotenv"

func main() {
	// redundant? could just use envconfig?
	if err := godotenv.Load(); err != nil {
		// fatal...
	}

	// hexagonal...

}
