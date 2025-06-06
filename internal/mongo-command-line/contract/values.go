package contract

const (
	// Xfs file system
	Xfs = 1481003842

	MongoUser = "mongod"

	PackagePath = "/data"

	Profile    = "/etc/profile"
	ProfileEnv = "$PATH:/usr/local/mongodb/bin"

	LinkPath = "/usr/local/mongodb"

	MongoDPath = "/usr/local/mongodb/bin/mongod"
	MongoSPath = "/usr/local/mongodb/bin/mongos"
	MongoPath  = "/usr/local/mongodb/bin/mongo"

	BaseDataDirFormat = "/data1/mg%d/"
	DropDataDirFormat = "/data1/mg%d_%s"

	// SystemdConfigPath /usr/lib/systemd/system/mongod_27017.service
	SystemdConfigPath = "/usr/lib/systemd/system"
	//SystemdHugepageConfigPath  = "/lib/systemd/system/"
	SystemdHugepageServiceName = "mongodb-hugepage-fix.service"

	SystemdMongoDServiceName = "mongod_%d.service"
	SystemdMongoSServiceName = "mongos_%d.service"

	MongoDataDirFormat = "data%d"
	MongoLogDirFormat  = "log%d"
)

const (
	MongoS              = "mongos"
	Repl                = "replset"
	MongoD              = "mongod"
	Cluster             = "cluster"
	DBName              = "mongo"
	Mongo3_4_24         = "3.4.24"
	Mongo4_2_10         = "4.2.10"
	CentOsLinuxName     = "CentOS Linux"
	CentOsLinuxVersion7 = "7"
)
