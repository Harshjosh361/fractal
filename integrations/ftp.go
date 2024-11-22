package integrations

import (
	"errors"
	"log"

	"github.com/SkySingh04/fractal/interfaces"
	"github.com/SkySingh04/fractal/registry"
	"github.com/jlaffaye/ftp"
)

// FTPSource implements the DataSource interface
type FTPSource struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// FTPDestination implements the DataDestination interface
type FTPDestination struct {
	URL      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// FetchData fetches data from an FTP server
func (f FTPSource) FetchData(req interfaces.Request) (interface{}, error) {
	if err := validateFTPRequest(req, true); err != nil {
		return nil, err
	}
	log.Println("Fetching data from FTP...")
	conn, err := ftp.Dial(req.FTPURL)
	if err != nil {
		return nil, err
	}
	err = conn.Login(req.FTPUser, req.FTPPassword)
	if err != nil {
		return nil, err
	}
	// Fetch the file or data from FTP
	return "FTPData", nil
}

// SendData sends data to an FTP server
func (f FTPDestination) SendData(data interface{}, req interfaces.Request) error {
	if err := validateFTPRequest(req, false); err != nil {
		return err
	}
	log.Println("Sending data to FTP...")
	conn, err := ftp.Dial(req.FTPURL)
	if err != nil {
		return err
	}
	err = conn.Login(req.FTPUser, req.FTPPassword)
	if err != nil {
		return err
	}
	// Send the file or data to FTP
	return nil
}

// validateFTPRequest validates the request fields for FTP
func validateFTPRequest(req interfaces.Request, isSource bool) error {
	if isSource && req.FTPURL == "" {
		return errors.New("missing FTP URL for source")
	}
	if !isSource && req.FTPURL == "" {
		return errors.New("missing FTP URL for destination")
	}
	if req.FTPUser == "" {
		return errors.New("missing FTP user")
	}
	if req.FTPPassword == "" {
		return errors.New("missing FTP password")
	}
	return nil
}

func init() {
	registry.RegisterSource("FTP", FTPSource{})
	registry.RegisterDestination("FTP", FTPDestination{})
}
