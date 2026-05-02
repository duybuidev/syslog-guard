package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/moby/moby/client"
	"github.com/moby/moby/api/types/events"
)

const discordWebhookURL = "https://discord.com/api/webhooks/1499298659228192798/OhRYIYPLvgRlr9Bt86r7nez__mOFOsG7VCUYQztmzzVlFQ7eYuHH20c-_jXXedabMbJh"

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
	fmt.Println("🔍 Điểm danh hệ thống: Đang kiểm tra các container có sẵn...")

	// API mới: ContainerList trả về ContainerListResult thay vì []container.Summary
	result, err := cli.ContainerList(context.Background(), client.ContainerListOptions{})
	if err != nil {
		log.Printf("Lỗi khi quét container: %v", err)
		return
	}

	containers := result.Items
	activeCount := len(containers)
	fmt.Printf("✅ Tìm thấy %d containers đang hoạt động.\n", activeCount)

	discordMsg := fmt.Sprintf("🛡️ **SysWatch Agent Bootup!**\n✅ Hệ thống hiện đang có **%d** containers hoạt động bình thường.", activeCount)

	for _, c := range containers {
		fmt.Printf("   - 🟢 Đang chạy: %s\n", c.Names[0][1:])
	}

	sendDiscordAlert(discordMsg)
	fmt.Println("---------------------------------------------------")
}

func main() {
	fmt.Println("🛡️ SysWatch Agent starting...")

	// API mới: dùng client.New() thay vì client.NewClientWithOpts()
	cli, err := client.New(client.FromEnv)
	if err != nil {
		log.Fatalf("Fatal: Cannot connect to Docker daemon: %v", err)
	}
	defer cli.Close()

	checkCurrentState(cli)

	fmt.Println("📡 Listening for NEW container events in real-time...\n")

	ctx := context.Background()
	// API mới: dùng client.EventsOptions thay vì types.EventsOptions
	result, errs := cli.Events(ctx, client.EventsOptions{})

	for {
		select {
		case err := <-errs:
			log.Printf("Error from Docker events: %v\n", err)
		case msg := <-result.Items:
			if msg.Type == events.ContainerEventType && msg.Action == "start" {
				name := msg.Actor.Attributes["name"]
				fmt.Printf("🟢 [%s] Container %s STARTED.\n", time.Now().Format("15:04:05"), name)
				sendDiscordAlert(fmt.Sprintf("🟢 **[STARTED]** Container `%s` is up!", name))
			}
			if msg.Type == events.ContainerEventType && msg.Action == "die" {
				name := msg.Actor.Attributes["name"]
				exitCode := msg.Actor.Attributes["exitCode"]
				fmt.Printf("🚨 [%s] Container %s CRASHED! (Code: %s)\n", time.Now().Format("15:04:05"), name, exitCode)
				sendDiscordAlert(fmt.Sprintf("🚨 **[CRASHED]** Container `%s` has DIED! (Exit: %s)", name, exitCode))
			}
		}
	}
}
