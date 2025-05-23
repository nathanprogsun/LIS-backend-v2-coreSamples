package service

import (
	"bytes"
	"context"
	"coresamples/common"
	"coresamples/dbutils"
	"coresamples/ent"
	"coresamples/ent/rbacroles"
	"coresamples/util"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/casbin/casbin/v2"
	casbin_constant "github.com/casbin/casbin/v2/constant"
	casbin_model "github.com/casbin/casbin/v2/model"
	"github.com/casbin/ent-adapter"
	"github.com/klauspost/compress/zstd"
)

const (
	InvalidClinicId      = 0
	ActPos               = 1
	ObjPos               = 2
	ResourceRole         = "ROLE"
	ResourceExternalRole = "EXTERNALROLE"
	ResourceClinicRole   = "CLINICROLE"
	ActionCreate         = "CREATE"
	ActionDelete         = "DELETE"
	ActionEdit           = "EDIT"
	ActionAssign         = "ASSIGN"
	ActionRemove         = "REMOVE"
)

//go:embed rbac_model.conf
var configb string

type IRBACService interface {
	CheckPermission(accountId int32, action string, resource string, ctx context.Context) (bool, error)
	CheckCreateRolePermission(roleTypeName string, callerAccountId int32, ctx context.Context) (bool, error)
	CheckDeleteRolePermission(roleName string, clinicId int32, callerAccountId int32, ctx context.Context) (bool, error)
	CreateRole(roleName string, roleTypeName string, clinicId int32, callerAccountId int32, ctx context.Context) error
	DeleteRole(roleName string, clinicId int32, callerAccountId int32, ctx context.Context) (int, error)
	AddPermissionToRole(actionName string, resourceName string, roleName string, clinicId int32, callerAccountId int32, ctx context.Context) error
	AddPermissionToUser(actionName string, resourceName string, accountId int32, callerAccountId int32, ctx context.Context) error
	DeletePermissionFromRole(actionName string, resourceName string, roleName string, clinicId int32, callerAccountId int32, ctx context.Context) error
	DeletePermissionFromUser(actionName string, resourceName string, accountId int32, callerAccountId int32, ctx context.Context) error
	GetAccountRoleInternalNames(accountId int32, ctx context.Context) ([]string, error)
	GetAccountRoles(accountId int32, ctx context.Context) ([]*ent.RBACRoles, error)
	GetAccountPermissions(accountId int32, ctx context.Context) (map[string][]string, error)
	GetRolePermissions(roleName string, clinicId int32, ctx context.Context) (map[string][]string, error)
	GetNonPrivateRoles(ctx context.Context) ([]*ent.RBACRoles, error)
	GetRolesByType(roleTypeName string, ctx context.Context) ([]*ent.RBACRoles, error)
	GetRoleTypeByInternalName(name string, ctx context.Context) (string, error)
	CheckAssignRolePermission(roleName string, clinicId int32, callerAccountId int32, ctx context.Context) (bool, error)
	CheckRemoveRolePermission(roleName string, clinicId int32, callerAccountId int32, ctx context.Context) (bool, error)
	AssignRoleToAccount(roleName string, clinicId int32, accountId int32, callerAccountId int32, ctx context.Context) error
	RemoveRoleFromAccount(roleName string, clinicId int32, accountId int32, callerAccountId int32, ctx context.Context) error
	CheckResourcePermission(callerAccountId int32, ctx context.Context) (bool, error)
	CreateResource(name string, description string, callerAccountId int32, ctx context.Context) error
	DeleteResourceByName(name string, callerAccountId int32, ctx context.Context) (int, error)
	DeleteResourceById(id int32, callerAccountId int32, ctx context.Context) (int, error)
	GetResources(ctx context.Context) ([]*ent.RBACResources, error)
	GetResourceById(id int32, ctx context.Context) (*ent.RBACResources, error)
	GetResourceDescription(name string, ctx context.Context) (string, error)
	UpdateResourceDescription(name string, description string, callerAccountId int32, ctx context.Context) error
	CheckActionPermission(callerAccountId int32, ctx context.Context) (bool, error)
	CreateAction(name string, callerAccountId int32, ctx context.Context) error
	DeleteAction(name string, callerAccountId int32, ctx context.Context) error
	GetActions(ctx context.Context) ([]*ent.RBACActions, error)
	GetDefaultPermissions(ctx context.Context) (string, error)
	CheckAddPermissionPermission(callerAccountId int32, resourceName string, roleName string, clinicId int32, ctx context.Context) (bool, error)
	CheckDeletePermissionPermission(callerAccountId int32, resourceName string, roleName string, clinicId int32, ctx context.Context) (bool, error)
	GetAccountRolesInDomain(accountId int32, domainType string, domainID int32, ctx context.Context) ([]*ent.RBACRoles, error)
	GetAccountPermissionsInDomain(accountId int32, domainType string, domainID int32, ctx context.Context) (map[string][]string, error)
	AssignRoleToAccountInDomain(roleName string, accountId int32, domainType string, domainID int32, callerAccountId int32, ctx context.Context) error
	RemoveRoleFromAccountInDomain(roleName string, accountId int32, domainType string, domainID int32, callerAccountId int32, ctx context.Context) error
	CheckPermissionInDomain(accountId int32, action string, resource string, domainType string, domainID int32, ctx context.Context) (bool, error)
}

type RBACService struct {
	Service
	enforcer *casbin.Enforcer
}

func newRBACService(dbClient *ent.Client, redisClient *common.RedisClient, driverName string, dataSource string) IRBACService {
	s := &RBACService{
		Service: InitService(dbClient, redisClient),
	}

	adapter, err := entadapter.NewAdapter(driverName, dataSource)
	if err != nil {
		common.Fatal(err)
	}

	m, err := casbin_model.NewModelFromString(configb)
	if err != nil {
		common.Fatal(err)
	}

	s.enforcer, err = casbin.NewEnforcer(m, adapter)
	if err != nil {
		common.Fatal(err)
	}

	s.enforcer.SetFieldIndex("p", casbin_constant.DomainIndex, 3)

	s.enforcer.EnableAutoSave(true)
	return s
}

func (s *RBACService) CreateRole(roleName string, roleTypeName string, clinicId int32, callerAccountId int32, ctx context.Context) error {
	roleType, err := util.GetRoleTypeEnum(roleTypeName)
	if err != nil {
		return err
	}
	if (roleType == rbacroles.TypeClinic && clinicId == InvalidClinicId) ||
		(roleType != rbacroles.TypeClinic && clinicId != InvalidClinicId) {
		return errors.New("clinic role should have valid clinic id, non-clinic role shouldn't have valid clinic id")
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionCreate, roleType, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to create " + roleTypeName + " role")
	}
	internalName := getInternalName(roleName, clinicId)
	return dbutils.CreateRole(roleName, internalName, clinicId, roleType, s.dbClient, ctx)
}

func (s *RBACService) CheckCreateRolePermission(roleTypeName string, callerAccountId int32, ctx context.Context) (bool, error) {
	roleType, err := util.GetRoleTypeEnum(roleTypeName)
	if err != nil {
		return false, err
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionCreate, roleType, ctx)
	if err != nil {
		return false, err
	}
	if !granted {
		return false, nil
	}
	return true, nil
}

func (s *RBACService) DeleteRole(roleName string, clinicId int32, callerAccountId int32, ctx context.Context) (int, error) {
	roleType, err := s.getRoleType(roleName, clinicId, ctx)
	if err != nil || roleType.String() == "" {
		return 0, err
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionDelete, roleType, ctx)
	if err != nil {
		return 0, err
	}
	if !granted {
		return 0, errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to delete role " + roleName)
	}
	internalName := getInternalName(roleName, clinicId)

	// Ideally when we delete role in rule table and in role resource table,
	// we would want to them to be wrapped in the same transaction.
	// But we are forced to use different ent client, so we'll have to make do

	// the worst thing that can happen would be the rules are deleted but the resource entry is not,
	// in that case we'll have to fix it manually
	if _, err = s.enforcer.DeleteRole(internalName); err != nil {
		return 0, err
	}
	return dbutils.DeleteRoleByInternalName(internalName, s.dbClient, ctx)
}

func (s *RBACService) CheckDeleteRolePermission(roleName string, clinicId int32, callerAccountId int32, ctx context.Context) (bool, error) {
	roleType, err := s.getRoleType(roleName, clinicId, ctx)
	if err != nil || roleType.String() == "" {
		return false, err
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionDelete, roleType, ctx)
	if err != nil {
		return false, err
	}
	if !granted {
		return false, nil
	}
	return true, nil
}

func (s *RBACService) AddPermissionToRole(actionName string, resourceName string, roleName string, clinicId int32, callerAccountId int32, ctx context.Context) error {
	act, err := dbutils.FindActionByName(actionName, s.dbClient, ctx)
	if err != nil {
		return err
	}
	if act == nil {
		return errors.New("action " + actionName + " does not exist")
	}
	rsc, err := dbutils.FindResourceByName(resourceName, s.dbClient, ctx)
	if err != nil {
		return err
	}
	if rsc == nil {
		return errors.New("resource " + resourceName + " does not exist")
	}
	if (actionName == ActionAssign || actionName == ActionRemove) &&
		(resourceName != ResourceRole && resourceName != ResourceExternalRole && resourceName != ResourceClinicRole) {
		return errors.New("action " + actionName + " can only be granted on role related resource, use other action instead")
	}
	internalName := getInternalName(roleName, clinicId)
	role, err := dbutils.FindRoleByInternalName(internalName, s.dbClient, ctx)
	if role == nil {
		return errors.New("role " + roleName + " does not exist")
	}

	roleType := role.Type
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionEdit, roleType, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to edit role " + roleName)
	}

	if resourceName == ResourceRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		if !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource Role")
		}
	} else if resourceName == ResourceExternalRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		externalRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceExternalRole, ctx)
		if err != nil {
			return err
		}
		if !externalRoleResourceGranted && !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource ExternalRole")
		}
	} else if resourceName == ResourceClinicRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		clinicRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceClinicRole, ctx)
		if err != nil {
			return err
		}
		if !clinicRoleResourceGranted && !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource ClinicRole")
		}
	}

	_, err = s.enforcer.AddPermissionForUser(internalName, actionName, resourceName)
	return err
}

func (s *RBACService) AddPermissionToUser(actionName string, resourceName string, accountId int32, callerAccountId int32, ctx context.Context) error {
	act, err := dbutils.FindActionByName(actionName, s.dbClient, ctx)
	if err != nil {
		return err
	}
	if act == nil {
		return errors.New("action " + actionName + " does not exist")
	}
	rsc, err := dbutils.FindResourceByName(resourceName, s.dbClient, ctx)
	if err != nil {
		return err
	}
	if rsc == nil {
		return errors.New("resource " + resourceName + " does not exist")
	}
	if resourceName == ResourceRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		if !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource Role")
		}
	} else if resourceName == ResourceExternalRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		externalRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceExternalRole, ctx)
		if err != nil {
			return err
		}
		if !externalRoleResourceGranted && !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource ExternalRole")
		}
	} else if resourceName == ResourceClinicRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		clinicRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceClinicRole, ctx)
		if err != nil {
			return err
		}
		if !clinicRoleResourceGranted && !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource ClinicRole")
		}
	}

	_, err = s.enforcer.AddPermissionForUser(strconv.Itoa(int(accountId)), actionName, resourceName)
	return err
}

func (s *RBACService) DeletePermissionFromUser(actionName string, resourceName string, accountId int32, callerAccountId int32, ctx context.Context) error {
	if resourceName == ResourceRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		if !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource Role")
		}
	} else if resourceName == ResourceExternalRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		externalRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceExternalRole, ctx)
		if err != nil {
			return err
		}
		if !externalRoleResourceGranted && !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource ExternalRole")
		}
	} else if resourceName == ResourceClinicRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		clinicRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceClinicRole, ctx)
		if err != nil {
			return err
		}
		if !clinicRoleResourceGranted && !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource ClinicRole")
		}
	}
	_, err := s.enforcer.DeletePermissionForUser(strconv.Itoa(int(accountId)), actionName, resourceName)
	return err
}

func (s *RBACService) DeletePermissionFromRole(actionName string, resourceName string, roleName string, clinicId int32, callerAccountId int32, ctx context.Context) error {
	internalName := getInternalName(roleName, clinicId)
	roleType, err := s.getRoleType(roleName, clinicId, ctx)
	if err != nil {
		return err
	}
	if roleType.String() == "" {
		return errors.New("role " + roleName + " does not exist")
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionEdit, roleType, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to edit role " + roleName)
	}

	if resourceName == ResourceRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		if !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource Role")
		}
	} else if resourceName == ResourceExternalRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		externalRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceExternalRole, ctx)
		if err != nil {
			return err
		}
		if !externalRoleResourceGranted && !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource ExternalRole")
		}
	} else if resourceName == ResourceClinicRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return err
		}
		clinicRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceClinicRole, ctx)
		if err != nil {
			return err
		}
		if !clinicRoleResourceGranted && !roleResourceGranted {
			return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to add permission on resource ClinicRole")
		}
	}
	_, err = s.enforcer.DeletePermissionForUser(internalName, actionName, resourceName)
	return err
}

func (s *RBACService) GetAccountRoleInternalNames(accountId int32, ctx context.Context) ([]string, error) {
	roleNames, err := s.enforcer.GetImplicitRolesForUser(strconv.Itoa(int(accountId)))
	if err != nil {
		return nil, err
	}
	return roleNames, err
}

func (s *RBACService) GetAccountRoles(accountId int32, ctx context.Context) ([]*ent.RBACRoles, error) {
	names, err := s.GetAccountRoleInternalNames(accountId, ctx)
	if err != nil {
		return nil, err
	}
	return dbutils.FindRolesByInternalNames(names, s.dbClient, ctx)
}

func (s *RBACService) GetAccountPermissions(accountId int32, ctx context.Context) (map[string][]string, error) {
	plist, err := s.enforcer.GetImplicitPermissionsForUser(strconv.Itoa(int(accountId)))
	if err != nil {
		return map[string][]string{}, err
	}
	permissions := map[string][]string{}
	resources := map[string]map[string]bool{}
	for _, permission := range plist {
		act := permission[ActPos]
		obj := permission[ObjPos]
		if _, exist := resources[act]; !exist {
			resources[act] = make(map[string]bool)
		}
		if _, exist := resources[act][obj]; !exist {
			permissions[act] = append(permissions[act], obj)
			resources[act][obj] = true
		}
	}
	return permissions, nil
}

func (s *RBACService) GetRolePermissions(roleName string, clinicId int32, ctx context.Context) (map[string][]string, error) {
	internalName := getInternalName(roleName, clinicId)
	plist, err := s.enforcer.GetImplicitPermissionsForUser(internalName)
	if err != nil {
		return map[string][]string{}, err
	}
	permissions := map[string][]string{}
	for _, permission := range plist {
		act := permission[ActPos]
		obj := permission[ObjPos]
		permissions[act] = append(permissions[act], obj)
	}
	return permissions, nil
}

func (s *RBACService) GetNonPrivateRoles(ctx context.Context) ([]*ent.RBACRoles, error) {
	roles, err := dbutils.FindSharedRoles(s.dbClient, ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (s *RBACService) GetRolesByType(roleTypeName string, ctx context.Context) ([]*ent.RBACRoles, error) {
	roleType, err := util.GetRoleTypeEnum(roleTypeName)
	if err != nil {
		return nil, err
	}
	return dbutils.FindRolesByType(roleType, s.dbClient, ctx)
}

// GetRoleTypeByInternalName essentially this only works on non-clinic roles, with name identical with internal name
// if you know its internal name, you already knew its type, didn't you?
func (s *RBACService) GetRoleTypeByInternalName(name string, ctx context.Context) (string, error) {
	roleType, err := dbutils.GetRoleTypeByInternalName(name, s.dbClient, ctx)
	return roleType.String(), err
}

func (s *RBACService) CheckPermission(accountId int32, action string, resource string, ctx context.Context) (bool, error) {
	granted, err := s.enforcer.Enforce(strconv.Itoa(int(accountId)), action, resource)
	if err != nil {
		common.Error(err)
		return false, err
	}
	return granted, err
}

func (s *RBACService) CheckAssignRolePermission(roleName string, clinicId int32, callerAccountId int32, ctx context.Context) (bool, error) {
	roleType, err := s.getRoleType(roleName, clinicId, ctx)
	if err != nil || roleType.String() == "" {
		return false, err
	}
	return s.checkPermissionOnRole(callerAccountId, ActionAssign, roleType, ctx)
}

func (s *RBACService) AssignRoleToAccount(roleName string, clinicId int32, accountId int32, callerAccountId int32, ctx context.Context) error {
	internalName := getInternalName(roleName, clinicId)
	role, err := dbutils.FindRoleByInternalName(internalName, s.dbClient, ctx)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role " + roleName + " does not exist")
	}
	roleType := role.Type
	// caller has to be allowed to create a role to be allowed to assign it to a user
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionAssign, roleType, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to assign role " + roleName)
	}
	_, err = s.enforcer.AddRoleForUser(strconv.Itoa(int(accountId)), internalName)
	return err
}

func (s *RBACService) CheckRemoveRolePermission(roleName string, clinicId int32, callerAccountId int32, ctx context.Context) (bool, error) {
	roleType, err := s.getRoleType(roleName, clinicId, ctx)
	if err != nil || roleType.String() == "" {
		return false, err
	}
	return s.checkPermissionOnRole(callerAccountId, ActionRemove, roleType, ctx)
}

func (s *RBACService) RemoveRoleFromAccount(roleName string, clinicId int32, accountId int32, callerAccountId int32, ctx context.Context) error {
	internalName := getInternalName(roleName, clinicId)
	roleType, err := s.getRoleType(internalName, clinicId, ctx)
	if err != nil || roleType.String() == "" {
		return err
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionRemove, roleType, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to remove role " + roleName)
	}
	_, err = s.enforcer.DeleteRoleForUser(strconv.Itoa(int(accountId)), internalName)
	return err
}

func (s *RBACService) CheckResourcePermission(callerAccountId int32, ctx context.Context) (bool, error) {
	return s.checkPermissionOnAction(callerAccountId, ActionEdit, ctx)
}

func (s *RBACService) CreateResource(name string, description string, callerAccountId int32, ctx context.Context) error {
	granted, err := s.checkPermissionOnResource(callerAccountId, ActionEdit, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to create resource")
	}
	return dbutils.CreateResource(name, description, s.dbClient, ctx)
}

func (s *RBACService) DeleteResourceByName(name string, callerAccountId int32, ctx context.Context) (int, error) {
	if name == ResourceRole || name == ResourceClinicRole || name == ResourceExternalRole {
		return 0, errors.New("role related resource cannot be deleted")
	}
	granted, err := s.checkPermissionOnResource(callerAccountId, ActionEdit, ctx)
	if err != nil {
		return 0, err
	}
	if !granted {
		return 0, errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to create resource")
	}
	return dbutils.DeleteResource(name, s.dbClient, ctx)
}

func (s *RBACService) DeleteResourceById(id int32, callerAccountId int32, ctx context.Context) (int, error) {
	r, err := s.GetResourceById(id, ctx)
	if err != nil {
		return 0, err
	}
	if r == nil {
		return 0, nil
	}
	name := r.Name
	if name == ResourceRole || name == ResourceClinicRole || name == ResourceExternalRole {
		return 0, errors.New("role related resource cannot be deleted")
	}
	granted, err := s.checkPermissionOnResource(callerAccountId, ActionEdit, ctx)
	if err != nil {
		return 0, err
	}
	if !granted {
		return 0, errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to create resource")
	}
	return dbutils.DeleteResourceById(id, s.dbClient, ctx)
}

func (s *RBACService) GetResources(ctx context.Context) ([]*ent.RBACResources, error) {
	return dbutils.GetAllResources(s.dbClient, ctx)
}

func (s *RBACService) GetResourceById(id int32, ctx context.Context) (*ent.RBACResources, error) {
	return dbutils.FindResourceById(id, s.dbClient, ctx)
}

func (s *RBACService) GetResourceDescription(name string, ctx context.Context) (string, error) {
	return dbutils.GetResourceDescription(name, s.dbClient, ctx)
}

func (s *RBACService) UpdateResourceDescription(name string, description string, callerAccountId int32, ctx context.Context) error {
	granted, err := s.checkPermissionOnResource(callerAccountId, ActionEdit, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to edit resource")
	}
	return dbutils.UpdateResourceDescription(name, description, s.dbClient, ctx)
}

func (s *RBACService) CheckActionPermission(callerAccountId int32, ctx context.Context) (bool, error) {
	return s.checkPermissionOnResource(callerAccountId, ActionEdit, ctx)
}

func (s *RBACService) CreateAction(name string, callerAccountId int32, ctx context.Context) error {
	granted, err := s.checkPermissionOnResource(callerAccountId, ActionEdit, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to create action")
	}
	return dbutils.CreateAction(name, s.dbClient, ctx)
}

func (s *RBACService) DeleteAction(name string, callerAccountId int32, ctx context.Context) error {
	granted, err := s.checkPermissionOnResource(callerAccountId, ActionEdit, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to create action")
	}
	err = dbutils.DeleteAction(name, s.dbClient, ctx)
	if err != nil {
		return err
	}
	_, err = s.enforcer.DeletePermission(name)
	if err != nil {
		return err
	}
	return nil
}

func (s *RBACService) GetActions(ctx context.Context) ([]*ent.RBACActions, error) {
	return dbutils.FindAllActions(s.dbClient, ctx)
}

func (s *RBACService) GetDefaultPermissions(ctx context.Context) (string, error) {
	permissions, err := s.GetRolePermissions("patients", InvalidClinicId, ctx)
	bytesArr, err := json.Marshal(permissions)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	gw, err := zstd.NewWriter(&buf)
	if err != nil {
		return "", err
	}
	_, err = gw.Write(bytesArr)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := gw.Close(); err != nil {
			common.Error(err)
		}
	}()
	compressed := buf.Bytes()
	str := hex.EncodeToString(compressed)
	return str, nil
}

func (s *RBACService) checkPermissionOnRole(callerAccountId int32, action string, roleType rbacroles.Type, ctx context.Context) (bool, error) {
	roleResourceGranted, err := s.CheckPermission(callerAccountId, action, ResourceRole, ctx)
	if err != nil {
		return false, err
	}
	switch roleType {
	case rbacroles.TypeInternal:
		return roleResourceGranted, nil
	case rbacroles.TypeClinic:
		clinicRoleResourceGranted, err := s.CheckPermission(callerAccountId, action, ResourceClinicRole, ctx)
		if err != nil {
			return false, err
		}
		return clinicRoleResourceGranted || roleResourceGranted, nil
	case rbacroles.TypeExternal:
		externalResourceGranted, err := s.CheckPermission(callerAccountId, action, ResourceExternalRole, ctx)
		if err != nil {
			return false, err
		}
		return externalResourceGranted || roleResourceGranted, nil
	default:
		return false, errors.New("unidentified role type " + roleType.String())
	}
}

func (s *RBACService) CheckAddPermissionPermission(callerAccountId int32, resourceName string, roleName string, clinicId int32, ctx context.Context) (bool, error) {
	roleType, err := s.getRoleType(roleName, clinicId, ctx)
	if err != nil || roleType.String() == "" {
		return false, err
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionEdit, roleType, ctx)
	if err != nil {
		return false, err
	}
	if !granted {
		return false, nil
	}

	if resourceName == ResourceRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return false, err
		}
		if !roleResourceGranted {
			return false, nil
		}
	} else if resourceName == ResourceExternalRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return false, err
		}
		externalRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceExternalRole, ctx)
		if err != nil {
			return false, err
		}
		if !externalRoleResourceGranted && !roleResourceGranted {
			return false, nil
		}
	} else if resourceName == ResourceClinicRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return false, err
		}
		clinicRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceClinicRole, ctx)
		if err != nil {
			return false, err
		}
		if !clinicRoleResourceGranted && !roleResourceGranted {
			return false, nil
		}
	}
	return true, nil
}

func (s *RBACService) CheckDeletePermissionPermission(callerAccountId int32, resourceName string, roleName string, clinicId int32, ctx context.Context) (bool, error) {
	roleType, err := s.getRoleType(roleName, clinicId, ctx)
	if err != nil || roleType.String() == "" {
		return false, err
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionEdit, roleType, ctx)
	if err != nil {
		return false, err
	}
	if !granted {
		return false, nil
	}

	if resourceName == ResourceRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return false, err
		}
		if !roleResourceGranted {
			return false, nil
		}
	} else if resourceName == ResourceExternalRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return false, err
		}
		externalRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceExternalRole, ctx)
		if err != nil {
			return false, err
		}
		if !externalRoleResourceGranted && !roleResourceGranted {
			return false, nil
		}
	} else if resourceName == ResourceClinicRole {
		roleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceRole, ctx)
		if err != nil {
			return false, err
		}
		clinicRoleResourceGranted, err := s.CheckPermission(callerAccountId, ActionEdit, ResourceClinicRole, ctx)
		if err != nil {
			return false, err
		}
		if !clinicRoleResourceGranted && !roleResourceGranted {
			return false, nil
		}
	}

	return true, nil
}

func (s *RBACService) GetAccountRolesInDomain(accountId int32, domainType string, domainID int32, ctx context.Context) ([]*ent.RBACRoles, error) {
	// domain roles must be public (i.e. non-clinic, so visible to all) to be shared among domains
	domain := getDomain(domainType, domainID)
	// had to use some very low level API because we are actually having two set of rules
	roleNames, err := s.enforcer.GetModel()["g"]["g2"].RM.GetRoles(strconv.Itoa(int(accountId)), domain)
	if err != nil {
		return nil, err
	}
	return dbutils.FindRolesByInternalNames(roleNames, s.dbClient, ctx)
}

func (s *RBACService) GetAccountPermissionsInDomain(accountId int32, domainType string, domainID int32, ctx context.Context) (map[string][]string, error) {
	domain := getDomain(domainType, domainID)
	// had to use some very low level API because we are actually having two set of rules
	roles, err := s.enforcer.GetModel()["g"]["g2"].RM.GetRoles(strconv.Itoa(int(accountId)), domain)
	if err != nil {
		return map[string][]string{}, err
	}
	permissions := map[string][]string{}
	resources := map[string]map[string]bool{}
	for _, role := range roles {
		permission, err := s.GetRolePermissions(role, InvalidClinicId, ctx)
		if err != nil {
			return permissions, err
		}
		for act, objs := range permission {
			for _, obj := range objs {
				if _, exist := resources[act]; !exist {
					resources[act] = make(map[string]bool)
				}
				if _, exist := resources[act][obj]; !exist {
					permissions[act] = append(permissions[act], obj)
					resources[act][obj] = true
				}
			}
		}
	}
	return permissions, nil
}

func (s *RBACService) AssignRoleToAccountInDomain(roleName string, accountId int32, domainType string, domainID int32, callerAccountId int32, ctx context.Context) error {
	domain := getDomain(domainType, domainID)
	role, err := dbutils.FindRoleByInternalName(roleName, s.dbClient, ctx)
	if err != nil {
		return err
	}
	if role == nil {
		return errors.New("role " + roleName + " does not exist")
	}
	roleType := role.Type
	// caller has to be allowed to create a role to be allowed to assign it to a user
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionAssign, roleType, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to assign role " + roleName)
	}
	_, err = s.enforcer.AddNamedGroupingPolicy("g2", strconv.Itoa(int(accountId)), roleName, domain)
	return err
}

func (s *RBACService) RemoveRoleFromAccountInDomain(roleName string, accountId int32, domainType string, domainID int32, callerAccountId int32, ctx context.Context) error {
	domain := getDomain(domainType, domainID)
	roleType, err := s.getRoleType(roleName, InvalidClinicId, ctx)
	if err != nil || roleType.String() == "" {
		return err
	}
	granted, err := s.checkPermissionOnRole(callerAccountId, ActionRemove, roleType, ctx)
	if err != nil {
		return err
	}
	if !granted {
		return errors.New("account " + strconv.Itoa(int(callerAccountId)) + " is not authorized to remove role " + roleName)
	}

	_, err = s.enforcer.RemoveNamedGroupingPolicy("g2", strconv.Itoa(int(accountId)), roleName, domain)
	return err
}

func (s *RBACService) CheckPermissionInDomain(accountId int32, action string, resource string, domainType string, domainID int32, ctx context.Context) (bool, error) {
	domain := getDomain(domainType, domainID)
	enforceContext := casbin.EnforceContext{
		RType: "r2",
		PType: "p",
		EType: "e",
		MType: "m2",
	}
	granted, err := s.enforcer.Enforce(enforceContext, strconv.Itoa(int(accountId)), action, resource, domain)
	if err != nil {
		common.Error(err)
		return false, err
	}
	return granted, err
}

func (s *RBACService) checkPermissionOnAction(callerAccountId int32, action string, ctx context.Context) (bool, error) {
	return s.CheckPermission(callerAccountId, action, ResourceRole, ctx)
}

func (s *RBACService) checkPermissionOnResource(callerAccountId int32, action string, ctx context.Context) (bool, error) {
	return s.CheckPermission(callerAccountId, action, ResourceRole, ctx)
}

// getRoleType essentially this only *really* works on non-clinic roles, with name identical with internal name
// if you know its internal name, you already knew its type
// This does not guarantee returned role type is valid
func (s *RBACService) getRoleType(internalName string, clinicId int32, ctx context.Context) (rbacroles.Type, error) {
	if clinicId != InvalidClinicId {
		return rbacroles.TypeClinic, nil
	}
	return dbutils.GetRoleTypeByInternalName(internalName, s.dbClient, ctx)
}

func getInternalName(name string, clinicId int32) string {
	if clinicId != InvalidClinicId {
		return strconv.Itoa(int(clinicId)) + "_clinic_" + name
	}
	return name
}

func getDomain(domainType string, domainID int32) string {
	domain := domainType
	if domainID != 0 {
		domain = domain + "_" + strconv.Itoa(int(domainID))
	}
	return domain
}
