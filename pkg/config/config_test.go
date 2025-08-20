// pkg/config/config_test.go
package config_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"crypto-observer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestMustLoad_OK(t *testing.T) {
	// создаём временную директорию, где будет ./configs/config.yaml
	td := t.TempDir()
	cfgDir := filepath.Join(td, "configs")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	// валидный YAML
	yaml := `
server:
  addr: ":9090"
db:
  dsn: "postgres://user:pass@localhost:5432/crypto?sslmode=disable"
collector:
  default_period_seconds: 7
coingecko:
  base_url: "https://api.coingecko.com/api/v3"
  timeout_s: 3
log:
  level: "debug"
`
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(yaml), 0o644))

	// меняем рабочую директорию на temp, чтобы относительный путь "configs/config.yaml" совпал
	chdir(t, td)

	// загрузка
	got := config.MustLoad()
	require.NotNil(t, got)

	// проверки
	require.Equal(t, ":9090", got.Server.Addr)
	require.Equal(t, "postgres://user:pass@localhost:5432/crypto?sslmode=disable", got.DB.DSN)
	require.Equal(t, 7, got.Collector.DefaultPeriodSeconds)
	require.Equal(t, "https://api.coingecko.com/api/v3", got.Coingecko.BaseURL)
	require.Equal(t, 3, got.Coingecko.TimeoutSec)
	require.Equal(t, "debug", got.Log.Level)

	// убедимся, что глобальный getter возвращает тот же объект
	require.Equal(t, got, config.C())
}

func TestMustLoad_NoFile_Panics(t *testing.T) {
	td := t.TempDir()
	chdir(t, td) // здесь нет каталога configs/ — ожидаем панику
	require.Panics(t, func() {
		_ = config.MustLoad()
	}, "ожидалась паника при отсутствии configs/config.yaml")
}

func TestMustLoad_BadYAML_Panics(t *testing.T) {
	td := t.TempDir()
	cfgDir := filepath.Join(td, "configs")
	require.NoError(t, os.MkdirAll(cfgDir, 0o755))

	// битый YAML
	bad := `server: { addr: ":8080"  db: { dsn: "x" }` // пропущены закрывающие скобки/отступы
	require.NoError(t, os.WriteFile(filepath.Join(cfgDir, "config.yaml"), []byte(bad), 0o644))

	chdir(t, td)
	require.Panics(t, func() {
		_ = config.MustLoad()
	}, "ожидалась паника при невалидном YAML")
}

// chdir меняет рабочую директорию и автоматически возвращает её назад в конце теста
func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(dir))
	t.Cleanup(func() {
		_ = os.Chdir(old)
	})
}

func repoRoot() string {
	_, file, _, _ := runtime.Caller(0)
	// pkg/config/config_test.go -> подняться на 2 уровня
	return filepath.Clean(filepath.Join(filepath.Dir(file), "../.."))
}
