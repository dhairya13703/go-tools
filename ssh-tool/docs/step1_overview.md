# Step 1 Overview - Database Upgrade Process

## Purpose
Step 1 is the initialization and preparation phase of the database upgrade process. It ensures the environment is properly configured and prepared for the upcoming changes.

## Key Operations

### 1. Environment Setup and Validation
- Initializes build environment settings
- Validates user connection and permissions
- Confirms build ID and version information
- Verifies tablespace configurations
- Sets up logging infrastructure

### 2. Database Preparation
- Performs three compilation cycles to ensure database stability
- Configures session parameters
- Disables recyclebin
- Modifies optimizer settings
- Turns off configuration event mode

### 3. Safety Checks
- Validates build ID against previously applied builds
- Verifies connection details and instance information
- Confirms database settings and parameters
- Provides user warnings and confirmation points

### 4. System Configuration
- Stops FM queue processing if active
- Drops all configuration (CFG) triggers
  - These will be recreated in Step 3
  - Prevents conflicts during data manipulation
- Sets up index tablespace configurations

### 5. Initial Schema Changes
- Drops requested program units and objects
- Executes base alterations
- Applies custom alterations
- Prepares for subsequent steps

## Important Notes
- Multiple validation checkpoints for safety
- Built-in rollback capabilities
- Comprehensive logging of all operations
- User confirmation required at critical points

## Mappings
### 1. Environment Setup and Validation maps to:

- Section 2: "Initialization Scripts" (RE_Bld_Init.sql, -RE_Bld_Time_Define.sql, RE_Bld_Log_Step.sql)
- Section 3: "Validation Scripts" (all validation scripts)


### 2. Database Preparation maps to:

- Section 1: "Compilation Scripts" (cmpx.sql, RE_Bld_cmpUpg.sql)
"Multiple Compilation Cycles" under Important Notes


### 3. Safety Checks maps to:

- RE_Bld_Notice_User.sql
- RE_Bld_Confirm_Build_Info.sql
- RE_Bld_Validate_Build_ID.sql
- "Safety Features" sections throughout


### 4. System Configuration maps to:

- Section 4: "System Configuration Scripts"
- RE_Bld_FM_Queue_Stop.sql
- RE_Bld_CFG_Triggers_Drop.sql
- RE_Bld_Index_Tablespace_Set.sql


### 5. Initial Schema Changes maps to:

- Section 5: "Build Scripts"
