package settings

type AppConfig struct {
	DBHost, DBName, DBUser, DBPassword, LogPrefix, UploadsDir, ApiPathPrefix, PublicPathPrefix, PublicDir string
}

var Config *AppConfig
