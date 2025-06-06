package contract

import (
	"context"
	"github.com/coreos/go-systemd/v22/dbus"
)

type InstallModule interface {
	Prepare
	BaseInstall
}

type Prepare interface {
	PackageCheck
	DirCheck
	CreateDir
}

type BaseInstall interface {
	AddUser
	FileCheck
	DownloadFile
	Extract
	CreateSymlink
	AddEnv
	CreateDir
	FileOwner
	StartService
	ReloadService
	EnableService
}

type SystemdStartService interface {
	StartService
	EnableService
	ReloadService
}

type CreateDir interface {
	CreateDir(path string) error
}

type DirCheck interface {
	CheckDirEmpty(path string) error
}

type FileCheck interface {
	CheckFileExist(path string) error
}

type PackageCheck interface {
	GetPackage(databaseInfo DatabaseInfo) (DbPackage, error)
}

type AddUser interface {
	AddNoLoginUser(username, shell string) error
}

type DownloadFile interface {
	DownloadFile(ctx context.Context, url string, descFilepath string) error
}

type Extract interface {
	Extract(ctx context.Context, path, descDir string) error
}

type CreateSymlink interface {
	CreateSymlink(target, linkName string) error
}

type AddEnv interface {
	AddEnvToFile(profileFile, variable, value string) error
}

type FileOwner interface {
	SetUserOwner(uName, path string) error
}

type StartService interface {
	StartService(ctx context.Context, serviceName string) error
}

type StopService interface {
	StopService(ctx context.Context, serviceName string) error
}

type EnableService interface {
	EnableService(ctx context.Context, unitFiles []string) error
}

type DisableService interface {
	DisableService(ctx context.Context, unitFiles []string) error
}

type ReloadService interface {
	SystemdReload(ctx context.Context) error
}

type StatusService interface {
	StatusService(ctx context.Context, unit, property string) (*dbus.Property, error)
}
