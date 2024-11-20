# Oasis Database Upgrade Documentation - Step 3 (Pat)

## Overview
Step 3 (Pat) of the Oasis database upgrade process focuses on pattern scripts deployment, configuration management, and final system setup. This step ensures all objects are properly compiled and the system is ready for production use.

## Process Flow
1. Initial Setup & Compilation
2. Pattern Scripts Deployment
3. Configuration Management
4. System Finalization

## Script Execution Flow

### 1. Initial Compilation Phase
Three rounds of compilation are performed to ensure all objects are
- Repeated 3 times for dependency resolution
- Uses both targeted (cmpUpg) and full schema (cmpx) compilation

### 2. Environment Setup
- **RE_Bld_Init.SQL**: Initializes build environment
- **RE_Bld_Time_Define.SQL**: Sets up timing variables
- **RE_Bld_Log_Step.SQL**: Configures logging

### 3. Core Setup Phase
1. **RE_Bld_Notice_Step_Begin.SQL**: Marks step commencement
2. **RE_Bld_Check_Aliens_Init.SQL**: Sets up audit tracking
3. **RE_Bld_Config_Prop_Off.SQL**: Disables configuration propagation
4. **RE_Bld_Oasis_Title.SQL**: Updates version information

### 4. Pattern Scripts Deployment
- **OASIS_BUILD_BASE_PAT**: Deploys base pattern scripts
- **OASIS_BUILD_CUST_PAT**: Deploys custom pattern scripts
- Environment restoration via RE_ENV_BASE

### 5. Configuration Management
- **RE_Bld_CFG_Triggers.sql**: Creates configuration triggers
  - Requires valid OASIS_CONFIG package
  - Depends on Config_Event_Util data from Pat scripts

### 6. Build Finalization
1. **RE_Bld_Build_Applied.SQL**: Records build in Build_Applied table
2. **RE_Bld_Compile.SQL**: Final compilation phase

### 7. Completion Phase
1. **RE_Bld_Notice_Step_End.SQL**: Marks step completion
2. **RE_Bld_Time_Report.SQL**: Generates execution timing report
3. **RE_Bld_Notice_Build.SQL**: Provides build completion notice
4. **RE_Bld_Cleanup.SQL**: Performs cleanup operations

## Key Features and Notes

### Session Configuration
- RECYCLEBIN is turned OFF for the session
```sql
ALTER SESSION SET RECYCLEBIN = OFF;
```

### Compilation Strategy
- Multiple compilation cycles ensure dependency resolution
- Both targeted and full schema compilation used
- Final compilation after configuration triggers creation

### Configuration Management
- Configuration propagation disabled during upgrade
- CFG triggers created after pattern scripts
- Requires valid OASIS_CONFIG package

### FM Queue Management
- FM Queue start is commented out (as of 2014.03.19)
- Prevents premature queue activation
- Waits for TIPS application

