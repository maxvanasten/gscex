package git

import (
	"fmt"
	"os"
	"os/exec"
)

type Client struct {
	repo   string
	branch string
	path   string
}

func New(repo, branch, path string) *Client {
	return &Client{
		repo:   repo,
		branch: branch,
		path:   path,
	}
}

func (c *Client) Clone() error {
	if _, err := os.Stat(c.path); err == nil {
		return fmt.Errorf("directory already exists: %s", c.path)
	}

	cmd := exec.Command("git", "clone", "--depth", "1", "--branch", c.branch, c.repo, c.path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) Pull() error {
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		return fmt.Errorf("directory does not exist: %s", c.path)
	}

	cmd := exec.Command("git", "-C", c.path, "pull", "origin", c.branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (c *Client) Exists() bool {
	_, err := os.Stat(c.path)
	return err == nil
}
