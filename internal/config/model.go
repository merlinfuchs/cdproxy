package config

type Config struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	PublicURL       string `yaml:"public_url"`
	DBFileName      string `yaml:"db_file_name"`
	DownloadTimeout int    `yaml:"download_timeout"`

	SFTPUser               string `yaml:"sftp_user"`
	SFTPPassword           string `yaml:"sftp_password"`
	SFTPHost               string `yaml:"sftp_host"`
	BrotliCompressionLevel int    `yaml:"brotli_compression_level"`

	DefaultMaxSize        int `yaml:"default_max_size"`
	DefaultExpiry         int `yaml:"default_expiry"`
	DefaultOriginalExpiry int `yaml:"default_original_expiry"`

	MaxQueueSize int `yaml:"max_queue_size"`
	NumWorkers   int `yaml:"num_workers"`
}
