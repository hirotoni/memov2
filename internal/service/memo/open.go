package memo

func (uc memo) Open(path string) error {
	path, err := resolveToMemosDir(uc.config.MemosDir(), path)
	if err != nil {
		return err
	}
	return uc.editor.Open(uc.config.BaseDir(), path)
}
