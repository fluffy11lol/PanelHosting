package models

type Metadata struct {
	ID        string `json:"id"`         // Уникальный ID файла/директории
	UserID    string `json:"user_id"`    // ID пользователя
	Name      string `json:"name"`       // Имя файла/директории
	Path      string `json:"path"`       // Путь в структуре
	MimeType  string `json:"mime_type"`  // Тип содержимого (файл или директория)
	URL       string `json:"url"`        // Ссылка на файл/директорию
	CreatedAt string `json:"created_at"` // Время создания
	UpdatedAt string `json:"updated_at"` // Время последнего обновления
}
