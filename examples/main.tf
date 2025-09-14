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

# Example 1: Simple time adjustment
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
  description = "Adjusted morning schedule: ${data.cronmath_calculate.morning_schedule.result}"
}

# Example 2: Multiple operations
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
  description = "Complex schedule result: ${data.cronmath_calculate.complex_schedule.result}"
}

# Example 3: Resource for persistent schedule
resource "cronmath_schedule" "backup_schedule" {
  name        = "database-backup"
  base_cron   = "0 2 * * *"  # 2:00 AM
  description = "Database backup schedule with timezone adjustment"
  
  adjustments {
    type  = "sub"
    value = 5
    unit  = "hours"
  }
}

output "backup_schedule_final" {
  value       = cronmath_schedule.backup_schedule.final_cron
  description = "Backup schedule final time: ${cronmath_schedule.backup_schedule.final_cron}"
}

# Example 4: Staggered schedules
locals {
  base_sync_time = "0 3 * * *"  # 3:00 AM
}

resource "cronmath_schedule" "primary_sync" {
  name        = "primary-sync"
  base_cron   = local.base_sync_time
  description = "Primary synchronization job"
}

resource "cronmath_schedule" "secondary_sync" {
  name        = "secondary-sync"
  base_cron   = local.base_sync_time
  description = "Secondary synchronization job - 30 minutes after primary"
  
  adjustments {
    type  = "add"
    value = 30
    unit  = "minutes"
  }
}

resource "cronmath_schedule" "tertiary_sync" {
  name        = "tertiary-sync"
  base_cron   = local.base_sync_time
  description = "Tertiary synchronization job - 1 hour after primary"
  
  adjustments {
    type  = "add"
    value = 1
    unit  = "hours"
  }
}

# Example 5: Timezone adjustments
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

output "timezone_schedules" {
  value = {
    utc = var.utc_schedule
    est = data.cronmath_calculate.est_schedule.result
    pst = data.cronmath_calculate.pst_schedule.result
    jst = data.cronmath_calculate.jst_schedule.result
  }
  description = "Schedules in different timezones"
}
