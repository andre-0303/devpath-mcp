# Roadmap Golang — Referência Completa

18 tópicos em ordem de progressão. Use estes nomes exatos ao chamar as ferramentas MCP (case-insensitive).

---

## Mapa de dependências

```
Variables ──┬──► Constants ──► Functions ──┬──► Conditionals ──► Loops ──┬──► Arrays ──► Slices
            │                              │                              │
            │                              └──► Structs ──► Methods ──►  │   Maps ◄────────────┘
            │                                      │          │           │
            │                              ────────┘          └──► Interfaces
            │
            └──► (base de tudo)

Functions ──► Error Handling ──► Packages ──┐
                                            ├──► REST API
Structs ────────────────────────────────────┤
                                            │
Functions + Loops ──► Concurrency ──► Channels ──► Context ──┘

Functions + Packages ──► Testing
```

---

## Tabela de tópicos

| # | Nome | Dificuldade | Tempo | Pré-requisitos |
|---|------|-------------|-------|----------------|
| 1 | Variables | beginner | 20 min | — |
| 2 | Constants | beginner | 15 min | Variables |
| 3 | Functions | beginner | 30 min | Variables, Constants |
| 4 | Conditionals | beginner | 20 min | Variables, Functions |
| 5 | Loops | beginner | 25 min | Variables, Conditionals |
| 6 | Arrays | beginner | 20 min | Variables, Loops |
| 7 | Slices | beginner | 35 min | Arrays |
| 8 | Maps | beginner | 30 min | Variables, Loops |
| 9 | Structs | intermediate | 40 min | Variables, Functions |
| 10 | Methods | intermediate | 35 min | Structs |
| 11 | Interfaces | intermediate | 45 min | Methods |
| 12 | Error Handling | intermediate | 40 min | Functions |
| 13 | Packages | intermediate | 30 min | Functions, Error Handling |
| 14 | Concurrency | advanced | 60 min | Functions, Loops |
| 15 | Channels | advanced | 50 min | Concurrency |
| 16 | Context | advanced | 45 min | Channels |
| 17 | Testing | intermediate | 40 min | Functions, Packages |
| 18 | REST API | advanced | 90 min | Structs, Error Handling, Packages, Context |

**Total:** ~615 min (~10h de estudo)

---

## Detalhes por tópico

### 1. Variables — beginner (20 min)

**Desafios:**
1. Declare variáveis de 5 tipos diferentes (`int`, `float64`, `string`, `bool`, `byte`) e imprima cada uma.
2. Troque dois inteiros sem variável temporária.
3. Use `:=` dentro de uma função e explique por que não funciona em package level.

**Hints:**
- `fmt.Printf` com `%T` imprime o tipo da variável.
- Tuple assignment: `a, b = b, a`
- `:=` é syntax sugar para `var` + atribuição; o compilador infere o tipo.

---

### 2. Constants — beginner (15 min)

**Pré-requisitos:** Variables

**Desafios:**
1. Enum com `iota` para dias da semana (Sunday=0 a Saturday=6).
2. Constante tipada para `pi`; mostrar que não pode ser reatribuída.
3. Flags de permissão com `iota` e bit-shifting (Read=1, Write=2, Execute=4).

**Hints:**
- `iota` reseta para 0 a cada bloco `const`.
- `const Pi = 3.14159` — descomentar `Pi = 3.0` gera erro de compilação.
- `const ( Read = 1 << iota; Write; Execute )` — shift duplica o valor.

---

### 3. Functions — beginner (30 min)

**Pré-requisitos:** Variables, Constants

**Desafios:**
1. Função que retorna quociente e resto (múltiplos retornos).
2. Função variádica `sum(...int) int`.
3. Higher-order function que aplica `func(int) int` a cada elemento de um slice.

**Hints:**
- `func divide(a, b int) (int, int) { return a / b, a % b }`
- Dentro da função variádica, `nums` é um `[]int`.
- `func apply(nums []int, f func(int) int) []int` — `range` sobre nums, chama `f` em cada elemento.

---

### 4. Conditionals — beginner (20 min)

**Pré-requisitos:** Variables, Functions

**Desafios:**
1. Função que classifica número: negativo, zero, positivo, grande (>1000).
2. `switch` mapeando mês (1-12) para estação do ano.
3. FizzBuzz 1-100 com `switch` sem condição.

**Hints:**
- `if-else` em Go exige chaves mesmo para uma linha. Sem parênteses na condição.
- `case 12, 1, 2:` — múltiplos valores por case.
- `switch { case n%15==0: ... }` — primeiro match vence.

---

### 5. Loops — beginner (25 min)

**Pré-requisitos:** Variables, Conditionals

**Desafios:**
1. Fatorial de n com `for` (tratar n=0).
2. Contar frequência de caracteres em string com `range`.
3. Loop estilo `while` lendo de channel até sentinel.

**Hints:**
- Go só tem `for` — sem `while` ou `do-while`.
- `range` sobre string retorna runes (Unicode), não bytes.
- `for val := range ch` — sai quando o channel é fechado.

---

### 6. Arrays — beginner (20 min)

**Pré-requisitos:** Variables, Loops

**Desafios:**
1. Array de 5 inteiros com quadrados (1,4,9,16,25).
2. Reverter `[5]int` in-place sem memória extra.
3. Máximo e mínimo em passagem única.

**Hints:**
- Tamanho é parte do tipo: `[5]int` e `[6]int` são tipos diferentes.
- Swap `arr[i]` e `arr[len(arr)-1-i]` enquanto `i < len(arr)/2`.
- Inicializar `max = arr[0]`, `min = arr[0]`, iterar a partir do índice 1.

---

### 7. Slices — beginner (35 min)

**Pré-requisitos:** Arrays

**Desafios:**
1. Remover duplicatas preservando ordem.
2. Stack (push, pop, peek) com slice como storage.
3. Flatten de `[][]int` em `[]int`.

**Hints:**
- `map[int]bool` como seen-set; append só se `!seen[v]`.
- Pop: `v, stack = stack[len-1], stack[:len-1]`.
- Pré-alocar com `make([]int, 0, totalLen)`.

---

### 8. Maps — beginner (30 min)

**Pré-requisitos:** Variables, Loops

**Desafios:**
1. Frequência de palavras em string com `map[string]int`.
2. Inverter `map[string]string` (troca keys e values).
3. Cache in-memory com get/set/delete.

**Hints:**
- `strings.Fields` para split; `freq[word]++` por palavra.
- Valores duplicados: o último key vence — ou coletar colisões em `[]string`.
- `val, ok := cache[key]` — `ok` é false se key não existe.

---

### 9. Structs — intermediate (40 min)

**Pré-requisitos:** Variables, Functions

**Desafios:**
1. `BankAccount` com owner, balance, accountNumber e construtor `NewBankAccount`.
2. Embedded struct: `Person` (name/age) → `Employee` (Person + role/salary).
3. Deep-copy de struct com campo slice.

**Hints:**
- Pointer receiver para métodos que modificam estado: `func (a *BankAccount) Deposit(amount float64)`.
- `type Employee struct { Person; Role string }` — acesso direto: `e.Name`.
- `copy(newSlice, original.Items)` para deep-copy de slice.

---

### 10. Methods — intermediate (35 min)

**Pré-requisitos:** Structs

**Desafios:**
1. `Deposit`, `Withdraw`, `Balance` no BankAccount; Withdraw retorna error se saldo insuficiente.
2. Implementar `String() string` em tipo customizado para `fmt.Println` formatar bem.
3. `Rectangle` e `Circle` com métodos `Area()` e `Perimeter()`.

**Hints:**
- `func (a *BankAccount) Withdraw(amount float64) error`
- `func (b BankAccount) String() string` faz o tipo implementar `fmt.Stringer` automaticamente.
- Value receiver para métodos read-only é idiomático quando struct é pequena.

---

### 11. Interfaces — intermediate (45 min)

**Pré-requisitos:** Methods

**Desafios:**
1. Interface `Shape` com `Area()` e `Perimeter()`; implementar para Rectangle e Circle.
2. Função `PrintShapeInfo(s Shape)` chamada com ambos os tipos.
3. Type switch para detectar Rectangle, Circle, ou tipo desconhecido.

**Hints:**
- Interfaces satisfeitas implicitamente — sem keyword `implements`.
- `switch v := s.(type) { case Rectangle: ... }`

---

### 12. Error Handling — intermediate (40 min)

**Pré-requisitos:** Functions

**Desafios:**
1. `SafeDivide(a, b float64) (float64, error)` retornando error para divisão por zero.
2. Tipo de erro customizado `ValidationError` com `Field` e `Message`.
3. Wrap com `fmt.Errorf("context: %w", err)` e unwrap com `errors.As` e `errors.Is`.

**Hints:**
- Interface `error` tem um método: `Error() string`. Retornar `nil` para sucesso.
- `errors.As(err, &target)` preenche target se a chain contém `*ValidationError`.
- `errors.Is` verifica identidade.

---

### 13. Packages — intermediate (30 min)

**Pré-requisitos:** Functions, Error Handling

**Desafios:**
1. Extrair BankAccount para package `bank`; exportar só o necessário.
2. Função helper unexported; verificar que não é acessível de fora.
3. `init()` para validar configuração no startup.

**Hints:**
- Exported = maiúscula (`BankAccount`). Unexported = minúscula (`validateAmount`).
- Acesso a unexported de outro package = erro de compilação.
- `func init()` roda antes de `main()`; usar para setup, não lógica de negócio.

---

### 14. Concurrency — advanced (60 min)

**Pré-requisitos:** Functions, Loops

**Desafios:**
1. 5 goroutines com índice; `WaitGroup` para aguardar todas.
2. Identificar data race em counter; corrigir com `sync.Mutex`.
3. `sync.Once` para inicializar recurso compartilhado exatamente uma vez.

**Hints:**
- `go run -race` detecta data races.
- `var wg sync.WaitGroup; wg.Add(5); go func() { defer wg.Done(); ... }(); wg.Wait()`
- `var once sync.Once; once.Do(func() { ... })` — seguro para chamadas concorrentes.

---

### 15. Channels — advanced (50 min)

**Pré-requisitos:** Concurrency

**Desafios:**
1. Pipeline: generator → square → printer, conectados por channels.
2. Fan-out: um input channel, múltiplos workers.
3. `select` com `time.After` para timeout de 2 segundos.

**Hints:**
- Cada stage lê do input channel e envia para output channel; fechar output quando input fecha.
- Múltiplas goroutines podem ler do mesmo channel — Go distribui sends entre receivers.
- `select { case result := <-work: ...; case <-time.After(2*time.Second): return errors.New("timeout") }`

---

### 16. Context — advanced (45 min)

**Pré-requisitos:** Channels

**Desafios:**
1. `context.WithTimeout` em função que simula trabalho; cancelar e observar.
2. Propagar request ID via `context.WithValue` por cadeia de chamadas.
3. Função que respeita `ctx.Done()` em select.

**Hints:**
- `ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()`
- Usar typed key (não string) para evitar colisão: `type ctxKey string`
- `select { case <-ctx.Done(): return ctx.Err(); case result := <-workCh: ... }`

---

### 17. Testing — intermediate (40 min)

**Pré-requisitos:** Functions, Packages

**Desafios:**
1. Testes table-driven para `SafeDivide`: caso normal, divisão por zero, negativos.
2. `t.Run` para subtests; `t.Parallel()` para rodá-los concorrentemente.
3. Benchmark `BenchmarkSafeDiv` comparando com e sem o zero-check.

**Hints:**
- `tests := []struct{ a, b float64; wantErr bool }{ ... }; for _, tt := range tests { t.Run(...) }`
- `t.Parallel()` deve ser chamado no início de cada subtest; o outer test NÃO chama.
- `go test -bench=.` para rodar benchmarks.

---

### 18. REST API — advanced (90 min)

**Pré-requisitos:** Structs, Error Handling, Packages, Context

**Desafios:**
1. CRUD de Todo com `net/http` (sem framework): POST/GET/GET/{id}/DELETE `/todos`.
2. Validação de JSON: retornar 400 com body de erro estruturado.
3. Middleware chain: logging → auth → handler, composto manualmente.

**Hints:**
- `http.NewServeMux()` e `mux.HandleFunc("/todos/", handler)`.
- `json.NewDecoder(r.Body).Decode(&req)` retorna error se body inválido.
- `type Middleware func(http.Handler) http.Handler`; compor: `handler = logging(auth(realHandler))`.
