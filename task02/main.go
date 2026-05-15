// Задание 2: Таймаут для медленной операции
//
// context.WithTimeout - это автоматическая отмена через заданное время.
// Очень часто используется для HTTP-запросов, запросов в БД и т.д.
//
// Разница от WithCancel:
//   - WithCancel: ты вручную вызываешь cancel()
//   - WithTimeout: контекст отменится сам через N времени (и тоже можно cancel())
//
// ВАЖНО: всегда вызывай defer cancel() даже с таймаутом!
//   Если операция завершится ДО таймаута - cancel() освободит ресурсы немедленно.
//   Без defer cancel() ресурсы освободятся только когда сработает таймаут.
//
// ──────────────────────────────────────────────────────────────────────────────
//
// Напиши функцию downloadFile(ctx context.Context, url string) (string, error),
// которая имитирует скачивание файла:
//   - использует select с двумя случаями:
//       * time.After(2 * time.Second): "скачивание" заняло 2 секунды, возвращает
//         ("содержимое " + url, nil)
//       * ctx.Done(): возвращает ("", ctx.Err())
//
// В main() вызови downloadFile дважды:
//
//   Вызов 1: таймаут 3 секунды → должен успеть
//     ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
//     defer cancel()
//     result, err := downloadFile(ctx, "https://example.com/file.txt")
//     // ожидаемый вывод: скачано: содержимое https://example.com/file.txt
//
//   Вызов 2: таймаут 1 секунда → не успеет, таймаут
//     // ожидаемый вывод: ошибка: context deadline exceeded
//
// После каждого вызова выводи время выполнения через time.Since(start).
//
// Ожидаемый вывод:
//   скачано: содержимое https://example.com/file.txt (за ~2s)
//   ошибка: context deadline exceeded (за ~1s)
//
// Запусти: go run main.go

package main

import (
	"context"
	"fmt"
	"time"
)

func timeTrack() func() {
	start := time.Now()
	return func() {
		fmt.Println("за", time.Since(start))
	}
}

// TODO: напиши функцию downloadFile(ctx context.Context, url string) (string, error)
func downloadFile(ctx context.Context, url string) (string, error) {
	start := time.Now()
	defer func() {
		fmt.Printf("за %v ", time.Since(start))
	}()
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(2 * time.Second):
		return fmt.Sprint("содержимое ", url), nil
	}
}

func main() {
	// Вызов 1: таймаут 3 секунды - должен успеть
	// TODO: создай контекст с таймаутом 3s, вызови downloadFile, выведи результат
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := downloadFile(ctx, "https://example.com/file.txt")
	if err != nil {
		fmt.Printf("ошибка: %s\n", err)
	} else {
		fmt.Printf("скачано: %s\n", res)
	}

	// Вызов 2: таймаут 1 секунда - не успеет
	// TODO: создай контекст с таймаутом 1s, вызови downloadFile, выведи ошибку
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel2()

	res, err = downloadFile(ctx2, "https://example.com/file.txt")
	if err != nil {
		fmt.Println("ошибка:", err)
	} else {
		fmt.Println("скачано:", res)
	}
}
