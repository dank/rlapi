# RL Version

Extracts the game version (build ID) and feature set string from a Rocket League PE binary. Use this to update the `gameVersion` and `featureSet` constants in `psynet.go` after a game patch.

## How it works

### 1. Scan `GPsyonixBuildID` export
Locates the named export in the PE export directory. The export contains a pointer to a UTF-16LE string; the tool dereferences it using the section table to convert the virtual address to a file offset and reads the null-terminated string (e.g. `260506.26700.517210`).

### 2. Scan `.rdata` for `PrimeUpdate`
Searches the `.rdata` section for all UTF-16LE occurrences of `PrimeUpdate`, reads the trailing version suffix after each match, deduplicates by `(major, suffix)`, and returns the entry with the highest major version number (e.g. `PrimeUpdate58_1`).

## Usage

```bash
go run . <path/to/RocketLeague.exe>
```