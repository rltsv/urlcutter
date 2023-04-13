package entity

type Link struct {
	LinkID   string `json:"link_id"`
	UserID   string `json:"user_id"`
	LongURL  string `json:"long_link"`
	ShortURL string `json:"short_link"`
}

type CreateLinkDTO struct {
	UserID  string
	LongURL string
}

type GetLinkDTO struct {
	UserID string
	LinkID string
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

type InputData struct {
	URL string `json:"url"`
}

type OutputData struct {
	Response string `json:"result"`
}
