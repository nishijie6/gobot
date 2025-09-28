package main

import (
<<<<<<< HEAD
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
=======
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
>>>>>>> 399e239 (修改配置文件)
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"google.golang.org/genai"
)

<<<<<<< HEAD
=======
type Config struct {
	ApiKey string `yaml:"API_KEY"`
	Prompt string `yaml:"prompt"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

>>>>>>> 399e239 (修改配置文件)
var result []string

// 配置 upgrader（允许跨域示例，生产请根据需求限制）
var upgrader = websocket.Upgrader{
	// 允许所有来源（开发方便），生产请做校验
	CheckOrigin: func(r *http.Request) bool { return true },
}

<<<<<<< HEAD
=======
func callMCP(msg string) string {
	ctx := context.Background()

	config, err := LoadConfig("./config.yaml")
	if err != nil {
		log.Fatal("读取配置失败:", err)
	}

	// 配置文件读取apikey
	apiKey := config.ApiKey
	if apiKey == "" {
		log.Fatal("请在配置文件中配置 API_KEY")
	}

	// 建立 client
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
		// Backend 默認是 Gemini API；如果要用 Vertex AI 或特別地區可以設定 Backend 或其他參數
	})
	if err != nil {
		log.Fatalf("NewClient 錯誤: %v", err)
	}
	thinkingBudgetVal := int32(-1)
	// 發送 request，這裡用文字輸入
	resp, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash", // model name
		genai.Text(msg),    // prompt
		&genai.GenerateContentConfig{
			SystemInstruction: genai.NewContentFromText(config.Prompt, genai.RoleUser),
			ThinkingConfig: &genai.ThinkingConfig{
				ThinkingBudget:  &thinkingBudgetVal,
				IncludeThoughts: true,
				// Turn off thinking:
				// ThinkingBudget: int32(0),
				// Turn on dynamic thinking:
				// ThinkingBudget: int32(-1),
			}, // 可選配置，例如溫度、最大 token 數等
		})
	for _, part := range resp.Candidates[0].Content.Parts {
		if part.Text != "" {
			if part.Thought {
				fmt.Println("Thoughts Summary:")
				fmt.Println(part.Text)
			} else {
				fmt.Println("Answer:")
				//fmt.Println(part.Text)
			}
		}
	}
	if err != nil {
		log.Fatalf("GenerateContent 錯誤: %v", err)
	}

	// 输出結果
	return resp.Text()
}

func containsAny(text string, fields []string) bool {
	for _, f := range fields {
		if strings.Contains(text, f) {
			return true
		}
	}
	return false
}

func IsResponse(msg string) bool {
	text := []string{"{\"type\":\"at\",\"data\":{\"qq\":\"1911407507\"}}", "holly", "Holly"}
	if containsAny(msg, text) {
		return true
	} else {
		rand.Seed(time.Now().UnixNano())
		x := rand.Float64()
		if x < 1 {
			fmt.Println(x)
			return true
		}
	}
	return false
}

func wsReMsg(msg string) {
	url := "http://localhost:16000/send_group_msg"

	payload := map[string]interface{}{
		"group_id": "253631878",
		"message": []interface{}{
			map[string]interface{}{
				"type": "text",
				"data": map[string]interface{}{
					"text": msg,
				},
			},
		},
	}
	jsonData, _ := json.Marshal(payload)

	resp, _ := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}

// 连接napcatqq获取消息
>>>>>>> 399e239 (修改配置文件)
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
<<<<<<< HEAD
	conn.SetReadDeadline(time.Now().Add(60 * time.Second))
=======
	conn.SetReadDeadline(time.Now().Add(600 * time.Second))
>>>>>>> 399e239 (修改配置文件)
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
<<<<<<< HEAD
				message := gjson.Get(smsg, "message").String()
				if len(message) > 0 {
					group_name := gjson.Get(smsg, "group_name").String()
=======
				//fmt.Println(smsg)
				message := gjson.Get(smsg, "raw_message").String()
				group_name := gjson.Get(smsg, "group_name").String()
				if len(message) > 0 && group_name == "" {
>>>>>>> 399e239 (修改配置文件)
					nickname := gjson.Get(smsg, "sender.nickname").String()
					userId := gjson.Get(smsg, "sender.user_id").String()

					fmt.Println(group_name + ":" + nickname + "(" + userId + "):" + message)
<<<<<<< HEAD
					result = append(result, group_name+":"+nickname+"("+userId+"):"+message)
=======
					//result = append(result, group_name+":"+nickname+"("+userId+"):"+message)
					if IsResponse(message) {
						// 调用大模型
						res := callMCP(nickname + "(" + userId + "):" + message)
						// 返回结果
						//wsReMsg(res)
						fmt.Println(res)
					}
>>>>>>> 399e239 (修改配置文件)
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
<<<<<<< HEAD

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
=======
	//
	//url := "http://localhost:16000/send_private_msg"
	//
	//payload := map[string]interface{}{
	//	"user_id": "357020327",
	//	"message": []interface{}{
	//		map[string]interface{}{
	//			"type": "text",
	//			"data": map[string]interface{}{
	//				"text": "napcat",
	//			},
	//		},
	//	},
	//}
	//jsonData, _ := json.Marshal(payload)
	//
	//// 创建 POST 请求
	//req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	//req.Header.Set("Content-Type", "application/json")
	//
	//client := &http.Client{}
	//resp, err := client.Do(req)
	//if err != nil {
	//	panic(err)
	//}
	//defer resp.Body.Close()
	//
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
>>>>>>> 399e239 (修改配置文件)
}
