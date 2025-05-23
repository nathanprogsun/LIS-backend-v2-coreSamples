#!/bin/bash

# 1. 环境配置 - 包含必要的环境变量
curl -X PUT -d '{
  "log_level": "debug",
  "order_service_name": "go.micro.lis.service.order.v2",
  "dry_run": false
}' http://localhost:8500/v1/kv/micro/config/lis/env

# 2. 端点配置 - 使用实际的服务地址
curl -X PUT -d '{
  "getTestGroup": "http://localhost:8080/api/v1/test-group",
  "getTestGroupMapping": "http://localhost:8080/api/v1/test-group-mapping",
  "accounting": "http://localhost:8080/api/v1/accounting",
  "charging": "http://localhost:8080/api/v1/charging",
  "order": "http://localhost:8080/api/v1/order"
}' http://localhost:8500/v1/kv/micro/config/lis/endpointsStaging

# 3. 密钥配置 - 使用实际的密钥
curl -X PUT -d '{
  "jwt_secret": "${JWT_SECRET}",
  "secret": "dev_secret_key_${RUN_ENV}",
  "secret_staging": "staging_secret_key_${RUN_ENV}"
}' http://localhost:8500/v1/kv/micro/config/lis/secrets

# 4. 允许的服务列表 - 实际的服务列表
curl -X PUT -d '[
  "go.micro.lis.service.order.v2",
  "go.micro.lis.service.coresamples.v2",
  "go.micro.lis.service.accounting.v2",
  "go.micro.lis.service.charging.v2"
]' http://localhost:8500/v1/kv/micro/config/lis/allowedServices

# 5. 允许的诊所列表 - 实际的诊所列表
curl -X PUT -d '[
  "clinic_001",
  "clinic_002",
  "clinic_003"
]' http://localhost:8500/v1/kv/micro/config/lis/allowedClinics

echo "Consul configuration completed!" 