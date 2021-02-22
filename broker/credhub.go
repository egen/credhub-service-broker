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
}

func NewCredHub(config config.Config, logger lager.Logger) (CredHub, error) {

	ch, err := chcli.New(
		config.Credhub.Server,
		chcli.SkipTLSValidation(true), // TODO use CA
		chcli.Auth(auth.UaaClientCredentials(config.Credhub.Client, config.Credhub.Secret)),
	)

	if err != nil {
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
func (ch *credhub) WriteSecret(instanceId string, value values.JSON) error {
	ch.logger.Info("Writing secret")
	_, err := ch.client.SetJSON(ch.MakePath(instanceId), value)

	return err
}

func (ch *credhub) DeleteSecret(instanceId string) error {
	return ch.client.Delete(ch.MakePath(instanceId))
}
