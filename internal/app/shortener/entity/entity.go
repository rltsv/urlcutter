package entity

type Link struct {
	LinkID   string `json:"link_id"`
	UserID   string `json:"user_id"`
	LongURL  string `json:"long_url,omitempty"`
	ShortURL string `json:"short_url,omitempty"`
}

type CreateLinkDTO struct {
	UserID  string `json:"user_id"`
	LongURL string `json:"long_url"`
}

type GetLinkDTO struct {
	UserID string
	LinkID string
}

type GetAllLinksDTO struct {
	UserID string
}

type SendLinkDTO struct {
	ShortURL string `json:"short_url,omitempty"`
	LongURL  string `json:"long_url,omitempty"`
}

func NewLink(dto CreateLinkDTO) Link {
	return Link{
		UserID:  dto.UserID,
		LongURL: dto.LongURL,
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
