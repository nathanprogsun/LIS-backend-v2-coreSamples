package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/rbacroles"
	"errors"
)

const InvalidClinicId = 0

func CreateRole(name string, internalName string, clinicId int32, roleType rbacroles.Type, client *ent.Client, ctx context.Context) error {
	role, err := FindRoleByInternalName(internalName, client, ctx)
	if err != nil {
		return err
	}
	if role != nil {
		return errors.New("role with internal name " + internalName + " already exists")
	}
	_, err = client.RBACRoles.
		Create().
		SetName(name).
		SetInternalName(internalName).
		SetType(roleType).
		SetClinicID(clinicId).
		Save(ctx)
	return err
}

func FindRoleByInternalName(internalName string, client *ent.Client, ctx context.Context) (*ent.RBACRoles, error) {
	roles, err := client.RBACRoles.
		Query().
		Where(rbacroles.InternalNameEqualFold(internalName)).
		Limit(1).
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return nil, nil
	}
	return roles[0], nil
}

func FindRolesByInternalNames(internalNames []string, client *ent.Client, ctx context.Context) ([]*ent.RBACRoles, error) {
	if internalNames == nil || len(internalNames) == 0 {
		return []*ent.RBACRoles{}, nil
	}

	roles, err := client.RBACRoles.Query().Where(
		rbacroles.InternalNameIn(internalNames...),
	).
		All(ctx)

	if err != nil {
		return nil, err
	}
	return roles, nil
}

func FindRolesByName(name string, client *ent.Client, ctx context.Context) ([]*ent.RBACRoles, error) {
	roles, err := client.RBACRoles.Query().Where(rbacroles.NameEqualFold(name)).All(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func FindRolesByNames(names []string, client *ent.Client, ctx context.Context) ([]*ent.RBACRoles, error) {
	if names == nil || len(names) == 0 {
		return []*ent.RBACRoles{}, nil
	}
	roles, err := client.RBACRoles.
		Query().
		Where(rbacroles.NameIn(names...)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func FindRolesByType(roleType rbacroles.Type, client *ent.Client, ctx context.Context) ([]*ent.RBACRoles, error) {
	roles, err := client.RBACRoles.Query().Where(rbacroles.TypeEQ(roleType)).All(ctx)
	if err != nil {
		return nil, err
	}
	return roles, err
}

// FindSharedRoles currently internal + external roles
func FindSharedRoles(client *ent.Client, ctx context.Context) ([]*ent.RBACRoles, error) {
	roles, err := client.RBACRoles.Query().Where(
		rbacroles.Or(
			rbacroles.TypeEQ(rbacroles.TypeInternal),
			rbacroles.TypeEQ(rbacroles.TypeExternal),
		)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// FindClinicRoles Find clinic typed roles within a clinic
func FindClinicRoles(clinicId int32, client *ent.Client, ctx context.Context) ([]*ent.RBACRoles, error) {
	if clinicId == InvalidClinicId {
		return nil, errors.New("invalid clinic ID")
	}
	roles, err := client.RBACRoles.
		Query().
		Where(rbacroles.And(
			rbacroles.TypeEQ(rbacroles.TypeClinic),
			rbacroles.ClinicIDEQ(clinicId),
		)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func DeleteRoleByInternalName(internalName string, client *ent.Client, ctx context.Context) (int, error) {
	txClient, err := client.Tx(ctx)
	if err != nil {
		return 0, err
	}
	num, err := txClient.RBACRoles.Delete().Where(rbacroles.InternalNameEqualFold(internalName)).Exec(ctx)
	if err != nil {
		return 0, Rollback(txClient, err)
	}
	return num, txClient.Commit()
}

func GetRoleTypeByInternalName(internalName string, client *ent.Client, ctx context.Context) (rbacroles.Type, error) {
	roles, err := client.RBACRoles.
		Query().
		Where(rbacroles.InternalNameEqualFold(internalName)).
		Limit(1).
		All(ctx)
	if err != nil {
		return "", err
	}
	if len(roles) == 0 {
		return "", nil
	}
	// assume at most one matching result
	return roles[0].Type, nil
}
