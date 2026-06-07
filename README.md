# conduit
![Alt text](/screenshots/Conduit_JE1_BE1.webp)

```bash
go install github.com/kikikian/conduit@latest
```

---

## the problem

polyglot projects make you duplicate constants. a port number lives in your Go server, your Python script, and your TypeScript frontend — three places to keep in sync, three places to get wrong.

conduit fixes that. define a variable once and it automatically appears in every language you use, with the right type.

---

## quick start

```bash
# 1. initialize your project
conduit init

# 2. add a variable
conduit add --name PORT --type int --value 3000

# 3. start watching
conduit watch
```

that's it. conduit monitors `.conduit` and rewrites your target files every time a variable changes.

---

## how it works

mark the spot in your code where you want a variable injected:

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

when conduit runs, those lines become real typed declarations:

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

## commands

**`conduit init`**
interactive setup. select your target languages and map them to file paths. creates `conduit.config.json`.

**`conduit add`**
add a variable. fully interactive or via flags:
```bash
conduit add --name DATABASE_URL --type string --value postgresql://localhost/mydb
conduit add --name DEBUG --type bool --value true
conduit add --name PORT --type int --value 5432
```

**`conduit watch`**
watch `.conduit` for changes and regenerate all target files instantly. run this in the background while you develop.

```bash
conduit watch
conduit watch --file python:./app.py   # one-off override
```

---

## types

| conduit | python | typescript | go |
|---------|--------|------------|----|
| `int` | `int` | `number` | `int` |
| `string` | `str` | `string` | `string` |
| `bool` | `bool` | `boolean` | `bool` |

---

## config

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

commit this file — your whole team shares the same setup.

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

requires go 1.21+