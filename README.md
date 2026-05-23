# 🌸 HINA MD — WhatsApp Bot (Railway Deploy)

## Railway pe Deploy karne ka tarika

### Step 1: GitHub pe Upload karo
```bash
git init
git add .
git commit -m "Hina MD Bot"
git remote add origin https://github.com/USERNAME/hina-md.git
git push -u origin main
```

### Step 2: Railway pe Deploy
1. railway.app pe jayen
2. "New Project" → "Deploy from GitHub"
3. hina-md repo select karein
4. Deploy ho jayega automatically ✅

### Step 3: WhatsApp Connect karo
Deploy hone ke baad Railway URL milega, browser mein yeh open karein:

```
https://AAPKI-RAILWAY-URL/pair/PHONE_NUMBER
```

**Example:**
```
https://hina-md.up.railway.app/pair/923001234567
```

Yahan se **Pairing Code** milega.

### Step 4: WhatsApp mein Code Enter karo
1. WhatsApp kholen
2. Settings → Linked Devices
3. "Link a Device" → "Link with phone number instead"
4. Code enter karein ✅

---

## Endpoints

| URL | Kaam |
|-----|------|
| `/` | Bot status |
| `/status` | JSON status |
| `/pair/NUMBER` | Pairing code lo |
| `/logout` | Session delete karo |

---

## Commands

| Command | Description |
|---------|-------------|
| .menu | Full menu |
| .ping | Speed check |
| .owner | Owner info |
| .alive | Bot status |
| .id | User/Chat ID |

---
_🌸 Powered by Hina MD v1.0_
