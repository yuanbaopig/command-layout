package module

import (
	"context"
	"net"
)

// PortAvailable 判断指定地址和端口是否可用
func PortAvailable(ctx context.Context, address string) bool {
	//address := fmt.Sprintf("%s:%d", host, port)
	dialer := net.Dialer{}

	// 使用 DialContext 方法，支持超时和取消
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return false // 端口不可用
	}
	defer conn.Close()

	return true // 端口可用
}
