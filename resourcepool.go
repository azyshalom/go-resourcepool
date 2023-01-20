package resourcepool

type CreateCallback func() interface{}
type CleanCallback func(i interface{})

type ResourcePool struct {
	resources      chan interface{}
	createCallback CreateCallback
	cleanCallback  CleanCallback
}

// Creates a new pool of resources.
func New(createCallback CreateCallback, cleanCallback CleanCallback, max int) *ResourcePool {

	rp := &ResourcePool{
		resources:      make(chan interface{}, max),
		createCallback: createCallback,
		cleanCallback:  cleanCallback,
	}

	return rp
}

// Get a resource from the pool.
func (rp *ResourcePool) Get() interface{} {

	// Get a resource from pool or create new
	// if the pool is empty

	var resource interface{}
	select {
	case resource = <-rp.resources:
	default:
		resource = rp.createCallback()
	}
	return resource
}

// Returns a resource to the pool.
func (rp *ResourcePool) Put(resource interface{}) {

	// Return a resource to pool or clean it
	// in case the pool is already full

	select {
	case rp.resources <- resource:
	default:
		rp.cleanCallback(resource)
	}
}
