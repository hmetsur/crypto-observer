package logger

import (
	"os"
	"testing"

	"crypto-observer/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

// helper: форсируем re-init под нужный уровень
func reinit(t *testing.T, cfgLevel, envLevel string) {
	t.Helper()
	// правим конфиг в памяти, без файлов
	config.C().Log.Level = cfgLevel

	if envLevel == "" {
		_ = os.Unsetenv("LOG_LEVEL")
	} else {
		t.Setenv("LOG_LEVEL", envLevel)
	}

	// сбрасываем синглтон
	log = nil
	Init()
}

func Test_parseLevel_table(t *testing.T) {
	cases := []struct {
		in   string
		want logrus.Level
	}{
		{"debug", logrus.DebugLevel},
		{"DEBUG", logrus.DebugLevel},
		{"warn", logrus.WarnLevel},
		{"Warning", logrus.WarnLevel},
		{"error", logrus.ErrorLevel},
		{"fatal", logrus.FatalLevel},
		{"", logrus.InfoLevel},
		{"unknown", logrus.InfoLevel},
	}
	for _, tc := range cases {
		got := parseLevel(tc.in)
		require.Equalf(t, tc.want, got, "parseLevel(%q)", tc.in)
	}
}

func Test_Init_DefaultFromConfig_JSONFormatter(t *testing.T) {
	reinit(t, "info", "")
	require.NotNil(t, L())
	require.Equal(t, logrus.InfoLevel, L().GetLevel())
	_, ok := L().Formatter.(*logrus.JSONFormatter)
	require.True(t, ok, "logger must use JSON formatter")
}

func Test_Init_EnvOverridesConfig(t *testing.T) {
	// В конфиге warn, но окружение задаёт debug, берём окружение
	reinit(t, "warn", "debug")
	require.Equal(t, logrus.DebugLevel, L().GetLevel())

	// Очищаем окружение, берём уже конфиг
	reinit(t, "error", "")
	require.Equal(t, logrus.ErrorLevel, L().GetLevel())
}

func Test_L_IsSafeWithoutExplicitInit(t *testing.T) {
	// имитируем ситуацию, где кто-то вызвал L() до Init()
	log = nil
	_ = os.Unsetenv("LOG_LEVEL")
	config.C().Log.Level = "info"

	// не должно паниковать и должен быть JSON форматтер
	require.NotPanics(t, func() { _ = L() })
	_, ok := L().Formatter.(*logrus.JSONFormatter)
	require.True(t, ok)
}
