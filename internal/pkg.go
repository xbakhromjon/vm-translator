package main

import "os"

func GetArg(name string) string {
	args := os.Args
	next := false
	for _, a := range args {
		if next {
			return a
		}
		if a == name {
			next = true
		}
	}
	return ""
}
