package collector

import (
	"io"
	"net"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"github.com/uesleicarvalhoo/go-collector-service/internal/domain/models"
	"github.com/uesleicarvalhoo/go-collector-service/pkg/logger"
	"golang.org/x/crypto/ssh"
)

type SFTPCollector struct {
	config       Config
	KeyExchanges []string
	sshClient    *ssh.Client
	sftpClient   *sftp.Client
}

func NewSFTPCollector(cfg Config, keyExchanges ...string) (*SFTPCollector, error) {
	client := &SFTPCollector{
		config:       cfg,
		KeyExchanges: keyExchanges,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (sc *SFTPCollector) GetFiles() (fileList []models.File, err error) {
	if err := sc.connect(); err != nil {
		return nil, err
	}

	collectedFiles, _ := sc.sftpClient.Glob(sc.config.MatchPattern)
	logger.Debugf("[SFTPCollector] Pattern '%s' find %d files\n", sc.config.MatchPattern, len(collectedFiles))

	for _, fp := range collectedFiles {
		fileName := filepath.Base(fp)
		fileList = append(
			fileList,
			models.NewFile(fileName, fp, sc),
		)
	}

	return fileList, nil
}

func (sc *SFTPCollector) GetFileReader(filePath string) (io.ReadSeekCloser, error) {
	if err := sc.connect(); err != nil {
		return nil, err
	}

	return sc.sftpClient.Open(filePath)
}

func (sc *SFTPCollector) RemoveFile(file models.File) error {
	if err := sc.connect(); err != nil {
		return err
	}

	return sc.sftpClient.Remove(file.FilePath)
}

func (sc *SFTPCollector) newFileModel(fp string) models.File {
	return models.NewFile(filepath.Base(fp), fp, sc)
}

func (sc *SFTPCollector) connect() error {
	if sc.sshClient != nil {
		_, _, err := sc.sshClient.SendRequest("keepalive", true, nil)
		if err == nil {
			return nil
		}
	}

	auth := ssh.Password(sc.config.Password)

	if sc.config.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(sc.config.PrivateKey))
		if err != nil {
			return errors.Wrapf(err, "ssh parse private key: %w", ErrConnectionFailed)
		}

		auth = ssh.PublicKeys(signer)
	}

	cfg := &ssh.ClientConfig{
		User:            sc.config.User,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: hostKeyCallback,
	}

	sshclient, err := ssh.Dial("tcp", sc.config.Server, cfg)
	if err != nil {
		return errors.Wrapf(err, "ssh dial: %w", ErrConnectionFailed)
	}

	SFTPCollector, err := sftp.NewClient(sshclient)
	if err != nil {
		return errors.Wrapf(err, "sftp new client: %w", ErrConnectionFailed)
	}

	sc.sshClient = sshclient
	sc.sftpClient = SFTPCollector

	return nil
}

func hostKeyCallback(_ string, _ net.Addr, _ ssh.PublicKey) error {
	return nil
}
