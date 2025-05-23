package common

import "github.com/go-redsync/redsync/v4"

const (
	CreateOrderLock = "lis::core_service_v2::order::create_order_lock"
)

// RedLock The wrapper that enables easier locking and so that we can test the function without
// actually use a lock that connects to redis
func RedLock(lockName string, rs *redsync.Redsync) (*redsync.Mutex, error) {
	if rs == nil {
		return nil, nil
	}
	lock := rs.NewMutex(lockName)
	if err := lock.Lock(); err != nil {
		return lock, err
	}
	return lock, nil
}

// RedUnlock A wrapper so that we can test the function without actually use a lock that connects to redis
func RedUnlock(lock *redsync.Mutex) (bool, error) {
	if lock == nil {
		return true, nil
	}
	return lock.Unlock()
}
