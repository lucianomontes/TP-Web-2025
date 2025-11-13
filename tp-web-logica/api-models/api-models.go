package apimodels

type CreateGameReq struct {
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Categoria   string `json:"categoria"`
	Fecha       string `json:"fecha"`
	Estado      string `json:"estado"`
	Imagen      string `json:"imagen"`
}

type UpdateGameReq struct {
	ID          int32  `json:"id"`
	Titulo      string `json:"titulo"`
	Descripcion string `json:"descripcion"`
	Categoria   string `json:"categoria"`
	Fecha       string `json:"fecha"`
	Estado      string `json:"estado"`
	Imagen      string `json:"imagen"`
}
