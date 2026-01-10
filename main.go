package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultTimeout = 500 * time.Millisecond
)

type ScanResult struct {
	Port int
}

func worker(host string, jobs <-chan int, results chan<- ScanResult, wg *sync.WaitGroup, timeout time.Duration, retries int, verbose bool, scannedPorts *int64) {
	defer wg.Done()

	for port := range jobs {
		address := fmt.Sprintf("%s:%d", host, port)
		var conn net.Conn
		var err error

		// Retry logic
		for attempt := 0; attempt <= retries; attempt++ {
			conn, err = net.DialTimeout("tcp", address, timeout)
			if err == nil {
				break
			}
			if verbose && attempt < retries {
				fmt.Printf("[RETRY] TCP %d (tentativo %d/%d)\n", port, attempt+1, retries)
			}
		}

		if err == nil {
			conn.Close()
			results <- ScanResult{Port: port}
			if verbose {
				fmt.Printf("[FOUND] Porta %d aperta\n", port)
			}
		}

		atomic.AddInt64(scannedPorts, 1)
	}
}

func progressBar(current, total int64) {
	percent := float64(current) / float64(total) * 100
	fmt.Printf("\r[Progress] %.1f%% (%d/%d)", percent, current, total)
}

func main() {
	// 1. Definizione dei Flag
	hostPtr := flag.String("host", "127.0.0.1", "Host da scansionare")
	startPtr := flag.Int("start", 1, "Porta iniziale")
	endPtr := flag.Int("end", 1024, "Porta finale")
	workersPtr := flag.Int("workers", 100, "Numero di worker concorrenti")
	timeoutPtr := flag.Int("timeout", 500, "Timeout in millisecondi")
	retriesPtr := flag.Int("retries", 0, "Numero di retry per timeout")
	verbosePtr := flag.Bool("verbose", false, "Output verbose")
	progressPtr := flag.Bool("progress", true, "Mostra progress bar")

	// 2. Parsing dei flag
	flag.Parse()

	targetHost := *hostPtr
	startPort := *startPtr
	endPort := *endPtr
	numWorkers := *workersPtr
	timeout := time.Duration(*timeoutPtr) * time.Millisecond
	retries := *retriesPtr
	verbose := *verbosePtr
	showProgress := *progressPtr

	// 3. Validazioni
	if startPort > endPort {
		fmt.Fprintln(os.Stderr, "Errore: La porta iniziale non può essere maggiore della finale.")
		os.Exit(1)
	}

	if startPort < 1 || endPort > 65535 {
		fmt.Fprintln(os.Stderr, "Errore: Le porte devono essere comprese tra 1 e 65535.")
		os.Exit(1)
	}

	if numWorkers < 1 || numWorkers > 10000 {
		fmt.Fprintln(os.Stderr, "Errore: Il numero di worker deve essere compreso tra 1 e 10000.")
		os.Exit(1)
	}

	totalPorts := int64(endPort - startPort + 1)

	fmt.Printf("Avvio scansione TCP su: %s (Porte %d-%d, Workers: %d)\n",
		targetHost, startPort, endPort, numWorkers)
	startTime := time.Now()

	// Setup canali
	jobs := make(chan int, numWorkers)
	results := make(chan ScanResult)
	var wg sync.WaitGroup
	var scannedPorts int64

	// Avvio Worker
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(targetHost, jobs, results, &wg, timeout, retries, verbose, &scannedPorts)
	}

	// Invio Job
	go func() {
		for i := startPort; i <= endPort; i++ {
			jobs <- i
		}
		close(jobs)
	}()

	// Monitoraggio chiusura
	go func() {
		wg.Wait()
		close(results)
	}()

	// Raccolta risultati con progress bar
	var openPorts []ScanResult

	// Progress bar in goroutine separata (disabilitata se verbose)
	stopProgress := make(chan bool)
	if showProgress && !verbose {
		go func() {
			ticker := time.NewTicker(100 * time.Millisecond)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					current := atomic.LoadInt64(&scannedPorts)
					progressBar(current, totalPorts)
				case <-stopProgress:
					return
				}
			}
		}()
	}

	// Raccolta risultati
	for result := range results {
		openPorts = append(openPorts, result)
	}

	// Stop progress bar
	if showProgress && !verbose {
		stopProgress <- true
		progressBar(totalPorts, totalPorts)
		fmt.Println()
	}

	// Ordina risultati
	sort.Slice(openPorts, func(i, j int) bool {
		return openPorts[i].Port < openPorts[j].Port
	})

	// Output
	fmt.Println("\n--- Scansione Completata ---")
	fmt.Printf("Target: %s\n", targetHost)
	if len(openPorts) == 0 {
		fmt.Println("Nessuna porta aperta trovata nel range specificato.")
	} else {
		fmt.Printf("Trovate %d porte aperte:\n", len(openPorts))
		for _, result := range openPorts {
			fmt.Printf("- Porta %d aperta\n", result.Port)
		}
	}

	fmt.Printf("\nTempo totale: %s\n", time.Since(startTime))
}
