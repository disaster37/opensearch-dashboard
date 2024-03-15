package api

import (
	"github.com/go-resty/resty/v2"
)

// Api is the main API interface
type Api interface {
	SavedObject() SavedObjectApi
	ShortenUrl() ShortenUrlApi
	Status() StatusApi
	Client() *resty.Client
}

// SavedObjectApi is the interface to call with saved object API
type SavedObjectApi interface {
	Get(tenantId string, objectType string, objectId string) (object *Object, err error)
	BulkGet(tenantId string, bulks []SavedObjectBulkGetOption) (objects []Object, err error)
	Find(tenantId string, option SavedObjectFindOption) (objects []Object, err error)
	Create(tenantId string, object *Object, overwrite bool) (createdObject *Object, err error)
	BulkCreate(tenantId string, overwrite bool, bulks []Object) (createdObjects []Object, err error)
	Update(tenantId string, object *Object) (updatedObject *Object, err error)
	Delete(tenantId string, objectType string, objectId string, force bool) (err error)
	Export(tenantId string, option SavedObjectExportOption) (data []byte, err error)
	Import(tenantId string, option SavedObjectImportOption, data []byte) (res map[string]any, err error)
}

// ShortenUrlApi is the interface to call with shorten API
type ShortenUrlApi interface {
	Create(url string) (shortUrl string, err error)
}

type StatusApi interface {
	Status() (status map[string]any, err error)
}

// DefaultApi is the default implementation of Api interface
type DefaultApi struct {
	client         *resty.Client
	shortenUrlApi  ShortenUrlApi
	savedObjectApi SavedObjectApi
	statusApi      StatusApi
}

// New permit to get the default implementation of Api interface
func New(client *resty.Client) Api {
	return &DefaultApi{
		client:         client,
		shortenUrlApi:  NewShortenUrlApi(client),
		savedObjectApi: NewSavedObjectApi(client),
		statusApi:      NewStatusApi(client),
	}
}

func (h DefaultApi) Client() *resty.Client {
	return h.client
}

func (h DefaultApi) ShortenUrl() ShortenUrlApi {
	return h.shortenUrlApi
}

func (h DefaultApi) SavedObject() SavedObjectApi {
	return h.savedObjectApi
}

func (h DefaultApi) Status() StatusApi {
	return h.statusApi
}
