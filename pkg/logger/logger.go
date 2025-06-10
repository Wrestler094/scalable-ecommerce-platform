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

	// DSL-style методы
	WithOp(op string) Logger
	WithRequestID(id string) Logger
	WithUserID(id int64) Logger
	WithError(err error) Logger
}
