package api

import (
	"bytes"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"
)

func (s *ApiTestSuite) TestSaveObject() {
	// Create
	object := &Object{
		Type: "index-pattern",
		Id:   "test",
		Attributes: map[string]any{
			"title": "test-pattern",
		},
	}
	resCreate, err := s.Api.SavedObject().Create("", object, true)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotNil(s.T(), resCreate)

	// Create on tenant
	object = &Object{
		Type: "index-pattern",
		Id:   "test-tenant",
		Attributes: map[string]any{
			"title": "test-pattern",
		},
	}
	resCreateTenant, err := s.Api.SavedObject().Create("test", object, true)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotNil(s.T(), resCreateTenant)

	// Bulks create some objects
	bulksCreate := []Object{
		{
			Type: "index-pattern",
			Id:   "test2",
			Attributes: map[string]any{
				"title": "test-pattern2",
			},
		},
		{
			Type: "dashboard",
			Id:   "be3733a0-9efe-11e7-acb3-3dab96693fab",
			Attributes: map[string]any{
				"title": "test",
			},
		},
	}
	resBulkCreate, err := s.Api.SavedObject().BulkCreate("", true, bulksCreate)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 2, len(resBulkCreate))

	// Bulks create some objects on tenant
	bulksCreate = []Object{
		{
			Type: "index-pattern",
			Id:   "test2_tenant",
			Attributes: map[string]any{
				"title": "test-pattern2",
			},
		},
		{
			Type: "dashboard",
			Id:   "be3733a0-9efe-11e7-acb3-3dab96693fab_tenant",
			Attributes: map[string]any{
				"title": "test",
			},
		},
	}
	resBulkCreateTenant, err := s.Api.SavedObject().BulkCreate("test", true, bulksCreate)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 2, len(resBulkCreateTenant))

	// Get object
	resGet, err := s.Api.SavedObject().Get("", "dashboard", "be3733a0-9efe-11e7-acb3-3dab96693fab")
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotNil(s.T(), resGet)

	// Get object on tenant
	resGetTenant, err := s.Api.SavedObject().Get("test", "dashboard", "be3733a0-9efe-11e7-acb3-3dab96693fab_tenant")
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotNil(s.T(), resGetTenant)

	// Get bulk
	resGetBulk, err := s.Api.SavedObject().BulkGet("", []SavedObjectBulkGetOption{
		{
			Type: "index-pattern",
			Id:   "test2",
		},
		{
			Type: "dashboard",
			Id:   "be3733a0-9efe-11e7-acb3-3dab96693fab",
		},
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 2, len(resGetBulk))

	// Get bulk on tenant
	resGetBulkTenant, err := s.Api.SavedObject().BulkGet("test", []SavedObjectBulkGetOption{
		{
			Type: "index-pattern",
			Id:   "test2_tenant",
		},
		{
			Type: "dashboard",
			Id:   "be3733a0-9efe-11e7-acb3-3dab96693fab_tenant",
		},
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 2, len(resGetBulkTenant))

	// Find
	resFind, err := s.Api.SavedObject().Find("", SavedObjectFindOption{
		Type:         "index-pattern",
		SearchFields: []string{"title"},
		Search:       "index-pattern",
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 1, len(resFind))

	// Find on tenant
	resFindTenant, err := s.Api.SavedObject().Find("test", SavedObjectFindOption{
		Type:         "index-pattern",
		SearchFields: []string{"title"},
		Search:       "index-pattern",
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 1, len(resFindTenant))

	// Update object
	resGet.Attributes["title"] = "test2"
	resUpdate, err := s.Api.SavedObject().Update("", resGet)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotNil(s.T(), resUpdate)

	// Update object on tenant
	resGetTenant.Attributes["title"] = "test2"
	resUpdateTenant, err := s.Api.SavedObject().Update("test", resGetTenant)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotNil(s.T(), resUpdateTenant)

	// Export object
	resExport, err := s.Api.SavedObject().Export("", SavedObjectExportOption{
		Type:                 "index-pattern",
		ExcludeExportDetails: ptr.To[bool](true),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 2, len(bytes.Split(resExport, []byte("\n"))))

	resExport, err = s.Api.SavedObject().Export("", SavedObjectExportOption{
		Objects: []SavedObjectBulkGetOption{
			{
				Type: "index-pattern",
				Id:   "test2",
			},
		},
		ExcludeExportDetails: ptr.To[bool](true),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 1, len(bytes.Split(resExport, []byte("\n"))))

	// Export object on tenant
	resExportTenant, err := s.Api.SavedObject().Export("test", SavedObjectExportOption{
		Type:                 "index-pattern",
		ExcludeExportDetails: ptr.To[bool](true),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 2, len(bytes.Split(resExportTenant, []byte("\n"))))

	resExportTenant, err = s.Api.SavedObject().Export("test", SavedObjectExportOption{
		Objects: []SavedObjectBulkGetOption{
			{
				Type: "index-pattern",
				Id:   "test2_tenant",
			},
		},
		ExcludeExportDetails: ptr.To[bool](true),
	})
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), 1, len(bytes.Split(resExportTenant, []byte("\n"))))

	// Import object
	resImport, err := s.Api.SavedObject().Import("", SavedObjectImportOption{Overwrite: true}, resExport)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotEmpty(s.T(), resImport)

	// Import object on tenant
	resImportTenant, err := s.Api.SavedObject().Import("test", SavedObjectImportOption{Overwrite: true}, resExportTenant)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.NotEmpty(s.T(), resImportTenant)

	// Delete object
	if err = s.Api.SavedObject().Delete("", "index-pattern", "test", false); err != nil {
		s.T().Fatal(err)
	}

	// Delete object on tenant
	if err = s.Api.SavedObject().Delete("test", "index-pattern", "test", false); err != nil {
		s.T().Fatal(err)
	}
}
