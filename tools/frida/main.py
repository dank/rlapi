import frida
import sys

def main():
    session = frida.attach("RocketLeague.exe", persist_timeout=None)

    with open("dist/index.js", "r", encoding="utf-8") as f:
        script_code = f.read()

    script = session.create_script(script_code)
    script.on("message", print)

    script.load()

    sys.stdin.read()

if __name__ == "__main__":
    main()