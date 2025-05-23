package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/rbacactions"
	"errors"
)

func FindActionByName(name string, client *ent.Client, ctx context.Context) (*ent.RBACActions, error) {
	actions, err := client.RBACActions.Query().Where(rbacactions.NameEqualFold(name)).Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(actions) == 0 {
		return nil, nil
	}
	return actions[0], nil
}

func CreateAction(name string, client *ent.Client, ctx context.Context) error {
	a, err := FindActionByName(name, client, ctx)
	if err != nil {
		return err
	}
	if a != nil {
		return errors.New("action " + name + "already exists")
	}
	_, err = client.RBACActions.Create().SetName(name).Save(ctx)
	return err
}

func DeleteAction(name string, client *ent.Client, ctx context.Context) error {
	_, err := client.RBACActions.Delete().Where(rbacactions.NameEQ(name)).Exec(ctx)
	return err
}

func FindAllActions(client *ent.Client, ctx context.Context) ([]*ent.RBACActions, error) {
	return client.RBACActions.Query().All(ctx)
}
