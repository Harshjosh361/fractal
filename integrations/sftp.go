package integrations

import (
	"errors"
	"log"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SFTPSource implements the DataSource interface
type SFTPSource struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// SFTPDestination implements the DataDestination interface
type SFTPDestination struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// FetchData fetches data from an SFTP server
func (s SFTPSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateSFTPRequest(req, true); err != nil {
		return nil, err
	}
	log.Println("Fetching data from SFTP...")
	_, err := connectSFTP(req.SFTPURL, req.SFTPUser, req.SFTPPassword)
	if err != nil {
		return nil, err
	}
	// Implement SFTP file fetching logic
	return "SFTPData", nil
}

// SendData sends data to an SFTP server
func (s SFTPDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateSFTPRequest(req, false); err != nil {
		return err
	}
	log.Println("Sending data to SFTP...")
	_, err := connectSFTP(req.SFTPURL, req.SFTPUser, req.SFTPPassword)
	if err != nil {
		return err
	}
	// Implement SFTP file sending logic
	return nil
}

// connectSFTP establishes an SFTP connection
func connectSFTP(url, user, password string) (*sftp.Client, error) {
	// SSH client setup
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", url, config)
	if err != nil {
		return nil, err
	}
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return nil, err
	}
	return sftpClient, nil
}

// validateSFTPRequest validates the request fields for SFTP
func validateSFTPRequest(req interfaces.Request, isSource bool) error {
	if isSource && req.SFTPURL == "" {
		return errors.New("missing SFTP URL for source")
	}
	if !isSource && req.SFTPURL == "" {
		return errors.New("missing SFTP URL for destination")
	}
	if req.SFTPUser == "" {
		return errors.New("missing SFTP user")
	}
	if req.SFTPPassword == "" {
		return errors.New("missing SFTP password")
	}
	return nil
}

func init() {
	registry.RegisterSource("SFTP", SFTPSource{})
	registry.RegisterDestination("SFTP", SFTPDestination{})
}
