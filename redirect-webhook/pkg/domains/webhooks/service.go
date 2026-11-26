package webhooks

import (
	"context"
	"encoding/json"
	"net/url"

	"conductor.com.br/pier/hooks/pkg/domains/encryption"
	"conductor.com.br/pier/hooks/pkg/domains/event_controller"
	"conductor.com.br/pier/hooks/pkg/domains/issuer_events"
	"conductor.com.br/pier/hooks/pkg/domains/tenants"
	"conductor.com.br/pier/hooks/pkg/helpers"
	"conductor.com.br/pier/hooks/pkg/helpers/constants"
	"github.com/jinzhu/gorm"
)

type Service interface {
	Create(key *tenants.TenantKey, hook WebHook) WebHook
	Update(key *tenants.TenantKey, id string, updateMap map[string]interface{}) WebHook
	FindByIdWithPreload(key *tenants.TenantKey, id string) WebHook
	FetchWithPreload(key *tenants.TenantKey, id string) (WebHook, error)
	FindByTypeWithPreload(key *tenants.TenantKey, eventType string) (WebHook, error)
	FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string) helpers.Paginator
	Trigger(ctx context.Context, webhook WebHook, payload string, forceMockPayload bool, eventType string, publicKey string) (statusCode int, responseBody string, encryptedAesKey string)
	Delete(tenant *tenants.TenantKey, id string)
	FindAllByTypeAndIssuerId(context context.Context, tenantKey *tenants.TenantKey, eventType string, issuerId int) ([]WebHook, error)
}

type ServiceImpl struct {
	repository          Repository
	repositoryHeader    RepositoryHeader
	encryptionService   encryption.Service
	issuerEventsService issuer_events.Service
	eventController     event_controller.Controller
}

func NewService(repository Repository, encryptionService encryption.Service, issuerEventsService issuer_events.Service, eventController event_controller.Controller, repositoryHeader RepositoryHeader) ServiceImpl {
	return ServiceImpl{repository: repository, encryptionService: encryptionService, issuerEventsService: issuerEventsService, eventController: eventController, repositoryHeader: repositoryHeader}
}

func (s ServiceImpl) Create(key *tenants.TenantKey, webHook WebHook) WebHook {

	issuerEvent := s.issuerEventsService.FindByEventNameAndIssuerId(key, webHook.EventType, key.IssuerId)
	webHook.EventType = issuerEvent.Name
	parsedUrl, err := url.Parse(webHook.URL)
	if err != nil || parsedUrl == nil {
		Error{}.BadRequest("invalid URL").WithChild(err).Throw()
	}

	if parsedUrl.Host == "" {
		Error{}.BadRequest("invalid URL").WithChild(err).Throw()
	}

	existingEventWithSameUrl, err := s.FindByTypeAndUrl(context.Background(), key, webHook.EventType, webHook.URL)
	if err != nil && err != gorm.ErrRecordNotFound {
		Error{}.HandleDatabaseError(err)
	}
	if existingEventWithSameUrl.ID != "" {
		Error{}.BadRequest("another webhook for this event with this URL already exists").Throw()
	}

	savedWebHook, e := s.repository.Save(key, webHook)
	Error{}.HandleDatabaseError(e)
	_ = s.eventController.UpdateEventByIssuerAndName(key, key.IssuerId, webHook.EventType)
	return savedWebHook
}

func (s ServiceImpl) Update(key *tenants.TenantKey, id string, updateMap map[string]interface{}) WebHook {
	newEventType := updateMap["event_type"].(string)
	s.issuerEventsService.FindByEventNameAndIssuerId(key, newEventType, key.IssuerId)
	s.FindByIdWithPreload(key, id)

	if newUrl, inMap := updateMap["url"].(string); inMap {
		existingEventWithSameUrl, err := s.FindByTypeAndUrl(context.Background(), key, newEventType, newUrl)
		if err != nil && err != gorm.ErrRecordNotFound {
			Error{}.HandleDatabaseError(err)
		}
		if existingEventWithSameUrl.ID != "" {
			Error{}.BadRequest("another webhook for this event with this URL already exists").Throw()
		}
	}

	e := s.repository.Update(key, id, updateMap)
	Error{}.HandleDatabaseErrorWithMessage(e, "error updating webhook")

	if hdrs, ok := updateMap["Headers"]; ok {
		var headers []Header
		bHeader, _ := json.Marshal(hdrs)
		json.Unmarshal(bHeader, &headers)

		err := s.repositoryHeader.DeleteByWebhookId(key, id)
		Error{}.HandleDatabaseErrorWithMessage(err, "error deleting headers")

		for i := 0; i < len(headers); i++ {
			headers[i].WebHookId = id
			_, err = s.repositoryHeader.Save(key, headers[i])
			Error{}.HandleDatabaseErrorWithMessage(err, "error save header")
		}

	}

	updatedWebhook, e := s.FetchWithPreload(key, id)
	Error{}.HandleDatabaseErrorWithMessage(e, "error finding webhook after update")
	_ = s.eventController.UpdateEventByIssuerAndName(key, key.IssuerId, newEventType)
	return updatedWebhook
}

func (s ServiceImpl) FindByIdWithPreload(key *tenants.TenantKey, id string) WebHook {
	webHook, e := s.FetchWithPreload(key, id)
	Error{}.HandleSearch(e)
	return webHook
}

func (s ServiceImpl) FetchWithPreload(key *tenants.TenantKey, id string) (WebHook, error) {
	return s.repository.FindByIdWithPreload(key, id)
}

func (s ServiceImpl) FindByTypeWithPreload(key *tenants.TenantKey, eventType string) (WebHook, error) {
	return s.repository.FindByTypeWithPreload(key, eventType)
}

func (s ServiceImpl) FindAllWithPagination(key *tenants.TenantKey, page int, limit int, sort []string) helpers.Paginator {
	paginator, e := s.repository.FindAllWithPagination(key, page, limit, sort, nil)
	Error{}.HandlePagination(paginator, e)
	return paginator
}

func (s ServiceImpl) Trigger(ctx context.Context, webhook WebHook, payload string, forceMockPayload bool, eventType string, publicKey string) (statusCode int, responseBody string, encryptedAesKey string) {
	newPayload := payload

	if forceMockPayload {
		if eventType == "" {
			eventType = "blank"
		}
		newPayload = constants.MockOfEventType(eventType)
	}

	if publicKey != "" {
		newPayload, encryptedAesKey = s.createKeyAndEncryptPayload(newPayload, publicKey)
		webhook = webhook.WithHeader("encryptedAesKey", encryptedAesKey)
	}

	httpResponse, err := webhook.Call(ctx, newPayload)
	Error{}.HandleHttpRequest(err)
	return httpResponse.Status, httpResponse.Body, encryptedAesKey
}

func (s ServiceImpl) createKeyAndEncryptPayload(dataToEncrypt string, publicKey string) (encryptedPayload string, encryptedAesKey string) {
	var err error
	keyVault := s.encryptionService.CreateMockKeyVault(publicKey)

	encryptedPayload, err = keyVault.EncryptWithAes(dataToEncrypt)
	Error{}.HandlePayloadEncryption(err)

	encryptedAesKey, err = keyVault.EncryptAesWithPublicKeyMock()
	Error{}.HandleAesEncryption(err)

	return encryptedPayload, encryptedAesKey
}
func (s ServiceImpl) Delete(tenant *tenants.TenantKey, id string) {
	webHook := s.FindByIdWithPreload(tenant, id)
	err := s.repository.Delete(tenant, id, webHook.EventType)
	Error{}.HandleDatabaseErrorWithMessage(err, "error deleting webhook")
	_ = s.eventController.UpdateEventByIssuerAndName(tenant, tenant.IssuerId, webHook.EventType)
}

func (s ServiceImpl) FindAllByTypeAndIssuerId(context context.Context, tenantKey *tenants.TenantKey, eventType string, issuerId int) ([]WebHook, error) {
	return s.repository.FindAllByTypeAndIssuerId(context, tenantKey, eventType, issuerId)
}

func (s ServiceImpl) FindByTypeAndUrl(context context.Context, tenantKey *tenants.TenantKey, eventType string, url string) (WebHook, error) {
	return s.repository.FindByTypeAndUrl(context, tenantKey, eventType, url)
}
