package s3

import (
	"os"

	"github.com/minio/minio-go/v7/pkg/credentials"
)

// FileSecretsCredentials is an implementation of Minio's credentials.Provider,
// allowing to read credentials from Kubernetes or Docker secrets, as described in
// https://kubernetes.io/docs/tasks/inject-data-application/distribute-credentials-secure
// and https://docs.docker.com/engine/swarm/secrets.
type FileSecretsCredentials struct {
	// Path to the fiel containing the access key,
	// e.g. /etc/s3-secrets/access-key.
	AccessKeyFilename string `json:"access_key_file"`

	// Path to the fiel containing the secret key,
	// e.g. /etc/s3-secrets/secret-key.
	SecretKeyFilename string `json:"secret_key_file"`
}

// IsExpired implements credentials.Provider.
// As there is no totally reliable way to tell
// if a file was modified accross all filesystems except opening it,
// we always return true, and p.Retrieve will open it regardless.
func (*FileSecretsCredentials) IsExpired() bool {
	return true
}

// Retrieve implements credentials.Provider.
// It reads files pointed to by p.AccessKeyFilename and p.SecretKeyFilename.
func (c *FileSecretsCredentials) Retrieve() (value credentials.Value, err error) {
	load := func(filename string, dst *string) error {
		var b []byte
		b, err = os.ReadFile(filename)
		if err != nil {
			return err
		}

		*dst = string(b)
		return nil
	}

	err = load(c.AccessKeyFilename, &value.AccessKeyID)
	if err != nil {
		return value, err
	}
	err = load(c.SecretKeyFilename, &value.SecretAccessKey)

	return value, err
}

var _ credentials.Provider = new(FileSecretsCredentials)
