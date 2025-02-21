// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package clickhouselogsexporter

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configretry"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/otelcol/otelcoltest"
)

func TestLoadConfig(t *testing.T) {
	factories, err := otelcoltest.NopFactories()
	require.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[component.MustNewType(typeStr)] = factory
	cfg, err := otelcoltest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, len(cfg.Exporters), 3)

	defaultCfg := factory.CreateDefaultConfig()
	defaultCfg.(*Config).DSN = "tcp://127.0.0.1:9000/?dial_timeout=5s"
	r0 := cfg.Exporters[component.NewID(component.MustNewType(typeStr))]
	assert.Equal(t, r0, defaultCfg)

	r1 := cfg.Exporters[component.NewIDWithName(component.MustNewType(typeStr), "full")].(*Config)
	assert.Equal(t, r1, &Config{
		DSN: "tcp://127.0.0.1:9000/?dial_timeout=5s",
		TimeoutConfig: exporterhelper.TimeoutConfig{
			Timeout: 5 * time.Second,
		},
		BackOffConfig: configretry.BackOffConfig{
			Enabled:             true,
			InitialInterval:     5 * time.Second,
			MaxInterval:         30 * time.Second,
			MaxElapsedTime:      300 * time.Second,
			RandomizationFactor: 0.7,
			Multiplier:          1.3,
		},
		QueueConfig: exporterhelper.QueueConfig{
			Enabled:      true,
			NumConsumers: 10,
			QueueSize:    100,
		},
		AttributesLimits: AttributesLimits{
			FetchKeysInterval: 10 * time.Minute,
			MaxDistinctValues: 25000,
		},
	})

	defaultCfg.(*Config).UseNewSchema = true
	r2 := cfg.Exporters[component.NewIDWithName(component.MustNewType(typeStr), "new_schema")]
	assert.Equal(t, r2, defaultCfg)
}
