package dbutils

import (
	"coresamples/common"
	"coresamples/ent"
	"fmt"
)

func Rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		common.Errorf("Error occurred when rolling back", fmt.Errorf("%v, %v", err, rerr))
	}
	return err
}
