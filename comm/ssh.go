// cloudflared_rwc.go
package comm

import (
    "io"
    "os/exec"

)

// CloudflaredSSH 包装 cloudflared SSH 为 io.ReadWriteCloser
type CloudflaredSSH struct {
    cmd    *exec.Cmd
    stdin  io.WriteCloser
    stdout io.ReadCloser
}

// 创建新的 Cloudflared SSH 连接
func NewCloudflaredSSH(hostname string) (*CloudflaredSSH, error) {
    cmd := exec.Command("cloudflared", "access", "ssh", "--hostname", hostname)
    
    stdin, err := cmd.StdinPipe()
    if err != nil {
        return nil, err
    }
    
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return nil, err
    }

    
    if err := cmd.Start(); err != nil {
        return nil, err
    }
    
    return &CloudflaredSSH{
        cmd:    cmd,
        stdin:  stdin,
        stdout: stdout,
    }, nil
}

// 实现 io.Reader 接口（读取 stdout）
func (c *CloudflaredSSH) Read(p []byte) (n int, err error) {

    return c.stdout.Read(p)
}

// 实现 io.Writer 接口（写入 stdin）
func (c *CloudflaredSSH) Write(p []byte) (n int, err error) {
    return c.stdin.Write(p)
}

// 实现 io.Closer 接口
func (c *CloudflaredSSH) Close() error {
   
    
    // 关闭 stdin
    c.stdin.Close()
    
    // 等待进程结束
    return c.cmd.Wait()
}

