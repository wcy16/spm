package config

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// read the .config files in config/ folder
// the file name exclude the file type
func Read(name string) map[string]string {
	file, err := os.Open("config/" + name + ".config")
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	cfg := make(map[string]string)
	for scanner.Scan() {
		vals := strings.SplitN(scanner.Text(), "=", 2)
		if len(vals) < 2 {
			continue
		}
		cfg[vals[0]] = vals[1]
	}
	return cfg
}
