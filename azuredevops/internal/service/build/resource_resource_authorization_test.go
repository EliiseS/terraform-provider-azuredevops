// +build all resource_resource_authorization
// +build !exclude_resource_authorization

package build

// The tests in this file use the mock clients in mock_client.go to mock out
// the Azure DevOps client operations.

import (
	"context"
	"errors"
	"testing"

	"github.com/terraform-providers/terraform-provider-azuredevops/azdosdkmocks"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/microsoft/azure-devops-go-api/azuredevops/build"
	"github.com/stretchr/testify/require"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/client"
	"github.com/terraform-providers/terraform-provider-azuredevops/azuredevops/internal/utils/converter"
)

var projectId = "projectid"
var endpointId = uuid.New()

var resourceReferenceAuthorized = build.DefinitionResourceReference{
	Authorized: converter.Bool(true),
	Id:         converter.String(endpointId.String()),
	Name:       nil,
	Type:       converter.String("endpoint"),
}

var resourceReferenceNotAuthorized = build.DefinitionResourceReference{
	Authorized: converter.Bool(false),
	Id:         converter.String(endpointId.String()),
	Name:       nil,
	Type:       converter.String("endpoint"),
}

func TestResourceAuthorization_FlattenExpand_RoundTrip(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, ResourceResourceAuthorization().Schema, nil)
	flattenAuthorizedResource(resourceData, &resourceReferenceAuthorized, projectId)

	resourceReferenceAfterRoundtrip, projectIdAfterRoundtrip, err := expandAuthorizedResource(resourceData)
	require.Nil(t, err)
	require.Equal(t, resourceReferenceAuthorized, *resourceReferenceAfterRoundtrip)
	require.Equal(t, projectId, projectIdAfterRoundtrip)
}

func TestResourceAuthorization_Create_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r, resourceData, clients := prepareForCreateOrUpdate(t, ctrl, "CreateResourceAuthorization() Failed")

	err := r.Create(resourceData, clients)
	require.Contains(t, err.Error(), "CreateResourceAuthorization() Failed")
}

func TestResourceAuthorization_Update_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r, resourceData, clients := prepareForCreateOrUpdate(t, ctrl, "UpdateResourceAuthorization() Failed")

	err := r.Update(resourceData, clients)
	require.Contains(t, err.Error(), "UpdateResourceAuthorization() Failed")
}

func prepareForCreateOrUpdate(t *testing.T, ctrl *gomock.Controller, expectedMessage string) (*schema.Resource, *schema.ResourceData, *client.AggregatedClient) {
	r := ResourceResourceAuthorization()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuthorizedResource(resourceData, &resourceReferenceAuthorized, projectId)

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{BuildClient: buildClient, Ctx: context.Background()}

	expectedArgs := build.AuthorizeProjectResourcesArgs{
		Resources: &[]build.DefinitionResourceReference{resourceReferenceAuthorized},
		Project:   &projectId,
	}
	buildClient.
		EXPECT().
		AuthorizeProjectResources(clients.Ctx, expectedArgs).
		Return(nil, errors.New(expectedMessage)).
		Times(1)
	return r, resourceData, clients
}

func TestResourceAuthorization_Read_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceResourceAuthorization()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuthorizedResource(resourceData, &resourceReferenceAuthorized, projectId)

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{BuildClient: buildClient, Ctx: context.Background()}

	expectedArgs := build.GetProjectResourcesArgs{
		Project: &projectId,
		Type:    resourceReferenceAuthorized.Type,
		Id:      resourceReferenceAuthorized.Id,
	}
	buildClient.
		EXPECT().
		GetProjectResources(clients.Ctx, expectedArgs).
		Return(nil, errors.New("ReadResourceAuthorization() Failed")).
		Times(1)

	err := r.Read(resourceData, clients)
	require.Contains(t, err.Error(), "ReadResourceAuthorization() Failed")
}

func TestResourceAuthorization_Delete_DoesNotSwallowError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	r := ResourceResourceAuthorization()
	resourceData := schema.TestResourceDataRaw(t, r.Schema, nil)
	flattenAuthorizedResource(resourceData, &resourceReferenceNotAuthorized, projectId)

	buildClient := azdosdkmocks.NewMockBuildClient(ctrl)
	clients := &client.AggregatedClient{BuildClient: buildClient, Ctx: context.Background()}

	expectedArgs := build.AuthorizeProjectResourcesArgs{
		Resources: &[]build.DefinitionResourceReference{resourceReferenceNotAuthorized},
		Project:   &projectId,
	}
	buildClient.
		EXPECT().
		AuthorizeProjectResources(clients.Ctx, expectedArgs).
		Return(nil, errors.New("DeleteResourceAuthorization() Failed")).
		Times(1)

	err := r.Delete(resourceData, clients)
	require.Contains(t, err.Error(), "DeleteResourceAuthorization() Failed")
}
