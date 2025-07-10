import yaml
from pathlib import Path

def load_config(filename="config/config.yaml"):
    """
    Load YAML config file relative to project root.
    """
    # Resolve project root based on this file's location
    project_root = Path(__file__).resolve().parent.parent

    config_path = project_root / filename

    with open(config_path, "r") as f:
        return yaml.safe_load(f)
