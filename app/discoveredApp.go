package app

import "sync"

type discoveredApp struct {
	Address string
}

type discoveredAppList struct {
	mut   *sync.RWMutex
	items map[discoveredApp]bool
}

func newDiscoveredAppList() discoveredAppList {
	return discoveredAppList{
		mut:   &sync.RWMutex{},
		items: make(map[discoveredApp]bool),
	}
}

func (list *discoveredAppList) add(app discoveredApp) bool {
	list.mut.Lock()
	defer list.mut.Unlock()

	_, ok := list.items[app]
	if !ok {
		list.items[app] = true
		return true
	}

	return false
}

func (list *discoveredAppList) remove(app discoveredApp) {
	delete(list.items, app)
}

func (list *discoveredAppList) iter() <-chan discoveredApp {
	c := make(chan discoveredApp)

	iter := func() {
		list.mut.RLock()
		defer list.mut.RUnlock()

		for app := range list.items {
			c <- app
		}

		close(c)
	}

	go iter()

	return c
}
