package server

import (
	"fmt"
	"html/template"
)

type RenderContext struct {
	AssetsPath string
	ImportMap  template.HTML
}

const importMapTpl = `
<script type="importmap">
%s
</script>
`

func (s *Server) renderContext() RenderContext {

	// we only need the import map when serving assets from the local FS
	importMap := template.HTML(fmt.Sprintf(importMapTpl, s.cfg.ImportMap))
	if !s.cfg.UseFilesystemAssets {
		importMap = template.HTML("")
	}

	return RenderContext{
		AssetsPath: s.assetsPath,
		ImportMap:  importMap,
	}
}
