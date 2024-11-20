# Step 2 Overview - Database Upgrade Process

## Purpose
Step 2 focuses on creating and updating database program units, including PLSQL objects, triggers, and views. It handles the core database object creation and validation phase of the upgrade process.

## Key Operations

### 1. Environment Setup
- Initializes build environment
- Disables recyclebin for the session
- Sets up logging infrastructure
- Configures time tracking for the process

### 2. Audit and Version Control
- Creates audit table and trigger if requested
- Recreates Get_Oasis_Version function
- Sets up version tracking mechanisms

### 3. Base Object Creation
Executes base master scripts in sequence:
- Types (OASIS_BUILD_BASE_TYP)
- PL/SQL Objects (OASIS_BUILD_BASE_PLS)
- Views (OASIS_BUILD_BASE_VEW)
- Triggers (OASIS_BUILD_BASE_TRG)

### 4. Custom Object Creation
Executes custom master scripts in sequence:
- Types (OASIS_BUILD_CUST_TYP)
- PL/SQL Objects (OASIS_BUILD_CUST_PLS)
- Views (OASIS_BUILD_CUST_VEW)
- Triggers (OASIS_BUILD_CUST_TRG)

### 5. Validation and Compilation
- Performs alien objects check
- Verifies PLSQL objects invoker rights
- Executes final compilation of all objects
- Validates object states

## Important Notes
- Maintains strict object creation order
- Includes comprehensive validation checks
- Separates base and custom object handling
- Provides detailed logging and timing reports

## Next Steps
After successful completion of Step 2:
- All database objects are created
- Validations are completed
- System is ready for Step 3 of the upgrade process