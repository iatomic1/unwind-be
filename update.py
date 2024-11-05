import re
import yaml
import os
from pathlib import Path

def extract_validation_tags(sql_file):
    with open(sql_file, 'r') as f:
        content = f.read()
    
    # Regular expression to match column definitions with comments
    pattern = r'(\w+)\s+\w+(?:$$\d+$$)?\s+(?:NOT NULL\s+)?.*?--\s*tags:(.*?)(?:\n|$)'
    matches = re.findall(pattern, content, re.IGNORECASE | re.MULTILINE)
    
    validation_tags = {}
    for column, tags in matches:
        validation_tags[column] = tags.strip()
    
    return validation_tags

def update_sqlc_config(config_file, validation_tags):
    with open(config_file, 'r') as f:
        config = yaml.safe_load(f)
    
    # Ensure the necessary structure exists
    if 'sql' not in config:
        config['sql'] = [{'gen': {'go': {'overrides': []}}}]
    elif not isinstance(config['sql'], list):
        config['sql'] = [config['sql']]
    
    if 'gen' not in config['sql'][0]:
        config['sql'][0]['gen'] = {'go': {'overrides': []}}
    elif 'go' not in config['sql'][0]['gen']:
        config['sql'][0]['gen']['go'] = {'overrides': []}
    elif 'overrides' not in config['sql'][0]['gen']['go']:
        config['sql'][0]['gen']['go']['overrides'] = []
    
    overrides = config['sql'][0]['gen']['go']['overrides']
    
    # Update or add validation tags
    for column, tags in validation_tags.items():
        override = next((o for o in overrides if o.get('column') == column), None)
        if override:
            override['go_struct_tag'] = f'validate:"{tags}"'
        else:
            overrides.append({
                'column': column,
                'go_struct_tag': f'validate:"{tags}"'
            })
    
    # Write updated config back to file
    with open(config_file, 'w') as f:
        yaml.dump(config, f)

def main():
    migrations_dir = './internal/db/migrations'
    config_file = 'sqlc.yaml'
    
    all_validation_tags = {}
    
    # Process all SQL files in the migrations directory
    for sql_file in Path(migrations_dir).glob('*.sql'):
        table_name = sql_file.stem  # Use filename as table name
        file_tags = extract_validation_tags(sql_file)
        for column, tags in file_tags.items():
            all_validation_tags[f"{table_name}.{column}"] = tags
    
    update_sqlc_config(config_file, all_validation_tags)
    print("SQLC configuration updated successfully.")

if __name__ == "__main__":
    main()
