package entity

type Link struct {
	LinkID        string `json:"link_id"`
	UserID        string `json:"user_id"`
	OriginalURL   string `json:"original_url,omitempty"`
	ShortURL      string `json:"short_url,omitempty"`
	CorrelationID string `json:"correlation_id,omitempty"`
}

func NewLink(dto CreateLinkDTO) Link {
	return Link{
		UserID:        dto.UserID,
		OriginalURL:   dto.OriginalURL,
		CorrelationID: dto.CorrelationID,
	}
}

func GetLink(dto GetLinkDTO) Link {
	return Link{
		LinkID: dto.LinkID,
		UserID: dto.UserID,
	}
}

func GetAllLinks(dto GetAllLinksDTO) Link {
	return Link{
		UserID: dto.UserID,
	}
}