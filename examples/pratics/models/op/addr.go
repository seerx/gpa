package op

type Addr struct {
	Name string
}

func (a *Addr) Read(data []byte) error {
	a.Name = string(data)
	return nil
}
func (a *Addr) Write() ([]byte, error) {
	data := []byte(a.Name)
	return data, nil
}
