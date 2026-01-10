# PortScanner

![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)

Un semplice port scanner scritto in Go che utilizza goroutine per scansionare velocemente le porte TCP di un host target.

## Caratteristiche

- 🚀 Scansione concorrente con worker configurabili (default 100)
- ⚡ Timeout configurato a 500ms per connessione
- 🎯 Range di porte personalizzabile

## Installazione

```bash
go build -o portscanner.out main.go
```

## Utilizzo

```bash
./portscanner.out [flags]
```

### Flag Disponibili

| Flag | Tipo | Default | Descrizione |
|------|------|---------|-------------|
| `-host` | string | `127.0.0.1` | Host o indirizzo IP da scansionare |
| `-start` | int | `1` | Porta iniziale del range da scansionare |
| `-end` | int | `1024` | Porta finale del range da scansionare |
| `-workers` | int | `100` | Numero di worker concorrenti (goroutine) |
| `-timeout` | int | `500` | Timeout per connessione in millisecondi |
| `-retries` | int | `0` | Numero di tentativi su timeout |
| `-verbose` | bool | `false` | Mostra output dettagliato |
| `-progress` | bool | `true` | Mostra progress bar |

### Esempi

**Scansione base di localhost:**
```bash
./portscanner
```

**Scansione di un host remoto:**
```bash
./portscanner -host scanme.nmap.org
```

**Scansione completa di tutte le porte:**
```bash
./portscanner -start 1 -end 65535 -workers 500
```

**Scansione con retry e timeout maggiore:**
```bash
./portscanner -host example.com -timeout 1000 -retries 2
```

**Modalità verbose senza progress bar:**
```bash
./portscanner -verbose -start 20 -end 100
```

**Scansione veloce con più worker:**
```bash
./portscanner -host 192.168.1.1 -workers 300 -timeout 200
```

## Note
- Il timeout e i retry influenzano significativamente il tempo di scansione
- Usare più worker velocizza la scansione ma può sovraccaricare la rete
