# DevPath MCP

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MCP](https://img.shields.io/badge/MCP-Protocol-6D28D9?style=for-the-badge&logoColor=white)
![JSON](https://img.shields.io/badge/JSON-Storage-000000?style=for-the-badge&logo=json&logoColor=white)
![STDIO](https://img.shields.io/badge/Transport-STDIO-gray?style=for-the-badge&logoColor=white)

Servidor MCP em Go que guia desenvolvedores pelo aprendizado progressivo de linguagens. Em vez de um curso genérico, responde: **"Hoje, baseado no seu progresso, você deve estudar _isso_ por 1h."**

## Pré-requisitos

- Go 1.21+
- Qualquer cliente MCP: Claude Desktop, Claude Code CLI, Cursor, ou `mcptools`

## Instalação

```powershell
git clone https://github.com/andre-0303/devpath-mcp
cd devpath-mcp
go build -o devpath-mcp.exe ./cmd/server
```

O binário `devpath-mcp.exe` é gerado na raiz do projeto.

---

## Testando com Claude Code CLI

Adicione ao arquivo de configuração do Claude Code (`%APPDATA%\Claude\claude_desktop_config.json` ou `.claude/mcp.json` no projeto):

```json
{
  "mcpServers": {
    "devpath": {
      "command": "C:\\caminho\\absoluto\\para\\devpath-mcp.exe",
      "args": []
    }
  }
}
```

Reinicie o Claude Code. Confirme que as ferramentas carregaram com `/mcp` no prompt.

Frases de teste:

```
"O que devo estudar em Golang hoje? Tenho 45 minutos."
"Gera um exercício prático para o tópico Constants."
"Marca o tópico Variables como concluído."
"Mostra meu progresso em Golang."
```

---

## Testando com Claude Desktop

Mesmo arquivo de config: `%APPDATA%\Claude\claude_desktop_config.json`.

Reinicie o Claude Desktop. As 5 ferramentas aparecem automaticamente no ícone de ferramentas.

---

## Testando sem cliente MCP (raw stdin/stdout)

O servidor fala JSON-RPC 2.0 via STDIO. Útil para debugar sem abrir nenhum cliente.

**1. Inicializar a sessão:**

```powershell
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"0.0.1"}}}' | .\devpath-mcp.exe
```

**2. Chamar uma ferramenta diretamente:**

```powershell
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"get_today_topic","arguments":{"language":"golang","time_available":60}}}' | .\devpath-mcp.exe
```

**3. Listar ferramentas disponíveis:**

```powershell
echo '{"jsonrpc":"2.0","id":3,"method":"tools/list","params":{}}' | .\devpath-mcp.exe
```

> Logs de erro vão para `stderr` — nunca para `stdout`, que seria corromperia o transporte MCP.

---

## Testando com mcptools

`mcptools` é uma CLI para inspecionar e chamar servidores MCP sem cliente gráfico.

```powershell
# Instalar
go install github.com/f/mcptools/cmd/mcptools@latest

# Listar ferramentas
mcptools tools .\devpath-mcp.exe

# Chamar uma ferramenta
mcptools call get_today_topic --params '{"language":"golang","time_available":60}' -- .\devpath-mcp.exe

# Modo interativo (REPL)
mcptools shell .\devpath-mcp.exe
```

---

## Ferramentas disponíveis

| Ferramenta | Parâmetros | Descrição |
|-----------|-----------|-----------|
| `get_today_topic` | `language`, `time_available` (min) | Tópico ideal para hoje com base no progresso |
| `get_next_topic` | `language`, `current_topic` | Próximo tópico no roadmap |
| `generate_practice` | `topic`, `language` (opt) | Desafio prático para fixação |
| `mark_topic_complete` | `language`, `topic` | Marca tópico concluído e salva progresso |
| `review_progress` | `language` | Evolução completa: feitos, próximo, áreas fracas |

Todos os nomes de tópicos são **case-insensitive**.

---

## Onde o progresso fica salvo

`data/progress.json` na raiz do projeto (criado automaticamente na primeira chamada a `mark_topic_complete`).

```json
{
  "completed": {
    "golang": {
      "Variables": true,
      "Constants": true
    }
  }
}
```

**Para resetar o progresso durante testes:**

```powershell
Remove-Item data\progress.json
```

---

## Roadmap Golang (18 tópicos)

```
Variables → Constants → Functions → Conditionals → Loops →
Arrays → Slices → Maps → Structs → Methods →
Interfaces → Pointers → Error Handling → Packages →
Goroutines → Channels → Concurrency Patterns → REST API
```

---

## Erros comuns ao testar

| Erro | Causa | Fix |
|------|-------|-----|
| Ferramentas não aparecem no cliente | Caminho do exe errado no config | Usar caminho absoluto com barras duplas `\\` |
| `language 'xyz' not found` | Linguagem não suportada | Só `"golang"` no MVP |
| `time_available must be between 5 and 480` | Valor fora do range | Passar valor entre 5 e 480 |
| `topic 'xyz' not found` | Tópico não existe no roadmap | Verificar nome exato na lista acima |
| Servidor não inicia | `data/` não existe | Criar com `mkdir data` |

---

## Estrutura do projeto

```
cmd/server/main.go          # Ponto de entrada, wire-up
internal/mcp/server.go      # Factory do servidor MCP
internal/tools/             # Handlers por ferramenta
internal/service/service.go # Toda a lógica de negócio
internal/roadmap/golang.go  # Dados do roadmap (hardcoded)
internal/storage/storage.go # Leitura/escrita do progress.json
internal/models/models.go   # Tipos compartilhados
data/progress.json          # Progresso do usuário (gerado em runtime)
```

---

## Roadmap de versões

| Versão | Escopo |
|--------|--------|
| V1 (atual) | Golang, persistência JSON, 5 ferramentas MCP |
| V2 | Persistência SQLite, múltiplos usuários |
| V3 | React, Next.js, TypeScript, Docker, System Design |
| V4 | Modo objetivo ("quero ser backend engineer") |
| V5 | Curadoria de conteúdo (docs, artigos, vídeos) |

---

## Documentação adicional

- [Arquitetura](docs/ARCHITECTURE.md)
- [API — ferramentas MCP](docs/API.md)
- [Roadmap Golang — referência completa](docs/ROADMAP.md)
- [Contribuindo — adicionar tópicos, linguagens e tools](docs/CONTRIBUTING.md)
