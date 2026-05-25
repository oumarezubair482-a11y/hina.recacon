const {
    default: makeWASocket,
    useMultiFileAuthState,
    DisconnectReason,
    fetchLatestBaileysVersion,
    makeInMemoryStore,
} = require('@whiskeysockets/baileys');
const pino = require('pino');
const { Boom } = require('@hapi/boom');
const fs = require('fs');
const chalk = require('chalk');

// ══════════════════════════════════════════════
//              HINA MD — LOGO
// ══════════════════════════════════════════════
function showLogo() {
    console.clear();
    console.log(chalk.magentaBright('╔═══════════════════════════════════════════════╗'));
    console.log(chalk.magentaBright('║') + chalk.bold.cyanBright(' ██╗  ██╗██╗███╗  ██╗ █████╗               ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('║') + chalk.bold.cyanBright(' ██║  ██║██║████╗ ██║██╔══██╗              ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('║') + chalk.bold.cyanBright(' ███████║██║██╔██╗██║███████║              ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('║') + chalk.bold.cyanBright(' ██╔══██║██║██║╚████║██╔══██║              ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('║') + chalk.bold.cyanBright(' ██║  ██║██║██║ ╚███║██║  ██║   MD         ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('║') + chalk.bold.cyanBright(' ╚═╝  ╚═╝╚═╝╚═╝  ╚══╝╚═╝  ╚═╝            ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('╠═══════════════════════════════════════════════╣'));
    console.log(chalk.magentaBright('║') + chalk.yellow('   ⚡  WhatsApp Multi Device Bot  ⚡           ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('║') + chalk.green('   ✅  Version  : 2.0.0                        ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('║') + chalk.blueBright('   🔧  Prefix   : .                            ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('║') + chalk.redBright('   👑  Owner    : Hina                         ') + chalk.magentaBright('║'));
    console.log(chalk.magentaBright('╚═══════════════════════════════════════════════╝'));
    console.log('');
}

showLogo();

// ══════════════════════════════════════════════
//              CONFIG
// ══════════════════════════════════════════════
const CONFIG = {
    botName:     'Hina MD',
    prefix:      '.',
    ownerNumber: '923001234567',  // ← Apna number daalo
    sessionDir:  './hina_session',
    version:     '2.0.0',
};

// ══════════════════════════════════════════════
//              STORE
// ══════════════════════════════════════════════
const store = makeInMemoryStore({
    logger: pino().child({ level: 'silent', stream: 'store' })
});

// ══════════════════════════════════════════════
//              HELPERS
// ══════════════════════════════════════════════
function getBody(msg) {
    const m = msg.message;
    if (!m) return '';
    return (
        m.conversation ||
        m.extendedTextMessage?.text ||
        m.imageMessage?.caption ||
        m.videoMessage?.caption || ''
    );
}

function getSender(msg) {
    return msg.key.participant || msg.key.remoteJid || '';
}

function isOwner(msg) {
    return getSender(msg).replace(/\D/g, '') === CONFIG.ownerNumber.replace(/\D/g, '');
}

function isGroup(jid) {
    return jid?.endsWith('@g.us');
}

// ══════════════════════════════════════════════
//              COMMANDS
// ══════════════════════════════════════════════
const CMD = {

    // .menu
    menu: async (sock, msg, from) => {
        const text = `
╔════════════════════════════════╗
║    ⚡ *HINA MD — MENU v2* ⚡   ║
╠════════════════════════════════╣
║  🌐 *GENERAL*                  ║
║  .menu    » Yeh menu           ║
║  .ping    » Speed check        ║
║  .alive   » Bot alive?         ║
║  .info    » Bot info           ║
║  .id      » IDs dekho          ║
╠════════════════════════════════╣
║  👥 *GROUP*                    ║
║  .tagall  » Sab tag karo       ║
║  .kick    » Member hatao       ║
║  .promote » Admin banao        ║
║  .demote  » Admin hatao        ║
╠════════════════════════════════╣
║  🎮 *FUN*                      ║
║  .quote   » Random quote       ║
║  .joke    » Random joke        ║
║  .fact    » Random fact        ║
╠════════════════════════════════╣
║  🔒 *OWNER ONLY*               ║
║  .restart » Bot restart        ║
╚════════════════════════════════╝
📌 *Prefix:* ${CONFIG.prefix}
🤖 *${CONFIG.botName} v${CONFIG.version}*`.trim();
        await sock.sendMessage(from, { text }, { quoted: msg });
    },

    // .ping
    ping: async (sock, msg, from) => {
        const t1 = Date.now();
        await sock.sendMessage(from, { text: '🏓 *Pong!*' }, { quoted: msg });
        const ms = Date.now() - t1;
        await sock.sendMessage(from, { text: `⚡ *Speed:* ${ms}ms\n✅ *Status:* Online` });
    },

    // .alive
    alive: async (sock, msg, from) => {
        await sock.sendMessage(from, {
            text: `✅ *${CONFIG.botName}* zinda hai!\n⚡ Version: *${CONFIG.version}*\n🕐 Uptime: *${Math.floor(process.uptime())}s*`
        }, { quoted: msg });
    },

    // .info
    info: async (sock, msg, from) => {
        await sock.sendMessage(from, {
            text: `🤖 *Bot Info*\n━━━━━━━━━━━━━━━\n📛 Name    : *${CONFIG.botName}*\n🔖 Version : *${CONFIG.version}*\n⚙️ Prefix  : *${CONFIG.prefix}*\n📚 Library : *Baileys MD*\n👑 Owner   : *Hina*\n⏱️ Uptime  : *${Math.floor(process.uptime())}s*`
        }, { quoted: msg });
    },

    // .id
    id: async (sock, msg, from) => {
        const sender = getSender(msg);
        let text = `👤 *User ID:*\n${sender}\n\n💬 *Chat ID:*\n${from}`;
        if (isGroup(from)) {
            const meta = await sock.groupMetadata(from);
            text += `\n\n👥 *Group:* ${meta.subject}`;
        }
        await sock.sendMessage(from, { text }, { quoted: msg });
    },

    // .tagall
    tagall: async (sock, msg, from) => {
        if (!isGroup(from)) return sock.sendMessage(from, { text: '❌ Sirf group mein!' });
        const meta = await sock.groupMetadata(from);
        let text = `📢 *Tag All — ${CONFIG.botName}*\n\n`;
        const mentions = [];
        for (const m of meta.participants) {
            text += `• @${m.id.split('@')[0]}\n`;
            mentions.push(m.id);
        }
        await sock.sendMessage(from, { text, mentions }, { quoted: msg });
    },

    // .kick
    kick: async (sock, msg, from) => {
        if (!isGroup(from)) return sock.sendMessage(from, { text: '❌ Sirf group mein!' });
        if (!isOwner(msg)) return sock.sendMessage(from, { text: '🔒 Sirf owner!' });
        const target = msg.message?.extendedTextMessage?.contextInfo?.participant;
        if (!target) return sock.sendMessage(from, { text: '❌ Kisi ko reply karo!' });
        await sock.groupParticipantsUpdate(from, [target], 'remove');
        await sock.sendMessage(from, { text: `✅ @${target.split('@')[0]} kick!`, mentions: [target] });
    },

    // .promote
    promote: async (sock, msg, from) => {
        if (!isGroup(from)) return sock.sendMessage(from, { text: '❌ Sirf group mein!' });
        if (!isOwner(msg)) return sock.sendMessage(from, { text: '🔒 Sirf owner!' });
        const target = msg.message?.extendedTextMessage?.contextInfo?.participant;
        if (!target) return sock.sendMessage(from, { text: '❌ Kisi ko reply karo!' });
        await sock.groupParticipantsUpdate(from, [target], 'promote');
        await sock.sendMessage(from, { text: `⬆️ @${target.split('@')[0]} admin ban gaya!`, mentions: [target] });
    },

    // .demote
    demote: async (sock, msg, from) => {
        if (!isGroup(from)) return sock.sendMessage(from, { text: '❌ Sirf group mein!' });
        if (!isOwner(msg)) return sock.sendMessage(from, { text: '🔒 Sirf owner!' });
        const target = msg.message?.extendedTextMessage?.contextInfo?.participant;
        if (!target) return sock.sendMessage(from, { text: '❌ Kisi ko reply karo!' });
        await sock.groupParticipantsUpdate(from, [target], 'demote');
        await sock.sendMessage(from, { text: `⬇️ @${target.split('@')[0]} admin nahi raha!`, mentions: [target] });
    },

    // .quote
    quote: async (sock, msg, from) => {
        const list = [
            "Mushkilein insaan ko mazboot banati hain. 💪",
            "Har raat ke baad subah zaroor aati hai. 🌅",
            "Haar mat mano, koshish karo. 🔥",
            "Zindagi ek safar hai, enjoy karo. ✨",
            "Ummeed mat choro, waqt badalta hai. 🌟",
            "Kamyabi unhe milti hai jo koshish karte hain. 🏆",
        ];
        const q = list[Math.floor(Math.random() * list.length)];
        await sock.sendMessage(from, { text: `💬 *Quote*\n\n_"${q}"_` }, { quoted: msg });
    },

    // .joke
    joke: async (sock, msg, from) => {
        const list = [
            "Teacher: 1+1?\nStudent: Depends on situation sir! 😅",
            "Doctor se kaha neend nahi aati...\nDoctor: Meri fees dene ke baad aayegi! 😂",
            "Biwi: Mujhse pyaar?\nShohar: WiFi ki tarah connected rehta hun! 📶",
        ];
        const j = list[Math.floor(Math.random() * list.length)];
        await sock.sendMessage(from, { text: `😂 *Joke*\n\n${j}` }, { quoted: msg });
    },

    // .fact
    fact: async (sock, msg, from) => {
        const list = [
            "Octopus ke 3 dil hote hain! 🐙",
            "Shehad 3000 saal baad bhi kha sakte hain! 🍯",
            "Insaan ka dimagh 75% paani se bana hai. 🧠",
            "Butterflies pairon se taste karti hain! 🦋",
            "Dil ek din mein 100,000 baar dhakta hai! ❤️",
        ];
        const f = list[Math.floor(Math.random() * list.length)];
        await sock.sendMessage(from, { text: `🔍 *Fact*\n\n${f}` }, { quoted: msg });
    },

    // .restart
    restart: async (sock, msg, from) => {
        if (!isOwner(msg)) return sock.sendMessage(from, { text: '🔒 Sirf owner!' });
        await sock.sendMessage(from, { text: `🔄 *${CONFIG.botName} restart ho raha hai...* ⚡` });
        process.exit(0);
    },
};

// ══════════════════════════════════════════════
//              START BOT
// ══════════════════════════════════════════════
async function startBot() {
    const { state, saveCreds } = await useMultiFileAuthState(CONFIG.sessionDir);
    const { version } = await fetchLatestBaileysVersion();

    console.log(chalk.cyan(`📡 Baileys: ${version.join('.')}`));

    const sock = makeWASocket({
        version,
        logger: pino({ level: 'silent' }),
        printQRInTerminal: true,
        auth: state,
        browser: [CONFIG.botName, 'Chrome', CONFIG.version],
        getMessage: async (key) => {
            const m = await store.loadMessage(key.remoteJid, key.id);
            return m?.message || { conversation: '' };
        },
    });

    store.bind(sock.ev);

    sock.ev.on('connection.update', ({ connection, lastDisconnect, qr }) => {
        if (qr) console.log(chalk.yellowBright('\n📱 QR Code scan karo!\n'));
        if (connection === 'close') {
            const code = new Boom(lastDisconnect?.error)?.output?.statusCode;
            if (code === DisconnectReason.loggedOut) {
                console.log(chalk.red('❌ Logged out! Session delete karo.'));
                fs.rmSync(CONFIG.sessionDir, { recursive: true, force: true });
                process.exit(0);
            } else {
                console.log(chalk.yellow('🔄 Reconnecting...'));
                setTimeout(startBot, 3000);
            }
        }
        if (connection === 'open') {
            showLogo();
            console.log(chalk.greenBright('✅ Hina MD connected!\n'));
        }
    });

    sock.ev.on('creds.update', saveCreds);

    sock.ev.on('messages.upsert', async ({ messages, type }) => {
        if (type !== 'notify') return;
        for (const msg of messages) {
            try {
                if (!msg.message || msg.key.fromMe) continue;
                const from = msg.key.remoteJid;
                const body = getBody(msg);
                if (!body.startsWith(CONFIG.prefix)) continue;

                const args = body.slice(CONFIG.prefix.length).trim().split(/\s+/);
                const cmd  = args.shift().toLowerCase();

                console.log(chalk.cyan(`[CMD] .${cmd}`) + chalk.gray(` from ${from}`));

                if (CMD[cmd]) {
                    await CMD[cmd](sock, msg, from, args);
                } else {
                    await sock.sendMessage(from, {
                        text: `❓ *${CONFIG.prefix}${cmd}* nahi mila.\n\n*${CONFIG.prefix}menu* type karo!`
                    }, { quoted: msg });
                }
            } catch (e) {
                console.error(chalk.red('[ERR]'), e.message);
            }
        }
    });
}

startBot().catch(console.error);
