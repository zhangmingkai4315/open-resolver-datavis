package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/miekg/dns"
	"github.com/spf13/cobra"
	"github.com/zhangmingkai4315/open-resolver-datavis/config"
)

var tcpmode bool
var srcPath string
var destPath string

func init() {
	resolverCmd.Flags().BoolVarP(&tcpmode, "tcp", "t", false, "in tcp mode")
	resolverCmd.Flags().StringVarP(&srcPath, "path:src", "s", "", "source file path with raw ip address from zmap")
	resolverCmd.Flags().StringVarP(&destPath, "path:dest", "d", "", "output file path")
	rootCmd.AddCommand(resolverCmd)
}

var resolverCmd = &cobra.Command{
	Use:   "resolver",
	Short: "resolver will remove any ip from file not responese to dns requery or just refuse",
	Run: func(cmd *cobra.Command, args []string) {
		if srcPath == "" {
			log.Printf("Source file %s not set, please use --path:src to load file first", srcPath)
			os.Exit(1)
		}
		if destPath == "" {
			log.Printf("Dest file %s not set, please use --path:src to load file first", destPath)
			os.Exit(1)
		}
		log.Printf("Start send dns query packet to candidate open reslover and save to %s", destPath)
		if tcpmode == false {
			err := collectOpenDNSResolver(srcPath, destPath, false)
			if err != nil {
				log.Printf("Error:%s", err.Error())
				os.Exit(1)
			}
		} else {
			// query in tcp mode
			err := collectOpenDNSResolver(srcPath, destPath, true)
			if err != nil {
				log.Printf("Error:%s", err.Error())
				os.Exit(1)
			}

		}
		os.Exit(0)
	},
}

func queryWorker(tcpmode bool) func(id int, jobs <-chan string, results chan<- string) {
	c := new(dns.Client)
	if tcpmode == true {
		c.Net = "tcp"
	}
	return func(id int, jobs <-chan string, results chan<- string) {
		for ip := range jobs {
			m := new(dns.Msg)
			m.RecursionDesired = true
			m.SetQuestion("google.com.", dns.TypeA)
			in, _, err := c.Exchange(m, ip+":53")
			if err != nil {
				continue
			}
			if in.Rcode == dns.RcodeSuccess {
				results <- ip
			}
		}
	}
}

func collectOpenDNSResolver(srcPath string, destPath string, tcpmode bool) error {
	dnsworker := queryWorker(tcpmode)
	// open src file read million lines
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	// create dest file and write results to this file
	log.Printf("Create destination ip file %s\n", destPath)
	dest, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dest.Close()

	// create worker to query dns and save result to channes
	scanner := bufio.NewScanner(f)
	cfg := config.GetAppConfig()
	worker := cfg.GlobalConfig.Worker
	jobs := make(chan string, 100)
	results := make(chan string, 10000)

	log.Printf("Start %d workers to send dns query\n", worker)
	for w := 1; w <= worker; w++ {
		go dnsworker(w, jobs, results)
	}

	log.Printf("Create write channel to save resolver data\n")
	go func(dest *os.File, results chan string) {
		w := bufio.NewWriter(dest)
		counter := 0
		for r := range results {
			fmt.Fprintln(w, r)
			counter++
			if counter%1000 == 0 {
				w.Flush()
			}
		}
		w.Flush()
	}(dest, results)
	log.Printf("Start read raw files: %s", srcPath)
	for scanner.Scan() {
		line := scanner.Text()
		jobs <- line
	}

	close(jobs)
	close(results)
	return nil
}
