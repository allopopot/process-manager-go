package views

import (
	"fmt"
	"strings"
	"time"

	"allopopot.com/proc-man/config"
	"charm.land/lipgloss/v2"
	psnet "github.com/shirou/gopsutil/v4/net"
)

// prevCounters and prevTime hold the last snapshot for rate calculation.
var (
	prevCounters map[string]psnet.IOCountersStat
	prevTime     time.Time
)

func RenderNetworkInfo() string {
	const ifaceW, addrW, rateW, totalW = 14, 22, 14, 14

	now := time.Now()
	counters, _ := psnet.IOCounters(true)

	// Build current counter map.
	cc := make(map[string]psnet.IOCountersStat, len(counters))
	for _, s := range counters {
		cc[s.Name] = s
	}

	// Compute per-interface rates if we have a previous snapshot.
	type rates struct{ txBps, rxBps float64 }
	rateMap := make(map[string]rates)
	if prevCounters != nil {
		elapsed := now.Sub(prevTime).Seconds()
		if elapsed > 0 {
			for name, cur := range cc {
				if prev, ok := prevCounters[name]; ok {
					rateMap[name] = rates{
						txBps: float64(cur.BytesSent-prev.BytesSent) / elapsed,
						rxBps: float64(cur.BytesRecv-prev.BytesRecv) / elapsed,
					}
				}
			}
		}
	}

	// Save snapshot for next call.
	prevCounters = cc
	prevTime = now

	ifaces, _ := psnet.Interfaces()

	s := ""

	s += lipgloss.NewStyle().Foreground(config.ColorAccent).Render("Network")
	s += "\n"
	s += mutedStyle.Render(strings.Repeat("─", ifaceW+addrW+rateW+totalW+rateW+totalW+6))
	s += "\n"
	s += headerStyle.Render(
		fmt.Sprintf("%-*s %-*s %*s %*s %*s %*s",
			ifaceW, "Interface",
			addrW, "Address",
			rateW, "↑ Upload/s",
			rateW, "↓ Download/s",
			totalW, "Total Sent",
			totalW, "Total Recv",
		),
	)
	s += "\n"
	s += mutedStyle.Render(strings.Repeat("─", ifaceW+addrW+rateW+totalW+rateW+totalW+6))
	s += "\n"

	for _, iface := range ifaces {
		c := cc[iface.Name]
		r := rateMap[iface.Name]

		addrs := iface.Addrs
		if len(addrs) == 0 {
			// Show interface with no addresses.
			addrs = []psnet.InterfaceAddr{{Addr: "—"}}
		}

		for i, addr := range addrs {
			var line string
			if i == 0 {
				// First row: show interface name, rates, and totals.
				txStr := formatBytesRate(r.txBps)
				rxStr := formatBytesRate(r.rxBps)
				if prevCounters == nil || r.txBps == 0 && r.rxBps == 0 {
					txStr = "—"
					rxStr = "—"
				}
				line = fmt.Sprintf("%-*s %-*s %*s %*s %*s %*s",
					ifaceW, truncate(iface.Name, ifaceW),
					addrW, truncate(addr.Addr, addrW),
					rateW, txStr,
					rateW, rxStr,
					totalW, formatBytes(c.BytesSent),
					totalW, formatBytes(c.BytesRecv),
				)
			} else {
				// Subsequent rows: only show the extra address.
				line = fmt.Sprintf("%-*s %-*s %*s %*s %*s %*s",
					ifaceW, "",
					addrW, truncate(addr.Addr, addrW),
					rateW, "",
					rateW, "",
					totalW, "",
					totalW, "",
				)
			}
			s += textStyle.Render(line) + "\n"
		}
	}

	return s
}

func formatBytesRate(bps float64) string {
	const (
		KB = 1024.0
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bps >= GB:
		return fmt.Sprintf("%.2f GB/s", bps/GB)
	case bps >= MB:
		return fmt.Sprintf("%.2f MB/s", bps/MB)
	case bps >= KB:
		return fmt.Sprintf("%.2f KB/s", bps/KB)
	default:
		return fmt.Sprintf("%.0f B/s", bps)
	}
}
