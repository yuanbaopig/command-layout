package module

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// AddEnvToFile 检查并添加环境变量到 .bashrc 文件
func AddEnvToFile(profileFile, variable, value string) error {

	// 检查 .bashrc 文件是否已包含该变量
	file, err := os.Open(profileFile)
	if err != nil {
		return fmt.Errorf("failed to open .bashrc: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	variableDeclaration := fmt.Sprintf("export %s=%s", variable, value)

	for scanner.Scan() {
		// 如果找到该变量的定义，且值相同，说明不需要再添加
		//if strings.Contains(scanner.Text(), fmt.Sprintf("export %s=", variable)) {
		// 检查是否已定义相同的值
		if scanner.Text() == variableDeclaration {
			//fmt.Println("Environment variable already set with the same value.")
			return nil
		}

		// 已有定义但值不同，更新为新值
		//fmt.Println("Updating existing environment variable to new value.")
		//return updateFileRow(profileFile, variableDeclaration)	// PATH环境变量 当前需求，没有定义相同，但是值更新的逻辑
		//}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scanner error: %v", err)
	}

	// 如果文件中没有找到该变量定义，则追加新定义
	file, err = os.OpenFile(profileFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file for appending: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(variableDeclaration + "\n"); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	//fmt.Println("Environment variable added successfully.")
	return nil
}

// 按行读取和写入数据，避免占用大量内存，按行包含匹配，不是等值匹配
func updateFileRow(path, oldRow, newRow string) error {
	fi, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer fi.Close()

	reader := bufio.NewReader(fi)
	var currentOffset int64

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read line: %v", err)
		}

		// 当前行的开始位置为已累加的偏移量
		//fmt.Println("Current line start offset:", currentOffset)

		//fmt.Printf(line)
		// 移除行尾换行符，保证字符串处理的一致性
		line = strings.TrimSuffix(line, "\n")
		fmt.Println(line)
		// 检查是否包含 oldRow
		// 需要注意，这里是包含，不是等值匹配
		if strings.Contains(line, oldRow) {
			// 计算当前行的长度
			lineLen := len(line)

			// 判断新的行是否比源行要短
			if len(newRow) <= lineLen {
				// 原地更新

				// 游标到行开始位置
				if _, err := fi.Seek(currentOffset, io.SeekStart); err != nil {
					return fmt.Errorf("failed to seek to line start: %v", err)
				}
				// 如果 newRow 比原行短，用空格补全覆盖
				if _, err = fi.WriteString(newRow + strings.Repeat(" ", lineLen-len(newRow)) + "\n"); err != nil {
					return fmt.Errorf("failed to write new line with padding: %v", err)
				}
			} else {
				// 末尾追加

				// 如果 newRow 长度大于原行长度，则清空当前行

				// 游标到行开始位置
				if _, err := fi.Seek(currentOffset, io.SeekStart); err != nil {
					return fmt.Errorf("failed to seek to line start for clearing: %v", err)
				}

				// 整行使用空格覆盖
				if _, err = fi.WriteString(strings.Repeat(" ", lineLen) + "\n"); err != nil {
					return fmt.Errorf("failed to clear current line: %v", err)
				}

				// 游标到末尾
				if _, err = fi.Seek(0, io.SeekEnd); err != nil {
					return fmt.Errorf("failed to seek to end of file: %v", err)
				}

				if _, err = fi.WriteString(newRow + "\n"); err != nil {
					return fmt.Errorf("failed to append new row to end of file: %v", err)
				}
			}
			break // 假设只需要更新第一个匹配项，若需要更新所有匹配项，移除此行
		}

		// 更新偏移量，包含当前行的长度
		currentOffset += int64(len(line) + 1) // 假设换行符为单个字节
	}

	return nil

}
