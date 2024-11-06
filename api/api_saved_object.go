package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

const (
	basePathSavedObjectURL = "/api/saved_objects"
)

// Object is the object data
type Object struct {
	Type                 string           `json:"type,omitempty"`
	Id                   string           `json:"id,omitempty"`
	Attributes           map[string]any   `json:"attributes"`
	References           []map[string]any `json:"references,omitempty"`
	Version              string           `json:"version,omitempty"`
	CoreMigrationVersion string           `json:"coreMigrationVersion,omitempty"`
	CreatedAt            *time.Time       `json:"created_at,omitempty"`
	MigrationVersion     map[string]any   `json:"migrationVersion,omitempty"`
	UpdatedAt            *time.Time       `json:"updatedAt,omitempty"`
	OriginId             string           `json:"originId,omitempty"`
}

// SavedObjectBulkGetOption is the option when call bulk get API
type SavedObjectBulkGetOption struct {
	Type   string   `json:"type"`
	Id     string   `json:"id"`
	Fields []string `json:"fields,omitempty"`
}

// SavedObjectFindOption is the option when call find API
type SavedObjectFindOption struct {
	Type                  string   `json:"type"`
	PerPage               int      `json:"per_page,omitempty"`
	Page                  int      `json:"page,omitempty"`
	Search                string   `json:"search,omitempty"`
	DefaultSearchOperator string   `json:"default_search_operator,omitempty"`
	SearchFields          []string `json:"search_fields,omitempty"`
	Fields                []string `json:"fields,omitempty"`
	SortField             string   `json:"sort_field,omitempty"`
	HasReference          string   `json:"has_reference,omitempty"`
	Filter                string   `json:"filter,omitempty"`
}

// SavedObjectExportOption the option when call export API
type SavedObjectExportOption struct {
	Type                  string                     `json:"type,omitempty"`
	Objects               []SavedObjectBulkGetOption `json:"objects,omitempty"`
	IncludeReferencesDeep *bool                      `json:"includeReferencesDeep,omitempty"`
	ExcludeExportDetails  *bool                      `json:"excludeExportDetails,omitempty"`
}

// SavedObjectImportOption is the option when call import API
type SavedObjectImportOption struct {
	CreateNewCopies bool `json:"createNewCopies,omitempty"`
	Overwrite       bool `json:"overwrite,omitempty"`
}

// SavedObjectFindResponse is the response when cal find API
type SavedObjectFindResponse struct {
	Total        int64    `json:"total"`
	SavedObjects []Object `json:"saved_objects"`
}

// SavedObjectBulkCreateResponse is the repsonse when call bulk create API
type SavedObjectBulkResponse struct {
	SavedObjects []Object `json:"saved_objects"`
}

// DefaultSavedObjectApi is the default implementation of SavedObjectApi interface
type DefaultSavedObjectApi struct {
	client *resty.Client
}

// NewSavedObjectApi permit to get default implementation of SavedObjectApi interface
func NewSavedObjectApi(client *resty.Client) SavedObjectApi {
	return &DefaultSavedObjectApi{
		client: client,
	}
}

func (h DefaultSavedObjectApi) Get(tenantId string, objectType string, objectId string) (object *Object, err error) {
	if objectType == "" {
		return nil, NewAPIError(600, "You must provide the object type")
	}
	if objectId == "" {
		return nil, NewAPIError(600, "You must provide the object ID")
	}
	log.Debug("ObjectType: ", objectType)
	log.Debug("ID: ", objectId)
	log.Debug("Tenant: ", tenantId)

	// forge URL
	path := fmt.Sprintf("%s/%s/%s", basePathSavedObjectURL, objectType, objectId)

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// Fire
	resp, err := req.Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())

	}
	object = &Object{}
	err = json.Unmarshal(resp.Body(), object)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (h DefaultSavedObjectApi) BulkGet(tenantId string, bulks []SavedObjectBulkGetOption) (objects []Object, err error) {
	if len(bulks) == 0 {
		return nil, NewAPIError(600, "You must provide the bulks")
	}
	log.Debug("Tenant: ", tenantId)

	// forge URL
	path := fmt.Sprintf("%s/_bulk_get", basePathSavedObjectURL)

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// Fire
	resp, err := req.SetBody(bulks).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())

	}
	res := new(SavedObjectBulkResponse)
	err = json.Unmarshal(resp.Body(), res)
	if err != nil {
		return nil, err
	}

	return res.SavedObjects, nil
}

func (h DefaultSavedObjectApi) Find(tenantId string, option SavedObjectFindOption) (objects []Object, err error) {
	log.Debug("Tenant: ", tenantId)

	if option.Type == "" {
		return nil, NewAPIError(600, "You must provide the object type")
	}

	// forge URL
	path := fmt.Sprintf("%s/_find", basePathSavedObjectURL)

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// Handle parameter
	req.SetQueryParam("type", option.Type)
	if option.PerPage != 0 {
		req.SetQueryParam("per_page", strconv.Itoa(option.PerPage))
	}
	if option.Page != 0 {
		req.SetQueryParam("page", strconv.Itoa(option.Page))
	}
	if option.Search != "" {
		req.SetQueryParam("search", option.Search)
	}
	if option.DefaultSearchOperator != "" {
		req.SetQueryParam("default_search_operator", option.DefaultSearchOperator)
	}
	if len(option.SearchFields) > 0 {
		req.SetQueryParam("search_fields", strings.Join(option.SearchFields, ","))
	}
	if len(option.Fields) > 0 {
		req.SetQueryParam("fields", strings.Join(option.Fields, ","))
	}
	if option.SortField != "" {
		req.SetQueryParam("sort_field", option.SortField)
	}
	if option.HasReference != "" {
		req.SetQueryParam("has_reference", option.HasReference)
	}
	if option.Filter != "" {
		req.SetQueryParam("filter", option.Filter)
	}

	// Fire
	resp, err := req.Get(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())

	}
	findResponse := new(SavedObjectFindResponse)
	err = json.Unmarshal(resp.Body(), findResponse)
	if err != nil {
		return nil, err
	}

	return findResponse.SavedObjects, nil
}

func (h DefaultSavedObjectApi) Create(tenantId string, object *Object, overwrite bool) (createdObject *Object, err error) {
	log.Debug("Tenant: ", tenantId)
	log.Debug("Overwrite: ", overwrite)

	if object == nil {
		return nil, NewAPIError(600, "You must provide the object to create")
	}

	if object.Type == "" {
		return nil, NewAPIError(600, "You must provide the object type to create")
	}

	// forge URL
	var path string
	copyOject := &Object{
		Attributes: object.Attributes,
		References: object.References,
	}

	if object.Id == "" {
		path = fmt.Sprintf("%s/%s", basePathSavedObjectURL, object.Type)
	} else {
		path = fmt.Sprintf("%s/%s/%s", basePathSavedObjectURL, object.Type, object.Id)
	}

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// Handle parameter
	if overwrite {
		req.SetQueryParam("overwrite", "true")
	}

	// Fire
	resp, err := req.SetBody(copyOject).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())

	}
	createdObject = &Object{}
	err = json.Unmarshal(resp.Body(), createdObject)
	if err != nil {
		return nil, err
	}

	return createdObject, nil
}

func (h DefaultSavedObjectApi) BulkCreate(tenantId string, overwrite bool, bulks []Object) (createdObjects []Object, err error) {
	log.Debug("Tenant: ", tenantId)
	log.Debug("Overwrite: ", overwrite)

	if len(bulks) == 0 {
		return nil, NewAPIError(600, "You must provide the bulks")
	}

	// forge URL
	path := fmt.Sprintf("%s/_bulk_create", basePathSavedObjectURL)

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// Handle parameter
	if overwrite {
		req.SetQueryParam("overwrite", "true")
	}

	// Fire
	resp, err := req.SetBody(bulks).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())

	}
	res := new(SavedObjectBulkResponse)
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	return res.SavedObjects, nil
}

func (h DefaultSavedObjectApi) Update(tenantId string, object *Object) (updatedObject *Object, err error) {
	log.Debug("Tenant: ", tenantId)

	if object == nil {
		return nil, NewAPIError(600, "You must provide the object to update")
	}

	if object.Id == "" {
		return nil, NewAPIError(600, "You must provide the object ID to update")
	}

	if object.Type == "" {
		return nil, NewAPIError(600, "You must provide the object Type to update")
	}

	// forge URL
	path := fmt.Sprintf("%s/%s/%s", basePathSavedObjectURL, object.Type, object.Id)

	copyObject := &Object{
		Attributes: object.Attributes,
		References: object.References,
	}

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// Fire
	resp, err := req.SetBody(copyObject).Put(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())

	}
	updatedObject = &Object{}
	err = json.Unmarshal(resp.Body(), updatedObject)
	if err != nil {
		return nil, err
	}

	return updatedObject, nil
}

func (h DefaultSavedObjectApi) Delete(tenantId string, objectType string, objectId string, force bool) (err error) {
	log.Debug("Tenant: ", tenantId)

	if objectId == "" {
		return NewAPIError(600, "You must provide the object ID to update")
	}

	if objectType == "" {
		return NewAPIError(600, "You must provide the object Type to update")
	}

	// forge URL
	path := fmt.Sprintf("%s/%s/%s", basePathSavedObjectURL, objectType, objectId)

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// Handle parameters
	if force {
		req.SetQueryParam("force", "true")
	}

	// Fire
	resp, err := req.Delete(path)
	if err != nil {
		return err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil
		}
		return NewAPIError(resp.StatusCode(), resp.Status())

	}

	return nil
}

func (h DefaultSavedObjectApi) Export(tenantId string, option SavedObjectExportOption) (data []byte, err error) {
	log.Debug("Tenant: ", tenantId)

	// forge URL
	path := fmt.Sprintf("%s/_export", basePathSavedObjectURL)

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// Fire
	resp, err := req.SetBody(option).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())

	}

	return resp.Body(), nil
}

func (h DefaultSavedObjectApi) Import(tenantId string, option SavedObjectImportOption, data []byte) (res map[string]any, err error) {
	log.Debug("Tenant: ", tenantId)

	// forge URL
	path := fmt.Sprintf("%s/_import", basePathSavedObjectURL)

	// handle tenant
	req := h.client.R()
	if tenantId != "" {
		req.SetHeader("securitytenant", tenantId)
	}

	// handle parameters
	if option.CreateNewCopies {
		req.SetQueryParam("createNewCopies", "true")
	} else if option.Overwrite {
		req.SetQueryParam("overwrite", "true")
	}

	// Fire
	resp, err := req.SetFileReader("file", "file.ndjson", bytes.NewReader(data)).Post(path)
	if err != nil {
		return nil, err
	}
	log.Debug("Response: ", resp)
	if resp.StatusCode() >= 300 {
		if resp.StatusCode() == 404 {
			return nil, nil
		}
		return nil, NewAPIError(resp.StatusCode(), resp.Status())

	}

	res = map[string]any{}
	err = json.Unmarshal(resp.Body(), &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
