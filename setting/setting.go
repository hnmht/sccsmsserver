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
	StartTime       string                                        `mapstructure:"start_time" json:"startTime"`            // Timestamp when the application started. If you are using multiple application servers, you can use this field to record the start time each one
	MachineID       int64                                         `mapstructure:"machine_id" json:"machineID"`            // If you're using multiple application servers, you can use this field to record the machine identifier for each one
	TLS             bool                                          `mapstructure:"tls" json:"TLS"`                         // Enable TLS or not
	CertificateFile string                                        `mapstructure:"certificatefile" json:"certificateFile"` // TLS client certificate file name, required when TLS is enabled
	PrivateKeyFile  string                                        `mapstructure:"privatekeyfile" json:"privateKeyFile"`   // TLS private key certificate file name, required when TLS is enabled
	UserLockTh      int32                                         `mapstructure:"userlockth" json:"userLockTh"`           // Maximum number of allowed password failures within 30 minutes; exceeding this will lock the account
	IpLockTh        int32                                         `mapstructure:"iplockth" json:"ipLockTh"`               // Maxium number of login attempts with non-existent usernames within 30 minutes; exceeding this will result in the IP address being locked
	IpLockedMinutes int32                                         `mapstructure:"iplockedminutes" json:"ipLockedMinutes"` // Duration of IP address lock (in minutes)
	*LogConfig      `mapstructure:"log" json:"log"`               // Application log configuration
	*PqConfig       `mapstructure:"postgresql" json:"postgresql"` // Application PostgreSQL database configuration
	*S3Storage      `mapstructure:"s3storage" json:"s3storage"`   // AWS s3 object storage server configuration
	*RedisConfig    `mapstructure:"redis" json:"redis"`           // Redis cache configuration
}

// Application's log configuration structure
type LogConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	Filename   string `mapstructure:"filename" json:"filename"`
	MaxSize    int    `mapstructure:"max_size" json:"maxSize"`
	MaxAge     int    `mapstructure:"max_age" json:"maxAge"`
	MaxBackups int    `mapstructure:"max_backups" json:"maxBackups"`
}

// PostgreSQL database configuration
type PqConfig struct {
	Host         string `mapstructure:"host" json:"host"`
	Port         int    `mapstructure:"port" json:"port"`
	DbName       string `mapstructure:"dbname" json:"dbName"`
	Username     string `mapstructure:"username" json:"username"`
	Password     string `mapstructure:"password" json:"password"`
	MaxOpenConns int    `mapstructure:"max_open_conns" json:"maxOpenConns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" json:"maxIdleConns"`
	MaxRecord    int32  `mapstructure:"max_record" json:"maxRecord"`
}

// Minio Object Storage Configuration
type S3Storage struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint"`
	AccessKeyID     string `mapstructure:"accesskeyid" json:"accessKeyID"`
	SecretAccessKey string `mapstructure:"secretaccesskey" json:"secretAccessKey"`
	Secure          bool   `mapstructure:"secure" json:"secure"`
	SelfSigned      bool   `mapstructure:"selfsigned" json:"selfSigned"`
	DefaultBucket   string `mapstructure:"defaultbucket" json:"defaultBucket"`
	Location        string `mapstructure:"location" json:"location"`
}

// Redis cache configuration
type RedisConfig struct {
	Enabled  bool   `mapstructure:"enabled" json:"enabled"` // Define if Redis caching is enabled. This is only necessary for distributed and cluster deployments.
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Db       int    `mapstructure:"db" json:"db"`
	Password string `mapstructure:"password" json:"password"`
	PoolSize int    `mapstructure:"pool_size" json:"poolSize"`
}

// Read global configuration
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
	fmt.Println("Global configuration read successfully.")
	return
}
