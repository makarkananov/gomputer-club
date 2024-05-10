package main

import (
	"flag"
	"gomputerClub/internal/adapters/input"
	"gomputerClub/internal/adapters/output"
	"gomputerClub/internal/core/service"
	"log"
)

func main() {
	flag.Parse()
	filename := flag.Arg(0)
	if filename == "" {
		log.Fatal("filename is required")
	}

	ccm := service.NewComputerClubManager(output.NewConsole())

	fic := input.NewFileController(filename, ccm)
	incomeEvents, err := fic.Read()
	if err != nil {
		log.Fatalf("error reading input file: %v", err)
	}

	fic.Run(incomeEvents)
}
