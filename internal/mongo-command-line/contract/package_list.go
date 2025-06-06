package contract

var (
	MongoVersionList = []DbInstallPackageInfo{
		{
			S: SystemInfo{
				Name:    CentOsLinuxName,
				Version: CentOsLinuxVersion7,
			},
			D: DatabaseInfo{
				Name:    DBName,
				Version: Mongo3_4_24,
			},
			P: DbPackage{
				DirName:     "mongodb-linux-x86_64-CentOS7-3.4.24",
				PackageName: "mongodb-linux-x86_64-CentOS7-3.4.24.tgz",
				DownLoadURL: "http://opdownload.zlongame.com/DBA/mongodb-linux-x86_64-CentOS7-3.4.24.tgz",
			},
		},
		{
			S: SystemInfo{
				Name:    CentOsLinuxName,
				Version: CentOsLinuxVersion7,
			},
			D: DatabaseInfo{
				Name:    DBName,
				Version: Mongo4_2_10,
			},
			P: DbPackage{
				DirName:     "mongodb-linux-x86_64-CentOS7-4.2.10",
				PackageName: "mongodb-linux-x86_64-CentOS7-4.2.10.tgz",
				DownLoadURL: "http://opdownload.zlongame.com/DBA/mongodb-linux-x86_64-CentOS7-4.2.10.tgz",
			},
		},
	}
)

type DbInstallPackageInfo struct {
	S SystemInfo
	D DatabaseInfo
	P DbPackage
}

// SystemInfo 表示系统的版本信息
type SystemInfo struct {
	Name    string // e.g., "CentOS Linux"
	Version string // e.g., "7"
}

// DatabaseInfo 表示数据库的版本信息
type DatabaseInfo struct {
	Name    string // e.g., "mongo"
	Version string // e.g., "3.4.24"
}

type DbPackage struct {
	DirName     string // e.g., "mongodb-linux-x86_64-CentOS7-3.4.24"
	PackageName string // e.g., "mongodb-linux-x86_64-CentOS7-3.4.24.tgz"
	DownLoadURL string // e.g., "http://opdownload.zlongame.com/DBA/mongodb-linux-x86_64-CentOS7-3.4.24.tgz"
}
