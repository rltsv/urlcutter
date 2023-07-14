package entity

type CreateLinkDTO struct {
	UserID        string `json:"user_id"`
	OriginalURL   string `json:"original_url,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

type GetLinkDTO struct {
	UserID string
	LinkID string
}

type GetAllLinksDTO struct {
	UserID string
}

type SendLinkDTO struct {
	ShortURL      string `json:"short_url,omitempty"`
	OriginalURL   string `json:"original_url,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
}
