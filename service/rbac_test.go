package service

import (
	"context"
	"coresamples/common"
	"coresamples/ent/enttest"
	"coresamples/ent/rbacroles"
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/casbin/casbin/v2"
	casbin_constant "github.com/casbin/casbin/v2/constant"
	entadapter "github.com/casbin/ent-adapter"
	casbin_ent "github.com/casbin/ent-adapter/ent"
	casbin_enttest "github.com/casbin/ent-adapter/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/exp/slices"
)

const (
	allowedTestService   = "allowedTestService"
	unallowedTestService = "unallowedTestService"
)

var casbinDbClient *casbin_ent.Client

func setupRBAC(t *testing.T) (*RBACService, context.Context) {
	dataSource := "file:ent?mode=memory&_fk=1"
	dbClient := enttest.Open(t, "sqlite3", dataSource)

	err := dbClient.Schema.Create(context.Background())

	if err != nil {
		common.Fatalf("failed opening connection to MySQL", err)
	}

	//TODO: we'll come back and add redis client here after we use redis in RBAC service
	s := &RBACService{
		Service: InitService(dbClient, nil),
	}
	c := casbin_enttest.Open(t, "sqlite3", dataSource)
	adapter, err := entadapter.NewAdapterWithClient(c)
	casbinDbClient = c
	if err != nil {
		common.Fatal(err)
	}

	s.enforcer, err = casbin.NewEnforcer("rbac_model.conf", adapter)
	if err != nil {
		common.Fatal(err)
	}
	s.enforcer.SetFieldIndex("p", casbin_constant.DomainIndex, 3)
	s.enforcer.EnableAutoSave(true)
	ctx, err := setupTestDB(s)
	if err != nil {
		common.Fatal(err)
	}
	return s, ctx
}

func setupTestDB(s *RBACService) (context.Context, error) {
	file, err := os.Open("./test_input.csv")
	if err != nil {
		return nil, err
	}

	defer file.Close()
	ctx := context.Background()

	_, err = s.dbClient.RBACRoles.Create().
		SetName("IT").
		SetInternalName("IT").
		SetClinicID(InvalidClinicId).
		SetType(rbacroles.TypeInternal).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	resources := []string{"ROLE", "CLINICROLE", "EXTERNALROLE", "Password"}
	for _, rsc := range resources {
		_, err = s.dbClient.RBACResources.Create().SetName(rsc).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	// a subset of needed actions
	actions := []string{"EDIT", "VIEW", "DELETE", "CREATE", "ASSIGN", "REMOVE"}
	for _, act := range actions {
		_, err = s.dbClient.RBACActions.Create().SetName(act).Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		if _, err = s.enforcer.AddPermissionForUser(r[0], r[1], r[2]); err != nil {
			return nil, err
		}

	}
	if _, err = s.enforcer.AddRoleForUser("0", "IT"); err != nil {
		return nil, err
	}
	return ctx, nil
}

func cleanupRBAC(s *RBACService) {
	defer func() {
		if s.redisClient != nil {
			if err := s.redisClient.Close(); err != nil {
				common.Error(err)
			}
		}
		// we can't close the client in the enforcer, so let that be...
		if s.dbClient != nil {
			if err := s.dbClient.Close(); err != nil {
				common.Error(err)
			}

		}
		if casbinDbClient != nil {
			if err := casbinDbClient.Close(); err != nil {
				common.Error(err)
			}
			casbinDbClient = nil
		}
	}()
}

func TestCreateRole3(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)
	err := s.CreateRole("Admin", "superman", InvalidClinicId, 0, ctx)
	if err == nil {
		t.Fatal("Create role should not be allowed to non-existing role type")
	}
}

func TestCreateRole4(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)
	err := s.CreateRole("Admin", "internal", InvalidClinicId, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	roles, err := s.GetRolesByType("internal", ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range roles {
		if r.InternalName == "Admin" {
			return
		}
	}
	t.Fatal("Admin role not found")
}

func TestCreateRole5(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)
	err := s.CreateRole("Admin", "internal", InvalidClinicId, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.CreateRole("Admin", "internal", InvalidClinicId, 0, ctx)
	if err == nil {
		t.Fatal("duplicate roles should not be allowed")
	}
}

func TestCreateRole6(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("Admin", "clinic", InvalidClinicId, 0, ctx)
	if err == nil {
		t.Fatal("clinic type role with invalid clinicID should not be allowed")
	}
	err = s.CreateRole("Admin", "internal", 1, 0, ctx)
	if err == nil {
		t.Fatal("non-clinic type role with valid clinicID should not be allowed")
	}
}

func TestCheckCreateRolePermission(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicAdmin", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddPermissionForUser("1_clinic_ClinicAdmin", "CREATE", "CLINICROLE")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddRoleForUser("1", "1_clinic_ClinicAdmin")
	if err != nil {
		t.Fatal(err)
	}
	g, err := s.CheckCreateRolePermission("external", 1, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if g {
		t.Fatal("user 1 with only create ClinicRole permission should not be granted to create external role")
	}

	g, err = s.CheckCreateRolePermission("clinic", 1, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !g {
		t.Fatal("user 1 with create ClinicRole permission should be granted to create clinic role")
	}
}

func TestDeleteRole1(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("Admin", "internal", InvalidClinicId, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddPermissionForUser("Admin", "VIEW", "Password")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddRoleForUser("1", "Admin")
	g, err := s.CheckPermission(1, "VIEW", "Password", ctx)
	if !g {
		t.Fatal("user 1 should have view Password permission")
	}
	if err != nil {
		t.Fatal(err)
	}
	num, err := s.DeleteRole("Admin", InvalidClinicId, 0, ctx)
	if num != 1 {
		t.Fatal(fmt.Errorf("deleted %d role, should be 1", num))
	}
	g, err = s.CheckPermission(1, "VIEW", "Password", ctx)
	if g {
		t.Fatal("user 1 should not have view Password permission")
	}
}

func TestDeleteRole2(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicAdmin", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.enforcer.AddPermissionForUser("1_clinic_ClinicAdmin", "VIEW", "Password")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddRoleForUser("1", "1_clinic_ClinicAdmin")
	g, err := s.CheckPermission(1, "VIEW", "Password", ctx)
	if !g {
		t.Fatal("user 1 should have view Password permission")
	}
	if err != nil {
		t.Fatal(err)
	}

	num, err := s.DeleteRole("ClinicAdmin", 2, 0, ctx)
	if num != 0 {
		t.Fatal("should not delete ClinicAdmin from clinic 2")
	}

	num, err = s.DeleteRole("ClinicAdmin", 1, 0, ctx)
	if num != 1 {
		t.Fatalf("deleted %d role, should be 1", num)
	}
	g, err = s.CheckPermission(1, "VIEW", "Password", ctx)
	if g {
		t.Fatal("user 1 should not have view Password permission")
	}
}

func TestDeleteRole3(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicAdmin", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddPermissionForUser("1_clinic_ClinicAdmin", "DELETE", "CLINICROLE")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddRoleForUser("1", "1_clinic_ClinicAdmin")

	if _, err := s.DeleteRole("IT", InvalidClinicId, 1, ctx); err == nil {
		t.Fatal("user 1 with only delete ClinicRole permission cannot delete internal roles")
	}
}

func TestCheckPermission1(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	g, err := s.CheckPermission(0, "CREATE", "ROLE", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !g {
		t.Fatal("create Role permission for account 0 should be granted")
	}
	g, err = s.CheckPermission(0, "CREATE", "Password", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if g {
		t.Fatal("create Password permission for account 0 should not be granted")
	}
}

func TestAddPermissionToRole2(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.AddPermissionToRole("VIEW",
		"Password",
		"IT",
		InvalidClinicId,
		0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	g, err := s.CheckPermission(0, "VIEW", "Password", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !g {
		t.Fatal("user 0 should have been granted view Password permission")
	}
}

func TestAddPermissionToRole3(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.AddPermissionToRole("destroy",
		"Password",
		"IT",
		InvalidClinicId,
		0, ctx)
	if err == nil {
		t.Fatal("action destroy is not valid, thus should not be allowed")
	}
}

func TestAddPermissionToRole4(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.AddPermissionToRole("VIEW",
		"TopSecret",
		"IT",
		InvalidClinicId,
		0, ctx)
	if err == nil {
		t.Fatal("resource TopSecret is not valid, thus should not be allowed")
	}
}

func TestAddPermissionToRole5(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.AddPermissionToRole("VIEW",
		"TopSecret",
		"CEO",
		InvalidClinicId,
		0, ctx)
	if err == nil {
		t.Fatal("role CEO is not valid, thus should not be allowed")
	}
}

func TestAddPermissionToRole6(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.AddPermissionToRole("ASSIGN",
		"Password",
		"IT",
		InvalidClinicId,
		0, ctx)
	if err == nil {
		t.Fatal("assign action should only be allowed on Role related resources")
	}
	err = s.AddPermissionToRole("REMOVE",
		"Password",
		"IT",
		InvalidClinicId,
		0, ctx)
	if err == nil {
		t.Fatal("remove action should only be allowed on Role related resources")
	}
}

func TestAddPermissionToRole7(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicAdmin", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.CreateRole("ClinicUser", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	// even though ClinicAdmin can edit ClinicRole such as ClinicUser, it cannot assign creat Role permission to it
	_, err = s.enforcer.AddPermissionForUser("1_clinic_ClinicAdmin", "EDIT", "CLINICROLE")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddRoleForUser("1", "1_clinic_ClinicAdmin")
	if err = s.AddPermissionToRole("CREATE", "ROLE", "ClinicUser", 1, 1, ctx); err == nil {
		t.Fatal("user without edit Role permission cannot add Role related permission")
	}
	fmt.Println(err)
	if err = s.AddPermissionToRole("CREATE", "CLINICROLE", "ClinicUser", 1, 1, ctx); err != nil {
		t.Fatal(err)
	}
}

func TestAddPermissionToRole9(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicAdmin", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	// even though ClinicAdmin can edit ClinicRole such as ClinicUser, it cannot assign creat Role permission to it
	_, err = s.enforcer.AddPermissionForUser("1_clinic_ClinicAdmin", "EDIT", "CLINICROLE")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddRoleForUser("1", "1_clinic_ClinicAdmin")
	if err = s.AddPermissionToRole("CREATE", "ROLE", "IT", InvalidClinicId, 1, ctx); err == nil {
		t.Fatal("ClinicAdmin is not granted to edit IT role, which is internal")
	}
	fmt.Println(err)
}

func TestAddPermissionToRole8(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicAdmin", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.CreateRole("ClinicUser", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.enforcer.AddPermissionForUser("1_clinic_ClinicAdmin", "EDIT", "CLINICROLE")
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddRoleForUser("1", "1_clinic_ClinicAdmin")
	if err = s.AddPermissionToRole("CREATE", "Password", "ClinicUser", 1, 1, ctx); err != nil {
		t.Fatal(err)
	}
	if _, err = s.enforcer.AddRoleForUser("2", "1_clinic_ClinicUser"); err != nil {
		t.Fatal(err)
	}
	g, err := s.CheckPermission(2, "CREATE", "Password", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !g {
		t.Fatal("ClinicUser should be able to create password")
	}
}

func TestDeletePermissionFromRole2(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.DeletePermissionFromRole("CREATE", "ROLE", "IT", InvalidClinicId, 0, ctx)
	if err != nil {
		t.Fatal("Role related permission cannot be deleted once assigned")
	}
}

func TestDeletePermissionFromRole3(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.AddPermissionToRole("VIEW", "Password", "IT", InvalidClinicId, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	g, err := s.CheckPermission(0, "VIEW", "Password", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !g {
		t.Fatal("user 0 should have view Password permission")
	}

	err = s.DeletePermissionFromRole("VIEW", "Password", "IT", InvalidClinicId, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	g, err = s.CheckPermission(0, "VIEW", "Password", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if g {
		t.Fatal("user 0 should not have view Password permission")
	}
}

func TestDeletePermissionFromRole4(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicAdmin", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddPermissionToRole("EDIT", "CLINICROLE", "ClinicAdmin", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AddPermissionToRole("VIEW", "EXTERNALROLE", "ClinicAdmin", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AddPermissionToRole("VIEW", "Password", "ClinicAdmin", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	_, err = s.enforcer.AddRoleForUser("1", "1_clinic_ClinicAdmin")
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddPermissionToRole("VIEW", "Password", "IT", InvalidClinicId, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.DeletePermissionFromRole("VIEW", "Password", "IT", InvalidClinicId, 1, ctx)
	if err == nil {
		t.Fatal("ClinicAdmin shouldn't be authorized to delete permission from IT")
	}

	err = s.DeletePermissionFromRole("VIEW", "EXTERNALROLE", "ClinicAdmin", 1, 1, ctx)
	if err == nil {
		t.Fatal("ClinicAdmin shouldn't be authorized to delete ExternalRole related permission from ClinicAdmin")
	}

	err = s.DeletePermissionFromRole("VIEW", "Password", "ClinicAdmin", 1, 1, ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetRolesAndPermissionsIntegration(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicUser", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AddPermissionToRole("VIEW", "Password", "ClinicUser", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AssignRoleToAccount("ClinicUser", 1, 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	roleNames, err := s.GetAccountRoleInternalNames(0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(roleNames, "1_clinic_ClinicUser") || !slices.Contains(roleNames, "IT") {
		t.Fatal("account 0 should have roles 1_clinic_ClinicUser, IT")
	}

	roles, err := s.GetAccountRoles(0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 2 {
		t.Fatalf("account 0 should have 2 roles, getting %d role", len(roles))
	}

	permissions, err := s.GetAccountPermissions(0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(permissions["CREATE"], "ROLE") ||
		!slices.Contains(permissions["EDIT"], "ROLE") ||
		!slices.Contains(permissions["DELETE"], "ROLE") ||
		!slices.Contains(permissions["ASSIGN"], "ROLE") ||
		!slices.Contains(permissions["REMOVE"], "ROLE") ||
		!slices.Contains(permissions["VIEW"], "Password") {
		fmt.Println(permissions)
		t.Fatal("permissions do not match")
	}

	clinicUserPermissions, err := s.GetRolePermissions("ClinicUser", 1, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(clinicUserPermissions["VIEW"]) != 1 || clinicUserPermissions["VIEW"][0] != "Password" {
		fmt.Println(clinicUserPermissions)
		t.Fatal("permissions do not match")
	}

	nonPrivateRoles, err := s.GetNonPrivateRoles(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(nonPrivateRoles) != 1 || nonPrivateRoles[0].Name != "IT" {
		t.Fatal("nonprivate roles do not match")
	}

	rolesByType, err := s.GetRolesByType("clinic", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rolesByType) != 1 || rolesByType[0].InternalName != "1_clinic_ClinicUser" {
		t.Fatal("roles by type do not match")
	}
}

func TestAssignRole2(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	// user 0: IT, user 1: ClinicAdmin for clinic 1
	err := s.CreateRole("ClinicAdmin", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AddPermissionToRole("ASSIGN", "CLINICROLE", "ClinicAdmin", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AddPermissionToRole("CREATE", "CLINICROLE", "ClinicAdmin", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AssignRoleToAccount("ClinicAdmin", 1, 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.CreateRole("ClinicUser", "clinic", 1, 1, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AssignRoleToAccount("ClinicUser", 1, 2, 1, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AssignRoleToAccount("IT", InvalidClinicId, 2, 1, ctx)
	if err == nil {
		t.Fatal("ClinicAdmin is not authorized to assign internal role")
	}
}

func TestAddPermissionToUser(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("ClinicUser", "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AssignRoleToAccount("ClinicUser", 1, 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	roleNames, err := s.GetAccountRoleInternalNames(0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(roleNames, "1_clinic_ClinicUser") || !slices.Contains(roleNames, "IT") {
		t.Fatal("account 0 should have roles 1_clinic_ClinicUser, IT")
	}

	roles, err := s.GetAccountRoles(0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(roles) != 2 {
		t.Fatalf("account 0 should have 2 roles, getting %d role", len(roles))
	}

	err = s.AddPermissionToUser("VIEW", "Password", 0, 1, ctx)
	if err != nil {
		t.Fatal(err)
	}

	permissions, err := s.GetAccountPermissions(0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Contains(permissions["CREATE"], "ROLE") ||
		!slices.Contains(permissions["EDIT"], "ROLE") ||
		!slices.Contains(permissions["DELETE"], "ROLE") ||
		!slices.Contains(permissions["ASSIGN"], "ROLE") ||
		!slices.Contains(permissions["REMOVE"], "ROLE") ||
		!slices.Contains(permissions["VIEW"], "Password") {
		fmt.Println(permissions)
		t.Fatal("permissions do not match")
	}

	clinicUserPermissions, err := s.GetRolePermissions("ClinicUser", 1, ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(clinicUserPermissions["VIEW"]) != 0 {
		fmt.Println(clinicUserPermissions)
		t.Fatal("permissions do not match")
	}

	nonPrivateRoles, err := s.GetNonPrivateRoles(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if len(nonPrivateRoles) != 1 || nonPrivateRoles[0].Name != "IT" {
		t.Fatal("nonprivate roles do not match")
	}

	rolesByType, err := s.GetRolesByType("clinic", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(rolesByType) != 1 || rolesByType[0].InternalName != "1_clinic_ClinicUser" {
		t.Fatal("roles by type do not match")
	}

	g, err := s.CheckPermission(0, "VIEW", "Password", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !g {
		t.Fatal("user 0 should have been granted view Password permission")
	}
}

func TestDeletePermissionFromUser(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)
	err := s.CreateRole("User", "external", 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AddPermissionToRole("Edit", "Password", "User", 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	s.AssignRoleToAccount("User", 0, 1, 0, ctx)
	s.AddPermissionToUser("View", "Password", 1, 0, ctx)

	g, err := s.CheckPermission(1, "View", "Password", ctx)
	if !g {
		t.Fatal("user 1 should have view Password permission")
	}

	g, err = s.CheckPermission(1, "Edit", "Password", ctx)
	if !g {
		t.Fatal("user 1 should have Edit Password permission")
	}

	s.DeletePermissionFromUser("View", "Password", 1, 0, ctx)
	g, err = s.CheckPermission(1, "View", "Password", ctx)
	if g {
		t.Fatal("user 1 shouldn't have view Password permission")
	}

	err = s.DeletePermissionFromUser("Edit", "Password", 1, 0, ctx)
	g, err = s.CheckPermission(1, "Edit", "Password", ctx)
	if !g {
		t.Fatal("user 1 should still have Edit Password permission")
	}

	p, err := s.GetRolePermissions("User", 0, ctx)
	if len(p["Edit"]) == 0 || p["Edit"][0] != "Password" {
		t.Fatal("role permission shouldn't change")
	}
}

func TestCreateResource2(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateResource("Password", "", 0, ctx)
	if err == nil {
		t.Fatal("Duplicate resource should not be allowed")
	}
}

func TestCreateResource3(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateResource("Password", "", 1, ctx)
	if err == nil {
		t.Fatal("Unauthorized account should not be allowed")
	}
}

func TestCreateResource4(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateResource("Password2", "", 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateAction1(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateAction("EDIT", 0, ctx)
	if err == nil {
		t.Fatal("duplicate action should not be allowed")
	}
}

func TestCreateAction3(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateAction("edit2", 1, ctx)
	if err == nil {
		t.Fatal("unauthorized user should not be allowed")
	}
}

func TestCreateAction4(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateAction("edit2", 0, ctx)
	if err != nil {
		t.Fatal(nil)
	}
}

func TestDomainRoles(t *testing.T) {
	s, ctx := setupRBAC(t)
	defer cleanupRBAC(s)

	err := s.CreateRole("User", "external", 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AddPermissionToRole("Edit", "Password", "User", 0, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}
	err = s.AssignRoleToAccount("User", 0, 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	err = s.AssignRoleToAccountInDomain("User", 2, "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	roles, err := s.GetAccountRolesInDomain(2, "clinic", 1, ctx)
	if roles == nil || err != nil || len(roles) == 0 || roles[0].Name != "User" {
		t.Fatal("unmatched roles found in domain clinic_1")
	}

	persPublic, err := s.GetAccountPermissions(2, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(persPublic) != 0 {
		t.Fatal("account 2 shouldn't have public roles")
	}

	persInDomain, err := s.GetAccountPermissionsInDomain(2, "clinic", 1, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(persInDomain) != 1 || !slices.Contains(persInDomain["Edit"], "Password") {
		t.Fatal("account 2 should have Edit password permission in clinic1")
	}

	grantedPublic1, err := s.CheckPermission(1, "Edit", "Password", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !grantedPublic1 {
		t.Fatal("user 1 should be able to edit password globally")
	}

	grantedPublic2, err := s.CheckPermission(2, "Edit", "Password", ctx)
	if err != nil {
		t.Fatal(err)
	}
	if grantedPublic2 {
		t.Fatal("user 2 should not be able to edit password globally")
	}

	grantedClinic1, err := s.CheckPermissionInDomain(2, "Edit", "Password", "clinic", 1, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if !grantedClinic1 {
		t.Fatal("user 2 should be able to edit password in clinic 1")
	}

	grantedClinic2, err := s.CheckPermissionInDomain(2, "Edit", "Password", "clinic", 2, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if grantedClinic2 {
		t.Fatal("user 2 should not be able to edit password in clinic 2")
	}

	err = s.RemoveRoleFromAccountInDomain("User", 2, "clinic", 1, 0, ctx)
	if err != nil {
		t.Fatal(err)
	}

	grantedClinic1, err = s.CheckPermissionInDomain(2, "Edit", "Password", "clinic", 1, ctx)
	if err != nil {
		t.Fatal(err)
	}
	if grantedClinic1 {
		t.Fatal("user 2 should not be able to edit password in clinic 1 after removing the role")
	}
}
