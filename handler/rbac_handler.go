package handler

import (
	"context"
	"coresamples/ent"
	pb "coresamples/proto"
	"coresamples/service"
	"errors"
	errors2 "go-micro.dev/v4/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"strings"
)

type RBACHandler struct {
	RBACService     service.IRBACService
	allowedServices []string
}

func (rh *RBACHandler) CheckPermission(ctx context.Context, request *pb.CheckPermissionRequest, response *pb.CheckPermissionResponse) error {
	var err error
	var granted bool
	action := strings.ToUpper(request.GetAction())
	resource := strings.ToUpper(request.GetResource())
	if request.DomainType != "" {
		granted, err = rh.RBACService.CheckPermissionInDomain(request.GetAccountId(), action, resource, strings.ToLower(request.DomainType), request.DomainId, ctx)
	} else {
		granted, err = rh.RBACService.CheckPermission(request.GetAccountId(), action, resource, ctx)
	}
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CheckCreateRolePermission(ctx context.Context, request *pb.CheckCreateRoleRequest, response *pb.CheckPermissionResponse) error {
	granted, err := rh.RBACService.CheckCreateRolePermission(request.RoleTypeName, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CheckDeleteRolePermission(ctx context.Context, request *pb.CheckDeleteRoleRequest, response *pb.CheckPermissionResponse) error {
	roleName := strings.ToUpper(request.RoleName)
	granted, err := rh.RBACService.CheckDeleteRolePermission(roleName, request.ClinicId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CheckAddPermissionPermission(ctx context.Context, request *pb.CheckEditRolePermissionRequest, response *pb.CheckPermissionResponse) error {
	resource := strings.ToUpper(request.GetResourceName())
	roleName := strings.ToUpper(request.RoleName)
	granted, err := rh.RBACService.CheckAddPermissionPermission(request.CallerAccountId, resource, roleName, request.ClinicId, ctx)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CheckDeletePermissionPermission(ctx context.Context, request *pb.CheckEditRolePermissionRequest, response *pb.CheckPermissionResponse) error {
	resource := strings.ToUpper(request.GetResourceName())
	roleName := strings.ToUpper(request.RoleName)
	granted, err := rh.RBACService.CheckDeletePermissionPermission(request.CallerAccountId, resource, roleName, request.ClinicId, ctx)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CheckAssignRolePermission(ctx context.Context, request *pb.CheckAssignRolePermissionRequest, response *pb.CheckPermissionResponse) error {
	roleName := strings.ToUpper(request.RoleName)
	granted, err := rh.RBACService.CheckAssignRolePermission(roleName, request.ClinicId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CheckRemoveRolePermission(ctx context.Context, request *pb.CheckRemoveRolePermissionRequest, response *pb.CheckPermissionResponse) error {
	roleName := strings.ToUpper(request.RoleName)
	granted, err := rh.RBACService.CheckRemoveRolePermission(roleName, request.ClinicId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CheckResourcePermission(ctx context.Context, request *pb.CheckResourcePermissionRequest, response *pb.CheckPermissionResponse) error {
	granted, err := rh.RBACService.CheckResourcePermission(request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CheckActionPermission(ctx context.Context, request *pb.CheckActionRequest, response *pb.CheckPermissionResponse) error {
	granted, err := rh.RBACService.CheckActionPermission(request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
	} else {
		response.Message = MsgSuccess
	}
	response.Granted = granted
	return err
}

func (rh *RBACHandler) CreateRole(ctx context.Context, request *pb.CreateRoleRequest, response *pb.SimpleResponse) error {
	roleName := strings.ToUpper(request.RoleName)
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	err := rh.RBACService.CreateRole(roleName, request.RoleTypeName, request.ClinicId, request.CallerAccountId, ctx)
	if err == nil {
		response.Message = MsgSuccess
		return nil
	}

	response.Message = err.Error()
	return err
}

func (rh *RBACHandler) DeleteRole(ctx context.Context, request *pb.DeleteRoleRequest, response *pb.DeleteRoleResponse) error {
	roleName := strings.ToUpper(request.RoleName)
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	num, err := rh.RBACService.DeleteRole(roleName, request.ClinicId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		response.NumDeleted = 0
		return err
	}
	response.Message = MsgSuccess
	response.NumDeleted = int32(num)
	return nil
}

func (rh *RBACHandler) AddPermissionToRole(ctx context.Context, request *pb.AddPermissionToRoleRequest, response *pb.SimpleResponse) error {
	resourceName := strings.ToUpper(request.GetResourceName())
	roleName := strings.ToUpper(request.RoleName)
	actionName := strings.ToUpper(request.GetActionName())
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	err := rh.RBACService.AddPermissionToRole(actionName, resourceName, roleName, request.ClinicId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) AddPermissionToUser(ctx context.Context, request *pb.AddPermissionToUserRequest, response *pb.SimpleResponse) error {
	resourceName := strings.ToUpper(request.GetResourceName())
	actionName := strings.ToUpper(request.GetActionName())
	err := rh.RBACService.AddPermissionToUser(actionName, resourceName, request.AccountId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) DeletePermissionFromUser(ctx context.Context, request *pb.DeletePermissionFromUserRequest, response *pb.SimpleResponse) error {
	resourceName := strings.ToUpper(request.GetResourceName())
	actionName := strings.ToUpper(request.GetActionName())
	err := rh.RBACService.DeletePermissionFromUser(actionName, resourceName, request.AccountId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) DeletePermissionFromRole(ctx context.Context, request *pb.DeletePermissionFromRoleRequest, response *pb.SimpleResponse) error {
	resourceName := strings.ToUpper(request.GetResourceName())
	roleName := strings.ToUpper(request.RoleName)
	actionName := strings.ToUpper(request.GetActionName())
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	err := rh.RBACService.DeletePermissionFromRole(actionName, resourceName, roleName, request.ClinicId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) GetAccountRoleInternalNames(ctx context.Context, request *pb.GetAccountRoleInternalNamesRequest, response *pb.GetAccountRoleInternalNamesResponse) error {
	names, err := rh.RBACService.GetAccountRoleInternalNames(request.AccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		response.InternalNames = names
		return err
	}
	response.Message = MsgSuccess
	response.InternalNames = names
	return nil
}

func (rh *RBACHandler) GetAccountRoles(ctx context.Context, request *pb.GetAccountRolesRequest, response *pb.GetRolesResponse) error {
	var roles []*ent.RBACRoles
	var err error
	if request.DomainType != "" {
		roles, err = rh.RBACService.GetAccountRolesInDomain(request.AccountId, strings.ToLower(request.DomainType), request.DomainId, ctx)
	} else {
		roles, err = rh.RBACService.GetAccountRoles(request.AccountId, ctx)
	}

	if err != nil {
		response.Message = err.Error()
		response.Roles = []*pb.Role{}
		return err
	}
	var outRoles []*pb.Role
	for _, r := range roles {
		outRole := &pb.Role{
			Id:           int32(r.ID),
			Name:         r.Name,
			InternalName: r.InternalName,
			ClinicId:     r.ClinicID,
			Type:         r.Type.String(),
		}
		outRoles = append(outRoles, outRole)
	}
	response.Message = MsgSuccess
	response.Roles = outRoles
	return nil
}

func (rh *RBACHandler) GetAccountPermissions(ctx context.Context, request *pb.GetAccountPermissionsRequest, response *pb.GetPermissionsResponse) error {
	var permissionMap map[string][]string
	var err error
	if request.DomainType != "" {
		permissionMap, err = rh.RBACService.GetAccountPermissionsInDomain(request.AccountId, strings.ToLower(request.DomainType), request.DomainId, ctx)
	} else {
		permissionMap, err = rh.RBACService.GetAccountPermissions(request.AccountId, ctx)
	}

	var permissions []*pb.Permissions
	if err != nil {
		response.Permissions = permissions
		response.Message = err.Error()
		return err
	}
	for action, resources := range permissionMap {
		p := &pb.Permissions{
			Action:    action,
			Resources: resources,
		}
		permissions = append(permissions, p)
	}

	response.Message = MsgSuccess
	response.Permissions = permissions
	return nil
}

func (rh *RBACHandler) GetRolePermissions(ctx context.Context, request *pb.GetRolePermissionsRequest, response *pb.GetPermissionsResponse) error {
	roleName := strings.ToUpper(request.RoleName)
	permissionMap, err := rh.RBACService.GetRolePermissions(roleName, request.ClinicId, ctx)
	if err != nil {
		response.Permissions = []*pb.Permissions{}
		response.Message = err.Error()
		return err
	}

	response.Message = MsgSuccess
	response.Permissions = rh.convertPermissions(permissionMap)
	return nil
}

func (rh *RBACHandler) GetNonPrivateRoles(ctx context.Context, request *emptypb.Empty, response *pb.GetRolesResponse) error {
	roles, err := rh.RBACService.GetNonPrivateRoles(ctx)
	if err != nil {
		response.Message = err.Error()
		response.Roles = []*pb.Role{}
		return err
	}

	response.Message = MsgSuccess
	response.Roles = rh.convertRoles(roles)
	return nil
}

func (rh *RBACHandler) GetRolesByType(ctx context.Context, request *pb.GetRolesByTypeRequest, response *pb.GetRolesResponse) error {
	roles, err := rh.RBACService.GetRolesByType(request.TypeName, ctx)
	if err != nil {
		response.Message = err.Error()
		response.Roles = []*pb.Role{}
		return err
	}

	response.Message = MsgSuccess
	response.Roles = rh.convertRoles(roles)
	return nil
}

func (rh *RBACHandler) GetRoleTypeByInternalName(ctx context.Context, request *pb.GetRoleTypeByInternalNameRequest, response *pb.GetRoleTypeByInternalNameResponse) error {
	roleType, err := rh.RBACService.GetRoleTypeByInternalName(request.InternalName, ctx)
	if err != nil {
		response.Message = err.Error()
		response.RoleType = ""
		return err
	}

	response.Message = MsgSuccess
	response.RoleType = roleType
	return nil
}

func (rh *RBACHandler) AssignRoleToAccount(ctx context.Context, request *pb.AssignRoleToAccountRequest, response *pb.SimpleResponse) error {
	var err error
	roleName := strings.ToUpper(request.RoleName)
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}

	if request.DomainType != "" {
		err = rh.RBACService.AssignRoleToAccountInDomain(roleName, request.AccountId, strings.ToLower(request.DomainType), request.DomainId, request.CallerAccountId, ctx)
	} else {
		err = rh.RBACService.AssignRoleToAccount(roleName, request.ClinicId, request.AccountId, request.CallerAccountId, ctx)
	}

	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) RemoveRoleFromAccount(ctx context.Context, request *pb.RemoveRoleFromAccountRequest, response *pb.SimpleResponse) error {
	var err error
	roleName := strings.ToUpper(request.RoleName)
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}

	if request.DomainType != "" {
		err = rh.RBACService.RemoveRoleFromAccountInDomain(roleName, request.AccountId, strings.ToLower(request.DomainType), request.DomainId, request.CallerAccountId, ctx)
	} else {
		err = rh.RBACService.RemoveRoleFromAccount(roleName, request.ClinicId, request.AccountId, request.CallerAccountId, ctx)
	}

	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) CreateResource(ctx context.Context, request *pb.CreateResourceRequest, response *pb.SimpleResponse) error {
	resourceName := strings.ToUpper(request.Name)
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors2.Unauthorized(strconv.Itoa(int(request.CallerAccountId)), response.Message)
	}
	err := rh.RBACService.CreateResource(resourceName, request.Description, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		return errors2.InternalServerError(strconv.Itoa(int(request.CallerAccountId)), err.Error())
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) DeleteResourceByName(ctx context.Context, request *pb.DeleteResourceByNameRequest, response *pb.DeleteResourceResponse) error {
	resourceName := strings.ToUpper(request.GetResourceName())
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	num, err := rh.RBACService.DeleteResourceByName(resourceName, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		response.NumDeleted = 0
		return err
	}
	response.NumDeleted = int32(num)
	response.Message = MsgSuccess
	return err
}

func (rh *RBACHandler) DeleteResourceById(ctx context.Context, request *pb.DeleteResourceByIdRequest, response *pb.DeleteResourceResponse) error {
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	num, err := rh.RBACService.DeleteResourceById(request.ResourceId, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		response.NumDeleted = 0
		return err
	}
	response.NumDeleted = int32(num)
	response.Message = MsgSuccess
	return err
}

func (rh *RBACHandler) GetResources(ctx context.Context, request *emptypb.Empty, response *pb.GetResourcesResponse) error {
	resources, err := rh.RBACService.GetResources(ctx)
	if err != nil {
		response.Message = err.Error()
		response.Resources = []*pb.Resource{}
		return err
	}
	response.Message = MsgSuccess
	response.Resources = rh.convertResources(resources)
	return nil
}

func (rh *RBACHandler) GetResourceById(ctx context.Context, request *pb.GetResourceByIdRequest, response *pb.GetResourceResponse) error {
	resource, err := rh.RBACService.GetResourceById(request.ResourceId, ctx)
	if err != nil {
		response.Resource = nil
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	if resource != nil {
		response.Resource = &pb.Resource{
			Id:          int32(resource.ID),
			Name:        resource.Name,
			Description: resource.Description,
		}
	} else {
		response.Resource = nil
	}
	return nil
}

func (rh *RBACHandler) GetResourceDescription(ctx context.Context, request *pb.GetResourceDescriptionRequest, response *pb.GetResourceDescriptionResponse) error {
	resource := strings.ToUpper(request.ResourceName)
	description, err := rh.RBACService.GetResourceDescription(resource, ctx)
	if err != nil {
		response.Message = err.Error()
		response.Description = description
		return err
	}
	response.Message = MsgSuccess
	response.Description = description
	return nil
}

func (rh *RBACHandler) UpdateResourceDescription(ctx context.Context, request *pb.UpdateResourceDescriptionRequest, response *pb.SimpleResponse) error {
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	resourceName := strings.ToUpper(request.GetResourceName())
	err := rh.RBACService.UpdateResourceDescription(resourceName, request.Description, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) CreateAction(ctx context.Context, request *pb.ActionRequest, response *pb.SimpleResponse) error {
	actionName := strings.ToUpper(request.GetActionName())
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	err := rh.RBACService.CreateAction(actionName, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) DeleteAction(ctx context.Context, request *pb.ActionRequest, response *pb.SimpleResponse) error {
	actionName := strings.ToUpper(request.GetActionName())
	if !rh.serviceAllowed(request.ServiceName) {
		response.Message = "Service " + request.ServiceName + " not allowed"
		return errors.New(response.Message)
	}
	err := rh.RBACService.DeleteAction(actionName, request.CallerAccountId, ctx)
	if err != nil {
		response.Message = err.Error()
		return err
	}
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) GetActions(ctx context.Context, request *emptypb.Empty, response *pb.GetActionsResponse) error {
	actions, err := rh.RBACService.GetActions(ctx)
	if err != nil {
		response.Message = err.Error()
		response.Actions = []*pb.Action{}
		return err
	}
	response.Message = MsgSuccess
	response.Actions = rh.convertActions(actions)
	return nil
}

func (rh *RBACHandler) GetDefaultPermissions(ctx context.Context, request *emptypb.Empty, response *pb.GetDefaultPermissionsResponse) error {
	str, err := rh.RBACService.GetDefaultPermissions(ctx)
	if err != nil {
		response.Message = err.Error()
		response.CompressedPermissions = ""
		return err
	}
	response.CompressedPermissions = str
	response.Message = MsgSuccess
	return nil
}

func (rh *RBACHandler) convertRoles(roles []*ent.RBACRoles) []*pb.Role {
	var outRoles []*pb.Role
	for _, r := range roles {
		outRole := &pb.Role{
			Id:           int32(r.ID),
			Name:         r.Name,
			InternalName: r.InternalName,
			ClinicId:     r.ClinicID,
			Type:         r.Type.String(),
		}
		outRoles = append(outRoles, outRole)
	}
	return outRoles
}

func (rh *RBACHandler) convertPermissions(permissionMap map[string][]string) []*pb.Permissions {
	var permissions []*pb.Permissions
	for action, resources := range permissionMap {
		p := &pb.Permissions{
			Action:    action,
			Resources: resources,
		}
		permissions = append(permissions, p)
	}
	return permissions
}

func (rh *RBACHandler) convertResources(resources []*ent.RBACResources) []*pb.Resource {
	var outResources []*pb.Resource
	for _, r := range resources {
		outResource := &pb.Resource{
			Id:          int32(r.ID),
			Name:        r.Name,
			Description: r.Description,
		}
		outResources = append(outResources, outResource)
	}
	return outResources
}

func (rh *RBACHandler) convertActions(actions []*ent.RBACActions) []*pb.Action {
	var outActions []*pb.Action
	for _, a := range actions {
		outAction := &pb.Action{
			Id:   int32(a.ID),
			Name: a.Name,
		}
		outActions = append(outActions, outAction)
	}
	return outActions
}

func (h *RBACHandler) serviceAllowed(serviceName string) bool {
	for _, name := range h.allowedServices {
		if name == serviceName {
			return true
		}
	}
	return false
}
