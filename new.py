import yaml
import os

def parse_sql_file(sql_file):
    validations = []
    table_name = None
    
    with open(sql_file, 'r') as f:
        lines = f.readlines()
        
    for line in lines:
        if 'CREATE TABLE IF NOT EXISTS' in line:
            table_name = line.split('CREATE TABLE IF NOT EXISTS')[1].strip().split()[0]
            continue
            
        if table_name and '-- tags:' in line:
            # Split the line into column definition and tags
            parts = line.split('-- tags:')
            # Get column name (first word in the column definition)
            column_name = parts[0].strip().split()[0]
            # Extract validation rules
            validation = parts[1].strip().strip('`')
            
            validations.append({
                'column': f'{table_name}.{column_name}',
                'go_struct_tag': validation
            })
            
    return validations

def update_sqlc_config():
    migrations_dir = "./internal/db/migrations"
    sqlc_config = "sqlc.yaml"
    
    all_validations = []
    for file in os.listdir(migrations_dir):
        if file.endswith('.sql'):
            file_path = os.path.join(migrations_dir, file)
            all_validations.extend(parse_sql_file(file_path))
    
    # Update sqlc.yaml
    with open(sqlc_config, 'r') as f:
        config = yaml.safe_load(f)
    
    # Get existing overrides or create new ones
    overrides = config['sql'][0]['gen']['go'].get('overrides', [])
    
    # Filter out existing column validations
    overrides = [o for o in overrides if 'column' not in o]
    overrides.extend(all_validations)
    
    config['sql'][0]['gen']['go']['overrides'] = overrides
    
    with open(sqlc_config, 'w') as f:
        yaml.dump(config, f, sort_keys=False)
    
    print(f"Added {len(all_validations)} validations:")
    for v in all_validations:
        print(f"  {v['column']}: {v['go_struct_tag']}")

if __name__ == "__main__":
    update_sqlc_config()
