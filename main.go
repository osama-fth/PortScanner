package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

// Configurazioni fisse
const (
	timeout = 500 * time.Millisecond
	workers = 100
)

func worker(host string, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for port := range jobs {
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, timeout)

		if err != nil {
			continue
		}

		conn.Close()
		results <- port
	}
}

func main() {
	// 1. Definizione dei Flag
	// Host (Stringa) - Default: 127.0.0.1
	hostPtr := flag.String("host", "127.0.0.1", "Host da scansionare")

	// Range Porte (Interi) - Default: 1 e 1024
	startPtr := flag.Int("start", 1, "Porta iniziale")
	endPtr := flag.Int("end", 1024, "Porta finale")

	// 2. Parsing dei flag
	flag.Parse()

	// Dereferenziamo i puntatori per ottenere i valori
	targetHost := *hostPtr
	startPort := *startPtr
	endPort := *endPtr

	// 3. Validazioni (Exit Failure in caso di errore)
	if startPort > endPort {
		fmt.Fprintln(os.Stderr, "Errore: La porta iniziale non può essere maggiore della finale.")
		os.Exit(1)
	}

	if startPort < 1 || endPort > 65535 {
		fmt.Fprintln(os.Stderr, "Errore: Le porte devono essere comprese tra 1 e 65535.")
		os.Exit(1)
	}

	fmt.Printf("Avvio scansione su: %s (Porte %d-%d)\n", targetHost, startPort, endPort)
	startTime := time.Now()

	// Setup canali e WaitGroup
	jobs := make(chan int, workers)
	results := make(chan int)
	var wg sync.WaitGroup

	// Avvio Worker
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go worker(targetHost, jobs, results, &wg)
	}

	// Invio Job (Produttore)
	// Utilizziamo le variabili startPort ed endPort ottenute dai flag
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

	// Raccolta risultati
	var openPorts []int
	for p := range results {
		openPorts = append(openPorts, p)
	}

	sort.Ints(openPorts)

	// Output
	fmt.Println("\n--- Scansione Completata ---")
	fmt.Printf("Target: %s\n", targetHost)
	if len(openPorts) == 0 {
		fmt.Println("Nessuna porta aperta trovata nel range specificato.")
	} else {
		for _, port := range openPorts {
			fmt.Printf("- Porta %d aperta\n", port)
		}
	}

	fmt.Printf("\nTempo totale: %s\n", time.Since(startTime))
}
