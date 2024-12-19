package lorca

func NewPage(source string, op ...Option) Page {
	return &page{
		Source:  dealSource(source),
		Options: op,
	}
}

type Page interface {
	Switch(app APP) error
}

type page struct {
	Source  string
	Options []Option
}

func (this *page) Switch(app APP) error {
	if err := app.Load(this.Source); err != nil {
		return err
	}
	for _, v := range this.Options {
		if err := v(app); err != nil {
			return err
		}
	}
	return nil
}
