# Contribuindo com DevPath MCP

---

## Como adicionar um novo tópico ao roadmap Golang

Edite `internal/roadmap/golang.go`. Cada tópico segue a struct `models.Topic`:

```go
{
    Name:             "Generics",           // nome exato (case-insensitive nas tools)
    Difficulty:       "advanced",           // "beginner" | "intermediate" | "advanced"
    EstimatedMinutes: 50,
    Prerequisites:    []string{"Interfaces", "Functions"},
    PracticeChallenges: []string{
        "Write a generic Stack[T any] with Push, Pop, and Peek methods.",
        "Implement a generic Map[K comparable, V any] wrapper with Get and Set.",
        "Write a generic Filter[T any](slice []T, pred func(T) bool) []T function.",
    },
    Hints: []string{
        "type Stack[T any] struct { items []T } — T is the type parameter.",
        "comparable constraint allows == comparisons; any allows any type.",
        "func Filter[T any](s []T, f func(T) bool) []T — range and append only if f(v).",
    },
},
```

**Regras:**
- Inserir na posição correta do slice (ordem encoda progressão).
- `Prerequisites` devem referenciar `Name` exato de tópicos existentes.
- Ao menos 3 `PracticeChallenges` e 3 `Hints` correspondentes (mesmo índice).
- `EstimatedMinutes` realista — usado para checar se cabe no `time_available`.

---

## Como adicionar uma nova linguagem

### 1. Criar o arquivo de roadmap

`internal/roadmap/python.go`:

```go
package roadmap

import "github.com/andre-0303/devpath-mcp/internal/models"

var PythonRoadmap = []models.Topic{
    {
        Name:             "Variables",
        Difficulty:       "beginner",
        EstimatedMinutes: 15,
        Prerequisites:    []string{},
        PracticeChallenges: []string{
            "Declare 5 variables of different types and use type() to inspect them.",
        },
        Hints: []string{
            "Python is dynamically typed — no explicit type declaration needed.",
        },
    },
    // ... demais tópicos
}
```

### 2. Registrar no Registry

`internal/roadmap/golang.go` (ou criar `internal/roadmap/registry.go`):

```go
var Registry = map[string][]models.Topic{
    "golang": GolangRoadmap,
    "python": PythonRoadmap,  // adicionar aqui
}
```

### 3. Atualizar o Enum nas tools

Em cada arquivo de tool que aceita `language` como parâmetro, atualizar o campo `Enum`:

```go
// internal/tools/get_today_topic.go
mcp.WithString("language",
    mcp.Required(),
    mcp.Description("Programming language to study"),
    mcp.Enum("golang", "python"),  // adicionar aqui
),
```

Arquivos a atualizar:
- `internal/tools/get_today_topic.go`
- `internal/tools/get_next_topic.go`
- `internal/tools/mark_topic_complete.go`
- `internal/tools/review_progress.go`

`generate_practice` aceita `language` como opcional (default `"golang"`) — atualizar também.

### 4. Nenhuma outra mudança

Toda a lógica de progressão em `internal/service/service.go` é genérica — funciona para qualquer linguagem registrada.

---

## Como adicionar uma nova ferramenta MCP

### 1. Criar o handler

`internal/tools/my_tool.go`:

```go
package tools

import (
    "context"
    "fmt"

    "github.com/andre-0303/devpath-mcp/internal/service"
    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

// NewMyTool defines the MCP tool schema.
func NewMyTool() mcp.Tool {
    return mcp.NewTool("my_tool",
        mcp.WithDescription("Short description of what this tool does."),
        mcp.WithString("language",
            mcp.Required(),
            mcp.Description("Programming language"),
            mcp.Enum("golang"),
        ),
        mcp.WithString("my_param",
            mcp.Required(),
            mcp.Description("What this param does"),
        ),
    )
}

// HandleMyTool is the handler — extract params, call service, format text response.
func HandleMyTool(svc *service.Service) server.ToolHandlerFunc {
    return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        language, err := req.RequireString("language")
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }
        myParam, err := req.RequireString("my_param")
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }

        result, err := svc.MyServiceMethod(language, myParam)
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }

        text := fmt.Sprintf("Result: %s", result)
        return mcp.NewToolResultText(text), nil
    }
}
```

**Regras do handler:**
- Só extrai parâmetros, chama service, formata texto. Sem lógica de negócio.
- Erros retornam `mcp.NewToolResultError(msg)` — nunca `error` Go como segundo retorno (isso indica falha do protocolo, não do usuário).
- Output é texto formatado legível, não JSON.

### 2. Adicionar lógica no service

`internal/service/service.go`:

```go
func (s *Service) MyServiceMethod(language string, param string) (string, error) {
    rm := roadmap.Get(language)
    if rm == nil {
        return "", fmt.Errorf("language %q not found. Supported: %s",
            language, strings.Join(roadmap.SupportedLanguages(), ", "))
    }
    // lógica pura aqui — sem I/O, sem MCP, sem formatação de texto para o usuário
    return "resultado", nil
}
```

### 3. Registrar a tool

`internal/tools/registry.go`:

```go
func RegisterAll(s *server.MCPServer, svc *service.Service) {
    s.AddTool(NewGetTodayTopicTool(), HandleGetTodayTopic(svc))
    s.AddTool(NewGetNextTopicTool(), HandleGetNextTopic(svc))
    s.AddTool(NewGeneratePracticeTool(), HandleGeneratePractice(svc))
    s.AddTool(NewMarkTopicCompleteTool(), HandleMarkTopicComplete(svc))
    s.AddTool(NewReviewProgressTool(), HandleReviewProgress(svc))
    s.AddTool(NewMyTool(), HandleMyTool(svc))  // adicionar aqui
}
```

### 4. Documentar em docs/API.md

Seguir o formato existente: parâmetros, output esperado, erros possíveis, exemplo JSON.

---

## Convenções de código

| Área | Convenção |
|------|-----------|
| Nomes de tópicos | PascalCase canônico no roadmap; case-insensitive nas tools |
| Erros | Sempre `fmt.Errorf("context: %w", err)` para wrap; mensagens em inglês (output ao usuário) |
| Output das tools | Texto formatado legível; não JSON |
| Receivers | Pointer receiver para mutação; value receiver para leitura |
| Logs | Sempre para `stderr` — qualquer byte no `stdout` corrompe o transporte MCP |

---

## Rodando após mudanças

```powershell
go build -o devpath-mcp.exe ./cmd/server
go vet ./...
```

Reinicie o cliente MCP após rebuild — o servidor é um processo separado e não recarrega automaticamente.

---

## Testando a lógica de negócio

O service é testável sem MCP — usar `noopSave` como closure:

```go
package service_test

import (
    "testing"

    "github.com/andre-0303/devpath-mcp/internal/models"
    "github.com/andre-0303/devpath-mcp/internal/service"
)

func newTestService() *service.Service {
    progress := &models.UserProgress{
        Completed: map[string]map[string]bool{},
    }
    noopSave := func(*models.UserProgress) error { return nil }
    return service.New(progress, noopSave)
}

func TestGetTodayTopic_FirstTopic(t *testing.T) {
    svc := newTestService()
    resp, err := svc.GetTodayTopic("golang", 60)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp.Topic.Name != "Variables" {
        t.Errorf("expected Variables, got %q", resp.Topic.Name)
    }
}

func TestGetTodayTopic_SkipsCompleted(t *testing.T) {
    svc := newTestService()
    _ = svc.MarkTopicComplete("golang", "Variables")

    resp, err := svc.GetTodayTopic("golang", 60)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if resp.Topic.Name == "Variables" {
        t.Error("should have skipped Variables")
    }
}
```

```powershell
go test ./internal/service/...
go test -race ./...
```
