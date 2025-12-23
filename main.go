package main

import "net"
import "log"
import "time"
import "io"
import "dosgo/cloudflareSSH/comm"

var stopChan chan struct{}



func main() {
	startProxy(":8024" ,"armbian.16v16.com")
}

func startProxy(tcpPort string, hostName string) {
	// 启动 TCP 服务器
	listener, err := net.Listen("tcp", tcpPort)
	if err != nil {
		log.Fatalf("TCP监听失败: %v", err)
	}
	defer listener.Close()

	log.Printf("TCP服务器启动在 %s，等待连接...", tcpPort)


	stopChan = make(chan struct{})


	for {
		select {
		case <-stopChan:
			log.Println("收到关闭信号，退出")
			return
		default:
			// 设置超时以便能检查信号
			listener.(*net.TCPListener).SetDeadline(time.Now().Add(1 * time.Second))
			tcpConn, err := listener.Accept()
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue
				}
				log.Printf("接受连接失败: %v", err)
				continue
			}
			log.Printf("客户端连接: %s", tcpConn.RemoteAddr())
			// 处理连接
			go handleConnection(tcpConn,hostName)
		}
	}
}
func stopProxy() {
	if stopChan != nil {
		close(stopChan)
	}
}

func handleConnection(tcpConn net.Conn,hostName string) {
	defer tcpConn.Close()
	sshPorxy,err:=comm.NewCloudflaredSSH(hostName)
	if err != nil {
		log.Printf("创建Cloudflared SSH连接失败: %v", err)
		return
	}
	defer sshPorxy.Close()
	// TCP → 串口
	go func() {
		_, err = io.Copy(sshPorxy, tcpConn)
		if err != nil {
			log.Printf("TCP→串口转发错误: %v", err)
		}

	}()
	_, err = io.Copy(tcpConn, sshPorxy)
	if err != nil {
		log.Printf("串口→TCP转发错误: %v", err)
	}
	log.Printf("连接断开: %s", tcpConn.RemoteAddr())
}
