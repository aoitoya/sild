package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/toyaAoi/sild/codegen"
	"github.com/toyaAoi/sild/parser"
	"github.com/toyaAoi/sild/scanner"
)

func printError(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func main() {
	var outFileName string
	var isDebug bool

	flag.StringVar(&outFileName, "out", "", "Output file name")
	flag.StringVar(&outFileName, "o", "", "Output file name")
	flag.BoolVar(&isDebug, "debug", false, "Enable debug mode")
	flag.Parse()

	inputFile := os.Args[len(os.Args) - 1]

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Error: No input file specified\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [-o output_file] <input_file>\n", os.Args[0])
		os.Exit(1)
	}

	if isDebug {
		fmt.Fprintf(os.Stderr, "Debug: Input file: '%s'\n", inputFile)
		fmt.Fprintf(os.Stderr, "Debug: Output file: '%s'\n", outFileName)
	}
    
    if inputFile == "" {
        fmt.Fprintf(os.Stderr, "Error: No input file specified\n")
        fmt.Fprintf(os.Stderr, "Usage: %s [-o output_file] <input_file>\n", os.Args[0])
        os.Exit(1)
    }

	if isDebug {
		printError("Debug: Reading input file: %s\n", inputFile)
	}
    file, err := os.ReadFile(inputFile)
    if err != nil {
        printError("Error reading file: %v\n", err)
        os.Exit(1)
    }

    input := string(file)

    scanner := scanner.New(strings.NewReader(input))
    parser := parser.New(scanner)
    gen := codegen.New()

    program := parser.ParseProgram()
    output := gen.Generate(program)

	if isDebug {
		fmt.Fprintf(os.Stderr, "Debug: Checking if we should write to file...\n")
	}
    if outFileName != "" {
		if isDebug {
			fmt.Fprintf(os.Stderr, "Debug: Output file specified: %s\n", outFileName)
			fmt.Fprintf(os.Stderr, "Debug: Output file name before processing: %s\n", outFileName)
		}
        
        if !strings.HasSuffix(outFileName, ".go") {
            outFileName += ".go"
        }
        
        absPath, err := filepath.Abs(outFileName)
        if err != nil {
            printError("Error getting absolute path: %v\n", err)
            os.Exit(1)
        }
        
		if isDebug {
			fmt.Fprintf(os.Stderr, "Debug: Writing %d bytes to: %s\n", len(output), absPath)
		}
        
        dir := filepath.Dir(absPath)
        if _, err := os.Stat(dir); os.IsNotExist(err) {
            printError("Error: Directory does not exist: %s\n", dir)
            os.Exit(1)
        }
        
		if isDebug {
			fmt.Fprintf(os.Stderr, "Debug: Attempting to create file at: %s\n", absPath)
		}

        file, err := os.Create(absPath)
        if err != nil {
            printError("Error creating file: %v\n", err)
            dir := filepath.Dir(absPath)
            if stat, err := os.Stat(dir); err == nil {
                printError("Directory info: %+v, Permissions: %v\n", stat, stat.Mode().String())
            } else {
                printError("Could not stat directory: %v\n", err)
            }
            os.Exit(1)
        }
        
        n, err := file.WriteString(output)
        if err != nil {
            file.Close()
            printError("Error writing to file: %v\n", err)
            os.Exit(1)
        }
        
        if err := file.Close(); err != nil {
            printError("Error closing file: %v\n", err)
            os.Exit(1)
        }
 
		if isDebug {
			fmt.Fprintf(os.Stderr, "Debug: Successfully wrote %d bytes to: %s\n", n, absPath)
		}
    } else {
        fmt.Println(output)
    }
}