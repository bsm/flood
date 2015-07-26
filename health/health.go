package health

// Objects implementing the Check can be registered
// with the env to indicate overall health status.
//
// IsHealthy is called on every bid request and must
// therefore return as quickly as possible. Avoid I/O operations
// at all costs and ensure the health check logic itself is trivial.
//
// If you want e.g. to register a health check for a remote data store
// connection, do not ping the data store on IsHealthy(). Instead, avoid
// I/O and use the ContinuousPing helper.
type Check interface {
	IsHealthy() bool
}

// CheckFunc can be registered as a HealthCheck
type CheckFunc func() bool

func (f CheckFunc) IsHealthy() bool { return f() }
