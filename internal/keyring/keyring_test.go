package keyring

import (
	"errors"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestKeyring_Flow(t *testing.T) {
	keyring.MockInit()

	k := NewKeyring("svlt-test-app")
	label := "prod-database"
	secret := "my-highly-secure-password-123"

	// --- 1. ТЕСТ: Получение несуществующего секрета ---
	_, err := k.Get(label)
	if !errors.Is(err, ErrSecreteNotFound) {
		t.Errorf("ожидалась ошибка %v, получено: %v", ErrSecreteNotFound, err)
	}

	// --- 2. ТЕСТ: Валидация пустых аргументов (ErrKeyringOperation) ---
	if err := k.Set("", secret); !errors.Is(err, ErrKeyringOperation) {
		t.Errorf("Set с пустой меткой: ожидалась ошибка %v, получено: %v", ErrKeyringOperation, err)
	}
	if err := k.Set(label, ""); !errors.Is(err, ErrKeyringOperation) {
		t.Errorf("Set с пустым секретом: ожидалась ошибка %v, получено: %v", ErrKeyringOperation, err)
	}
	if _, err := k.Get(""); !errors.Is(err, ErrKeyringOperation) {
		t.Errorf("Get с пустой меткой: ожидалась ошибка %v, получено: %v", ErrKeyringOperation, err)
	}
	if err := k.Delete(""); !errors.Is(err, ErrKeyringOperation) {
		t.Errorf("Delete с пустой меткой: ожидалась ошибка %v, получено: %v", ErrKeyringOperation, err)
	}

	// --- 3. ТЕСТ: Успешная запись и чтение секрета ---
	if err := k.Set(label, secret); err != nil {
		t.Fatalf("ошибка при записи валидного секрета: %v", err)
	}

	foundSecret, err := k.Get(label)
	if err != nil {
		t.Fatalf("ошибка при чтении существующего секрета: %v", err)
	}
	if foundSecret != secret {
		t.Errorf("ожидался секрет '%s', получен '%s'", secret, foundSecret)
	}

	// --- 4. ТЕСТ: Перезапись (обновление) существующего ключа ---
	newSecret := "updated-secure-password-456"
	if err := k.Set(label, newSecret); err != nil {
		t.Fatalf("ошибка при обновлении существующего секрета: %v", err)
	}

	foundNewSecret, _ := k.Get(label)
	if foundNewSecret != newSecret {
		t.Errorf("секрет не обновился, ожидалось '%s', получено '%s'", newSecret, foundNewSecret)
	}

	// --- 5. ТЕСТ: Успешное удаление ---
	if err := k.Delete(label); err != nil {
		t.Fatalf("ошибка при удалении существующего секрета: %v", err)
	}

	// Проверяем, что после удаления секрет снова возвращает ErrSecreteNotFound
	_, err = k.Get(label)
	if !errors.Is(err, ErrSecreteNotFound) {
		t.Errorf("после удаления ожидалась ошибка %v, получено: %v", ErrSecreteNotFound, err)
	}

	// --- 6. ТЕСТ: Повторное удаление несуществующего секрета ---
	if err := k.Delete(label); err != nil {
		t.Errorf("повторное удаление не должно возвращать ошибку, получено: %v", err)
	}
}
