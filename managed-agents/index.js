const fs = require('node:fs')
const path = require('node:path')
const crypto = require('node:crypto')
const readline = require('node:readline/promises')

const TODO_FILE = path.join(__dirname, 'todo.txt')
const SESSION_FILE = path.join(__dirname, 'session.jsonl')

function getTodo() {
    try {
        if (!fs.existsSync(TODO_FILE)) return { data: [], error: null }
        const todos = fs.readFileSync(TODO_FILE, 'utf8').split('\n').filter(Boolean)
        return { data: todos, error: null }
    } catch (err) {
        return { data: null, error: err.message }
    }
}

function setTodo(input) {
    try {
        if (!input?.todo) return { data: null, error: "Missing 'todo' in input" }
        const todo = input.todo.trim()

        const { data: todos, error } = getTodo()
        if (error) return { data: null, error }
        if (todos.some(t => t.toLowerCase() === todo.toLowerCase())) {
            return { data: `Already exists: ${todo}`, error: null }
        }

        fs.appendFileSync(TODO_FILE, todo + '\n')
        return { data: `Added: ${todo}`, error: null }
    } catch (err) {
        return { data: null, error: err.message }
    }
}

const tools = {
    'getTodo': {
        name: 'getTodo',
        description: 'List all todo items',
        parameters: { type: 'object', properties: {} },
        exec: getTodo,
    },
    'setTodo': {
        name: 'setTodo',
        description: 'Add a todo item',
        parameters: {
            type: 'object', required: ['todo'],
            properties: { todo: { type: 'string', description: 'The todo text' } }
        },
        exec: setTodo,
    }
}

function logSessions(sessionId, message){
    const event = { sessionId, ts: new Date().toISOString(), message }
    fs.appendFileSync(SESSION_FILE, JSON.stringify(event) + '\n')
}

const ollamaTools = Object.values(tools).map(({ name, description, parameters }) => ({
    type: 'function',
    function: { name, description, parameters },
}))

// hands: execute(name, input): string
function execute(name, input) {
    if (!Object.hasOwn(tools, name)) {
        return JSON.stringify({ data: null, error: `Unknown tool: ${name}` })
    }
    try {
        return JSON.stringify(tools[name].exec(input))
    } catch (err) {
        return JSON.stringify({ data: null, error: err.message })
    }
}


const SYSTEM_PROMPT = `You are a todo assistant. You manage the user's todo list using tools.
- To see todos, call getTodo.
- To add a todo, call setTodo with a short, clear todo text.
- Never invent todos; only report what getTodo returns.
- After using tools, reply to the user in one or two short sentences.`

const MAX_TURNS = 5

// brain
async function chat(messages) {
    const res = await fetch("http://localhost:11434/api/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            model: "qwen3:4b",
            messages,
            tools: ollamaTools,
            stream: false,
            think: true,
            options: { num_predict: 4096 },
        }),
        signal: AbortSignal.timeout(60_000),
    });

    if (!res.ok) throw new Error(`Ollama ${res.status}: ${await res.text()}`);

    const data = await res.json();
    return data.message
}

function createSession() {
    const session = { id: crypto.randomUUID(), messages: [] }
    addMessage(session, { role: "system", content: SYSTEM_PROMPT })
    return session
}

function loadSession(id) {
    if (!fs.existsSync(SESSION_FILE)) return null

    const messages = []
    for (const line of fs.readFileSync(SESSION_FILE, 'utf8').split('\n')) {
        if (!line.trim()) continue
        try {
            const event = JSON.parse(line)
            if (event.sessionId === id) messages.push(event.message)
        } catch {
            // skip malformed lines
        }
    }

    return messages.length ? { id, messages } : null
}

// what the brain sees, built from the full log
function buildContext(messages) {
    return messages.map(({ role, content, tool_calls, tool_name }) =>
        ({ role, content, tool_calls, tool_name }))
}

function addMessage(session, message) {
    session.messages.push(message)
    logSessions(session.id, message)
}

// tool calls from the last assistant message that have no result logged yet
function pendingToolCalls(messages) {
    const i = messages.findLastIndex(m => m.role === 'assistant')
    if (i === -1) return []
    const calls = messages[i].tool_calls ?? []
    const done = messages.length - 1 - i // tool results logged after it
    return calls.slice(done)
}

// harness loop: works out the next step purely from the log
async function run(session) {
    let turns = 0
    while (true) {
        for (const call of pendingToolCalls(session.messages)) {
            const { name, arguments: args } = call.function
            addMessage(session, { role: "tool", tool_name: name, content: execute(name, args) })
        }

        const last = session.messages.at(-1)
        if (last.role === 'assistant') return last.content // final answer, nothing pending

        if (turns++ >= MAX_TURNS) return `Stopped after ${MAX_TURNS} turns without a final answer.`
        addMessage(session, await chat(buildContext(session.messages)))
    }
}

async function harness(session, userText) {
    addMessage(session, { role: "user", content: userText })
    return run(session)
}

// replace a dead harness: load the log and carry on where it stopped
async function wake(sessionId) {
    const session = loadSession(sessionId)
    if (!session) throw new Error(`Session not found: ${sessionId}`)
    return { session, reply: await run(session) }
}

async function main() {
    const [command, sessionId] = process.argv.slice(2)

    let session
    if (command === 'resume') {
        if (!sessionId) return console.error('Usage: node index.js resume <session_id>')
        try {
            const woken = await wake(sessionId)
            session = woken.session
            console.log(`resumed session: ${session.id} (${session.messages.length} messages)`)
            console.log(woken.reply) // the finished interrupted run, or the last reply
        } catch (err) {
            return console.error(err.message)
        }
    } else if (!command) {
        session = createSession()
        console.log(`new session: ${session.id}`)
    } else {
        return console.error('Usage: node index.js [resume <session_id>]')
    }

    const rl = readline.createInterface({ input: process.stdin, output: process.stdout })
    console.log('Type a message, or "exit" to quit.')
    rl.setPrompt('> ')
    rl.prompt()

    for await (const line of rl) {
        const userText = line.trim()
        if (userText === 'exit') break

        if (userText) {
            try {
                console.log(await harness(session, userText))
            } catch (err) {
                console.error(err.message)
            }
        }
        rl.prompt()
    }

    rl.close()
    console.log(`To continue later: node index.js resume ${session.id}`)
}

main()