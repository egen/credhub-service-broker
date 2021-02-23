package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	"code.cloudfoundry.org/lager"
	brokerapi "github.com/pivotal-cf/brokerapi/domain"
	"github.com/starkandwayne/credhub-service-broker/config"
)

type ConfigServerBroker struct {
	Config config.Config
	Logger lager.Logger
}

func (broker *ConfigServerBroker) Services(ctx context.Context) ([]brokerapi.Service, error) {
	planList := []brokerapi.ServicePlan{
		brokerapi.ServicePlan{
			ID:          broker.Config.BasicPlanId,
			Name:        broker.Config.BasicPlanName,
			Description: broker.Config.Description,
			Metadata: &brokerapi.ServicePlanMetadata{
				DisplayName: "Basic",
			},
		}}

	return []brokerapi.Service{
		brokerapi.Service{
			ID:          broker.Config.ServiceID,
			Name:        broker.Config.ServiceName,
			Description: broker.Config.Description,
			Bindable:    true,
			Plans:       planList,
			Metadata: &brokerapi.ServiceMetadata{
				DisplayName:         broker.Config.DisplayName,
				LongDescription:     broker.Config.LongDescription,
				DocumentationUrl:    broker.Config.DocumentationURL,
				SupportUrl:          broker.Config.SupportURL,
				ImageUrl:            fmt.Sprintf("data:image/png;base64,%s", broker.Config.IconImage),
				ProviderDisplayName: broker.Config.ProviderDisplayName,
			},
			Tags: []string{
				"snw",
				"credhub",
			},
		},
	}, nil
}

func (broker *ConfigServerBroker) Provision(ctx context.Context, instanceID string, serviceDetails brokerapi.ProvisionDetails, asyncAllowed bool) (spec brokerapi.ProvisionedServiceSpec, err error) {
	spec = brokerapi.ProvisionedServiceSpec{}

	value := values.JSON{}
	err = json.Unmarshal(serviceDetails.RawParameters, &value)

	if err != nil {
		return spec, err
	}
	broker.Logger.Info("provision")

	chcli, err := NewCredHub(
		broker.Config,
		broker.Logger,
	)
	if err != nil {
		return spec, err
	}

	err = chcli.WriteSecret(instanceID, value)

	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (broker *ConfigServerBroker) Deprovision(ctx context.Context, instanceID string, details brokerapi.DeprovisionDetails, asyncAllowed bool) (brokerapi.DeprovisionServiceSpec, error) {
	spec := brokerapi.DeprovisionServiceSpec{}

	broker.Logger.Info("deprovision")

	chcli, err := NewCredHub(
		broker.Config,
		broker.Logger,
	)
	if err != nil {
		return spec, err
	}

	err = chcli.DeleteSecret(instanceID)

	if err != nil {
		return spec, err
	}

	return spec, nil
}

func (broker *ConfigServerBroker) Unbind(ctx context.Context, instanceID, bindingID string, details brokerapi.UnbindDetails, asyncAllowed bool) (brokerapi.UnbindSpec, error) {
	unbind := brokerapi.UnbindSpec{}

	return unbind, nil
}

func (broker *ConfigServerBroker) Bind(ctx context.Context, instanceID, bindingID string, details brokerapi.BindDetails, asyncAllowed bool) (brokerapi.Binding, error) {
	binding := brokerapi.Binding{}
	chcli, err := NewCredHub(
		broker.Config,
		broker.Logger,
	)
	if err != nil {
		return binding, err
	}
	binding = brokerapi.Binding{
		Credentials: map[string]string{
			"credhub-ref": chcli.MakePath(instanceID),
		},
	}

	return binding, nil
}

// LastOperation ...
// If the broker provisions asynchronously, the Cloud Controller will poll this endpoint
// for the status of the provisioning operation.
func (broker *ConfigServerBroker) LastOperation(ctx context.Context, instanceID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	return brokerapi.LastOperation{}, errors.New("not implemented")
}

func (broker *ConfigServerBroker) Update(cxt context.Context, instanceID string, details brokerapi.UpdateDetails, asyncAllowed bool) (brokerapi.UpdateServiceSpec, error) {
	return brokerapi.UpdateServiceSpec{}, errors.New("not implemented")
}

func (broker *ConfigServerBroker) GetBinding(ctx context.Context, instanceID, bindingID string) (brokerapi.GetBindingSpec, error) {
	return brokerapi.GetBindingSpec{}, errors.New("not implemented")
}

func (broker *ConfigServerBroker) GetInstance(ctx context.Context, instanceID string) (brokerapi.GetInstanceDetailsSpec, error) {
	return brokerapi.GetInstanceDetailsSpec{}, errors.New("not implemented")
}

func (broker *ConfigServerBroker) LastBindingOperation(ctx context.Context, instanceID, bindingID string, details brokerapi.PollDetails) (brokerapi.LastOperation, error) {
	//create client

	return brokerapi.LastOperation{}, errors.New("not implemented")
}
