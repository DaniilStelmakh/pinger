package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	pinger "github.com/DaniilStelmakh/pinger/src"
	"github.com/DaniilStelmakh/pinger/src/dto"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	// Загружаем переменные окружения из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Устанавливаем интервал для пинга в секундах.
	intervalStr := os.Getenv("INTERVAL")
	if intervalStr == "" {
		log.Fatal("invalid is required")
	}
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		log.Fatal("invalid internal: %v", err)
	}

	ips := os.Getenv("IPS")
	if ips == "" {
		log.Fatal("ip address is required")
	}

	// Разделяем IP-адреса на срез.
	ipPasre := strings.Split(ips, ",")

	// Бесконечный цикл.
	for {
		// Цикл по каждому IP-адресу.
		for _, ip := range ipPasre {
			// Создаем новый пингер для IP-адреса.
			pinger := (&pinger.LocalPinger{ListenAddr: ip})

			// Пингуем IP-адрес.
			dst, dur, err := pinger.Ping(ip)
			if err != nil {
				log.Printf("Ошибка при пинге IP-адреса %s: %v", ip, err)
				continue
			}

			// Создаем новый PingInfo с информацией о пинге.
			ping := &dto.PingInfo{Ip: dst.String(), PingTime: dur.Seconds(), LastSeen: time.Now()}

			// Преобразуем PingInfo в JSON.
			jsn, err := json.Marshal(&ping)
			if err != nil {
				log.Printf("Ошибка при маршалинге PingInfo: %v", err)
				return
			}

			// Отправляем JSON на сервер.
			url := fmt.Sprintf("http://%s:%s/pings", os.Getenv("API_HOST"), os.Getenv("API_PORT"))
			res, err := http.Post(url, "application/json", bytes.NewBuffer(jsn))
			if err != nil {
				log.Printf("Ошибка при отправке JSON на сервер: %v", err)
				return
			}

			// Закрываем тело ответа.
			err = res.Body.Close()
			if err != nil {
				log.Printf("Ошибка при закрытии тела ответа: %v", err)
				return
			}
		}

		// Ждем указанный интервал перед следующим пингом.
		time.Sleep(time.Duration(interval) * time.Second)
	}
}
