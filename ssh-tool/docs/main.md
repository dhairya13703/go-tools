# Database Upgrade Scripts Documentation

## Overview
This documentation covers the database upgrade process implemented through a series of SQL scripts. The process follows a structured approach with multiple steps and configuration options.

## Main Script Structure

### 1. Build_Master.SQL
**Purpose:** Main orchestrator script that controls the entire build process
**Location:** N:\Groups\RelEng\2023.1.0\release\2023.1.0
**Key Functions:**
- Initializes build configuration
- Configures inter-step processing
- Executes the build process in three sequential steps

**Script Flow:**
1. Loads Build_Config.SQL
2. Configures inter-step processing (INTER_STEP_SCRIPT)
3. Executes in sequence:
   - Build_Master_Step1.SQL
   - Build_Master_Step2.SQL
   - Build_Master_Step3.SQL

### 2. Build_Config.SQL
**Purpose:** Central configuration script that defines all build variables and paths
**Key Configurations:**
- Build ID: 2023.1.0
- Operating System specific configurations
- Build root folder and subfolders
- Compilation modes
- Verification settings

**Important Variables:**
```sql
OASIS_BUILD_ROOT="N:\Groups\RelEng\2023.1.0\release\2023.1.0\"
OASIS_BUILD_ID='2023.1.0'
```

**Mode Settings:**
- Invoker Rights Verification: ON
- Alien Objects Verification: ON
- Compile Mode: CUSTOM
- Index Tablespace Auto Mode: ON

### 3. Step-wise Build Scripts
Each step script (Step1, Step2, Step3) follows a similar pattern:
- Loads configuration from Build_Config.SQL
- Executes specific build step using RelEng (RE) scripts

#### Build_Master_Step1.SQL
- Loads configuration
- Executes `&RE_BUILD_MASTER_STEP1`

#### Build_Master_Step2.SQL
- Loads configuration
- Executes `&RE_BUILD_MASTER_STEP2`

#### Build_Master_Step3.SQL
- Loads configuration
- Executes `&RE_BUILD_MASTER_STEP3`

## Script Dependencies
The scripts reference several other scripts that need to be documented:

### Base Scripts:
- Build_Base_Alt.SQL
- Build_Base_Typ.SQL
- Build_Base_Pls.SQL
- Build_Base_Trg.SQL
- Build_Base_Vew.SQL
- Build_Base_Pat.SQL

### Custom Scripts:
- Build_Cust_Alt.SQL
- Build_Cust_Typ.SQL
- Build_Cust_Pls.SQL
- Build_Cust_Trg.SQL
- Build_Cust_Vew.SQL
- Build_Cust_Pat.SQL

### Other Referenced Scripts:
- Build_Drop_Obj.SQL
- RE_Bld_Config.SQL

## Next Steps
To complete this documentation, we need information about:
1. The specific operations performed in each RE_BUILD_MASTER_STEP
2. The contents and purpose of each Base and Custom script
3. The relationship between these scripts and the Python dbupgrade script
4. Any error handling or rollback procedures
5. Typical execution timeframes
6. Common issues and troubleshooting steps

Would you like to provide information about any of these aspects to expand the documentation further?