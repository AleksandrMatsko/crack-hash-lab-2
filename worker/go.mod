module distributed.systems.labs/worker

go 1.21

require (
	github.com/gorilla/mux v1.8.1
	github.com/spf13/viper v1.18.2
)

require (
	distributed.systems.labs/shared/pkg/alphabet v0.0.0-unpublished
	distributed.systems.labs/shared/pkg/cartesian-gen v0.0.0-unpublished
	distributed.systems.labs/shared/pkg/communication v0.0.0-unpublished
	distributed.systems.labs/shared/pkg/contracts v0.0.0-unpublished
	distributed.systems.labs/shared/pkg/middlewares v0.0.0-unpublished
)

replace distributed.systems.labs/shared/pkg/middlewares v0.0.0-unpublished => ./../shared/pkg/middlewares

replace distributed.systems.labs/shared/pkg/contracts v0.0.0-unpublished => ./../shared/pkg/contracts

replace (
	distributed.systems.labs/shared/pkg/alphabet v0.0.0-unpublished => ./../shared/pkg/alphabet
	distributed.systems.labs/shared/pkg/cartesian-gen v0.0.0-unpublished => ./../shared/pkg/cartesian-gen
	distributed.systems.labs/shared/pkg/communication v0.0.0-unpublished => ./../shared/pkg/communication
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/exp v0.0.0-20230905200255-921286631fa9 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
