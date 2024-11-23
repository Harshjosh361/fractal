package integrations

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPSource implements the DataSource interface
type SFTPSource struct {
	URL          string `json:"url"`
	User         string `json:"user"`
	Password     string `json:"password"`
	SFTPFILEPATH string `json:"file_path"`
}

// SFTPDestination implements the DataDestination interface
type SFTPDestination struct {
	URL          string `json:"url"`
	User         string `json:"user"`
	Password     string `json:"password"`
	SFTPFILEPATH string `json:"file_path"`
}

// FetchData fetches data from an SFTP server
func (s SFTPSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateSFTPRequest(req, true); err != nil {
		return nil, err
	}
	logger.Infof("Connecting to SFTP server at %s...", req.SFTPURL)

	client, err := dialSFTP(req.SFTPURL, req.SFTPUser, req.SFTPPassword)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	logger.Infof("Downloading file from SFTP: %s", req.SFTPFILEPATH)
	file, err := client.Open(req.SFTPFILEPATH)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file from SFTP: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from SFTP response: %w", err)
	}

	logger.Infof("Successfully fetched data from SFTP.")
	return data, nil
}

// SendData sends data to an SFTP server
func (s SFTPDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateSFTPRequest(req, false); err != nil {
		return err
	}
	logger.Infof("Connecting to SFTP server at %s...", req.SFTPURL)

	client, err := dialSFTP(req.SFTPURL, req.SFTPUser, req.SFTPPassword)
	if err != nil {
		return err
	}
	defer client.Close()

	logger.Infof("Uploading file to SFTP: %s", req.SFTPFILEPATH)
	dataBytes, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("invalid data format; expected []byte, got %T", data)
	}

	file, err := client.Create(req.SFTPFILEPATH)
	if err != nil {
		return fmt.Errorf("failed to create file on SFTP server: %w", err)
	}
	defer file.Close()

	_, err = file.Write(dataBytes)
	if err != nil {
		return fmt.Errorf("failed to write file to SFTP server: %w", err)
	}

	logger.Infof("Successfully sent data to SFTP.")
	return nil
}

// dialSFTP creates and authenticates an SFTP connection
func dialSFTP(url, user, password string) (*sftp.Client, error) {
	// Remove "sftp://" prefix if present
	url = strings.TrimPrefix(url, "sftp://")

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	conn, err := ssh.Dial("tcp", url, config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SFTP server: %w", err)
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}

	return client, nil
}

// validateSFTPRequest validates the request fields for SFTP
func validateSFTPRequest(req interfaces.Request, isSource bool) error {
	if req.SFTPURL == "" {
		return errors.New("missing SFTP URL")
	}
	if req.SFTPUser == "" {
		return errors.New("missing SFTP user")
	}
	if req.SFTPPassword == "" {
		return errors.New("missing SFTP password")
	}
	if req.SFTPFILEPATH == "" {
		return errors.New("missing file path")
	}
	if !strings.HasPrefix(req.SFTPURL, "sftp://") {
		return fmt.Errorf("invalid SFTP URL: %s", req.SFTPURL)
	}
	return nil
}

func init() {
	registry.RegisterSource("SFTP", SFTPSource{})
	registry.RegisterDestination("SFTP", SFTPDestination{})
}
