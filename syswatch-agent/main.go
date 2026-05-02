package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const discordWebhookURL = "https://discord.com/api/webhooks/1499298659228192798/OhRYIYPLvgRlr9Bt86r7nez__mOFOsG7VCUYQztmzzVlFQ7eYuHH20c-_jXXedabMbJh"

// Docker API structs
type Container struct {
	Names []string `json:"Names"`
}

type DockerEvent struct {
	Type   string `json:"Type"`
	Action string `json:"Action"`
	Actor  struct {
		Attributes map[string]string `json:"Attributes"`
	} `json:"Actor"`
}

type DiscordPayload struct {
	Content string `json:"content"`
}

// HTTP client giao tiếp qua Unix socket
func newDockerClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				return (&net.Dialer{}).DialContext(ctx, "unix", "/var/run/docker.sock")
			},
		},
	}
}

func sendDiscordAlert(message string) {
	payload := DiscordPayload{Content: message}
	jsonValue, _ := json.Marshal(payload)
	resp, err := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(jsonValue))
	if err == nil {
		resp.Body.Close()
	}
}

func checkCurrentState(client *http.Client, expectedCount int) {
	fmt.Println("🔍 Điểm danh hệ thống: Đang kiểm tra các container có sẵn...")

	// Retry tối đa 30 lần, mỗi lần cách nhau 2 giây
	for i := 0; i < 30; i++ {
		resp, err := client.Get("http://localhost/containers/json")
		if err != nil {
			log.Printf("Lỗi khi quét container: %v", err)
			time.Sleep(2 * time.Second)
			continue
		}

		var containers []Container
		json.NewDecoder(resp.Body).Decode(&containers)
		resp.Body.Close()

		activeCount := len(containers)
		fmt.Printf("⏳ Lần %d: Tìm thấy %d containers...\n", i+1, activeCount)

		// Chờ đủ số container mong đợi (không tính syswatch-agent)
		if activeCount >= expectedCount {
			fmt.Printf("✅ Đủ %d containers đang hoạt động!\n", activeCount)

			for _, c := range containers {
				if len(c.Names) > 0 {
					fmt.Printf("   - 🟢 Đang chạy: %s\n", c.Names[0][1:])
				}
			}

			discordMsg := fmt.Sprintf("🛡️ **SysWatch Agent Bootup!**\n✅ Hệ thống hiện đang có **%d** containers hoạt động bình thường.", activeCount)
			sendDiscordAlert(discordMsg)
			fmt.Println("---------------------------------------------------")
			return
		}

		time.Sleep(2 * time.Second)
	}

	log.Println("⚠️ Timeout: Không đủ container sau 60 giây")
}

func listenEvents(client *http.Client) {
	resp, err := client.Get("http://localhost/events")
	if err != nil {
		log.Fatalf("Không thể kết nối Docker events: %v", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		var event DockerEvent
		if err := json.Unmarshal(line, &event); err != nil {
			continue
		}

		if event.Type != "container" {
			continue
		}

		name := event.Actor.Attributes["name"]
		timestamp := time.Now().Format("15:04:05")

		switch event.Action {
		case "start":
			fmt.Printf("🟢 [%s] Container %s STARTED.\n", timestamp, name)
			sendDiscordAlert(fmt.Sprintf("🟢 **[STARTED]** Container `%s` is up!", name))
		case "die":
			exitCode := event.Actor.Attributes["exitCode"]
			fmt.Printf("🚨 [%s] Container %s CRASHED! (Code: %s)\n", timestamp, name, exitCode)
			sendDiscordAlert(fmt.Sprintf("🚨 **[CRASHED]** Container `%s` has DIED! (Exit: %s)", name, exitCode))
		}
	}
}

func main() {
	fmt.Println("🛡️ SysWatch Agent starting...")

	client := newDockerClient()

	// Truyền số container mong đợi (không tính syswatch-agent)
	// order(1) + shipping(3) + auth(1) = 5
	checkCurrentState(client, 5)

	fmt.Println("📡 Listening for NEW container events in real-time...\n")

	listenEvents(client)
}
