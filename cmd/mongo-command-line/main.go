package main

import mongocommandline "DatabaseManage/internal/mongo-command-line"

func main() {
	mongocommandline.New("mongo-command-line").Run()
}
