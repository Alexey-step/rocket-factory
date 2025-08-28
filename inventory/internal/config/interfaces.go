package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
	DisabledInitMockParts() bool
}

type IamGRPCConfig interface {
	Address() string
}
