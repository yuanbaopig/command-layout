package module

import (
	"context"
	"fmt"
	"github.com/coreos/go-systemd/v22/dbus"
)

const (
	PropertyActiveState = "ActiveState" // inactive,active
	Active              = "active"
	Inactive            = "inactive"
	PropertySubState    = "SubState"
	SubSateRunning      = "running"
	SubStateDead        = "dead"
	PropertyLoadState   = "LoadState"
)

func StartService(ctx context.Context, serviceName string) error {
	conn, err := dbus.NewSystemConnectionContext(ctx) // 建立到 systemd 的连接
	if err != nil {
		return fmt.Errorf("failed to connect to systemd: %v", err)
	}
	defer conn.Close()

	ch := make(chan string)
	// 请求启动服务
	_, err = conn.StartUnitContext(ctx, serviceName, "replace", ch)
	if err != nil {
		return fmt.Errorf("failed to start service of %s: %v", serviceName, err)
	}
	//fmt.Printf("Started job ID: %d\n", jobID)
	// 等待服务启动完成
	jobStatus := <-ch
	//<-ch
	if jobStatus != "done" {
		return fmt.Errorf("failed to start service of %s: status %s", serviceName, jobStatus)
	}
	return nil
}

func StopService(ctx context.Context, serviceName string) error {
	conn, err := dbus.NewSystemConnectionContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to systemd: %v", err)
	}
	defer conn.Close()

	ch := make(chan string)
	_, err = conn.StopUnitContext(ctx, serviceName, "replace", ch)
	if err != nil {
		return fmt.Errorf("failed to stop service of %s: %v", serviceName, err)
	}
	//fmt.Printf("Stoped jobStatus ID: %d\n", jobID)
	jobStatus := <-ch

	if jobStatus != "done" {
		return fmt.Errorf("failed to stop service of %s: status %s", serviceName, jobStatus)
	}
	//fmt.Println("Service stop status:", jobStatus)
	return nil
}

func EnableService(ctx context.Context, unitFiles []string) error {
	conn, err := dbus.NewSystemConnectionContext(ctx) // 建立到 systemd 的连接
	if err != nil {
		return fmt.Errorf("failed to connect to systemd: %v", err)
	}
	defer conn.Close()

	_, _, err = conn.EnableUnitFilesContext(ctx, unitFiles, false, true)
	if err != nil {
		return fmt.Errorf("failed to enable unit files: %v", err)
	}

	// 打印结果
	//fmt.Printf("Runtime Flag: %v\n", runtimeFlag)
	//for _, change := range changes {
	//	fmt.Printf("Type: %s, Link: %s, Destination: %s\n", change.Type, change.Filename, change.Destination)
	//}
	return nil
}

func DisableService(ctx context.Context, unitFiles []string) error {
	conn, err := dbus.NewSystemConnectionContext(ctx) // 建立到 systemd 的连接
	if err != nil {
		return fmt.Errorf("failed to connect to systemd: %v", err)
	}
	defer conn.Close()

	_, err = conn.DisableUnitFilesContext(ctx, unitFiles, false)
	if err != nil {
		return fmt.Errorf("failed to disable units: %v", err)
	}

	//fmt.Println("Changes made during disabling:")
	//for _, change := range changes {
	//	fmt.Printf("Type: %s, Link: %s, Destination: %s\n", change.Type, change.Filename, change.Destination)
	//}
	return nil
}

func StatusService(ctx context.Context, unit, property string) (*dbus.Property, error) {
	conn, err := dbus.NewSystemConnectionContext(ctx) // 建立到 systemd 的连接
	if err != nil {
		return nil, fmt.Errorf("failed to connect to systemd: %v", err)
	}
	defer conn.Close()

	//property := "ActiveState"

	return conn.GetUnitPropertyContext(ctx, unit, property)
	//prop, err := conn.GetUnitPropertyContext(ctx, unit, property)
	//if err != nil {
	//	return nil, fmt.Errorf("failed to get property %s of %s: %v", property, unit, err)
	//}

	// 打印属性值
	//fmt.Printf("Property %s of unit %s: %v\n", property, unit, prop.Value.Value())

	//return nil
}

func SystemdReload(ctx context.Context) error {
	conn, err := dbus.NewSystemConnectionContext(ctx) // 建立到 systemd 的连接
	if err != nil {
		return fmt.Errorf("failed to connect to systemd: %v", err)
	}
	defer conn.Close()

	return conn.ReloadContext(ctx)
}
