package registry

import "sync"

type Record struct {
	IP     string
	Target string
}

type Registry struct {
	m            sync.RWMutex
	r            map[string]Record
	containerMap sync.Map // containerID -> fqdn mapping
}

func NewRegistry() *Registry {
	return &Registry{
		r: make(map[string]Record),
	}
}

func (rg *Registry) AddA(name, ip string) {
	rg.m.Lock()
	defer rg.m.Unlock()

	rg.r[name] = Record{IP: ip}
}

func (rg *Registry) AddCNAME(name, tgt string) {
	rg.m.Lock()
	defer rg.m.Unlock()

	rg.r[name] = Record{Target: tgt}
}

func (rg *Registry) Del(name string) {
	rg.m.Lock()
	defer rg.m.Unlock()

	delete(rg.r, name)
}

func (rg *Registry) RemoveA(name, ip string) {
	rg.m.Lock()
	defer rg.m.Unlock()

	if rec, exists := rg.r[name]; exists && rec.IP == ip {
		delete(rg.r, name)
	}
}

func (rg *Registry) AddContainer(containerID, fqdn string) {
	rg.containerMap.Store(containerID, fqdn)
}

func (rg *Registry) RemoveContainer(containerID string) (string, bool) {
	fqdn, exists := rg.containerMap.LoadAndDelete(containerID)
	if !exists {
		return "", false
	}
	return fqdn.(string), true
}

func (rg *Registry) Resolve(q string) (ips []string, ok bool) {
	rg.m.RLock()
	defer rg.m.RUnlock()

	rec, ok := rg.r[q]
	if !ok {
		return nil, false
	}

	if rec.IP != "" {
		return []string{rec.IP}, true
	}

	if rec.Target != "" {
		ips, ok = rg.Resolve(rec.Target)
		return ips, ok
	}

	return nil, false
}
