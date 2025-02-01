package server

type RenderContext struct {
	AssetsPath string
}

func (s *Server) renderContext() RenderContext {
	return RenderContext{
		AssetsPath: s.assetsPath,
	}
}
