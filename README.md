# Pinger

**Pinger** is a lightweight, high-performance network utility written in Go. 
It enables users to send ICMP echo requests (pings) to specified hosts, facilitating network diagnostics and monitoring.

## Features

* Send ICMP echo requests to specified hosts.
* Measure round-trip time (RTT) for each ping.
* Support for both IPv4 and IPv6 addresses.
* Configurable number of ping attempts and intervals.
* Lightweight and efficient, suitable for scripting and automation.

## Installation

### Prerequisites

* Go (version 1.20 or later) installed on your system.
* make

### Steps

1. Clone the repository:

```bash
git clone https://github.com/pnx/pinger.git
cd pinger
```

2. Build the application:

```bash
make
```

## Usage

See `./pinger -h`

### Example

```bash
./pinger --udp -t 30s example.com
```

This command sends UDP echo requests to example.com, stopping after 30 seconds

```bash
./pinger -c 8 -i 1s example.com
```

This command sends 8 ICMP echo requests to example.com, with a 1-second interval between pings.
NOTE: this requires root privileges

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
