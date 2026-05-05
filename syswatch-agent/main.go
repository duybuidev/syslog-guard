package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

const discordWebhookURL = "https://discord.com/api/webhooks/xxxx/yyyy"

type DiscordPayload struct {
	Content string `json:"content"`
}

func sendDiscordAlert(message string) {
	if discordWebhookURL == "https://discord.com/api/webhooks/xxxx/yyyy" {
		return
	}
	payload := DiscordPayload{Content: message}
	jsonValue, _ := json.Marshal(payload)
	resp, err := http.Post(discordWebhookURL, "application/json", bytes.NewBuffer(jsonValue))
	if err == nil {
		resp.Body.Close()
	}
}

func checkCurrentState(cli *client.Client) {
	fmt.Println("🔍 Điểm danh hệ thống...")
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		log.Printf("Lỗi khi quét container: %v", err)
		return
	}
	activeCount := len(containers)
	fmt.Printf("✅ Tìm thấy %d containers đang hoạt động.\n", activeCount)
	for _, c := range containers {
		fmt.Printf("   - 🟢 %s\n", c.Names[0][1:])
	}
	sendDiscordAlert(fmt.Sprintf("🛡️ **SysWatch Bootup!**\n✅ **%d** containers đang hoạt động.", activeCount))
	fmt.Println("---------------------------------------------------")
}

func main() {
	fmt.Println("🛡️ SysWatch Agent starting...")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("Fatal: Cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	checkCurrentState(cli)

	fmt.Println("📡 Listening for container events...")
	ctx := context.Background()
	msgs, errs := cli.Events(ctx, types.EventsOptions{})

	for {
		select {
		case err := <-errs:
			log.Printf("Error: %v\n", err)
		case msg := <-msgs:
			if msg.Type == events.ContainerEventType && msg.Action == "start" {
				name := msg.Actor.Attributes["name"]
				fmt.Printf("🟢 [%s] %s STARTED\n", time.Now().Format("15:04:05"), name)
				sendDiscordAlert(fmt.Sprintf("🟢 **[STARTED]** `%s` is up!", name))
			}
			if msg.Type == events.ContainerEventType && msg.Action == "die" {
				name := msg.Actor.Attributes["name"]
				exitCode := msg.Actor.Attributes["exitCode"]
				fmt.Printf("🚨 [%s] %s CRASHED! (Code: %s)\n", time.Now().Format("15:04:05"), name, exitCode)
				sendDiscordAlert(fmt.Sprintf("🚨 **[CRASHED]** `%s` died! (Exit: %s)", name, exitCode))
			}
		}
	}
}
