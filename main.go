package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
	_ "github.com/mattn/go-sqlite3"
)

var (
	client    *whatsmeow.Client
	container *sqlstore.Container
)

const (
	BOT_NAME    = "Hina MD"
	BOT_VERSION = "v1.0"
	BOT_PREFIX  = "."
	OWNER_NAME  = "Hina"
)

var LOGO = strings.Join([]string{
	"",
	"╔══════════════════════════╗",
	"║   🌸  H I N A  M D  🌸  ║",
	"║     WhatsApp Bot v1.0    ║",
	"╚══════════════════════════╝",
	"",
}, "\n")

// ── Menu ─────────────────────────────────────────────────────────────────────

func menuText() string {
	return strings.Join([]string{
		LOGO,
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"🤖 *" + BOT_NAME + "*  |  " + BOT_VERSION,
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"",
		"👑 *Owner:*   " + OWNER_NAME,
		"⚡ *Prefix:*  [ " + BOT_PREFIX + " ]",
		"🕐 *Time:*    " + time.Now().Format("02 Jan 2006 | 15:04:05"),
		"",
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"📋 *GENERAL COMMANDS*",
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"",
		"  🌸 *.menu*    ➤ Yeh menu dikhao",
		"  🏓 *.ping*    ➤ Bot speed check",
		"  🪪 *.id*      ➤ Apna WhatsApp ID",
		"  🤖 *.info*    ➤ Bot ki info",
		"  🕐 *.time*    ➤ Date aur time",
		"  💚 *.alive*   ➤ Bot online check",
		"  👑 *.owner*   ➤ Owner ki info",
		"  ⚡ *.speed*   ➤ Speed test",
		"",
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"😍 *REACTION COMMANDS*",
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"",
		"  💬 *.react* 😍  ➤ Kisi msg pe react karo",
		"                  _(reply karke use karo)_",
		"",
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"👥 *GROUP COMMANDS*",
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"",
		"  👥 *.gc*       ➤ Group ki info",
		"  📢 *.tagall*   ➤ Sab members tag karo",
		"",
		"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		"",
		"_⚡ Har command pe auto-reaction hoti hai!_",
		"",
		"> 💬 _Powered by *" + BOT_NAME + "* • Made with ❤️_",
		"",
	}, "\n")
}

// ── Send Helpers ──────────────────────────────────────────────────────────────

func sendText(chat types.JID, text string) {
	if client == nil || !client.IsConnected() {
		return
	}
	_, _ = client.SendMessage(context.Background(), chat, &waProto.Message{
		Conversation: proto.String(text),
	})
}

func sendReply(chat types.JID, stanzaID, participant, text string) {
	if client == nil || !client.IsConnected() {
		return
	}
	_, _ = client.SendMessage(context.Background(), chat, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:    proto.String(stanzaID),
				Participant: proto.String(participant),
				QuotedMessage: &waProto.Message{
					Conversation: proto.String(""),
				},
			},
		},
	})
}

func sendReaction(chat types.JID, msgID, sender, emoji string) {
	if client == nil || !client.IsConnected() {
		return
	}
	_, _ = client.SendMessage(context.Background(), chat, &waProto.Message{
		ReactionMessage: &waProto.ReactionMessage{
			Key: &waProto.MessageKey{
				RemoteJid:   proto.String(chat.String()),
				FromMe:      proto.Bool(false),
				Id:          proto.String(msgID),
				Participant: proto.String(sender),
			},
			Text:              proto.String(emoji),
			SenderTimestampMs: proto.Int64(time.Now().UnixMilli()),
		},
	})
}

// ── Command Handler ───────────────────────────────────────────────────────────

func handleMessage(evt *events.Message) {
	if evt.Info.IsFromMe {
		return
	}

	text := evt.Message.GetConversation()
	if text == "" && evt.Message.ExtendedTextMessage != nil {
		text = evt.Message.ExtendedTextMessage.GetText()
	}
	text = strings.TrimSpace(text)

	if !strings.HasPrefix(text, BOT_PREFIX) {
		return
	}

	parts  := strings.Fields(text)
	cmd    := strings.ToLower(strings.TrimPrefix(parts[0], BOT_PREFIX))
	chat   := evt.Info.Chat
	sender := evt.Info.Sender.ToNonAD().String()
	msgID  := evt.Info.ID

	fmt.Printf("CMD [%s] from %s\n", cmd, sender)

	// ✅ Har command pe auto reaction
	go sendReaction(chat, msgID, sender, "⚡")

	switch cmd {

	case "menu", "help":
		go sendReaction(chat, msgID, sender, "🌸")
		sendText(chat, menuText())

	case "ping":
		start := time.Now()
		go sendReaction(chat, msgID, sender, "🏓")
		sendText(chat, "🏓 *Pong!*\n⚡ *Speed:* "+time.Since(start).String()+"\n✅ Bot is alive!")

	case "id":
		go sendReaction(chat, msgID, sender, "🪪")
		sendReply(chat, msgID, sender,
			"👤 *Your ID:*\n`"+sender+"`\n\n💬 *Chat ID:*\n`"+chat.ToNonAD().String()+"`")

	case "info":
		go sendReaction(chat, msgID, sender, "🤖")
		sendText(chat, strings.Join([]string{
			LOGO,
			"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
			"🤖 *Bot Information*",
			"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
			"",
			"*Name:*    " + BOT_NAME,
			"*Version:* " + BOT_VERSION,
			"*Owner:*   " + OWNER_NAME,
			"*Prefix:*  " + BOT_PREFIX,
			"*Library:* whatsmeow (Go)",
			"*Status:*  Online ✅",
			"",
			"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		}, "\n"))

	case "time", "date":
		go sendReaction(chat, msgID, sender, "🕐")
		now := time.Now()
		sendText(chat, strings.Join([]string{
			"🕐 *Date & Time*",
			"",
			"📅 *Date:* " + now.Format("Monday, 02 January 2006"),
			"⏰ *Time:* " + now.Format("15:04:05"),
			"🌍 *Zone:* " + now.Format("MST"),
		}, "\n"))

	case "alive":
		go sendReaction(chat, msgID, sender, "💚")
		sendText(chat, strings.Join([]string{
			LOGO,
			"💚 *Yes! I am Alive!*",
			"",
			"🤖 *" + BOT_NAME + "* is running perfectly!",
			"⚡ *Status:* Online ✅",
		}, "\n"))

	case "owner":
		go sendReaction(chat, msgID, sender, "👑")
		sendText(chat, strings.Join([]string{
			"👑 *Bot Owner*",
			"",
			"*Name:* " + OWNER_NAME,
			"*Bot:*  " + BOT_NAME + " " + BOT_VERSION,
			"",
			"_Support ke liye owner se rabta karein._",
		}, "\n"))

	case "speed", "test":
		t1 := time.Now()
		go sendReaction(chat, msgID, sender, "⚡")
		sendText(chat, "⚡ *Speed Test*\n\n✅ Response: "+time.Since(t1).String())

	case "react":
		if len(parts) < 2 {
			sendText(chat, "❓ *Usage:* Kisi msg pe reply karo aur likho:\n*.react* 😍")
			return
		}
		emoji := parts[1]
		var targetID, targetSender string
		if evt.Message.ExtendedTextMessage != nil && evt.Message.ExtendedTextMessage.ContextInfo != nil {
			ci := evt.Message.ExtendedTextMessage.ContextInfo
			if ci.StanzaId != nil {
				targetID = *ci.StanzaId
			}
			if ci.Participant != nil {
				targetSender = *ci.Participant
			}
		}
		if targetID == "" {
			sendText(chat, "❓ Pehle kisi message pe *reply* karo, phir *.react 😍* likho.")
			return
		}
		sendReaction(chat, targetID, targetSender, emoji)

	case "gc", "groupinfo":
		if chat.Server != "g.us" {
			go sendReaction(chat, msgID, sender, "❌")
			sendText(chat, "❌ Yeh command sirf *groups* mein use hoti hai!")
			return
		}
		groupInfo, err := client.GetGroupInfo(chat)
		if err != nil {
			go sendReaction(chat, msgID, sender, "❌")
			sendText(chat, "❌ Group info fetch nahi hui.")
			return
		}
		go sendReaction(chat, msgID, sender, "👥")
		sendText(chat, strings.Join([]string{
			"👥 *Group Information*",
			"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
			"",
			"*Name:*    " + groupInfo.Name,
			"*ID:*      " + chat.ToNonAD().String(),
			"*Members:* " + fmt.Sprintf("%d", len(groupInfo.Participants)),
			"*Created:* " + groupInfo.GroupCreated.Format("02 Jan 2006"),
			"",
			"*━━━━━━━━━━━━━━━━━━━━━━━━━━*",
		}, "\n"))

	case "tagall", "everyone", "all":
		if chat.Server != "g.us" {
			go sendReaction(chat, msgID, sender, "❌")
			sendText(chat, "❌ Yeh command sirf *groups* mein use hoti hai!")
			return
		}
		groupInfo, err := client.GetGroupInfo(chat)
		if err != nil {
			go sendReaction(chat, msgID, sender, "❌")
			sendText(chat, "❌ Members list nahi mili.")
			return
		}
		go sendReaction(chat, msgID, sender, "📢")
		var sb strings.Builder
		sb.WriteString("📢 *Tag All — " + BOT_NAME + "*\n\n")
		for _, p := range groupInfo.Participants {
			sb.WriteString("@" + p.JID.User + "\n")
		}
		sendText(chat, sb.String())

	default:
		go sendReaction(chat, msgID, sender, "❓")
		sendText(chat, "❓ *Unknown:* *"+BOT_PREFIX+cmd+"*\n\n_Type_ *.menu* _to see all commands._")
	}
}

// ── Event Handler ─────────────────────────────────────────────────────────────

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go handleMessage(v)
	case *events.Connected:
		fmt.Println("✅ WhatsApp Connected!")
	case *events.LoggedOut:
		fmt.Println("❌ Logged out!")
	case *events.Disconnected:
		fmt.Println("⚠️  Disconnected! Reconnecting in 5s...")
		go func() {
			time.Sleep(5 * time.Second)
			if client != nil {
				_ = client.Connect()
			}
		}()
	}
}

// ── HTTP Endpoints ────────────────────────────────────────────────────────────

func handlePair(w http.ResponseWriter, r *http.Request) {
	number := strings.NewReplacer("+", "", " ", "", "-", "").Replace(
		strings.TrimPrefix(r.URL.Path, "/pair/"))
	if len(number) < 10 {
		http.Error(w, `{"error":"Invalid number. Use: /pair/923001234567"}`, 400)
		return
	}
	if container == nil {
		http.Error(w, `{"error":"DB not ready, retry"}`, 500)
		return
	}
	if client != nil && client.IsConnected() {
		client.Disconnect()
		time.Sleep(2 * time.Second)
	}
	tmp := whatsmeow.NewClient(container.NewDevice(), waLog.Stdout("Pair", "ERROR", false))
	tmp.AddEventHandler(eventHandler)
	if err := tmp.Connect(); err != nil {
		http.Error(w, `{"error":"Connect failed"}`, 500)
		return
	}
	time.Sleep(3 * time.Second)

	code, err := tmp.PairPhone(context.Background(), number, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		tmp.Disconnect()
		http.Error(w, `{"error":"Pair failed: `+err.Error()+`"}`, 500)
		return
	}
	go func() {
		for i := 0; i < 60; i++ {
			time.Sleep(time.Second)
			if tmp.Store.ID != nil {
				fmt.Println("✅ Paired: " + number)
				client = tmp
				return
			}
		}
		fmt.Println("⚠️  Pair timeout")
		tmp.Disconnect()
	}()
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"success":true,"code":"%s","number":"%s"}`, code, number)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	status := "offline"
	if client != nil && client.IsConnected() && client.IsLoggedIn() {
		status = "online"
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"bot":"%s","version":"%s","status":"%s"}`, BOT_NAME, BOT_VERSION, status)
}

// ── Main ──────────────────────────────────────────────────────────────────────

func main() {
	fmt.Println(LOGO)
	fmt.Println("Starting " + BOT_NAME + " " + BOT_VERSION + "...")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/pair/", handlePair)

	go func() {
		fmt.Println("🌐 HTTP: 0.0.0.0:" + port)
		_ = http.ListenAndServe("0.0.0.0:"+port, nil)
	}()

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "file:hinamd.db?_foreign_keys=on"
	}

	var err error
	container, err = sqlstore.New(context.Background(), "sqlite3", dbPath, waLog.Stdout("DB", "ERROR", false))
	if err != nil {
		fmt.Println("❌ DB Error: " + err.Error())
		os.Exit(1)
	}

	dev, err := container.GetFirstDevice(context.Background())
	if err == nil {
		client = whatsmeow.NewClient(dev, waLog.Stdout("WA", "INFO", true))
		client.AddEventHandler(eventHandler)
		if client.Store.ID != nil {
			if err := client.Connect(); err == nil {
				fmt.Println("✅ Session restored!")
			}
		} else {
			fmt.Println("⚠️  No session. Pair via: /pair/YOURNUMBER")
		}
	}

	fmt.Println("✅ " + BOT_NAME + " is running!")
	fmt.Println("📱 Pair: http://localhost:" + port + "/pair/YOURNUMBER")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	fmt.Println("Shutting down...")
	if client != nil {
		client.Disconnect()
	}
}
