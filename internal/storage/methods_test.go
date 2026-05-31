package storage

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
)

// Вспомогательный метод инициализации хранилища для каждого теста
func setupTestStorage(t *testing.T) (*Storage, func()) {
	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test_bbolt.db")

	// Используем ваш конструктор, который сам создает и инициализирует бакеты
	storage, err := NewStorage(dbPath)
	if err != nil {
		slog.Error("failed to create storage", slog.Any("error", err))
		t.Fatalf("не удалось инициализировать хранилище: %v", err)
	}

	return storage, func() {
		storage.Close()
		os.Remove(dbPath)
	}
}

func TestStorage_AllMethodsFlow(t *testing.T) {
	storage, teardown := setupTestStorage(t)
	defer teardown()

	// --- 1. ТЕСТ: Добавление соединений (AddConn) ---
	conn1 := &UpsertConnectionDto{
		Label:      "  Prod-Server-1 ", // Нарочно добавляем пробелы и разный регистр
		Host:       "192.168.1.1",
		Port:       22,
		User:       "root",
		AuthMethod: PassphraseMethod,
		KeyPath:    "~/.ssh/id_ed25519",
	}

	if err := storage.AddConn(conn1); err != nil {
		t.Fatalf("ошибка при добавлении conn1: %v", err)
	}

	// Проверка на дубликаты (даже если регистр и пробелы отличаются)
	duplicateConn := &UpsertConnectionDto{
		Label: "prod-server-1",
		Host:  "10.0.0.1",
	}
	if err := storage.AddConn(duplicateConn); !errors.Is(err, ErrConnectionAlreadyExists) {
		t.Errorf("ожидалась ошибка дубликата ErrConnectionAlreadyExists, получено: %v", err)
	}

	// --- 2. ТЕСТ: Получение по имени (GetOneByName) и точное O(1) совпадение ---
	// Ищем нормализованную строку без пробелов и в нижнем регистре
	found, err := storage.GetOneByName("prod-server-1")
	if err != nil {
		t.Fatalf("ошибка получения по имени: %v", err)
	}
	if found.Label != "  Prod-Server-1 " {
		t.Errorf("ожидался оригинальный сохраненный лейбл, получено: '%s'", found.Label)
	}
	if found.ID == "" {
		t.Error("ID (UUID) нового соединения пустой")
	}

	// Проверка поиска несуществующего элемента
	_, err = storage.GetOneByName("unknown-server")
	if !errors.Is(err, ErrConnectionNotFound) {
		t.Errorf("ожидалась ошибка ErrConnectionNotFound, получено: %v", err)
	}

	// --- 3. ТЕСТ: Обновление без смены имени (Update) ---
	originalID := found.ID
	updateDto := &UpsertConnectionDto{
		Label:      "  Prod-Server-1 ", // Имя то же самое
		Host:       "192.168.1.99",     // Меняем хост
		Port:       2222,               // Меняем порт
		User:       "admin",
		AuthMethod: PasswordMethod,
	}

	if err := storage.Update("prod-server-1", updateDto); err != nil {
		t.Fatalf("ошибка обновления соединения: %v", err)
	}

	// Проверяем, что данные изменились, но UUID остался прежним
	updatedConn, err := storage.GetOneByName("prod-server-1")
	if err != nil {
		t.Fatalf("ошибка получения обновленного соединения: %v", err)
	}
	if updatedConn.Host != "192.168.1.99" || updatedConn.Port != 2222 {
		t.Error("данные внутри соединения не обновились")
	}
	if updatedConn.ID != originalID {
		t.Errorf("КРИТИЧЕСКИЙ БАГ: ID (UUID) изменился при обновлении! Было: %s, стало: %s", originalID, updatedConn.ID)
	}

	// --- 4. ТЕСТ: Обновление со сменой имени ---
	renameDto := &UpsertConnectionDto{
		Label: "New-Awesome-Server",
		Host:  "192.168.1.99",
		Port:  2222,
		User:  "admin",
	}

	if err := storage.Update("prod-server-1", renameDto); err != nil {
		t.Fatalf("ошибка при переименовании соединения: %v", err)
	}

	// Проверяем, что по старому имени больше ничего нет
	_, err = storage.GetOneByName("prod-server-1")
	if !errors.Is(err, ErrConnectionNotFound) {
		t.Error("старый ключ не удалился из BoltDB после переименования")
	}

	// Проверяем, что по новому имени данные доступны и ID сохранен
	renamedConn, err := storage.GetOneByName("new-awesome-server")
	if err != nil {
		t.Fatalf("не удалось найти соединение по новому имени: %v", err)
	}
	if renamedConn.ID != originalID {
		t.Error("ID (UUID) потерялся при переименовании")
	}

	// --- 5. ТЕСТ: Получение массового списка (GetConns) ---
	// Добавим еще одно соединение, чтобы в списке стало 2 элемента
	_ = storage.AddConn(&UpsertConnectionDto{Label: "Second-Server", Host: "8.8.8.8"})

	list, err := storage.GetConns()
	if err != nil {
		t.Fatalf("ошибка получения списка всех соединений: %v", err)
	}
	if len(list) != 2 {
		t.Errorf("ожидалось 2 соединения в базе, найдено: %d", len(list))
	}

	// --- 6. ТЕСТ: Удаление (DeleteConn) ---
	if err := storage.DeleteConn("second-server"); err != nil {
		t.Fatalf("ошибка при удалении соединения: %v", err)
	}

	// Проверяем, что список действительно уменьшился до 1
	listAfterDelete, _ := storage.GetConns()
	if len(listAfterDelete) != 1 {
		t.Errorf("ожидался 1 элемент после удаления, найдено: %d", len(listAfterDelete))
	}

	// Ошибка при попытке удалить уже несуществующее соединение
	if err := storage.DeleteConn("second-server"); !errors.Is(err, ErrConnectionNotFound) {
		t.Errorf("ожидалась ошибка ErrConnectionNotFound при повторном удалении, получено: %v", err)
	}
}
