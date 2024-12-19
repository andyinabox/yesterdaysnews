package textprocessor

func (p *Processor) ExportModel() ([]byte, error) {
	return p.chain.Save()
}
