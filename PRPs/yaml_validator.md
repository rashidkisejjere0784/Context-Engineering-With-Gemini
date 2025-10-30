# Product Requirements Prompt (PRP)
## 1. Overview
- **Feature Name:** YAML Schema Validator

- **Objective:** Build a Rust command-line tool that validates YAML data files against provided schema definitions.

- **Why:** To ensure data integrity and consistency across various YAML configuration files, providing developers with a robust validation tool for early error detection.

## 2. Success Criteria
- [ ] The code compiles and runs without errors.
- [ ] All unit tests pass.
- [ ] The tool correctly parses both the data and schema YAML files.
- [ ] The tool accurately validates data against the schema.
- [ ] Clear and informative error messages are displayed for invalid data.
- [ ] The tool adheres to Rust best practices and code style guidelines.
- [ ] The tool uses `clap`, `serde`, `serde_yaml`, and `jsonschema` correctly.

## 3. Context & Resources
### External Documentation:
- **Resource:** <https://crates.io/crates/clap>
   - **Purpose:**  To learn how to use the `clap` crate to create a user-friendly CLI for the tool.  Specifically, how to define arguments, options, and help messages.
- **Resource:** <https://crates.io/crates/serde>
   - **Purpose:** To understand how to use `serde` and `serde_yaml` for YAML parsing and serialization in Rust. This includes defining structs that map to the YAML structure.
- **Resource:** <https://crates.io/crates/serde_yaml>
   - **Purpose:**  To understand the specifics of using `serde_yaml` for YAML serialization and deserialization.
- **Resource:** <https://crates.io/crates/jsonschema>
   - **Purpose:** To learn how to use `jsonschema` to perform schema validation on the parsed YAML data, and how to interpret its error messages effectively.


### Internal Codebase Patterns:
_None. This is the first Rust feature in this project._


### Known Pitfalls:
- **Error Handling:** The `jsonschema` crate's error messages might be complex.  The tool must translate these into user-friendly messages.
- **Rust Ownership:** Pay close attention to Rust's borrowing and ownership rules when working with parsed YAML data to avoid compiler errors.
- **YAML Schema Complexity:** The schema itself could be complex; ensure the validator correctly handles nested structures, arrays, and different data types.

## 4. Implementation Blueprint
### Proposed File Structure:
```
src/
+-- main.rs          (new)
+-- validator.rs     (new)
+-- error_handler.rs (new)
tests/
+-- validator_test.rs (new)
```

### Task Breakdown:
**Task 1: Project Setup**

- Create a new Rust project using Cargo.
- Add necessary dependencies (`clap`, `serde`, `serde_yaml`, `jsonschema`) to `Cargo.toml`.

**Task 2: CLI Implementation (using `clap`)**

- Define command-line arguments for schema and data file paths using `clap`.
- Handle argument parsing and potential errors.

**Task 3: YAML Parsing (using `serde` and `serde_yaml`)**

- Define Rust structs representing the schema and data structures.
- Implement functions to parse YAML files into these structs using `serde_yaml`.
- Handle parsing errors gracefully.

**Task 4: Schema Validation (using `jsonschema`)**

- Implement a function to validate the parsed data against the parsed schema using `jsonschema`.

**Task 5: Error Handling and Reporting**

- Create a module (`error_handler.rs`) to translate `jsonschema` errors into user-friendly messages.
- Implement robust error handling throughout the application.

**Task 6: Main Function**

- Orchestrate the workflow: parse files, validate data, and report results via the CLI.


## 5. Validation Plan
### Unit Tests:
- `test_valid_data():`  Verify successful validation with valid data and schema.
- `test_missing_key():` Check error reporting when a required key is missing.
- `test_invalid_type():` Verify error reporting when a key's data type is incorrect.
- `test_invalid_yaml():`  Handle cases with malformed YAML input files (both data and schema).
- `test_missing_file():`  Handle cases where the data or schema file is missing.

**Manual Test Command:**
```bash
cargo run --release -- <path_to_data.yaml> <path_to_schema.yaml>
```

**Expected Output (Success):**
```
Data is valid.
```

**Expected Output (Failure):**
```
Error: Key 'version' is missing.
```
