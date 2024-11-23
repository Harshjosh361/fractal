package integrations

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/logger"
	"github.com/SkySingh04/fractal/registry"
	"github.com/jlaffaye/ftp"
)

// FTPSource implements the DataSource interface
type FTPSource struct {
	URL         string `json:"url"`
	User        string `json:"user"`
	Password    string `json:"password"`
	FTPFILEPATH string `json:"file_path"`
}

// FTPDestination implements the DataDestination interface
type FTPDestination struct {
	URL         string `json:"url"`
	User        string `json:"user"`
	Password    string `json:"password"`
	FTPFILEPATH string `json:"file_path"`
}

// FetchData fetches data from an FTP server
func (f FTPSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateFTPRequest(req, true); err != nil {
		return nil, err
	}
	logger.Infof("Connecting to FTP server at %s...", req.FTPURL)

	conn, err := dialFTP(req.FTPURL, req.FTPUser, req.FTPPassword)
	if err != nil {
		return nil, err
	}
	defer conn.Quit()

	logger.Infof("Downloading file from FTP: %s", req.FTPFILEPATH)
	resp, err := conn.Retr(req.FTPFILEPATH)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve file from FTP: %w", err)
	}
	defer resp.Close()

	data, err := io.ReadAll(resp)
	if err != nil {
		return nil, fmt.Errorf("failed to read data from FTP response: %w", err)
	}

	logger.Infof("Successfully fetched data from FTP.")
	return data, nil
}

// SendData sends data to an FTP server
func (f FTPDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateFTPRequest(req, false); err != nil {
		return err
	}
	logger.Infof("Connecting to FTP server at %s...", req.FTPURL)

	conn, err := dialFTP(req.FTPURL, req.FTPUser, req.FTPPassword)
	if err != nil {
		return err
	}
	defer conn.Quit()

	logger.Infof("Uploading file to FTP: %s", req.FTPFILEPATH)
	dataBytes, ok := data.([]byte)
	if !ok {
		return fmt.Errorf("invalid data format; expected []byte, got %T", data)
	}

	err = conn.Stor(req.FTPFILEPATH, bytes.NewReader(dataBytes))
	if err != nil {
		return fmt.Errorf("failed to store file to FTP: %w", err)
	}

	logger.Infof("Successfully sent data to FTP.")
	return nil
}

// dialFTP creates and authenticates an FTP connection
func dialFTP(url, user, password string) (*ftp.ServerConn, error) {
	// Remove "ftp://" prefix if present
	url = strings.TrimPrefix(url, "ftp://")

	conn, err := ftp.Dial(url, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FTP server: %w", err)
	}

	err = conn.Login(user, password)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with FTP server: %w", err)
	}
	return conn, nil
}

func init() {
	registry.RegisterSource("FTP", FTPSource{})
	registry.RegisterDestination("FTP", FTPDestination{})
}

// Predefined FTP errors
var (
	ErrFTPConnectionFailed = errors.New("failed to connect to FTP server")
	ErrFTPLoginFailed      = errors.New("failed to login to FTP server")
	ErrFTPFileNotFound     = errors.New("file not found on FTP server")
	ErrFTPFileUploadFailed = errors.New("failed to upload file to FTP server")
)

// validateFTPRequest validates the request fields for FTP
func validateFTPRequest(req interfaces.Request, isSource bool) error {
	if req.FTPURL == "" {
		return errors.New("missing FTP URL")
	}
	if req.FTPUser == "" {
		return errors.New("missing FTP user")
	}
	if req.FTPPassword == "" {
		return errors.New("missing FTP password")
	}
	if req.FTPFILEPATH == "" {
		return errors.New("missing file path")
	}
	if !strings.HasPrefix(req.FTPURL, "ftp://") {
		return fmt.Errorf("invalid FTP URL: %s", req.FTPURL)
	}
	return nil
}
