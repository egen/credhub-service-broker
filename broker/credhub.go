package broker

import (
	"fmt"

	chcli "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
	"code.cloudfoundry.org/lager"
	"github.com/starkandwayne/credhub-service-broker/config"

	// "code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
)

type CredHub interface {
	WriteSecret(string, values.JSON) error
	DeleteSecret(string) error
	MakePath(string) string
	AddReadPermission(instanceId, appGuid string) error
	DeletePermission(instanceId, appGuid string) error
}

func NewCredHub(config config.Config, logger lager.Logger) (CredHub, error) {

	ch, err := chcli.New(
		config.Credhub.Server,
		chcli.SkipTLSValidation(true), // TODO use CA
		chcli.Auth(auth.UaaClientCredentials(config.Credhub.Client, config.Credhub.Secret)),
	)

	if err != nil {
		logger.Error("could-not-create-credhub", err)
		return nil, err
	}
	return &credhub{client: ch, logger: logger, config: config}, nil
}

type credhub struct {
	config config.Config
	client *chcli.CredHub
	logger lager.Logger
}

func (ch *credhub) MakePath(instanceId string) string {
	return fmt.Sprintf("/%v/%v/%v/credentials", ch.config.ServiceName, ch.config.ServiceID, instanceId)
}

func MakeActor(appGuid string) string {
	return fmt.Sprintf("mtls-app:%s", appGuid)
}

func (ch *credhub) WriteSecret(instanceId string, value values.JSON) error {
	ch.logger.Info("Writing secret")
	_, err := ch.client.SetJSON(ch.MakePath(instanceId), value)

	return err
}

func (ch *credhub) DeleteSecret(instanceId string) error {
	return ch.client.Delete(ch.MakePath(instanceId))
}

func (ch *credhub) AddReadPermission(instanceId, appGuid string) error {
	_, err := ch.client.AddPermission(ch.MakePath(instanceId),
		MakeActor(appGuid),
		[]string{"read"},
	)
	return err
}

func (ch *credhub) DeletePermission(instanceId, appGuid string) error {
	permission, err := ch.client.GetPermissionByPathActor(
		ch.MakePath(instanceId), MakeActor(appGuid))
	if err != nil {
		return fmt.Errorf("failed to lookup permissions for instance: %s got: %s",
			instanceId, err)
	}

	_, err = ch.client.DeletePermission(permission.UUID)
	return err
}
