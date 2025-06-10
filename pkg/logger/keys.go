package logger

const (
	// Общее
	LogKeyOp    = "op"    // Название операции, например "repo.CreateUser"
	LogKeyError = "error" // Ошибка (всегда в конце)

	// Request context
	LogKeyRequestID = "request_id" // ID запроса (обычно из middleware)
	LogKeyTraceID   = "trace_id"   // ID трейса (если включён OpenTelemetry)
	LogKeySpanID    = "span_id"    // ID спана (опционально)

	// Пользователь
	LogKeyUserID    = "user_id"    // ID текущего пользователя
	LogKeySessionID = "session_id" // Сессия пользователя (опционально)

	// HTTP-запрос
	LogKeyStatus     = "status"      // HTTP-статус ответа (200, 404, ...)
	LogKeyMethod     = "method"      // HTTP-метод (GET, POST, ...)
	LogKeyPath       = "path"        // Человекочитаемый путь
	LogKeyBytes      = "bytes"       // Кол-во байт, отправленных клиенту
	LogKeyUserAgent  = "user_agent"  // User-Agent клиента
	LogKeyRemoteAddr = "remote_addr" // IP клиента

	// Тайминги
	LogKeyDurationMS = "duration_ms" // Продолжительность операции (в миллисекундах)
)
