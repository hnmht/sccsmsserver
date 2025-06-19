package setting

import (
	"fmt"

	"github.com/spf13/viper"
)

// Global scope variable
var Conf = new(AppConfig)

// AppConfig defines the entire application's configuration structure
type AppConfig struct {
	Name            string                                        `mapstructure:"name" json:"name"`                       // Application's name
	Addr            string                                        `mapstructure:"addr" json:"addr"`                       // IP address of the application
	Mode            string                                        `mapstructure:"mode" json:"mode"`                       // The application's run mode, can be set to "debug" or "release"
	Port            int32                                         `mapstructure:"port" json:"port"`                       // Port number of the application
	StartTime       string                                        `mapstructure:"start_time" json:"start_time"`           // Timestamp when the application started. If you are using multiple application servers, you can use this field to record the start time each one
	MachineID       int64                                         `mapstructure:"machine_id" json:"machine_id"`           // If you're using multiple application servers, you can use this field to record the machine identifier for each one
	TLS             bool                                          `mapstructure:"tls" json:"tls"`                         // Enable TLS or not
	CertificateFile string                                        `mapstructure:"certificatefile" json:"certificatefile"` // TLS client certificate file name, required when TLS is enabled
	PrivateKeyFile  string                                        `mapstructure:"privatekeyfile" json:"privatekeyfile"`   // TLS private key certificate file name, required when TLS is enabled
	UserLockTh      int32                                         `mapstructure:"userlockth" json:"userlockth"`           // Maximum number of allowed password failures within 30 minutes; exceeding this will lock the account
	IpLockTh        int32                                         `mapstructure:"iplockth" json:"iplockth"`               // Maxium number of login attempts with non-existent usernames within 30 minutes; exceeding this will result in the IP address being locked
	IpLockedMinutes int32                                         `mapstructure:"iplockedminutes" json:"iplockedminutes"` // Duration of IP address lock (in minutes)
	*LogConfig      `mapstructure:"log" json:"log"`               // Application log configuration
	*PqConfig       `mapstructure:"postgresql" json:"postgresql"` // Application PostgreSQL database configuration
	*MinioConfig    `mapstructure:"minio" json:"minio"`           // Minio object storage server configuration
	*RedisConfig    `mapstructure:"redis" json:"redis"`           // Redis cache configuration
}

// Application's log configuration structure
type LogConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	Filename   string `mapstructure:"filename" json:"filename"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
}

// PostgreSQL database configuration
type PqConfig struct {
	Host         string `mapstructure:"host" json:"host"`
	Port         int    `mapstructure:"port" json:"port"`
	DbName       string `mapstructure:"dbname" json:"dbname"`
	Username     string `mapstructure:"username" json:"username"`
	Password     string `mapstructure:"password" json:"password"`
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"max_idle_conns"`
	MaxRecord    int32  `mapstructure:"max_record" json:"max_record"`
}

// Minio Object Storage Configuration
type MinioConfig struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint"`
	AccessKeyID     string `mapstructure:"accesskeyid" json:"accesskeyid"`
	SecretAccessKey string `mapstructure:"secretaccesskey" json:"secretaccesskey"`
	Secure          bool   `mapstructure:"secure" json:"secure"`
	SelfSigned      bool   `mapstructure:"selfsigned" json:"selfsigned"`
	DefaultBucket   string `mapstructure:"defaultbucket" json:"defaultbucket"`
}

// Redis cache configuration
type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled" json:"enabled"` // Define if Redis caching is enabled. This is only necessary for distributed and cluster deployments.
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Db       int    `mapstructure:"db" json:"db"`
	Password string `mapstructure:"password" json:"password"`
	PoolSize int    `mapstructure:"pool_size" json:"pool_size"`
}

// Global variable initialization
func Init() (err error) {
	viper.SetConfigFile("config.yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("Init viper.ReadInConfig failed:", err)
		return
	}

	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Println("Init viper.Unmarshal(Conf) failed:", err)
		return
	}

	return
}
