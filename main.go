package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"google.golang.org/genai"
)

var result []string

// 配置 upgrader（允许跨域示例，生产请根据需求限制）
var upgrader = websocket.Upgrader{
	// 允许所有来源（开发方便），生产请做校验
	CheckOrigin: func(r *http.Request) bool { return true },
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer conn.Close()
	log.Printf("客户端已连接: %s\n", r.RemoteAddr)

	// 设置 pong handler 与读超时（用于心跳）
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	conn.SetPongHandler(func(appData string) error {
		conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// 循环读取消息
	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			// 区分常见断开/超时错误
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读消息错误: %v\n", err)
			} else {
				log.Printf("连接已关闭: %v\n", err)
			}
			break
		}

		// 根据消息类型处理（文本/二进制）
		if len(msg) > 0 {
			switch msgType {

			case websocket.TextMessage:
				//var result map[string]json.RawMessage
				//if err := json.Unmarshal([]byte(msg), &result); err != nil {
				//	panic(err)
				//}
				smsg := string(msg)
				message := gjson.Get(smsg, "message").String()
				if len(message) > 0 {
					group_name := gjson.Get(smsg, "group_name").String()
					nickname := gjson.Get(smsg, "sender.nickname").String()
					userId := gjson.Get(smsg, "sender.user_id").String()

					fmt.Println(group_name + ":" + nickname + "(" + userId + "):" + message)
					result = append(result, group_name+":"+nickname+"("+userId+"):"+message)
				}
				//if raw, ok := result["sender"]; ok {
				//	fmt.Println(string(raw))
				//	usr := string(raw)
				//	if err := json.Unmarshal([]byte(usr), &result); err != nil {
				//		panic(err)
				//	}
				//	fields := []string{"user_id", "nickname"}
				//	for _, f := range fields {
				//		if _, ok := result[f]; ok {
				//			fmt.Printf(f)
				//		}
				//	}
				//}

				//fields := []string{"nickname", "user_id", "msg"}
				//for _, f := range fields {
				//	if _, ok := result[f]; ok {
				//		fmt.Printf("包含字段: %s\n", f)
				//	}
				//}
				//// 发送简单 ACK 回客户端
				//ack := fmt.Sprintf("ACK: received %d bytes", len(msg))
				//if err := conn.WriteMessage(websocket.TextMessage, []byte(ack)); err != nil {
				//	log.Println("写 ACK 错误:", err)
				//	break
				//}
			case websocket.BinaryMessage:
				log.Printf("收到二进制消息 (%d bytes)\n", len(msg))
				// 可根据需要处理二进制
				if err := conn.WriteMessage(websocket.TextMessage, []byte("ACK: binary received")); err != nil {
					log.Println("写 ACK 错误:", err)
					break
				}
			case websocket.CloseMessage:
				log.Println("收到关闭消息")
				return
			default:
				log.Println("收到未知消息类型:", msgType)
			}
		}
	}
}

func main() {

	http.HandleFunc("/ws", wsHandler)

	addr := ":2280"
	log.Printf("WebSocket 服务启动，监听 %s，ws 地址 ws://localhost%s/ws\n", addr, addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	ctx := context.Background()

	// 從環境變數讀 API Key（比較安全）
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("請設定環境變數 GEMINI_API_KEY")
	}

	// 建立 client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
		// Backend 默認是 Gemini API；如果要用 Vertex AI 或特別地區可以設定 Backend 或其他參數
	})
	if err != nil {
		log.Fatalf("NewClient 錯誤: %v", err)
	}

	// 發送 request，這裡用文字輸入
	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash", // model name
		genai.Text("你好"),   // prompt
		nil,                // 可選配置，例如溫度、最大 token 數等
	)
	if err != nil {
		log.Fatalf("GenerateContent 錯誤: %v", err)
	}

	// 印出結果
	fmt.Println("回應內容:", resp.Text())
}
