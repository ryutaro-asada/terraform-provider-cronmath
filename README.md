# Terraform Provider CronMath

[![Release](https://img.shields.io/github/release/ryutaro-asada/terraform-provider-cronmath.svg)](https://github.com/ryutaro-asada/terraform-provider-cronmath/releases/latest)
[![Registry](https://img.shields.io/badge/terraform-registry-623CE4)](https://registry.terraform.io/providers/ryutaro-asada/cronmath)
[![Go Report Card](https://goreportcard.com/badge/github.com/ryutaro-asada/terraform-provider-cronmath)](https://goreportcard.com/report/github.com/ryutaro-asada/terraform-provider-cronmath)
[![Tests](https://github.com/ryutaro-asada/terraform-provider-cronmath/actions/workflows/test.yml/badge.svg)](https://github.com/ryutaro-asada/terraform-provider-cronmath/actions/workflows/test.yml)
[![License: MPL-2.0](https://img.shields.io/badge/License-MPL--2.0-blue.svg)](https://opensource.org/licenses/MPL-2.0)

The CronMath provider allows you to perform time arithmetic operations on cron expressions within your Terraform configurations. It uses the [cronmath](https://github.com/ryutaro-asada/cronmath) library for cron expression manipulation.

## Features

- ➕ Add minutes or hours to cron expressions
- ➖ Subtract minutes or hours from cron expressions
- 🔄 Handle day boundary transitions automatically
- 🌍 Perfect for timezone adjustments
- ⏰ Create staggered schedules easily

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)

## Installation

### From Terraform Registry

```hcl
terraform {
  required_providers {
    cronmath = {
      source  = "ryutaro-asada/cronmath"
      version = "~> 1.0"
    }
  }
}

provider "cronmath" {
  # No configuration required
}
```

### Manual Installation

1. Download the latest release from the [releases page](https://github.com/ryutaro-asada/terraform-provider-cronmath/releases)
2. Extract the archive
3. Move the binary to `~/.terraform.d/plugins/registry.terraform.io/ryutaro-asada/cronmath/1.0.0/[OS]_[ARCH]/`

## Usage

### Data Source: cronmath_calculate

Calculate a new cron expression by applying time operations:

```hcl
data "cronmath_calculate" "morning_job" {
  input = "30 9 * * *"  # 9:30 AM
  
  operations {
    type  = "sub"
    value = 30
    unit  = "minutes"
  }
}

output "adjusted_schedule" {
  value = data.cronmath_calculate.morning_job.result  # "0 9 * * *"
}
```

## Examples

### Timezone Adjustment

```hcl
# Convert UTC to EST (UTC-5)
data "cronmath_calculate" "est_schedule" {
  input = "0 15 * * *"  # 3:00 PM UTC
  
  operations {
    type  = "sub"
    value = 5
    unit  = "hours"
  }
}

# Convert UTC to JST (UTC+9)
data "cronmath_calculate" "jst_schedule" {
  input = "0 10 * * *"  # 10:00 AM UTC
  
  operations {
    type  = "add"
    value = 9
    unit  = "hours"
  }
}
```

### Building the Provider

```bash
# Clone the repository
git clone https://github.com/ryutaro-asada/terraform-provider-cronmath.git
cd terraform-provider-cronmath

# Download dependencies
go mod download

# Build the provider
make build

# Install locally for testing
make install
```

### Running Tests

```bash
# Run unit tests
make test

# Run acceptance tests
make testacc

# Run specific test
go test -v -run TestCronMathDataSource ./internal/provider/
```


## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For bugs and feature requests, please [open an issue](https://github.com/ryutaro-asada/terraform-provider-cronmath/issues/new).

## Related Projects

- [cronmath](https://github.com/ryutaro-asada/cronmath) - The underlying library for cron expression manipulation
- [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) - Framework used to build this provider

## Acknowledgments

- Built with [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework)
- Cron expression manipulation powered by [cronmath](https://github.com/ryutaro-asada/cronmath)
