# Arquitetura — DevPath MCP

## Diagrama de Camadas

```
┌─────────────────────────────────────┐
│         MCP Client                  │
│  (Claude Desktop / Cursor / CLI)    │
└──────────────┬──────────────────────┘
               │ JSON-RPC 2.0
               ▼
┌─────────────────────────────────────┐
│         STDIO Transport             │
│     (github.com/mark3labs/mcp-go)   │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│         cmd/server/main.go          │
│   wire-up: storage → service →      │
│            tools → mcp server       │
└──────────────┬──────────────────────┘
               │
       ┌───────┴────────┐
       ▼                ▼
┌────────────┐   ┌──────────────┐
│  internal/ │   │  internal/   │
│  mcp/      │   │  tools/      │
│  server.go │   │  registry.go │
│  (factory) │   │  + handlers  │
└────────────┘   └──────┬───────┘
                        │ chama
                        ▼
               ┌─────────────────┐
               │ internal/       │
               │ service/        │
               │ service.go      │
               │ (lógica pura)   │
               └────────┬────────┘
                        │ lê
               ┌────────┴────────┐
               ▼                 ▼
     ┌──────────────┐   ┌──────────────┐
     │ internal/    │   │ internal/    │
     │ roadmap/     │   │ storage/     │
     │ golang.go    │   │ storage.go   │
     │ (dados)      │   │ (JSON I/O)   │
     └──────────────┘   └──────────────┘
```

## Responsabilidade de cada pacote

### `cmd/server/`
Ponto de entrada. Wire-up apenas: carrega progresso, cria service, registra tools, inicia servidor STDIO.
**NÃO contém**: lógica de negócio, acesso a dados.

### `internal/mcp/`
Factory que cria o `*server.MCPServer` configurado com nome, versão e instruções.
**NÃO contém**: ferramentas registradas (isso é responsabilidade de `tools/`).

### `internal/tools/`
Handlers MCP. Cada arquivo define um `NewXxxTool()` (schema) e `HandleXxx()` (handler).
**NÃO contém**: lógica de negócio — apenas extrai parâmetros, chama `service`, formata texto de resposta.

### `internal/service/`
Toda a lógica de negócio. Recebe dados e retorna respostas.
**NÃO contém**: I/O de arquivo, código MCP, formatação de texto para o usuário.

### `internal/roadmap/`
Dados hardcoded dos roadmaps. Registry de linguagens.
**NÃO contém**: lógica de progressão, persistência.

### `internal/storage/`
Leitura e escrita do `progress.json`. Sem lógica de negócio.
**NÃO contém**: conhecimento do roadmap, estrutura de response.

### `internal/models/`
Tipos compartilhados entre todos os pacotes. Sem lógica.

---

## Fluxo de uma Tool Call

```
1. Cliente envia JSON-RPC: {"method":"tools/call","params":{"name":"get_today_topic",...}}
2. mcp-go desserializa → chama handler registrado
3. Handler extrai params (language, time_available)
4. Handler chama service.GetTodayTopic(language, timeAvailable)
5. Service consulta roadmap.Get(language) → slice de tópicos
6. Service filtra tópicos: pula completados, pula com pré-requisitos não feitos
7. Service retorna *StudyPlanResponse
8. Handler formata resposta em texto legível
9. Handler retorna mcp.NewToolResultText(texto)
10. mcp-go serializa → envia JSON-RPC response via stdout
```

---

## Estratégia de Persistência

- Arquivo único: `data/progress.json`
- Carregado **uma vez** no startup em `cmd/server/main.go`
- Escrito **após cada `mark_topic_complete`** via closure `saveFn`
- **Sem mutex**: STDIO é serial por definição — o servidor processa exatamente uma tool call por vez, impossibilitando escritas concorrentes

Estrutura do arquivo:
```json
{
  "completed": {
    "golang": {
      "Variables": true,
      "Structs": true
    }
  }
}
```

---

## Algoritmo de "Weak Areas"

Weak areas são tópicos **não concluídos** que bloqueiam muitos outros tópicos downstream.

```
Para cada tópico T no roadmap:
    dependentCount[T.Name] = quantidade de outros tópicos que listam T como pré-requisito

weakAreas = [T para T no roadmap se:
    T não está completo
    E dependentCount[T.Name] >= 2]

Ordenar por dependentCount decrescente
Retornar top 3
```

Exemplo com Go roadmap: `Functions` é pré-requisito de 7 tópicos. Se não completado, aparece como weak area prioritária.

---

## Como adicionar uma nova linguagem

1. Criar `internal/roadmap/python.go` (ou outra linguagem)
2. Definir o slice `PythonRoadmap []models.Topic` com todos os tópicos
3. Registrar em `roadmap.Registry`:
   ```go
   var Registry = map[string][]models.Topic{
       "golang": GolangRoadmap,
       "python": PythonRoadmap,  // adicionar aqui
   }
   ```
4. Atualizar o `Enum` nos tools que aceitam `language` para incluir `"python"`
5. Nenhuma outra mudança necessária — toda a lógica de progressão é genérica

---

## Decisões de Design

**Por que texto formatado no output (não JSON)?**
Clientes MCP exibem `TextContent` direto ao usuário. JSON bruto seria ilegível. Os tipos de response em `models.go` existem para estrutura interna; o handler é responsável por formatar para leitura humana.

**Por que `saveFn` closure em vez de injeção de interface `Storage`?**
O service permanece puro e testável com `noopSave := func(...) error { return nil }`. Menos abstração desnecessária para o MVP.

**Por que `[]models.Topic` (slice ordenado) e não `map`?**
Ordem encoda progressão. O índice do slice serve como chave de ordenação natural. `GetNextTopic` varre para frente a partir da posição atual — O(n) sobre 18 elementos é negligível.

**Por que `time.Now().YearDay() % len(challenges)` para rotação de desafios?**
Sem banco, sem contadores, sem aleatoriedade. Mesmo tópico gera mesmo desafio no mesmo dia (reproduzível), rotaciona naturalmente dia a dia, incentiva o usuário a voltar amanhã.
