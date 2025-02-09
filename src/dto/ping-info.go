package dto

import "time"

// PingInfo представляет информацию о событии пинга.
type PingInfo struct {
	// Ip - это IP-адрес устройства, которое было пропинговано.
	// PingTime - это время, которое потребовалось для пинга.
	// LastSeen - это время, когда произошел пинг.
	Ip       string    `json:"ip_address"`
	PingTime float64   `json:"ping_time"`
	LastSeen time.Time `json:"last_seen"`
}
