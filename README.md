# Learn Go

> 💡 “Golang” is just SEO. The language is called **Go**. The gopher is the unofficial mascot.

| **Description** | **Choice** |
| --- | --- |
| Build toolchain | `go` (compiler, linker, toolchain) |
| Dependency management | **Go Modules** (`go.mod`, `go.sum`) |
| LSP / Editor brain | `gopls` |
| Testing | `testing` (built-in), `go test` |
| Linting (popular) | `staticcheck`, `golangci-lint` |
| HTTP router | `chi`, `gin`, or stdlib `net/http` |
| CLI scaffolding | `cobra` (+ `viper` for config) |
| Serialization | `encoding/json` (stdlib), `proto`/`gRPC` |
| Observability | `pprof`, `expvar`, OpenTelemetry |
| Concurrency prims | goroutines, channels, `sync` |
| Release packaging | static binaries, `-ldflags` |
| Task runner (light) | `go run` / `go generate` |

## A brief history of Go

- **2007–2009**: Designed at Google by Robert Griesemer, Rob Pike, Ken Thompson. Goals: compile fast, easy concurrency, batteries included standard library.  
- **Go 1 (2012)**: Compatibility promise: 1.x should keep working.  
- **Modules (2018)**: `go mod` replaces GOPATH.  
- **Generics + fuzzing (2022)**: 1.18 adds type parameters and fuzz tests.  
- Ongoing: faster compiler, tighter GC pauses, language tweaks (`any`, `comparable`, `min`/`max` in `constraints`)

## Why Go?

- **Fast compile + static binaries** (easy deploys, tiny containers).  
- **Great concurrency model** (goroutines, channels, `select`).  
- **Predictable performance** (GC optimized for servers).  
- **Excellent tooling** (formatter, vet, race detector, pprof).  
- **Cross-compilation** with `GOOS`/`GOARCH`.

## How Go works behind the scenes

- **Goroutines**: Lightweight functions (KB stacks) scheduled by the runtime.  
- **Scheduler (G–M–P)**: Goroutines (G) on OS threads (M) via logical processors (P).  
- **Channels**: Typed pipes; blocking send/recv by default.  
- **GC**: Concurrent, tri-color; short STW pauses.  
- **Maps**: Hash maps; **not safe for concurrent writes**.

## Fundamental concepts

### Modules, packages, files
- One **module** = one `go.mod`.  
- One **package** per directory.  
- `package main` builds an executable; anything else is a library.  
- Workspaces: `go work` to stitch multiple modules.

### Types & zero values
- Zero values: `0`, `""`, `false`, `nil` (for ptr/slice/map/chan/func).
- **Arrays** are fixed length.
- **Slices** are views over arrays; capacity grows geometrically.  
- **Maps**: `make(map[K]V)`; nil maps panic on write.  
- **Strings**: immutable UTF-8; use `rune` for Unicode.

### Structs, methods, interfaces
```go
type User struct{ ID int; Name string }
func (u *User) Rename(n string) { u.Name = n } // pointer receiver mutates

type Stringer interface{ String() string } // interfaces satisfied implicitly
```
- Prefer **small interfaces** at call sites (e.g., `io.Reader`).  
- Use **embedding** for composition.

### Generics (1.18+)
```go
func Min[T constraints.Ordered](a, b T) T {
    if a < b { return a }
    return b
}
```
- Great for containers/algorithms; keep APIs simple.

### Errors, panic, defer
```go
f, err := os.Open(p)
if err != nil { return fmt.Errorf("open %s: %w", p, err) }
defer f.Close()
```
- **Errors are values**. Wrap with `%w` (think wtf); use `errors.Is/As`.  
- `panic` for truly exceptional states; `recover` only at boundaries.  
- `defer` runs LIFO; args evaluated at **defer time**.

### Context (cancellation & deadlines)
```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()
```
- Pass `ctx context.Context` as the **first** param for I/O-y funcs.  
- Cancel what you start; avoid leaking goroutines.

## Concurrency patterns

### WaitGroup fan-out
```go
var wg sync.WaitGroup
for _, j := range jobs {
    wg.Add(1)
    go func(j Job){ defer wg.Done(); _ = do(j) }(j)
}
wg.Wait()
```

### Channels & select
```go
select {
case v := <-ch:
    use(v)
case <-ctx.Done():
    return ctx.Err()
}
```

### Sync tools
- `sync.Mutex` / `RWMutex` for shared data.  
- `sync.Map` for special high-concurrency read-heavy cases.  
- `atomic` for counters/flags.

## HTTP & JSON (stdlib wins)

```go
mux := http.NewServeMux()
mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    _ = json.NewEncoder(w).Encode(map[string]string{"msg":"world"})
})
srv := &http.Server{ Addr: ":8080", Handler: mux }
log.Fatal(srv.ListenAndServe())
```

- Prefer a **single shared `http.Client`** with timeouts.  
- JSON struct tags:
```go
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name,omitempty"`
}
```

## Tooling (built-in superpowers)

```bash
go fmt ./...                  # format
go vet ./...                  # static checks
go test ./...                 # run tests
go test -race ./...           # data race detector
go test -bench=. -benchmem    # benchmarks
go test -fuzz=Fuzz -fuzztime=10s  # fuzzing
go mod tidy                   # prune/add deps
go build ./cmd/app            # build
go run ./cmd/app              # build+run
go doc net/http               # docs in terminal
```

## Modules & versions

- **Semantic Import Versioning**: `v2+` adds `/v2` suffix in import path.  
- `replace` in `go.mod` to pin local forks.  
- Private repos: set `GOPRIVATE` and configure auth.  
- Vendoring: `go mod vendor` if you need vendored deps.

## Testing style

- **Table-driven tests**
```go
func TestSum(t *testing.T) {
    cases := []struct{ in []int; want int }{
        {[]int{}, 0}, {[]int{1,2,3}, 6},
    }
    for _, c := range cases {
        if got := Sum(c.in...); got != c.want {
            t.Fatalf("got %d want %d", got, c.want)
        }
    }
}
```
- Benchmarks with `testing.B`, fuzz with `FuzzXxx`.  
- Coverage: `go test -coverprofile=cover.out && go tool cover -html=cover.out`.

## Build & release tips

- **Cross-compile**: `GOOS=linux GOARCH=arm64 go build`.  
- **Static-ish**: `CGO_ENABLED=0` for pure-Go binaries (no C dependencies, so can run in `FROM scratch`/`FROM alpine` Docker container).  
- **Embed files**:
```go
//go:embed static/*
var staticFS embed.FS
```
- **Inject version info**:
```bash
go build -ldflags "-X main.version=$(git describe --tags) -s -w"
```

## Performance & memory notes

- Minimize allocations (reuse buffers; `bytes.Buffer`; `sync.Pool` cautiously).  
- **Escape analysis**: `go build -gcflags=all=-m`.  
- Preallocate: `make([]T, 0, n)`.  
- Avoid copying large structs; use pointers.  
- **pprof**: `import _ "net/http/pprof"` and hit `/debug/pprof`.

## Common gotchas

- **Concurrent map writes** panic; use a mutex.  
- **Nil maps** can’t be written to.  
- **Loop var capture**:
```go
for _, v := range xs {
    v := v // shadow before goroutine
    go func(){ fmt.Println(v) }()
}
```
- **Slice aliasing**: reslicing shares backing array; mutations may leak.  
- **`defer` in tight loops** can balloon; close promptly.  
- **Interface-nil trap**: typed nil ≠ untyped nil.  
- **Time layouts** use `2006-01-02 15:04:05`.  
- Always **close `resp.Body`** and set **timeouts**.

## Minimal project layout

Looks arbitrary until you’ve lived through the pain.

```
.
├─ cmd/          # thin entrypoints
│  └─ api/       # each subfolder is a separate executable
│     └─ main.go
├─ internal/     # private visibility boundary enforced by the compiler for reusable code
│  ├─ http/      # handlers, routers
│  └─ store/     # db access
├─ pkg/          # (optional) public helpers
├─ go.mod
└─ go.sum
```

## Handy stdlib you’ll actually use

- **I/O**: `io`, `os`, `bufio`, `bytes`, `strings`  
- **Net**: `net`, `net/http`, `net/url`, `net/netip`  
- **Time**: `time`, `context`  
- **Encoders**: `encoding/json`, `encoding/xml`, `encoding/csv`, `encoding/binary`  
- **Crypto**: `crypto/*`, `hash/*`  
- **Sync**: `sync`, `sync/atomic`  
- **Text**: `regexp`, `unicode/utf8`, `text/template`, `html/template`  
- **DB**: `database/sql` + driver (e.g., `lib/pq`, `mysql`, `sqlite`)  

## Tiny HTTP server (test with curl)

```go
package main

import (
  "encoding/json"
  "log"
  "net/http"
)

func main() {
  http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
    _ = json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
  })
  log.Fatal(http.ListenAndServe(":8080", nil))
}
```

## SQL with `database/sql` (idiomatic sketch)

```go
db, err := sql.Open("postgres", dsn) // use context for real calls
if err != nil { log.Fatal(err) }
defer db.Close()

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

var name string
err = db.QueryRowContext(ctx, `SELECT name FROM users WHERE id=$1`, 7).Scan(&name)
if err != nil { /* handle */ }
```

## Reference patterns

**Functional options**
```go
type Server struct{ addr string }
type Option func(*Server)
func WithAddr(a string) Option { return func(s *Server){ s.addr = a } }
func New(opts ...Option) *Server { s := &Server{addr:":8080"}; for _,o := range opts{o(s)}; return s }
```

**Logger interface**
```go
type Logger interface{ Printf(string, ...any) }
```

## Observability quickies

- **Race detector**: `go test -race ./...`  
- **pprof** CPU profile:
```bash
go test -run=NONE -bench=. -benchmem -cpuprofile=cpu.out
go tool pprof cpu.out
```
- **OpenTelemetry**: instrument HTTP, gRPC, DB; export to Jaeger/OTLP.

## Cheat-sheet: one-liners

```bash
go mod init example.com/x
go get github.com/acme/lib@v1
go mod tidy
go fmt ./... && go vet ./...
go test ./... -race
go test -bench=. -benchmem
GOOS=linux GOARCH=arm64 go build -o app-linux-arm64 .
```

## Random notes

- Prefer **small, focused packages**. Fewer exported symbols.  
- Return **concrete types**, accept **interfaces**.  
- Keep APIs boring; **clever code ages poorly**.  
- Start with stdlib; add frameworks only when needed.  
- Measure before optimizing; **profilers tell the truth**.

> 💡 The **Go 1 compatibility promise** means upgrading Go should rarely break your code.

## Trivia

- Function signatures of the same type `(s1 string, s2 string)` can be shortened to `(s1, s2 string)`
