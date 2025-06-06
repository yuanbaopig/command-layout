package common

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/mongo-command-line/module/ini"
	"DatabaseManage/internal/pkg/log"
	"context"
)

type SystemdServiceConfig struct {
	FileName    string
	ServiceName string
	Config      interface{}
}

func ServiceStartAndEnable(ctx context.Context, c SystemdServiceConfig, s contract.SystemdStartService) error {

	if err := ini.CreateIniConfig(c.FileName, c.Config); err != nil {
		log.Debug(err)
		return err
	}

	log.Debug("systemd daemon reload")
	if err := s.SystemdReload(ctx); err != nil {
		log.Debug(err)
		return err
	}

	log.Debugf("start %s service", c.ServiceName)
	if err := s.StartService(ctx, c.ServiceName); err != nil {
		log.Debug(err)
		return err
	}

	// 加入自启动
	if err := s.EnableService(ctx, []string{c.ServiceName}); err != nil {
		log.Debug(err)
		return err
	}
	return nil
}
