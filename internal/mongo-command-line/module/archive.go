package module

import (
	"context"
	"fmt"
	"github.com/mholt/archiver/v4"
	"io"
	"os"
	"path/filepath"
)

// Extract archiver解压压缩包
func Extract(ctx context.Context, path, descDir string) error {
	f, err := os.Open(path)
	if err != nil {
		//log.Errorf("%s file open failed error: %v", path, err)
		return fmt.Errorf("%s file open failed error: %v", path, err)
	}

	format, stream, err := archiver.Identify(path, f)
	if err != nil {
		//log.Errorf("archiver format identify error: %v", err)
		return fmt.Errorf("archiver format identify error: %v", err)
	}

	extractor, ok := format.(archiver.Extractor)
	if !ok {
		//log.Error("archiver format error, unsupported extractor interface")
		return fmt.Errorf("archiver format error, unsupported extractor interface")
	}

	//switch extractor.(type) {
	//case archiver.Zip:
	//	//fmt.Println("archiver.Zip")
	//	log.Debugf("%s archiver type is zip", f.Name())
	//	extractor = archiver.Zip{}
	//
	//case archiver.Tar:
	//	//fmt.Println("archiver.tar")
	//	log.Debugf("%s archiver type is tar", f.Name())
	//	extractor = archiver.Tar{}
	//
	//case archiver.CompressedArchive:
	//	log.Debugf("%s archiver type is compressed archive", f.Name())
	//	extractor = archiver.CompressedArchive{
	//		Compression: archiver.Gz{},
	//		Archival:    archiver.Tar{},
	//	}
	//
	//default:
	//	log.Errorf("%s unsupported compression algorithm, archiver type is %T", f.Name(), extractor)
	//	return fmt.Errorf("unsupported compression algorithm")
	//}

	handler := func(ctx context.Context, f archiver.File) error {
		// 遇到空目录就跳过（有些打包格式会把空目录当作一个文件）
		if f.IsDir() {
			return nil
		}

		// 获取文件在压缩包中的路径
		filePath := f.NameInArchive

		// 构建解压缩后的文件路径
		// 这里使用了 filepath.Join 函数，避免路径拼接错误
		targetPath := filepath.Join(descDir, filePath)

		// 获取文件的路径
		dir := filepath.Dir(targetPath)
		// 创建文件目录
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create directory %s, failed: %v", dir, err)
		}

		// 创建目标文件，并设置权限
		//outFile, err := os.Create(targetPath)
		// 打开文件，如果文件中有内容，则truncate掉
		outFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("file create %s failed: %v", targetPath, err)
		}
		defer outFile.Close()

		body, err := f.Open()
		if err != nil {
			return fmt.Errorf("file open %s failed: %v", targetPath, err)
		}
		defer body.Close()

		// 将响应内容复制到文件
		_, err = io.Copy(outFile, body)
		if err != nil {
			return fmt.Errorf("file copy fialed: %v", err)
		}
		return nil
	}

	err = extractor.Extract(ctx, stream, nil, handler)
	if err != nil {
		//log.Errorf("archiver extract failed, error: %v", err)
		return fmt.Errorf("archiver extract failed, error: %v", err)
	}

	return nil
}
