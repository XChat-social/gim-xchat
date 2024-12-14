package main

import (
	"context"
	"fmt"
	pb "gim/pkg/protocol/pb" // 引入生成的 gRPC 包
	"log"
	"net/http"
	"net/url"
	"strconv"

	"google.golang.org/grpc"
)

const (
	grpcAddress = "localhost:8020" // gRPC 服务地址
	httpAddress = ":8080"          // HTTP 服务地址
)

// TwitterSignInHandler 处理 Twitter 回调请求
func TwitterSignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// 提取回调参数
	authorizationCode := r.URL.Query().Get("code")
	if authorizationCode == "" {
		http.Error(w, "Authorization code is required", http.StatusBadRequest)
		return
	}

	state := r.URL.Query().Get("state")
	if state == "" {
		http.Error(w, "State is required", http.StatusBadRequest)
		return
	}

	// 调用 gRPC 服务
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
		log.Printf("gRPC connection error: %v\n", err)
		return
	}
	defer conn.Close()

	client := pb.NewBusinessExtClient(conn)
	grpcResp, err := client.TwitterSignIn(context.Background(), &pb.TwitterSignInReq{
		AuthorizationCode: authorizationCode,
		State:             state,
	})
	if err != nil {
		http.Error(w, "Failed to call gRPC service: "+err.Error(), http.StatusInternalServerError)
		log.Printf("gRPC call error: %v\n", err)
		return
	}

	//// 提取 gRPC 响应中的用户信息
	//response := map[string]interface{}{
	//	"token":           grpcResp.Token,
	//	"userId":          grpcResp.UserId,
	//	"nickname":        grpcResp.UserInfo.Nickname,
	//	"avatarUrl":       grpcResp.UserInfo.AvatarUrl,
	//	"twitterUsername": grpcResp.UserInfo.TwitterUsername,
	//}
	//
	//// 设置响应头为 JSON
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//
	//// 返回 JSON 响应
	//if err := json.NewEncoder(w).Encode(response); err != nil {
	//	http.Error(w, "Failed to encode response: "+err.Error(), http.StatusInternalServerError)
	//	log.Printf("JSON encoding error: %v\n", err)
	//	return
	//}

	// 提取 gRPC 响应中的用户信息
	token := grpcResp.Token
	userId := grpcResp.UserId
	nickname := grpcResp.UserInfo.Nickname
	avatarUrl := grpcResp.UserInfo.AvatarUrl
	twitterUsername := grpcResp.UserInfo.TwitterUsername

	// 构造插件页面的跳转 URL，附加用户信息
	redirectURL := fmt.Sprintf(
		"https://x.com?redirect=redirectx&token=%s&userId=%s&nickname=%s&avatarUrl=%s&twitterUsername=%s",
		url.QueryEscape(token),
		url.QueryEscape(strconv.FormatInt(userId, 10)),
		url.QueryEscape(nickname),
		url.QueryEscape(avatarUrl),
		url.QueryEscape(twitterUsername),
	)

	// 重定向到插件页面
	http.Redirect(w, r, redirectURL, http.StatusFound)
	log.Printf("Redirecting to: %s\n", redirectURL)
}

func main() {
	http.HandleFunc("/twitter/signin", TwitterSignInHandler)

	log.Printf("HTTP server is running on %s\n", httpAddress)
	if err := http.ListenAndServe(httpAddress, nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
