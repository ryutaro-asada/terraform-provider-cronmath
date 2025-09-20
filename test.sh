#!/bin/bash

set -e

echo "🚀 Terraform Provider CronMath - Quick Test"
echo "================================================"

# 1. Build
echo "📦 Building provider..."
go build -o terraform-provider-cronmath

# 2. Install
echo "📥 Installing provider..."
OS=$(go env GOOS)
ARCH=$(go env GOARCH)
INSTALL_PATH="$HOME/.terraform.d/plugins/registry.terraform.io/ryutaro-asada/cronmath/1.0.0/${OS}_${ARCH}"
mkdir -p "$INSTALL_PATH"
cp terraform-provider-cronmath "$INSTALL_PATH/"

# 3. Copy example configuration
cp examples/main.tf ./

# 4. Initialize Terraform
echo "⚙️ Initializing Terraform..."
terraform init

# 5. Apply configuration
echo "⚡ Running Terraform..."
terraform apply -auto-approve

# 6. Show results
echo ""
echo "✅ Test complete! Results:"
echo "================================"
terraform output

# 7. Cleanup
echo ""
echo "🧹 Cleaning up..."
terraform destroy -auto-approve
rm -rf .terraform .terraform.lock.hcl terraform.tfstate terraform.tfstate.backup
rm main.tf
rm terraform-provider-cronmath

echo ""
echo "🎉 All tests passed successfully!"
