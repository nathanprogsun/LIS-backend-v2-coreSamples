package dbutils

import (
	"context"
	"coresamples/ent"
	"coresamples/ent/rbacresources"
	"errors"
)

func FindResourceByName(name string, client *ent.Client, ctx context.Context) (*ent.RBACResources, error) {
	resources, err := client.RBACResources.Query().Where(rbacresources.NameEqualFold(name)).Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		return nil, nil
	}
	return resources[0], nil
}

func FindResourceById(id int32, client *ent.Client, ctx context.Context) (*ent.RBACResources, error) {
	resources, err := client.RBACResources.Query().Where(rbacresources.IDEQ(int(id))).Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}

	if len(resources) == 0 {
		return nil, nil
	}

	return resources[0], nil
}

func CreateResource(name string, description string, client *ent.Client, ctx context.Context) error {
	r, err := FindResourceByName(name, client, ctx)
	if err != nil {
		return nil
	}
	if r != nil {
		return errors.New("resource " + name + " already exists")
	}

	if description == "" {
		_, err = client.RBACResources.Create().SetName(name).Save(ctx)
	} else {
		_, err = client.RBACResources.Create().SetName(name).SetDescription(description).Save(ctx)
	}

	return err
}

func GetAllResources(client *ent.Client, ctx context.Context) ([]*ent.RBACResources, error) {
	resources, err := client.RBACResources.Query().All(ctx)
	if err != nil {
		return nil, err
	}
	return resources, nil
}

func DeleteResource(name string, client *ent.Client, ctx context.Context) (int, error) {
	num, err := client.RBACResources.Delete().Where(rbacresources.NameEqualFold(name)).Exec(ctx)
	return num, err
}

func DeleteResourceById(id int32, client *ent.Client, ctx context.Context) (int, error) {
	num, err := client.RBACResources.Delete().Where(rbacresources.IDEQ(int(id))).Exec(ctx)
	return num, err
}

func GetResourceDescription(name string, client *ent.Client, ctx context.Context) (string, error) {
	resource, err := FindResourceByName(name, client, ctx)
	if err != nil {
		return "", nil
	}
	if resource == nil {
		return "", errors.New("resource " + name + " not found")
	}

	return resource.Description, nil
}

func UpdateResourceDescription(name string, description string, client *ent.Client, ctx context.Context) error {
	r, err := FindResourceByName(name, client, ctx)
	if err != nil {
		return err
	}
	if r == nil {
		return errors.New("resource " + name + " does not exist")
	}
	_, err = client.RBACResources.UpdateOne(r).SetDescription(description).Save(ctx)
	return err
}
