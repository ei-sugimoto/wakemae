package yaml_test

import (
	"testing"

	"github.com/ei-sugimoto/wakemae/internal/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		data    []byte
		want    *yaml.Config
		wantErr bool
	}{
		{
			name:    "empty",
			data:    []byte(``),
			want:    &yaml.DefaultConfig,
			wantErr: false,
		},
		{
			name: "valid",
			data: []byte(`dns:
  bind_address: 0.0.0.0:53
  upstream: 1.1.1.1:53
  timeout: 10s
`),
			want: &yaml.Config{
				DNS: yaml.DNS{
					BindAddress: "0.0.0.0:53",
					Upstream:    "1.1.1.1:53",
					Timeout:     "10s",
				},
			},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := yaml.Parse(test.data)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			if err != nil {
				t.Fatalf("Parse() error = %v, wantErr %v", err, test.wantErr)
			}

			assert.Equal(t, test.want.DNS.BindAddress, got.DNS.BindAddress)
			assert.Equal(t, test.want.DNS.Upstream, got.DNS.Upstream)
			assert.Equal(t, test.want.DNS.Timeout, got.DNS.Timeout)
		})
	}
}
