package logger

// Logger описывает универсальный интерфейс для логирования.
// Позволяет логировать сообщения различных уровней и создавать "scoped" логгеры через With(...).
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any) // Завершает приложение с кодом 1
	With(args ...any) Logger       // Возвращает логгер с дополнительным контекстом
}
