package common

type Creds struct {
	AccessKey string
	SecretKey string
	Region    string
}

func GetCredentials() Creds {
	return Creds{
		AccessKey: "XXXXXXXX",
		SecretKey: "XXXXXXXX",
		Region:    "us-east-1",
	}
}
