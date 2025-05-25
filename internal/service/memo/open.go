package memo

func (uc memo) Open(path string) error {
	path = resolveToMemosDir(uc.config.MemosDir(), path)
	return uc.editor.Open(uc.config.BaseDir(), path)
}
