package discovery

import (
	"github.com/wanghongfei/go-eureka-client/eureka"
)

func QueryAll() ([]eureka.Application, error) {
	apps, err := euClient.GetApplications()
	if nil != err {
		return nil, err
	}

	return apps.Applications, nil
}

func QueryApp(appId string) ([]eureka.InstanceInfo, error) {
	app, err := euClient.GetApplication(appId)
	if nil != err {
		return nil, err
	}

	return app.Instances, nil
}