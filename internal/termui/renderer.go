package termui

type Renderer struct {
	cfg Config
	screen [][]byte
}

func NewRenderer(cfg Config, size int) *Renderer {
	return &Renderer{
		cfg: cfg,
		screen: make([][]byte, size),
	}
}

func (r *Renderer) Render(){

}
