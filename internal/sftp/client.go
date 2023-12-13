package sftp

import (
	"fmt"
	"io"
	"os"

	"github.com/merlinfuchs/cdproxy/internal/config"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPClient struct {
	client *sftp.Client
}

func New() (*SFTPClient, error) {
	var auths []ssh.AuthMethod
	if config.C.SFTPPassword != "" {
		auths = append(auths, ssh.Password(config.C.SFTPPassword))
	}

	cfg := ssh.ClientConfig{
		User:            config.C.SFTPUser,
		Auth:            auths,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", config.C.SFTPHost, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to connect: %w", err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("unable to start sftp subsystem: %w", err)
	}

	return &SFTPClient{
		client: client,
	}, nil
}

func (c *SFTPClient) Close() {
	c.client.Close()
}

func (c *SFTPClient) ReadFile(path string) ([]byte, error) {
	file, err := c.client.OpenFile(path, os.O_RDONLY)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	body, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return body, nil
}

func (c *SFTPClient) WriteFile(path string, body []byte) error {
	file, err := c.client.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (c *SFTPClient) DeleteFile(path string) error {
	err := c.client.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}
