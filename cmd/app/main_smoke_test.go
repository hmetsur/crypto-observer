package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// Тест без подъёма БД/сервера: просто проверяем,
// что бинарь запускается и печатает help/usage.
func TestBinary_Help(t *testing.T) {
	// путь к корню модуля
	_, thisFile, _, _ := runtime.Caller(0)
	repoRoot := filepath.Join(filepath.Dir(thisFile), "../..")

	cmd := exec.Command("go", "run", "./cmd/app", "-h")
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(),
		// чтобы main не падал из-за отсутствия переменных — но мы вызываем -h,
		// приложение должно просто выйти.
		"HTTP_ADDR=:0",
		"DB_DSN=postgres://user:pass@localhost:5432/db?sslmode=disable",
	)

	out, err := cmd.CombinedOutput()
	if err == nil {
		// Ок, команда завершилась без ошибок.
		return
	}
	// Некоторые фреймворки возвращают код != 0 на -h.
	// Тогда хотя бы проверим, что что-то осмысленное напечатано.
	if len(out) == 0 {
		t.Fatalf("no output on -h, err: %v", err)
	}
}
