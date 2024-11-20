# Database Upgrade Step 1 - Detailed Documentation

## 1. Compilation Scripts

### cmpx.sql
**Purpose:** Executes schema-wide compilation of database objects.
**Key Operations:**
- Spools compilation output to `uerror.lis`
- Uses `dbms_utility.compile_schema` to compile all objects for the current user
- Sets echo on for detailed logging
- Reloads environment settings from rcenv.sql

### RE_Bld_cmpUpg.sql
**Purpose:** Generates and executes targeted recompilation of invalid objects.
**Key Operations:**
- Creates a dynamic compilation script
- Specifically targets invalid objects of types:
  - PACKAGE BODY
  - PACKAGE
  - FUNCTION
  - MATERIALIZED VIEW
  - PROCEDURE
  - TRIGGER
  - VIEW
  - TYPE
  - SYNONYM
  - TYPE BODY
  - JAVA CLASS
- Generates individual ALTER statements for each invalid object
- Orders compilation by object_type and object_name

## 2. Initialization Scripts

### RE_Bld_Init.sql
**Purpose:** Initializes the build environment and sets up OS-specific commands.
**Key Configurations:**
- Preserves original SQL*Plus settings
- Defines OS-specific commands:
  ```sql
  OS_COPY     ='copy'
  OS_MOVE     ='move'
  OS_DELETE   ='del /F /Q'
  OS_TYPE     ='type'
  OS_READONLY ='attrib +R /S'
  ```
- Manages temporary files for build steps
- Defines SQL*Plus compatibility settings
- Determines spool append support based on SQL*Plus version

### RE_Bld_Time_Define.sql
**Purpose:** Manages timestamp creation for build process tracking.
**Key Operations:**
- Creates formatted timestamps for logging
- Uses format: 'mm/dd/yyyy hh24:mi:ss'
- Maintains environment settings

### RE_Bld_Log_Step.sql
**Purpose:** Sets up logging infrastructure for build steps.
**Key Features:**
- Creates standardized log file names:
  - Format: `Step#-yyyy.mm.dd-hh.mi.ss.log`
  - Example: `Step1-2008.03.27-11.39.58.log`
- Handles log file archiving
- Manages log append modes based on SQL*Plus version
- Includes OS user information in log names

## Important Notes

### Multiple Compilation Cycles
The step performs three compilation cycles because:
1. Initial compilation handles straightforward dependencies
2. Second pass resolves secondary dependencies
3. Final pass ensures complete resolution of complex dependencies

### Environment Management
- Original SQL*Plus settings are preserved
- Environment settings are consistently reloaded after operations
- OS-specific commands are abstracted for portability

### Logging Infrastructure
- Comprehensive logging system with:
  - Unique timestamped files
  - User identification
  - Version-specific handling
  - Archive management

## 3. Validation Scripts

### RE_Bld_Notice_User.sql
**Purpose:** Provides initial user notification and safety checks.
**Key Features:**
- Displays important notices about Release Notes and Installation Instructions
- Provides instructions for safe termination using Ctrl-C
- Includes cleanup procedures after termination
- References:
  - RE_BUILD_ABORT_CLEANUP for cleanup
  - RE_ENV_INIT for environment reset
  - RE_ENV_BASE for predefined build state

### RE_Bld_Confirm_Build_Info.sql
**Purpose:** Validates and displays build configuration for user confirmation.
**Key Validations:**
- Build ID verification
- Compilation mode check
- Invoker rights validation
- Alien objects verification
- Index tablespace settings
- Connection details including:
  - Username
  - Server host
  - Instance name
**Features:**
- Cross-version compatibility (Oracle 9i support)
- Detailed parameter display with explanations
- Safety prompts for incorrect settings

### RE_Bld_Validate_Build_ID.sql
**Purpose:** Ensures build ID integrity and manages build tracking.
**Key Operations:**
- Manages BUILD_APPLIED table:
  - Creates if non-existent
  - Verifies structure
  - Validates build ID uniqueness
- Handles temporary utility tables
- Provides multiple validation attempts
- Includes error handling and user prompts
**Technical Implementation:**
- Uses global temporary tables for validation
- Implements proper cleanup procedures
- Handles multiple validation scenarios

### RE_Bld_Index_Tablespace_Set.sql
**Purpose:** Manages and validates index tablespace configuration.
**Key Features:**
- Supports both automatic and manual tablespace setting modes
- Multiple validation attempts for tablespace names
- Enforces proper tablespace existence
- Provides user interaction for corrections
**Safety Features:**
- Two-step validation process
- Forced exit on invalid settings
- Clear user warnings and instructions
- Proper cleanup procedures

## 4. System Configuration Scripts

### RE_Bld_FM_Queue_Stop.sql
**Purpose:** Manages the stopping of FM Queue processing during the build process.
**Key Operations:**
- Validates and stops FM Policy Interface Queue
- Comprehensive error handling and status checking
- Queue management through DBMS_AQ

**Validation Steps:**
1. Queue Setup Verification
   - Checks DBA_QUEUE_TABLES for queue table existence
   - Validates schema ownership

2. Queue Status Checks
   - Verifies dequeue_enabled status
   - Confirms queue accessibility

3. Job Management
   - Checks for existing submitted jobs
   - Handles queue message processing

**Error Handling:**
- Detailed error messages for each failure point
- Specific handling for:
  - Environment not set up
  - Queue status issues
  - Missing queues
  - Non-running jobs

### RE_Bld_CFG_Triggers_Drop.sql
**Purpose:** Manages the removal of all CFG triggers before the build process.
**Key Features:**
- Identifies and drops all triggers ending with '_CFG'
- Provides detailed execution logging
- Includes failure tracking and reporting

**Process Flow:**
1. Initial Count
   - Counts existing CFG triggers
   - Displays summary header

2. Trigger Removal
   - Systematically drops each trigger
   - Tracks successful and failed operations

3. Final Reporting
   - Reports total triggers dropped
   - Indicates any failed operations

**Technical Implementation:**
- Uses dynamic SQL for trigger removal
- Exception handling for each trigger
- Maintains environment settings

## 5. Build Scripts

### OASIS_BUILD_DROP_OBJ (Build_Drop_Obj.SQL)
**Purpose:** Manages the systematic removal of obsolete or unnecessary database objects.
**Key Features:**
- Handles multiple object types:
  - FUNCTIONS
  - PACKAGES
  - PROCEDURES
  - TRIGGERS
  - VIEWS
  - MATERIALIZED VIEWS
- Provides detailed execution logging
- Includes comprehensive error handling

**Process Flow:**
1. Object Collection
   - Maintains list of objects to be dropped
   - Supports comments and documentation for each object
   - Enforces naming and format rules

2. Execution Management
   - Tracks statistics (total, processed, failed)
   - Provides detailed progress reporting
   - Handles exceptions gracefully

3. Reporting
   - Shows successful drops
   - Reports failed operations
   - Provides summary statistics

### OASIS_BUILD_BASE_ALT (Build_Base_Alt.SQL)
**Purpose:** Executes base alterations scripts in chronological order.
**Key Features:**
- Runs ALT scripts from specified base directory
- Maintains version control through dated files
- Handles database structural changes

**Execution Flow:**
- Processes scripts in chronological order
- Includes various types of alterations:
  - DDL changes
  - Sequence modifications
  - Table structure updates
- Maintains consistent environment settings

### OASIS_BUILD_CUST_ALT (Build_Cust_Alt.SQL)
**Purpose:** Handles custom alterations to the database structure.
**Current Status:** No objects to apply in this version.
