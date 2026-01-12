# Load the restart_process extension
load('ext://restart_process', 'docker_build_with_restart')


### K8s Config ###
# use secrets
k8s_yaml('./infra/development/k8s/secrets.yaml')
k8s_yaml('./infra/development/k8s/app-config.yaml')
k8s_yaml('./infra/development/k8s/postgres-db.yaml')

# MongoDB
k8s_yaml('./infra/development/k8s/mongo-db.yaml')

# PG
k8s_resource('postgres-deployment', labels=['database'], port_forwards='5432:5432')

# Mongo resources (DB + admin UI)
k8s_resource('mongo', labels=['database'], port_forwards='27017:27017')
k8s_resource('mongo-express', labels=['database'], port_forwards='8082:8081')

### API GATEWAY ###
gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway/cmd/main.go'
if os.name == 'nt':
  gateway_compile_cmd = './infra/development/docker/api-gateway-build.bat'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], labels="compiles")

docker_build_with_restart(
  'yaak-kaii/api-gateway',
  '.',
  entrypoint=['/app/build/api-gateway'],
  dockerfile='./infra/development/docker/api-gateway.Dockerfile',
  only=[
    './build/api-gateway',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8081,
             resource_deps=['api-gateway-compile'], labels="services")
### End of API Gateway ###

### AUTH SERVICE ###
auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/auth-service ./services/auth-service/cmd/main.go'
if os.name == 'nt':
  auth_compile_cmd = './infra/development/docker/auth-service-build.bat'

local_resource(
  'auth-service-compile',
  auth_compile_cmd,
  deps=['./services/auth-service', './shared'], labels="compiles")

docker_build_with_restart(
  'yaak-kaii/auth-service',
  '.',
  entrypoint=['/app/build/auth-service'],
  dockerfile='./infra/development/docker/auth-service.Dockerfile',
  only=[
    './build/auth-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)

k8s_yaml('./infra/development/k8s/auth-service-deployment.yaml')
k8s_resource('auth-service', resource_deps=['auth-service-compile'], labels="services")
### End of AUTH SERVICE ###

### START OF PRODUCT SERVICE ###
auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/product-service ./services/product-service/cmd/main.go'
if os.name == 'nt':
  auth_compile_cmd = './infra/development/docker/product-service-build.bat'

local_resource(
  'product-service-compile',
  auth_compile_cmd,
  deps=['./services/product-service', './shared'], labels="compiles")

docker_build_with_restart(
  'yaak-kaii/product-service',
  '.',
  entrypoint=['/app/build/product-service'],
  dockerfile='./infra/development/docker/product-service.Dockerfile',
  only=[
    './build/product-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)
k8s_yaml('./infra/development/k8s/product-service-deployment.yaml')
k8s_resource('product-service', resource_deps=['product-service-compile'], labels="services")
### End OF PRODUCT SERVICE ###

### START OF PRODUCT CATALOG SERVICE ###
auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/product-catalog-service ./services/product-catalog-service/cmd/main.go'
if os.name == 'nt':
  auth_compile_cmd = './infra/development/docker/product-catalog-service-build.bat'

local_resource(
  'product-catalog-service-compile',
  auth_compile_cmd,
  deps=['./services/product-catalog-service', './shared'], labels="compiles")

docker_build_with_restart(
  'yaak-kaii/product-catalog-service',
  '.',
  entrypoint=['/app/build/product-catalog-service'],
  dockerfile='./infra/development/docker/product-catalog-service.Dockerfile',
  only=[
    './build/product-catalog-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)
k8s_yaml('./infra/development/k8s/product-catalog-service-deployment.yaml')
k8s_resource('product-catalog-service', resource_deps=['product-catalog-service-compile'], labels="services")
### End OF PRODUCT CATALOG SERVICE ###


### START OF ORDER SERVICE ###
auth_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/order-service ./services/order-service/cmd/main.go'
if os.name == 'nt':
  auth_compile_cmd = './infra/development/docker/order-service-build.bat'

local_resource(
  'order-service-compile',
  auth_compile_cmd,
  deps=['./services/order-service', './shared'], labels="compiles")

docker_build_with_restart(
  'yaak-kaii/order-service',
  '.',
  entrypoint=['/app/build/order-service'],
  dockerfile='./infra/development/docker/order-service.Dockerfile',
  only=[
    './build/order-service',
    './shared',
  ],
  live_update=[
    sync('./build', '/app/build'),
    sync('./shared', '/app/shared'),
  ],
)
k8s_yaml('./infra/development/k8s/order-service-deployment.yaml')
k8s_resource('order-service', resource_deps=['order-service-compile'], labels="services")
### End OF ORDER SERVICE ###
