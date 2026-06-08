# conduit
![conduit](/screenshots/Conduit_JE1_BE1.webp)

---
> Conduit (noun.): a channel or pipe used to convey something, such as water, electrical wires, or information

Conduit is a channel for your variables across languages.

```bash
go install github.com/kikikian/conduit@latest
```

## The problem

Polyglot projects make you duplicate constants. A port number lives in your Go server, your Python script, and your TypeScript frontend. That isthree places to keep in sync, three places to get wrong.

Conduit fixes that. define a variable once and it automatically appears in every language you use, with the right type.

---

## Supported Languages (for now)
-> Golang

-> Python

-> Typesript

---

## Quick start

```bash
# 1. initialize your project
conduit init

# 2. add a variable
conduit add --name PORT --type int --value 3000

# 3. start watching
conduit watch
```

Conduit monitors `.conduit` and rewrites your target files every time a variable changes.

---

## How it works

Mark the spot in your code where you want a variable injected:

**python**
```python
# conduit:import PORT
```

**typescript**
```typescript
// conduit:import PORT
```

**go**
```go
// conduit:import PORT
```

When conduit runs, those lines become typed declarations:

```python
PORT: int = 3000
```
```typescript
const PORT: number = 3000
```
```go
var PORT int = 3000
```

---

## Commands

**`conduit init`**
Select your target languages and map them to file paths -> creates `conduit.config.json`.

**`conduit add`**
Add a variable. Interactive or via flags:
```bash
conduit add --name DATABASE_URL --type string --value postgresql://localhost/mydb
conduit add --name DEBUG --type bool --value true
conduit add --name PORT --type int --value 5432
```

**`conduit watch`**
Watch `.conduit` for changes and regenerate all target files instantly. Run this in the background while you develop.

```bash
conduit watch
conduit watch --file python:./app.py   # one-off override
```

---

## Types

| conduit | python | typescript | go |
|---------|--------|------------|----|
| `int` | `int` | `number` | `int` |
| `string` | `str` | `string` | `string` |
| `bool` | `bool` | `boolean` | `bool` |

---

## Config

`conduit init` generates a `conduit.config.json` in your project root:

```json
{
  "targets": [
    { "lang": "python", "filePath": "./app.py" },
    { "lang": "typescript", "filePath": "./src/config.ts" },
    { "lang": "go", "filePath": "./pkg/config.go" }
  ]
}
```

---

## project files

| file | purpose |
|------|---------|
| `.conduit` | stores all variable definitions |
| `conduit.config.json` | maps languages to target files |

---

## build from source

```bash
git clone https://github.com/kikikian/conduit.git
cd conduit
go build
```

### requires go 1.21+
