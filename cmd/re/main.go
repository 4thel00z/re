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
)

const maxFileSize = 10 * 1024 * 1024

type Config struct {
	Needle      string
	Replacement string
	Directories []string
	FlagApply   bool
	Patterns    FilePattern
	Replaced    int
}

type FilePattern struct {
	Excludes []regexp.Regexp
	Includes []regexp.Regexp
}

func main() {
	config := parseFlags()
	if !config.FlagApply {
		fmt.Println("No changes will be applied unless -f is given.")
	}

	for _, dir := range config.Directories {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			return walk(path, info, err, config)
		}); err != nil {
			log.Printf("Cannot scan directory: %s; Skipping.\n", err)
		}
	}

	msg := "file(s) WOULD have been updated."
	if config.FlagApply {
		msg = "file(s) were updated."
	}
	fmt.Printf("%d %s\n", config.Replaced, msg)
}

func parseFlags() *Config {

	flagApply := flag.Bool("f", false, "Apply changes")
	flagExcludes := flag.String("e", ".bzr,CVS,.git,.hg,.svn", "Comma-separated list of excluded files, wildcards supported")
	flagIncludes := flag.String("i", "", "Comma-separated list of included files, wildcards supported, e.g \"*.js,*.html,*index.*\"")

	flag.Usage = func() {
		fmt.Println("Syntax: re [options] SEARCH REPLACEMENT [DIR ...]")
		fmt.Println()
		fmt.Println("Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	dirs := []string{"."}
	if len(flag.Args()) > 2 {
		dirs = flag.Args()[2:]
	}

	return &Config{
		Needle:      flag.Arg(0),
		Replacement: flag.Arg(1),
		Directories: dirs,
		FlagApply:   *flagApply,
		Patterns: FilePattern{
			Excludes: parseRegexList(*flagExcludes),
			Includes: parseRegexList(*flagIncludes),
		},
	}
}

func parseRegexList(s string) []regexp.Regexp {
	var result []regexp.Regexp
	for _, str := range strings.Split(s, ",") {
		if str == "" {
			continue
		}
		pattern := regexp.QuoteMeta(str)
		pattern = strings.Replace(pattern, "\\*", ".*", -1)
		pattern = strings.Replace(pattern, "\\?", ".", -1)
		pattern = "^" + pattern + "$"

		regex, err := regexp.Compile(pattern)
		if err != nil {
			panic(err)
		}
		result = append(result, *regex)
	}
	return result
}

func walk(path string, info os.FileInfo, err error, config *Config) error {
	if err != nil {
		log.Printf("Error accessing path %s: %s", path, err)
		return nil
	}

	if info.Size() >= maxFileSize || info.IsDir() {
		return nil
	}

	if fileExcluded(path, config.Patterns.Excludes) {
		return nil
	}

	if len(config.Patterns.Includes) > 0 && !fileIncluded(path, config.Patterns.Includes) {
		return nil
	}

	contents, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error reading %s. Skipping.", path)
		return nil
	}

	oldContents := string(contents)
	newContents := strings.Replace(oldContents, config.Needle, config.Replacement, -1)

	if oldContents == newContents {
		return nil
	}

	if config.FlagApply {
		if err := ioutil.WriteFile(path, []byte(newContents), info.Mode()); err != nil {
			log.Printf("Error writing file %s: %s.", path, err)
		}
	}
	config.Replaced++
	fmt.Println("+", path)

	return nil
}

func fileExcluded(filename string, excludes []regexp.Regexp) bool {
	for _, exclude := range excludes {
		if exclude.MatchString(filename) {
			log.Printf("Skipping excluded file %s", filename)
			return true
		}
	}
	return false
}

func fileIncluded(filename string, includes []regexp.Regexp) bool {
	for _, include := range includes {
		if include.MatchString(filename) {
			return true
		}
	}
	return false
}
