package main

import "github.com/amir79esmaeili/go-tel-money/cmd"

func main() {
	if err := cmd.New().Execute(); err != nil {
		panic(err)
	}
}
