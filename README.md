# Server Performance Stats

A lightweight, cross-platform CLI tool written in Go to monitor server performance.

## 📦 Features

- CPU, memory, disk usage, load average
- JSON or human-readable output
- Log to file with timestamp
- Export to CSV
- Threshold-based warnings
- Interval-based monitoring (`--interval`)
- Graceful shutdown on Ctrl+C (SIGINT/SIGTERM)

## 🚀 Installation

```bash
git clone https://github.com/rizkilabs/server-performance-stats.git
cd server-performance-stats
go build -o server-performance-stats
````

## 🛠️ Usage

```bash
./server-performance-stats [flags]
```

### Example:

```bash
./server-performance-stats --interval=5 --cpu-threshold=80 --json --log --export=stats.csv
```

## 📋 Flags

| Flag               | Description                        | Default |
| ------------------ | ---------------------------------- | ------- |
| `--interval`       | Run every N seconds (0 = run once) | 0       |
| `--json`           | Output stats in JSON format        | false   |
| `--log`            | Log stats to `monitor.log`         | false   |
| `--export`         | Export stats to CSV file           | ""      |
| `--cpu-threshold`  | CPU warning threshold (%)          | 80      |
| `--mem-threshold`  | Memory warning threshold (%)       | 90      |
| `--disk-threshold` | Disk usage warning threshold (%)   | 90      |

## 🧪 Testing

```bash
go test -v
```

## 🛑 Graceful Exit

Press `Ctrl+C` to exit the monitor safely. It will flush logs and close the CSV file cleanly.

## 📝 License

MIT License © 2025 Mochamad Rizki
```
