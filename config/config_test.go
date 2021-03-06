package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func tempFile(t *testing.T) string {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		t.Fatalf("tempfile: %v", err)
	}
	defer f.Close()
	return f.Name()
}

func TestConfigDefaults(t *testing.T) {
	fpath := tempFile(t)
	defer os.Remove(fpath)
	cfg, err := ReadFile(fpath)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if cfg.Node.Agent.SocksAddr != "127.0.0.1:9250" {
		t.Errorf("default node.agent.socksport = %d", cfg.Node.Agent.SocksAddr)
	}
	if cfg.Node.Agent.ControlAddr != "127.0.0.1:9251" {
		t.Errorf("default node.agent.controlport = %d", cfg.Node.Agent.ControlAddr)
	}
}

func TestRoundTrip(t *testing.T) {
	fpath := tempFile(t)
	defer os.Remove(fpath)
	config := Config{
		Dir:  filepath.Dir(fpath),
		Path: fpath,
		Node: Node{
			Agent: Agent{
				TorBinaryPath: "/usr/bin/tor",
				SocksAddr:     "127.0.0.1:9050",
				ControlAddr:   "127.0.0.1:9051",
			},
			Service: Service{
				Exports: []string{
					"127.0.0.1:80",
				},
				Clients: []Client{{
					Name:    "bob",
					Address: "qwertyuiop.onion",
				}},
			},
		},
	}
	err := WriteFile(&config, fpath)
	if err != nil {
		t.Fatalf("failed to write config: %v", err)
	}
	config2, err := ReadFile(fpath)
	if err != nil {
		t.Fatalf("failed to read config: %v", err)
	}
	assert.Equal(t, &config, config2)
}
