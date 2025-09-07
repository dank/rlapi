const SteamUser = require("steam-user");

// Steam auth session ticket generator for Rocket League
// Usage: USERNAME=username PASSWORD=password node steam-ticket.js
// Generates a hex auth session ticket that can be used with ExchangeEOSTokenFromSteam

const user = new SteamUser();

user.logOn({
    "accountName": process.env.USERNAME ?? (() => { throw new Error("env USERNAME required") })(),
    "password": process.env.PASSWORD ?? (() => { throw new Error("env PASSWORD required") })(),
});

user.on("loggedOn", async () => {
    console.log("Logged into Steam successfully");
    
    try {
        const session = await user.createAuthSessionTicket(252950); // Rocket League app ID
        const ticket = Buffer.from(session.sessionTicket).toString("hex").toUpperCase();
        
        console.log(ticket);
        console.log("Copy the ticket above and paste it into the Go Steam example");
        
        process.exit(0);
    } catch (error) {
        console.error("Failed to create auth session ticket:", error);
        process.exit(1);
    }
});

user.on("error", (err) => {
    console.error("Steam login error:", err);
    process.exit(1);
});

user.on("steamGuard", (domain, callback) => {
    console.log("Steam Guard code needed");
    const readline = require("readline");
    const rl = readline.createInterface({
        input: process.stdin,
        output: process.stdout
    });
    
    const prompt = domain ? `Enter Steam Guard code from email (${domain}): ` : "Enter Steam Guard code from mobile app: ";
    rl.question(prompt, (code) => {
        rl.close();
        callback(code);
    });
});