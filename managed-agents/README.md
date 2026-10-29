# managed-agents

A minimal todo agent in plain Node.js, built to learn the ideas in Anthropic's [Managed Agents](https://www.anthropic.com/research/managed-agents) post, using a local Ollama model.

## Goal

Keep the agent's parts separate, so any of them can be replaced:

- **Brain**: `chat()` calls the model through Ollama's `/api/chat`.
- **Hands**: `execute(name, input) → string` runs the tools (`getTodo`, `setTodo`) and never throws.
- **Session**: `session.jsonl` is an append-only log of every message, tagged with a `sessionId`.
- **Harness**: `run()` decides the next step only from what's in the log, so a crashed process can be replaced with `wake(sessionId)`.

## Setup

You need Node 18+ and Docker.

```bash
docker run -d --gpus=all -v ollama:/root/.ollama -p 11434:11434 --name ollama ollama/ollama
docker exec -it ollama ollama pull qwen3:4b
```

Leave out `--gpus=all` if you don't have an NVIDIA GPU.

## Run

```bash
node index.js                     # start a new session (prints its id)
node index.js resume <session_id> # continue an existing session
```

Type a message at the `>` prompt. Type `exit` or press Ctrl+D to quit.

## Cases to try

| Case | How | Expected |
|---|---|---|
| Tool use | `add buy milk, then show my todos` | Calls `setTodo`, then `getTodo`, then replies |
| Duplicate todo | Add the same todo twice | The second time returns `Already exists: ...` |
| Memory within a session | Ask `what did I ask you first?` | Answers using the earlier messages |
| Resume | `exit`, then `node index.js resume <id>` and ask a follow-up | The conversation continues with full history |
| Crash recovery | Add `process.exit(1)` right after the tool-result `addMessage` in `run()`. Ask `add gym and groceries`. Remove the line, then `resume <id>` | Only the unfinished tool call runs, then the final answer |
| Unsupported action | `delete buy milk` | The model says it can't do that (there's no delete tool) |

## Notes

- `qwen3:4b` is a thinking-only model. Keep `think: true`, or its reasoning ends up in the reply. For replies without reasoning, use `qwen3:4b-instruct` with `think: false`.
- Don't use `temperature: 0` with `qwen3:4b`. Its reasoning tends to loop and never finish.
- `session.jsonl` and `todo.txt` are created at runtime and are gitignored.
