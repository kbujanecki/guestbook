package controller

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"sync"
)

// WatchedResourcePredicateProvider holds watched resource state and provides predicates.
type WatchedResourcePredicateProvider struct {
	watched map[types.NamespacedName]bool // Map of namespaced names being tracked
	mu      sync.RWMutex                  // Read-Write Mutex for safe concurrent access
}

// NewWatchedResourcePredicateProvider creates a new WatchedResourcePredicateProvider instance.
func NewWatchedResourcePredicateProvider() *WatchedResourcePredicateProvider {
	return &WatchedResourcePredicateProvider{
		watched: make(map[types.NamespacedName]bool),
	}
}

func (u *WatchedResourcePredicateProvider) ForWrite() predicate.Funcs {
	return predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			u.mu.Lock()
			defer u.mu.Unlock()
			key := types.NamespacedName{Name: e.Object.GetName(), Namespace: e.Object.GetNamespace()}
			u.watched[key] = true
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			u.mu.Lock()
			defer u.mu.Unlock()
			key := types.NamespacedName{Name: e.ObjectNew.GetName(), Namespace: e.ObjectNew.GetNamespace()}
			u.watched[key] = true
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			u.mu.Lock()
			defer u.mu.Unlock()
			key := types.NamespacedName{Name: e.Object.GetName(), Namespace: e.Object.GetNamespace()}
			delete(u.watched, key)
			return true
		},
	}
}

func (u *WatchedResourcePredicateProvider) ForRead() predicate.Funcs {
	return predicate.Funcs{
		GenericFunc: func(e event.GenericEvent) bool {
			u.mu.RLock()
			defer u.mu.RUnlock()
			key := types.NamespacedName{Name: e.Object.GetName(), Namespace: e.Object.GetNamespace()}
			_, exists := u.watched[key]
			return exists
		},
	}
}
