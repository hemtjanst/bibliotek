package feature

type Fake struct {
	Err         error
	FeatureName string
}

func (f *Fake) Name() string {
	return f.FeatureName
}
func (f *Fake) Min() int {
	return 0
}
func (f *Fake) Max() int {
	return 1
}
func (f *Fake) Step() int {
	return 1
}
func (f *Fake) Set(string) error {
	return f.Err
}
func (f *Fake) Update(string) error {
	return f.Err
}
func (f *Fake) OnSet() (chan string, error) {
	return nil, f.Err
}
func (f *Fake) OnUpdate() (chan string, error) {
	return nil, f.Err
}
func (f *Fake) Exists() bool {
	return false
}
