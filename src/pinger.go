package pinger

import (
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// Pinger - это интерфейс, который определяет метод Ping.
type Pinger interface {
	Ping(ip string) (*net.IPAddr, time.Duration, error)
}

// LocalPinger - это структура, которая реализует интерфейс Pinger.
type LocalPinger struct {
	ListenAddr string
}

// ListenAddr - это адрес по умолчанию для прослушивания ICMP-пакетов.
var ListenAddr = "0.0.0.0"

// Ping отправляет ICMP-запрос эхо-запроса на указанный IP-адрес и возвращает ответ.
func (pinger LocalPinger) Ping(ip string) (*net.IPAddr, time.Duration, error) {
	// Прослушиваем ICMP-пакеты на указанном адресе.
	con, err := icmp.ListenPacket("ip4:icmp", ListenAddr)
	if err != nil {
		log.Printf("Ошибка при прослушивании ICMP-пакетов: %v", err)
		return nil, 0, err
	}
	defer con.Close()

	// Разрешаем IP-адрес назначения.
	dst, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		log.Printf("Ошибка при разрешении IP-адреса: %v", err)
		return nil, 0, err
	}

	// Создаем пакет ICMP-запроса эхо.
	icmpPacket := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte(""),
		},
	}

	// Преобразуем пакет ICMP в байты.
	icmpPacketBytes, err := icmpPacket.Marshal(nil)
	if err != nil {
		log.Printf("Ошибка при маршалинге ICMP-пакета: %v", err)
		return dst, 0, err
	}

	// Отправляем ICMP-пакет и записываем время начала.
	start := time.Now()
	n, err := con.WriteTo(icmpPacketBytes, dst)
	if err != nil {
		log.Printf("Ошибка при записи ICMP-пакета: %v", err)
		return dst, 0, err
	} else if n != len(icmpPacketBytes) {
		log.Printf("Ошибка при записи ICMP-пакета: %v", err)
		return dst, 0, err
	}

	// Создаем буфер для хранения ответа.
	response := make([]byte, 1500)

	// Устанавливаем крайний срок чтения для ответа.
	err = con.SetReadDeadline(time.Now().Add(5 * time.Second))
	if err != nil {
		log.Printf("Ошибка при установке крайнего срока чтения: %v", err)
		return dst, 0, err
	}

	// Читаем ответ.
	n, _, err = con.ReadFrom(response)
	if err != nil {
		log.Printf("Ошибка при чтении ICMP-ответа: %v", err)
		return dst, 0, nil
	}

	// Вычисляем продолжительность пинга.
	duration := time.Since(start)

	// Разбираем ответ.
	rm, err := icmp.ParseMessage(1, response[:n])
	if err != nil {
		log.Printf("Ошибка при разборе ICMP-ответа: %v", err)
		return dst, 0, err
	}

	// Проверяем тип ответа.
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		return dst, duration, err
	default:
		log.Printf("Неожиданный тип ICMP-ответа: %v", rm.Type)
		return dst, 0, err
	}
}
