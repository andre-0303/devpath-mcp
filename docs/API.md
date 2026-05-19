# API — DevPath MCP Tools

Todas as ferramentas retornam **texto formatado** legível diretamente no cliente MCP.

---

## `get_today_topic`

Retorna o tópico ideal para estudar hoje com base no progresso atual e tempo disponível.

### Parâmetros

| Nome | Tipo | Obrigatório | Descrição |
|------|------|-------------|-----------|
| `language` | string | sim | Linguagem de programação. Valores: `"golang"` |
| `time_available` | number | sim | Minutos disponíveis para estudo (5–480) |

### Output

```
Today's Topic: Slices [beginner — 35 min]
Reason: You've completed Arrays. Slices are the natural next step.
Fits in your 60-minute window: YES

Practice suggestion: Build a function that removes duplicates from an integer slice.
```

Se o tópico não cabe no tempo disponível:
```
Today's Topic: Concurrency [advanced — 60 min]
Reason: You've completed Loops and Functions.
Fits in your 20-minute window: NO (needs 60 min — consider a longer session)

Practice suggestion: Create 3 goroutines that print numbers concurrently.
```

Se o roadmap está completo:
```
Congratulations! You've completed all 18 topics in the golang roadmap.
Consider reviewing weak areas or starting a new language.
```

### Erros

- `"language 'xyz' not found. Supported: golang"` — linguagem não reconhecida
- `"time_available must be between 5 and 480 minutes"` — valor fora do range

### Exemplo

```json
{
  "name": "get_today_topic",
  "arguments": {
    "language": "golang",
    "time_available": 45
  }
}
```

---

## `get_next_topic`

Retorna o próximo tópico no roadmap após o tópico atual.

### Parâmetros

| Nome | Tipo | Obrigatório | Descrição |
|------|------|-------------|-----------|
| `language` | string | sim | Linguagem de programação. Valores: `"golang"` |
| `current_topic` | string | sim | Nome do tópico atual (case-insensitive) |

### Output

```
Next topic after 'Structs' in golang:

  Methods [intermediate — 35 min]
  Prerequisites: Structs ✓

Start with: "Create methods on a BankAccount struct for deposit and withdrawal."
```

Se o roadmap está completo:
```
'REST API' is the last topic in the golang roadmap.
Roadmap complete! You've mastered all 18 topics.
```

### Erros

- `"topic 'xyz' not found in golang roadmap"` — tópico não existe
- `"language 'xyz' not found. Supported: golang"` — linguagem não reconhecida

### Exemplo

```json
{
  "name": "get_next_topic",
  "arguments": {
    "language": "golang",
    "current_topic": "structs"
  }
}
```

---

## `generate_practice`

Gera um desafio prático para fixação de um tópico. O desafio rotaciona por dia (mesmo tópico = mesmo desafio no mesmo dia calendário).

### Parâmetros

| Nome | Tipo | Obrigatório | Descrição |
|------|------|-------------|-----------|
| `topic` | string | sim | Nome do tópico (case-insensitive) |
| `language` | string | não | Linguagem (default: `"golang"`) |

### Output

```
Practice Challenge — Structs [intermediate]

Challenge: Create a BankAccount struct with fields for owner, balance,
and account number. Add methods for deposit, withdraw, and balance inquiry.
Include validation (no negative balance).

Hint: Start with the struct definition, then add methods one by one.
Use a constructor function NewBankAccount() to enforce initial state.
```

### Erros

- `"topic 'xyz' not found in golang roadmap"` — tópico não existe

### Exemplo

```json
{
  "name": "generate_practice",
  "arguments": {
    "topic": "structs"
  }
}
```

---

## `mark_topic_complete`

Marca um tópico como concluído e salva o progresso no arquivo JSON.

### Parâmetros

| Nome | Tipo | Obrigatório | Descrição |
|------|------|-------------|-----------|
| `language` | string | sim | Linguagem de programação. Valores: `"golang"` |
| `topic` | string | sim | Nome do tópico (case-insensitive) |

### Output

```
✓ Marked 'Structs' as complete for golang!

Progress: 9/18 topics completed (50.0%)
Next up: Methods [intermediate — 35 min]
```

Se chamado para tópico já marcado:
```
'Structs' was already marked as complete for golang.

Progress: 9/18 topics completed (50.0%)
Next up: Methods [intermediate — 35 min]
```

### Erros

- `"topic 'xyz' not found in golang roadmap"` — tópico não existe
- `"failed to save progress: ..."` — erro de I/O ao salvar arquivo

### Exemplo

```json
{
  "name": "mark_topic_complete",
  "arguments": {
    "language": "golang",
    "topic": "structs"
  }
}
```

---

## `review_progress`

Exibe a evolução completa do estudante para uma linguagem.

### Parâmetros

| Nome | Tipo | Obrigatório | Descrição |
|------|------|-------------|-----------|
| `language` | string | sim | Linguagem de programação. Valores: `"golang"` |

### Output

```
Progress Report — golang
════════════════════════════════════

Completed: 7/18 topics (38.9%)

Completed topics:
  ✓ Variables      ✓ Constants      ✓ Functions
  ✓ Conditionals   ✓ Loops          ✓ Arrays
  ✓ Slices

Next recommended:
  Maps [beginner — 30 min]
  Prerequisites: Variables ✓, Loops ✓

Weak areas (uncompleted blockers):
  1. Error Handling — blocks 4 downstream topics
  2. Structs        — blocks 3 downstream topics
  3. Packages       — blocks 2 downstream topics
```

Se nenhum tópico foi completado:
```
Progress Report — golang
════════════════════════════════════

Completed: 0/18 topics (0%)

No topics completed yet. Start with:
  Variables [beginner — 20 min]

Tip: Use get_today_topic to get a personalized study plan for today.
```

### Erros

- `"language 'xyz' not found. Supported: golang"` — linguagem não reconhecida

### Exemplo

```json
{
  "name": "review_progress",
  "arguments": {
    "language": "golang"
  }
}
```

---

## Notas gerais

- Todos os nomes de tópicos são **case-insensitive** na entrada; armazenados com capitalização canônica do roadmap
- O progresso é **por linguagem** — estudar Golang não afeta progresso futuro em Python
- A tool `generate_practice` usa `time.Now().YearDay() % len(challenges)` para rotação determinística — mesmo resultado para múltiplas chamadas no mesmo dia
