terraform {
  required_providers {
    cronmath = {
      source  = "registry.terraform.io/ryutaro-asada/cronmath"
      version = "1.0.0"
    }
  }
}

provider "cronmath" {
  # No configuration required
}

# ==================================================
# Example 1: Simple time adjustment
# ==================================================
data "cronmath_calculate" "morning_schedule" {
  input = "5 9 * * *"  # 9:05 AM
  
  operations {
    type  = "sub"
    value = 5
    unit  = "minutes"
  }
}

output "morning_schedule_result" {
  value       = data.cronmath_calculate.morning_schedule.result
  description = "Adjusted morning schedule"
}

# ==================================================
# Example 2: Multiple operations
# ==================================================
data "cronmath_calculate" "complex_schedule" {
  input = "30 10 * * *"  # 10:30 AM
  
  operations {
    type  = "add"
    value = 2
    unit  = "hours"
  }
  
  operations {
    type  = "sub"
    value = 15
    unit  = "minutes"
  }
}

output "complex_schedule_result" {
  value       = data.cronmath_calculate.complex_schedule.result
  description = "Complex schedule result"
}

# ==================================================
# Example 5: Timezone adjustments
# ==================================================
variable "utc_schedule" {
  default     = "0 15 * * *"  # 3:00 PM UTC
  description = "Base schedule in UTC"
}

data "cronmath_calculate" "est_schedule" {
  input = var.utc_schedule
  
  operations {
    type  = "sub"
    value = 5
    unit  = "hours"  # EST is UTC-5
  }
}

data "cronmath_calculate" "pst_schedule" {
  input = var.utc_schedule
  
  operations {
    type  = "sub"
    value = 8
    unit  = "hours"  # PST is UTC-8
  }
}

data "cronmath_calculate" "jst_schedule" {
  input = var.utc_schedule
  
  operations {
    type  = "add"
    value = 9
    unit  = "hours"  # JST is UTC+9
  }
}

# ==================================================
# Outputs
# ==================================================

output "simple_calculations" {
  value = {
    morning = {
      input  = data.cronmath_calculate.morning_schedule.input
      result = data.cronmath_calculate.morning_schedule.result
    }
    complex = {
      input  = data.cronmath_calculate.complex_schedule.input
      result = data.cronmath_calculate.complex_schedule.result
    }
  }
  description = "Simple calculation results"
}

output "timezone_schedules" {
  value = {
    utc = var.utc_schedule
    est = data.cronmath_calculate.est_schedule.result
    pst = data.cronmath_calculate.pst_schedule.result
    jst = data.cronmath_calculate.jst_schedule.result
  }
  description = "Schedules in different timezones"
}

# output "summary" {
#   value = format(
#     "\n========================================\n%s\n========================================\n%s\n%s\n%s\n%s\n\n%s\n%s\n%s\n%s\n========================================",
#     "CronMath Provider Test Results",
#     "Data Source Tests:",
#     "  Morning: ${data.cronmath_calculate.morning_schedule.input} → ${data.cronmath_calculate.morning_schedule.result}",
#     "  Complex: ${data.cronmath_calculate.complex_schedule.input} → ${data.cronmath_calculate.complex_schedule.result}",
#     "  Timezones: UTC(${var.utc_schedule}) → JST(${data.cronmath_calculate.jst_schedule.result})",
#   )
#   description = "Test execution summary"
# }
