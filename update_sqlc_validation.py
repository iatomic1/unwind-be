import re
import yaml
import glob
import os

def parse_sql_files(migrations_dir):
    tag_pattern = re.compile(r'--\s*tags:`([^`]+)`')
    column_pattern = re.compile(r'CREATE TABLE IF NOT EXISTS (\w+)\s*\((.*?)\);', re.DOTALL)
    column_def_pattern = re.compile(r'(\w+)[^,]*?(--.*?)?(?=,|\))', re.DOTALL)
    
    table_validations = []
    
    # Read all .sql files in migrations directory
    for sql_file in glob.glob(os.path.join(migrations_dir, "*.sql")):
        print(f"Processing file: {sql_file}")
        with open(sql_file, 'r') as f:
            content = f.read()
            
        # Find all CREATE TABLE statements
        for table_match in column_pattern.finditer(content):
            table_name = table_match.group(1)
            table_body = table_match.group(2)
            print(f"\nFound table: {table_name}")
            
            # Find all column definitions
            for column_match in column_def_pattern.finditer(table_body):
                column_name = column_match.group(1)
                comment = column_match.group(2) or ''
                print(f"  Column: {column_name}")
                print(f"  Comment: {comment}")
                
                # Extract validation tags if present
                tag_match = tag_pattern.search(comment)
                if tag_match:
                    validation_rules = tag_match.group(1)
                    print(f"  Found validation: {validation_rules}")
                    table_validations.append({
                        'column': f'{table_name}.{column_name}',
                        'go_struct_tag': f'validate:"{validation_rules}"'
                    })
    
    return table_validations

def update_sqlc_config(config_file, validations):
    # Read existing config
    with open(config_file, 'r') as f:
        config = yaml.safe_load(f)
    
    # Get existing overrides or create new ones
    overrides = config['sql'][0]['gen']['go'].get('overrides', [])
    
    # Filter out existing column validations to avoid duplicates
    overrides = [override for override in overrides 
                if 'column' not in override]
    
    # Add new validations
    overrides.extend(validations)
    
    # Update config
    config['sql'][0]['gen']['go']['overrides'] = overrides
    
    # Write updated config
    with open(config_file, 'w') as f:
        yaml.dump(config, f, sort_keys=False)
    
    # Print the validations for debugging
    print("\nValidations being added:")
    for validation in validations:
        print(f"  {validation['column']}: {validation['go_struct_tag']}")

if __name__ == "__main__":
    migrations_dir = "./internal/db/migrations"
    sqlc_config = "sqlc.yaml"
    
    validations = parse_sql_files(migrations_dir)
    update_sqlc_config(sqlc_config, validations)
    print(f"\nUpdated {sqlc_config} with {len(validations)} validation rules")
