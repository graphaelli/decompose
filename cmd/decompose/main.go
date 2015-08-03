package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/variadico/decompose"
)

const usageText = `DECOMPOSE
    Convert Docker Compose YAML into Docker commands.

USAGE
    decompose [options] path

OPTIONS
    -d, -detach
        Generate docker run commands with detach flag. This flag can't be
        used at the same time as -r.

    -p, -prefix
        Set a string to prefix names with. Default is directory name of
        input file. For no prefix, provide an empty string.

    -r, -rm
        Generate docker run commands with the remove flag. This flag can't
        be used at the same time as -d.

    -s, -service
        Only generate docker commands for a particular service. A service can
        be specified by name or index number. Indexes start at 1.

    -u, -units
        Generate systemd unit files. NOT IMPLEMENTED.

    -h, -help
        Print help and exit.
`

func main() {
	log.SetFlags(0)

	detach := flag.Bool("d", false, "")
	flag.BoolVar(detach, "detach", false, "")
	prefix := flag.String("p", "docker-compose dir", "")
	flag.StringVar(prefix, "prefix", "docker-compose dir", "")
	remove := flag.Bool("r", false, "")
	flag.BoolVar(remove, "rm", false, "")
	service := flag.String("s", "", "")
	flag.StringVar(service, "service", "", "")
	units := flag.Bool("u", false, "")
	flag.BoolVar(units, "units", false, "")
	help := flag.Bool("h", false, "")
	flag.BoolVar(help, "help", false, "")
	flag.Usage = func() { log.Println(usageText) }
	flag.Parse()

	if *help {
		fmt.Println(usageText)
		return
	}

	if *detach && *remove {
		log.Println("can't use -detach and -rm at the same time")
		log.Fatal(usageText)
	}

	if len(flag.Args()) != 1 {
		log.Println("invalid path argument count")
		log.Fatal(usageText)
	}

	inputFile := flag.Args()[0]

	if *prefix == "docker-compose dir" {
		// prefix is the directory that inputFile is in
		finAbs, err := filepath.Abs(inputFile)
		if err != nil {
			log.Fatal(err)
		}

		*prefix = filepath.Base(filepath.Dir(finAbs)) + "_"
	}

	decompose.Prefix = *prefix
	decompose.Detach = *detach
	decompose.Remove = *remove

	services, err := decompose.ParseComposeFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	if *service == "" {
		for _, s := range services {
			fmt.Println(s)
		}
		return
	}

	if servsIndex, err := strconv.Atoi(*service); err == nil {
		// service correctly parsed to an int, get index
		if servsIndex > 0 && servsIndex <= len(services) {
			fmt.Println(services[servsIndex-1])
		} else {
			log.Println("service index out of bounds")
			log.Fatal(usageText)
		}
	} else {
		// service wasn't an int, get name
		for _, s := range services {
			if s.Name == *service {
				fmt.Println(s)
				return
			}
		}

		log.Fatal("service name not found")
	}

}
