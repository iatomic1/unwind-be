import json
import os
from pathlib import Path

def update_api_paths(file_path):
    """
    Read a JSON file, update all paths by adding '/api/v1' prefix,
    and save the modified JSON back to file.
    
    Args:
        file_path (str): Path to the JSON file
    """
    try:
        # Convert to absolute path and resolve any relative path components
        abs_file_path = Path(file_path).resolve()
        print(f"Attempting to read file: {abs_file_path}")
        
        # Read the JSON file
        with open(abs_file_path, 'r') as file:
            data = json.load(file)
        
        # Create a new dictionary with updated paths
        updated_paths = {}
        for path in data['paths']:
            new_path = f"/api/v1{path}"
            updated_paths[new_path] = data['paths'][path]
        
        # Update the original data with new paths
        data['paths'] = updated_paths
        
        # Write the modified JSON back to file
        with open(abs_file_path, 'w') as file:
            json.dump(data, file, indent=2)
            
        print(f"Successfully updated API paths in {abs_file_path}!")
        
    except FileNotFoundError:
        print(f"Error: File '{abs_file_path}' not found")
        print("Current working directory:", os.getcwd())
    except json.JSONDecodeError:
        print("Error: Invalid JSON format in file")
    except KeyError:
        print("Error: JSON file doesn't contain 'paths' key")
    except Exception as e:
        print(f"An unexpected error occurred: {str(e)}")

if __name__ == "__main__":
    # Assuming you're running the script from the 'scripts' directory,
    # navigate up one level and then to the openapi.json file
    file_path = Path(__file__).parent.parent / "internal" / "docs" / "openapi.json"
    update_api_paths(str(file_path))
