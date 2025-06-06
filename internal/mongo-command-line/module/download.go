package module

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(ctx context.Context, url string, descFilepath string) error {
	// 发送 HTTP GET 请求
	//resp, err := http.Get(url)
	//if err != nil {
	//	return fmt.Errorf("URL %s request fiald: %v", url, err)
	//}
	//defer resp.Body.Close()

	// 创建带有 context 的请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		//log.Errorf("could not create request: %w", err)
		return fmt.Errorf("could not create request: %w", err)
	}

	// 执行请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		//log.Errorf("error during http request: %v", err)
		return fmt.Errorf("error during http request: %v", err)
	}
	defer resp.Body.Close()

	// 检查请求是否成功
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed, status code: %d", resp.StatusCode)
	}

	// 创建文件
	outFile, err := os.Create(descFilepath)
	if err != nil {
		//log.Errorf("file create failed: %v", err)
		return fmt.Errorf("file create failed: %v", err)
	}
	defer outFile.Close()

	// 将响应内容复制到文件
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		//log.Errorf("file copy fialed: %v", err)
		return fmt.Errorf("file copy fialed: %v", err)
	}

	return nil
}
