package container

type TestContainer struct {
	Container     Container
	InitFunc      func(Container, string) error
	SelectAllFunc func(Container, string) (map[string][]byte, error)
	InsertFunc    func(Container, string, string, []byte) error
	DeleteFunc    func(Container, string, string) (bool, error)
}

func NewTestContainer(container Container) *TestContainer {
	return &TestContainer{
		Container: container,
		InitFunc:  func(container Container, tableID string) error { 
			return container.Init(tableID) 
		},
		SelectAllFunc: func(container Container, tableID string) (map[string][]byte, error) {
			return container.SelectAll(tableID)
		},
		InsertFunc: func(container Container, tableID, key string, data []byte) error {
			return container.Insert(tableID, key, data)
		},
		DeleteFunc:  func(container Container, tableID, key string) (bool, error) {
			return container.Delete(tableID, key)
		},
	}
}

func (this *TestContainer) Init(tableID string) error {
	return this.InitFunc(this.Container, tableID)
}

func (this *TestContainer) SelectAll(tableID string) (map[string][]byte, error) {
	return this.SelectAllFunc(this.Container, tableID)
}

func (this *TestContainer) Insert(tableID, key string, entry []byte) error {
	return this.InsertFunc(this.Container, tableID, key, entry)
}

func (this *TestContainer) Delete(tableID, key string) (bool, error) {
	return this.DeleteFunc(this.Container, tableID, key)
}
