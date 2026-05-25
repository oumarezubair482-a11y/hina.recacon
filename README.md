# 🌸 Hina MD — WhatsApp Bot

> WhatsApp Bot built with Go + whatsmeow

---

## ⚡ Commands

| Command | Description |
|---------|-------------|
| `.menu` | Show full menu |
| `.ping` | Check bot speed |
| `.id` | Get your WhatsApp ID |
| `.info` | Bot information |
| `.time` | Current date & time |
| `.gc` | Group info (groups only) |
| `.tagall` | Tag all members (groups only) |

---

## 🚀 Deploy on Railway

### Step 1 — GitHub pe upload karo
```bash
git init
git add .
git commit -m "Hina MD Bot"
git remote add origin YOUR_GITHUB_REPO_URL
git push -u origin main
```

### Step 2 — Railway setup
1. [railway.app](https://railway.app) pe jao
2. **New Project → Deploy from GitHub repo** select karo
3. Apna repo select karo
4. Railway auto-detect karega Dockerfile
5. Deploy ho jayega ✅

### Step 3 — Bot pair karo
Deploy hone ke baad Railway aapko ek URL dega, jaise:
```
https://hinamd-xxxx.railway.app
```

Browser mein yeh open karo:
```
https://hinamd-xxxx.railway.app/pair/923001234567
```
- Number mein `+` ya spaces mat dalo
- Pakistan number example: `923001234567`

Yeh aapko ek **8-digit code** dega jaise `ABCD-EFGH`

### Step 4 — WhatsApp mein link karo
1. WhatsApp open karo
2. **Settings → Linked Devices → Link a Device**
3. **"Link with phone number instead"** pe tap karo
4. Code enter karo — Done! ✅

---

## 💻 Local Run (Test ke liye)

```bash
# Dependencies install
go mod tidy

# Run
go run .
```

Phir browser mein:
```
http://localhost:8080/pair/YOURNUMBER
```

---

## 📁 File Structure

```
hinamd/
├── main.go        # Main bot code
├── go.mod         # Go dependencies
├── Dockerfile     # Railway deployment
└── README.md      # Ye file
```

---

## ℹ️ Status Check

```
https://hinamd-xxxx.railway.app/
```
Returns:
```json
{"bot":"Hina MD","version":"v1.0","status":"online"}
```

---

_🌸 Hina MD — Made with ❤️ in Go_
