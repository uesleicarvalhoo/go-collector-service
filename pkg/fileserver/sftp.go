package fileserver

import (
	"context"
	"io"
	"io/fs"
	"net"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTPFileServer struct {
	config       Config
	KeyExchanges []string
	sshClient    *ssh.Client
	sftpClient   *sftp.Client
}

func NewSFTP(cfg Config, keyExchanges ...string) (*SFTPFileServer, error) {
	client := &SFTPFileServer{
		config:       cfg,
		KeyExchanges: keyExchanges,
	}

	if err := client.connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (fs *SFTPFileServer) Glob(ctx context.Context, pattern string) ([]string, error) {
	if err := fs.connect(); err != nil {
		return nil, err
	}

	files := []string{}

	matchs, err := fs.sftpClient.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, match := range matchs {
		if f, _ := fs.sftpClient.Stat(match); !f.IsDir() {
			files = append(files, match)
		}
	}

	return files, nil
}

func (fs *SFTPFileServer) Open(ctx context.Context, filePath string) (io.ReadSeekCloser, error) {
	if err := fs.connect(); err != nil {
		return nil, err
	}

	return fs.sftpClient.Open(filePath)
}

func (fs *SFTPFileServer) Remove(ctx context.Context, filePath string) error {
	if err := fs.connect(); err != nil {
		return err
	}

	return fs.sftpClient.Remove(filePath)
}

func (fs *SFTPFileServer) Move(ctx context.Context, oldname, newname string) error {
	if err := fs.connect(); err != nil {
		return err
	}

	dirName, _ := filepath.Split(newname)
	if err := fs.sftpClient.MkdirAll(dirName); err != nil {
		return err
	}

	if err := fs.sftpClient.MkdirAll(dirName); err != nil {
		return err
	}

	return fs.sftpClient.Rename(oldname, newname)
}

func (fs *SFTPFileServer) Stat(ctx context.Context, filePath string) (fs.FileInfo, error) {
	if err := fs.connect(); err != nil {
		return nil, err
	}

	return fs.sftpClient.Stat(filePath)
}

func (fs *SFTPFileServer) Lock(ctx context.Context, filePath string) error {
	return nil
}

func (fs *SFTPFileServer) Unlock(ctx context.Context, filepath string) error {
	return nil
}

func (fs *SFTPFileServer) connect() error {
	if fs.sshClient != nil {
		_, _, err := fs.sshClient.SendRequest("keepalive", true, nil)
		if err == nil {
			return nil
		}
	}

	auth := ssh.Password(fs.config.Password)

	if fs.config.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(fs.config.PrivateKey))
		if err != nil {
			return errors.Wrapf(err, "ssh parse private key: %s", ErrConnectionFailed)
		}

		auth = ssh.PublicKeys(signer)
	}

	cfg := &ssh.ClientConfig{
		User:            fs.config.User,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: hostKeyCallback,
	}

	sshclient, err := ssh.Dial("tcp", fs.config.Server, cfg)
	if err != nil {
		return errors.Wrapf(err, "ssh dial: %s", ErrConnectionFailed)
	}

	SFTPFileServer, err := sftp.NewClient(sshclient)
	if err != nil {
		return errors.Wrapf(err, "sftp new client: %s", ErrConnectionFailed)
	}

	fs.sshClient = sshclient
	fs.sftpClient = SFTPFileServer

	return nil
}

func hostKeyCallback(_ string, _ net.Addr, _ ssh.PublicKey) error {
	return nil
}
