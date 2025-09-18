const rl = Process.findModuleByName("RocketLeague.exe")!;

console.log(`RocketLeague.exe: ${rl.base}`);

function findPattern(pattern: string): NativePointer {
    const result = Memory.scanSync(rl.base, rl.size, pattern);
    if (result.length === 0) {
        throw new Error(`Pattern not found: ${pattern}`);
    }
    console.log(`[+] Found pattern ${pattern} at: ${result[0].address}`);
    return result[0].address;
}

const curlEasyInit = findPattern("48 83 EC ?? 83 3D ?? ?? ?? ?? ?? 0F 85 8D");
const curlSetOpts = findPattern("89 54 24 ?? 4C 89 44 24 ?? 4C 89 4C 24 ?? 48 83 EC ?? 48 85");
const x509VerifyCert = findPattern("40 53 57 B8 ?? ?? ?? ?? E8 ?? ?? ?? ?? 48 2B E0 48 8B D9");

Interceptor.attach(curlEasyInit, {
    onLeave(retval) {
        // console.log("[+] curl_easy_init handle:", retval);

        const curlSetOptsFn = new NativeFunction(curlSetOpts, "int", ["pointer", "int", "int"]);
        curlSetOptsFn(retval, 64, 0); // CURLOPT_SSL_VERIFYPEER 
        curlSetOptsFn(retval, 81, 0); // CURLOPT_SSL_VERIFYHOST  
    }
});

Interceptor.attach(curlSetOpts, {
    onEnter(args) {
        const option = args[1].toInt32();
        const param = args[2];

        if (option === 10002) { // CURLOPT_URL
            let url = param.readUtf8String()!;
            // console.log("[+] CURLOPT_URL:", url);

            const targetHost = "https://api.rlpp.psynet.gg";
            if (url.startsWith(targetHost)) {
                const replacement = "https://127.0.0.1";
                const path = url.slice(targetHost.length);
                const newUrl = replacement + path;

                const newUrlPtr = Memory.allocUtf8String(newUrl);
                args[2] = newUrlPtr;

                console.log(`[+] CURLOPT_URL: ${url} => ${newUrl}`);
            }
        }

        // console.log(`[+] curl_easy_setopt(${args[0]}, option=${option}, param=${param})`)
    }
});

Interceptor.replace(x509VerifyCert, new NativeCallback(function (_a1) {
    console.log('[+] Bypassing X509_verify_cert');
    return 1;
}, 'int64', ['int64']));