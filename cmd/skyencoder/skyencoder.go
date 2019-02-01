package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/skycoin/skyencoder"
)

/* TODO

- float32/float64 support

- generate test files that compare against the encoder for an encoded struct
	- how can we automate the creation of data though?
		- could have complicated way, where we build a data generator
			- write random bytes, except for length prefixes, which must be valid
	- can generate with an empty object, and let the user fill in the details?
		- regeneration would overwrite the filled in object

IN SKYCOIN:

- add skycoin tests (verify entire db can be loaded)

TO DOCUMENT:

- (wiki) document encoding format
- (wiki) document usage in golang (maxlen, omitempty, "-", ignores unexported fields)
- decoding specifics (empty map/slice are always left nil)

*/

const debug = false

func debugPrintln(args ...interface{}) {
	if debug {
		fmt.Println(args...)
	}
}

func debugPrintf(msg string, args ...interface{}) {
	if debug {
		fmt.Printf(msg, args...)
	}
}

var (
	typeName       = flag.String("type", "", "type name, must be set")
	outputFilename = flag.String("output-file", "", "output file name; default <type_name>_skyencoder.go")
	outputPath     = flag.String("output-path", "", "output path; defaults to the package's path, or the file's containing folder")
	buildTags      = flag.String("tags", "", "comma-separated list of build tags to apply")
	destPackage    = flag.String("package", "", "package name for the output; if not provided, defaults to the type's package")
	silent         = flag.Bool("silent", false, "disable all non-error log output")
	noTest         = flag.Bool("no-test", false, "disable generating the _test.go file (test files require github.com/google/go-cmp/cmp)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of skyencoder:\n")
	fmt.Fprintf(os.Stderr, "\tskyencoder [flags] -type T [go import path e.g. github.com/skycoin/skycoin/src/coin]\n")
	fmt.Fprintf(os.Stderr, "\tskyencoder [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("skyencoder: ")

	flag.Usage = usage
	flag.Parse()

	if *typeName == "" {
		flag.Usage()
		os.Exit(2)
	}

	var tags []string
	if len(*buildTags) > 0 {
		tags = strings.Split(*buildTags, ",")
	}

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	program, err := skyencoder.LoadProgram(args, tags)
	if err != nil {
		log.Fatal("skyencoder.LoadProgram failed:", err)
	}

	debugPrintln("args:", args)

	typeInfo, err := skyencoder.FindTypeInfoInProgram(program, *typeName)
	if err != nil {
		log.Fatalf("Program did not contain valid type for name %s: %v", *typeName, err)
	}
	if typeInfo == nil {
		log.Fatal("Program does not contain type:", *typeName)
	}

	// Determine if the arg is a directory or multiple files
	// If it is a directory, construct an artificial filename in that directory for goimports formatting,
	// otherwise use the first filename specified (they must all be in the same package)
	fmtFilename := args[0]
	destPath := filepath.Dir(args[0])

	stat, err := os.Stat(args[0])
	if err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
		// argument is a import path e.g. "github.com/skycoin/skycoin/src/coin"
		destPath, err = skyencoder.FindDiskPathOfImport(typeInfo.Package.Path())
		if err != nil {
			log.Fatal(err)
		}
		fmtFilename = filepath.Join(typeInfo.Package.Path(), "foo123123123123999.go")
	} else if stat.IsDir() {
		destPath = args[0]
		fmtFilename = filepath.Join(args[0], "foo123123123123999.go")
	}

	src, err := skyencoder.BuildTypeEncoder(typeInfo, *destPackage, fmtFilename)
	if err != nil {
		log.Fatal("skyencoder.BuildTypeEncoder failed:", err)
	}

	var testSrc []byte
	if !*noTest {
		testSrc, err = skyencoder.BuildTypeEncoderTest(typeInfo, *destPackage, fmtFilename)
		if err != nil {
			log.Fatal("skyencoder.BuildTypeEncoderTest failed:", err)
		}
	}

	debugPrintln(string(src))

	outputFn := *outputFilename
	if outputFn == "" {
		// If the input is a filename, put next to the file
		// If the input is a package, put in the package
		outputFn = fmt.Sprintf("%s_skyencoder.go", toSnakeCase(*typeName))
	}

	outputPth := *outputPath
	if outputPth == "" {
		outputPth = destPath
	}
	outputFn = filepath.Join(outputPth, outputFn)

	if !*silent {
		log.Printf("Writing skyencoder for type %q to file %q", *typeName, outputFn)
	}

	if err := ioutil.WriteFile(outputFn, src, 0644); err != nil {
		log.Fatal("ioutil.WriteFile failed:", err)
	}

	if !*noTest {
		outputExt := filepath.Ext(outputFn)
		base := outputFn[:len(outputFn)-len(outputExt)]
		testOutputFn := fmt.Sprintf("%s_test%s", base, outputExt)

		if !*silent {
			log.Printf("Writing skyencoder tests for type %q to file %q", *typeName, testOutputFn)
		}

		if err := ioutil.WriteFile(testOutputFn, testSrc, 0644); err != nil {
			log.Fatal("ioutil.WriteFile failed:", err)
		}
	}
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
