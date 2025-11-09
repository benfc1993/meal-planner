package main

import (
	"embed"
	"fmt"
	"log"
	"meal-choices/db"
	"meal-choices/routes"
	"meal-choices/templates"
	"net"
	"net/http"
	"os"
	"sync"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

//go:embed pages/*
//go:embed static/*
var f embed.FS

func main() {

	opts := getOpts()

	templates, err := templates.CreateTemplates(f, "pages")

	if err != nil {
		log.Fatal("Error parsing templates")
		log.Fatal(err)
	}

	/* api.Handlers() */
	routes.Handlers(f, templates)

	wg := new(sync.WaitGroup)

	wg.Go(func() {
		if opts.Https {
			log.Fatal(http.ListenAndServeTLS(fmt.Sprintf("%v:%v", opts.Host, opts.Port), "", "", nil))
		} else {
			log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", opts.Host, opts.Port), nil))
		}
	})

	prefix := "http"
	if opts.Https {
		prefix += "s"
	}
	fmt.Printf("Server running %v://%v:%v\n", prefix, opts.Host, opts.Port)

	db.Init()

	wg.Wait()
}

type Options struct {
	Port  string
	Host  string
	Https bool
}

func getOpts() *Options {
	options := &Options{Port: "5000", Host: "localhost", Https: false}

	for i := 0; i < len(os.Args); i++ {
		arg := os.Args[i]

		if string(arg[0]) == "-" {
			for {
				if string(arg[0]) == "-" {
					arg = arg[1:]
				} else {
					break
				}
			}

			switch arg {
			case "h", "help":
				data, _ := f.ReadFile("static/help.txt")
				fmt.Print(string(data))
				os.Exit(0)
			case "p", "port":
				options.Port = os.Args[i+1]
				i++
			case "u", "host":
				options.Host = os.Args[i+1]
				i++
			case "s", "https":
				options.Https = true
			default:
				log.Printf("Unknown option: %v\n", arg)
			}

		}

	}

	return options

}
