package discovery

import (
	"fmt"
	"os"

	"github.com/ArthurHlt/go-eureka-client/eureka"
)

func QueryAll() {
	apps, err := euClient.GetApplications()
	if nil != err {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, app := range apps.Applications {
		fmt.Println(app)
	}
}

func QueryApp(appId string) ([]eureka.InstanceInfo, error) {
	app, err := euClient.GetApplication(appId)
	if nil != err {
		return nil, err
	}

	return app.Instances, nil
}