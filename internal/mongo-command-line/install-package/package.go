package install_package

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"fmt"
)

var i = &InstallerMapping{
	Mapping: make(map[string]map[string]contract.DbPackage),
}

// InstallerMapping 表示系统和数据库版本与安装包之间的映射
type InstallerMapping struct {
	Mapping map[string]map[string]contract.DbPackage
}

// RegistryDbPackages 初始化一个 InstallerMapping
func RegistryDbPackages(d contract.DbInstallPackageInfo) {
	i.AddMapping(d)
}

func GetDbPackageMapping() *InstallerMapping {
	return i
}

// AddMapping 添加一个系统、数据库版本和安装包路径的映射
func (im *InstallerMapping) AddMapping(d contract.DbInstallPackageInfo) {
	sysKey := fmt.Sprintf("%s-%s", d.S.Name, d.S.Version)
	dbKey := fmt.Sprintf("%s-%s", d.D.Name, d.D.Version)

	// 初始化嵌套映射
	if _, ok := im.Mapping[sysKey]; !ok {
		im.Mapping[sysKey] = make(map[string]contract.DbPackage)
	}
	im.Mapping[sysKey][dbKey] = d.P
}

// GetPackage 查找对应的安装包路径
func (im *InstallerMapping) GetPackage(s contract.SystemInfo, d contract.DatabaseInfo) (contract.DbPackage, bool) {
	var p = contract.DbPackage{}
	sysKey := fmt.Sprintf("%s-%s", s.Name, s.Version)
	dbKey := fmt.Sprintf("%s-%s", d.Name, d.Version)

	if dbMap, ok := im.Mapping[sysKey]; ok {
		if p, ok = dbMap[dbKey]; ok {
			return p, true
		}
	}
	return p, false
}

func PrintPackageInfo() map[string][]string {

	var l = make(map[string][]string)

	for sysKey, v := range i.Mapping {
		for dbKey, _ := range v {
			l[sysKey] = append(l[sysKey], dbKey)

		}
	}
	return l
}

func DbPackageInfoFormat(s contract.SystemInfo, d contract.DatabaseInfo) string {
	sysKey := fmt.Sprintf("%s-%s", s.Name, s.Version)
	dbKey := fmt.Sprintf("%s-%s", d.Name, d.Version)
	return fmt.Sprintf("%s:%s", sysKey, dbKey)
}
