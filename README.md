# PortScanner

![Go Version](https://img.shields.io/badge/Go-1.16+-00ADD8?style=flat&logo=go)

Un semplice port scanner scritto in Go che utilizza goroutine per scansionare velocemente le porte TCP di un host target.

## Caratteristiche

- 🚀 Scansione concorrente con 100 worker paralleli
- ⚡ Timeout configurato a 500ms per connessione
- 🎯 Range di porte personalizzabile
- 📊 Output ordinato delle porte aperte

## Installazione

```bash
go build -o portscanner main.go
```

## Utilizzo

```bash
./portscanner [flags]
```

### Flag Disponibili

| Flag | Tipo | Default | Descrizione |
|------|------|---------|-------------|
| `-host` | string | `127.0.0.1` | Host o indirizzo IP da scansionare |
| `-start` | int | `1` | Porta iniziale del range da scansionare |
| `-end` | int | `1024` | Porta finale del range da scansionare |

### Esempi

**Scansione di localhost (porte di default):**
```bash
./portscanner
```

**Scansione di un host specifico:**
```bash
./portscanner -host scanme.nmap.org
```

**Scansione di un range limitato di porte:**
```bash
./portscanner -host scanme.nmap.org -start 1 -end 1000
```

**Scansione delle porte comuni:**
```bash
./portscanner -host scanme.nmap.org -start 20 -end 443
```

## Test consigliato

Per testare lo scanner, si consiglia di utilizzare **scanme.nmap.org**, un host pubblico messo a disposizione dal team di Nmap per testare scanner di rete:

```bash
./portscanner -host scanme.nmap.org -start 1 -end 1000
```

⚠️ **Nota importante**: Scansionare host senza autorizzazione può essere illegale. Utilizza questo strumento solo su sistemi di cui possiedi l'autorizzazione o su host pubblici come scanme.nmap.org.

## Output

Lo scanner visualizzerà:
- Host target e range di porte
- Elenco delle porte aperte trovate
- Tempo totale di esecuzione

Esempio di output:
```
Avvio scansione su: scanme.nmap.org (Porte 1-1000)

--- Scansione Completata ---
Target: scanme.nmap.org
- Porta 22 aperta
- Porta 80 aperta

Tempo totale: 8.5s
```

## Requisiti

- Go 1.16 o superiore

---

