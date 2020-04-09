package discovery

var DoNothingClient = new(EmptyClient)

type EmptyClient struct{

}

func (e EmptyClient) QueryServices() ([]*InstanceInfo, error) {
	return nil, nil
}

func (e EmptyClient) Register() error {
	return nil
}

func (e EmptyClient) UnRegister() error {
	return nil
}

func (e EmptyClient) Get(string) []*InstanceInfo {
	return nil
}

func (e EmptyClient) StartPeriodicalRefresh() error {
	return nil
}

func (e EmptyClient) GetInternalRegistryStore() *InsInfoArrSyncMap {
	return nil
}

func (e EmptyClient) SetInternalRegistryStore(*InsInfoArrSyncMap) {
}



